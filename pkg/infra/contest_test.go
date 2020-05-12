package infra

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func NewMockContainer() *Container {
	_ = os.Setenv("DB_MASTER_URL", ":memory:")
	_ = os.Setenv("DB_SLAVE_URL", ":memory:")
	ctn, err := NewContainer("../../config.yaml")

	if nil != err {
		panic(err)
	}

	return ctn
}

func TestContainer(t *testing.T) {
	ass := assert.New(t)
	container := NewMockContainer()
	id := container.Identifier()

	ass.NotNil(t, id)

	sv, err := container.modules.User()
	ass.NoError(err)
	ass.NotNil(sv)
}
