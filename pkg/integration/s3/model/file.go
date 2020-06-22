package model

import "time"

type File struct {
	ID            string    `json:"id"`
	Version       string    `json:"version"`
	ApplicationId string    `json:"applicationId"`
	Path          string    `json:"path"`
	Size          float32   `json:"size"`
	IsActive      bool      `json:"isActive"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}
