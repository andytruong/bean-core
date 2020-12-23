package space

import (
	"context"
	"time"

	"gorm.io/gorm"

	"bean/pkg/space/model"
	"bean/pkg/space/model/dto"
)

type ConfigService struct {
	bundle *SpaceBundle
}

func (service *ConfigService) CreateFeatures(tx *gorm.DB, space *model.Space, in dto.SpaceCreateInput) error {
	if in.Object.Features.Register {
		return service.CreateFeature(tx, space, "default", "register", []byte("true"))
	} else {
		return service.CreateFeature(tx, space, "default", "register", []byte("false"))
	}
}

func (service *ConfigService) CreateFeature(
	tx *gorm.DB,
	space *model.Space, bucket string, key string, value []byte,
) error {
	config := model.SpaceConfig{
		Id:        service.bundle.id.MustULID(),
		Version:   service.bundle.id.MustULID(),
		SpaceId:   space.ID,
		Bucket:    bucket,
		Key:       key,
		Value:     value,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return tx.Create(&config).Error
}

func (service *ConfigService) List(ctx context.Context, space *model.Space) (*model.SpaceFeatures, error) {
	features := &model.SpaceFeatures{
		Register: false,
	}

	var configList []model.SpaceConfig
	err := service.bundle.db.Find(&configList, "space_id = ?", space.ID).Error
	if nil != err {
		return nil, err
	}

	for _, config := range configList {
		switch config.Bucket {
		case "default":
			switch config.Key {
			case "register":
				if string(config.Value) == "true" {
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

func (service *ConfigService) updateFeatures(tx *gorm.DB, obj *model.Space, in dto.SpaceUpdateInput) error {
	if nil != in.Object.Features.Register {
		if *in.Object.Features.Register {
			return service.updateFeature(tx, obj, "default", "register", []byte("true"))
		} else {
			return service.updateFeature(tx, obj, "default", "register", []byte("false"))
		}
	}

	return nil
}

func (service *ConfigService) updateFeature(
	tx *gorm.DB,
	obj *model.Space, bucket string, key string, value []byte,
) error {
	return tx.
		Where("space_id = ? AND bucket = ? AND key = ?", obj.ID, bucket, key).
		Updates(&model.SpaceConfig{
			Version:   service.bundle.id.MustULID(),
			Value:     value,
			UpdatedAt: time.Now(),
		}).
		Error
}
