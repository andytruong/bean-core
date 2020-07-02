package namespace

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
	"bean/pkg/namespace/model"
	"bean/pkg/namespace/model/dto"
	"bean/pkg/user"
	"bean/pkg/util"
)

func NewNamespaceBean(
	db *gorm.DB, logger *zap.Logger, id *unique.Identifier,
	bUser *user.UserBean,
	genetics *Genetic,
) *NamespaceBean {
	this := &NamespaceBean{
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

type NamespaceBean struct {
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

func (this *NamespaceBean) MembershipResolver() MembershipResolver {
	return this.CoreMember.Resolver
}

func (this *NamespaceBean) Dependencies() []module.Bean {
	return []module.Bean{this.user}
}

func (this NamespaceBean) Migrate(tx *gorm.DB, driver string) error {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return nil
	}

	runner := migrate.Runner{
		Tx:     tx,
		Logger: this.logger,
		Driver: driver,
		Bean:   "namespace",
		Dir:    path.Dir(filename) + "/model/migration/",
	}

	return runner.Run()
}

func (this NamespaceBean) Namespace(ctx context.Context, filters dto.NamespaceFilters) (*model.Namespace, error) {
	return this.Core.Find(ctx, filters)
}

func (this NamespaceBean) Load(ctx context.Context, id string) (*model.Namespace, error) {
	return this.Core.Load(ctx, id)
}

func (this NamespaceBean) NamespaceCreate(ctx context.Context, in dto.NamespaceCreateInput) (*dto.NamespaceCreateOutcome, error) {
	txn := this.db.WithContext(ctx).Begin()
	out, err := this.Core.Create(txn, in)

	if nil != err {
		txn.Rollback()
		return nil, err
	} else {
		return out, txn.Commit().Error
	}
}

func (this NamespaceBean) NamespaceUpdate(ctx context.Context, in dto.NamespaceUpdateInput) (*dto.NamespaceCreateOutcome, error) {
	namespace, err := this.Load(ctx, in.NamespaceID)
	if nil != err {
		return nil, err
	}

	txn := this.db.WithContext(ctx).Begin()
	out, err := this.Core.Update(txn, namespace, in)

	if nil != err {
		txn.Rollback()
		return nil, err
	} else {
		return out, txn.Commit().Error
	}
}

func (this NamespaceBean) NamespaceMembershipCreate(
	ctx context.Context,
	input dto.NamespaceMembershipCreateInput,
) (*dto.NamespaceMembershipCreateOutcome, error) {
	namespace, err := this.Load(ctx, input.NamespaceID)
	if nil != err {
		return nil, err
	}

	user, err := this.user.Resolvers.Query.User(ctx, input.UserID)
	if nil != err {
		return nil, err
	}

	features, err := this.Resolvers.Object.Features(ctx, namespace)
	if nil != err {
		return nil, err
	}

	if !features.Register {
		return nil, errors.Wrap(util.ErrorConfig, "register is off")
	}

	tx := this.db.WithContext(ctx).Begin()
	outcome, err := this.CoreMember.Create(tx, input, namespace, user)

	if nil != err {
		tx.Rollback()
		return nil, err
	} else {
		return outcome, tx.Commit().Error
	}
}

func (this NamespaceBean) NamespaceMembershipUpdate(ctx context.Context, input dto.NamespaceMembershipUpdateInput) (*dto.NamespaceMembershipCreateOutcome, error) {
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
