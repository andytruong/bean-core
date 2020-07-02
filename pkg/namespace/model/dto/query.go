package dto

import (
	"bean/components/scalar"
)

type NamespaceFilters struct {
	ID     *string     `json:"id"`
	Domain *scalar.Uri `json:"domain"`
}
