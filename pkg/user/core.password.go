package user

import (
	"time"

	"gorm.io/gorm"

	"bean/pkg/user/model"
	"bean/pkg/user/model/dto"
)

type CorePassword struct {
	bean *UserBean
}

func (this *CorePassword) create(tx *gorm.DB, user *model.User, in *dto.UserPasswordInput) error {
	if nil == in {
		return nil
	}

	pass := &model.UserPassword{
		ID:          this.bean.id.MustULID(),
		UserId:      user.ID,
		Algorithm:   in.Algorithm,
		HashedValue: in.HashedValue,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		IsActive:    true,
	}

	{
		err := tx.Create(pass).Error
		if nil != err {
			return err
		}
	}

	// set other passwords to inactive
	return tx.
		Where("user_id == ?", pass.UserId).
		Where("id != ?", pass.ID).
		Updates(model.UserPassword{IsActive: false}).
		Error
}
