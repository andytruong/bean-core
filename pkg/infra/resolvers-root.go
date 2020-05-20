package infra

import "bean/pkg/infra/gql"

func (this *root) Mutation() gql.MutationResolver {
	return this.container.gql.getResolver()
}

func (this *root) Query() gql.QueryResolver {
	return this.container.gql.getResolver()
}

func (this *root) Session() gql.SessionResolver {
	return this.container.gql.getSession()
}

func (this *root) User() gql.UserResolver {
	module, _ := this.container.modules.User()

	return module
}

func (this *root) Namespace() gql.NamespaceResolver {
	module, _ := this.container.modules.Namespace()

	return module
}

func (this *root) UserEmail() gql.UserEmailResolver {
	module, _ := this.container.modules.User()

	return module
}
