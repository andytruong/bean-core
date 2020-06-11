package namespace

import (
	"context"
	"database/sql"
	"path"
	"runtime"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"bean/pkg/namespace/api/handler"
	"bean/pkg/namespace/model"
	"bean/pkg/namespace/model/dto"
	"bean/pkg/user"
	"bean/pkg/util"
	"bean/pkg/util/connect"
	"bean/pkg/util/migrate"
)

func NewNamespaceBean(
	db *gorm.DB, logger *zap.Logger, id *util.Identifier,
	bUser *user.UserBean,
	config *Config,
) *NamespaceBean {
	this := &NamespaceBean{
		logger: logger,
		db:     db,
		id:     id,
		user:   bUser,
		config: config,
	}

	this.Core = &NamespaceBeanCore{bean: this}
	this.Member = &NamespaceBeanMembership{
		bean:     this,
		Resolver: newMembershipResolver(this, bUser),
	}

	return this
}

type NamespaceBean struct {
	logger *zap.Logger
	db     *gorm.DB
	id     *util.Identifier
	user   *user.UserBean
	config *Config

	Core   *NamespaceBeanCore
	Member *NamespaceBeanMembership
}

func (this *NamespaceBean) MembershipResolver() MembershipResolver {
	return this.Member.Resolver
}

func (this *NamespaceBean) Dependencies() []util.Bean {
	return []util.Bean{this.user}
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
	hdl := handler.NamespaceQueryHandler{DB: this.db}

	return hdl.Handle(ctx, filters)
}

func (this NamespaceBean) Load(ctx context.Context, id string) (*model.Namespace, error) {
	hdl := handler.NamespaceQueryHandler{DB: this.db}

	return hdl.Load(ctx, id)
}

func (this NamespaceBean) NamespaceCreate(ctx context.Context, input dto.NamespaceCreateInput) (*dto.NamespaceCreateOutcome, error) {
	txn := this.db.BeginTx(ctx, &sql.TxOptions{})
	outcome, err := this.Core.Create(txn, input)

	if nil != err {
		txn.Rollback()
		return nil, err
	} else {
		return outcome, txn.Commit().Error
	}
}

func (this NamespaceBean) DomainNames(ctx context.Context, namespace *model.Namespace) (*model.DomainNames, error) {
	hdl := handler.DomainQueryHandler{DB: this.db}
	return hdl.DomainNames(ctx, namespace)
}

func (this NamespaceBean) Features(ctx context.Context, namespace *model.Namespace) (*model.NamespaceFeatures, error) {
	hdl := handler.NamespaceQueryFeaturesHandler{DB: this.db}
	return hdl.Features(ctx, namespace)
}

func (this NamespaceBean) NamespaceUpdate(ctx context.Context, input dto.NamespaceUpdateInput) (*bool, error) {
	namespace, err := this.Load(ctx, input.NamespaceID)
	if nil != err {
		return nil, err
	}

	tx := this.db.BeginTx(ctx, &sql.TxOptions{})
	hdl := handler.NamespaceUpdateHandler{ID: this.id}

	outcome, err := hdl.NamespaceUpdate(tx, namespace, input)

	if nil != err {
		tx.Rollback()
		return nil, err
	} else {
		return outcome, tx.Commit().Error
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

	user, err := this.user.User(ctx, input.UserID)
	if nil != err {
		return nil, err
	}

	features, err := this.Features(ctx, namespace)
	if nil != err {
		return nil, err
	}

	if !features.Register {
		return nil, errors.Wrap(util.ErrorConfig, "register is off")
	}

	tx := this.db.BeginTx(ctx, &sql.TxOptions{})
	outcome, err := this.Member.Create(tx, input, namespace, user)

	if nil != err {
		tx.Rollback()
		return nil, err
	} else {
		return outcome, tx.Commit().Error
	}
}

func (this NamespaceBean) NamespaceMembershipUpdate(ctx context.Context, input dto.NamespaceMembershipUpdateInput) (*dto.NamespaceMembershipCreateOutcome, error) {
	membership, err := this.Membership(ctx, input.Id, util.NilString(input.Version))
	if nil != err {
		return nil, err
	}

	tx := this.db.BeginTx(ctx, &sql.TxOptions{})
	outcome, err := this.Member.Update(tx, input, membership)

	if nil != err {
		tx.Rollback()
		return nil, err
	} else {
		return outcome, tx.Commit().Error
	}
}

func (this NamespaceBean) Membership(ctx context.Context, id string, version *string) (*model.Membership, error) {
	obj := &model.Membership{}

	err := this.db.
		Table(connect.TableNamespaceMemberships).
		First(&obj, "id = ?", id).
		Error

	if nil != err {
		return nil, err
	} else if nil != version {
		if obj.Version != *version {
			return nil, util.ErrorVersionConflict
		}
	}

	return obj, nil
}

func (this NamespaceBean) Memberships(ctx context.Context, first int, after *string, filters dto.MembershipsFilter) (*model.MembershipConnection, error) {
	return this.Member.Find(first, after, filters)
}

func (this NamespaceBean) Parent(ctx context.Context, obj *model.Namespace) (*model.Namespace, error) {
	if nil == obj.ParentID {
		return nil, nil
	}

	return this.Load(ctx, *obj.ParentID)
}
