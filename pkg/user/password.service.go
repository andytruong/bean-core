package user

import (
	"time"

	"gorm.io/gorm"

	"bean/pkg/user/model"
	"bean/pkg/user/model/dto"
)

type PasswordService struct {
	bundle *UserBundle
}

func (this *PasswordService) create(tx *gorm.DB, user *model.User, in *dto.UserPasswordInput) error {
	if nil == in {
		return nil
	}

	pass := &model.UserPassword{
		ID:          this.bundle.id.MustULID(),
		UserId:      user.ID,
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