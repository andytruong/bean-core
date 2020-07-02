package namespace

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"bean/pkg/namespace/model"
	"bean/pkg/namespace/model/dto"
	mUser "bean/pkg/user/model"
	"bean/pkg/util"
	"bean/pkg/util/api/scalar"
	"bean/pkg/util/connect"
)

type CoreMember struct {
	bean     *NamespaceBean
	Resolver MembershipResolver
}

func (this CoreMember) Find(first int, after *string, filters dto.MembershipsFilter) (*model.MembershipConnection, error) {
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

func (this CoreMember) findUnlimited(afterRaw *string, filters dto.MembershipsFilter) (*gorm.DB, error) {
	query := this.bean.db.
		Table(connect.TableNamespaceMemberships).
		Where("namespace_memberships.user_id = ?", filters.UserID).
		Where("namespace_memberships.is_active = ?", filters.IsActive).
		Order("namespace_memberships.logged_in_at DESC")

	if nil != filters.Namespace {
		if nil != filters.Namespace.Title {
			query = query.
				Joins("INNER JOIN namespaces ON namespace_memberships.namespace_id = namespaces.id").
				Where("namespaces.title LIKE ?", "%"+*filters.Namespace.Title+"%")
		}

		if nil != filters.Namespace.DomainName {
			query = query.
				Joins("INNER JOIN namespace_domains ON namespace_domains.namespace_id = namespaces.id").
				Where("namespaces.title value ?", "%"+*filters.Namespace.DomainName+"%")
		}
	}

	if nil != filters.ManagerId {
		query = query.
			Joins("INNER JOIN namespace_manager_edge ON namespace_manager_edge.user_member_id = namespace_memberships.id").
			Where("namespace_manager_edge.manager_member_id = ?", *filters.ManagerId)
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

		query = query.Where("namespace_memberships.logged_in_at > ?", after.Value)
	}

	return query, nil
}

func (this CoreMember) Create(
	tx *gorm.DB,
	in dto.NamespaceMembershipCreateInput,
	namespace *model.Namespace,
	user *mUser.User,
) (*dto.NamespaceMembershipCreateOutcome, error) {
	membership, err := this.doCreate(tx, namespace.ID, user.ID, in.IsActive)

	if nil != err {
		return nil, err
	}

	errorList, err := this.createRelationships(tx, membership, in.ManagerMemberIds)
	if nil != err {
		return nil, err
	}

	return &dto.NamespaceMembershipCreateOutcome{
		Errors:     errorList,
		Membership: membership,
	}, nil
}

func (this CoreMember) doCreate(tx *gorm.DB, namespaceId string, userId string, isActive bool) (*model.Membership, error) {
	membership := &model.Membership{
		ID:          this.bean.id.MustULID(),
		Version:     this.bean.id.MustULID(),
		NamespaceID: namespaceId,
		UserID:      userId,
		IsActive:    isActive,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := tx.Table(connect.TableNamespaceMemberships).Create(&membership).Error; nil != err {
		return nil, err
	}

	return membership, nil
}

func (this CoreMember) createRelationships(tx *gorm.DB, obj *model.Membership, managerMemberIds []string) ([]*util.Error, error) {
	if len(managerMemberIds) > this.bean.genetic.Manager.MaxNumberOfManager {
		return util.NewErrors(util.ErrorQueryTooMuch, []string{"input", "managerMemberIds"}, "exceeded limitation"), nil
	}

	// validate manager in same namespace
	{
		var counter int64
		err := tx.
			Table(connect.TableNamespaceMemberships).
			Where("namespace_id = ?", obj.NamespaceID).
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

func (this CoreMember) createRelationship(tx *gorm.DB, obj *model.Membership, managerMemberId string) error {
	relationship := model.ManagerRelationship{
		ID:              this.bean.id.MustULID(),
		Version:         this.bean.id.MustULID(),
		UserMemberId:    obj.ID,
		ManagerMemberId: managerMemberId,
		IsActive:        true,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	return tx.Table(connect.TableManagerEdge).Save(&relationship).Error
}

func (this CoreMember) Update(tx *gorm.DB, in dto.NamespaceMembershipUpdateInput, obj *model.Membership) (*dto.NamespaceMembershipCreateOutcome, error) {
	obj.Version = this.bean.id.MustULID()
	obj.IsActive = in.IsActive

	err := tx.
		Table(connect.TableNamespaceMemberships).
		Save(&obj).
		Error

	if nil != err {
		return nil, err
	} else {
		// TODO: remove manager
		// …

		// TODO: add manager
		// …
	}

	return &dto.NamespaceMembershipCreateOutcome{Errors: nil, Membership: obj}, nil
}
