package app

import (
	"context"
	"path"
	"runtime"
	
	"go.uber.org/zap"
	
	"bean/components/connect"
	"bean/components/module"
	"bean/components/scalar"
	"bean/components/util"
	"bean/pkg/config"
	"bean/pkg/space"
)

func NewApplicationBundle(
	idr *scalar.Identifier,
	lgr *zap.Logger,
	hook *module.Hook,
	spaceBundle *space.SpaceBundle,
	configBundle *config.ConfigBundle,
) (*AppBundle, error) {
	bundle := &AppBundle{
		idr:          idr,
		lgr:          lgr,
		hook:         hook,
		spaceBundle:  spaceBundle,
		configBundle: configBundle,
	}
	
	bundle.resolvers = bundle.newResolvers()
	bundle.Service = &AppService{bundle: bundle}
	
	return bundle, nil
}

const (
	ErrorInactiveApplication = util.Err("inactive application")
)

type AppBundle struct {
	module.AbstractBundle
	
	Service      *AppService
	spaceBundle  *space.SpaceBundle
	configBundle *config.ConfigBundle
	idr          *scalar.Identifier
	lgr          *zap.Logger
	hook         *module.Hook
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

func (bundle AppBundle) Migrate(ctx context.Context, driver string) error {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return nil
	}
	
	runner := connect.Runner{
		Logger: bundle.lgr,
		Driver: driver,
		Bundle: "bundle",
		Dir:    path.Dir(filename) + "/model/migration/",
	}
	
	return runner.Run(ctx)
}

func (bundle *AppBundle) GraphqlResolver() map[string]interface{} {
	return bundle.resolvers
}
