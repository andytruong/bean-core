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
	spaceBundle *space.Bundle,
	configBundle *config.Bundle,
) (*Bundle, error) {
	bundle := &Bundle{
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

type Bundle struct {
	module.AbstractBundle

	Service      *AppService
	spaceBundle  *space.Bundle
	configBundle *config.Bundle
	idr          *scalar.Identifier
	lgr          *zap.Logger
	hook         *module.Hook
	resolvers    map[string]interface{}
}

func (Bundle) Name() string {
	return "App"
}

func (bundle Bundle) Dependencies() []module.Bundle {
	if nil != bundle.spaceBundle {
		return []module.Bundle{bundle.spaceBundle}
	}

	return nil
}

func (bundle Bundle) Migrate(ctx context.Context, driver string) error {
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

func (bundle *Bundle) GraphqlResolver() map[string]interface{} {
	return bundle.resolvers
}
