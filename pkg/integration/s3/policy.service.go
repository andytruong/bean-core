package s3

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	"gorm.io/gorm"

	"bean/components/scalar"
	"bean/components/util"
	configModel "bean/pkg/config/model"
	configDto "bean/pkg/config/model/dto"
	"bean/pkg/integration/s3/model"
	"bean/pkg/integration/s3/model/dto"
)

const (
	policyConfigSlug   = `bean.s3.policy.schema.v1`
	policyConfigSchema = `{
		"type":       "object",
		"required":   [],
		"properties": {
			"fileExtensions": {
				"type":  "array",
				"items": { "type": "string", "maxLength": 32 }
			},
			"rateLimit":      {
				"type":  "array",
				"items": {
					"type":       "object",
					"required":   ["value", "object", "interval"],
					"properties": {
						"value":    { "type": "string" },
						"object":   { "type": "string", "enum": ["user", "space"] },
						"interval": {
							"type": "string",
							"pattern": "^(\\d+) (minute|minutes|hour|hours|day|days)$"
						}
					}
				}
			}
		}
	}`
)

type policyService struct {
	bundle *S3Bundle
}

func (srv *policyService) load(ctx context.Context, appId string) (*model.S3UploadPolicy, error) {
	var (
		err      error
		bucket   *configModel.ConfigBucket
		variable *configModel.ConfigVariable
		policy   = &model.S3UploadPolicy{}
	)

	// load current bucket
	bucket, err = srv.bundle.configBundle.BucketService.Load(ctx, configDto.BucketKey{Slug: policyConfigSlug})
	if nil != err {
		return nil, err
	}

	// load current variable
	variable, err = srv.bundle.configBundle.VariableService.Load(ctx, configDto.VariableKey{BucketId: bucket.Id, Name: appId})
	if nil != err {
		return nil, err
	}

	err = json.Unmarshal([]byte(variable.Value), policy)
	if nil != err {
		return nil, err
	}

	policy.Id = variable.Id
	policy.Version = variable.Version

	return policy, nil
}

func (srv *policyService) save(ctx context.Context, in dto.UploadPolicyInput) (*model.S3UploadPolicy, error) {
	var (
		err    error
		bucket *configModel.ConfigBucket
		policy *model.S3UploadPolicy
	)

	// load current bucket
	bucket, err = srv.bundle.configBundle.BucketService.Load(ctx, configDto.BucketKey{Slug: policyConfigSlug})
	if nil != err {
		return nil, errors.Wrap(err, "bucket load error")
	}

	policy, err = srv.load(ctx, in.ApplicationId)
	if nil != err {
		if err != gorm.ErrRecordNotFound {
			return nil, err
		}
	}

	newPolicy := model.S3UploadPolicy{
		FileExtensions: in.FileExtensions,
		RateLimit:      []model.UploadRateLimitPolicy{},
	}

	for _, input := range in.RateLimit {
		newPolicy.RateLimit = append(newPolicy.RateLimit, model.UploadRateLimitPolicy{
			Value:    input.Value,
			Object:   input.Object,
			Interval: input.Interval,
		})
	}

	newPolicyBytes, err := json.Marshal(newPolicy)
	if nil != err {
		return nil, err
	}

	var out *configDto.VariableMutationOutcome

	if nil == policy {
		out, err = srv.bundle.configBundle.VariableService.Create(ctx, configDto.VariableCreateInput{
			BucketId: bucket.Id,
			Name:     in.ApplicationId,
			Value:    string(newPolicyBytes),
			IsLocked: scalar.NilBool(false),
		})

		if nil != err {
			return nil, err
		}
	} else {
		useless, err := policy.EqualTo(newPolicy)
		if nil != err {
			return nil, err
		} else if useless {
			return nil, util.ErrorUselessInput
		}

		out, err = srv.bundle.configBundle.VariableService.Update(ctx, configDto.VariableUpdateInput{
			Id:       policy.Id,
			Version:  in.Version,
			Value:    scalar.NilString(string(newPolicyBytes)),
			IsLocked: scalar.NilBool(false),
		})

		if nil != err {
			return nil, err
		}
	}

	return &model.S3UploadPolicy{
		Id:             out.Variable.Id,
		Version:        out.Variable.Version,
		FileExtensions: newPolicy.FileExtensions,
		RateLimit:      newPolicy.RateLimit,
	}, nil
}
