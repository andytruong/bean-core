package connect

import (
	"database/sql"

	"github.com/mattn/go-sqlite3"
)

func Driver(db *sql.DB) string {
	switch db.Driver().(type) {
	case *sqlite3.SQLiteDriver:
		return SQLite

	default:
		return Postgres
	}
}
