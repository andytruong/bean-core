package model

import "time"

type ManagerRelationship struct {
	ID              string    `json:"id"`
	Version         string    `json:"version"`
	UserMemberId    string    `json:"userMemberId"`
	ManagerMemberId string    `json:"managerMemberId"`
	IsActive        bool      `json:"isActive"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

// TODO: Review this
func (ManagerRelationship) TableName() string {
	return "TableManagerEdge"
}
