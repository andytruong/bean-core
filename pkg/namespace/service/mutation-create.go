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

func (this *NamespaceCreateAPI) Create(tx *gorm.DB, input *dto.NamespaceCreateInput) (*dto.NamespaceCreateOutcome, error) {
	namespace := &model.Namespace{
		Title:     input.Object.Title,
		IsActive:  input.Object.Status,
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

	if err := tx.Create(namespace).Error; nil != err {
		return nil, err
	}

	// create domain
	if err := this.createDomain(tx, namespace); nil != err {
		return nil, err
	}

	// create membership
	if err := this.createMembership(tx, namespace); nil != err {
		return nil, err
	}

	return nil, nil
}

func (this *NamespaceCreateAPI) createDomain(tx *gorm.DB, namespace *model.Namespace) error {
	if true {
		panic("wip ::createDomain()")
	}

	return nil
}

func (this *NamespaceCreateAPI) createMembership(tx *gorm.DB, namespace *model.Namespace) error {
	if true {
		panic("wip ::createMembership()")
	}
	return nil
}
