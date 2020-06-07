package model

import (
	"time"

	"bean/pkg/util/connect"
)

type (
	Membership struct {
		ID          string     `json:"id"`
		Version     string     `json:"version"`
		NamespaceID string     `json:"namespaceId"`
		UserID      string     `json:"userId"`
		IsActive    bool       `json:"isActive"`
		CreatedAt   time.Time  `json:"createdAt"`
		UpdatedAt   time.Time  `json:"updatedAt"`
		LoggedInAt  *time.Time `json:"loggedInAt"`
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

func MembershipNodeCursor(node Membership) string {
	after := connect.Cursor{
		Entity:   "Membership",
		Property: "logged_in_at",
		Value:    node.LoggedInAt.String(),
	}

	return after.Encode()
}
