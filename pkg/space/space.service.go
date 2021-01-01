package space

import (
	"context"
	"errors"
	"time"

	"bean/components/claim"
	"bean/components/connect"
	"bean/components/scalar"
	"bean/components/util"
	"bean/pkg/infra/api"
	"bean/pkg/space/model"
	"bean/pkg/space/model/dto"
)

type SpaceService struct {
	bundle *SpaceBundle
}

func (srv SpaceService) Load(ctx context.Context, id string) (*model.Space, error) {
	obj := &model.Space{}
	err := connect.ContextToDB(ctx).First(&obj, "id = ?", id).Error
	if nil != err {
		return nil, err
	}

	return obj, nil
}

func (srv SpaceService) FindOne(ctx context.Context, filters dto.SpaceFilters) (*model.Space, error) {
	if nil != filters.ID {
		return srv.Load(ctx, *filters.ID)
	} else if nil != filters.Domain {
		domain := &model.DomainName{}
		err := connect.ContextToDB(ctx).Where("value = ?", filters.Domain).First(&domain).Error
		if nil != err {
			return nil, err
		} else if !domain.IsActive {
			return nil, errors.New("domain name is not active")
		}

		return srv.Load(ctx, domain.SpaceId)
	}

	return nil, nil
}

func (srv *SpaceService) Create(ctx context.Context, in dto.SpaceCreateInput) (*dto.SpaceCreateOutcome, error) {
	space, err := srv.create(ctx, in)
	if nil != err {
		return nil, err
	} else {
		err := srv.createRelationships(ctx, space, in)
		if nil != err {
			return nil, err
		}

		return &dto.SpaceCreateOutcome{Errors: nil, Space: space}, nil
	}
}

func (srv *SpaceService) create(ctx context.Context, in dto.SpaceCreateInput) (*model.Space, error) {
	db := connect.ContextToDB(ctx)
	space := &model.Space{
		ID:        srv.bundle.idr.MustULID(),
		Version:   srv.bundle.idr.MustULID(),
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
		claims := claim.ContextToPayload(db.Statement.Context)
		if nil != claims {
			parentSpaceId := claims.SpaceId()
			space.ParentID = &parentSpaceId
		}
	}

	if err := db.Create(&space).Error; nil != err {
		return nil, err
	}

	return space, nil
}

func (srv *SpaceService) createRelationships(ctx context.Context, space *model.Space, in dto.SpaceCreateInput) error {
	if err := srv.bundle.domainNameService.createMultiple(ctx, space, in); nil != err {
		return err
	}

	// space configuration
	{
		if err := srv.bundle.configService.CreateFeatures(ctx, space, in); nil != err {
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

		if ownerRole, err := srv.create(ctx, ownerRoleInput); nil != err {
			return err
		} else {
			// membership of user -> organisation
			_, err = srv.bundle.MemberService.doCreate(ctx, space.ID, claims.UserId(), true)
			if nil != err {
				return err
			}

			// membership of user -> owner role
			_, err = srv.bundle.MemberService.doCreate(ctx, ownerRole.ID, claims.UserId(), true)
			if nil != err {
				return err
			}
		}
	}

	return nil
}

func (srv SpaceService) Update(ctx context.Context, obj model.Space, in dto.SpaceUpdateInput) (*dto.SpaceCreateOutcome, error) {
	tx := connect.ContextToDB(ctx)

	// check version for conflict
	if in.SpaceVersion != obj.Version {
		return nil, util.ErrorVersionConflict
	}

	if nil != in.Object.Language {
		obj.Language = *in.Object.Language
	}

	// change version
	obj.Version = srv.bundle.idr.MustULID()
	if err := tx.Save(obj).Error; nil != err {
		return nil, err
	}

	err := srv.bundle.configService.updateFeatures(tx, &obj, in)
	if nil != err {
		return nil, err
	}

	return &dto.SpaceCreateOutcome{Errors: nil, Space: &obj}, nil
}
