package infra

import (
	"context"

	"bean/pkg/access/model"
	namespace_model "bean/pkg/namespace/model"
	user_model "bean/pkg/user/model"
)

func (this sessionResolver) User(ctx context.Context, obj *model.Session) (*user_model.User, error) {
	panic("implement me")
}

func (this sessionResolver) Namespace(ctx context.Context, obj *model.Session) (*namespace_model.Namespace, error) {
	panic("implement me")
}
