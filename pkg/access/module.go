package access

import (
	"github.com/jinzhu/gorm"
)

func NewAccessModule() *AccessModule {
	return &AccessModule{}
}

type AccessModule struct {
}

func (this AccessModule) MutationResolver() (*AccessMutationResolver, error) {
	return &AccessMutationResolver{}, nil
}

func (this AccessModule) Install(tx *gorm.DB) error {
	panic("implement me")
}
