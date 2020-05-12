package infra

import (
	"context"
	"database/sql"

	"bean/pkg/util"
)

func (this *Container) Install(ctx context.Context) error {
	if db, err := this.dbs.master(); nil != err {
		return err
	} else {
		// start transaction
		tx := db.BeginTx(ctx, &sql.TxOptions{})

		// create migration table if not existing
		if !tx.HasTable(util.Migration{}) {
			if err := tx.CreateTable(util.Migration{}).Error; nil != err {
				tx.Rollback()
				return err
			}
		}

		// loop through modules
		for _, mod := range this.modules.List() {
			if err := mod.Install(tx); nil != err {
				tx.Rollback()
				return err
			}
		}

		return tx.Commit().Error
	}
}
