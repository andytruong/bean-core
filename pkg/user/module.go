package user

import (
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"

	"bean/pkg/user/service"
	"bean/pkg/util"
)

func NewUserModule(db *gorm.DB, logger *zap.Logger, id *util.Identifier) (*UserModule, error) {
	if err := util.NilPointerErrorValidate(db, logger, id); nil != err {
		return nil, err
	} else {
		var err error
		module := &UserModule{logger: logger}
		module.Mutation, err = NewUserMutationResolver(db, id)
		if nil != err {
			return nil, err
		}

		module.Query = UserQueryResolver{db: db}
		module.Model = UserModelResolver{db: db}
		module.Email = UserEmailResolver{db: db}

		return module, nil
	}
}

type UserModule struct {
	logger   *zap.Logger
	Mutation *UserMutationResolver
	Query    UserQueryResolver
	Model    UserModelResolver
	Email    UserEmailResolver
}

func (this UserModule) Install(tx *gorm.DB, driver string) error {
	api := service.NewUserInstallAPI(this.logger)

	return api.Run(tx, driver)
}
