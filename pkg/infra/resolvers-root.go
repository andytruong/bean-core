package infra

import "bean/pkg/infra/gql"

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
	module, _ := this.container.modules.User()

	return module
}

func (this *rootResolver) Namespace() gql.NamespaceResolver {
	module, _ := this.container.modules.Namespace()

	return module
}

func (this *rootResolver) UserEmail() gql.UserEmailResolver {
	module, _ := this.container.modules.User()

	return module
}
