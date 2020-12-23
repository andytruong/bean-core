package scalar

import (
	"fmt"
	"io"
	"net/mail"
	"strings"
)

type EmailAddress string

func (address EmailAddress) LowerCaseValue() EmailAddress {
	stringRaw := string(address)
	stringLower := strings.ToLower(stringRaw)

	return EmailAddress(stringLower)
}

func (address *EmailAddress) UnmarshalGQL(v interface{}) error {
	if in, ok := v.(string); !ok {
		return fmt.Errorf("email-address must be strings")
	} else {
		_, err := mail.ParseAddress(in)
		if nil != err {
			return err
		}

		*address = EmailAddress(in)
	}

	return nil
}

func (address EmailAddress) MarshalGQL(w io.Writer) {
	fmt.Fprintf(w, `"%s"`, address)
}
