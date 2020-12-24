package space

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"bean/components/claim"
	"bean/components/scalar"
	"bean/components/util"
	"bean/components/util/connect"
	"bean/pkg/infra/api"
	"bean/pkg/space/model"
	"bean/pkg/space/model/dto"
)

type SpaceService struct {
	bundle *SpaceBundle
}

func (service SpaceService) Load(ctx context.Context, id string) (*model.Space, error) {
	obj := &model.Space{}
	err := service.bundle.db.WithContext(ctx).First(&obj, "id = ?", id).Error
	if nil != err {
		return nil, err
	}

	return obj, nil
}

func (service SpaceService) FindOne(ctx context.Context, filters dto.SpaceFilters) (*model.Space, error) {
	if nil != filters.ID {
		return service.Load(ctx, *filters.ID)
	} else if nil != filters.Domain {
		domain := &model.DomainName{}
		err := service.bundle.db.WithContext(ctx).Where("value = ?", filters.Domain).First(&domain).Error
		if nil != err {
			return nil, err
		} else if !domain.IsActive {
			return nil, errors.New("domain name is not active")
		}

		return service.Load(ctx, domain.SpaceId)
	}

	return nil, nil
}

func (service *SpaceService) Create(ctx context.Context, in dto.SpaceCreateInput) (*dto.SpaceCreateOutcome, error) {
	space, err := service.create(ctx, in)
	if nil != err {
		return nil, err
	} else {
		err := service.createRelationships(ctx, space, in)
		if nil != err {
			return nil, err
		}

		return &dto.SpaceCreateOutcome{Errors: nil, Space: space}, nil
	}
}

func (service *SpaceService) create(ctx context.Context, in dto.SpaceCreateInput) (*model.Space, error) {
	tx := connect.ContextToDB(ctx)
	space := &model.Space{
		ID:        service.bundle.idr.MustULID(),
		Version:   service.bundle.idr.MustULID(),
		Kind:      in.Object.Kind,
		Title:     *in.Object.Title,
		Language:  in.Object.Language,
		IsActive:  in.Object.IsActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if nil != in.Object.ParentId {
		space.ParentID = in.Object.ParentId
	} else {
		claims := claim.ContextToPayload(tx.Statement.Context)
		if nil != claims {
			parentSpaceId := claims.SpaceId()
			space.ParentID = &parentSpaceId
		}
	}

	if err := tx.Table("spaces").Create(&space).Error; nil != err {
		return nil, err
	}

	return space, nil
}

func (service *SpaceService) createRelationships(ctx context.Context, space *model.Space, in dto.SpaceCreateInput) error {
	tx := connect.ContextToDB(ctx)
	if err := service.bundle.domainNameService.createMultiple(tx, space, in); nil != err {
		return err
	}

	// space configuration
	{
		if err := service.bundle.configService.CreateFeatures(tx, space, in); nil != err {
			return err
		}
	}

	claims := claim.ContextToPayload(ctx)
	if nil != claims {
		// setup roles
		//  - create 'owner' role for the new space
		//  - grant 'owner' role to actor
		ownerRoleInput := dto.SpaceCreateInput{
			Object: dto.SpaceCreateInputObject{
				ParentId: &space.ID,
				Kind:     model.SpaceKindRole,
				Title:    scalar.NilString("owner"),
				Language: api.LanguageDefault,
				IsActive: true,
			},
		}

		if ownerRole, err := service.create(ctx, ownerRoleInput); nil != err {
			return err
		} else {
			// membership of user -> organisation
			_, err = service.bundle.MemberService.doCreate(tx, space.ID, claims.UserId(), true)
			if nil != err {
				return err
			}

			// membership of user -> owner role
			_, err = service.bundle.MemberService.doCreate(tx, ownerRole.ID, claims.UserId(), true)
			if nil != err {
				return err
			}
		}
	}

	return nil
}

func (service SpaceService) Update(tx *gorm.DB, obj model.Space, in dto.SpaceUpdateInput) (*dto.SpaceCreateOutcome, error) {
	// check version for conflict
	if in.SpaceVersion != obj.Version {
		return nil, util.ErrorVersionConflict
	}

	if nil != in.Object.Language {
		obj.Language = *in.Object.Language
	}

	// change version
	obj.Version = service.bundle.idr.MustULID()
	if err := tx.Save(obj).Error; nil != err {
		return nil, err
	}

	err := service.bundle.configService.updateFeatures(tx, &obj, in)
	if nil != err {
		return nil, err
	}

	return &dto.SpaceCreateOutcome{Errors: nil, Space: &obj}, nil
}
