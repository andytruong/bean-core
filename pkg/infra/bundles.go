package infra

import (
	"bean/pkg/access"
	"bean/pkg/integration/mailer"
	"bean/pkg/integration/s3"
	"bean/pkg/space"
	"bean/pkg/user"
)

type (
	bundles struct {
		container *Container
		user      *user.UserBundle
		space     *space.SpaceBundle
		access    *access.AccessBundle
		s3        *s3.S3IntegrationBundle
		mailer    *mailer.MailerIntegrationBundle
	}
)

func (this *bundles) User() (*user.UserBundle, error) {
	var err error

	if nil == this.user {
		db, err := this.container.dbs.master()
		if nil != err {
			return nil, err
		}

		this.user = user.NewUserBundle(
			db,
			this.container.logger,
			this.container.Identifier(),
		)
	}

	return this.user, err
}

func (this *bundles) Space() (*space.SpaceBundle, error) {
	var err error

	if nil == this.space {
		db, err := this.container.dbs.master()
		if nil != err {
			return nil, err
		}

		mUser, err := this.User()
		if nil != err {
			return nil, err
		}

		this.space = space.NewSpaceBundle(
			db,
			this.container.logger,
			this.container.Identifier(),
			mUser,
			this.container.Bundles.Space,
		)
	}

	return this.space, err
}

func (this *bundles) Access() (*access.AccessBundle, error) {
	if nil == this.access {
		db, err := this.container.dbs.master()
		if nil != err {
			return nil, err
		}

		mUser, err := this.User()
		if nil != err {
			return nil, err
		}

		mSpace, err := this.Space()
		if nil != err {
			return nil, err
		}

		this.access = access.NewAccessBundle(
			db,
			this.container.Identifier(),
			this.container.logger,
			mUser,
			mSpace,
			this.container.Bundles.Access,
		)
	}

	return this.access, nil
}

func (this *bundles) Mailer() (*mailer.MailerIntegrationBundle, error) {
	if nil == this.mailer {
		this.mailer = mailer.NewMailerIntegration(this.container.Bundles.Integration.Mailer, this.container.logger)
	}

	return this.mailer, nil
}

func (this *bundles) S3() (*s3.S3IntegrationBundle, error) {
	if nil == this.s3 {
		db, err := this.container.dbs.master()
		if nil != err {
			return nil, err
		}

		this.s3 = s3.NewS3Integration(db, this.container.Identifier(), this.container.logger, this.container.Bundles.Integration.S3)
	}

	return this.s3, nil
}
