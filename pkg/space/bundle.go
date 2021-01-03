package space

import (
	"context"
	"path"
	"runtime"

	"go.uber.org/zap"

	"bean/components/connect"
	"bean/components/module"
	"bean/components/scalar"
	"bean/pkg/user"
)

func NewSpaceBundle(
	lgr *zap.Logger, idr *scalar.Identifier,
	userBundle *user.Bundle,
	config *Config,
) *Bundle {
	this := &Bundle{
		lgr:        lgr,
		idr:        idr,
		userBundle: userBundle,
		config:     config,
	}

	this.resolvers = this.newResolvers()
	this.Service = &SpaceService{bundle: this}
	this.configService = &ConfigService{bundle: this}
	this.domainNameService = &DomainNameService{bundle: this}
	this.MemberService = &MemberService{bundle: this}

	return this
}

type Bundle struct {
	module.AbstractBundle

	Service       *SpaceService
	MemberService *MemberService

	// Internal services
	config            *Config
	lgr               *zap.Logger
	idr               *scalar.Identifier
	userBundle        *user.Bundle
	resolvers         map[string]interface{}
	configService     *ConfigService
	domainNameService *DomainNameService
}

func (Bundle) Name() string {
	return "Space"
}

func (bundle *Bundle) Dependencies() []module.Bundle {
	return []module.Bundle{bundle.userBundle}
}

func (bundle Bundle) Migrate(ctx context.Context, driver string) error {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return nil
	}

	runner := connect.Runner{
		Logger: bundle.lgr,
		Driver: driver,
		Bundle: "space",
		Dir:    path.Dir(filename) + "/model/migration/",
	}

	return runner.Run(ctx)
}

func (bundle *Bundle) GraphqlResolver() map[string]interface{} {
	return bundle.resolvers
}
