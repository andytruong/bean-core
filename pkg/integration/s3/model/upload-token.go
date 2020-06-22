package model

import "time"

type UploadToken struct {
	ID          string    `json:"id"`
	NamespaceId string    `json:"namespaceId"`
	UserId      string    `json:"userId"`
	FilePath    string    `json:"filePath"`
	CreatedAt   time.Time `json:"createdAt"`
}
