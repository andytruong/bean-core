package infra

import (
	"sync"
	
	"bean/pkg/access"
	"bean/pkg/infra/gql"
	"bean/pkg/integration/mailer"
	"bean/pkg/integration/s3"
	"bean/pkg/space"
	"bean/pkg/user"
)

type (
	graph struct {
		can       *Can
		mutex     *sync.Mutex
		resolvers *resolvers
	}
	
	resolvers struct {
		*user.UserBundle
		*user.UserQueryResolver
		*user.UserMutationResolver
		*space.SpaceBundle
		*space.SpaceQueryResolver
		*access.AccessBundle
		*s3.ApplicationResolver
		*mailer.MailerResolver
	}
)

func (this *graph) Application() gql.ApplicationResolver {
	bundle, _ := this.can.bundles.S3()
	
	return bundle.AppService.Resolver
}

func (this *graph) MembershipConnection() gql.MembershipConnectionResolver {
	bundle, _ := this.can.bundles.Space()
	
	return bundle.MembershipResolver()
}

func (this *graph) Membership() gql.MembershipResolver {
	bundle, _ := this.can.bundles.Space()
	
	return bundle.MembershipResolver()
}

func (this *graph) Mutation() gql.MutationResolver {
	return this.getResolvers()
}

func (this *graph) Query() gql.QueryResolver {
	return this.getResolvers()
}

func (this *graph) Session() gql.SessionResolver {
	bundle, _ := this.can.bundles.Access()
	
	return bundle.SessionResolver
}

func (this *graph) User() gql.UserResolver {
	bundle, _ := this.can.bundles.User()
	
	return bundle.Resolvers.Object
}

func (this *graph) Space() gql.SpaceResolver {
	bundle, _ := this.can.bundles.Space()
	
	return bundle.Resolvers.Object
}

func (this *graph) UserEmail() gql.UserEmailResolver {
	bundle, _ := this.can.bundles.User()
	
	return bundle.Resolvers.Object
}

func (this *graph) getResolvers() *resolvers {
	if nil == this.resolvers {
		this.mutex.Lock()
		defer this.mutex.Unlock()
		
		userBundle, _ := this.can.bundles.User()
		accessBundle, _ := this.can.bundles.Access()
		spaceBundle, _ := this.can.bundles.Space()
		s3Bundle, _ := this.can.bundles.S3()
		mailerBundle, _ := this.can.bundles.Mailer()
		mailerBundle.Resolver()
		
		this.resolvers = &resolvers{
			UserBundle:           userBundle,
			UserQueryResolver:    userBundle.Resolvers.Query,
			UserMutationResolver: userBundle.Resolvers.Mutation,
			AccessBundle:         accessBundle,
			SpaceBundle:          spaceBundle,
			SpaceQueryResolver:   spaceBundle.Resolvers.Query,
			ApplicationResolver:  s3Bundle.AppService.Resolver,
			MailerResolver:       mailerBundle.Resolver(),
		}
	}
	
	return this.resolvers
}
