package user

import (
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"

	"bean/pkg/util"
)

func NewUserService(db *gorm.DB, logger *zap.Logger, id *util.Identifier) *UserModule {
	return &UserModule{
		db:     db,
		logger: logger,
		id:     id,
	}
}

type UserModule struct {
	db        *gorm.DB
	logger    *zap.Logger
	id        *util.Identifier
	rMutation *UserMutationResolver
}

func (this *UserModule) MutationResolver() *UserMutationResolver {
	if nil == this.rMutation {
		this.rMutation = &UserMutationResolver{
			db: this.db,
			id: this.id,
		}
	}

	return this.rMutation
}
