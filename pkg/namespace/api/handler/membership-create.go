package handler

import (
	"time"

	"github.com/jinzhu/gorm"

	"bean/pkg/namespace/model"
	"bean/pkg/namespace/model/dto"
	mUser "bean/pkg/user/model"
	"bean/pkg/util"
	"bean/pkg/util/connect"
)

type MembershipCreateHandler struct {
	ID         *util.Identifier
	MaxManager int
}

func (this MembershipCreateHandler) NamespaceMembershipCreate(
	tx *gorm.DB,
	input dto.NamespaceMembershipCreateInput,
	namespace *model.Namespace,
	user *mUser.User,
) (*dto.NamespaceMembershipCreateOutcome, error) {
	membership := &model.Membership{
		ID:          this.ID.MustULID(),
		Version:     this.ID.MustULID(),
		NamespaceID: namespace.ID,
		UserID:      user.ID,
		IsActive:    input.IsActive,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := tx.Table(connect.TableNamespaceMemberships).Create(&membership).Error
	if nil != err {
		return nil, err
	} else if errors, err := this.createRelationships(tx, membership, input.ManagerMemberIds); nil != err {
		return nil, err
	} else {
		return &dto.NamespaceMembershipCreateOutcome{
			Errors:     errors,
			Membership: membership,
		}, nil
	}
}

func (this MembershipCreateHandler) createRelationships(tx *gorm.DB, membership *model.Membership, managerMemberIds []string) ([]*util.Error, error) {
	if len(managerMemberIds) > this.MaxManager {
		return util.NewErrors(util.ErrorQueryTooMuch, []string{"input", "managerMemberIds"}, "exceeded limitation"), nil
	}

	// validate manager in same namespace
	{
		counter := 0
		err := tx.
			Table(connect.TableNamespaceMemberships).
			Where("namespace_id = ?", membership.NamespaceID).
			Where("id IN (?)", managerMemberIds).
			Where("is_active = ?", true).
			Count(&counter).
			Error

		if nil != err {
			return nil, err
		} else if counter != len(managerMemberIds) {
			return util.NewErrors(util.ErrorQueryTooMuch, []string{"input", "managerMemberIds"}, "one ore more IDs are invalid"), nil
		}
	}

	// create relationship with managers
	for _, managerMemberId := range managerMemberIds {
		err := this.createRelationship(tx, membership, managerMemberId)
		if nil != err {
			return nil, err
		}
	}

	return nil, nil
}

func (this MembershipCreateHandler) createRelationship(tx *gorm.DB, membership *model.Membership, managerMemberId string) error {
	relationship := model.ManagerRelationship{
		ID:              this.ID.MustULID(),
		Version:         this.ID.MustULID(),
		UserMemberId:    membership.ID,
		ManagerMemberId: managerMemberId,
		IsActive:        true,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	return tx.Table(connect.TableManagerEdge).Save(&relationship).Error
}
