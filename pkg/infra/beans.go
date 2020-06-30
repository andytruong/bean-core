package infra

import (
	"bean/pkg/access"
	"bean/pkg/integration/s3"
	"bean/pkg/namespace"
	"bean/pkg/user"
	"bean/pkg/util"
)

type (
	beans struct {
		can       *Can
		user      *user.UserBean
		namespace *namespace.NamespaceBean
		access    *access.AccessBean
		s3        *s3.S3IntegrationBean
	}
)

func (this *beans) List() []util.Bean {
	mUser, _ := this.User()
	mNamespace, _ := this.Namespace()
	mAccess, _ := this.Access()

	return []util.Bean{mUser, mNamespace, mAccess}
}

func (this *beans) User() (*user.UserBean, error) {
	var err error

	if nil == this.user {
		db, err := this.can.dbs.master()
		if nil != err {
			return nil, err
		}

		this.user = user.NewUserBean(
			db,
			this.can.logger,
			this.can.Identifier(),
		)
	}

	return this.user, err
}

func (this *beans) Namespace() (*namespace.NamespaceBean, error) {
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

		this.namespace = namespace.NewNamespaceBean(
			db,
			this.can.logger,
			this.can.Identifier(),
			mUser,
			this.can.Beans.Namespace,
		)
	}

	return this.namespace, err
}

func (this *beans) Access() (*access.AccessBean, error) {
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

		this.access = access.NewAccessBean(
			db,
			this.can.Identifier(),
			this.can.logger,
			mUser,
			mNamespace,
			this.can.Beans.Access,
		)
	}

	return this.access, nil
}

func (this *beans) S3() (*s3.S3IntegrationBean, error) {
	if nil == this.s3 {
		db, err := this.can.dbs.master()
		if nil != err {
			return nil, err
		}

		this.s3 = s3.NewS3Integration(db, this.can.Identifier(), this.can.logger, this.can.Beans.S3)
	}

	return this.s3, nil
}
