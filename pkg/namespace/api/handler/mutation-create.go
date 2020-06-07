package handler

import (
	"time"

	"github.com/jinzhu/gorm"

	"bean/pkg/namespace/model"
	"bean/pkg/namespace/model/dto"
	"bean/pkg/util"
)

type NamespaceCreateHandler struct {
	ID *util.Identifier
}

func (this *NamespaceCreateHandler) Create(tx *gorm.DB, input dto.NamespaceCreateInput) (*dto.NamespaceCreateOutcome, error) {
	namespace := &model.Namespace{
		Title:     input.Object.Title,
		Language:  input.Object.Language,
		IsActive:  input.Object.IsActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Generate user Identifier.
	if id, err := this.ID.ULID(); nil != err {
		return nil, err
	} else if version, err := this.ID.ULID(); nil != err {
		return nil, err
	} else {
		namespace.ID = id
		namespace.Version = version
	}

	if err := tx.Create(&namespace).Error; nil != err {
		return nil, err
	}

	// create domain
	if err := this.createDomains(tx, namespace, input); nil != err {
		return nil, err
	}

	// create membership
	if err := this.createMembership(tx, namespace, input); nil != err {
		return nil, err
	}

	// save features
	if err := this.createFeatures(tx, namespace, input); nil != err {
		return nil, err
	}

	return &dto.NamespaceCreateOutcome{
		Errors:    nil,
		Namespace: namespace,
	}, nil
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
	id, err := this.ID.ULID()
	if nil != err {
		return err
	}

	domain := model.DomainName{
		ID:          id,
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
	if nil != input.Context {
		id, err := this.ID.ULID()
		if nil != err {
			return err
		}

		version, err := this.ID.ULID()
		if nil != err {
			return err
		}

		membership := model.Membership{
			ID:          id,
			Version:     version,
			NamespaceID: namespace.ID,
			UserID:      input.Context.UserID,
			IsActive:    true,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		if err := tx.Table("namespace_memberships").Create(membership).Error; nil != err {
			return err
		}
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
	id, err := this.ID.ULID()
	if nil != err {
		return err
	}

	version, err := this.ID.ULID()
	if nil != err {
		return err
	}

	config := model.NamespaceConfig{
		Id:          id,
		Version:     version,
		NamespaceId: namespace.ID,
		Bucket:      bucket,
		Key:         key,
		Value:       value,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	return tx.Table("namespace_config").Create(&config).Error
}
