package dto

import (
	"bean/components/scalar"
	"bean/components/util"
	"bean/pkg/infra/api"
	"bean/pkg/integration/mailer/model"
)

// Query
type (
	MailerQuery struct {
		Account *MailerQueryAccount `json:"account"`
	}

	MailerQueryAccount struct {
	}

	MailerAccountConnection struct {
		PageInfo *MailerAccountPageInfo `json:"pageInfo"`
		Edges    []*MailerAccountEdge   `json:"edges"`
	}

	MailerAccountPageInfo struct {
		EndCursor   *string `json:"endCursor"`
		HasNextPage bool    `json:"hasNextPage"`
		StartCursor *string `json:"startCursor"`
	}

	MailerAccountEdge struct {
		Cursor string               `json:"cursor"`
		Node   *model.MailerAccount `json:"node"`
	}
)

// Mutation
type (
	MailerMutation struct {
		Account  *MailerAccountMutation  `json:"account"`
		Template *MailerTemplateMutation `json:"template"`
	}

	MailerAccountMutation struct {
	}

	MailerTemplateMutation struct {
	}

	MailerAccountMutationOutcome struct {
		Account *model.MailerAccount `json:"account"`
		Errors  []*util.Error        `json:"errors"`
	}

	// account.create
	MailerAccountCreateInput struct {
		SpaceID       string                        `json:"spaceId"`
		IsActive      bool                          `json:"isActive"`
		ConnectionURL string                        `json:"connectionUrl"`
		Sender        *MailerAccountSenderInput     `json:"sender"`
		Attachment    *MailerAccountAttachmentInput `json:"attachment"`
	}

	MailerAccountSenderInput struct {
		Name  string              `json:"name"`
		Email scalar.EmailAddress `json:"email"`
	}

	MailerAccountAttachmentInput struct {
		SizeLimit     *int           `json:"sizeLimit"`
		SizeLimitEach *int           `json:"sizeLimitEach"`
		FileTypes     []api.FileType `json:"fileTypes"`
	}

	// account.update
	MailerAccountUpdateAttachmentInput struct {
		SizeLimit     *int           `json:"sizeLimit"`
		SizeLimitEach *int           `json:"sizeLimitEach"`
		FileTypes     []api.FileType `json:"fileTypes"`
	}

	MailerAccountUpdateInput struct {
		ID      string                         `json:"id"`
		Version string                         `json:"version"`
		Values  *MailerAccountUpdateValueInput `json:"values"`
	}

	MailerAccountUpdateSenderInput struct {
		Name  *string              `json:"name"`
		Email *scalar.EmailAddress `json:"email"`
	}

	MailerAccountUpdateValueInput struct {
		Status        *model.MailerAccountStatus      `json:"status"`
		ConnectionURL *string                         `json:"connectionUrl"`
		Sender        *MailerAccountUpdateSenderInput `json:"sender"`
		Attachment    *MailerAccountUpdateSenderInput `json:"attachment"`
	}
)
