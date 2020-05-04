package model

type Namespace struct {
	ID          string       `json:"id"`
	Title       *string      `json:"title"`
	DomainNames *DomainNames `json:"domainNames"`
}
