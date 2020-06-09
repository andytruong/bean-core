package model

import "time"

type ConfigVariable struct {
	Id        string    `json:"id"`
	Version   string    `json:"version"`
	BucketId  string    `json:"bucketId"`
	Name      string    `json:"name"`
	Value     []byte    `json:"value"`
	CreateAt  time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
