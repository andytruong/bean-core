package model

import (
	"fmt"
	"io"
	"strconv"
)

type NamespaceKind string

const (
	NamespaceKindOrganisation NamespaceKind = "Organisation"
	NamespaceKindRole         NamespaceKind = "Role"
)

var AllNamespaceKind = []NamespaceKind{
	NamespaceKindOrganisation,
	NamespaceKindRole,
}

func (e NamespaceKind) IsValid() bool {
	switch e {
	case NamespaceKindOrganisation, NamespaceKindRole:
		return true
	}
	return false
}

func (e NamespaceKind) String() string {
	return string(e)
}

func (e *NamespaceKind) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = NamespaceKind(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid NamespaceKind", str)
	}
	return nil
}

func (e NamespaceKind) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
