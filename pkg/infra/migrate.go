package infra

import (
	"context"
	"database/sql"

	"bean/pkg/util/connect"
	"bean/pkg/util/migrate"
)

func (this *Can) Migrate(ctx context.Context) error {
	if db, err := this.dbs.master(); nil != err {
		return err
	} else {
		// start transaction
		driver := connect.Driver(db)
		tx := db.BeginTx(ctx, &sql.TxOptions{})

		// create migration table if not existing
		if !tx.HasTable(migrate.Migration{}) {
			if err := tx.CreateTable(migrate.Migration{}).Error; nil != err {
				tx.Rollback()
				return err
			}
		}

		// loop through modules
		for _, module := range this.modules.List() {
			if err := module.Migrate(tx, driver); nil != err {
				tx.Rollback()
				return err
			}
		}

		return tx.Commit().Error
	}
}
