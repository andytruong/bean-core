package space

import (
	"time"

	"gorm.io/gorm"

	"bean/pkg/space/model"
	"bean/pkg/space/model/dto"
)

type DomainNameService struct {
	bean *SpaceBundle
}

func (this *DomainNameService) createMultiple(tx *gorm.DB, space *model.Space, in dto.SpaceCreateInput) error {
	if nil != in.Object.DomainNames.Primary {
		err := this.create(tx, space, in.Object.DomainNames.Primary, true)
		if nil != err {
			return err
		}
	}

	if nil != in.Object.DomainNames.Secondary {
		for _, in := range in.Object.DomainNames.Secondary {
			err := this.create(tx, space, in, false)
			if nil != err {
				return err
			}
		}
	}

	return nil
}

func (this *DomainNameService) create(tx *gorm.DB, space *model.Space, in *dto.DomainNameInput, isPrimary bool) error {
	domain := model.DomainName{
		ID:         this.bean.id.MustULID(),
		SpaceId:    space.ID,
		IsVerified: *in.Verified,
		Value:      *in.Value,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		IsPrimary:  isPrimary,
		IsActive:   *in.IsActive,
	}

	return tx.Table("space_domains").Create(&domain).Error
}

func (this *DomainNameService) Find(space *model.Space) (*model.DomainNames, error) {
	out := &model.DomainNames{
		Primary:   nil,
		Secondary: nil,
	}

	var domainNames []*model.DomainName
	err := this.bean.db.
		Where("space_id = ?", space.ID).
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