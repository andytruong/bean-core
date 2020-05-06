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
	ass.NotNil(t, container.modules.User())
}
