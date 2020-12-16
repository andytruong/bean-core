package dto

import (
	util2 "bean/components/util"
	"bean/pkg/space/model"
)

type SpaceMembershipCreateInput struct {
	SpaceID          string   `json:"spaceId"`
	UserID           string   `json:"userId"`
	IsActive         bool     `json:"isActive"`
	ManagerMemberIds []string `json:"managerMemberIds"`
}

type SpaceMembershipCreateOutcome struct {
	Errors     []*util2.Error    `json:"errors"`
	Membership *model.Membership `json:"membership"`
}
