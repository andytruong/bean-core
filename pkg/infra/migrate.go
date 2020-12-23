package infra

import (
	"context"

	"bean/components/module/migrate"
	connect2 "bean/components/util/connect"
)

func (container *Container) Migrate(ctx context.Context) error {
	if db, err := container.dbs.master(); nil != err {
		return err
	} else if con, err := db.DB(); nil != err {
		return err
	} else {
		// start transaction
		driver := connect2.Driver(con)
		tx := db.WithContext(ctx).Begin()

		// create migration table if not existing
		if !tx.Migrator().HasTable(migrate.Migration{}) {
			if err := tx.Migrator().CreateTable(migrate.Migration{}); nil != err {
				tx.Rollback()
				return err
			}
		}

		// loop through bundles
		for _, bundle := range container.BundleList() {
			if err := bundle.Migrate(tx, driver); nil != err {
				tx.Rollback()
				return err
			}
		}

		return tx.Commit().Error
	}
}
