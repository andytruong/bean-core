package user

import (
	"path"
	"runtime"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"bean/components/module"
	"bean/components/module/migrate"
	"bean/components/unique"
	util2 "bean/components/util"
)

func NewUserBundle(db *gorm.DB, logger *zap.Logger, id *unique.Identifier) *UserBundle {
	if err := util2.NilPointerErrorValidate(db, logger, id); nil != err {
		panic(err)
	}

	this := &UserBundle{
		logger:                   logger,
		db:                       db,
		id:                       id,
		maxSecondaryEmailPerUser: 20,
	}

	this.Service = &UserService{bundle: this}
	this.nameService = &NameService{bundle: this}
	this.emailService = &EmailService{bundle: this}
	this.passwordService = &PasswordService{bundle: this}
	this.resolvers = newResolvers(this)

	return this
}

type UserBundle struct {
	module.AbstractBundle

	Service *UserService

	// Internal services
	logger                   *zap.Logger
	db                       *gorm.DB
	id                       *unique.Identifier
	maxSecondaryEmailPerUser uint8
	resolvers                map[string]interface{}
	nameService              *NameService
	emailService             *EmailService
	passwordService          *PasswordService
}

func (this UserBundle) Migrate(tx *gorm.DB, driver string) error {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return nil
	}

	runner := migrate.Runner{
		Tx:     tx,
		Logger: this.logger,
		Driver: driver,
		Bean:   "user",
		Dir:    path.Dir(filename) + "/model/migration/",
	}

	return runner.Run()
}

func (this *UserBundle) GraphqlResolver() map[string]interface{} {
	return this.resolvers
}
