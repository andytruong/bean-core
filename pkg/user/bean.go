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

func NewUserBean(db *gorm.DB, logger *zap.Logger, id *unique.Identifier) *UserBean {
	if err := util.NilPointerErrorValidate(db, logger, id); nil != err {
		panic(err)
	}

	this := &UserBean{
		logger:                   logger,
		db:                       db,
		id:                       id,
		maxSecondaryEmailPerUser: 20,
	}

	this.Resolvers = newResolvers(this)
	this.Core = &Core{bean: this}
	this.CoreName = &CoreName{bean: this}
	this.CoreEmail = &CoreEmail{bean: this}
	this.CorePassword = &CorePassword{bean: this}

	return this
}

type UserBean struct {
	logger                   *zap.Logger
	db                       *gorm.DB
	id                       *unique.Identifier
	maxSecondaryEmailPerUser uint8
	Resolvers                *Resolvers
	Core                     *Core
	CoreName                 *CoreName
	CoreEmail                *CoreEmail
	CorePassword             *CorePassword
}

func (this UserBean) Dependencies() []module.Bean {
	return nil
}

func (this UserBean) Migrate(tx *gorm.DB, driver string) error {
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
