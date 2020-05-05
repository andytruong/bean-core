package util

import (
	"fmt"
	"io"
	"net/mail"
	"net/url"
)

type (
	UUID         string
	ULID         string
	EmailAddress string
	Uri          string

	// Useful for simple DB query.
	ValueWrapper struct {
		Value string
	}
)

func (this *EmailAddress) UnmarshalGQL(v interface{}) error {
	if input, ok := v.(string); !ok {
		return fmt.Errorf("points must be strings")
	} else {
		_, err := mail.ParseAddress(input)
		if nil != err {
			return err
		}

		*this = EmailAddress(input)
	}

	return nil
}

func (this EmailAddress) MarshalGQL(w io.Writer) {
	fmt.Fprintf(w, `"%s"`, this)
}

func (this *Uri) UnmarshalGQL(v interface{}) error {
	if input, ok := v.(string); !ok {
		return fmt.Errorf("points must be strings")
	} else {
		url, err := url.ParseRequestURI(input)
		if nil != err {
			return err
		}

		if "" == url.RequestURI() {
			return fmt.Errorf("missing request URI")
		}

		*this = Uri(input)
	}

	return nil
}

func (this Uri) MarshalGQL(w io.Writer) {
	fmt.Fprintf(w, `"%s"`, this)
}
