package infra

import (
	"context"
	"database/sql"

	"github.com/jinzhu/gorm"
	"github.com/mattn/go-sqlite3"

	"bean/pkg/util/migrate"
)

func (this *Container) Migrate(ctx context.Context) error {
	if db, err := this.dbs.master(); nil != err {
		return err
	} else {
		// start transaction
		driver := driver(db)
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

func driver(db *gorm.DB) string {
	switch db.DB().Driver().(type) {
	case *sqlite3.SQLiteDriver:
		return "sqlite3"

	default:
		return "postgres"
	}
}
