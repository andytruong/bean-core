package api

import (
	"fmt"
	"io"
	"strconv"
)

type Language string

const (
	// See more at https://en.wikipedia.org/wiki/List_of_ISO_639-1_codes
	LanguageAU      Language = "AU"
	LanguageUS      Language = "US"
	LanguageUK      Language = "UK"
	LanguageVN      Language = "VN"
	LanguageDefault Language = "US"
)

var AllLanguage = []Language{
	LanguageAU,
	LanguageUS,
	LanguageUK,
	LanguageVN,
}

func (lang Language) IsValid() bool {
	switch lang {
	case LanguageAU, LanguageUS, LanguageUK, LanguageVN:
		return true
	}
	return false
}

func (lang Language) String() string {
	return string(lang)
}

func (lang Language) Nil() *Language {
	return &lang
}

func (lang *Language) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*lang = Language(str)
	if !lang.IsValid() {
		return fmt.Errorf("%s is not a valid Language", str)
	}
	return nil
}

func (lang Language) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(lang.String()))
}
