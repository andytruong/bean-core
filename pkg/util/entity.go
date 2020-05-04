package util

import (
	"fmt"
	"io"
	"strconv"
)

type Error struct {
	Code   *ErrorCode `json:"code"`
	Fields []string   `json:"fields"`
}

type ErrorCode string

const (
	ErrorCodeInut         ErrorCode = "Inut"
	ErrorCodeConfig       ErrorCode = "Config"
	ErrorCodeRuntime      ErrorCode = "Runtime"
	ErrorCodeDbTimeout    ErrorCode = "DB_Timeout"
	ErrorCodeDbConstraint ErrorCode = "DB_Constraint"
)

var AllErrorCode = []ErrorCode{
	ErrorCodeInut,
	ErrorCodeConfig,
	ErrorCodeRuntime,
	ErrorCodeDbTimeout,
	ErrorCodeDbConstraint,
}

func (e ErrorCode) IsValid() bool {
	switch e {
	case ErrorCodeInut, ErrorCodeConfig, ErrorCodeRuntime, ErrorCodeDbTimeout, ErrorCodeDbConstraint:
		return true
	}
	return false
}

func (e ErrorCode) String() string {
	return string(e)
}

func (e *ErrorCode) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ErrorCode(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ErrorCode", str)
	}
	return nil
}

func (e ErrorCode) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
