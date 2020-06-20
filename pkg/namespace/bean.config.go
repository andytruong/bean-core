package namespace

import (
	"context"
	"time"

	"gorm.io/gorm"

	"bean/pkg/namespace/model"
	"bean/pkg/namespace/model/dto"
	"bean/pkg/util/connect"
)

type CoreConfig struct {
	bean *NamespaceBean
}

func (this *CoreConfig) CreateFeatures(tx *gorm.DB, namespace *model.Namespace, in dto.NamespaceCreateInput) error {
	if in.Object.Features.Register {
		return this.CreateFeature(tx, namespace, "default", "register", []byte("true"))
	} else {
		return this.CreateFeature(tx, namespace, "default", "register", []byte("false"))
	}
}

func (this *CoreConfig) CreateFeature(
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

func (this *CoreConfig) List(ctx context.Context, namespace *model.Namespace) (*model.NamespaceFeatures, error) {
	features := &model.NamespaceFeatures{
		Register: false,
	}

	var configList []model.NamespaceConfig
	err := this.bean.db.
		Table(connect.TableNamespaceConfig).
		Find(&configList, "namespace_id = ?", namespace.ID).Error

	if nil != err {
		return nil, err
	}

	for _, config := range configList {
		switch config.Bucket {
		case "default":
			switch config.Key {
			case "register":
				if "true" == string(config.Value) {
					features.Register = true
				}

			default:
				panic("unknown bucket: " + config.Bucket)
			}

		default:
			panic("unknown bucket: " + config.Bucket)
		}
	}

	return features, nil
}

func (this *CoreConfig) updateFeatures(tx *gorm.DB, obj *model.Namespace, in dto.NamespaceUpdateInput) error {
	if nil != in.Object.Features.Register {
		if *in.Object.Features.Register {
			return this.updateFeature(tx, obj, "default", "register", []byte("true"))
		} else {
			return this.updateFeature(tx, obj, "default", "register", []byte("false"))
		}
	}

	return nil
}

func (this *CoreConfig) updateFeature(
	tx *gorm.DB,
	obj *model.Namespace, bucket string, key string, value []byte,
) error {
	return tx.
		Table(connect.TableNamespaceConfig).
		Where("namespace_id = ? AND bucket = ? AND key = ?", obj.ID, bucket, key).
		Updates(&model.NamespaceConfig{
			Version:   this.bean.id.MustULID(),
			Value:     value,
			UpdatedAt: time.Now(),
		}).
		Error
}
