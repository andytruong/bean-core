package infra

import (
	"context"

	"bean/components/connect"
	"bean/components/module/migrate"
)

func (c *Container) Migrate(ctx context.Context) error {
	if db, err := c.dbs.Master(); nil != err {
		return err
	} else if con, err := db.DB(); nil != err {
		return err
	} else {
		// start transaction
		driver := connect.Driver(con)
		tx := db.WithContext(ctx).Begin()

		// create migration table if not existing
		if !tx.Migrator().HasTable(migrate.Migration{}) {
			if err := tx.Migrator().CreateTable(migrate.Migration{}); nil != err {
				tx.Rollback()
				return err
			}
		}

		// loop through bundles
		bundles := c.BundleList()
		for _, bundle := range bundles.Get() {
			if err := bundle.Migrate(tx, driver); nil != err {
				tx.Rollback()
				return err
			}
		}

		return tx.Commit().Error
	}
}
