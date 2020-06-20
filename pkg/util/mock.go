package util

import (
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"bean/pkg/util/connect"
	"bean/pkg/util/migrate"
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
	if !tx.Migrator().HasTable(migrate.Migration{}) {
		if err := tx.Migrator().CreateTable(migrate.Migration{}); nil != err {
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
