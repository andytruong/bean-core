package model

import (
	"time"

	"bean/pkg/util/api"
)

type Namespace struct {
	ID        string        `json:"id"`
	ParentID  *string       `json:"parentId"`
	Version   string        `json:"version"`
	Kind      NamespaceKind `json:"kind"`
	Title     string        `json:"title"`
	IsActive  bool          `json:"isActive"`
	CreatedAt time.Time     `json:"createdAt"`
	UpdatedAt time.Time     `json:"updatedAt"`
	Language  api.Language  `json:"language"`
}
