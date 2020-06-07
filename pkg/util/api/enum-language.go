package api

import (
	"fmt"
	"io"
	"strconv"
)

type Language string

const (
	// See more at https://en.wikipedia.org/wiki/List_of_ISO_639-1_codes
	LanguageAustralia     Language = "AU"
	LanguageUnitedStates  Language = "US"
	LanguageUnitedKingdom Language = "UK"
	LanguageVietnam       Language = "VN"
)

var AllLanguage = []Language{
	LanguageAustralia,
	LanguageUnitedStates,
	LanguageUnitedKingdom,
	LanguageVietnam,
}

func (this Language) IsValid() bool {
	switch this {
	case LanguageAustralia, LanguageUnitedStates, LanguageUnitedKingdom, LanguageVietnam:
		return true
	}
	return false
}

func (this Language) String() string {
	return string(this)
}

func (this Language) Nil() *Language {
	return &this
}

func (this *Language) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*this = Language(str)
	if !this.IsValid() {
		return fmt.Errorf("%s is not a valid Language", str)
	}
	return nil
}

func (this Language) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(this.String()))
}
