package dto

import (
	"time"

	"bean/pkg/space/model"
	mUser "bean/pkg/user/model"
	"bean/pkg/util/api"
)

type MailerTemplate struct {
	ID        string                 `json:"id"`
	Version   string                 `json:"version"`
	Space     *model.Space           `json:"space"`
	IsActive  bool                   `json:"isActive"`
	CreatedAt time.Time              `json:"createdAt"`
	UpdatedAt time.Time              `json:"updatedAt"`
	DeletedAt *time.Time             `json:"deletedAt"`
	Message   *MailerTemplateMessage `json:"message"`
}

type MailerTemplateEvent struct {
	ID       string                  `json:"id"`
	Template *MailerTemplate         `json:"template"`
	User     *mUser.User             `json:"user"`
	Key      *MailerTemplateEventKey `json:"key"`
	Payload  string                  `json:"payload"`
}

type MailerTemplateMessage struct {
	Title    string       `json:"title"`
	Language api.Language `json:"language"`
	BodyHTML string       `json:"bodyHTML"`
	BodyText *string      `json:"bodyText"`
}

type MailerTemplateMutation struct {
}
