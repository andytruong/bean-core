package s3

import (
	"context"
	
	"bean/pkg/integration/s3/model"
)

func newGraphqlResolver() map[string]interface{} {
	return map[string]interface{}{
		"Query": map[string]interface{}{},
		"Application": map[string]interface{}{
			"Polices": func(ctx context.Context, obj *model.Application) ([]*model.Policy, error) {
				panic("wip")
			},
			
			"Credentials": func(ctx context.Context, obj *model.Application) (*model.Credentials, error) {
				panic("no implementation")
			},
		},
	}
}
