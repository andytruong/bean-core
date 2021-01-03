package mailer

import (
	"context"

	"bean/pkg/integration/mailer/model/dto"
)

func newResoler(bundle *Bundle) map[string]interface{} {
	return map[string]interface{}{
		"Query": map[string]interface{}{
			"MailerQuery": func(ctx context.Context) (*dto.MailerQuery, error) {
				return &dto.MailerQuery{}, nil
			},
		},
		"Mutation": map[string]interface{}{
			"MailerMutation": func(ctx context.Context) (*dto.MailerMutation, error) {
				return &dto.MailerMutation{}, nil
			},
		},
		"MailerQuery": map[string]interface{}{},
		"MailerMutation": map[string]interface{}{
			"Account": func(ctx context.Context) (*dto.MailerAccountMutation, error) {
				return &dto.MailerAccountMutation{}, nil
			},
			"Template": func(ctx context.Context) (*dto.MailerTemplateMutation, error) {
				return &dto.MailerTemplateMutation{}, nil
			},
		},
		"MailerAccountQuery": map[string]interface{}{
			"Get": func(ctx context.Context, obj *dto.MailerQueryAccount, id string) (*dto.MailerAccount, error) {
				panic("no implementation")
			},
			"GetMultiple": func(ctx context.Context, obj *dto.MailerQueryAccount, first int, after *string) ([]*dto.MailerAccount, error) {
				panic("no implementation")
			},
		},
		"MailerAccountMutation": map[string]interface{}{
			"Create": func(ctx context.Context, obj *dto.MailerAccountMutation, input dto.MailerAccountCreateInput) (*dto.MailerAccountMutationOutcome, error) {
				panic("no implementation")
			},
			"Update": func(ctx context.Context, obj *dto.MailerAccountMutation, input dto.MailerAccountUpdateInput) (*dto.MailerAccountMutationOutcome, error) {
				panic("no implementation")
			},
			"Verify": func(ctx context.Context, obj *dto.MailerAccountMutation, id string, version string) (*dto.MailerAccountMutationOutcome, error) {
				panic("no implementation")
			},
		},
		"MailerTemplateMutation": map[string]interface{}{},
	}
}
