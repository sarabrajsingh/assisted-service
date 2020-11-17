package migrations

import (
	"errors"

	"github.com/google/uuid"
	"github.com/openshift/assisted-service/internal/common"
	"github.com/openshift/assisted-service/internal/events"
	"github.com/openshift/assisted-service/models"

	gormigrate "github.com/go-gormigrate/gormigrate/v2"
	"github.com/go-openapi/strfmt"
	"gorm.io/gorm"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("changeOverridesToText", func() {
	var (
		db        *gorm.DB
		gm        *gormigrate.Gormigrate
		overrides string
	)

	BeforeEach(func() {
		db = common.PrepareTestDB("change_overrides_to_text", &events.Event{})

		overrides = `{"ignition": {"version": "3.1.0"}, "storage": {"files": [{"path": "/tmp/example", "contents": {"source": "data:text/plain;base64,aGVscGltdHJhcHBlZGluYXN3YWdnZXJzcGVj"}}]}}`

		gm = gormigrate.New(db, gormigrate.DefaultOptions, all())
		err := gm.MigrateTo("20201019194303")
		Expect(err).ToNot(HaveOccurred())
	})

	It("Migrates down and up", func() {
		t, err := columnType(db)
		Expect(err).ToNot(HaveOccurred())
		Expect(t).To(Equal("text"))
		applyAndExpectOverride(db, overrides)

		err = gm.RollbackMigration(changeOverridesToText())
		Expect(err).ToNot(HaveOccurred())

		t, err = columnType(db)
		Expect(err).ToNot(HaveOccurred())
		Expect(t).To(Equal("varchar(2048)"))
		applyAndExpectOverride(db, overrides)

		err = gm.MigrateTo("20201019194303")
		Expect(err).ToNot(HaveOccurred())

		t, err = columnType(db)
		Expect(err).ToNot(HaveOccurred())
		Expect(t).To(Equal("text"))
		applyAndExpectOverride(db, overrides)
	})
})

func columnType(db *gorm.DB) (string, error) {
	rows, err := db.Model(&common.Cluster{}).Rows()
	Expect(err).NotTo(HaveOccurred())

	colTypes, err := rows.ColumnTypes()
	Expect(err).NotTo(HaveOccurred())

	for _, colType := range colTypes {
		if colType.Name() == "install_config_overrides" {
			return colType.DatabaseTypeName(), nil
		}
	}
	return "", errors.New("Failed to find install_config_overrides column in clusters")
}

func applyAndExpectOverride(db *gorm.DB, overrides string) {

	// create a new UUID
	clusterID := strfmt.UUID(uuid.New().String())

	// create a Model with overrides
	cluster := common.Cluster{Cluster: models.Cluster{
		ID:                      &clusterID,
		IgnitionConfigOverrides: overrides,
	}}

	// apply override against db
	err := db.Create(&cluster).Error
	Expect(err).NotTo(HaveOccurred())

	var c common.Cluster

	// search for override
	err = db.First(&c, "id = ?", clusterID).Error
	Expect(err).ShouldNot(HaveOccurred())
	Expect(c.IgnitionConfigOverrides).To(Equal(overrides))
}
