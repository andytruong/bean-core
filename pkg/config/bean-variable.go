package config

import (
	"context"
	"time"

	"github.com/jinzhu/gorm"

	"bean/pkg/config/model"
	"bean/pkg/config/model/dto"
	"bean/pkg/util"
	"bean/pkg/util/connect"
)

type ConfigVariableBean struct {
	bean *ConfigBean
}

func (this ConfigVariableBean) Load(ctx context.Context, db *gorm.DB, id string) (*model.ConfigVariable, error) {
	variable := &model.ConfigVariable{}

	err := db.Table(connect.TableConfigVariable).First(&variable, "id = ?", id).Error
	if nil != err {
		return nil, err
	}

	return variable, nil
}

func (this ConfigVariableBean) Create(ctx context.Context, tx *gorm.DB, input dto.VariableCreateInput) (*dto.VariableMutationOutcome, error) {
	variable := &model.ConfigVariable{
		Id:          this.bean.id.MustULID(),
		Version:     this.bean.id.MustULID(),
		BucketId:    input.BucketId,
		Name:        input.Name,
		Description: input.Description,
		Value:       input.Value,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := tx.Table(connect.TableConfigVariable).Save(&variable).Error

	if nil != err {
		return nil, err
	} else {
		return &dto.VariableMutationOutcome{
			Errors:   nil,
			Variable: variable,
		}, nil
	}
}

func (this ConfigVariableBean) Update(ctx context.Context, tx *gorm.DB, input dto.VariableUpdateInput) (*dto.VariableMutationOutcome, error) {
	variable, err := this.Load(ctx, tx, input.Id)
	if nil != err {
		return nil, err
	}

	if variable.Version != input.Version {
		return nil, util.ErrorVersionConflict
	} else {
		changed := false

		if nil != input.Description {
			if variable.Description != input.Description {
				changed = true
				variable.Description = input.Description
			}
		}

		if input.Value != nil {
			if variable.Value != *input.Value {
				changed = true
				variable.Value = *input.Value
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

func (this ConfigVariableBean) Delete(ctx context.Context, tx *gorm.DB, input dto.VariableDeleteInput) (*dto.VariableMutationOutcome, error) {
	variable, err := this.Load(ctx, tx, input.Id)
	if nil != err {
		return nil, err
	} else {
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
