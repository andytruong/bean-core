package config

import (
	"path"
	"runtime"

	"github.com/jinzhu/gorm"
	"go.uber.org/zap"

	"bean/pkg/util"
	"bean/pkg/util/migrate"
)

func NewConfigBean(id *util.Identifier, logger *zap.Logger) *ConfigBean {
	this := &ConfigBean{
		id:     id,
		logger: logger,
	}

	this.CoreBucket = &ConfigBucketBean{bean: this}
	this.CoreVariable = &ConfigVariableBean{bean: this}

	return this
}

type ConfigBean struct {
	id           *util.Identifier
	logger       *zap.Logger
	CoreBucket   *ConfigBucketBean
	CoreVariable *ConfigVariableBean
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
