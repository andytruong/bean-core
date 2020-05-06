package infra

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func NewMockContainer() *Container {
	container, err := NewContainer("../../config.yaml")

	if nil != err {
		panic(err)
	}

	return container
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
