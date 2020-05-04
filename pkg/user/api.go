package user

import (
	"context"

	"bean/pkg/user/dto"
)

type (
	UserMutationResolver struct {
	}
)

func (this *UserMutationResolver) UserCreate(ctx context.Context, input *dto.UserCreateInput) (*dto.UserCreateOutcome, error) {
	// validate email address
	// validate avatar URI
	// how to ID like youtube?
	//  - https://blog.codinghorror.com/url-shortening-hashes-in-practice/
	//  - http://www.fileformat.info/tool/hash.htm
	//  - https://hashids.org/
	// create base record
	// create emails
	// save name object
	// outcome

	panic("not implemented")
}
