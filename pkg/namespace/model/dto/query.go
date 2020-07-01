package dto

import (
	"bean/pkg/util/api/scalar"
)

type NamespaceFilters struct {
	ID     *string     `json:"id"`
	Domain *scalar.Uri `json:"domain"`
}
