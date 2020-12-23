package model

import "time"

type SpaceConfig struct {
	Id        string    `json:"id"`
	Version   string    `json:"version"`
	SpaceId   string    `json:"spaceId"`
	Bucket    string    `json:"bucket"`
	Key       string    `json:"key"`
	Value     []byte    `json:"value"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (SpaceConfig) TableName() string {
	return "space_config"
}
