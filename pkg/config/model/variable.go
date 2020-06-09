package model

import "time"

type ConfigVariable struct {
	Id          string    `json:"id"`
	Version     string    `json:"version"`
	BucketId    string    `json:"bucketId"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	Value       string    `json:"value"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
