package handler

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"

	"bean/pkg/access/model"
	"bean/pkg/access/model/dto"
	"bean/pkg/user"
	"bean/pkg/util"
)

type SessionCreateHandler struct {
	ID         *util.Identifier
	UserModule *user.UserModule
}

func (this SessionCreateHandler) SessionCreate(ctx context.Context, tx *gorm.DB, input *dto.LoginInput) (*dto.LoginOutcome, error) {
	if true {
		if id, err := this.ID.ULID(); nil != err {
			return nil, err
		} else {
			// find user object
			user, err := this.UserModule.User(ctx, input.Username)
			if nil != err {
				return nil, errors.New("user not found")
			}

			fmt.Println("USER: ", user)

			// compare the input password
			// find membership
			// create session
			// input.NamespaceID
			// input.Username
			// input.HashedPassword

			session := model.Session{
				ID:          id,
				HashedToken: "",
				Scopes:      nil,
				Context:     nil,
				IsActive:    false,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
				ExpiredAt:   time.Now(),
			}

			fmt.Println("session", session)
		}
	}

	panic("wip")
}
