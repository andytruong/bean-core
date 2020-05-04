package model

type UserName struct {
	FirstName     *string `json:"firstName"`
	LastName      *string `json:"lastName"`
	PreferredName *string `json:"preferredName"`
}
