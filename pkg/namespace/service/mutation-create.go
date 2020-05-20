package service

import (
	"time"

	"github.com/jinzhu/gorm"

	"bean/pkg/namespace/model"
	"bean/pkg/namespace/model/dto"
	"bean/pkg/util"
)

type NamespaceCreateAPI struct {
	ID *util.Identifier
}

func (this *NamespaceCreateAPI) Create(tx *gorm.DB, input dto.NamespaceCreateInput) (*dto.NamespaceCreateOutcome, error) {
	namespace := &model.Namespace{
		Title:     input.Object.Title,
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

	return &dto.NamespaceCreateOutcome{
		Errors:    nil,
		Namespace: namespace,
	}, nil
}

func (this *NamespaceCreateAPI) createDomains(tx *gorm.DB, namespace *model.Namespace, input dto.NamespaceCreateInput) error {
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

func (this *NamespaceCreateAPI) createDomain(tx *gorm.DB, namespace *model.Namespace, input *dto.DomainNameInput, isPrimary bool) error {
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

func (this *NamespaceCreateAPI) createMembership(tx *gorm.DB, namespace *model.Namespace, input dto.NamespaceCreateInput) error {
	if nil != input.Context {
		id, err := this.ID.ULID()
		if nil != err {
			return err
		}

		membership := model.Membership{
			ID:          id,
			NamespaceID: namespace.ID,
			UserID:      input.Context.UserID,
			IsActive:    false,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		if err := tx.Table("namespace_memberships").Create(membership).Error; nil != err {
			return err
		}
	}

	return nil
}
