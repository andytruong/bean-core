package dto

import (
	"bean/components/util"
)

type ValidationInput struct {
	HashedToken string `json:"hashedToken"`
}

type ValidationOutcome struct {
	Status bool          `json:"status"`
	Errors []*util.Error `json:"errors"`
}
