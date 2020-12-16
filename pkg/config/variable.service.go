package config

import (
	"context"
	"time"

	"gorm.io/gorm"

	"bean/components/claim"
	"bean/components/scalar"
	util2 "bean/components/util"
	"bean/pkg/config/model"
	"bean/pkg/config/model/dto"
)

type VariableService struct {
	bundle *ConfigBundle
}

func (this VariableService) access(ctx context.Context, db *gorm.DB, bucketId string, action string) (bool, error) {
	bucket, err := this.bundle.BucketService.Load(ctx, db, bucketId)
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

func (this VariableService) Load(ctx context.Context, db *gorm.DB, id string) (*model.ConfigVariable, error) {
	variable := &model.ConfigVariable{}

	err := db.First(&variable, "id = ?", id).Error
	if nil != err {
		return nil, err
	}

	if access, err := this.access(ctx, db, variable.BucketId, "read"); nil != err {
		return nil, err
	} else if !access {
		return nil, util2.ErrorAccessDenied
	}

	return variable, nil
}

func (this VariableService) Create(ctx context.Context, tx *gorm.DB, in dto.VariableCreateInput) (*dto.VariableMutationOutcome, error) {
	if access, err := this.access(ctx, tx, in.BucketId, "write"); nil != err {
		return nil, err
	} else if !access {
		return nil, util2.ErrorAccessDenied
	}

	variable := &model.ConfigVariable{
		Id:          this.bundle.id.MustULID(),
		Version:     this.bundle.id.MustULID(),
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

func (this VariableService) Update(ctx context.Context, tx *gorm.DB, in dto.VariableUpdateInput) (*dto.VariableMutationOutcome, error) {
	variable, err := this.Load(ctx, tx, in.Id)
	if nil != err {
		return nil, err
	}

	if access, err := this.access(ctx, tx, variable.BucketId, "write"); nil != err {
		return nil, err
	} else if !access {
		return nil, util2.ErrorAccessDenied
	}

	if variable.Version != in.Version {
		return nil, util2.ErrorVersionConflict
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
				return nil, util2.ErrorLocked
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
			variable.Version = this.bundle.id.MustULID()
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

func (this VariableService) Delete(ctx context.Context, tx *gorm.DB, in dto.VariableDeleteInput) (*dto.VariableMutationOutcome, error) {
	variable, err := this.Load(ctx, tx, in.Id)
	if nil != err {
		return nil, err
	} else if variable.IsLocked {
		return nil, util2.ErrorLocked
	} else {
		if access, err := this.access(ctx, tx, variable.BucketId, "delete"); nil != err {
			return nil, err
		} else if !access {
			return nil, util2.ErrorAccessDenied
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