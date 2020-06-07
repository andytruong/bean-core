package dto

type MembershipsFilter struct {
	UserID    string                      `json:"userId"`
	Namespace *MembershipsFilterNamespace `json:"namespace"`
	IsActive  bool                        `json:"isActive"`
	ManagerId *string                     `json:"managerId"`
}

type MembershipsFilterNamespace struct {
	Title      *string `json:"title"`
	DomainName *string `json:"domainName"`
}
