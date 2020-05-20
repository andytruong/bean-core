package access

import (
	"context"
	"database/sql"

	"github.com/jinzhu/gorm"
	"go.uber.org/zap"

	"bean/pkg/access/api/handler"
	"bean/pkg/access/model"
	"bean/pkg/access/model/dto"
	namespace_model "bean/pkg/namespace/model"
	user_model "bean/pkg/user/model"
	"bean/pkg/util"
)

func NewAccessModule() *AccessModule {
	return &AccessModule{
	}
}

type AccessModule struct {
	logger          *zap.Logger
	db              *gorm.DB
	id              *util.Identifier
	SessionResolver ModuleResolver
}

func (this AccessModule) Migrate(tx *gorm.DB, driver string) error {
	return nil
}

func (this *AccessModule) SessionCreate(ctx context.Context, input *dto.LoginInput) (*dto.LoginOutcome, error) {
	hdl := handler.SessionCreateHandler{ID: this.id}
	txn := this.db.BeginTx(ctx, &sql.TxOptions{})
	outcome, err := hdl.SessionCreate(ctx, txn, input)
	if nil != err {
		txn.Rollback()

		return nil, err
	}

	return outcome, txn.Commit().Error
}

func (this *AccessModule) SessionDelete(ctx context.Context, input *dto.LoginInput) (*dto.LogoutOutcome, error) {
	panic("not implemented")
}

func (this AccessModule) LoadSession(ctx context.Context, input *dto.ValidationInput) (*dto.ValidationOutcome, error) {
	panic("wip")
}

type ModuleResolver struct {
}

func (this ModuleResolver) User(ctx context.Context, obj *model.Session) (*user_model.User, error) {
	panic("implement me")
}

func (this ModuleResolver) Namespace(ctx context.Context, obj *model.Session) (*namespace_model.Namespace, error) {
	panic("implement me")
}
