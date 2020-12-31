package connect

import (
	"gorm.io/gorm"

	"bean/components/module"
	"bean/components/module/migrate"
)

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

	if err := bean.Migrate(tx, SQLite); nil != err {
		tx.Rollback()
		panic(err)
	}

	return nil
}
