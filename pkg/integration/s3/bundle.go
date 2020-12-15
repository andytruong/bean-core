package s3

import (
	"context"
	"path"
	"runtime"
	
	"go.uber.org/zap"
	"gorm.io/gorm"
	
	"bean/components/module"
	"bean/components/module/migrate"
	"bean/components/unique"
	"bean/pkg/integration/s3/model"
	"bean/pkg/integration/s3/model/dto"
)

func NewS3Integration(
	db *gorm.DB,
	id *unique.Identifier,
	logger *zap.Logger,
	conf *S3Configuration,
) *S3IntegrationBundle {
	this := &S3IntegrationBundle{
		db:     db,
		id:     id,
		logger: logger,
		config: conf,
	}
	
	this.AppService = &ApplicationService{bundle: this}
	this.credentialService = &credentialService{bundle: this}
	this.policyService = &policyService{bundle: this}
	this.resolvers = this.newResolver()
	
	return this
}

type S3IntegrationBundle struct {
	module.AbstractBundle
	
	db     *gorm.DB
	id     *unique.Identifier
	logger *zap.Logger
	config *S3Configuration
	
	AppService        *ApplicationService
	credentialService *credentialService
	policyService     *policyService
	resolvers         map[string]interface{}
}

func (this S3IntegrationBundle) Migrate(tx *gorm.DB, driver string) error {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return nil
	}
	
	runner := migrate.Runner{
		Tx:     tx,
		Logger: this.logger,
		Driver: driver,
		Bean:   "integration.s3",
		Dir:    path.Dir(filename) + "/model/migration/",
	}
	
	return runner.Run()
}

func (this *S3IntegrationBundle) GraphqlResolver() map[string]interface{} {
	return this.resolvers
}

func (this *S3IntegrationBundle) newResolver() map[string]interface{} {
	return map[string]interface{}{
		"Mutation": map[string]interface{}{
			"S3ApplicationCreate": func(ctx context.Context, input *dto.S3ApplicationCreateInput) (*dto.S3ApplicationMutationOutcome, error) {
				return this.AppService.Create(ctx, input)
			},
			"S3ApplicationUpdate": func(ctx context.Context, input *dto.S3ApplicationUpdateInput) (*dto.S3ApplicationMutationOutcome, error) {
				return this.AppService.Update(ctx, input)
			},
			"S3UploadToken": func(ctx context.Context, input dto.S3UploadTokenInput) (map[string]interface{}, error) {
				return this.AppService.S3UploadToken(ctx, input)
			},
		},
		"Application": map[string]interface{}{
			"Polices": func(ctx context.Context, obj *model.Application) ([]*model.Policy, error) {
				return this.policyService.loadByApplicationId(ctx, obj.ID)
			},
			"Credentials": func(ctx context.Context, obj *model.Application) (*model.Credentials, error) {
				return this.credentialService.loadByApplicationId(ctx, obj.ID)
			},
		},
	}
}
