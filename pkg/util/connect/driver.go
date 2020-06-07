package connect

import (
	"github.com/jinzhu/gorm"
	"github.com/mattn/go-sqlite3"
)

func Driver(db *gorm.DB) string {
	switch db.DB().Driver().(type) {
	case *sqlite3.SQLiteDriver:
		return SQLite

	default:
		return Postgres
	}
}
