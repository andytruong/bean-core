package dto

import (
	util2 "bean/components/util"
)

type ValidationInput struct {
	HashedToken string `json:"hashedToken"`
}

type ValidationOutcome struct {
	Status bool           `json:"status"`
	Errors []*util2.Error `json:"errors"`
}
