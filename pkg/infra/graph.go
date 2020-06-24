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
		can       *Can
		mu        *sync.Mutex
		resolvers *resolvers
	}

	resolvers struct {
		*api.Resolver
		*user.UserBean
		*user.UserQueryResolver
		*user.UserMutationResolver
		*namespace.NamespaceBean
		*namespace.NamespaceQueryResolver
		*access.AccessBean
	}
)

func (this *graph) MembershipConnection() gql.MembershipConnectionResolver {
	bean, _ := this.can.beans.Namespace()

	return bean.MembershipResolver()
}

func (this *graph) Membership() gql.MembershipResolver {
	bean, _ := this.can.beans.Namespace()

	return bean.MembershipResolver()
}

func (this *graph) Mutation() gql.MutationResolver {
	return this.getResolvers()
}

func (this *graph) Query() gql.QueryResolver {
	return this.getResolvers()
}

func (this *graph) Session() gql.SessionResolver {
	bean, _ := this.can.beans.Access()

	return bean.SessionResolver
}

func (this *graph) User() gql.UserResolver {
	bean, _ := this.can.beans.User()

	return bean.Resolvers.Object
}

func (this *graph) Namespace() gql.NamespaceResolver {
	bean, _ := this.can.beans.Namespace()

	return bean.Resolvers.Object
}

func (this *graph) UserEmail() gql.UserEmailResolver {
	bean, _ := this.can.beans.User()

	return bean.Resolvers.Object
}

func (this *graph) getResolvers() *resolvers {
	if nil == this.resolvers {
		this.mu.Lock()
		defer this.mu.Unlock()

		bUser, _ := this.can.beans.User()
		bAccess, _ := this.can.beans.Access()
		bNamespace, _ := this.can.beans.Namespace()

		this.resolvers = &resolvers{
			UserBean:               bUser,
			UserQueryResolver:      bUser.Resolvers.Query,
			UserMutationResolver:   bUser.Resolvers.Mutation,
			AccessBean:             bAccess,
			NamespaceBean:          bNamespace,
			NamespaceQueryResolver: bNamespace.Resolvers.Query,
		}
	}

	return this.resolvers
}
