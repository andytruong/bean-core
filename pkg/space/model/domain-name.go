package model

import "time"

type DomainName struct {
	ID         string    `json:"id"`
	SpaceId    string    `json:"spaceId"`
	Value      string    `json:"value"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
	IsPrimary  bool      `json:"isPrimary"`
	IsVerified bool      `json:"isVerified"`
	IsActive   bool      `json:"isActive"`
}

func (DomainName) TableName() string {
	return "space_domains"
}

type DomainNames struct {
	Primary   *DomainName   `json:"primary"`
	Secondary []*DomainName `json:"secondary"`
}
