package module

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Dispatcher(t *testing.T) {
	log := []string{}

	hook := NewHook()
	hook.Hook(
		func(ctx context.Context, event interface{}) error {
			log = append(log, "from listener 2: "+event.(string))

			return nil
		},
		2,
	)

	hook.Hook(
		func(ctx context.Context, event interface{}) error {
			log = append(log, "from listener 1: "+event.(string))

			return nil
		},
		1,
	)

	ass := assert.New(t)
	err := hook.Invoke(context.Background(), "hi")
	ass.NoError(err)
	ass.Len(log, 2)
	ass.Equal("from listener 1: hi", log[0])
	ass.Equal("from listener 2: hi", log[1])
}
