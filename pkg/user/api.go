package user

import (
	"context"
	"fmt"

	"bean/pkg/user/dto"
	"bean/pkg/user/model"
)

type (
	UserMutationResolver struct {
	}
)

func (this *UserMutationResolver) UserCreate(ctx context.Context, input *dto.UserCreateInput) (*dto.UserCreateOutcome, error) {
	// validate email address
	// validate avatar URI
	// create base record
	user := model.User{
		// how to ID like youtube?
		//  - https://blog.codinghorror.com/url-shortening-hashes-in-practice/
		//  - http://www.fileformat.info/tool/hash.htm
		//  - https://hashids.org/
		ID:        "TODO",
		AvatarURI: input.AvatarURI,
		IsActive:  input.IsActive,
	}

	// db := ctx.DB()
	// db.Create(user)
	fmt.Println("WIP", user)

	// create emails
	// save name object
	// outcome

	panic("not implemented")
}
