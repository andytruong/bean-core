package namespace

import (
	"context"
	"database/sql"

	"github.com/jinzhu/gorm"

	"bean/pkg/namespace/handler"
	"bean/pkg/namespace/model"
	"bean/pkg/namespace/model/dto"
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

func (this NamespaceQueryResolver) Namespace(ctx context.Context, id string) (*model.Namespace, error) {
	obj := &model.Namespace{}
	err := this.db.First(&obj, "id = ?", id).Error
	if nil != err {
		return nil, err
	}

	return obj, nil
}

func (this NamespaceMutationResolver) NamespaceCreate(ctx context.Context, input dto.NamespaceCreateInput) (*dto.NamespaceCreateOutcome, error) {
	hdl := handler.NamespaceCreateHandler{ID: this.id}
	txn := this.db.BeginTx(ctx, &sql.TxOptions{})
	out, err := hdl.Create(txn, input)

	if nil != err {
		txn.Rollback()

		return nil, err
	} else {
		return out, txn.Commit().Error
	}
}

func (this NamespaceModelResolver) DomainNames(ctx context.Context, namespace *model.Namespace) (*model.DomainNames, error) {
	out := &model.DomainNames{
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
			out.Primary = domainName
		} else {
			out.Secondary = append(out.Secondary, domainName)
		}
	}

	return out, nil
}
