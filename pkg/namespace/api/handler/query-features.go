package handler

import (
	"context"

	"github.com/jinzhu/gorm"

	"bean/pkg/namespace/model"
)

type NamespaceQueryFeaturesHandler struct {
	DB *gorm.DB
}

func (this NamespaceQueryFeaturesHandler) Features(ctx context.Context, namespace *model.Namespace) (*model.NamespaceFeatures, error) {
	features := &model.NamespaceFeatures{
		Register: false,
	}

	var configList []model.NamespaceConfig
	err := this.DB.
		Table("namespace_config").
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
