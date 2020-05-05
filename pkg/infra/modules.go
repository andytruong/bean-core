package infra

import (
	"bean/pkg/access"
	"bean/pkg/user"
)

type (
	modules struct {
		container *Container
		user      *user.UserModule
		access    *access.AccessModule
	}
)

func (this *modules) User() *user.UserModule {
	if nil == this.user {
		this.user = user.NewUserService(this.container.DB, this.container.Logger)
	}

	return this.user
}

func (this *modules) Access() *access.AccessModule {
	if nil == this.access {
		this.access = access.NewAccessModule()
	}

	return this.access
}
