package api

import (
	"fmt"
	"io"
	"net/mail"
	"strings"
)

type EmailAddress string

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
