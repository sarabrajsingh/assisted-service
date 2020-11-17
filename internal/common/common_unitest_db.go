package common

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	. "github.com/onsi/gomega"
	"github.com/openshift/assisted-service/models"
)

func PrepareTestDB(dbName string, extrasSchemas ...interface{}) *gorm.DB {

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})

	Expect(err).ShouldNot(HaveOccurred())
	// db = db.Debug()
	db.AutoMigrate(&models.Host{}, &Cluster{})

	if len(extrasSchemas) > 0 {
		for _, schema := range extrasSchemas {
			db.AutoMigrate(schema)
			Expect(db.Error).ShouldNot(HaveOccurred())
		}
	}
	return db
}
