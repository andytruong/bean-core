package namespace

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"bean/pkg/namespace/model"
	"bean/pkg/namespace/model/dto"
	"bean/pkg/util"
	"bean/pkg/util/api"
	"bean/pkg/util/api/scalar"
	"bean/pkg/util/connect"
)

type Core struct {
	bean *NamespaceBean
}

func (this Core) Load(ctx context.Context, id string) (*model.Namespace, error) {
	obj := &model.Namespace{}
	err := this.bean.db.First(&obj, "id = ?", id).Error
	if nil != err {
		return nil, err
	}

	return obj, nil
}

func (this Core) Find(ctx context.Context, filters dto.NamespaceFilters) (*model.Namespace, error) {
	if nil != filters.ID {
		return this.Load(ctx, *filters.ID)
	} else if nil != filters.Domain {
		domain := &model.DomainName{}
		err := this.bean.db.
			Table(connect.TableNamespaceDomains).
			Where("value = ?", filters.Domain).
			First(&domain).
			Error

		if nil != err {
			return nil, err
		} else if !domain.IsActive {
			return nil, errors.New("domain name is not active")
		}

		return this.Load(ctx, domain.NamespaceId)
	}

	return nil, nil
}

func (this *Core) Create(tx *gorm.DB, in dto.NamespaceCreateInput) (*dto.NamespaceCreateOutcome, error) {
	namespace, err := this.create(tx, in)
	if nil != err {
		return nil, err
	} else {
		err := this.createRelationships(tx, namespace, in)
		if nil != err {
			return nil, err
		}

		return &dto.NamespaceCreateOutcome{Errors: nil, Namespace: namespace}, nil
	}
}

func (this *Core) create(tx *gorm.DB, in dto.NamespaceCreateInput) (*model.Namespace, error) {
	namespace := &model.Namespace{
		ID:        this.bean.id.MustULID(),
		Version:   this.bean.id.MustULID(),
		Kind:      in.Object.Kind,
		Title:     *in.Object.Title,
		Language:  in.Object.Language,
		IsActive:  in.Object.IsActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if nil != in.Object.ParentId {
		namespace.ParentID = in.Object.ParentId
	} else {
		claims := util.CxtKeyClaims.Actor(tx.Statement.Context)
		if nil != claims {
			parentNamespaceId := claims.NamespaceId()
			namespace.ParentID = &parentNamespaceId
		}
	}

	if err := tx.Create(&namespace).Error; nil != err {
		return nil, err
	}

	return namespace, nil
}

func (this *Core) createRelationships(tx *gorm.DB, namespace *model.Namespace, in dto.NamespaceCreateInput) error {
	if err := this.bean.CoreDomainName.createMultiple(tx, namespace, in); nil != err {
		return err
	}

	// namespace configuration
	{
		if err := this.bean.CoreConfig.CreateFeatures(tx, namespace, in); nil != err {
			return err
		}
	}

	claims := util.CxtKeyClaims.Actor(tx.Statement.Context)
	if nil != claims {
		// setup roles
		//  - create 'owner' role for the new namespace
		//  - grant 'owner' role to actor
		ownerRoleInput := dto.NamespaceCreateInput{
			Object: dto.NamespaceCreateInputObject{
				ParentId: &namespace.ID,
				Kind:     model.NamespaceKindRole,
				Title:    scalar.NilString("owner"),
				Language: api.LanguageDefault,
				IsActive: true,
			},
		}

		if ownerRole, err := this.create(tx, ownerRoleInput); nil != err {
			return err
		} else {
			// membership of user -> organisation
			_, err = this.bean.CoreMember.doCreate(tx, namespace.ID, claims.UserId(), true)
			if nil != err {
				return err
			}

			// membership of user -> owner role
			_, err = this.bean.CoreMember.doCreate(tx, ownerRole.ID, claims.UserId(), true)
			if nil != err {
				return err
			}
		}
	}

	return nil
}

func (this Core) Update(tx *gorm.DB, obj *model.Namespace, in dto.NamespaceUpdateInput) (*bool, error) {
	// check version for conflict
	if in.NamespaceVersion != obj.Version {
		return nil, util.ErrorVersionConflict
	}

	if nil != in.Object.Language {
		obj.Language = *in.Object.Language
	}

	// change version
	obj.Version = this.bean.id.MustULID()
	if err := tx.Save(obj).Error; nil != err {
		return nil, err
	}

	err := this.bean.CoreConfig.updateFeatures(tx, obj, in)
	if nil != err {
		return nil, err
	}

	return scalar.NilBool(true), nil
}
