package module

import (
	"github.com/99designs/gqlgen/codegen"
	"gorm.io/gorm"
)

type Bundle interface {
	Migrate(tx *gorm.DB, driver string) error
	Dependencies() []Bundle
	GetGraphqlResolver() GraphqlResolver
}

type GraphqlResolver interface {
	Aware(o *codegen.Object, f *codegen.Field) bool
}

type AbstractBundle struct {
}

func (this AbstractBundle) GetGraphqlResolver() GraphqlResolver {
	return nil
}
