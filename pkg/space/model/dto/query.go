package dto

import (
	"bean/components/scalar"
)

type SpaceFilters struct {
	ID     *string     `json:"id"`
	Domain *scalar.Uri `json:"domain"`
}
