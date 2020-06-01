package infra

import (
	"sync"

	"bean/pkg/access"
	"bean/pkg/infra/gql"
	"bean/pkg/namespace"
	"bean/pkg/user"
	"bean/pkg/util/api"
)

type (
	graph struct {
		container *Container
		mu        *sync.Mutex
		resolvers *resolvers
	}

	resolvers struct {
		*api.Resolver
		*user.UserModule
		*namespace.NamespaceModule
		*access.AccessModule
	}
)

func (this *graph) SessionCreateOutcome() gql.SessionCreateOutcomeResolver {
	module, _ := this.container.modules.Access()

	return module
}

func (this *graph) Membership() gql.MembershipResolver {
	module, _ := this.container.modules.Namespace()

	return module.MembershipResolver()
}

func (this *graph) Mutation() gql.MutationResolver {
	return this.getResolvers()
}

func (this *graph) Query() gql.QueryResolver {
	return this.getResolvers()
}

func (this *graph) Session() gql.SessionResolver {
	module, _ := this.container.modules.Access()

	return module.SessionResolver
}

func (this *graph) User() gql.UserResolver {
	module, _ := this.container.modules.User()

	return module
}

func (this *graph) Namespace() gql.NamespaceResolver {
	module, _ := this.container.modules.Namespace()

	return module
}

func (this *graph) UserEmail() gql.UserEmailResolver {
	module, _ := this.container.modules.User()

	return module
}

func (this *graph) getResolvers() *resolvers {
	if nil == this.resolvers {
		this.mu.Lock()
		defer this.mu.Unlock()

		modUser, _ := this.container.modules.User()
		modAccess, _ := this.container.modules.Access()
		modNamespace, _ := this.container.modules.Namespace()

		this.resolvers = &resolvers{
			UserModule:      modUser,
			AccessModule:    modAccess,
			NamespaceModule: modNamespace,
		}
	}

	return this.resolvers
}
