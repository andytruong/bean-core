package namespace

import (
	"context"
	"database/sql"
	"path"
	"runtime"

	"github.com/jinzhu/gorm"
	"go.uber.org/zap"

	"bean/pkg/namespace/api/handler"
	"bean/pkg/namespace/model"
	"bean/pkg/namespace/model/dto"
	"bean/pkg/util"
	"bean/pkg/util/migrate"
)

func NewNamespaceModule(db *gorm.DB, logger *zap.Logger, id *util.Identifier) *NamespaceModule {
	return &NamespaceModule{
		logger: logger,
		db:     db,
		id:     id,
	}
}

type NamespaceModule struct {
	logger *zap.Logger
	db     *gorm.DB
	id     *util.Identifier
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

	return hdl.NamespaceUpdate(tx, namespace, input)
}
