package space

import (
	"path"
	"runtime"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"bean/components/module"
	"bean/components/module/migrate"
	"bean/components/scalar"
	"bean/pkg/user"
)

func NewSpaceBundle(
	db *gorm.DB, logger *zap.Logger, idr *scalar.Identifier,
	userBundle *user.UserBundle,
	config *SpaceConfiguration,
) *SpaceBundle {
	this := &SpaceBundle{
		logger:     logger,
		db:         db,
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

type SpaceBundle struct {
	module.AbstractBundle

	Service       *SpaceService
	MemberService *MemberService

	// Internal services
	config            *SpaceConfiguration
	logger            *zap.Logger
	db                *gorm.DB
	idr               *scalar.Identifier
	userBundle        *user.UserBundle
	resolvers         map[string]interface{}
	configService     *ConfigService
	domainNameService *DomainNameService
}

func (bundle *SpaceBundle) Dependencies() []module.Bundle {
	return []module.Bundle{bundle.userBundle}
}

func (bundle SpaceBundle) Migrate(tx *gorm.DB, driver string) error {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return nil
	}

	runner := migrate.Runner{
		Tx:     tx,
		Logger: bundle.logger,
		Driver: driver,
		Bundle: "space",
		Dir:    path.Dir(filename) + "/model/migration/",
	}

	return runner.Run()
}

func (bundle *SpaceBundle) GraphqlResolver() map[string]interface{} {
	return bundle.resolvers
}
