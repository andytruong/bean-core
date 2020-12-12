package dto

import (
	"time"
)

type MailerAuditLog struct {
	ID              string          `json:"id"`
	Account         *MailerAccount  `json:"account"`
	SpanID          string          `json:"spanId"`
	Template        *MailerTemplate `json:"template"`
	CreatedAt       time.Time       `json:"createdAt"`
	RecipientHash   string          `json:"recipientHash"`
	ContextHash     string          `json:"contextHash"`
	ErrorCode       *int            `json:"errorCode"`
	ErrorMessage    *string         `json:"errorMessage"`
	WarningMessagse *string         `json:"warningMessagse"`
}
