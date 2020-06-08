package handler

import (
	"context"
	"errors"

	"github.com/jinzhu/gorm"

	"bean/pkg/namespace/model"
	"bean/pkg/namespace/model/dto"
	"bean/pkg/util/connect"
)

type NamespaceQueryHandler struct {
	DB *gorm.DB
}

func (this NamespaceQueryHandler) Handle(ctx context.Context, filters dto.NamespaceFilters) (*model.Namespace, error) {
	if nil != filters.ID {
		return this.Load(ctx, *filters.ID)
	} else if nil != filters.Domain {
		domain := &model.DomainName{}
		err := this.DB.
			Table(connect.TableNamespaceDomains).
			Where("value = ?", filters.Domain).
			First(&domain).
			Error

		if nil != err {
			return nil, err
		} else if !domain.IsActive {
			return nil, errors.New("domain name is not active")
		}

		return this.Load(ctx, domain.NamespaceId)
	}

	return nil, nil
}

func (this NamespaceQueryHandler) Load(ctx context.Context, id string) (*model.Namespace, error) {
	obj := &model.Namespace{}
	err := this.DB.First(&obj, "id = ?", id).Error
	if nil != err {
		return nil, err
	}

	return obj, nil
}
