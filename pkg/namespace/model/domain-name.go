package model

import "time"

type DomainName struct {
	ID          string    `json:"id"`
	NamespaceId string    `json:"namespaceId"`
	Value       string    `json:"value"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	IsPrimary   bool      `json:"isPrimary"`
	IsVerified  bool      `json:"isVerified"`
	IsActive    bool      `json:"isActive"`
}

func (this DomainName) TableName() string {
	return "namespace_domains"
}

type DomainNames struct {
	Primary   *DomainName   `json:"primary"`
	Secondary []*DomainName `json:"secondary"`
}
