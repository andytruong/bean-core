package model

import (
	"fmt"
	"io"
	"strconv"
)

type MailerAccountStatus string

const (
	MailerAccountStatusInactiveUnverified MailerAccountStatus = "INACTIVE_UNVERIFIED"
	MailerAccountStatusInactiveVerified   MailerAccountStatus = "INACTIVE_VERIFIED"
	MailerAccountStatusActiveUnverified   MailerAccountStatus = "ACTIVE_UNVERIFIED"
	MailerAccountStatusActiveVerified     MailerAccountStatus = "ACTIVE_VERIFIED"
)

var AllMailerAccountStatus = []MailerAccountStatus{
	MailerAccountStatusInactiveUnverified,
	MailerAccountStatusInactiveVerified,
	MailerAccountStatusActiveUnverified,
	MailerAccountStatusActiveVerified,
}

func (e MailerAccountStatus) IsValid() bool {
	switch e {
	case MailerAccountStatusInactiveUnverified, MailerAccountStatusInactiveVerified, MailerAccountStatusActiveUnverified, MailerAccountStatusActiveVerified:
		return true
	}
	return false
}

func (e MailerAccountStatus) String() string {
	return string(e)
}

func (e *MailerAccountStatus) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = MailerAccountStatus(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid MailerAccountStatus", str)
	}
	return nil
}

func (e MailerAccountStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
