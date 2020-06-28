package s3

import (
	"gorm.io/gorm"

	"bean/pkg/integration/s3/model"
	"bean/pkg/integration/s3/model/dto"
)

type corePolicy struct {
	bean *S3IntegrationBean
}

func (this *corePolicy) onAppCreate(tx *gorm.DB, app *model.Application, in []dto.S3ApplicationPolicyCreateInput) error {
	return nil // WIP
}
