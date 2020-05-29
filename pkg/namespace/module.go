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
	"bean/pkg/util/migrate"
)

func NewNamespaceModule(
	db *gorm.DB, logger *zap.Logger, id *util.Identifier,
	userModule *user.UserModule,
) *NamespaceModule {
	this := &NamespaceModule{
		logger:     logger,
		db:         db,
		id:         id,
		userModule: userModule,
	}

	this.membership = newMembershipResolver(this, userModule)

	return this
}

type NamespaceModule struct {
	logger     *zap.Logger
	db         *gorm.DB
	id         *util.Identifier
	userModule *user.UserModule
	membership MembershipResolver
}

func (this *NamespaceModule) MembershipResolver() MembershipResolver {
	return this.membership
}

func (this *NamespaceModule) Dependencies() []util.Module {
	return []util.Module{this.userModule}
}

func (this NamespaceModule) Migrate(tx *gorm.DB, driver string) error {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return nil
	}

	runner := migrate.Runner{
		Tx:     tx,
		Logger: this.logger,
		Driver: driver,
		Module: "namespace",
		Dir:    path.Dir(filename) + "/model/migration/",
	}

	return runner.Run()
}

func (this NamespaceModule) Namespace(ctx context.Context, id string) (*model.Namespace, error) {
	obj := &model.Namespace{}
	err := this.db.First(&obj, "id = ?", id).Error
	if nil != err {
		return nil, err
	}

	return obj, nil
}

func (this NamespaceModule) NamespaceCreate(ctx context.Context, input dto.NamespaceCreateInput) (*dto.NamespaceCreateOutcome, error) {
	hdl := handler.NamespaceCreateHandler{ID: this.id}
	txn := this.db.BeginTx(ctx, &sql.TxOptions{})
	outcome, err := hdl.Create(txn, input)

	if nil != err {
		txn.Rollback()
		return nil, err
	} else {
		return outcome, txn.Commit().Error
	}
}

func (this NamespaceModule) DomainNames(ctx context.Context, namespace *model.Namespace) (*model.DomainNames, error) {
	hdl := handler.DomainQueryHandler{DB: this.db}
	return hdl.DomainNames(ctx, namespace)
}

func (this NamespaceModule) Features(ctx context.Context, namespace *model.Namespace) (*model.NamespaceFeatures, error) {
	hdl := handler.NamespaceQueryFeaturesHandler{DB: this.db}
	return hdl.Features(ctx, namespace)
}

func (this NamespaceModule) NamespaceUpdate(ctx context.Context, input dto.NamespaceUpdateInput) (*bool, error) {
	namespace, err := this.Namespace(ctx, input.NamespaceID)
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

func (this NamespaceModule) NamespaceMembershipCreate(
	ctx context.Context,
	input dto.NamespaceMembershipCreateInput,
) (*dto.NamespaceMembershipCreateOutcome, error) {
	namespace, err := this.Namespace(ctx, input.NamespaceID)
	if nil != err {
		return nil, err
	}

	user, err := this.userModule.User(ctx, input.UserID)
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

	hdl := handler.MembershipCreateHandler{
		ID: this.id,
		DB: this.db,
	}

	return hdl.NamespaceMembershipCreate(ctx, input, namespace, user)
}

func (this NamespaceModule) Membership(ctx context.Context, id string, version *string) (*model.Membership, error) {
	obj := &model.Membership{}

	err := this.db.
		Table("namespace_memberships").
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

func (this NamespaceModule) NamespaceMembershipUpdate(ctx context.Context, input dto.NamespaceMembershipUpdateInput) (*dto.NamespaceMembershipCreateOutcome, error) {
	membership, err := this.Membership(ctx, input.Id, util.NilString(input.Version))
	if nil != err {
		return nil, err
	}

	hdl := handler.MembershipUpdateHandler{
		ID: this.id,
		DB: this.db,
	}

	return hdl.NamespaceMembershipUpdate(ctx, input, membership)
}
