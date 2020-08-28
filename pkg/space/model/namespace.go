package model

import (
	"time"

	"bean/pkg/util/api"
)

type Space struct {
	ID        string       `json:"id"`
	ParentID  *string      `json:"parentId"`
	Version   string       `json:"version"`
	Kind      SpaceKind    `json:"kind"`
	Title     string       `json:"title"`
	IsActive  bool         `json:"isActive"`
	CreatedAt time.Time    `json:"createdAt"`
	UpdatedAt time.Time    `json:"updatedAt"`
	Language  api.Language `json:"language"`
}

func (this Space) TableName() string {
	return "spaces"
}