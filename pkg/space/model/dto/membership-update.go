package dto

import (
	"bean/pkg/space/model"
	"bean/pkg/util"
	"bean/pkg/util/api"
)

type SpaceMembershipUpdateInput struct {
	Id       string        `json:"id"`
	Version  string        `json:"version"`
	IsActive bool          `json:"isActive"`
	Language *api.Language `json:"language"`
}

type SpaceMembershipUpdateOutcome struct {
	Errors     []*util.Error     `json:"errors"`
	Membership *model.Membership `json:"membership"`
}
