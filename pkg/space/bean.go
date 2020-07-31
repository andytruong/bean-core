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

func NewSpaceBean(
	db *gorm.DB, logger *zap.Logger, id *unique.Identifier,
	bUser *user.UserBean,
	genetics *Genetic,
) *SpaceBean {
	this := &SpaceBean{
		logger:  logger,
		db:      db,
		id:      id,
		user:    bUser,
		genetic: genetics,
	}

	this.Resolvers = newResolver(this)
	this.Core = &Core{bean: this}
	this.CoreConfig = &CoreConfig{bean: this}
	this.CoreDomainName = &CoreDomainName{bean: this}
	this.CoreMember = &CoreMember{
		bean:     this,
		Resolver: newMembershipResolver(this, bUser),
	}

	return this
}

type SpaceBean struct {
	genetic *Genetic
	logger  *zap.Logger
	db      *gorm.DB
	id      *unique.Identifier
	user    *user.UserBean
	config  *config.ConfigBean

	Resolvers      *Resolvers
	Core           *Core
	CoreConfig     *CoreConfig
	CoreMember     *CoreMember
	CoreDomainName *CoreDomainName
}

func (this *SpaceBean) MembershipResolver() MembershipResolver {
	return this.CoreMember.Resolver
}

func (this *SpaceBean) Dependencies() []module.Bean {
	return []module.Bean{this.user}
}

func (this SpaceBean) Migrate(tx *gorm.DB, driver string) error {
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

func (this SpaceBean) Space(ctx context.Context, filters dto.SpaceFilters) (*model.Space, error) {
	return this.Core.Find(ctx, filters)
}

func (this SpaceBean) Load(ctx context.Context, id string) (*model.Space, error) {
	return this.Core.Load(ctx, id)
}

func (this SpaceBean) SpaceCreate(ctx context.Context, in dto.SpaceCreateInput) (*dto.SpaceCreateOutcome, error) {
	txn := this.db.WithContext(ctx).Begin()
	out, err := this.Core.Create(txn, in)

	if nil != err {
		txn.Rollback()
		return nil, err
	} else {
		return out, txn.Commit().Error
	}
}

func (this SpaceBean) SpaceUpdate(ctx context.Context, in dto.SpaceUpdateInput) (*dto.SpaceCreateOutcome, error) {
	space, err := this.Load(ctx, in.SpaceID)
	if nil != err {
		return nil, err
	}

	txn := this.db.WithContext(ctx).Begin()
	out, err := this.Core.Update(txn, space, in)

	if nil != err {
		txn.Rollback()
		return nil, err
	} else {
		return out, txn.Commit().Error
	}
}

func (this SpaceBean) SpaceMembershipCreate(
	ctx context.Context,
	input dto.SpaceMembershipCreateInput,
) (*dto.SpaceMembershipCreateOutcome, error) {
	space, err := this.Load(ctx, input.SpaceID)
	if nil != err {
		return nil, err
	}

	user, err := this.user.Resolvers.Query.User(ctx, input.UserID)
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
	outcome, err := this.CoreMember.Create(tx, input, space, user)

	if nil != err {
		tx.Rollback()
		return nil, err
	} else {
		return outcome, tx.Commit().Error
	}
}

func (this SpaceBean) SpaceMembershipUpdate(ctx context.Context, input dto.SpaceMembershipUpdateInput) (*dto.SpaceMembershipCreateOutcome, error) {
	membership, err := this.Resolvers.Query.Membership(ctx, input.Id, scalar.NilString(input.Version))
	if nil != err {
		return nil, err
	}

	tx := this.db.WithContext(ctx).Begin()
	outcome, err := this.CoreMember.Update(tx, input, membership)

	if nil != err {
		tx.Rollback()
		return nil, err
	} else {
		return outcome, tx.Commit().Error
	}
}
