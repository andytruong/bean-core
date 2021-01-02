package model

import (
	"time"
)

type Application struct {
	ID        string     `json:"id"`
	Version   string     `json:"version"`
	SpaceId   string     `json:"spaceId"`
	IsActive  bool       `json:"isActive"`
	Title     *string    `json:"title"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt"`
}

func (Application) TableName() string {
	return "applications"
}
