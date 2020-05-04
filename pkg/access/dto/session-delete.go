package dto

import (
	"bean/pkg/util"
)

type LogoutInput struct {
	HashedToken string `json:"hashedToken"`
}

type LogoutPayload struct {
	Errors []*util.Error `json:"errors"`
	Result *bool         `json:"result"`
}
