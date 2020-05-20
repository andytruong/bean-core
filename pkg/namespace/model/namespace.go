package model

import "time"

type Namespace struct {
	ID        string    `json:"id"`
	Version   string    `json:"version"`
	Title     *string   `json:"title"`
	IsActive  bool      `json:"isActive"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
