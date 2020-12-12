package module

import (
	"gorm.io/gorm"
)

type Bundle interface {
	Migrate(tx *gorm.DB, driver string) error
	Dependencies() []Bundle
}
