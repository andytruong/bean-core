package mailer

import (
	"context"

	"bean/pkg/integration/mailer/model/dto"
)

type MailerResolver struct {
}

func (this MailerResolver) MailerMutation(ctx context.Context) (*dto.MailerMutation, error) {
	panic("wip")
}

func (this MailerResolver) MailerQuery(ctx context.Context) (*dto.MailerQuery, error) {
	panic("wip")
}
