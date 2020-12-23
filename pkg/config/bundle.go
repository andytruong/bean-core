package config

import (
	"path"
	"runtime"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"bean/components/module"
	"bean/components/module/migrate"
	"bean/components/unique"
)

func NewConfigBundle(id *unique.Identifier, logger *zap.Logger) *ConfigBundle {
	this := &ConfigBundle{
		id:     id,
		logger: logger,
	}

	this.BucketService = &BucketService{bundle: this}
	this.VariableService = &VariableService{bundle: this}

	return this
}

type ConfigBundle struct {
	module.AbstractBundle

	id              *unique.Identifier
	logger          *zap.Logger
	BucketService   *BucketService
	VariableService *VariableService
}

func (bundle ConfigBundle) Migrate(tx *gorm.DB, driver string) error {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return nil
	}

	runner := migrate.Runner{
		Tx:     tx,
		Logger: bundle.logger,
		Driver: driver,
		Bundle: "config",
		Dir:    path.Dir(filename) + "/model/migration/",
	}

	return runner.Run()
}

func (bundle ConfigBundle) Dependencies() []module.Bundle {
	return nil
}
