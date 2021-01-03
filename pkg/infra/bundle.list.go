package infra

// TODO: Generate this code

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

func (bundles *BundleList) Get() []module.Bundle {
	userBundle, _ := bundles.User()
	spaceBundle, _ := bundles.Space()
	appBundle, _ := bundles.App()
	accessBundle, _ := bundles.Access()
	s3Bundle, _ := bundles.S3()
	mailerBundle, _ := bundles.Mailer()

	return []module.Bundle{userBundle, spaceBundle, appBundle, accessBundle, s3Bundle, mailerBundle}
}

func (bundles *BundleList) User() (*user.Bundle, error) {
	var err error

	if nil == bundles.user {
		bundles.user = user.NewUserBundle(
			bundles.container.logger,
			bundles.container.identifier,
		)
	}

	return bundles.user, err
}

func (bundles *BundleList) Space() (*space.Bundle, error) {
	var err error

	if nil == bundles.space {
		mUser, err := bundles.User()
		if nil != err {
			return nil, err
		}

		bundles.space = space.NewSpaceBundle(
			bundles.container.logger,
			bundles.container.identifier,
			mUser,
			bundles.container.Config.Bundles.Space,
		)
	}

	return bundles.space, err
}

func (bundles *BundleList) Config() (*config.Bundle, error) {
	if nil == bundles.config {
		bundles.config = config.NewConfigBundle(
			bundles.container.identifier,
			bundles.container.logger,
		)
	}

	return bundles.config, nil
}

func (bundles *BundleList) Access() (*access.Bundle, error) {
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
			bundles.container.identifier,
			bundles.container.logger,
			userBundle,
			spaceBundle,
			bundles.container.Config.Bundles.Access,
		)
	}

	return bundles.access, nil
}

func (bundles *BundleList) Mailer() (*mailer.Bundle, error) {
	if nil == bundles.mailer {
		bundles.mailer = mailer.NewMailerBundle(bundles.container.Config.Bundles.Integration.Mailer, bundles.container.logger)
	}

	return bundles.mailer, nil
}

func (bundles *BundleList) App() (*app.Bundle, error) {
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
			bundles.container.identifier,
			bundles.container.logger,
			bundles.container.hook,
			spaceBundle,
			configBundle,
		)

		if nil != err {
			return nil, err
		}
	}

	return bundles.app, nil
}

func (bundles *BundleList) S3() (*s3.Bundle, error) {
	if nil == bundles.s3 {
		appBundle, err := bundles.App()
		if nil != err {
			return nil, err
		}

		configBundle, _ := bundles.Config()

		bundles.s3 = s3.NewS3Integration(
			bundles.container.identifier,
			bundles.container.logger,
			bundles.container.Config.Bundles.Integration.S3,
			appBundle,
			configBundle,
		)
	}

	return bundles.s3, nil
}
