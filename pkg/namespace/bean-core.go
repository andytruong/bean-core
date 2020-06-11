package namespace

import (
	"time"

	"github.com/jinzhu/gorm"

	"bean/pkg/namespace/model"
	"bean/pkg/namespace/model/dto"
	"bean/pkg/util"
	"bean/pkg/util/api"
	"bean/pkg/util/connect"
)

type NamespaceBeanCore struct {
	bean *NamespaceBean
}

func (this *NamespaceBeanCore) Create(tx *gorm.DB, input dto.NamespaceCreateInput) (*dto.NamespaceCreateOutcome, error) {
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

func (this *NamespaceBeanCore) doCreate(tx *gorm.DB, input dto.NamespaceCreateInput) (*model.Namespace, error) {
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

func (this *NamespaceBeanCore) createRelationships(tx *gorm.DB, namespace *model.Namespace, input dto.NamespaceCreateInput) error {
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
	if ownerRole, err := this.doCreate(tx, ownerRoleInput); nil != err {
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

func (this *NamespaceBeanCore) createDomains(tx *gorm.DB, namespace *model.Namespace, input dto.NamespaceCreateInput) error {
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

func (this *NamespaceBeanCore) createDomain(tx *gorm.DB, namespace *model.Namespace, input *dto.DomainNameInput, isPrimary bool) error {
	domain := model.DomainName{
		ID:          this.bean.id.MustULID(),
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

func (this *NamespaceBeanCore) createMembership(tx *gorm.DB, namespace *model.Namespace, input dto.NamespaceCreateInput) error {
	membership := model.Membership{
		ID:          this.bean.id.MustULID(),
		Version:     this.bean.id.MustULID(),
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

func (this *NamespaceBeanCore) createFeatures(tx *gorm.DB, namespace *model.Namespace, input dto.NamespaceCreateInput) error {
	if input.Object.Features.Register {
		return this.createFeature(tx, namespace, "default", "register", []byte("true"))
	} else {
		return this.createFeature(tx, namespace, "default", "register", []byte("false"))
	}
}

func (this *NamespaceBeanCore) createFeature(
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
