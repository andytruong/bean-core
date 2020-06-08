package handler

import (
	"github.com/jinzhu/gorm"

	"bean/pkg/namespace/model"
	"bean/pkg/namespace/model/dto"
	"bean/pkg/util"
	"bean/pkg/util/connect"
)

type MembershipUpdateHandler struct {
	ID *util.Identifier
}

func (this MembershipUpdateHandler) NamespaceMembershipUpdate(
	tx *gorm.DB,
	input dto.NamespaceMembershipUpdateInput,
	membership *model.Membership,
) (*dto.NamespaceMembershipCreateOutcome, error) {
	membership.Version = this.ID.MustULID()
	membership.IsActive = input.IsActive

	err := tx.
		Table(connect.TableNamespaceMemberships).
		Save(&membership).
		Error

	if nil != err {
		return nil, err
	} else {
		// TODO: remove manager
		// …

		// TODO: add manager
		// …
	}

	return &dto.NamespaceMembershipCreateOutcome{
		Errors:     nil,
		Membership: membership,
	}, nil
}
