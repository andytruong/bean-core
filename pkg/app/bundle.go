package app

import (
	"path"
	"runtime"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"bean/components/module"
	"bean/components/module/migrate"
	"bean/components/scalar"
	"bean/pkg/config"
	"bean/pkg/space"
)

func NewApplicationBundle(
	idr *scalar.Identifier,
	logger *zap.Logger,
	spaceBundle *space.SpaceBundle,
	configBundle *config.ConfigBundle,
) (*AppBundle, error) {
	bundle := &AppBundle{
		idr:          idr,
		logger:       logger,
		spaceBundle:  spaceBundle,
		configBundle: configBundle,
	}

	bundle.resolvers = bundle.newResolvers()
	bundle.Service = &AppService{bundle: bundle}

	return bundle, nil
}

type AppBundle struct {
	module.AbstractBundle

	Service      *AppService
	spaceBundle  *space.SpaceBundle
	configBundle *config.ConfigBundle
	idr          *scalar.Identifier
	logger       *zap.Logger
	resolvers    map[string]interface{}
}

func (AppBundle) Name() string {
	return "App"
}

func (bundle AppBundle) Dependencies() []module.Bundle {
	if nil != bundle.spaceBundle {
		return []module.Bundle{bundle.spaceBundle}
	}

	return nil
}

func (bundle AppBundle) Migrate(tx *gorm.DB, driver string) error {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return nil
	}

	runner := migrate.Runner{
		Tx:     tx,
		Logger: bundle.logger,
		Driver: driver,
		Bundle: "bundle",
		Dir:    path.Dir(filename) + "/model/migration/",
	}

	return runner.Run()
}

func (bundle *AppBundle) GraphqlResolver() map[string]interface{} {
	return bundle.resolvers
}
