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
		root      *root
		resolver  *resolver
	}

	root struct {
		container *Container
	}

	resolver struct {
		*user.UserModule
		*namespace.NamespaceModule
		*access.AccessModule
	}
)

func (this *resolvers) getRoot() *root {
	if nil == this.root {
		this.mu.Lock()
		defer this.mu.Unlock()

		this.root = &root{container: this.container}
	}

	return this.root
}

func (this *resolvers) getResolver() *resolver {
	if nil == this.resolver {
		this.mu.Lock()
		defer this.mu.Unlock()

		modUser, _ := this.container.modules.User()
		modAccess, _ := this.container.modules.Access()
		modNamespace, _ := this.container.modules.Namespace()

		this.resolver = &resolver{
			UserModule:      modUser,
			AccessModule:    modAccess,
			NamespaceModule: modNamespace,
		}
	}

	return this.resolver
}

func (this resolver) Ping(ctx context.Context) (string, error) {
	return "pong", nil
}
