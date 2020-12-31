package s3

import (
	"context"
	"time"

	"bean/components/connect"
	"bean/components/util"
	appModel "bean/pkg/app/model"
	"bean/pkg/integration/s3/model"
	"bean/pkg/integration/s3/model/dto"
)

type policyService struct {
	bundle *S3Bundle
}

func (srv *policyService) load(ctx context.Context, appId string, id string) (*model.Policy, error) {
	policy := &model.Policy{}
	err := connect.ContextToDB(ctx).
		Where("application_id = ?", appId).
		Where("id = ?", id).
		First(policy).Error
	if nil != err {
		return nil, err
	}

	return policy, nil
}

func (srv *policyService) create(ctx context.Context, appId string, kind model.PolicyKind, value string) (*model.Policy, error) {
	policy := &model.Policy{
		ID:            srv.bundle.idr.MustULID(),
		ApplicationId: appId,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		Kind:          kind,
		Value:         value,
	}

	err := connect.ContextToDB(ctx).Create(&policy).Error
	if nil != err {
		return nil, err
	}

	return policy, nil
}

func (srv *policyService) loadByApplicationId(ctx context.Context, appId string) ([]*model.Policy, error) {
	db := connect.ContextToDB(ctx)
	policies := []*model.Policy{}
	err := db.Where("application_id = ?", appId).Find(&policies).Error
	if nil != err {
		return nil, err
	}

	return policies, nil
}

func (srv *policyService) onAppCreate(ctx context.Context, app *appModel.Application, in []dto.S3ApplicationPolicyCreateInput) error {
	for _, input := range in {
		_, err := srv.create(ctx, app.ID, input.Kind, input.Value)
		if nil != err {
			return err
		}
	}

	return nil
}

func (srv *policyService) onAppUpdate(ctx context.Context, app *appModel.Application, in *dto.S3ApplicationPolicyMutationInput) error {
	if nil == in {
		return nil
	}

	db := connect.ContextToDB(ctx)
	useless := true

	if nil != in.Create {
		for _, input := range in.Create {
			_, err := srv.create(ctx, app.ID, input.Kind, input.Value)
			if nil != err {
				return err
			}

			useless = false
		}
	}

	if nil != in.Update {
		for _, input := range in.Update {
			if policy, err := srv.load(ctx, app.ID, input.Id); nil != err {
				return err
			} else {
				policy.Value = input.Value
				policy.UpdatedAt = time.Now()
				err := db.Save(policy).Error
				if nil != err {
					return err
				}

				useless = false
			}
		}
	}

	if nil != in.Delete {
		for _, input := range in.Delete {
			if policy, err := srv.load(ctx, app.ID, input.Id); nil != err {
				return err
			} else {
				err := db.Delete(policy).Error
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
