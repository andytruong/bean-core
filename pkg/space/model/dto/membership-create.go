package dto

import (
	"bean/pkg/space/model"
	"bean/pkg/util"
)

type SpaceMembershipCreateInput struct {
	SpaceID          string   `json:"spaceId"`
	UserID           string   `json:"userId"`
	IsActive         bool     `json:"isActive"`
	ManagerMemberIds []string `json:"managerMemberIds"`
}

type SpaceMembershipCreateOutcome struct {
	Errors     []*util.Error     `json:"errors"`
	Membership *model.Membership `json:"membership"`
}
