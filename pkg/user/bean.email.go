package user

import (
	"context"

	"bean/pkg/user/model"
)

type UserBeanEmail struct {
	bean *UserBean
}

// TODO: need a better resolver, we not always load secondary emails.
//       see: https://gqlgen.com/reference/field-collection/
func (this UserBeanEmail) List(ctx context.Context, user *model.User) (*model.UserEmails, error) {
	emails := &model.UserEmails{}

	var rows []*model.UserEmail
	err := this.bean.db.
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
