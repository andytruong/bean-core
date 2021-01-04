package model

import (
	"time"

	"bean/components/scalar"
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
