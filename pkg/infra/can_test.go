package infra

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func NewMockCan() *Can {
	_ = os.Setenv("DB_MASTER_URL", ":memory:")
	_ = os.Setenv("DB_SLAVE_URL", ":memory:")
	ctn, err := NewCan("../../config.yaml")

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

	sv, err := can.beans.User()
	ass.NoError(err)
	ass.NotNil(sv)
	ass.Equal("128h0m0s", can.Beans.Access.SessionTimeout.String())
	ass.Equal(100, can.Beans.Namespace.Manager.MaxNumberOfManager)
	ass.Equal("01EBWB516AP6BQD7", can.Beans.Integration.S3.Key)
}
