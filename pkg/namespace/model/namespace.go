package model

import (
	"time"

	"bean/pkg/util/api"
)

type Namespace struct {
	ID        string       `json:"id"`
	Version   string       `json:"version"`
	Title     *string      `json:"title"`
	IsActive  bool         `json:"isActive"`
	CreatedAt time.Time    `json:"createdAt"`
	UpdatedAt time.Time    `json:"updatedAt"`
	Language  api.Language `json:"language"`
}
