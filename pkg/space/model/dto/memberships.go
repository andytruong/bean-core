package dto

type MembershipsFilter struct {
	UserID    string                  `json:"userId"`
	Space     *MembershipsFilterSpace `json:"space"`
	IsActive  bool                    `json:"isActive"`
	ManagerId *string                 `json:"managerId"`
}

type MembershipsFilterSpace struct {
	Title      *string `json:"title"`
	DomainName *string `json:"domainName"`
}
