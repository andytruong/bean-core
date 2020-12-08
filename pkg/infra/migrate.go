package infra

import (
	"context"
	
	"bean/components/module/migrate"
	"bean/pkg/util/connect"
)

func (this *Can) Migrate(ctx context.Context) error {
	if db, err := this.dbs.master(); nil != err {
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
		
		// loop through beans
		for _, bean := range this.beans.List() {
			if err := bean.Migrate(tx, driver); nil != err {
				tx.Rollback()
				return err
			}
		}
		
		return tx.Commit().Error
	}
}
