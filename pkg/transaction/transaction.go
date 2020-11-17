package transaction

import "gorm.io/gorm"

func AddForUpdateQueryOption(db *gorm.DB) *gorm.DB {
	if db.Dialector.Name() != "sqlite3" {
		// return a new object and not overwrite pointer value because GORM have a pointer to parent
		return db.Set("gorm:query_option", "FOR UPDATE")
	}
	return db
}
