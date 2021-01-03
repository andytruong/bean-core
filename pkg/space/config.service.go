package space

import (
	"context"
	"time"

	"gorm.io/gorm"

	"bean/components/connect"
	"bean/pkg/space/model"
	"bean/pkg/space/model/dto"
)

type ConfigService struct {
	bundle *Bundle
}

func (srv *ConfigService) CreateFeatures(ctx context.Context, space *model.Space, in dto.SpaceCreateInput) error {
	value := []byte("false")
	if in.Object.Features.Register {
		value = []byte("true")
	}

	return srv.CreateFeature(ctx, space, "default", "register", value)
}

func (srv *ConfigService) CreateFeature(
	ctx context.Context,
	space *model.Space, bucket string, key string, value []byte,
) error {
	config := model.SpaceConfig{
		Id:        srv.bundle.idr.MustULID(),
		Version:   srv.bundle.idr.MustULID(),
		SpaceId:   space.ID,
		Bucket:    bucket,
		Key:       key,
		Value:     value,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return connect.ContextToDB(ctx).Create(&config).Error
}

func (srv *ConfigService) List(ctx context.Context, space *model.Space) (*model.SpaceFeatures, error) {
	db := connect.ContextToDB(ctx)
	features := &model.SpaceFeatures{Register: false}
	var configList []model.SpaceConfig

	err := db.Find(&configList, "space_id = ?", space.ID).Error
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

func (srv *ConfigService) updateFeatures(tx *gorm.DB, obj *model.Space, in dto.SpaceUpdateInput) error {
	if nil != in.Object.Features.Register {
		if *in.Object.Features.Register {
			return srv.updateFeature(tx, obj, "default", "register", []byte("true"))
		} else {
			return srv.updateFeature(tx, obj, "default", "register", []byte("false"))
		}
	}

	return nil
}

func (srv *ConfigService) updateFeature(
	tx *gorm.DB,
	obj *model.Space, bucket string, key string, value []byte,
) error {
	return tx.
		Where("space_id = ? AND bucket = ? AND key = ?", obj.ID, bucket, key).
		Updates(&model.SpaceConfig{
			Version:   srv.bundle.idr.MustULID(),
			Value:     value,
			UpdatedAt: time.Now(),
		}).
		Error
}
