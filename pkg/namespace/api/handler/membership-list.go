package handler

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"bean/pkg/namespace/model"
	"bean/pkg/namespace/model/dto"
	"bean/pkg/util"
	"bean/pkg/util/connect"
)

type MembershipsQueryHandler struct {
	DB     *gorm.DB
	Logger *zap.Logger
}

func (this MembershipsQueryHandler) Memberships(first int, after *string, filters dto.MembershipsFilter) (*model.MembershipConnection, error) {
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

	query, err := this.buildNoLimitQuery(after, filters)
	if nil != err {
		return nil, err
	} else {
		err := query.Limit(first).Find(&con.Nodes).Error
		if nil != err {
			return nil, err
		}

		counter := 0
		if err := query.Count(&counter).Error; nil != err {
			return nil, err
		} else {
			con.PageInfo.HasNextPage = counter > len(con.Nodes)

			if len(con.Nodes) > 0 {
				startEntity := con.Nodes[0]
				if startEntity.LoggedInAt != nil {
					con.PageInfo.StartCursor = util.NilString(startEntity.LoggedInAt.String())
				}

				endEntity := con.Nodes[len(con.Nodes)-1]
				if nil != endEntity.LoggedInAt {
					con.PageInfo.EndCursor = util.NilString(endEntity.LoggedInAt.String())
				}
			}
		}
	}

	return con, nil
}

func (this MembershipsQueryHandler) buildNoLimitQuery(
	afterRaw *string,
	filters dto.MembershipsFilter,
) (*gorm.DB, error) {
	query := this.DB.
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
