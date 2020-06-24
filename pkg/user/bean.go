package user

import (
	"path"
	"runtime"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"bean/pkg/util"
	"bean/pkg/util/migrate"
)

func NewUserBean(db *gorm.DB, logger *zap.Logger, id *util.Identifier) *UserBean {
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
	id                       *util.Identifier
	maxSecondaryEmailPerUser uint8
	Resolvers                *Resolvers
	Core                     *Core
	CoreName                 *CoreName
	CoreEmail                *CoreEmail
	CorePassword             *CorePassword
}

func (this UserBean) Dependencies() []util.Bean {
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
