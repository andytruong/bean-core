package space

import (
	"context"
	"errors"
	"fmt"
	"time"
	
	"gorm.io/gorm"
	
	"bean/components/scalar"
	"bean/components/util"
	"bean/components/util/connect"
	"bean/pkg/space/model"
	"bean/pkg/space/model/dto"
)

type MemberService struct {
	bundle *SpaceBundle
}

func (this MemberService) load(ctx context.Context, id string, version *string) (*model.Membership, error) {
	obj := &model.Membership{}
	err := this.bundle.db.WithContext(ctx).First(&obj, "id = ?", id).Error
	if nil != err {
		return nil, err
	} else if nil != version {
		if obj.Version != *version {
			return nil, util.ErrorVersionConflict
		}
	}
	
	return obj, nil
}

func (this MemberService) Find(first int, after *string, filters dto.MembershipsFilter) (*model.MembershipConnection, error) {
	if first > 100 {
		return nil, errors.New(util.ErrorQueryTooMuch.String())
	}
	
	con := &model.MembershipConnection{
		Nodes: []model.Membership{},
		PageInfo: model.MembershipInfo{
			EndCursor:   nil,
			HasNextPage: false,
			StartCursor: nil,
		},
	}
	
	query, err := this.findUnlimited(after, filters)
	if nil != err {
		return nil, err
	} else {
		err := query.Limit(first).Find(&con.Nodes).Error
		if nil != err {
			return nil, err
		}
		
		var counter int64
		if err := query.Count(&counter).Error; nil != err {
			return nil, err
		} else {
			con.PageInfo.HasNextPage = int(counter) > len(con.Nodes)
			
			if len(con.Nodes) > 0 {
				startEntity := con.Nodes[0]
				if startEntity.LoggedInAt != nil {
					con.PageInfo.StartCursor = scalar.NilString(startEntity.LoggedInAt.String())
				}
				
				endEntity := con.Nodes[len(con.Nodes)-1]
				if nil != endEntity.LoggedInAt {
					con.PageInfo.EndCursor = scalar.NilString(endEntity.LoggedInAt.String())
				}
			}
		}
	}
	
	return con, nil
}

func (this MemberService) findUnlimited(afterRaw *string, filters dto.MembershipsFilter) (*gorm.DB, error) {
	query := this.bundle.db.
		Where("space_memberships.user_id = ?", filters.UserID).
		Where("space_memberships.is_active = ?", filters.IsActive).
		Order("space_memberships.logged_in_at DESC")
	
	if nil != filters.Space {
		if nil != filters.Space.Title {
			query = query.
				Joins("INNER JOIN spaces ON space_memberships.space_id = spaces.id").
				Where("spaces.title LIKE ?", "%"+*filters.Space.Title+"%")
		}
		
		if nil != filters.Space.DomainName {
			query = query.
				Joins("INNER JOIN space_domains ON space_domains.space_id = spaces.id").
				Where("spaces.title value ?", "%"+*filters.Space.DomainName+"%")
		}
	}
	
	if nil != filters.ManagerId {
		query = query.
			Joins("INNER JOIN space_manager_edge ON space_manager_edge.user_member_id = space_memberships.id").
			Where("space_manager_edge.manager_member_id = ?", *filters.ManagerId)
	}
	
	// Pagination -> after
	if nil != afterRaw {
		after, err := connect.DecodeCursor(*afterRaw)
		
		if nil != err {
			return nil, err
		}
		
		if after.Entity != "Membership" {
			return nil, errors.New("unsupported sorting entity")
		}
		
		if after.Property != "logged_in_at" {
			return nil, errors.New("unsupported sorting property")
		}
		
		query = query.Where("space_memberships.logged_in_at > ?", after.Value)
	}
	
	return query, nil
}

