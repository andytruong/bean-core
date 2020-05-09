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
		db, err := this.container.dbs.get("master")
		if nil != err {
			return nil, err
		}

		this.user, err = user.NewUserModule(
			db,
			this.container.logger,
			this.container.Identifier(),
		)
	}

	return this.user, err
}

func (this *modules) Access() *access.AccessModule {
	if nil == this.access {
		this.access = access.NewAccessModule()
	}

	return this.access
}
