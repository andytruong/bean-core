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

func NewConfigBundle(idr *scalar.Identifier, lgr *zap.Logger) *ConfigBundle {
	bundle := &ConfigBundle{idr: idr, lgr: lgr}
	bundle.BucketService = &BucketService{bundle: bundle}
	bundle.VariableService = &VariableService{bundle: bundle}

	return bundle
}

type ConfigBundle struct {
	module.AbstractBundle

	idr             *scalar.Identifier
	lgr             *zap.Logger
	BucketService   *BucketService
	VariableService *VariableService
}

func (ConfigBundle) Name() string {
	return "Config"
}

func (bundle ConfigBundle) Migrate(tx *gorm.DB, driver string) error {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return nil
	}

	runner := migrate.Runner{
		Tx:     tx,
		Logger: bundle.lgr,
		Driver: driver,
		Bundle: "config",
		Dir:    path.Dir(filename) + "/model/migration/",
	}

	return runner.Run()
}