func (this MemberService) Create(tx *gorm.DB, in dto.SpaceMembershipCreateInput) (*dto.SpaceMembershipCreateOutcome, error) {
	membership, err := this.doCreate(tx, in.SpaceID, in.UserID, in.IsActive)
	
	if nil != err {
		return nil, err
	}
	
	errorList, err := this.createRelationships(tx, membership, in.ManagerMemberIds)
	if nil != err {
		return nil, err
	}
	
	return &dto.SpaceMembershipCreateOutcome{
		Errors:     errorList,
		Membership: membership,
	}, nil
}

func (this MemberService) doCreate(tx *gorm.DB, spaceId string, userId string, isActive bool) (*model.Membership, error) {
	membership := &model.Membership{
		ID:        this.bundle.id.MustULID(),
		Version:   this.bundle.id.MustULID(),
		SpaceID:   spaceId,
		UserID:    userId,
		IsActive:  isActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	
	if err := tx.Create(&membership).Error; nil != err {
		return nil, err
	}
	
	return membership, nil
}

func (this MemberService) createRelationships(tx *gorm.DB, obj *model.Membership, managerMemberIds []string) ([]*util.Error, error) {
	if len(managerMemberIds) > this.bundle.config.Manager.MaxNumberOfManager {
		return util.NewErrors(util.ErrorQueryTooMuch, []string{"input", "managerMemberIds"}, "exceeded limitation"), nil
	}
	
	// validate manager in same space
	{
		var counter int64
		err := tx.
			Model(&model.Membership{}).
			Where("space_id = ?", obj.SpaceID).
			Where("id IN (?)", managerMemberIds).
			Where("is_active = ?", true).
			Count(&counter).
			Error
		
		if nil != err {
			return nil, err
		} else if int(counter) != len(managerMemberIds) {
			return util.NewErrors(util.ErrorQueryTooMuch, []string{"input", "managerMemberIds"}, "one ore more IDs are invalid"), nil
		}
	}
	
	// create relationship with managers
	for _, managerMemberId := range managerMemberIds {
		err := this.createRelationship(tx, obj, managerMemberId)
		if nil != err {
			return nil, err
		}
	}
	
	return nil, nil
}

func (this MemberService) createRelationship(tx *gorm.DB, obj *model.Membership, managerMemberId string) error {
	relationship := model.ManagerRelationship{
		ID:              this.bundle.id.MustULID(),
		Version:         this.bundle.id.MustULID(),
		UserMemberId:    obj.ID,
		ManagerMemberId: managerMemberId,
		IsActive:        true,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	
	return tx.Save(&relationship).Error
}

func (this MemberService) Update(tx *gorm.DB, in dto.SpaceMembershipUpdateInput, obj *model.Membership) (*dto.SpaceMembershipCreateOutcome, error) {
	obj.Version = this.bundle.id.MustULID()
	obj.IsActive = in.IsActive
	
	err := tx.Save(&obj).Error
	if nil != err {
		return nil, err
	} else {
		// TODO: remove manager
		// …
		
		// TODO: add manager
		// …
	}
	
	return &dto.SpaceMembershipCreateOutcome{Errors: nil, Membership: obj}, nil
}

func (this MemberService) FindRoles(ctx context.Context, userId string, spaceId string) ([]*model.Space, error) {
	var roles []*model.Space
	
	err := this.bundle.db.
		WithContext(ctx).
		Joins(
			fmt.Sprintf(
				"INNER JOIN %s ON %s.space_id = %s.id AND %s.user_id = ?",
				connect.TableSpaceMemberships,
				connect.TableSpaceMemberships,
				connect.TableSpace,
				connect.TableSpaceMemberships,
			),
			userId,
		).
		Where("kind = ?", model.SpaceKindRole).
		Where("parent_id = ?", spaceId).
		Find(&roles).
		Error
	
	if nil != err {
		return nil, err
	}
	
	return roles, nil
}

func (this MemberService) UpdateLastLoginTime(db *gorm.DB, membership *model.Membership) error {
	membership.LoggedInAt = scalar.NilTime(time.Now())
	
	return db.Save(&membership).Error
}
