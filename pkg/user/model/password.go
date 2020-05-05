package model

import "time"

type UserPassword struct {
	ID          string    `json:"id"`
	Algorithm   string    `json:"algorithm"`
	HashedValue string    `json:"hashedValue"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	IsActive    bool      `json:"isActive"`
}
