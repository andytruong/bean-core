package util

import (
	"gorm.io/gorm"
)

type Bean interface {
	Migrate(tx *gorm.DB, driver string) error
	Dependencies() []Bean
}
