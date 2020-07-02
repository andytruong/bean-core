package scalar

import (
	"time"
)

func NilBool(in bool) *bool {
	return &in
}

func NilString(in string) *string {
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
