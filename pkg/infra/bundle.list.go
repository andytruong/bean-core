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
func (bundles *BundleList) Get() []module.Bundle {
	userBundle, _ := bundles.User()
	spaceBundle, _ := bundles.Space()
	appBundle, _ := bundles.App()
	accessBundle, _ := bundles.Access()
	s3Bundle, _ := bundles.S3()
	mailerBundle, _ := bundles.Mailer()

	return []module.Bundle{userBundle, spaceBundle, appBundle, accessBundle, s3Bundle, mailerBundle}
}

// TODO: Generate this code
func (bundles *BundleList) User() (*user.UserBundle, error) {
	var err error

	if nil == bundles.user {
		db, err := bundles.container.dbs.Master()
		if nil != err {
			return nil, err
		}

		bundles.user = user.NewUserBundle(
			db,
			bundles.container.logger,
			bundles.container.Identifier(),
		)
	}

	return bundles.user, err
}

// TODO: Generate this code
func (bundles *BundleList) Space() (*space.SpaceBundle, error) {
	var err error

	if nil == bundles.space {
		db, err := bundles.container.dbs.Master()
		if nil != err {
			return nil, err
		}

		mUser, err := bundles.User()
		if nil != err {
			return nil, err
		}

		bundles.space = space.NewSpaceBundle(
			db,
			bundles.container.logger,
			bundles.container.Identifier(),
			mUser,
			bundles.container.Bundles.Space,
		)
	}

	return bundles.space, err
}

// TODO: Generate this code
func (bundles *BundleList) Config() (*config.ConfigBundle, error) {
	var err error

	if nil == bundles.config {
		bundles.config = config.NewConfigBundle(
			bundles.container.Identifier(),
			bundles.container.logger,
		)
	}

	return bundles.config, err
}

// TODO: Generate this code
func (bundles *BundleList) Access() (*access.AccessBundle, error) {
	if nil == bundles.access {
		userBundle, err := bundles.User()
		if nil != err {
			return nil, err
		}

		spaceBundle, err := bundles.Space()
		if nil != err {
			return nil, err
		}

		bundles.access = access.NewAccessBundle(
			bundles.container.Identifier(),
			bundles.container.logger,
			userBundle,
			spaceBundle,
			bundles.container.Bundles.Access,
		)
	}

	return bundles.access, nil
}

// TODO: Generate this code
func (bundles *BundleList) Mailer() (*mailer.MailerBundle, error) {
	if nil == bundles.mailer {
		bundles.mailer = mailer.NewMailerBundle(bundles.container.Bundles.Integration.Mailer, bundles.container.logger)
	}

	return bundles.mailer, nil
}

// TODO: Generate this code
func (bundles *BundleList) App() (*app.AppBundle, error) {
	if nil == bundles.app {
		spaceBundle, err := bundles.Space()
		if nil != err {
			return nil, err
		}

		configBundle, err := bundles.Config()
		if nil != err {
			return nil, err
		}

		bundles.app, err = app.NewApplicationBundle(
			bundles.container.Identifier(),
			bundles.container.logger,
			spaceBundle,
			configBundle,
		)

		if nil != err {
			return nil, err
		}
	}

	return bundles.app, nil
}

// TODO: Generate this code
func (bundles *BundleList) S3() (*s3.S3Bundle, error) {
	if nil == bundles.s3 {
		appBundle, err := bundles.App()
		if nil != err {
			return nil, err
		}

		bundles.s3 = s3.NewS3Integration(
			bundles.container.Identifier(),
			bundles.container.logger,
			bundles.container.Bundles.Integration.S3,
			appBundle,
		)
	}

	return bundles.s3, nil
}
