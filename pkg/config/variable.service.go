package config

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"bean/components/claim"
	"bean/components/connect"
	"bean/components/scalar"
	"bean/components/util"
	"bean/pkg/config/model"
	"bean/pkg/config/model/dto"
)

type VariableService struct {
	bundle *Bundle
}

func (srv VariableService) canAccess(ctx context.Context, bucketId string, action string) (bool, error) {
	bucket, err := srv.bundle.BucketService.Load(ctx, dto.BucketKey{Id: bucketId})
	if nil != err {
		return false, err
	}

	if nil == bucket {
		return false, nil
	}

	return srv.canAccessBucket(ctx, bucket, action)
}

func (srv VariableService) canAccessBucket(ctx context.Context, bucket *model.ConfigBucket, action string) (bool, error) {
	claims := claim.ContextToPayload(ctx)
	isOwner := (nil != claims) && claims.UserId() == bucket.HostId
	isMember := (nil != claims) && claims.SpaceId() == bucket.HostId

	switch action {
	case "read":
		return bucket.Access.CanRead(isOwner, isMember), nil

	case "write":
		return bucket.Access.CanWrite(isOwner, isMember), nil

	case "delete":
		return bucket.Access.CanDelete(isOwner, isMember), nil
	}

	return false, nil
}

func (srv VariableService) Load(ctx context.Context, key dto.VariableKey) (*model.ConfigVariable, error) {
	if (key.Id == "") && (key.BucketId == "" || key.Name == "") {
		return nil, errors.New("invalid load key")
	}

	db := connect.DB(ctx)
	variable := &model.ConfigVariable{}
	if key.Id != "" {
		err := db.Take(&variable, "id = ?", key.Id).Error
		if nil != err {
			return nil, err
		}
	} else if key.BucketId != "" && key.Name != "" {
		err := db.Take(&variable, "bucket_id = ? AND name = ?", key.BucketId, key.Name).Error
		if nil != err {
			return nil, err
		}
	}

	if access, err := srv.canAccess(ctx, variable.BucketId, "read"); nil != err {
		return nil, err
	} else if !access {
		return nil, util.ErrorAccessDenied
	}

	return variable, nil
}

func (srv VariableService) Create(ctx context.Context, in dto.VariableCreateInput) (*dto.VariableMutationOutcome, error) {
	bucket, err := srv.bundle.BucketService.Load(ctx, dto.BucketKey{Id: in.BucketId})
	if nil != err {
		return nil, err
	} else if nil == bucket {
		return nil, errors.New("bucket not found")
	}

	if access, err := srv.canAccessBucket(ctx, bucket, "write"); nil != err {
		return nil, err
	} else if !access {
		return nil, util.ErrorAccessDenied
	}

	if reasons, err := bucket.Validate(ctx, in.Value); nil != err {
		return nil, err
	} else if len(reasons) > 0 {
		errList := []util.Error{}
		for _, reason := range reasons {
			err := util.NewError(util.ErrorCodeInput, []string{"VariableCreateInput.Value"}, reason)
			errList = append(errList, err)
		}

		return &dto.VariableMutationOutcome{Errors: errList}, nil
	}

	variable := &model.ConfigVariable{
		Id:          srv.bundle.idr.ULID(),
		Version:     srv.bundle.idr.ULID(),
		BucketId:    in.BucketId,
		Name:        in.Name,
		Description: in.Description,
		Value:       in.Value,
		IsLocked:    scalar.NotNilBool(in.IsLocked, false),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err = connect.DB(ctx).Create(&variable).Error
	if nil != err {
		return nil, err
	} else {
		return &dto.VariableMutationOutcome{Errors: nil, Variable: variable}, nil
	}
}

func (srv VariableService) Update(ctx context.Context, in dto.VariableUpdateInput) (*dto.VariableMutationOutcome, error) {
	tx := connect.DB(ctx)
	variable, err := srv.Load(ctx, dto.VariableKey{Id: in.Id})
	if nil != err {
		return nil, err
	}

	if access, err := srv.canAccess(ctx, variable.BucketId, "write"); nil != err {
		return nil, err
	} else if !access {
		return nil, util.ErrorAccessDenied
	}

	if variable.Version != in.Version {
		return nil, util.ErrorVersionConflict
	} else {
		changed := false

		if nil != in.Description {
			if variable.Description != in.Description {
				changed = true
				variable.Description = in.Description
			}
		}

		if in.Value != nil {
			if variable.Value != *in.Value {
				changed = true
				variable.Value = *in.Value
			}
		}

		if variable.IsLocked {
			if changed {
				return nil, util.ErrorLocked
			}
		}

		if nil != in.IsLocked {
			if variable.IsLocked != *in.IsLocked {
				changed = true
				variable.IsLocked = *in.IsLocked
			}
		}

		if changed {
			version := variable.Version
			variable.Version = srv.bundle.idr.ULID()
			err = tx.
				Where("version = ?", version).
				Save(&variable).
				Error
			if nil != err {
				return nil, err
			}
		}
	}

	return &dto.VariableMutationOutcome{
		Errors:   nil,
		Variable: variable,
	}, nil
}

func (srv VariableService) Delete(ctx context.Context, in dto.VariableDeleteInput) (*dto.VariableMutationOutcome, error) {
	tx := connect.DB(ctx)
	variable, err := srv.Load(ctx, dto.VariableKey{Id: in.Id})
	if nil != err {
		return nil, err
	} else if variable.IsLocked {
		return nil, util.ErrorLocked
	} else {
		if access, err := srv.canAccess(ctx, variable.BucketId, "delete"); nil != err {
			return nil, err
		} else if !access {
			return nil, util.ErrorAccessDenied
		}

		err := tx.Delete(variable, "id = ?", variable.Id).Error
		if nil != err {
			return nil, err
		}
	}

	return &dto.VariableMutationOutcome{
		Errors:   nil,
		Variable: variable,
	}, nil
}
