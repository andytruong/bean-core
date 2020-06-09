package util

import (
	"github.com/jinzhu/gorm"
)

type Bean interface {
	Migrate(tx *gorm.DB, driver string) error
	Dependencies() []Bean
}
