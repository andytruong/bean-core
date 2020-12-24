package infra

import (
	"bean/components/module"
	"bean/pkg/access"
	"bean/pkg/app"
	"bean/pkg/config"
	"bean/pkg/integration/mailer"
	"bean/pkg/integration/s3"
	"bean/pkg/space"
	"bean/pkg/user"
)

// TODO: Generate this code
func (list *bundles) List() []module.Bundle {
	userBundle, _ := list.User()
	spaceBundle, _ := list.Space()
	appBundle, _ := list.App()
	accessBundle, _ := list.Access()
	s3Bundle, _ := list.S3()
	mailerBundle, _ := list.Mailer()

	return []module.Bundle{userBundle, spaceBundle, appBundle, accessBundle, s3Bundle, mailerBundle}
}

// TODO: Generate this code
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

// TODO: Generate this code
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

// TODO: Generate this code
func (list *bundles) Config() (*config.ConfigBundle, error) {
	var err error

	if nil == list.config {
		list.config = config.NewConfigBundle(
			list.container.Identifier(),
			list.container.logger,
		)
	}

	return list.config, err
}

// TODO: Generate this code
func (list *bundles) Access() (*access.AccessBundle, error) {
	if nil == list.access {
		userBundle, err := list.User()
		if nil != err {
			return nil, err
		}

		spaceBundle, err := list.Space()
		if nil != err {
			return nil, err
		}

		list.access = access.NewAccessBundle(
			list.container.Identifier(),
			list.container.logger,
			userBundle,
			spaceBundle,
			list.container.Bundles.Access,
		)
	}

	return list.access, nil
}

// TODO: Generate this code
func (list *bundles) Mailer() (*mailer.MailerBundle, error) {
	if nil == list.mailer {
		list.mailer = mailer.NewMailerBundle(list.container.Bundles.Integration.Mailer, list.container.logger)
	}

	return list.mailer, nil
}

// TODO: Generate this code
func (list *bundles) App() (*app.AppBundle, error) {
	if nil == list.app {
		spaceBundle, err := list.Space()
		if nil != err {
			return nil, err
		}

		configBundle, err := list.Config()
		if nil != err {
			return nil, err
		}

		list.app, err = app.NewApplicationBundle(
			list.container.Identifier(),
			list.container.logger,
			spaceBundle,
			configBundle,
		)

		if nil != err {
			return nil, err
		}
	}

	return list.app, nil
}

// TODO: Generate this code
func (list *bundles) S3() (*s3.S3Bundle, error) {
	if nil == list.s3 {
		appBundle, err := list.App()
		if nil != err {
			return nil, err
		}

		list.s3 = s3.NewS3Integration(
			list.container.Identifier(),
			list.container.logger,
			list.container.Bundles.Integration.S3,
			appBundle,
		)
	}

	return list.s3, nil
}
