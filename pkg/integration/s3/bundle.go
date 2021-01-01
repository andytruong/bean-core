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
)

func NewS3Integration(
	idr *scalar.Identifier,
	logger *zap.Logger,
	conf *S3Configuration,
	appBundle *app.AppBundle,
	configBundle *config.ConfigBundle,
) *S3Bundle {
	this := &S3Bundle{
		idr:    idr,
		logger: logger,
		cnf:    conf,
	}

	this.appBundle = appBundle
	this.configBundle = configBundle
	this.AppService = &ApplicationService{bundle: this}
	this.credentialService = &credentialService{bundle: this}
	this.policyService = &policyService{bundle: this}
	this.resolvers = newResolvers(this)

	return this
}

type S3Bundle struct {
	module.AbstractBundle

	appBundle    *app.AppBundle
	configBundle *config.ConfigBundle

	idr               *scalar.Identifier
	logger            *zap.Logger
	cnf               *S3Configuration
	AppService        *ApplicationService
	credentialService *credentialService
	policyService     *policyService
	resolvers         map[string]interface{}
}

func (S3Bundle) Name() string {
	return "S3"
}

func (bundle S3Bundle) Dependencies() []module.Bundle {
	return []module.Bundle{bundle.configBundle, bundle.appBundle}
}

func (bundle S3Bundle) Migrate(ctx context.Context, driver string) error {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return nil
	}

	runner := connect.Runner{
		Logger: bundle.logger,
		Driver: driver,
		Bundle: "integration.s3",
		Dir:    path.Dir(filename) + "/model/migration/",
	}

	if err := runner.Run(ctx); nil != err {
		return err
	}

	// srv.bundle.configBundle.BucketService.Create()

	return nil
}

func (bundle *S3Bundle) GraphqlResolver() map[string]interface{} {
	return bundle.resolvers
}
