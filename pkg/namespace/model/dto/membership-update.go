package dto

import (
	"bean/pkg/namespace/model"
	"bean/pkg/util"
)

type NamespaceMembershipUpdateInput struct {
	Id       string `json:"id"`
	Version  string `json:"version"`
	IsActive bool   `json:"isActive"`
}

type NamespaceMembershipUpdateOutcome struct {
	Errors     []*util.Error     `json:"errors"`
	Membership *model.Membership `json:"membership"`
}
