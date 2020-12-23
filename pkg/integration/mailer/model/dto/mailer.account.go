package dto

import (
	"time"

	"bean/components/scalar"
	"bean/components/util"
	"bean/pkg/infra/api"
	"bean/pkg/space/model"
)

type MailerAccount struct {
	ID            string                   `json:"id"`
	Version       string                   `json:"version"`
	Space         *model.Space             `json:"space"`
	Status        MailerAccountStatus      `json:"status"`
	CreatedAt     time.Time                `json:"createdAt"`
	UpdatedAt     time.Time                `json:"updatedAt"`
	DeletedAt     *time.Time               `json:"deletedAt"`
	Sender        *MailerSender            `json:"sender"`
	ConnectionURL string                   `json:"connectionUrl"`
	Attachment    *MailerAccountAttachment `json:"attachment"`
}

type MailerSender struct {
	Name  string              `json:"name"`
	Email scalar.EmailAddress `json:"email"`
}

type MailerAccountAttachment struct {
	SizeLimit     *int           `json:"sizeLimit"`
	SizeLimitEach *int           `json:"sizeLimitEach"`
	FileTypes     []api.FileType `json:"fileTypes"`
}

type MailerAccountAttachmentInput struct {
	SizeLimit     *int           `json:"sizeLimit"`
	SizeLimitEach *int           `json:"sizeLimitEach"`
	FileTypes     []api.FileType `json:"fileTypes"`
}

type MailerAccountConnection struct {
	PageInfo *MailerAccountPageInfo `json:"pageInfo"`
	Edges    []*MailerAccountEdge   `json:"edges"`
}

type MailerAccountCreateInput struct {
	SpaceID       string                        `json:"spaceId"`
	IsActive      bool                          `json:"isActive"`
	Sender        *MailerAccountSenderInput     `json:"sender"`
	ConnectionURL string                        `json:"connectionUrl"`
	Attachment    *MailerAccountAttachmentInput `json:"attachment"`
}

type MailerAccountEdge struct {
	Cursor string         `json:"cursor"`
	Node   *MailerAccount `json:"node"`
}

type MailerAccountMutation struct {
}

type MailerAccountMutationOutcome struct {
	Account *MailerAccount `json:"account"`
	Errors  []*util.Error  `json:"errors"`
}

type MailerAccountPageInfo struct {
	EndCursor   *string `json:"endCursor"`
	HasNextPage bool    `json:"hasNextPage"`
	StartCursor *string `json:"startCursor"`
}

type MailerAccountSenderInput struct {
	Name  string              `json:"name"`
	Email scalar.EmailAddress `json:"email"`
}

type MailerAccountUpdateAttachmentInput struct {
	SizeLimit     *int           `json:"sizeLimit"`
	SizeLimitEach *int           `json:"sizeLimitEach"`
	FileTypes     []api.FileType `json:"fileTypes"`
}

type MailerAccountUpdateInput struct {
	ID      string                         `json:"id"`
	Version string                         `json:"version"`
	Values  *MailerAccountUpdateValueInput `json:"values"`
}

type MailerAccountUpdateSenderInput struct {
	Name  *string              `json:"name"`
	Email *scalar.EmailAddress `json:"email"`
}

type MailerAccountUpdateValueInput struct {
	Status        *MailerAccountStatus            `json:"status"`
	ConnectionURL *string                         `json:"connectionUrl"`
	Sender        *MailerAccountUpdateSenderInput `json:"sender"`
	Attachment    *MailerAccountUpdateSenderInput `json:"attachment"`
}
