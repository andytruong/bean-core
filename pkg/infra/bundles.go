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
		s3        *s3.S3Bundle
		mailer    *mailer.MailerBundle
	}
)

func (list *bundles) User() (*user.UserBundle, error) {
	var err error

	if nil == list.user {
		db, err := list.container.dbs.master()
		if nil != err {
			return nil, err
		}

		list.user = user.NewUserBundle(
			db,
			list.container.logger,
			list.container.Identifier(),
		)
	}

	return list.user, err
}

func (list *bundles) Space() (*space.SpaceBundle, error) {
	var err error

	if nil == list.space {
		db, err := list.container.dbs.master()
		if nil != err {
			return nil, err
		}

		mUser, err := list.User()
		if nil != err {
			return nil, err
		}

		list.space = space.NewSpaceBundle(
			db,
			list.container.logger,
			list.container.Identifier(),
			mUser,
			list.container.Bundles.Space,
		)
	}

	return list.space, err
}

func (list *bundles) Access() (*access.AccessBundle, error) {
	if nil == list.access {
		db, err := list.container.dbs.master()
		if nil != err {
			return nil, err
		}

		mUser, err := list.User()
		if nil != err {
			return nil, err
		}

		mSpace, err := list.Space()
		if nil != err {
			return nil, err
		}

		list.access = access.NewAccessBundle(
			db,
			list.container.Identifier(),
			list.container.logger,
			mUser,
			mSpace,
			list.container.Bundles.Access,
		)
	}

	return list.access, nil
}

func (list *bundles) Mailer() (*mailer.MailerBundle, error) {
	if nil == list.mailer {
		list.mailer = mailer.NewMailerBundle(list.container.Bundles.Integration.Mailer, list.container.logger)
	}

	return list.mailer, nil
}

func (list *bundles) S3() (*s3.S3Bundle, error) {
	if nil == list.s3 {
		db, err := list.container.dbs.master()
		if nil != err {
			return nil, err
		}

		list.s3 = s3.NewS3Integration(db, list.container.Identifier(), list.container.logger, list.container.Bundles.Integration.S3)
	}

	return list.s3, nil
}
