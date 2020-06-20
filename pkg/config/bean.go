package config

import (
	"path"
	"runtime"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"bean/pkg/util"
	"bean/pkg/util/migrate"
)

func NewConfigBean(id *util.Identifier, logger *zap.Logger) *ConfigBean {
	this := &ConfigBean{
		id:     id,
		logger: logger,
	}

	this.CoreBucket = &CoreBucket{bean: this}
	this.CoreVariable = &CoreVariable{bean: this}

	return this
}

type ConfigBean struct {
	id           *util.Identifier
	logger       *zap.Logger
	CoreBucket   *CoreBucket
	CoreVariable *CoreVariable
}

func (this ConfigBean) Migrate(tx *gorm.DB, driver string) error {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return nil
	}

	runner := migrate.Runner{
		Tx:     tx,
		Logger: this.logger,
		Driver: driver,
		Bean:   "config",
		Dir:    path.Dir(filename) + "/model/migration/",
	}

	return runner.Run()
}

func (this ConfigBean) Dependencies() []util.Bean {
	return nil
}
