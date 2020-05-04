package infra

import (
	"context"

	"bean/pkg/access"
	"bean/pkg/access/model"
	"bean/pkg/infra/gql"
	namespace_model "bean/pkg/namespace/model"
	"bean/pkg/user"
	user_model "bean/pkg/user/model"
)

type (
	rootResolver struct {
		container *Container
	}

	mutationResolver struct {
		*access.AccessMutationResolver
		*user.UserMutationResolver
	}

	queryResolver struct {
		*access.AccessQueryResolver
	}

	sessionResolver struct {
		container *Container
	}
)

func (this *rootResolver) Session() gql.SessionResolver {
	return this.container.gql.session
}

func (this *rootResolver) Mutation() gql.MutationResolver {
	return this.container.gql.mutation
}

func (this *rootResolver) Query() gql.QueryResolver {
	return this.container.gql.query
}

func (this mutationResolver) Ping(ctx context.Context) (string, error) {
	return "pong", nil
}

func (this queryResolver) Ping(ctx context.Context) (string, error) {
	return "pong", nil
}

func (this sessionResolver) User(ctx context.Context, obj *model.Session) (*user_model.User, error) {
	panic("implement me")
}

func (this sessionResolver) Namespace(ctx context.Context, obj *model.Session) (*namespace_model.Namespace, error) {
	panic("implement me")
}
