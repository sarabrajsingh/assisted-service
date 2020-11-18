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
	//db = db.Debug()
	err = db.AutoMigrate(&models.Host{}, &Cluster{})
	Expect(err).ShouldNot(HaveOccurred())

	if len(extrasSchemas) > 0 {
		for _, schema := range extrasSchemas {
			err = db.AutoMigrate(schema)
			Expect(err).ShouldNot(HaveOccurred())
		}
	}
	return db
}

func DeleteTestDB(db *gorm.DB, dbName string) {
	sqlDB, err := db.DB()
	Expect(err).ShouldNot(HaveOccurred())

	err = sqlDB.Close()
	Expect(err).ShouldNot(HaveOccurred())
}
