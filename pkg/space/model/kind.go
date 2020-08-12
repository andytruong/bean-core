package model

import (
	"fmt"
	"io"
	"strconv"
)

type SpaceKind string

const (
	SpaceKindOrganisation SpaceKind = "Organisation"
	SpaceKindRole         SpaceKind = "Role"
)

var AllSpaceKind = []SpaceKind{
	SpaceKindOrganisation,
	SpaceKindRole,
}

func (e SpaceKind) IsValid() bool {
	switch e {
	case SpaceKindOrganisation, SpaceKindRole:
		return true
	}
	return false
}

func (e SpaceKind) String() string {
	return string(e)
}

func (e *SpaceKind) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = SpaceKind(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid SpaceKind", str)
	}
	return nil
}

func (e SpaceKind) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
