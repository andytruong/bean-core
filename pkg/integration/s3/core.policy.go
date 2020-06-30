package s3

import (
	"time"

	"gorm.io/gorm"

	"bean/pkg/integration/s3/model"
	"bean/pkg/integration/s3/model/dto"
	"bean/pkg/util"
	"bean/pkg/util/connect"
)

type corePolicy struct {
	bean *S3IntegrationBean
}

func (this *corePolicy) load(tx *gorm.DB, appId string, id string) (*model.Policy, error) {
	policy := &model.Policy{}
	err := tx.Table(connect.TableIntegrationS3Policy).
		Where("application_id = ?", appId).
		Where("id = ?", id).
		First(policy).Error
	if nil != err {
		return nil, err
	}

	return policy, nil
}

func (this *corePolicy) create(tx *gorm.DB, appId string, kind model.PolicyKind, value string) (*model.Policy, error) {
	policy := &model.Policy{
		ID:            this.bean.id.MustULID(),
		ApplicationId: appId,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		Kind:          kind,
		Value:         value,
	}

	err := tx.Table(connect.TableIntegrationS3Policy).Create(&policy).Error
	if nil != err {
		return nil, err
	}

	return policy, nil
}

func (this *corePolicy) onAppCreate(tx *gorm.DB, app *model.Application, in []dto.S3ApplicationPolicyCreateInput) error {
	for _, input := range in {
		_, err := this.create(tx, app.ID, input.Kind, input.Value)
		if nil != err {
			return err
		}
	}

	return nil
}

func (this *corePolicy) onAppUpdate(tx *gorm.DB, app *model.Application, in *dto.S3ApplicationPolicyMutationInput) error {
	if nil == in {
		return nil
	}

	useless := true
	
	if nil != in.Create {
		for _, input := range in.Create {
			_, err := this.create(tx, app.ID, input.Kind, input.Value)
			if nil != err {
				return err
			}

			useless = false
		}
	}

	if nil != in.Update {
		for _, input := range in.Update {
			if policy, err := this.load(tx, app.ID, input.Id); nil != err {
				return err
			} else {
				policy.Value = input.Value
				policy.UpdatedAt = time.Now()
				err := tx.Table(connect.TableIntegrationS3Policy).Save(policy).Error
				if nil != err {
					return err
				}

				useless = false
			}
		}
	}

	if nil != in.Delete {
		for _, input := range in.Delete {
			if policy, err := this.load(tx, app.ID, input.Id); nil != err {
				return err
			} else {
				err := tx.Table(connect.TableIntegrationS3Policy).Delete(policy).Error
				if nil != err {
					return err
				}

				useless = false
			}
		}
	}

	if useless {
		return util.ErrorUselessInput
	}

	return nil
}
