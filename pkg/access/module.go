package access

import (
	"github.com/jinzhu/gorm"
)

func NewAccessModule() *AccessModule {
	return &AccessModule{
		Mutation: &AccessMutationResolver{},
	}
}

type AccessModule struct {
	Mutation *AccessMutationResolver
}

func (this AccessModule) Migrate(tx *gorm.DB, driver string) error {
	return nil
}
