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

func (this *modules) User() (*user.UserModule, error) {
	var err error

	if nil == this.user {
		this.user, err = user.NewUserService(this.container.DB, this.container.Logger, this.container.Identifier())
	}

	return this.user, err
}

func (this *modules) Access() *access.AccessModule {
	if nil == this.access {
		this.access = access.NewAccessModule()
	}

	return this.access
}
