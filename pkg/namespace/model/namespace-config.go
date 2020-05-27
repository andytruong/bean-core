package model

import "time"

type NamespaceConfig struct {
	Id          string    `json:"id"`
	Version     string    `json:"version"`
	NamespaceId string    `json:"namespaceId"`
	Bucket      string    `json:"bucket"`
	Key         string    `json:"key"`
	Value       []byte    `json:"value"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
