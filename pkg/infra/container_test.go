package infra

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func NewMockCan() *Container {
	_ = os.Setenv("ENV", "test")
	_ = os.Setenv("DB_MASTER_URL", ":memory:")
	_ = os.Setenv("DB_SLAVE_URL", ":memory:")
	ctn, err := NewContainer("../../config.yaml")

	if nil != err {
		panic(err)
	}

	return ctn
}

func Test(t *testing.T) {
	t.Parallel()

	ass := assert.New(t)
	can := NewMockCan()
	id := can.Identifier()
	ass.NotNil(t, id)

	sv, err := can.bundles.User()
	ass.NoError(err)
	ass.NotNil(sv)
	ass.Equal("128h0m0s", can.Bundles.Access.SessionTimeout.String())
	ass.Equal(100, can.Bundles.Space.Manager.MaxNumberOfManager)
	ass.Equal("01EBWB516AP6BQD7", can.Bundles.Integration.S3.Key)
}