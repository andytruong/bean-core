package infra

import (
	"sync"

	"bean/pkg/access"
	"bean/pkg/infra/gql"
	"bean/pkg/integration/s3"
	"bean/pkg/namespace"
	"bean/pkg/user"
)

type (
	graph struct {
		can       *Can
		mu        *sync.Mutex
		resolvers *resolvers
	}

	resolvers struct {
		*user.UserBean
		*user.UserQueryResolver
		*user.UserMutationResolver
		*namespace.NamespaceBean
		*namespace.NamespaceQueryResolver
		*access.AccessBean
		*s3.ApplicationResolver
	}
)

func (this *graph) Application() gql.ApplicationResolver {
	bean, _ := this.can.beans.S3()

	return bean.CoreApp.Resolver
}

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
		bS3, _ := this.can.beans.S3()

		this.resolvers = &resolvers{
			UserBean:               bUser,
			UserQueryResolver:      bUser.Resolvers.Query,
			UserMutationResolver:   bUser.Resolvers.Mutation,
			AccessBean:             bAccess,
			NamespaceBean:          bNamespace,
			NamespaceQueryResolver: bNamespace.Resolvers.Query,
			ApplicationResolver:    bS3.CoreApp.Resolver,
		}
	}

	return this.resolvers
}
