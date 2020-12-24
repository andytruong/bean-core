package user

import (
	"context"

	"gorm.io/gorm"

	connect2 "bean/components/util/connect"
	model2 "bean/pkg/space/model"
	"bean/pkg/user/model"
	"bean/pkg/user/model/dto"
)

func newResolvers(bundle *UserBundle) map[string]interface{} {
	return map[string]interface{}{
		"Query": map[string]interface{}{
			"UserQuery": func(ctx context.Context) (*dto.UserQuery, error) {
				return &dto.UserQuery{}, nil
			},
			"Membership": func(ctx context.Context, id string, version *string) (*model2.Membership, error) {
				panic("TODO")
			},
		},
		"Mutation": map[string]interface{}{
			"UserMutation": func(ctx context.Context) (*dto.UserMutation, error) {
				return &dto.UserMutation{}, nil
			},
		},
		"UserQuery": map[string]interface{}{
			"Load": func(ctx context.Context, _ *dto.UserQuery, id string) (*model.User, error) {
				return bundle.Service.Load(ctx, id)
			},
		},
		"UserMutation": map[string]interface{}{
			"Create": func(ctx context.Context, _ *dto.UserMutation, in *dto.UserCreateInput) (*dto.UserMutationOutcome, error) {
				var err error
				var out *dto.UserMutationOutcome

				err = connect2.Transaction(ctx, bundle.db, func(tx *gorm.DB) error {
					out, err = bundle.Service.Create(connect2.DBToContext(ctx, tx), in)

					return err
				})

				return out, err
			},
			"Update": func(ctx context.Context, _ *dto.UserMutation, input dto.UserUpdateInput) (*dto.UserMutationOutcome, error) {
				var err error
				var out *dto.UserMutationOutcome

				err = connect2.Transaction(ctx, bundle.db, func(tx *gorm.DB) error {
					out, err = bundle.Service.Update(connect2.DBToContext(ctx, tx), input)

					return err
				})

				return out, err
			},
		},
		"User": map[string]interface{}{
			"Name": func(ctx context.Context, user *model.User) (*model.UserName, error) {
				return bundle.nameService.load(bundle.db.WithContext(ctx), user.ID)
			},
			"Verified": func(ctx context.Context, obj *model.UserEmail) (bool, error) {
				return obj.IsVerified, nil
			},
			"Emails": func(ctx context.Context, obj *model.User) (*model.UserEmails, error) {
				return bundle.emailService.List(ctx, obj)
			},
		},
	}
}
