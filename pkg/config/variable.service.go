package config

import (
	"context"
	"time"

	"bean/components/claim"
	"bean/components/connect"
	"bean/components/scalar"
	"bean/components/util"
	"bean/pkg/config/model"
	"bean/pkg/config/model/dto"
)

type VariableService struct {
	bundle *ConfigBundle
}

func (service VariableService) access(ctx context.Context, bucketId string, action string) (bool, error) {
	db := connect.ContextToDB(ctx)
	bucket, err := service.bundle.BucketService.Load(connect.DBToContext(ctx, db), bucketId)
	if nil != err {
		return false, err
	}

	if nil == bucket {
		return false, nil
	}

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

func (service VariableService) Load(ctx context.Context, id string) (*model.ConfigVariable, error) {
	db := connect.ContextToDB(ctx)
	variable := &model.ConfigVariable{}

	err := db.First(&variable, "id = ?", id).Error
	if nil != err {
		return nil, err
	}

	if access, err := service.access(ctx, variable.BucketId, "read"); nil != err {
		return nil, err
	} else if !access {
		return nil, util.ErrorAccessDenied
	}

	return variable, nil
}

func (service VariableService) Create(ctx context.Context, in dto.VariableCreateInput) (*dto.VariableMutationOutcome, error) {
	tx := connect.ContextToDB(ctx)
	if access, err := service.access(ctx, in.BucketId, "write"); nil != err {
		return nil, err
	} else if !access {
		return nil, util.ErrorAccessDenied
	}

	variable := &model.ConfigVariable{
		Id:          service.bundle.idr.MustULID(),
		Version:     service.bundle.idr.MustULID(),
		BucketId:    in.BucketId,
		Name:        in.Name,
		Description: in.Description,
		Value:       in.Value,
		IsLocked:    scalar.NotNilBool(in.IsLocked, false),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := tx.Create(&variable).Error
	if nil != err {
		return nil, err
	} else {
		return &dto.VariableMutationOutcome{
			Errors:   nil,
			Variable: variable,
		}, nil
	}
}

func (service VariableService) Update(ctx context.Context, in dto.VariableUpdateInput) (*dto.VariableMutationOutcome, error) {
	tx := connect.ContextToDB(ctx)
	variable, err := service.Load(ctx, in.Id)
	if nil != err {
		return nil, err
	}

	if access, err := service.access(ctx, variable.BucketId, "write"); nil != err {
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
			variable.Version = service.bundle.idr.MustULID()
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

func (service VariableService) Delete(ctx context.Context, in dto.VariableDeleteInput) (*dto.VariableMutationOutcome, error) {
	tx := connect.ContextToDB(ctx)
	variable, err := service.Load(ctx, in.Id)
	if nil != err {
		return nil, err
	} else if variable.IsLocked {
		return nil, util.ErrorLocked
	} else {
		if access, err := service.access(ctx, variable.BucketId, "delete"); nil != err {
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
