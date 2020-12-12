package user

import (
	"path"
	"runtime"
	
	"go.uber.org/zap"
	"gorm.io/gorm"
	
	"bean/components/module"
	"bean/components/module/migrate"
	"bean/components/unique"
	"bean/pkg/util"
)

func NewUserBundle(db *gorm.DB, logger *zap.Logger, id *unique.Identifier) *UserBundle {
	if err := util.NilPointerErrorValidate(db, logger, id); nil != err {
		panic(err)
	}
	
	this := &UserBundle{
		logger:                   logger,
		db:                       db,
		id:                       id,
		maxSecondaryEmailPerUser: 20,
	}
	
	this.Resolvers = newResolvers(this)
	this.Service = &UserService{bundle: this}
	this.NameService = &NameService{bean: this}
	this.EmailService = &EmailService{bean: this}
	this.PasswordService = &PasswordService{bean: this}
	
	return this
}

type UserBundle struct {
	logger                   *zap.Logger
	db                       *gorm.DB
	id                       *unique.Identifier
	maxSecondaryEmailPerUser uint8
	Resolvers                *Resolvers
	Service                  *UserService
	NameService              *NameService
	EmailService             *EmailService
	PasswordService          *PasswordService
}

func (this UserBundle) Dependencies() []module.Bundle {
	return nil
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
