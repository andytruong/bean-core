package dto

import (
	util2 "bean/components/util"
)

type LogoutInput struct {
	HashedToken string `json:"hashedToken"`
}

type SessionArchiveOutcome struct {
	Errors []*util2.Error `json:"errors"`
	Result bool           `json:"result"`
}
