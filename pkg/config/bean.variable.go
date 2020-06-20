package config

import (
	"context"
	"time"

	"gorm.io/gorm"

	"bean/pkg/config/model"
	"bean/pkg/config/model/dto"
	"bean/pkg/util"
	"bean/pkg/util/connect"
)

type CoreVariable struct {
	bean *ConfigBean
}

func (this CoreVariable) access(ctx context.Context, db *gorm.DB, bucketId string, action string) (bool, error) {
	bucket, err := this.bean.CoreBucket.Load(ctx, db, bucketId)
	if nil != err {
		return false, err
	}

	if nil == bucket {
		return false, nil
	}

	actor := util.CxtKeyClaims.Actor(ctx)
	isOwner := (nil != actor) && actor.UserId() == bucket.HostId
	isMember := (nil != actor) && actor.NamespaceId() == bucket.HostId

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

func (this CoreVariable) Load(ctx context.Context, db *gorm.DB, id string) (*model.ConfigVariable, error) {
	variable := &model.ConfigVariable{}

	err := db.
		Table(connect.TableConfigVariable).
		First(&variable, "id = ?", id).
		Error

	if nil != err {
		return nil, err
	}

	if access, err := this.access(ctx, db, variable.BucketId, "read"); nil != err {
		return nil, err
	} else if !access {
		return nil, util.ErrorAccessDenied
	}

	return variable, nil
}

func (this CoreVariable) Create(ctx context.Context, tx *gorm.DB, in dto.VariableCreateInput) (*dto.VariableMutationOutcome, error) {
	if access, err := this.access(ctx, tx, in.BucketId, "write"); nil != err {
		return nil, err
	} else if !access {
		return nil, util.ErrorAccessDenied
	}

	variable := &model.ConfigVariable{
		Id:          this.bean.id.MustULID(),
		Version:     this.bean.id.MustULID(),
		BucketId:    in.BucketId,
		Name:        in.Name,
		Description: in.Description,
		Value:       in.Value,
		IsLocked:    util.NotNilBool(in.IsLocked, false),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := tx.Table(connect.TableConfigVariable).Create(&variable).Error

	if nil != err {
		return nil, err
	} else {
		return &dto.VariableMutationOutcome{
			Errors:   nil,
			Variable: variable,
		}, nil
	}
}

func (this CoreVariable) Update(ctx context.Context, tx *gorm.DB, in dto.VariableUpdateInput) (*dto.VariableMutationOutcome, error) {
	variable, err := this.Load(ctx, tx, in.Id)
	if nil != err {
		return nil, err
	}

	if access, err := this.access(ctx, tx, variable.BucketId, "write"); nil != err {
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
			variable.Version = this.bean.id.MustULID()
			err = tx.Table(connect.TableConfigVariable).
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

func (this CoreVariable) Delete(ctx context.Context, tx *gorm.DB, in dto.VariableDeleteInput) (*dto.VariableMutationOutcome, error) {
	variable, err := this.Load(ctx, tx, in.Id)
	if nil != err {
		return nil, err
	} else if variable.IsLocked {
		return nil, util.ErrorLocked
	} else {
		if access, err := this.access(ctx, tx, variable.BucketId, "delete"); nil != err {
			return nil, err
		} else if !access {
			return nil, util.ErrorAccessDenied
		}

		err := tx.
			Table(connect.TableConfigVariable).
			Delete(variable, "id = ?", variable.Id).
			Error
		if nil != err {
			return nil, err
		}
	}

	return &dto.VariableMutationOutcome{
		Errors:   nil,
		Variable: variable,
	}, nil
}
