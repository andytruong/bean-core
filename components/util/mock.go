package util

import (
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"bean/components/module"
	"bean/components/module/migrate"
	"bean/components/scalar"
	connect2 "bean/components/util/connect"
)

func MockDatabase() *gorm.DB {
	con := sqlite.Open(":memory:")
	db, err := gorm.Open(con, &gorm.Config{})
	if nil != err {
		panic(err)
	} else {
		db.Logger = db.Logger.LogMode(logger.Silent)
	}

	return db
}

func MockLogger() *zap.Logger {
	return zap.NewNop()
}

func MockIdentifier() *scalar.Identifier {
	return &scalar.Identifier{}
}

func MockInstall(bean module.Bundle, db *gorm.DB) {
	tx := db.Begin()

	err := mockInstall(bean, tx)
	if nil != err {
		tx.Rollback()
	} else {
		tx.Commit()
	}
}

func mockInstall(bean module.Bundle, tx *gorm.DB) error {
	if !tx.Migrator().HasTable(migrate.Migration{}) {
		if err := tx.Migrator().CreateTable(migrate.Migration{}); nil != err {
			return err
		}
	}

	dependencies := bean.Dependencies()
	if nil != dependencies {
		for _, dependency := range dependencies {
			err := mockInstall(dependency, tx)
			if nil != err {
				return err
			}
		}
	}

	if err := bean.Migrate(tx, connect2.SQLite); nil != err {
		tx.Rollback()
		panic(err)
	}

	return nil
}
