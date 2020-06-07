package handler

import (
	"context"

	"github.com/jinzhu/gorm"

	"bean/pkg/namespace/model"
	"bean/pkg/util/connect"
)

type DomainQueryHandler struct {
	DB *gorm.DB
}

func (this DomainQueryHandler) DomainNames(ctx context.Context, namespace *model.Namespace) (*model.DomainNames, error) {
	out := &model.DomainNames{
		Primary:   nil,
		Secondary: nil,
	}

	var domainNames []*model.DomainName
	err := this.DB.
		Table(connect.TableNamespaceDomains).
		Where("namespace_id = ?", namespace.ID).
		Find(&domainNames).
		Error
	if nil != err {
		return nil, err
	}

	for _, domainName := range domainNames {
		if domainName.IsPrimary {
			out.Primary = domainName
		} else {
			out.Secondary = append(out.Secondary, domainName)
		}
	}

	return out, nil
}
