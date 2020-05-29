package util

import (
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"

	"bean/pkg/util/migrate"
)

func MockDatabase() *gorm.DB {
	db, err := gorm.Open("sqlite3", ":memory:")
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

func MockInstall(module Module, db *gorm.DB) {
	tx := db.Begin()

	err := mockInstall(module, tx)
	if nil != err {
		tx.Rollback()
	} else {
		tx.Commit()
	}
}

func mockInstall(module Module, tx *gorm.DB) error {
	if !tx.HasTable(migrate.Migration{}) {
		if err := tx.CreateTable(migrate.Migration{}).Error; nil != err {
			return err
		}
	}

	if err := module.Migrate(tx, "sqlite3"); nil != err {
		tx.Rollback()
		panic(err)
	} else {
		dependencies := module.Dependencies()
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
