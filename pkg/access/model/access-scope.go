package model

import (
	"fmt"
	"io"
	"strconv"
)

type AccessScope string

const (
	AccessScopeAnonymous     AccessScope = "Anonymous"
	AccessScopeAuthenticated AccessScope = "Authenticated"
)

var AllAccessScope = []AccessScope{
	AccessScopeAnonymous,
	AccessScopeAuthenticated,
}

func (e AccessScope) IsValid() bool {
	switch e {
	case AccessScopeAnonymous, AccessScopeAuthenticated:
		return true
	}
	return false
}

func (e AccessScope) String() string {
	return string(e)
}

func (e *AccessScope) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = AccessScope(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid AccessScope", str)
	}
	return nil
}

func (e AccessScope) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
