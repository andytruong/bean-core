package dto

import (
	"bean/pkg/namespace/model"
	"bean/pkg/util"
	"bean/pkg/util/api"
)

type NamespaceMembershipUpdateInput struct {
	Id       string        `json:"id"`
	Version  string        `json:"version"`
	IsActive bool          `json:"isActive"`
	Language *api.Language `json:"language"`
}

type NamespaceMembershipUpdateOutcome struct {
	Errors     []*util.Error     `json:"errors"`
	Membership *model.Membership `json:"membership"`
}
