package user

import (
	"context"
	"time"

	"gorm.io/gorm"

	"bean/pkg/user/model"
	"bean/pkg/user/model/dto"
	"bean/pkg/util/connect"
)

type EmailService struct {
	bean *UserBundle
}

func (this EmailService) CreateBulk(tx *gorm.DB, user *model.User, in *dto.UserEmailsInput) error {
	if nil == in {
		return nil
	}

	if nil != in.Primary {
		err := this.bean.EmailService.Create(tx, user, *in.Primary, true)
		if nil != err {
			return err
		}
	}

	if nil != in.Secondary {
		for _, secondaryInput := range in.Secondary {
			err := this.bean.EmailService.Create(tx, user, *secondaryInput, false)
			if nil != err {
				return err
			}
		}
	}

	return nil
}

func (this EmailService) Create(tx *gorm.DB, user *model.User, in dto.UserEmailInput, isPrimary bool) error {
	table := connect.TableUserEmail
	if !in.Verified {
		table = connect.TableUserEmailUnverified
	}

	email := model.UserEmail{
		ID:        this.bean.id.MustULID(),
		UserId:    user.ID,
		Value:     in.Value.LowerCaseValue(),
		IsActive:  in.Verified,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		IsPrimary: isPrimary,
	}

	err := tx.Table(table).Create(&email).Error

	return err
}

// TODO: need a better resolver, we not always load secondary emails.
//       see: https://gqlgen.com/reference/field-collection/
func (this EmailService) List(ctx context.Context, user *model.User) (*model.UserEmails, error) {
	emails := &model.UserEmails{}

	var rows []*model.UserEmail
	err := this.bean.db.
		WithContext(ctx).
		Raw(`
			     SELECT *, 1 AS is_verified FROM user_emails            WHERE user_id = ?
		   UNION SELECT *, 0 AS is_verified FROM user_unverified_emails WHERE user_id = ?
		`, user.ID, user.ID).
		Find(&rows).
		Error

	if nil != err {
		return nil, err
	} else {
		for _, row := range rows {
			if row.IsPrimary {
				emails.Primary = row
			} else {
				emails.Secondary = append(emails.Secondary, row)
			}
		}
	}

	return emails, nil
}