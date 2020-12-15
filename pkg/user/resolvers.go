package user

import (
	"context"
	
	"gorm.io/gorm"
	
	"bean/pkg/user/model"
	"bean/pkg/user/model/dto"
	"bean/pkg/util/connect"
)

func newResolvers(bundle *UserBundle) map[string]interface{} {
	return map[string]interface{}{
		"Query": map[string]interface{}{
			"User": func(ctx context.Context, id string) (*model.User, error) {
				db := bundle.db.WithContext(ctx)
				
				return bundle.UserService.Load(db, id)
			},
		},
		"Mutation": map[string]interface{}{
			"UserCreate": func(ctx context.Context, in *dto.UserCreateInput) (*dto.UserMutationOutcome, error) {
				var err error
				var out *dto.UserMutationOutcome
				
				err = connect.Transaction(ctx, bundle.db, func(tx *gorm.DB) error {
					out, err = bundle.UserService.Create(tx, in)
					
					return err
				})
				
				return out, err
			},
			"UserUpdate": func(ctx context.Context, input dto.UserUpdateInput) (*dto.UserMutationOutcome, error) {
				var err error
				var out *dto.UserMutationOutcome
				
				err = connect.Transaction(ctx, bundle.db, func(tx *gorm.DB) error {
					out, err = bundle.UserService.Update(tx, input)
					
					return err
				})
				
				return out, err
			},
		},
		"User": map[string]interface{}{
			"Name": func(ctx context.Context, user *model.User) (*model.UserName, error) {
				return bundle.NameService.load(bundle.db.WithContext(ctx), user.ID)
			},
			"Verified": func(ctx context.Context, obj *model.UserEmail) (bool, error) {
				return obj.IsVerified, nil
			},
			"Emails": func(ctx context.Context, obj *model.User) (*model.UserEmails, error) {
				return bundle.EmailService.List(ctx, obj)
			},
		},
	}
}
