package namespace

import (
	"context"
	"errors"
	"time"

	"github.com/jinzhu/gorm"

	"bean/pkg/namespace/model"
	"bean/pkg/namespace/model/dto"
	"bean/pkg/util"
	"bean/pkg/util/api"
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


func (this *Core) Create(tx *gorm.DB, input dto.NamespaceCreateInput) (*dto.NamespaceCreateOutcome, error) {
	namespace, err := this.doCreate(tx, input)
	if nil != err {
		return nil, err
	} else {
		err := this.createRelationships(tx, namespace, input)
		if nil != err {
			return nil, err
		}

		return &dto.NamespaceCreateOutcome{
			Errors:    nil,
			Namespace: namespace,
		}, nil
	}
}

func (this *Core) doCreate(tx *gorm.DB, input dto.NamespaceCreateInput) (*model.Namespace, error) {
	namespace := &model.Namespace{
		ID:        this.bean.id.MustULID(),
		Version:   this.bean.id.MustULID(),
		ParentID:  input.Context.NamespaceID,
		Kind:      input.Object.Kind,
		Title:     *input.Object.Title,
		Language:  input.Object.Language,
		IsActive:  input.Object.IsActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := tx.Create(&namespace).Error; nil != err {
		return nil, err
	}

	return namespace, nil
}

func (this *Core) createRelationships(tx *gorm.DB, namespace *model.Namespace, input dto.NamespaceCreateInput) error {
	if err := this.bean.CoreDomainName.createMultiple(tx, namespace, input); nil != err {
		return err
	}

	if err := this.createFeatures(tx, namespace, input); nil != err {
		return err
	}

	// setup roles
	//  - create 'owner' role for the new namespace
	//  - grant 'owner' role to actor
	ownerRoleInput := dto.NamespaceCreateInput{
		Object: dto.NamespaceCreateInputObject{
			Kind:     model.NamespaceKindRole,
			Title:    util.NilString("owner"),
			Language: api.LanguageDefault,
			IsActive: true,
		},
		Context: input.Context,
	}

	ownerRoleInput.Context.NamespaceID = util.NilString(namespace.ID)
	if ownerRole, err := this.doCreate(tx, ownerRoleInput); nil != err {
		return err
	} else {
		// membership of user -> organisation
		_, err = this.bean.CoreMember.doCreate(tx, namespace.ID, input.Context.UserID, true)
		if nil != err {
			return err
		}

		// membership of user -> owner role
		_, err = this.bean.CoreMember.doCreate(tx, ownerRole.ID, input.Context.UserID, true)
		if nil != err {
			return err
		}
	}

	return nil
}

func (this *Core) createFeatures(tx *gorm.DB, namespace *model.Namespace, input dto.NamespaceCreateInput) error {
	if input.Object.Features.Register {
		return this.createFeature(tx, namespace, "default", "register", []byte("true"))
	} else {
		return this.createFeature(tx, namespace, "default", "register", []byte("false"))
	}
}

func (this *Core) createFeature(
	tx *gorm.DB,
	namespace *model.Namespace, bucket string, key string, value []byte,
) error {
	config := model.NamespaceConfig{
		Id:          this.bean.id.MustULID(),
		Version:     this.bean.id.MustULID(),
		NamespaceId: namespace.ID,
		Bucket:      bucket,
		Key:         key,
		Value:       value,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	return tx.Table(connect.TableNamespaceConfig).Create(&config).Error
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

	err := this.updateFeatures(tx, obj, in)
	if nil != err {
		return nil, err
	}

	return util.NilBool(true), nil
}

func (this *Core) updateFeatures(tx *gorm.DB, obj *model.Namespace, in dto.NamespaceUpdateInput) error {
	if nil != in.Object.Features.Register {
		if *in.Object.Features.Register {
			return this.updateFeature(tx, obj, "default", "register", []byte("true"))
		} else {
			return this.updateFeature(tx, obj, "default", "register", []byte("false"))
		}
	}

	return nil
}

func (this *Core) updateFeature(
	tx *gorm.DB,
	obj *model.Namespace, bucket string, key string, value []byte,
) error {
	return tx.
		Table(connect.TableNamespaceConfig).
		Where("namespace_id = ? AND bucket = ? AND key = ?", obj.ID, bucket, key).
		Update(&model.NamespaceConfig{
			Version:   this.bean.id.MustULID(),
			Value:     value,
			UpdatedAt: time.Now(),
		}).
		Error
}
