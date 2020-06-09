package model

import (
	"time"

	"bean/pkg/util/api"
)

type ConfigBucket struct {
	Id        string         `json:"id"`
	Version   string         `json:"version"`
	Title     string         `json:"title"`
	Slug      string         `json:"machineName"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	UserId    string         `json:"userId"`
	HostId    string         `json:"hostId"`
	Access    api.AccessMode `json:"access"`
}
