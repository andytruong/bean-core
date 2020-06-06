package api

import (
	"fmt"
	"io"
	"strconv"
)

type Sorts string

const (
	SortsCreatedAtAsc  Sorts = "CreatedAtAsc"
	SortsCreatedAtDesc Sorts = "CreatedAtDesc"
	SortsUpdatedAtAsc  Sorts = "UpdatedAtAsc"
	SortsUpdatedAtDesc Sorts = "UpdatedAtDesc"
)

var AllSorts = []Sorts{
	SortsCreatedAtAsc,
	SortsCreatedAtDesc,
	SortsUpdatedAtAsc,
	SortsUpdatedAtDesc,
}

func (e Sorts) IsValid() bool {
	switch e {
	case SortsCreatedAtAsc, SortsCreatedAtDesc, SortsUpdatedAtAsc, SortsUpdatedAtDesc:
		return true
	}
	return false
}

func (e Sorts) String() string {
	return string(e)
}

func (e *Sorts) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Sorts(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Sorts", str)
	}
	return nil
}

func (e Sorts) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
