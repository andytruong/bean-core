package api

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
	}

	// Generate user Identifier.
	if id, err := this.ID.Hash("User", time.Now()); nil != err {
		return nil, err
	} else {
		user.ID = id
	}

	if err := tx.Create(user).Error; nil != err {
		return nil, err
	}

	// create emails
	if err := this.createEmails(tx, &user, input); nil != err {
		tx.Rollback()

		return nil, err
	}

	// save name object
	if err := this.createName(tx, &user, input); nil != err {
		tx.Rollback()

		return nil, err
	}

	// outcome
	outcome := dto.UserCreateOutcome{
		User:   &user,
		Errors: nil,
	}

	return &outcome, nil
}

func (this *UserCreateAPI) createEmails(tx *gorm.DB, user *model.User, input *dto.UserCreateInput) error {
	if nil != input.Emails {
		if nil != input.Emails.Primary {
			id, _ := this.ID.ULID()
			email := model.UserEmail{
				ID:        id,
				UserId:    user.ID,
				Verified:  input.Emails.Primary.Verified,
				Value:     input.Emails.Primary.Value,
				IsActive:  input.Emails.Primary.IsActive,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			if err := tx.Create(email).Error; nil != err {
				return err
			}
		}

		if nil != input.Emails.Secondary {
			for _, secondaryInput := range input.Emails.Secondary {
				id, _ := this.ID.ULID()
				email := model.UserEmail{
					ID:        id,
					UserId:    user.ID,
					Verified:  secondaryInput.Verified,
					Value:     secondaryInput.Value,
					IsActive:  secondaryInput.IsActive,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}

				if err := tx.Create(email).Error; nil != err {
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
