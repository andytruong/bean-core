package infra

import (
	"context"
	"database/sql"

	"github.com/mattn/go-sqlite3"

	"bean/components/connect"
)

func dbDriver(db *sql.DB) string {
	switch db.Driver().(type) {
	case *sqlite3.SQLiteDriver:
		return connect.SQLite

	default:
		return connect.Postgres
	}
}

func (c *Container) Migrate(ctx context.Context) error {
	db, err := c.dbs.Master()
	if nil != err {
		return err
	}

	if con, err := db.DB(); nil != err {
		return err
	} else {
		// start transaction
		driver := dbDriver(con)
		tx := db.WithContext(ctx).Begin()
		ctx = connect.DBToContext(ctx, db)

		// create migration table if not existing
		if !tx.Migrator().HasTable(connect.Migration{}) {
			if err := tx.Migrator().CreateTable(connect.Migration{}); nil != err {
				tx.Rollback()
				return err
			}
		}

		// loop through bundles
		bundles := c.BundleList()
		for _, bundle := range bundles.Get() {
			if err := bundle.Migrate(ctx, driver); nil != err {
				tx.Rollback()
				return err
			}
		}

		return tx.Commit().Error
	}
}
