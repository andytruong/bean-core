package handler

import (
	"context"

	"github.com/jinzhu/gorm"

	"bean/pkg/namespace/model"
	"bean/pkg/namespace/model/dto"
	"bean/pkg/util"
	"bean/pkg/util/connect"
)

type MembershipUpdateHandler struct {
	ID *util.Identifier
	DB *gorm.DB
}

func (this MembershipUpdateHandler) NamespaceMembershipUpdate(
	ctx context.Context,
	input dto.NamespaceMembershipUpdateInput,
	membership *model.Membership,
) (*dto.NamespaceMembershipCreateOutcome, error) {
	// change version
	version, err := this.ID.ULID()
	if nil != err {
		return nil, err
	}

	membership.Version = version
	membership.IsActive = input.IsActive

	{
		err := this.DB.
			Table(connect.TableNamespaceMemberships).
			Save(&membership).
			Error

		if nil != err {
			return nil, err
		}
	}

	return &dto.NamespaceMembershipCreateOutcome{
		Errors:     nil,
		Membership: membership,
	}, nil
}
