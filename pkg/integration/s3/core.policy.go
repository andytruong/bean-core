package s3

import (
	"time"

	"gorm.io/gorm"

	"bean/pkg/integration/s3/model"
	"bean/pkg/integration/s3/model/dto"
	"bean/pkg/util/connect"
)

type corePolicy struct {
	bean *S3IntegrationBean
}

func (this *corePolicy) onAppCreate(tx *gorm.DB, app *model.Application, in []dto.S3ApplicationPolicyCreateInput) error {
	for _, input := range in {
		policy := model.Policy{
			ID:            this.bean.id.MustULID(),
			ApplicationId: app.ID,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
			Kind:          input.Kind,
			Value:         input.Value,
		}

		err := tx.Table(connect.TableIntegrationS3Policy).Create(&policy).Error
		if nil != err {
			return err
		}
	}

	return nil
}
