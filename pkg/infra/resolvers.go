package infra

import (
	"context"
	"sync"

	"bean/pkg/access"
	"bean/pkg/access/model"
	"bean/pkg/infra/gql"
	namespace_model "bean/pkg/namespace/model"
	"bean/pkg/user"
	user_model "bean/pkg/user/model"
)

type (
	// lazy access into detail resolver.
	resolvers struct {
		container *Container
		mu        *sync.Mutex
		root      *rootResolver
		query     *queryResolver
		mutation  *mutationResolver
		session   *sessionResolver
		user      *userResolver
	}

	rootResolver struct {
		container *Container
	}

	mutationResolver struct {
		*user.UserMutationResolver
		*access.AccessMutationResolver
	}

	queryResolver struct {
		*access.AccessQueryResolver
		user.UserQueryResolver
	}

	sessionResolver struct {
		container *Container
	}

	userResolver struct {
		container *Container
	}
)

func (this *rootResolver) UserEmail() gql.UserEmailResolver {
	panic("implement me")
}

func (this *resolvers) getRoot() *rootResolver {
	if nil == this.root {
		this.mu.Lock()
		defer this.mu.Unlock()

		this.root = &rootResolver{container: this.container}
	}

	return this.root
}

func (this *resolvers) getQuery() *queryResolver {
	if this.query == nil {
		this.mu.Lock()
		defer this.mu.Unlock()

		mUser, err := this.container.modules.User()
		if nil != err {
			panic(err)
		}

		this.query = &queryResolver{
			// AccessQueryResolver: this.container.modules.Access().MutationResolver(),
			UserQueryResolver: mUser.Query,
		}
	}

	return this.query
}

func (this *resolvers) getMutation() *mutationResolver {
	if nil == this.mutation {
		this.mu.Lock()
		defer this.mu.Unlock()

		modUser, err := this.container.modules.User()
		if nil != err {
			panic(err)
		}

		mAccess, err := this.container.modules.Access().MutationResolver()
		if nil != err {
			panic(err)
		}

		this.mutation = &mutationResolver{
			UserMutationResolver:   modUser.Mutation,
			AccessMutationResolver: mAccess,
		}
	}

	return this.mutation
}

func (this *resolvers) getSession() *sessionResolver {
	if nil == this.session {
		this.mu.Lock()
		defer this.mu.Unlock()

		this.session = &sessionResolver{
			container: this.container,
		}
	}

	return this.session
}

func (this *resolvers) getUser() *userResolver {
	if nil == this.user {
		this.mu.Lock()
		defer this.mu.Unlock()

		this.user = &userResolver{container: this.container}
	}

	return this.user
}

func (this *rootResolver) Mutation() gql.MutationResolver {
	return this.container.gql.getMutation()
}

func (this *rootResolver) Query() gql.QueryResolver {
	return this.container.gql.getQuery()
}

func (this *rootResolver) Session() gql.SessionResolver {
	return this.container.gql.getSession()
}

func (this *rootResolver) User() gql.UserResolver {
	return this.container.gql.getUser()
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

func (this userResolver) Name(ctx context.Context, obj *user_model.User) (*user_model.UserName, error) {
	panic("implement me")
}

func (this userResolver) Emails(ctx context.Context, obj *user_model.User) (*user_model.UserEmails, error) {
	panic("implement me")
}
