package config

import (
	"context"
	"encoding/json"
	"time"

	"github.com/pkg/errors"
	"github.com/qri-io/jsonschema"

	"bean/components/connect"
	"bean/components/scalar"
	"bean/components/util"
	"bean/pkg/config/model"
	"bean/pkg/config/model/dto"
)

type BucketService struct {
	bundle *Bundle
}

func (srv BucketService) Create(ctx context.Context, in dto.BucketCreateInput) (*dto.BucketMutationOutcome, error) {
	// make sure input schema is valid
	{
		rs := &jsonschema.Schema{}
		err := json.Unmarshal([]byte(in.Schema), rs)
		if nil != err {
			err := util.NewError(util.ErrorCodeInput, []string{"BucketCreateInput.Schema"}, err.Error())

			return &dto.BucketMutationOutcome{Errors: []util.Error{err}}, nil
		}
	}

	db := connect.DB(ctx)
	bucket := &model.ConfigBucket{
		Id:          srv.bundle.idr.ULID(),
		Version:     srv.bundle.idr.ULID(),
		Slug:        scalar.NotNilString(in.Slug, srv.bundle.idr.ULID()),
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

	err := db.Create(&bucket).Error
	if nil != err {
		return nil, err
	}

	return &dto.BucketMutationOutcome{Errors: nil, Bucket: bucket}, nil
}

func (srv BucketService) Update(ctx context.Context, in dto.BucketUpdateInput) (*dto.BucketMutationOutcome, error) {
	tx := connect.DB(ctx)
	bucket, err := srv.Load(ctx, dto.BucketKey{Id: in.Id})
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
		bucket.Version = srv.bundle.idr.ULID()
		err = tx.Save(&bucket).Error
		if nil != err {
			return nil, err
		}
	}

	return &dto.BucketMutationOutcome{
		Errors: nil,
		Bucket: bucket,
	}, nil
}

func (srv BucketService) Load(ctx context.Context, key dto.BucketKey) (*model.ConfigBucket, error) {
	bucket := &model.ConfigBucket{}
	db := connect.DB(ctx)

	if key.Id != "" {
		// load by ID
		err := db.Where("id = ?", key.Id).Take(&bucket).Error
		if nil != err {
			return nil, err
		}
	} else if key.Slug != "" {
		// load by slug
		err := db.Where(dto.BucketKey{Slug: key.Slug}).Take(&bucket).Error
		if nil != err {
			return nil, err
		}
	}

	return bucket, nil
}
