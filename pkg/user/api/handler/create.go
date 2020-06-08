package handler

import (
	"time"

	"github.com/jinzhu/gorm"

	"bean/pkg/user/model"
	"bean/pkg/user/model/dto"
	"bean/pkg/util"
)

type UserCreateHandler struct {
	ID *util.Identifier
}

func (this *UserCreateHandler) Handle(tx *gorm.DB, input *dto.UserCreateInput) (*dto.UserMutationOutcome, error) {
	// create base record
	user := model.User{
		ID:        this.ID.MustULID(),
		Version:   this.ID.MustULID(),
		AvatarURI: input.AvatarURI,
		IsActive:  input.IsActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
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
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			IsActive:    true,
		}

		if err := tx.Save(pass).Error; nil != err {
			return nil, err
		}
	}

	return &dto.UserMutationOutcome{User: &user}, nil
}

func (this *UserCreateHandler) createEmails(tx *gorm.DB, user *model.User, input *dto.UserCreateInput) error {
	if nil == input.Emails {
		return nil
	}

	if nil != input.Emails.Primary {
		table := "user_emails"

		if !input.Emails.Primary.Verified {
			table = "user_unverified_emails"
		}

		email := model.UserEmail{
			ID:        this.ID.MustULID(),
			UserId:    user.ID,
			Value:     input.Emails.Primary.Value.LowerCaseValue(),
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
			table := "user_emails"

			if !secondaryInput.Verified {
				table = "user_unverified_emails"
			}

			email := model.UserEmail{
				ID:        this.ID.MustULID(),
				UserId:    user.ID,
				Value:     secondaryInput.Value.LowerCaseValue(),
				IsActive:  secondaryInput.IsActive,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				IsPrimary: false,
			}

			if err := tx.Table(table).Create(&email).Error; nil != err {
				return err
			}
		}
	}

	return nil
}

func (this *UserCreateHandler) createName(tx *gorm.DB, user *model.User, input *dto.UserCreateInput) error {
	if nil != input.Name {
		name := model.UserName{
			ID:            this.ID.MustULID(),
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