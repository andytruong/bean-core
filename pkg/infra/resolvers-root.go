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
