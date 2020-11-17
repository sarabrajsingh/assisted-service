package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	. "github.com/onsi/gomega"
	"github.com/openshift/assisted-service/internal/common"
	"github.com/openshift/assisted-service/models"
	"gorm.io/gorm"
)

func changeOverridesToText() *gormigrate.Migration {
	migrate := func(tx *gorm.DB) error {
		err := tx.Migrator().DropTable(&models.Cluster{})
		Expect(err).ShouldNot(HaveOccurred())

		type Cluster struct {
			common.Cluster
			InstallConfigOverrides string `json:"install_config_overrides,omitempty" gorm:"type:text"`
		}
		return tx.AutoMigrate(&Cluster{})
	}

	rollback := func(tx *gorm.DB) error {
		err := tx.Migrator().DropTable(&models.Cluster{})
		Expect(err).ShouldNot(HaveOccurred())

		type Cluster struct {
			common.Cluster
			InstallConfigOverrides string `json:"install_config_overrides,omitempty" gorm:"type:varchar(2048)"`
		}
		return tx.AutoMigrate(&Cluster{})
	}

	return &gormigrate.Migration{
		ID:       "20201019194303",
		Migrate:  gormigrate.MigrateFunc(migrate),
		Rollback: gormigrate.RollbackFunc(rollback),
	}
}
