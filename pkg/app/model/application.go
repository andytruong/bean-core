package model

import (
	"time"
)

type Application struct {
	ID        string     `json:"id"`
	Version   string     `json:"version"`
	IsActive  bool       `json:"isActive"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt"`
}

func (Application) TableName() string {
	return "applications"
}
