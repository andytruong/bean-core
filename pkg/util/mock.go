package util

import (
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"

	"bean/pkg/util/connect"
	"bean/pkg/util/migrate"
)

func MockDatabase() *gorm.DB {
	db, err := gorm.Open(connect.SQLite, ":memory:")
	if nil != err {
		panic(err)
	}

	return db
}

func MockLogger() *zap.Logger {
	return zap.NewNop()
}

func MockIdentifier() *Identifier {
	return &Identifier{}
}

func MockInstall(bean Bean, db *gorm.DB) {
	tx := db.Begin()

	err := mockInstall(bean, tx)
	if nil != err {
		tx.Rollback()
	} else {
		tx.Commit()
	}
}

func mockInstall(bean Bean, tx *gorm.DB) error {
	if !tx.HasTable(migrate.Migration{}) {
		if err := tx.CreateTable(migrate.Migration{}).Error; nil != err {
			return err
		}
	}

	if err := bean.Migrate(tx, connect.SQLite); nil != err {
		tx.Rollback()
		panic(err)
	} else {
		dependencies := bean.Dependencies()
		if nil != dependencies {
			for _, dependency := range dependencies {
				err := mockInstall(dependency, tx)
				if nil != err {
					return err
				}
			}
		}
	}

	return nil
}
