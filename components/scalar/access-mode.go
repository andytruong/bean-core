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

func (this AccessMode) Owner() string {
	return string(string(this)[0])
}

func (this AccessMode) Group() string {
	return string(string(this)[1])
}

func (this AccessMode) Public() string {
	return string(string(this)[2])
}

func (this AccessMode) CanRead(isOwner, isMember bool) bool {
	return this.can(isOwner, isMember, 4)
}

func (this AccessMode) CanWrite(isOwner, isMember bool) bool {
	return this.can(isOwner, isMember, 2)
}

func (this AccessMode) CanDelete(isOwner, isMember bool) bool {
	return this.can(isOwner, isMember, 1)
}

func (this AccessMode) can(isOwner, isMember bool, check int) bool {
	// Check public access
	if this.check(this.Public(), check) {
		return true
	}

	if isMember && this.check(this.Group(), check) {
		return true
	}

	if isOwner && this.check(this.Owner(), check) {
		return true
	}

	return false
}

func (this AccessMode) check(scope string, check int) bool {
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

func (this *AccessMode) UnmarshalGQL(v interface{}) error {
	if in, ok := v.(string); !ok {
		return fmt.Errorf("access-mode must be strings")
	} else if len(in) != 3 {
		return fmt.Errorf("access-mode must be strings with length of 3")
	} else {
		*this = AccessMode(in)

		g := errgroup.Group{}
		g.Go(func() error { return this.isValidMode(this.Owner()) })
		g.Go(func() error { return this.isValidMode(this.Group()) })
		g.Go(func() error { return this.isValidMode(this.Public()) })

		return g.Wait()
	}
}

func (this *AccessMode) isValidMode(char string) error {
	val, err := strconv.Atoi(char)

	if nil != err {
		return err
	}

	if val < 0 || val > 7 {
		return fmt.Errorf("invalid access-mode: " + char)
	}

	return nil
}

func (this AccessMode) MarshalGQL(w io.Writer) {
	fmt.Fprintf(w, `"%s"`, this)
}
