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

func (bundles *bundleList) Get() []module.Bundle {
	userBundle, _ := bundles.User()
	spaceBundle, _ := bundles.Space()
	appBundle, _ := bundles.App()
	accessBundle, _ := bundles.Access()
	s3Bundle, _ := bundles.S3()
	mailerBundle, _ := bundles.Mailer()

	return []module.Bundle{userBundle, spaceBundle, appBundle, accessBundle, s3Bundle, mailerBundle}
}

func (bundles *bundleList) User() (*user.Bundle, error) {
	var err error

	if nil == bundles.user {
		bundles.user = user.NewUserBundle(
			bundles.container.logger,
			bundles.container.idr,
		)
	}

	return bundles.user, err
}

func (bundles *bundleList) Space() (*space.Bundle, error) {
	var err error

	if nil == bundles.space {
		mUser, err := bundles.User()
		if nil != err {
			return nil, err
		}

		bundles.space = space.NewSpaceBundle(
			bundles.container.logger,
			bundles.container.idr,
			mUser,
			bundles.container.Config.Bundles.Space,
		)
	}

	return bundles.space, err
}

func (bundles *bundleList) Config() (*config.Bundle, error) {
	if nil == bundles.config {
		bundles.config = config.NewConfigBundle(
			bundles.container.idr,
			bundles.container.logger,
		)
	}

	return bundles.config, nil
}

func (bundles *bundleList) Access() (*access.Bundle, error) {
	if nil == bundles.access {
		userBundle, err := bundles.User()
		if nil != err {
			return nil, err
		}

		spaceBundle, err := bundles.Space()
		if nil != err {
			return nil, err
		}

		bundles.access, err = access.NewAccessBundle(
			bundles.container.idr,
			bundles.container.logger,
			userBundle,
			spaceBundle,
			bundles.container.Config.Bundles.Access,
		)

		if nil != err {
			return nil, err
		}
	}

	return bundles.access, nil
}

func (bundles *bundleList) Mailer() (*mailer.Bundle, error) {
	if nil == bundles.mailer {
		bundles.mailer = mailer.NewMailerBundle(bundles.container.Config.Bundles.Integration.Mailer, bundles.container.logger)
	}

	return bundles.mailer, nil
}

func (bundles *bundleList) App() (*app.Bundle, error) {
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
			bundles.container.idr,
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

func (bundles *bundleList) S3() (*s3.Bundle, error) {
	if nil == bundles.s3 {
		appBundle, err := bundles.App()
		if nil != err {
			return nil, err
		}

		configBundle, _ := bundles.Config()

		bundles.s3 = s3.NewS3Integration(
			bundles.container.idr,
			bundles.container.logger,
			bundles.container.Config.Bundles.Integration.S3,
			appBundle,
			configBundle,
		)
	}

	return bundles.s3, nil
}
