package config

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"gorm.io/gorm"

	"bean/pkg/config/model"
	"bean/pkg/config/model/dto"
	"bean/pkg/util"
	"bean/pkg/util/api/scalar"
	"bean/pkg/util/connect"
)

type CoreBucket struct {
	bean *ConfigBean
}

func (this CoreBucket) Create(tx *gorm.DB, in dto.BucketCreateInput) (*dto.BucketMutationOutcome, error) {
	bucket := &model.ConfigBucket{
		Id:          this.bean.id.MustULID(),
		Version:     this.bean.id.MustULID(),
		Slug:        scalar.NotNilString(in.Slug, this.bean.id.MustULID()),
		Title:       scalar.NotNilString(in.Title, ""),
		Description: in.Description,
		Access:      "777",
		Schema:      in.Schema,
		HostId:      in.HostId,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		IsPublished: in.IsPublished,
	}

	if nil != in.Access {
		bucket.Access = *in.Access
	}

	err := tx.Table(connect.TableConfigBucket).Create(&bucket).Error
	if nil != err {
		return nil, err
	}

	return &dto.BucketMutationOutcome{Errors: nil, Bucket: bucket}, nil
}

func (this CoreBucket) Update(ctx context.Context, tx *gorm.DB, in dto.BucketUpdateInput) (*dto.BucketMutationOutcome, error) {
	bucket, err := this.Load(ctx, tx, in.Id)
	if nil != err {
		return nil, err
	}

	if bucket.Version != in.Version {
		return nil, util.ErrorVersionConflict
	}

	changed := false
	if in.Title != nil {
		if bucket.Title != *in.Title {
			changed = true
			bucket.Title = *in.Title
		}
	}

	if in.Description != nil {
		if bucket.Description != in.Description {
			changed = true
			bucket.Description = in.Description
		}
	}

	if in.Access != nil {
		if bucket.Access != *in.Access {
			changed = true
			bucket.Access = *in.Access
		}
	}

	if in.Schema != nil {
		if bucket.Schema != *in.Schema {
			if bucket.IsPublished {
				return nil, util.ErrorLocked
			}

			changed = true
			bucket.Schema = *in.Schema
		}
	}

	if nil != in.IsPublished {
		if *in.IsPublished != bucket.IsPublished {
			if bucket.IsPublished {
				return nil, errors.Wrap(util.ErrorLocked, "change not un-publish a published bucket")
			}

			bucket.IsPublished = *in.IsPublished
		}
	}

	if changed {
		bucket.Version = this.bean.id.MustULID()
		err = tx.
			Table(connect.TableConfigBucket).
			Save(&bucket).
			Error
		if nil != err {
			return nil, err
		}
	}

	return &dto.BucketMutationOutcome{
		Errors: nil,
		Bucket: bucket,
	}, nil
}

// TODO: need data-loader
func (this CoreBucket) Load(ctx context.Context, db *gorm.DB, id string) (*model.ConfigBucket, error) {
	bucket := &model.ConfigBucket{}

	err := db.
		WithContext(ctx).
		Table(connect.TableConfigBucket).
		Where("id = ?", id).
		First(&bucket).
		Error
	if nil != err {
		return nil, err
	}

	return bucket, nil
}
