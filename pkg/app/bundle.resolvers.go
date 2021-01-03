package app

import (
	"context"

	"bean/components/util"
	"bean/pkg/app/model"
	"bean/pkg/app/model/dto"
)

func (bundle *Bundle) newResolvers() map[string]interface{} {
	return map[string]interface{}{
		"Query": map[string]interface{}{
			"ApplicationQuery": map[string]interface{}{
				"Load": func(ctx context.Context, id string, version *string) (*model.Application, error) {
					app, err := bundle.Service.Load(ctx, id)

					if nil != err {
						return nil, err
					}

					if nil != version {
						if *version != app.Version {
							return nil, util.ErrorVersionConflict
						}
					}

					return app, nil
				},
			},
		},
		"Mutation": map[string]interface{}{
			"ApplicationMutation": map[string]interface{}{
				"Create": func(ctx context.Context, in *dto.ApplicationCreateInput) (*dto.ApplicationOutcome, error) {
					out, err := bundle.Service.Create(ctx, in)
					if nil != err {
						return nil, err
					}

					return out, nil
				},
				"Update": func(ctx context.Context, in *dto.ApplicationUpdateInput) (*dto.ApplicationOutcome, error) {
					out, err := bundle.Service.Update(ctx, in)
					if nil != err {
						return nil, err
					}

					return out, nil
				},
			},
		},
	}
}
