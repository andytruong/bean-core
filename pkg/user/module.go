package user

import (
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"

	"bean/pkg/util"
)

func NewUserModule(db *gorm.DB, logger *zap.Logger, id *util.Identifier) (*UserModule, error) {
	if err := util.NilPointerErrorValidate(db, logger, id); nil != err {
		return nil, err
	}

	module := &UserModule{logger: logger}
	module.Mutation = UserMutationResolver{db: db, id: id}
	module.Query = UserQueryResolver{db: db}

	return module, nil
}

type UserModule struct {
	logger   *zap.Logger
	Mutation UserMutationResolver
	Query    UserQueryResolver
}
