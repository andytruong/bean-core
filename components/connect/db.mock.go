package connect

import (
	"context"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"bean/components/module"
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

func MockInstall(ctx context.Context, bean module.Bundle) {
	db := DB(ctx)
	tx := db.Begin()
	err := mockInstall(DBToContext(ctx, tx), bean)
	if nil != err {
		tx.Rollback()
		panic(err)
	} else {
		tx.Commit()
	}
}

func mockInstall(ctx context.Context, bundle module.Bundle) error {
	db := DB(ctx)
	if !db.Migrator().HasTable(Migration{}) {
		if err := db.Migrator().CreateTable(Migration{}); nil != err {
			return err
		}
	}

	dependencies := bundle.Dependencies()
	if nil != dependencies {
		for _, dependency := range dependencies {
			err := mockInstall(ctx, dependency)
			if nil != err {
				return err
			}
		}
	}

	if err := bundle.Migrate(ctx, SQLite); nil != err {
		db.Rollback()
		panic(err)
	}

	return nil
}
