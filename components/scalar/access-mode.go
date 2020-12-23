package scalar

import (
	"fmt"
	"io"
	"strconv"

	"golang.org/x/sync/errgroup"
)

// 3 length string with each has value range from 0 to 7.
// Read: 4  Write: 2    Delete: 1
type AccessMode string

const (
	AccessModeSystemDisabled  AccessMode = "000"
	AccessModePrivateReadonly AccessMode = "400"
	AccessModePrivate         AccessMode = "700"
	AccessModePublicRead      AccessMode = "444"
)

func (mode AccessMode) Owner() string {
	return string(string(mode)[0])
}

func (mode AccessMode) Group() string {
	return string(string(mode)[1])
}

func (mode AccessMode) Public() string {
	return string(string(mode)[2])
}

func (mode AccessMode) CanRead(isOwner, isMember bool) bool {
	return mode.can(isOwner, isMember, 4)
}

func (mode AccessMode) CanWrite(isOwner, isMember bool) bool {
	return mode.can(isOwner, isMember, 2)
}

func (mode AccessMode) CanDelete(isOwner, isMember bool) bool {
	return mode.can(isOwner, isMember, 1)
}

func (mode AccessMode) can(isOwner, isMember bool, check int) bool {
	// Check public access
	if mode.check(mode.Public(), check) {
		return true
	}

	if isMember && mode.check(mode.Group(), check) {
		return true
	}

	if isOwner && mode.check(mode.Owner(), check) {
		return true
	}

	return false
}

func (mode AccessMode) check(scope string, check int) bool {
	actual, err := strconv.Atoi(scope)
	if nil != err {
		return false
	}

	if check < 4 {
		actual %= 4
	}

	if check < 2 {
		actual %= 2
	}

	if actual-check >= 0 {
		return true
	}

	return false
}

func (mode *AccessMode) UnmarshalGQL(v interface{}) error {
	if in, ok := v.(string); !ok {
		return fmt.Errorf("access-mode must be strings")
	} else if len(in) != 3 {
		return fmt.Errorf("access-mode must be strings with length of 3")
	} else {
		*mode = AccessMode(in)

		g := errgroup.Group{}
		g.Go(func() error { return mode.isValidMode(mode.Owner()) })
		g.Go(func() error { return mode.isValidMode(mode.Group()) })
		g.Go(func() error { return mode.isValidMode(mode.Public()) })

		return g.Wait()
	}
}

func (mode *AccessMode) isValidMode(char string) error {
	val, err := strconv.Atoi(char)

	if nil != err {
		return err
	}

	if val < 0 || val > 7 {
		return fmt.Errorf("invalid access-mode: " + char)
	}

	return nil
}

func (mode AccessMode) MarshalGQL(w io.Writer) {
	fmt.Fprintf(w, `"%s"`, mode)
}
