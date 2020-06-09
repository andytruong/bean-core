package api

import (
	"fmt"
	"io"
)

// 3 length string with each has value range from 0 to 7.
// Read: 4  Write: 2    Delete: 1
type AccessMode string

const (
	AccessModeSystemDisabled  AccessMode = "000"
	AccessModePrivateReadonly AccessMode = "400"
	AccessModePrivate         AccessMode = "600"
	AccessModePublicRead      AccessMode = "444"
	AccessModeInternalRead    AccessMode = "400"
)

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
