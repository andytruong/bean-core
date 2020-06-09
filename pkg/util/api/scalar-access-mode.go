package api

import (
	"fmt"
	"io"
)

// 3 length string with each has value range from 0 to 7.
type AccessMode string

func (this *AccessMode) UnmarshalGQL(v interface{}) error {
	if input, ok := v.(string); !ok {
		return fmt.Errorf("access-mode must be strings")
	} else {
		*this = AccessMode(input)
	}

	return nil
}

func (this AccessMode) MarshalGQL(w io.Writer) {
	fmt.Fprintf(w, `"%s"`, this)
}
