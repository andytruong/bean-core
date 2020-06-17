package namespace

import (
	"time"

	"github.com/jinzhu/gorm"

	"bean/pkg/namespace/model"
	"bean/pkg/namespace/model/dto"
	"bean/pkg/util/connect"
)

type CoreDomainName struct {
	bean *NamespaceBean
}

func (this *CoreDomainName) createMultiple(tx *gorm.DB, namespace *model.Namespace, input dto.NamespaceCreateInput) error {
	if nil != input.Object.DomainNames.Primary {
		err := this.create(tx, namespace, input.Object.DomainNames.Primary, true)
		if nil != err {
			return err
		}
	}

	if nil != input.Object.DomainNames.Secondary {
		for _, in := range input.Object.DomainNames.Secondary {
			err := this.create(tx, namespace, in, false)
			if nil != err {
				return err
			}
		}
	}

	return nil
}

func (this *CoreDomainName) create(tx *gorm.DB, namespace *model.Namespace, input *dto.DomainNameInput, isPrimary bool) error {
	domain := model.DomainName{
		ID:          this.bean.id.MustULID(),
		NamespaceId: namespace.ID,
		IsVerified:  *input.Verified,
		Value:       *input.Value,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		IsPrimary:   isPrimary,
		IsActive:    *input.IsActive,
	}

	return tx.Table("namespace_domains").Create(&domain).Error
}

func (this *CoreDomainName) Find(namespace *model.Namespace) (*model.DomainNames, error) {
	out := &model.DomainNames{
		Primary:   nil,
		Secondary: nil,
	}

	var domainNames []*model.DomainName
	err := this.bean.db.
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
