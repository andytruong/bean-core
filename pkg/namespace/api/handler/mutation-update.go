package handler

import (
	"time"

	"github.com/jinzhu/gorm"

	"bean/pkg/namespace/model"
	"bean/pkg/namespace/model/dto"
	"bean/pkg/util"
	"bean/pkg/util/connect"
)

type NamespaceUpdateHandler struct {
	ID *util.Identifier
}

func (this NamespaceUpdateHandler) NamespaceUpdate(
	tx *gorm.DB,
	namespace *model.Namespace,
	input dto.NamespaceUpdateInput,
) (*bool, error) {
	// check version for conflict
	if input.NamespaceVersion != namespace.Version {
		return nil, util.ErrorVersionConflict
	}

	// change version
	version, err := this.ID.ULID()
	if nil != err {
		return nil, err
	}

	if nil != input.Object.Language {
		namespace.Language = *input.Object.Language
	}

	namespace.Version = version
	if err := tx.Save(namespace).Error; nil != err {
		return nil, err
	}

	err = this.updateFeatures(tx, namespace, input)
	if nil != err {
		return nil, err
	}

	return util.NilBool(true), nil
}

func (this *NamespaceUpdateHandler) updateFeatures(
	tx *gorm.DB,
	namespace *model.Namespace, input dto.NamespaceUpdateInput,
) error {
	if nil != input.Object.Features.Register {
		if *input.Object.Features.Register {
			return this.updateFeature(tx, namespace, "default", "register", []byte("true"))
		} else {
			return this.updateFeature(tx, namespace, "default", "register", []byte("false"))
		}
	}

	return nil
}

func (this *NamespaceUpdateHandler) updateFeature(
	tx *gorm.DB,
	namespace *model.Namespace, bucket string, key string, value []byte,
) error {
	version, err := this.ID.ULID()
	if nil != err {
		return err
	}

	return tx.
		Table(connect.TableNamespaceConfig).
		Where("namespace_id = ? AND bucket = ? AND key = ?", namespace.ID, bucket, key).
		Update(&model.NamespaceConfig{
			Version:   version,
			Value:     value,
			UpdatedAt: time.Now(),
		}).
		Error
}
