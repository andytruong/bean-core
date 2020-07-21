package namespace

import (
	"time"

	"gorm.io/gorm"

	"bean/pkg/namespace/model"
	"bean/pkg/namespace/model/dto"
)

type CoreDomainName struct {
	bean *NamespaceBean
}

func (this *CoreDomainName) createMultiple(tx *gorm.DB, namespace *model.Namespace, in dto.NamespaceCreateInput) error {
	if nil != in.Object.DomainNames.Primary {
		err := this.create(tx, namespace, in.Object.DomainNames.Primary, true)
		if nil != err {
			return err
		}
	}

	if nil != in.Object.DomainNames.Secondary {
		for _, in := range in.Object.DomainNames.Secondary {
			err := this.create(tx, namespace, in, false)
			if nil != err {
				return err
			}
		}
	}

	return nil
}

func (this *CoreDomainName) create(tx *gorm.DB, namespace *model.Namespace, in *dto.DomainNameInput, isPrimary bool) error {
	domain := model.DomainName{
		ID:          this.bean.id.MustULID(),
		NamespaceId: namespace.ID,
		IsVerified:  *in.Verified,
		Value:       *in.Value,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		IsPrimary:   isPrimary,
		IsActive:    *in.IsActive,
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
