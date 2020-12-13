package space

import (
	"context"
	"path"
	"runtime"
	
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	
	"bean/components/module"
	"bean/components/module/migrate"
	"bean/components/scalar"
	"bean/components/unique"
	"bean/pkg/config"
	"bean/pkg/space/model"
	"bean/pkg/space/model/dto"
	"bean/pkg/user"
	"bean/pkg/util"
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
	
	this.Resolvers = newResolver(this)
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
	
	Resolvers         *Resolvers
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

func (this SpaceBundle) SpaceCreate(ctx context.Context, in dto.SpaceCreateInput) (*dto.SpaceCreateOutcome, error) {
	txn := this.db.WithContext(ctx).Begin()
	out, err := this.Service.Create(txn, in)
	
	if nil != err {
		txn.Rollback()
		return nil, err
	} else {
		return out, txn.Commit().Error
	}
}

func (this SpaceBundle) SpaceUpdate(ctx context.Context, in dto.SpaceUpdateInput) (*dto.SpaceCreateOutcome, error) {
	space, err := this.Load(ctx, in.SpaceID)
	if nil != err {
		return nil, err
	}
	
	txn := this.db.WithContext(ctx).Begin()
	out, err := this.Service.Update(txn, space, in)
	
	if nil != err {
		txn.Rollback()
		return nil, err
	} else {
		return out, txn.Commit().Error
	}
}

func (this SpaceBundle) SpaceMembershipCreate(
	ctx context.Context,
	input dto.SpaceMembershipCreateInput,
) (*dto.SpaceMembershipCreateOutcome, error) {
	space, err := this.Load(ctx, input.SpaceID)
	if nil != err {
		return nil, err
	}
	
	user, err := this.userBundle.Resolvers.Query.User(ctx, input.UserID)
	if nil != err {
		return nil, err
	}
	
	features, err := this.Resolvers.Object.Features(ctx, space)
	if nil != err {
		return nil, err
	}
	
	if !features.Register {
		return nil, errors.Wrap(util.ErrorConfig, "register is off")
	}
	
	tx := this.db.WithContext(ctx).Begin()
	outcome, err := this.MemberService.Create(tx, input, space, user)
	
	if nil != err {
		tx.Rollback()
		return nil, err
	} else {
		return outcome, tx.Commit().Error
	}
}

func (this SpaceBundle) SpaceMembershipUpdate(ctx context.Context, input dto.SpaceMembershipUpdateInput) (*dto.SpaceMembershipCreateOutcome, error) {
	membership, err := this.Resolvers.Query.Membership(ctx, input.Id, scalar.NilString(input.Version))
	if nil != err {
		return nil, err
	}
	
	tx := this.db.WithContext(ctx).Begin()
	outcome, err := this.MemberService.Update(tx, input, membership)
	
	if nil != err {
		tx.Rollback()
		return nil, err
	} else {
		return outcome, tx.Commit().Error
	}
}
