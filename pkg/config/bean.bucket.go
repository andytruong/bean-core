package config

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"gorm.io/gorm"

	"bean/pkg/config/model"
	"bean/pkg/config/model/dto"
	"bean/pkg/util"
	"bean/pkg/util/connect"
)

type CoreBucket struct {
	bean *ConfigBean
}

func (this CoreBucket) Create(ctx context.Context, tx *gorm.DB, input dto.BucketCreateInput) (*dto.BucketMutationOutcome, error) {
	bucket := &model.ConfigBucket{
		Id:          this.bean.id.MustULID(),
		Version:     this.bean.id.MustULID(),
		Slug:        util.NotNilString(input.Slug, this.bean.id.MustULID()),
		Title:       util.NotNilString(input.Title, ""),
		Description: input.Description,
		Access:      "777",
		Schema:      input.Schema,
		HostId:      input.HostId,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		IsPublished: input.IsPublished,
	}

	if nil != input.Access {
		bucket.Access = *input.Access
	}

	err := tx.Table(connect.TableConfigBucket).Save(&bucket).Error
	if nil != err {
		return nil, err
	}

	return &dto.BucketMutationOutcome{Errors: nil, Bucket: bucket}, nil
}

func (this CoreBucket) Update(ctx context.Context, tx *gorm.DB, input dto.BucketUpdateInput) (*dto.BucketMutationOutcome, error) {
	bucket, err := this.Load(ctx, tx, input.Id)
	if nil != err {
		return nil, err
	}

	if bucket.Version != input.Version {
		return nil, util.ErrorVersionConflict
	}

	changed := false
	if input.Title != nil {
		if bucket.Title != *input.Title {
			changed = true
			bucket.Title = *input.Title
		}
	}

	if input.Description != nil {
		if bucket.Description != input.Description {
			changed = true
			bucket.Description = input.Description
		}
	}

	if input.Access != nil {
		if bucket.Access != *input.Access {
			changed = true
			bucket.Access = *input.Access
		}
	}

	if input.Schema != nil {
		if bucket.Schema != *input.Schema {
			if bucket.IsPublished {
				return nil, util.ErrorLocked
			}

			changed = true
			bucket.Schema = *input.Schema
		}
	}

	if nil != input.IsPublished {
		if *input.IsPublished != bucket.IsPublished {
			if bucket.IsPublished {
				return nil, errors.Wrap(util.ErrorLocked, "change not un-publish a published bucket")
			}

			bucket.IsPublished = *input.IsPublished
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

	err := db.Table(connect.TableConfigBucket).First(&bucket, "id = ?", id).Error
	if nil != err {
		return nil, err
	}

	return bucket, nil
}
