package dto

import "bean/pkg/util/api/scalar"

type S3UploadTokenInput struct {
	ApplicationId string             `json:"applicationId"`
	FilePath      scalar.Uri         `json:"filePath"`
	ContentType   scalar.ContentType `json:"contentType"`

	// TODO: add custom tags
}
