package dto

import (
	"fmt"
	"io"
	"strconv"
)

type MailerTemplateEventKey string

const (
	MailerTemplateEventKeyCreate MailerTemplateEventKey = "CREATE"
	MailerTemplateEventKeyUpdate MailerTemplateEventKey = "UPDATE"
	MailerTemplateEventKeyDelete MailerTemplateEventKey = "DELETE"
)

var AllMailerTemplateEventKey = []MailerTemplateEventKey{
	MailerTemplateEventKeyCreate,
	MailerTemplateEventKeyUpdate,
	MailerTemplateEventKeyDelete,
}

func (e MailerTemplateEventKey) IsValid() bool {
	switch e {
	case MailerTemplateEventKeyCreate, MailerTemplateEventKeyUpdate, MailerTemplateEventKeyDelete:
		return true
	}
	return false
}

func (e MailerTemplateEventKey) String() string {
	return string(e)
}

func (e *MailerTemplateEventKey) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = MailerTemplateEventKey(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid MailerTemplateEventKey", str)
	}
	return nil
}

func (e MailerTemplateEventKey) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
