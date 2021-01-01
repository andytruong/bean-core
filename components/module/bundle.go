package module

import (
	"context"

	"github.com/99designs/gqlgen/codegen"
)

type Bundle interface {
	Name() string
	Migrate(ctx context.Context, driver string) error
	Dependencies() []Bundle
	GraphqlResolver() map[string]interface{}
	// TODO: Scopes
}

type GraphqlResolver interface {
	Aware(o *codegen.Object, f *codegen.Field) bool
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
