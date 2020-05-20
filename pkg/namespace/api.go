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
	outcome, err := sv.Create(tx, input)

	if nil != err {
		tx.Rollback()

		return nil, err
	} else {
		return outcome, tx.Commit().Error
	}
}

func (this NamespaceModelResolver) DomainNames(ctx context.Context, namespace *model.Namespace) (*model.DomainNames, error) {
	outcome := &model.DomainNames{
		Primary:   nil,
		Secondary: nil,
	}

	var domainNames []*model.DomainName
	err := this.db.
		Table("namespace_domains").
		Where("namespace_id = ?", namespace.ID).
		Find(&domainNames).
		Error
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
