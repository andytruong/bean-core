package infra

import (
	"bean/pkg/access"
	"bean/pkg/namespace"
	"bean/pkg/user"
	"bean/pkg/util"
)

type (
	modules struct {
		can       *Can
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
		db, err := this.can.dbs.master()
		if nil != err {
			return nil, err
		}

		this.user = user.NewUserModule(
			db,
			this.can.logger,
			this.can.Identifier(),
		)
	}

	return this.user, err
}

func (this *modules) Namespace() (*namespace.NamespaceModule, error) {
	var err error

	if nil == this.namespace {
		db, err := this.can.dbs.master()
		if nil != err {
			return nil, err
		}

		mUser, err := this.User()
		if nil != err {
			return nil, err
		}

		this.namespace = namespace.NewNamespaceModule(
			db,
			this.can.logger,
			this.can.Identifier(),
			mUser,
			this.can.Modules.Namespace,
		)
	}

	return this.namespace, err
}

func (this *modules) Access() (*access.AccessModule, error) {
	if nil == this.access {
		db, err := this.can.dbs.master()
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
			this.can.Identifier(),
			this.can.logger,
			mUser,
			mNamespace,
			this.can.Modules.Access,
		)
	}

	return this.access, nil
}
