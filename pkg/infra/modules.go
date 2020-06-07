package infra

import (
	"bean/pkg/access"
	"bean/pkg/namespace"
	"bean/pkg/user"
	"bean/pkg/util"
)

type (
	modules struct {
		container *Container
		user      *user.UserModule
		namespace *namespace.NamespaceModule
		access    *access.AccessModule
	}
)

func (this *modules) List() []util.Module {
	mUser, _ := this.User()
	mNamespace, _ := this.Namespace()
	mAccess, _ := this.Access()

	return []util.Module{mUser, mNamespace, mAccess}
}

func (this *modules) User() (*user.UserModule, error) {
	var err error

	if nil == this.user {
		db, err := this.container.dbs.master()
		if nil != err {
			return nil, err
		}

		this.user = user.NewUserModule(
			db,
			this.container.logger,
			this.container.Identifier(),
		)
	}

	return this.user, err
}

func (this *modules) Namespace() (*namespace.NamespaceModule, error) {
	var err error

	if nil == this.namespace {
		db, err := this.container.dbs.master()
		if nil != err {
			return nil, err
		}

		mUser, err := this.User()
		if nil != err {
			return nil, err
		}

		this.namespace = namespace.NewNamespaceModule(
			db,
			this.container.logger,
			this.container.Identifier(),
			mUser,
			this.container.Modules.Namespace,
		)
	}

	return this.namespace, err
}

func (this *modules) Access() (*access.AccessModule, error) {
	if nil == this.access {
		db, err := this.container.dbs.master()
		if nil != err {
			return nil, err
		}

		mUser, err := this.User()
		if nil != err {
			return nil, err
		}

		mNamespace, err := this.Namespace()
		if nil != err {
			return nil, err
		}

		this.access = access.NewAccessModule(
			db,
			this.container.Identifier(),
			this.container.logger,
			mUser,
			mNamespace,
			this.container.Modules.Access,
		)
	}

	return this.access, nil
}
