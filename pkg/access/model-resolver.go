package access

import (
	"context"

	"bean/pkg/access/model"
	namespace_model "bean/pkg/namespace/model"
	user_model "bean/pkg/user/model"
)

type ModelResolver struct {
	module *AccessModule
}

func (this ModelResolver) User(ctx context.Context, obj *model.Session) (*user_model.User, error) {
	return this.module.userModule.User(ctx, obj.UserId)
}

func (this ModelResolver) Context(ctx context.Context, obj *model.Session) (*model.SessionContext, error) {
	panic("implement me")
}

func (this ModelResolver) Scopes(ctx context.Context, obj *model.Session) ([]*model.AccessScope, error) {
	return obj.Scopes, nil
}

func (this ModelResolver) Namespace(ctx context.Context, obj *model.Session) (*namespace_model.Namespace, error) {
	return this.module.namespaceModule.Namespace(ctx, obj.NamespaceId)
}
