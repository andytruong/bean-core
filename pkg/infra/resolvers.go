package infra

import (
	"context"
	"sync"

	"bean/pkg/access"
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
		*user.UserModule
		*namespace.NamespaceModule
		*access.AccessModule
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

		mUser, _ := this.container.modules.User()
		mNamespace, _ := this.container.modules.Namespace()

		this.query = &queryResolver{
			UserModule:      mUser,
			NamespaceModule: mNamespace,
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
			UserModule:      modUser,
			AccessModule:    modAccess,
			NamespaceModule: modNamespace,
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

func (this mutationResolver) Ping(ctx context.Context) (string, error) {
	return "pong", nil
}
