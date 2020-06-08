package handler

import (
	"time"

	"github.com/jinzhu/gorm"

	"bean/pkg/user/model"
	"bean/pkg/user/model/dto"
	"bean/pkg/util"
	"bean/pkg/util/connect"
)

type UserUpdateHandler struct {
	ID *util.Identifier
}

func (this *UserUpdateHandler) Handle(tx *gorm.DB, user *model.User, input dto.UserUpdateInput) (*dto.UserMutationOutcome, error) {
	// validate version
	if user.Version != input.Version {
		errors := util.NewErrors(util.ErrorCodeConflict, []string{"version"}, "")

		return &dto.UserMutationOutcome{
			Errors: errors,
			User:   nil,
		}, nil
	}

	if input.Values.Password != nil {
		pass := model.UserPassword{
			ID:          this.ID.MustULID(),
			UserId:      user.ID,
			Algorithm:   input.Values.Password.Algorithm,
			HashedValue: input.Values.Password.HashedValue,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			IsActive:    true,
		}

		err := tx.Save(&pass).Error
		if nil != err {
			return nil, err
		} else {
			// disable other active password
			err := tx.Table(connect.TableAccessPassword).
				Where("user_id == ?", pass.UserId).
				Where("id != ?", pass.ID).
				Update(model.UserPassword{IsActive: false}).
				Error

			if nil != err {
				return nil, err
			}
		}
	}

	// bump new version
	user.Version = this.ID.MustULID()
	if err := tx.Save(user).Error; nil != err {
		return nil, err
	}

	return &dto.UserMutationOutcome{
		Errors: nil,
		User:   user,
	}, nil
}
