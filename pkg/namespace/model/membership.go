package model

import "time"

type Membership struct {
	ID          string     `json:"id"`
	Version     string     `json:"version"`
	NamespaceID string     `json:"namespaceId"`
	UserID      string     `json:"userId"`
	IsActive    bool       `json:"isActive"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	LoggedInAt  *time.Time `json:"loggedInAt"`
}

type MembershipConnection struct {
	Edges    []*MembershipEdge `json:"edges"`
	Nodes    []*MembershipEdge `json:"nodes"`
	PageInfo *MembershipInfo   `json:"pageInfo"`
}

type MembershipEdge struct {
	Node *Membership `json:"node"`
}

type MembershipInfo struct {
	EndCursor       *string `json:"endCursor"`
	HasNextPage     bool    `json:"hasNextPage"`
	HasPreviousPage bool    `json:"hasPreviousPage"`
	StartCursor     *string `json:"startCursor"`
}
