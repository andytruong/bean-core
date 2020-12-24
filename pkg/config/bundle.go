package config

import (
	"path"
	"runtime"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"bean/components/module"
	"bean/components/module/migrate"
	"bean/components/scalar"
)

func NewConfigBundle(idr *scalar.Identifier, logger *zap.Logger) *ConfigBundle {
	bundle := &ConfigBundle{
		idr:    idr,
		logger: logger,
	}

	bundle.BucketService = &BucketService{bundle: bundle}
	bundle.VariableService = &VariableService{bundle: bundle}

	return bundle
}

type ConfigBundle struct {
	module.AbstractBundle

	idr             *scalar.Identifier
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
