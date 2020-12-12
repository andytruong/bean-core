package user

import (
	"context"
	
	"gorm.io/gorm"
	
	"bean/pkg/user/model"
	"bean/pkg/user/model/dto"
	"bean/pkg/util/connect"
)

func newResolvers(bean *UserBundle) *Resolvers {
	return &Resolvers{
		Object:   &UserObjectResolver{bean: bean},
		Query:    &UserQueryResolver{bean: bean},
		Mutation: &UserMutationResolver{bean: bean},
	}
}

type (
	Resolvers struct {
		Object   *UserObjectResolver
		Query    *UserQueryResolver
		Mutation *UserMutationResolver
	}
	
	UserObjectResolver struct {
		bean *UserBundle
	}
	
	UserQueryResolver struct {
		bean *UserBundle
	}
	
	UserMutationResolver struct {
		bean *UserBundle
	}
)

func (this *UserQueryResolver) User(ctx context.Context, id string) (*model.User, error) {
	db := this.bean.db.WithContext(ctx)
	
	return this.bean.Service.Load(db, id)
}

func (this UserObjectResolver) Name(ctx context.Context, user *model.User) (*model.UserName, error) {
	return this.bean.NameService.load(this.bean.db.WithContext(ctx), user.ID)
}

func (this UserObjectResolver) Verified(ctx context.Context, obj *model.UserEmail) (bool, error) {
	return obj.IsVerified, nil
}

func (this UserObjectResolver) Emails(ctx context.Context, obj *model.User) (*model.UserEmails, error) {
	return this.bean.EmailService.List(ctx, obj)
}

func (this *UserMutationResolver) UserCreate(ctx context.Context, in *dto.UserCreateInput) (*dto.UserMutationOutcome, error) {
	var err error
	var out *dto.UserMutationOutcome
	
	err = connect.Transaction(ctx, this.bean.db, func(tx *gorm.DB) error {
		out, err = this.bean.Service.Create(tx, in)
		
		return err
	})
	
	return out, err
}

func (this *UserMutationResolver) UserUpdate(ctx context.Context, input dto.UserUpdateInput) (*dto.UserMutationOutcome, error) {
	var err error
	var out *dto.UserMutationOutcome
	
	err = connect.Transaction(ctx, this.bean.db, func(tx *gorm.DB) error {
		out, err = this.bean.Service.Update(tx, input)
		
		return err
	})
	
	return out, err
}
