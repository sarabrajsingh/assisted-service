package host

import (
	"context"
	"fmt"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/kelseyhightower/envconfig"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/openshift/assisted-service/internal/common"
	"github.com/openshift/assisted-service/internal/events"
	"github.com/openshift/assisted-service/models"
	"github.com/openshift/assisted-service/pkg/leader"
	"gorm.io/gorm"
)

var _ = Describe("monitor_disconnection", func() {
	var (
		ctx        = context.Background()
		db         *gorm.DB
		state      API
		host       models.Host
		ctrl       *gomock.Controller
		mockEvents *events.MockHandler
		dbName     = "monitor_disconnection"
	)

	BeforeEach(func() {
		db = common.PrepareTestDB(dbName, &events.Event{})
		ctrl = gomock.NewController(GinkgoT())
		mockEvents = events.NewMockHandler(ctrl)
		dummy := &leader.DummyElector{}
		state = NewManager(getTestLog(), db, mockEvents, nil, nil, createValidatorCfg(),
			nil, defaultConfig, dummy)
		clusterID := strfmt.UUID(uuid.New().String())
		host = getTestHost(strfmt.UUID(uuid.New().String()), clusterID, models.HostStatusDiscovering)
		cluster := getTestCluster(clusterID, "1.1.0.0/16")
		Expect(db.Save(&cluster).Error).ToNot(HaveOccurred())
		host.Inventory = workerInventory()
		err := state.RegisterHost(ctx, &host)
		Expect(err).ShouldNot(HaveOccurred())
		db.First(&host, "id = ? and cluster_id = ?", host.ID, host.ClusterID)
	})

	Context("host_disconnecting", func() {
		It("known_host_disconnects", func() {
			host.CheckedInAt = strfmt.DateTime(time.Now().Add(-4 * time.Minute))
			host.Status = swag.String(models.HostStatusKnown)
			db.Save(&host)
		})

		It("discovering_host_disconnects", func() {
			host.CheckedInAt = strfmt.DateTime(time.Now().Add(-4 * time.Minute))
			host.Status = swag.String(models.HostStatusDiscovering)
			db.Save(&host)
		})

		It("known_host_insufficient", func() {
			host.CheckedInAt = strfmt.DateTime(time.Now().Add(-4 * time.Minute))
			host.Status = swag.String(models.HostStatusInsufficient)
			db.Save(&host)
		})

		AfterEach(func() {
			mockEvents.EXPECT().AddEvent(gomock.Any(), host.ClusterID, host.ID, models.EventSeverityWarning,
				fmt.Sprintf("Host %s: updated status from \"%s\" to \"disconnected\" (Host has stopped communicating with the installation service)",
					host.ID.String(), *host.Status),
				gomock.Any())
			state.HostMonitoring()
			db.First(&host, "id = ? and cluster_id = ?", host.ID, host.ClusterID)
			Expect(*host.Status).Should(Equal(models.HostStatusDisconnected))
		})
	})

	Context("host_reconnecting", func() {
		It("host_connects", func() {
			host.CheckedInAt = strfmt.DateTime(time.Now())
			host.Inventory = ""
			host.Status = swag.String(models.HostStatusDisconnected)
			db.Save(&host)
		})

		AfterEach(func() {
			mockEvents.EXPECT().AddEvent(gomock.Any(), host.ClusterID, host.ID, models.EventSeverityInfo,
				fmt.Sprintf("Host %s: updated status from \"disconnected\" to \"discovering\" (Waiting for host to send hardware details)", host.ID.String()),
				gomock.Any())
			state.HostMonitoring()
			db.First(&host, "id = ? and cluster_id = ?", host.ID, host.ClusterID)
			Expect(*host.Status).Should(Equal(models.HostStatusDiscovering))
		})
	})

	AfterEach(func() {
		ctrl.Finish()

		sqliteDB, err := db.DB()
		Expect(err).ShouldNot(HaveOccurred())

		err = sqliteDB.Close()
		Expect(err).ShouldNot(HaveOccurred())
	})
})

var _ = Describe("TestHostMonitoring", func() {
	var (
		ctx        = context.Background()
		db         *gorm.DB
		state      API
		host       models.Host
		ctrl       *gomock.Controller
		cfg        Config
		mockEvents *events.MockHandler
		dbName     = "host_monitor_tests"
		clusterID  = strfmt.UUID(uuid.New().String())
	)

	BeforeEach(func() {
		db = common.PrepareTestDB(dbName, &events.Event{})
		ctrl = gomock.NewController(GinkgoT())
		mockEvents = events.NewMockHandler(ctrl)
		mockEvents.EXPECT().
			AddEvent(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			AnyTimes()
		Expect(envconfig.Process("myapp", &cfg)).ShouldNot(HaveOccurred())
		state = NewManager(getTestLog(), db, mockEvents, nil, nil, createValidatorCfg(),
			nil, &cfg, &leader.DummyElector{})
	})

	AfterEach(func() {

		ctrl.Finish()
	})

	Context("validate monitor in batches", func() {
		registerAndValidateDisconnected := func(nHosts int) {
			for i := 0; i < nHosts; i++ {
				if i%10 == 0 {
					clusterID = strfmt.UUID(uuid.New().String())
					cluster := getTestCluster(clusterID, "1.1.0.0/16")
					Expect(db.Save(&cluster).Error).ToNot(HaveOccurred())
				}
				host = getTestHost(strfmt.UUID(uuid.New().String()), clusterID, models.HostStatusDiscovering)
				host.Inventory = workerInventory()
				Expect(state.RegisterHost(ctx, &host)).ShouldNot(HaveOccurred())
				host.CheckedInAt = strfmt.DateTime(time.Now().Add(-4 * time.Minute))
				db.Save(&host)
			}
			state.HostMonitoring()
			var count int64
			Expect(db.Model(&models.Host{}).Where("status = ?", models.HostStatusDisconnected).Count(&count).Error).
				ShouldNot(HaveOccurred())
			Expect(count).Should(Equal(nHosts))
		}

		It("5 hosts all disconnected", func() {
			registerAndValidateDisconnected(5)
		})

		It("15 hosts all disconnected", func() {
			registerAndValidateDisconnected(15)
		})

		It("765 hosts all disconnected", func() {
			registerAndValidateDisconnected(765)
		})
	})

})
