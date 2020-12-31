package space

import (
	"context"
	"time"

	"bean/components/connect"
	"bean/pkg/space/model"
	"bean/pkg/space/model/dto"
)

type DomainNameService struct {
	bundle *SpaceBundle
}

func (service *DomainNameService) createMultiple(ctx context.Context, space *model.Space, in dto.SpaceCreateInput) error {
	if nil == in.Object.DomainNames {
		return nil
	}

	if nil != in.Object.DomainNames.Primary {
		err := service.create(ctx, space, in.Object.DomainNames.Primary, true)
		if nil != err {
			return err
		}
	}

	if nil != in.Object.DomainNames.Secondary {
		for _, in := range in.Object.DomainNames.Secondary {
			err := service.create(ctx, space, in, false)
			if nil != err {
				return err
			}
		}
	}

	return nil
}

func (service *DomainNameService) create(ctx context.Context, space *model.Space, in *dto.DomainNameInput, isPrimary bool) error {
	domain := model.DomainName{
		ID:         service.bundle.idr.MustULID(),
		SpaceId:    space.ID,
		IsVerified: *in.Verified,
		Value:      *in.Value,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		IsPrimary:  isPrimary,
		IsActive:   *in.IsActive,
	}

	return connect.ContextToDB(ctx).Create(&domain).Error
}

func (service *DomainNameService) Find(ctx context.Context, space *model.Space) (*model.DomainNames, error) {
	out := &model.DomainNames{
		Primary:   nil,
		Secondary: nil,
	}

	var domainNames []*model.DomainName

	err := connect.ContextToDB(ctx).
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
