package handler

import (
	"context"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"

	"bean/pkg/access/model"
	"bean/pkg/access/model/dto"
	"bean/pkg/util"
)

type SessionCreateHandler struct {
	ID *util.Identifier
}

func (this SessionCreateHandler) SessionCreate(ctx context.Context, tx *gorm.DB, input *dto.LoginInput) (*dto.LoginOutcome, error) {
	if true {
		if id, err := this.ID.ULID(); nil != err {
			return nil, err
		} else {
			// find user object
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
				CreatedAt:   time.Time{},
				UpdatedAt:   time.Time{},
				ExpiredAt:   time.Time{},
			}

			fmt.Println("session", session)
		}
	}

	panic("wip")
}
