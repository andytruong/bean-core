package util

import (
	"github.com/jinzhu/gorm"
)

type (
	Module interface {
		Migrate(tx *gorm.DB, driver string) error
		Dependencies() []Module
	}
)
