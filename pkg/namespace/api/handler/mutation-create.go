package handler

import (
	"time"

	"github.com/jinzhu/gorm"

	"bean/pkg/namespace/model"
	"bean/pkg/namespace/model/dto"
	"bean/pkg/util"
	"bean/pkg/util/api"
	"bean/pkg/util/connect"
)

type NamespaceCreateHandler struct {
	ID *util.Identifier
}

func (this *NamespaceCreateHandler) Create(tx *gorm.DB, input dto.NamespaceCreateInput) (*dto.NamespaceCreateOutcome, error) {
	namespace, err := this.create(tx, input)
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

func (this *NamespaceCreateHandler) create(tx *gorm.DB, input dto.NamespaceCreateInput) (*model.Namespace, error) {
	namespace := &model.Namespace{
		ID:        this.ID.MustULID(),
		Version:   this.ID.MustULID(),
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

func (this *NamespaceCreateHandler) createRelationships(tx *gorm.DB, namespace *model.Namespace, input dto.NamespaceCreateInput) error {
	if err := this.createDomains(tx, namespace, input); nil != err {
		return err
	}

	if err := this.createFeatures(tx, namespace, input); nil != err {
		return err
	}

	// setup roles
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
	if ownerRole, err := this.create(tx, ownerRoleInput); nil != err {
		return err
	} else {
		// membership of user -> organisation
		if err := this.createMembership(tx, namespace, input); nil != err {
			return err
		}

		// membership of user -> owner role
		if err := this.createMembership(tx, ownerRole, input); nil != err {
			return err
		}
	}

	return nil
}

func (this *NamespaceCreateHandler) createDomains(tx *gorm.DB, namespace *model.Namespace, input dto.NamespaceCreateInput) error {
	if nil != input.Object.DomainNames.Primary {
		err := this.createDomain(tx, namespace, input.Object.DomainNames.Primary, true)
		if nil != err {
			return err
		}
	}

	if nil != input.Object.DomainNames.Secondary {
		for _, in := range input.Object.DomainNames.Secondary {
			err := this.createDomain(tx, namespace, in, false)
			if nil != err {
				return err
			}
		}
	}

	return nil
}

func (this *NamespaceCreateHandler) createDomain(tx *gorm.DB, namespace *model.Namespace, input *dto.DomainNameInput, isPrimary bool) error {
	domain := model.DomainName{
		ID:          this.ID.MustULID(),
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

func (this *NamespaceCreateHandler) createMembership(tx *gorm.DB, namespace *model.Namespace, input dto.NamespaceCreateInput) error {
	membership := model.Membership{
		ID:          this.ID.MustULID(),
		Version:     this.ID.MustULID(),
		NamespaceID: namespace.ID,
		UserID:      input.Context.UserID,
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := tx.Table(connect.TableNamespaceMemberships).Create(membership).Error; nil != err {
		return err
	}

	return nil
}

func (this *NamespaceCreateHandler) createFeatures(tx *gorm.DB, namespace *model.Namespace, input dto.NamespaceCreateInput) error {
	if input.Object.Features.Register {
		return this.createFeature(tx, namespace, "default", "register", []byte("true"))
	} else {
		return this.createFeature(tx, namespace, "default", "register", []byte("false"))
	}
}

func (this *NamespaceCreateHandler) createFeature(
	tx *gorm.DB,
	namespace *model.Namespace, bucket string, key string, value []byte,
) error {
	config := model.NamespaceConfig{
		Id:          this.ID.MustULID(),
		Version:     this.ID.MustULID(),
		NamespaceId: namespace.ID,
		Bucket:      bucket,
		Key:         key,
		Value:       value,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	return tx.Table(connect.TableNamespaceConfig).Create(&config).Error
}
