package handler

import (
	"context"

	"github.com/jinzhu/gorm"

	"bean/pkg/user/model"
)

type EmailQueryHandler struct {
	DB *gorm.DB
}

// TODO: need a better resolver, we not always load secondary emails.
func (this EmailQueryHandler) Emails(ctx context.Context, user *model.User) (*model.UserEmails, error) {
	emails := &model.UserEmails{}

	var rows []*model.UserEmail
	err := this.DB.
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
