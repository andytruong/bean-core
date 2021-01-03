package access

import (
	"context"
	"path"
	"runtime"

	"go.uber.org/zap"

	"bean/components/connect"
	"bean/components/module"
	"bean/components/scalar"
	"bean/pkg/space"
	"bean/pkg/user"
)

func NewAccessBundle(
	idr *scalar.Identifier,
	lgr *zap.Logger,
	userBundle *user.Bundle,
	spaceBundle *space.Bundle,
	cnf *Config,
) *Bundle {
	this := &Bundle{
		cnf:         cnf.init(),
		lgr:         lgr,
		idr:         idr,
		userBundle:  userBundle,
		spaceBundle: spaceBundle,
	}

	this.sessionService = &SessionService{bundle: this}
	this.JwtService = &JwtService{bundle: this}
	this.resolvers = this.newResolves()

	return this
}

type (
	Bundle struct {
		module.AbstractBundle

		cnf            *Config
		lgr            *zap.Logger
		idr            *scalar.Identifier
		sessionService *SessionService
		JwtService     *JwtService

		// depends on userBundle bundle
		userBundle  *user.Bundle
		spaceBundle *space.Bundle
		resolvers   map[string]interface{}
	}
)

func (Bundle) Name() string {
	return "Access"
}

func (bundle Bundle) Dependencies() []module.Bundle {
	return []module.Bundle{
		bundle.userBundle,
		bundle.spaceBundle,
	}
}

func (bundle Bundle) Migrate(ctx context.Context, driver string) error {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return nil
	}

	runner := connect.Runner{
		Logger: bundle.lgr,
		Driver: driver,
		Bundle: "access",
		Dir:    path.Dir(filename) + "/model/migration/",
	}

	return runner.Run(ctx)
}

func (bundle Bundle) GraphqlResolver() map[string]interface{} {
	return bundle.resolvers
}
