package namespace

import (
	"context"
	"database/sql"

	"github.com/jinzhu/gorm"

	"bean/pkg/namespace/model"
	"bean/pkg/namespace/model/dto"
	"bean/pkg/namespace/service"
	"bean/pkg/util"
)

type (
	NamespaceMutationResolver struct {
		db *gorm.DB
		id *util.Identifier
	}

	NamespaceQueryResolver struct {
		db *gorm.DB
	}

	NamespaceModelResolver struct {
		db *gorm.DB
	}
)

func (this NamespaceMutationResolver) NamespaceCreate(ctx context.Context, input dto.NamespaceCreateInput) (*dto.NamespaceCreateOutcome, error) {
	sv := service.NamespaceCreateAPI{ID: this.id}
	tx := this.db.BeginTx(ctx, &sql.TxOptions{})

	return sv.Create(tx, input)
}

func (this NamespaceModelResolver) DomainNames(ctx context.Context, obj *model.Namespace) (*model.DomainNames, error) {
	outcome := &model.DomainNames{
		Primary:   nil,
		Secondary: nil,
	}

	domainNames := []*model.DomainName{}
	err := this.db.Where(model.DomainName{NamespaceId: obj.ID}).Find(&domainNames).Error
	if nil != err {
		return nil, err
	}

	for _, domainName := range domainNames {
		if domainName.IsPrimary {
			outcome.Primary = domainName
		} else {
			outcome.Secondary = append(outcome.Secondary, domainName)
		}
	}

	return outcome, nil
}
