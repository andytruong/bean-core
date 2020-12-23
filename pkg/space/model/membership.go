package model

import (
	"time"

	connect2 "bean/components/util/connect"
)

type (
	Membership struct {
		ID         string     `json:"id"`
		Version    string     `json:"version"`
		SpaceID    string     `json:"spaceId"`
		UserID     string     `json:"userId"`
		IsActive   bool       `json:"isActive"`
		CreatedAt  time.Time  `json:"createdAt"`
		UpdatedAt  time.Time  `json:"updatedAt"`
		LoggedInAt *time.Time `json:"loggedInAt"`
	}

	MembershipConnection struct {
		PageInfo MembershipInfo `json:"pageInfo"`
		Nodes    []Membership   `json:"nodes"`
	}

	MembershipEdge struct {
		Cursor string     `json:"cursor"`
		Node   Membership `json:"node"`
	}

	MembershipInfo struct {
		EndCursor   *string `json:"endCursor"`
		HasNextPage bool    `json:"hasNextPage"`
		StartCursor *string `json:"startCursor"`
	}
)

func (Membership) TableName() string {
	return "space_memberships"
}

func MembershipNodeCursor(node Membership) string {
	after := connect2.Cursor{
		Entity:   "Membership",
		Property: "logged_in_at",
		Value:    node.LoggedInAt.String(),
	}

	return after.Encode()
}
