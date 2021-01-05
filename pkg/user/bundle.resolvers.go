package user

import (
	"context"
	
	"gorm.io/gorm"
	
	"bean/components/connect"
	spaceModel "bean/pkg/space/model"
	"bean/pkg/user/model"
	"bean/pkg/user/model/dto"
)

func newResolvers(bundle *Bundle) map[string]interface{} {
	return map[string]interface{}{
		"Query": map[string]interface{}{
			"UserQuery": func(ctx context.Context) (*dto.UserQuery, error) {
				return &dto.UserQuery{}, nil
			},
			"Membership": func(ctx context.Context, id string, version *string) (*spaceModel.Membership, error) {
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
				return bundle.UserService.Load(ctx, id)
			},
		},
		"UserMutation": map[string]interface{}{
			"Create": func(ctx context.Context, _ *dto.UserMutation, in *dto.UserCreateInput) (
				*dto.UserMutationOutcome, error,
			) {
				var err error
				var out *dto.UserMutationOutcome
				
				err = connect.Transaction(
					ctx,
					func(tx *gorm.DB) error {
						out, err = bundle.UserService.Create(connect.DBToContext(ctx, tx), in)
						
						return err
					},
				)
				
				return out, err
			},
			"Update": func(ctx context.Context, _ *dto.UserMutation, input dto.UserUpdateInput) (
				*dto.UserMutationOutcome, error,
			) {
				var err error
				var out *dto.UserMutationOutcome
				
				err = connect.Transaction(
					ctx,
					func(tx *gorm.DB) error {
						out, err = bundle.UserService.Update(connect.DBToContext(ctx, tx), input)
						
						return err
					},
				)
				
				return out, err
			},
		},
		"User": map[string]interface{}{
			"Name": func(ctx context.Context, user *model.User) (*model.UserName, error) {
				return bundle.nameService.load(ctx, user.ID)
			},
			"Verified": func(ctx context.Context, obj *model.UserEmail) (bool, error) {
				return obj.IsVerified, nil
			},
			"Emails": func(ctx context.Context, obj *model.User) (*model.UserEmails, error) {
				return bundle.EmailService.List(ctx, obj)
			},
		},
		"UserEmail": map[string]interface{}{
			"Verified": func(ctx context.Context, obj model.UserEmail) (bool, error) {
				return obj.IsVerified, nil
			},
		},
	}
}
