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
	if in, ok := v.(string); !ok {
		return fmt.Errorf("email-address must be strings")
	} else {
		_, err := mail.ParseAddress(in)
		if nil != err {
			return err
		}

		*this = EmailAddress(in)
	}

	return nil
}

func (this EmailAddress) MarshalGQL(w io.Writer) {
	fmt.Fprintf(w, `"%s"`, this)
}

func (this *Uri) UnmarshalGQL(v interface{}) error {
	if in, ok := v.(string); !ok {
		return fmt.Errorf("URI must be strings")
	} else {
		uri, err := url.ParseRequestURI(in)
		if nil != err {
			return fmt.Errorf("invalid URI for request")
		}

		if "" == uri.RequestURI() {
			return fmt.Errorf("missing request URI")
		}

		*this = Uri(in)
	}

	return nil
}

func (this Uri) MarshalGQL(w io.Writer) {
	fmt.Fprintf(w, `"%s"`, this)
}

func NilBool(in bool) *bool {
	return &in
}

func NilString(in string) *string {
	return &in
}

func NilUri(in Uri) *Uri {
	return &in
}

func NilTime(time time.Time) *time.Time {
	return &time
}

func NotNilString(in *string, defaultValue string) string {
	if nil != in {
		return *in
	}

	return defaultValue
}

func NotNilBool(in *bool, defaultValue bool) bool {
	if nil != in {
		return *in
	}

	return defaultValue
}
