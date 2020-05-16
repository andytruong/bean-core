package model

import "time"

type DomainName struct {
	ID          string    `json:"id"`
	NamespaceId string    `json:"namespaceId"`
	Verified    bool      `json:"verified"`
	Value       string    `json:"value"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	IsActive    bool      `json:"isActive"`
}

type DomainNames struct {
	Primary   *DomainName   `json:"primary"`
	Secondary []*DomainName `json:"secondary"`
}
