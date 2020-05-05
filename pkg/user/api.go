package user

import (
	"context"
	"time"

	"github.com/jinzhu/gorm"

	"bean/pkg/user/dto"
	"bean/pkg/user/model"
	"bean/pkg/util"
)

type (
	UserMutationResolver struct {
		db *gorm.DB
		id *util.Identifier
	}
)

// TODO: Work in progress
func (this *UserMutationResolver) UserCreate(ctx context.Context, input *dto.UserCreateInput) (*dto.UserCreateOutcome, error) {
	// TODO: validate email address
	// TODO: validate avatar URI

	// create base record
	user := model.User{
		AvatarURI: input.AvatarURI,
		IsActive:  input.IsActive,
	}

	// Generate user Identifier.
	if id, err := this.id.Hash("User", time.Now()); nil != err {
		return nil, err
	} else {
		user.ID = id
	}

	if err := this.db.Create(user).Error; nil != err {
		return nil, err
	}

	// create emails
	if nil != input.Emails {
		if nil != input.Emails.Primary {
			id, _ := this.id.ULID()
			email := model.UserEmail{
				ID:        id,
				UserId:    user.ID,
				Verified:  input.Emails.Primary.Verified,
				Value:     input.Emails.Primary.Value,
				IsActive:  input.Emails.Primary.IsActive,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			this.db.Create(email)
		}

		if nil != input.Emails.Secondary {
			for _, secondaryInput := range input.Emails.Secondary {
				id, _ := this.id.ULID()
				email := model.UserEmail{
					ID:        id,
					UserId:    user.ID,
					Verified:  secondaryInput.Verified,
					Value:     secondaryInput.Value,
					IsActive:  secondaryInput.IsActive,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}

				this.db.Create(email)
			}
		}
	}

	// save name object
	// outcome

	panic("not implemented")
}
