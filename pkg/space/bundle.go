package space

import (
	"path"
	"runtime"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"bean/components/module"
	"bean/components/module/migrate"
	"bean/components/unique"
	"bean/pkg/config"
	"bean/pkg/user"
)

func NewSpaceBundle(
	db *gorm.DB, logger *zap.Logger, id *unique.Identifier,
	userBundle *user.UserBundle,
	config *SpaceConfiguration,
) *SpaceBundle {
	this := &SpaceBundle{
		logger:     logger,
		db:         db,
		id:         id,
		userBundle: userBundle,
		config:     config,
	}

	this.resolvers = this.newResolvers()
	this.Service = &SpaceService{bundle: this}
	this.ConfigService = &ConfigService{bundle: this}
	this.DomainNameService = &DomainNameService{bundle: this}
	this.MemberService = &MemberService{bundle: this}

	return this
}

type SpaceBundle struct {
	module.AbstractBundle

	config       *SpaceConfiguration
	logger       *zap.Logger
	db           *gorm.DB
	id           *unique.Identifier
	userBundle   *user.UserBundle
	configBundle *config.ConfigBundle
	resolvers    map[string]interface{}

	Service           *SpaceService
	ConfigService     *ConfigService
	MemberService     *MemberService
	DomainNameService *DomainNameService
}

func (this *SpaceBundle) Dependencies() []module.Bundle {
	return []module.Bundle{this.userBundle}
}

func (this SpaceBundle) Migrate(tx *gorm.DB, driver string) error {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return nil
	}

	runner := migrate.Runner{
		Tx:     tx,
		Logger: this.logger,
		Driver: driver,
		Bean:   "space",
		Dir:    path.Dir(filename) + "/model/migration/",
	}

	return runner.Run()
}

func (this *SpaceBundle) GraphqlResolver() map[string]interface{} {
	return this.resolvers
}
