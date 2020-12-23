package s3

import (
	"context"
	"time"

	"gorm.io/gorm"

	"bean/components/util"
	model2 "bean/pkg/app/model"
	"bean/pkg/integration/s3/model"
	"bean/pkg/integration/s3/model/dto"
)

type policyService struct {
	bundle *S3Bundle
}

func (service *policyService) load(tx *gorm.DB, appId string, id string) (*model.Policy, error) {
	policy := &model.Policy{}
	err := tx.
		Where("application_id = ?", appId).
		Where("id = ?", id).
		First(policy).Error
	if nil != err {
		return nil, err
	}

	return policy, nil
}

func (service *policyService) create(tx *gorm.DB, appId string, kind model.PolicyKind, value string) (*model.Policy, error) {
	policy := &model.Policy{
		ID:            service.bundle.idr.MustULID(),
		ApplicationId: appId,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		Kind:          kind,
		Value:         value,
	}

	err := tx.Create(&policy).Error
	if nil != err {
		return nil, err
	}

	return policy, nil
}

func (service *policyService) loadByApplicationId(ctx context.Context, appId string) ([]*model.Policy, error) {
	policies := []*model.Policy{}

	err := service.bundle.con.WithContext(ctx).
		Where("application_id = ?", appId).
		Find(&policies).
		Error
	if nil != err {
		return nil, err
	}

	return policies, nil
}

func (service *policyService) onAppCreate(tx *gorm.DB, app *model2.Application, in []dto.S3ApplicationPolicyCreateInput) error {
	for _, input := range in {
		_, err := service.create(tx, app.ID, input.Kind, input.Value)
		if nil != err {
			return err
		}
	}

	return nil
}

func (service *policyService) onAppUpdate(tx *gorm.DB, app *model2.Application, in *dto.S3ApplicationPolicyMutationInput) error {
	if nil == in {
		return nil
	}

	useless := true

	if nil != in.Create {
		for _, input := range in.Create {
			_, err := service.create(tx, app.ID, input.Kind, input.Value)
			if nil != err {
				return err
			}

			useless = false
		}
	}

	if nil != in.Update {
		for _, input := range in.Update {
			if policy, err := service.load(tx, app.ID, input.Id); nil != err {
				return err
			} else {
				policy.Value = input.Value
				policy.UpdatedAt = time.Now()
				err := tx.Save(policy).Error
				if nil != err {
					return err
				}

				useless = false
			}
		}
	}

	if nil != in.Delete {
		for _, input := range in.Delete {
			if policy, err := service.load(tx, app.ID, input.Id); nil != err {
				return err
			} else {
				err := tx.Delete(policy).Error
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
