package user

import (
	"context"
	"time"
	
	"bean/components/connect"
	"bean/components/scalar"
	"bean/pkg/user/model"
	"bean/pkg/user/model/dto"
)

type EmailService struct {
	bundle *Bundle
}

func (srv EmailService) Load(ctx context.Context, email scalar.EmailAddress) (*model.UserEmail, error) {
	entity := &model.UserEmail{}
	err := connect.ContextToDB(ctx).First(entity, "value = ?", email).Error
	
	if nil != err {
		return nil, err
	}
	
	return entity, nil
}

func (srv EmailService) CreateBulk(ctx context.Context, user *model.User, in *dto.UserEmailsInput) error {
	if nil == in {
		return nil
	}
	
	if nil != in.Primary {
		err := srv.bundle.EmailService.Create(ctx, user, *in.Primary, true)
		if nil != err {
			return err
		}
	}
	
	if nil != in.Secondary {
		for _, secondaryInput := range in.Secondary {
			err := srv.bundle.EmailService.Create(ctx, user, *secondaryInput, false)
			if nil != err {
				return err
			}
		}
	}
	
	return nil
}

func (srv EmailService) Create(ctx context.Context, user *model.User, in dto.UserEmailInput, isPrimary bool) error {
	table := connect.TableUserEmail
	if !in.Verified {
		table = connect.TableUserEmailUnverified
	}
	
	email := model.UserEmail{
		ID:        srv.bundle.idr.ULID(),
		UserId:    user.ID,
		Value:     in.Value.LowerCaseValue(),
		IsActive:  in.Verified,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		IsPrimary: isPrimary,
	}
	
	return connect.ContextToDB(ctx).Table(table).Create(&email).Error
}

// TODO: need a better resolver, we not always load secondary emails.
//       see: https://gqlgen.com/reference/field-collection/
func (srv EmailService) List(ctx context.Context, user *model.User) (*model.UserEmails, error) {
	emails := &model.UserEmails{}
	
	var rows []*model.UserEmail
	err := connect.ContextToDB(ctx).
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
