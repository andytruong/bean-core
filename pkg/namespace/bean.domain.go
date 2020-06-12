package namespace

import (
	"time"

	"github.com/jinzhu/gorm"

	"bean/pkg/namespace/model"
	"bean/pkg/namespace/model/dto"
)

type NamespaceBeanDomainName struct {
	bean *NamespaceBean
}

func (this *NamespaceBeanDomainName) createMultiple(tx *gorm.DB, namespace *model.Namespace, input dto.NamespaceCreateInput) error {
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

func (this *NamespaceBeanDomainName) create(tx *gorm.DB, namespace *model.Namespace, input *dto.DomainNameInput, isPrimary bool) error {
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
