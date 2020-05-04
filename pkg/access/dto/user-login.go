package dto

import (
	"bean/pkg/access/model"
	"bean/pkg/util"
)

type LoginOutcome struct {
	Errors  []*util.Error  `json:"errors"`
	Session *model.Session `json:"session"`
}
