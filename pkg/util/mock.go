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

	if !tx.HasTable(migrate.Migration{}) {
		if err := tx.CreateTable(migrate.Migration{}).Error; nil != err {
			tx.Rollback()
			panic(err)
		}
	}

	if err := module.Migrate(tx, "sqlite3"); nil != err {
		tx.Rollback()
		panic(err)
	}

	tx.Commit()
}
