package dto

import (
	"bean/components/util"
	"bean/pkg/app/model"
)

type ApplicationCreateInput struct {
	IsActive bool    `json:"isActive"`
	Title    *string `json:"title"`
}

type ApplicationUpdateInput struct {
	Id       string  `json:"id"`
	Version  string  `json:"version"`
	IsActive *bool   `json:"isActive"`
	Title    *string `json:"title"`
}

type ApplicationOutcome struct {
	App    *model.Application `json:"application"`
	Errors []*util.Error      `json:"errors"`
}
