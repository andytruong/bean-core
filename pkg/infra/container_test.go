package infra

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func NewMockContainer() *Container {
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
	can := NewMockContainer()
	idr := can.identifier
	ass.NotNil(t, idr)

	sv, err := can.bundles.User()
	ass.NoError(err)
	ass.NotNil(sv)
	ass.Equal("128h0m0s", can.Config.Bundles.Access.SessionTimeout.String())
	ass.Equal(100, can.Config.Bundles.Space.Manager.MaxNumberOfManager)
	ass.Equal("01EBWB516AP6BQD7", can.Config.Bundles.Integration.S3.Key)
}
