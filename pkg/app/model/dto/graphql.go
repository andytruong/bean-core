package dto

import (
	"bean/components/util"
	"bean/pkg/app/model"
)

type (
	ApplicationQuery    struct{}
	ApplicationMutation struct{}

	ApplicationOutcome struct {
		App    *model.Application `json:"application"`
		Errors []*util.Error      `json:"errors"`
	}

	ApplicationCreateInput struct {
		IsActive bool    `json:"isActive"`
		Title    *string `json:"title"`
	}

	ApplicationUpdateInput struct {
		Id       string  `json:"id"`
		Version  string  `json:"version"`
		IsActive *bool   `json:"isActive"`
		Title    *string `json:"title"`
	}

	ApplicationDeleteInput struct {
		Id      string `json:"id"`
		Version string `json:"version"`
	}
)
