package model

import "time"

type UploadToken struct {
	ID        string    `json:"id"`
	SpaceId   string    `json:"spaceId"`
	UserId    string    `json:"userId"`
	FilePath  string    `json:"filePath"`
	CreatedAt time.Time `json:"createdAt"`
}
