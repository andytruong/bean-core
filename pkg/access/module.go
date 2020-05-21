package access

import (
	"context"
	"database/sql"
	"path"
	"runtime"

	"github.com/jinzhu/gorm"
	"go.uber.org/zap"

	"bean/pkg/access/api/handler"
	"bean/pkg/access/model"
	"bean/pkg/access/model/dto"
	namespace_model "bean/pkg/namespace/model"
	user_model "bean/pkg/user/model"
	"bean/pkg/util"
	"bean/pkg/util/migrate"
)

func NewAccessModule(db *gorm.DB, id *util.Identifier, logger *zap.Logger) *AccessModule {
	return &AccessModule{
		logger: logger,
		db:     db,
		id:     id,
		SessionResolver: ModuleResolver{
		},
	}
}

type AccessModule struct {
	logger          *zap.Logger
	db              *gorm.DB
	id              *util.Identifier
	SessionResolver ModuleResolver
}

func (this AccessModule) Migrate(tx *gorm.DB, driver string) error {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return nil
	}

	runner := migrate.Runner{
		Tx:     tx,
		Logger: this.logger,
		Driver: driver,
		Module: "access",
		Dir:    path.Dir(filename) + "/model/migration/",
	}

	return runner.Run()
}

func (this *AccessModule) SessionCreate(ctx context.Context, input *dto.SessionCreateInput) (*dto.SessionCreateOutcome, error) {
	hdl := handler.SessionCreateHandler{ID: this.id}
	txn := this.db.BeginTx(ctx, &sql.TxOptions{})
	outcome, err := hdl.Handle(ctx, txn, input)
	if nil != err {
		txn.Rollback()

		return nil, err
	}

	return outcome, txn.Commit().Error
}

func (this *AccessModule) SessionDelete(ctx context.Context, input *dto.SessionCreateInput) (*dto.LogoutOutcome, error) {
	panic("not implemented")
}

func (this AccessModule) Session(ctx context.Context, token string) (*model.Session, error) {
	panic("wip")
}

type ModuleResolver struct {
}

func (this ModuleResolver) User(ctx context.Context, obj *model.Session) (*user_model.User, error) {
	panic("implement me")
}

func (this ModuleResolver) Context(ctx context.Context, obj *model.Session) (*model.SessionContext, error) {
	panic("implement me")
}

func (this ModuleResolver) Scopes(ctx context.Context, obj *model.Session) ([]*model.AccessScope, error) {
	return obj.Scopes, nil
}

func (this ModuleResolver) Namespace(ctx context.Context, obj *model.Session) (*namespace_model.Namespace, error) {
	panic("implement me")
}
