package model

type UserName struct {
	ID            string  `json:"id"`
	UserId        string  `json:"userId"`
	FirstName     *string `json:"firstName"`
	LastName      *string `json:"lastName"`
	PreferredName *string `json:"preferredName"`
}
