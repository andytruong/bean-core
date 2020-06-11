package util

import (
	"fmt"
	"io"
	"net/mail"
	"net/url"
	"strings"
	"time"

	"bean/pkg/util/migrate"
)

type (
	UUID         string
	ULID         string
	EmailAddress string
	Uri          string
	FilePath     string

	// Useful for simple DB query.
	ValueWrapper struct {
		Value string
	}
)

func (this FilePath) String() string {
	out := string(this)
	if strings.HasPrefix(out, "/") {
		return out
	}

	return migrate.RootDirectory() + "/" + out
}

func (this EmailAddress) LowerCaseValue() EmailAddress {
	stringRaw := string(this)
	stringLower := strings.ToLower(stringRaw)

	return EmailAddress(stringLower)
}

func (this *EmailAddress) UnmarshalGQL(v interface{}) error {
	if input, ok := v.(string); !ok {
		return fmt.Errorf("email-address must be strings")
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
		return fmt.Errorf("URI must be strings")
	} else {
		uri, err := url.ParseRequestURI(input)
		if nil != err {
			return fmt.Errorf("invalid URI for request")
		}

		if "" == uri.RequestURI() {
			return fmt.Errorf("missing request URI")
		}

		*this = Uri(input)
	}

	return nil
}

func (this Uri) MarshalGQL(w io.Writer) {
	fmt.Fprintf(w, `"%s"`, this)
}

func NilBool(input bool) *bool {
	return &input
}

func NilString(input string) *string {
	return &input
}

func NilUri(input Uri) *Uri {
	return &input
}

func NilTime(time time.Time) *time.Time {
	return &time
}

func NotNilString(input *string, defaultValue string) string {
	if nil != input {
		return *input
	}

	return defaultValue
}

func NotNilBool(input *bool, defaultValue bool) bool {
	if nil != input {
		return *input
	}

	return defaultValue
}
