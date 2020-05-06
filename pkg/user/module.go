package user

import (
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"

	"bean/pkg/util"
)

func NewUserService(db *gorm.DB, logger *zap.Logger, id *util.Identifier) (*UserModule, error) {
	if err := util.NilPointerErrorValidate(db, logger, id); nil != err {
		return nil, err
	}

	return &UserModule{
		db:     db,
		logger: logger,
		id:     id,
	}, nil
}

type UserModule struct {
	db       *gorm.DB
	logger   *zap.Logger
	id       *util.Identifier
	mutation *UserMutationResolver
	query    *UserQueryResolver
}

func (this *UserModule) MutationResolver() (*UserMutationResolver, error) {
	var err error

	if nil == this.mutation {
		this.mutation, err = NewUserMutationResolver(this.db, this.id)
	}

	return this.mutation, err
}

func (this *UserModule) QueryResolver() *UserQueryResolver {
	if nil == this.mutation {
		this.query = &UserQueryResolver{db: this.db}
	}

	return this.query
}
