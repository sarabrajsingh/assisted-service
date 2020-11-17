package host

import (
	"context"

	"github.com/openshift/assisted-service/internal/common"

	"github.com/go-openapi/strfmt"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/openshift/assisted-service/models"
	"gorm.io/gorm"
)

var _ = Describe("stop-podman", func() {
	ctx := context.Background()
	var host models.Host
	var db *gorm.DB
	var stopCmd *stopInstallationCmd
	var id, clusterId strfmt.UUID
	var stepReply []*models.Step
	var stepErr error
	dbName := "stop_podman"

	BeforeEach(func() {
		db = common.PrepareTestDB(dbName)
		stopCmd = NewStopInstallationCmd(getTestLog())

		id = strfmt.UUID(uuid.New().String())
		clusterId = strfmt.UUID(uuid.New().String())
		host = getTestHost(id, clusterId, models.HostStatusError)
		Expect(db.Create(&host).Error).ShouldNot(HaveOccurred())
	})

	It("get_step", func() {
		stepReply, stepErr = stopCmd.GetSteps(ctx, &host)
		Expect(stepReply[0].StepType).To(Equal(models.StepTypeExecute))
		Expect(stepErr).ShouldNot(HaveOccurred())
	})

	AfterEach(func() {
		// cleanup

		stepReply = nil
		stepErr = nil
	})
})
