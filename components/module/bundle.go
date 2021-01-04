package module

import (
	"context"
)

type Bundle interface {
	Name() string
	Migrate(ctx context.Context, driver string) error
	Dependencies() []Bundle
	GraphqlResolver() map[string]interface{}
}

type AbstractBundle struct {
}

func (AbstractBundle) Name() string {
	panic("not implemented")
}

func (AbstractBundle) Dependencies() []Bundle {
	return nil
}

func (AbstractBundle) GraphqlResolver() map[string]interface{} {
	return nil
}
