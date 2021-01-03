package s3

import (
	"context"
	"path"
	"runtime"

	"go.uber.org/zap"

	"bean/components/connect"
	"bean/components/module"
	"bean/components/scalar"
	"bean/pkg/app"
	"bean/pkg/config"
	"bean/pkg/config/model/dto"
)

const (
	// unique ID of bundle
	_id = "01EV1JTZ9FHTB1ZBFMEZ4PYG1N"
)

func NewS3Integration(
	idr *scalar.Identifier,
	lgr *zap.Logger,
	cnf *Config,
	appBundle *app.Bundle,
	cnfBundle *config.Bundle,
) *Bundle {
	bundle := &Bundle{
		idr:          idr,
		lgr:          lgr,
		cnf:          cnf,
		appBundle:    appBundle,
		configBundle: cnfBundle,
	}

	bundle.uploadSrv = &uploadService{bundle: bundle}
	bundle.configSrv = &configService{bundle: bundle}
	bundle.resolvers = newResolvers(bundle)

	return bundle
}

type Bundle struct {
	module.AbstractBundle

	appBundle    *app.Bundle
	configBundle *config.Bundle

	cnf       *Config
	idr       *scalar.Identifier
	lgr       *zap.Logger
	configSrv *configService
	uploadSrv *uploadService
	resolvers map[string]interface{}
}

func (Bundle) Name() string {
	return "S3"
}

func (bundle Bundle) Dependencies() []module.Bundle {
	return []module.Bundle{bundle.configBundle, bundle.appBundle}
}

func (bundle Bundle) Migrate(ctx context.Context, driver string) error {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return nil
	}

	runner := connect.Runner{
		Logger: bundle.lgr,
		Driver: driver,
		Bundle: "integration.s3",
		Dir:    path.Dir(filename) + "/model/migration/",
	}

	if err := runner.Run(ctx); nil != err {
		return err
	}

	// save configuration buckets: credentials, policies
	{
		access := scalar.AccessModePrivate

		// config bucket -> credentials
		{
			out, err := bundle.configBundle.BucketService.Create(ctx, dto.BucketCreateInput{
				Slug:        scalar.NilString(credentialsConfigSlug),
				Access:      &access,
				Schema:      credentialsConfigSchema,
				IsPublished: true,
				HostId:      _id,
			})

			if nil != err {
				return err
			} else if out.Errors != nil {
				panic(out.Errors)
			}
		}

		// config bucket -> policies
		{
			out, err := bundle.configBundle.BucketService.Create(ctx, dto.BucketCreateInput{
				Slug:        scalar.NilString(uploadPolicyConfigSlug),
				Access:      &access,
				Schema:      uploadPolicyConfigSchema,
				IsPublished: true,
				HostId:      _id,
			})

			if nil != err {
				return err
			} else if out.Errors != nil {
				panic(out.Errors)
			}
		}
	}

	return nil
}

func (bundle *Bundle) GraphqlResolver() map[string]interface{} {
	return bundle.resolvers
}
