package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_AccessMode(t *testing.T) {
	ass := assert.New(t)

	t.Run("valid", func(t *testing.T) {
		ass.Equal("4", AccessModePublicRead.Owner())
		ass.Equal("4", AccessModePublicRead.Group())
		ass.Equal("4", AccessModePublicRead.Public())

		t.Run("unmarshal", func(t *testing.T) {
			val := AccessMode("000")
			err := val.UnmarshalGQL("444")
			ass.NoError(err)
			ass.Equal("444", string(val))
		})
	})

	t.Run("invalid", func(t *testing.T) {
		t.Run("max length is 3", func(t *testing.T) {
			val := AccessMode("000")
			err := val.UnmarshalGQL("1234")
			ass.Error(err)
			ass.Contains(err.Error(), "must be strings with length of 3")
		})

		t.Run("range 0 to 7", func(t *testing.T) {
			val := AccessMode("000")
			err := val.UnmarshalGQL("248")
			ass.Error(err)
			ass.Contains(err.Error(), "invalid access-mode: 8")
		})
	})

	t.Run("access checking", func(t *testing.T) {
		t.Run("AccessModeSystemDisabled: all disabled", func(t *testing.T) {
			ass.False(AccessModeSystemDisabled.CanRead(true, true))
			ass.False(AccessModeSystemDisabled.CanRead(false, false))
			ass.False(AccessModeSystemDisabled.CanRead(false, true))
			ass.False(AccessModeSystemDisabled.CanWrite(true, true))
			ass.False(AccessModeSystemDisabled.CanWrite(false, false))
			ass.False(AccessModeSystemDisabled.CanWrite(false, true))
			ass.False(AccessModeSystemDisabled.CanDelete(true, true))
			ass.False(AccessModeSystemDisabled.CanDelete(false, false))
			ass.False(AccessModeSystemDisabled.CanDelete(false, true))
		})

		t.Run("AccessModePrivateReadonly: only available for user", func(t *testing.T) {
			ass.True(AccessModePrivateReadonly.CanRead(true, true))

			// all other: false
			ass.False(AccessModePrivateReadonly.CanRead(false, false))
			ass.False(AccessModePrivateReadonly.CanRead(false, true))
			ass.False(AccessModePrivateReadonly.CanWrite(true, true))
			ass.False(AccessModePrivateReadonly.CanWrite(false, false))
			ass.False(AccessModePrivateReadonly.CanWrite(false, true))
			ass.False(AccessModePrivateReadonly.CanDelete(true, true))
			ass.False(AccessModePrivateReadonly.CanDelete(false, false))
			ass.False(AccessModePrivateReadonly.CanDelete(false, true))
		})

		t.Run("AccessModePrivate: only available for user", func(t *testing.T) {
			ass.True(AccessModePrivate.CanRead(true, true))
			ass.True(AccessModePrivate.CanWrite(true, true))
			ass.True(AccessModePrivate.CanDelete(true, true))

			// all other: false
			ass.False(AccessModePrivate.CanRead(false, false))
			ass.False(AccessModePrivate.CanRead(false, true))
			ass.False(AccessModePrivate.CanWrite(false, false))
			ass.False(AccessModePrivate.CanWrite(false, true))
			ass.False(AccessModePrivate.CanDelete(false, false))
			ass.False(AccessModePrivate.CanDelete(false, true))
		})

		t.Run("owner", func(t *testing.T) {
			ass.False(AccessModeSystemDisabled.CanWrite(true, true))
			ass.True(AccessModePrivateReadonly.CanRead(true, false))
		})
	})
}
