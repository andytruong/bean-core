package space

import (
	"context"
	"path"
	"runtime"
	
	"go.uber.org/zap"
	"gorm.io/gorm"
	
	"bean/components/module"
	"bean/components/module/migrate"
	"bean/components/unique"
	"bean/pkg/config"
	"bean/pkg/space/model"
	"bean/pkg/space/model/dto"
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
	
	this.Service = &SpaceService{bundle: this}
	this.ConfigService = &ConfigService{bundle: this}
	this.DomainNameService = &DomainNameService{bean: this}
	this.MemberService = &MemberService{
		bundle:   this,
		Resolver: newMembershipResolver(this, userBundle),
	}
	
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
	
	Service           *SpaceService
	ConfigService     *ConfigService
	MemberService     *MemberService
	DomainNameService *DomainNameService
}

func (this *SpaceBundle) MembershipResolver() MembershipResolver {
	return this.MemberService.Resolver
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

func (this SpaceBundle) Space(ctx context.Context, filters dto.SpaceFilters) (*model.Space, error) {
	return this.Service.Find(ctx, filters)
}

func (this SpaceBundle) Load(ctx context.Context, id string) (*model.Space, error) {
	return this.Service.Load(ctx, id)
}
