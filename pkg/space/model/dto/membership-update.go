package dto

import (
	util2 "bean/components/util"
	"bean/pkg/infra/api"
	"bean/pkg/space/model"
)

type SpaceMembershipUpdateInput struct {
	Id       string        `json:"id"`
	Version  string        `json:"version"`
	IsActive bool          `json:"isActive"`
	Language *api.Language `json:"language"`
}

type SpaceMembershipUpdateOutcome struct {
	Errors     []*util2.Error    `json:"errors"`
	Membership *model.Membership `json:"membership"`
}
