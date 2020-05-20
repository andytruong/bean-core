package infra

import (
	"context"
	"sync"

	"bean/pkg/access"
	"bean/pkg/infra/gql"
	"bean/pkg/namespace"
	"bean/pkg/user"
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
	}

	rootResolver struct {
		container *Container
	}

	mutationResolver struct {
		*user.UserMutationResolver
		*namespace.NamespaceMutationResolver
		*access.AccessMutationResolver
	}

	queryResolver struct {
		*access.AccessQueryResolver
		user.UserQueryResolver
	}

	sessionResolver struct {
		container *Container
	}
)

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

		modUser, _ := this.container.modules.User()
		modAccess, _ := this.container.modules.Access()
		modNamespace, _ := this.container.modules.Namespace()

		this.mutation = &mutationResolver{
			UserMutationResolver:      modUser.Mutation,
			AccessMutationResolver:    modAccess.Mutation,
			NamespaceMutationResolver: modNamespace.Mutation,
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
	module, err := this.container.modules.User()
	if nil != err {
		panic(err)
	}

	return module.Model
}

func (this *rootResolver) Namespace() gql.NamespaceResolver {
	module, err := this.container.modules.Namespace()
	if nil != err {
		panic(err)
	}

	return module.Model
}

func (this *rootResolver) UserEmail() gql.UserEmailResolver {
	module, err := this.container.modules.User()
	if nil != err {
		panic(err)
	}

	return module.Email
}

func (this mutationResolver) Ping(ctx context.Context) (string, error) {
	return "pong", nil
}

func (this queryResolver) Ping(ctx context.Context) (string, error) {
	return "pong", nil
}
