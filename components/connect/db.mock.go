package connect

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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
