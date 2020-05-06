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
	db       *gorm.DB
	logger   *zap.Logger
	id       *util.Identifier
	mutation *UserMutationResolver
	query    *UserQueryResolver
}

func (this *UserModule) MutationResolver() *UserMutationResolver {
	if nil == this.mutation {
		this.mutation = &UserMutationResolver{
			db: this.db,
			id: this.id,
		}
	}

	return this.mutation
}

func (this *UserModule) QueryResolver() *UserQueryResolver {
	if nil == this.mutation {
		this.query = &UserQueryResolver{db: this.db}
	}

	return this.query
}
