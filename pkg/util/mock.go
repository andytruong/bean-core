package util

import (
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
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

	if !tx.HasTable(Migration{}) {
		if err := tx.CreateTable(Migration{}).Error; nil != err {
			tx.Rollback()
			panic(err)
		}
	}

	if err := module.Install(tx, "sqlite3"); nil != err {
		tx.Rollback()
		panic(err)
	}

	tx.Commit()
}
