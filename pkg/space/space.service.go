package space

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"bean/components/claim"
	"bean/components/scalar"
	"bean/pkg/space/model"
	"bean/pkg/space/model/dto"
	"bean/pkg/util"
	"bean/pkg/util/api"
)

type SpaceService struct {
	bundle *SpaceBundle
}

func (this SpaceService) Load(ctx context.Context, id string) (*model.Space, error) {
	obj := &model.Space{}
	err := this.bundle.db.First(&obj, "id = ?", id).Error
	if nil != err {
		return nil, err
	}

	return obj, nil
}

func (this SpaceService) FindOne(ctx context.Context, filters dto.SpaceFilters) (*model.Space, error) {
	if nil != filters.ID {
		return this.Load(ctx, *filters.ID)
	} else if nil != filters.Domain {
		domain := &model.DomainName{}
		err := this.bundle.db.Where("value = ?", filters.Domain).First(&domain).Error
		if nil != err {
			return nil, err
		} else if !domain.IsActive {
			return nil, errors.New("domain name is not active")
		}

		return this.Load(ctx, domain.SpaceId)
	}

	return nil, nil
}

func (this *SpaceService) Create(tx *gorm.DB, in dto.SpaceCreateInput) (*dto.SpaceCreateOutcome, error) {
	space, err := this.create(tx, in)
	if nil != err {
		return nil, err
	} else {
		err := this.createRelationships(tx, space, in)
		if nil != err {
			return nil, err
		}

		return &dto.SpaceCreateOutcome{Errors: nil, Space: space}, nil
	}
}

func (this *SpaceService) create(tx *gorm.DB, in dto.SpaceCreateInput) (*model.Space, error) {
	space := &model.Space{
		ID:        this.bundle.id.MustULID(),
		Version:   this.bundle.id.MustULID(),
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

func (this *SpaceService) createRelationships(tx *gorm.DB, space *model.Space, in dto.SpaceCreateInput) error {
	if err := this.bundle.DomainNameService.createMultiple(tx, space, in); nil != err {
		return err
	}

	// space configuration
	{
		if err := this.bundle.ConfigService.CreateFeatures(tx, space, in); nil != err {
			return err
		}
	}

	claims := claim.ContextToPayload(tx.Statement.Context)
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

		if ownerRole, err := this.create(tx, ownerRoleInput); nil != err {
			return err
		} else {
			// membership of user -> organisation
			_, err = this.bundle.MemberService.doCreate(tx, space.ID, claims.UserId(), true)
			if nil != err {
				return err
			}

			// membership of user -> owner role
			_, err = this.bundle.MemberService.doCreate(tx, ownerRole.ID, claims.UserId(), true)
			if nil != err {
				return err
			}
		}
	}

	return nil
}

func (this SpaceService) Update(tx *gorm.DB, obj model.Space, in dto.SpaceUpdateInput) (*dto.SpaceCreateOutcome, error) {
	// check version for conflict
	if in.SpaceVersion != obj.Version {
		return nil, util.ErrorVersionConflict
	}

	if nil != in.Object.Language {
		obj.Language = *in.Object.Language
	}

	// change version
	obj.Version = this.bundle.id.MustULID()
	if err := tx.Save(obj).Error; nil != err {
		return nil, err
	}

	err := this.bundle.ConfigService.updateFeatures(tx, &obj, in)
	if nil != err {
		return nil, err
	}

	return &dto.SpaceCreateOutcome{Errors: nil, Space: &obj}, nil
}
