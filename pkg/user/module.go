package user

import (
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
)

func NewUserService() *UserModule {
	return &UserModule{}
}

type UserModule struct {
	db        *gorm.DB
	logger    *zap.Logger
	rMutation *UserMutationResolver
}

func (this *UserModule) MutationResolver() *UserMutationResolver {
	if nil == this.rMutation {
		this.rMutation = &UserMutationResolver{
		}
	}

	return this.rMutation
}
