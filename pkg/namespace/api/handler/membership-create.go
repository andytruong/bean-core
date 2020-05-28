package handler

import (
	"context"
	"time"

	"github.com/jinzhu/gorm"

	"bean/pkg/namespace/model"
	"bean/pkg/namespace/model/dto"
	mUser "bean/pkg/user/model"
	"bean/pkg/util"
)

type MembershipCreateHandler struct {
	ID *util.Identifier
	DB *gorm.DB
}

func (this MembershipCreateHandler) NamespaceMembershipCreate(
	ctx context.Context,
	input dto.NamespaceMembershipCreateInput,
	namespace *model.Namespace,
	user *mUser.User,
) (*dto.NamespaceMembershipCreateOutcome, error) {
	id, err := this.ID.ULID()
	if nil != err {
		return nil, err
	}

	version, err := this.ID.ULID()
	if nil != err {
		return nil, err
	}

	membership := &model.Membership{
		ID:          id,
		Version:     version,
		NamespaceID: namespace.ID,
		UserID:      user.ID,
		IsActive:    input.IsActive,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err = this.DB.Table("namespace_memberships").Create(&membership).Error
	if nil != err {
		return nil, err
	}

	return &dto.NamespaceMembershipCreateOutcome{
		Errors:     nil,
		Membership: membership,
	}, nil
}
