package service

import (
	"time"

	"github.com/jinzhu/gorm"

	"bean/pkg/user/dto"
	"bean/pkg/user/model"
	"bean/pkg/util"
)

type (
	UserCreateAPI struct {
		ID *util.Identifier
	}
)

func (this *UserCreateAPI) Create(tx *gorm.DB, input *dto.UserCreateInput) (*dto.UserCreateOutcome, error) {
	// create base record
	user := model.User{
		AvatarURI: input.AvatarURI,
		IsActive:  input.IsActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Generate user Identifier.
	if id, err := this.ID.ULID(); nil != err {
		return nil, err
	} else if version, err := this.ID.ULID(); nil != err {
		return nil, err
	} else {
		user.ID = id
		user.Version = version
	}

	if err := tx.Create(&user).Error; nil != err {
		return nil, err
	}

	// create emails
	if err := this.createEmails(tx, &user, input); nil != err {
		return nil, err
	}

	// save name object
	if err := this.createName(tx, &user, input); nil != err {
		return nil, err
	}

	// save password
	if nil != input.Password {
		id, err := this.ID.UUID()
		if nil != err {
			return nil, err
		}

		pass := model.UserPassword{
			ID:          id,
			UserId:      user.ID,
			Algorithm:   input.Password.Algorithm,
			HashedValue: input.Password.HashedValue,
			CreatedAt:   time.Time{},
			UpdatedAt:   time.Time{},
			IsActive:    false,
		}

		tx.Save(pass)
	}

	return &dto.UserCreateOutcome{User: &user}, nil
}

func (this *UserCreateAPI) createEmails(tx *gorm.DB, user *model.User, input *dto.UserCreateInput) error {
	if nil != input.Emails {
		if nil != input.Emails.Primary {
			table := "user_email"
			id, _ := this.ID.ULID()

			if !input.Emails.Primary.Verified {
				table = "user_email_unverified"
			}

			email := model.UserEmail{
				ID:        id,
				UserId:    user.ID,
				Value:     input.Emails.Primary.Value,
				IsActive:  input.Emails.Primary.IsActive,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				IsPrimary: true,
			}

			if err := tx.Table(table).Create(&email).Error; nil != err {
				return err
			}
		}

		if nil != input.Emails.Secondary {
			for _, secondaryInput := range input.Emails.Secondary {
				table := "user_email"
				id, _ := this.ID.ULID()

				if !secondaryInput.Verified {
					table = "user_email_unverified"
				}

				email := model.UserEmail{
					ID:        id,
					UserId:    user.ID,
					Value:     secondaryInput.Value,
					IsActive:  secondaryInput.IsActive,
					IsPrimary: false,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}

				if err := tx.Table(table).Create(&email).Error; nil != err {
					return err
				}
			}
		}
	}

	return nil
}

func (this *UserCreateAPI) createName(tx *gorm.DB, user *model.User, input *dto.UserCreateInput) error {
	if nil != input.Name {
		id, _ := this.ID.ULID()
		name := model.UserName{
			ID:            id,
			UserId:        user.ID,
			FirstName:     input.Name.FirstName,
			LastName:      input.Name.LastName,
			PreferredName: input.Name.PreferredName,
		}

		if err := tx.Create(name).Error; nil != err {
			return err
		}
	}

	return nil
}
