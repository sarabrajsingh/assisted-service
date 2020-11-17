package migrations

import (
	"sort"

	gormigrate "github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	return gormigrate.New(db, gormigrate.DefaultOptions, all()).Migrate()
}

func all() []*gormigrate.Migration {
	allMigrations := []*gormigrate.Migration{
		changeOverridesToText(),
	}

	sort.SliceStable(allMigrations, func(i, j int) bool { return allMigrations[i].ID < allMigrations[j].ID })

	return allMigrations
}
