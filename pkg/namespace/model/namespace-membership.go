package model

import "time"

type Membership struct {
	ID          string     `json:"id"`
	NamespaceID string     `json:"namespaceId"`
	UserID      string     `json:"userId"`
	IsActive    bool       `json:"isActive"`
	CreatedAt   *time.Time `json:"createdAt"`
	UpdatedAt   *time.Time `json:"updatedAt"`
}
