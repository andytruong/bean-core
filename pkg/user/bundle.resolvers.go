package user

import (
	"context"

	"gorm.io/gorm"

	model2 "bean/pkg/space/model"
	"bean/pkg/user/model"
	"bean/pkg/user/model/dto"
	"bean/pkg/util/connect"
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
			"Load": func(ctx context.Context, id string) (*model.User, error) {
				return bundle.Service.Load(bundle.db.WithContext(ctx), id)
			},
		},
		"UserMutation": map[string]interface{}{
			"Create": func(ctx context.Context, in *dto.UserCreateInput) (*dto.UserMutationOutcome, error) {
				var err error
				var out *dto.UserMutationOutcome

				err = connect.Transaction(ctx, bundle.db, func(tx *gorm.DB) error {
					out, err = bundle.Service.Create(tx, in)

					return err
				})

				return out, err
			},
			"Update": func(ctx context.Context, input dto.UserUpdateInput) (*dto.UserMutationOutcome, error) {
				var err error
				var out *dto.UserMutationOutcome

				err = connect.Transaction(ctx, bundle.db, func(tx *gorm.DB) error {
					out, err = bundle.Service.Update(tx, input)

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
