package module

import (
	"github.com/99designs/gqlgen/codegen"
	"gorm.io/gorm"
)

type Bundle interface {
	Name() string
	Migrate(tx *gorm.DB, driver string) error
	Dependencies() []Bundle
	GraphqlResolver() map[string]interface{}
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
