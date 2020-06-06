package dto

type MembershipsFilter struct {
	UserID    string                      `json:"userId"`
	Namespace *MembershipsFilterNamespace `json:"namespace"`
}

type MembershipsFilterNamespace struct {
	Title      string `json:"title"`
	DomainName string `json:"domainName"`
}
