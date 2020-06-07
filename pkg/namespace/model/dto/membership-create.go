package dto

import (
	"bean/pkg/namespace/model"
	"bean/pkg/util"
)

type NamespaceMembershipCreateInput struct {
	NamespaceID      string   `json:"namespaceId"`
	UserID           string   `json:"userId"`
	IsActive         bool     `json:"isActive"`
	ManagerMemberIds []string `json:"managerMemberIds"`
}

type NamespaceMembershipCreateOutcome struct {
	Errors     []*util.Error     `json:"errors"`
	Membership *model.Membership `json:"membership"`
}
