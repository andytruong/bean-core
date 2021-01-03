package s3

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"bean/components/scalar"
	"bean/components/util"
	configModel "bean/pkg/config/model"
	configDto "bean/pkg/config/model/dto"
	"bean/pkg/integration/s3/model"
	"bean/pkg/integration/s3/model/dto"
)

const (
	credentialsConfigSlug   = `bean.s3.credentials.schema.v1`
	credentialsConfigSchema = `{
		"type":       "object",
		"required":   ["endpoint", "bucket", "accessKey", "secretKey", "isSecure"],
		"properties": {
			"endpoint":  { "type": "string", "maxLength": 255, "format": "uri" },
			"bucket":    { "type": "string", "maxLength": 64  },
			"accessKey": { "type": "string", "maxLength": 64  },
			"secretKey": { "type": "string", "maxLength": 128 },
			"isSecure":  { "type": "boolean" }
		}
	}`

	uploadPolicyConfigSlug   = `bean.s3.policy.schema.v1`
	uploadPolicyConfigSchema = `{
		"type":       "object",
		"required":   [],
		"properties": {
			"fileExtensions": {
				"type":  "array",
				"items": { "type": "string", "maxLength": 32 }
			},
			"rateLimit":      {
				"type":  "array",
				"items": {
					"type":       "object",
					"required":   ["value", "object", "interval"],
					"properties": {
						"value":    { "type": "string" },
						"object":   { "type": "string", "enum": ["user", "space"] },
						"interval": {
							"type": "string",
							"pattern": "^(\\d+) (minute|minutes|hour|hours|day|days)$"
						}
					}
				}
			}
		}
	}`
)

type configService struct {
	bundle    *Bundle
	transport http.RoundTripper
}

func (srv configService) loadVariable(ctx context.Context, bucketSlug string, appId string) (*configModel.ConfigVariable, error) {
	bucket, err := srv.bundle.configBundle.BucketService.Load(ctx, configDto.BucketKey{Slug: bucketSlug})
	if nil != err {
		return nil, err
	}

	return srv.bundle.configBundle.VariableService.Load(ctx, configDto.VariableKey{BucketId: bucket.Id, Name: appId})
}

func (srv configService) loadCredentials(ctx context.Context, appId string) (*model.S3Credentials, error) {
	// load current variable
	variable, err := srv.loadVariable(ctx, credentialsConfigSlug, appId)
	if nil != err {
		return nil, err
	}

	cre := &model.S3Credentials{}
	err = json.Unmarshal([]byte(variable.Value), cre)
	if nil != err {
		return nil, err
	}

	cre.Id = variable.Id
	cre.Version = variable.Version

	return cre, nil
}

func (srv *configService) saveCredentials(ctx context.Context, in dto.S3CredentialsInput) (*model.S3Credentials, error) {
	var (
		err    error
		bucket *configModel.ConfigBucket
	)

	// load current bucket
	bucket, err = srv.bundle.configBundle.BucketService.Load(ctx, configDto.BucketKey{Slug: credentialsConfigSlug})
	if nil != err {
		return nil, errors.Wrap(err, "bucket load error")
	}

	cre, err := srv.loadCredentials(ctx, in.ApplicationId)
	if nil != err {
		if err != gorm.ErrRecordNotFound {
			return nil, err
		}
	}

	newCre := model.S3Credentials{
		Endpoint:  in.Endpoint,
		Bucket:    in.Bucket,
		AccessKey: in.AccessKey,
		SecretKey: srv.encrypt(in.SecretKey),
		IsSecure:  in.IsSecure,
	}

	newCreBytes, err := json.Marshal(newCre)
	if nil != err {
		return nil, err
	}

	var out *configDto.VariableMutationOutcome

	if nil == cre {
		out, err = srv.bundle.configBundle.VariableService.Create(ctx, configDto.VariableCreateInput{
			BucketId: bucket.Id,
			Name:     in.ApplicationId,
			Value:    string(newCreBytes),
			IsLocked: scalar.NilBool(false),
		})

		if nil != err {
			return nil, err
		}
	} else {
		useless := cre.Endpoint == newCre.Endpoint &&
			cre.Bucket == newCre.Bucket &&
			cre.AccessKey == newCre.AccessKey &&
			srv.decrypt(cre.SecretKey) == in.SecretKey &&
			cre.IsSecure == newCre.IsSecure

		if useless {
			return nil, util.ErrorUselessInput
		}

		out, err = srv.bundle.configBundle.VariableService.Update(ctx, configDto.VariableUpdateInput{
			Id:       cre.Id,
			Version:  in.Version,
			Value:    scalar.NilString(string(newCreBytes)),
			IsLocked: scalar.NilBool(false),
		})

		if nil != err {
			return nil, err
		}
	}

	return &model.S3Credentials{
		Id:        out.Variable.Id,
		Version:   out.Variable.Version,
		Endpoint:  in.Endpoint,
		Bucket:    in.Bucket,
		IsSecure:  in.IsSecure,
		AccessKey: in.AccessKey,
		SecretKey: in.SecretKey,
	}, nil
}

func (srv configService) loadUploadPolicy(ctx context.Context, appId string) (*model.S3UploadPolicy, error) {
	variable, err := srv.loadVariable(ctx, uploadPolicyConfigSlug, appId)
	if nil != err {
		return nil, err
	}

	pol := &model.S3UploadPolicy{}
	err = json.Unmarshal([]byte(variable.Value), pol)
	if nil != err {
		return nil, err
	}

	pol.Id = variable.Id
	pol.Version = variable.Version

	return pol, nil
}

func (srv *configService) saveUploadPolicy(ctx context.Context, in dto.UploadPolicyInput) (*model.S3UploadPolicy, error) {
	var (
		err    error
		bucket *configModel.ConfigBucket
		policy *model.S3UploadPolicy
	)

	// load current bucket
	bucket, err = srv.bundle.configBundle.BucketService.Load(ctx, configDto.BucketKey{Slug: uploadPolicyConfigSlug})
	if nil != err {
		return nil, errors.Wrap(err, "bucket load error")
	}

	policy, err = srv.loadUploadPolicy(ctx, in.ApplicationId)
	if nil != err {
		if err != gorm.ErrRecordNotFound {
			return nil, err
		}
	}

	newPolicy := model.S3UploadPolicy{
		FileExtensions: in.FileExtensions,
		RateLimit:      []model.UploadRateLimitPolicy{},
	}

	for _, input := range in.RateLimit {
		newPolicy.RateLimit = append(newPolicy.RateLimit, model.UploadRateLimitPolicy{
			Value:    input.Value,
			Object:   input.Object,
			Interval: input.Interval,
		})
	}

	newPolicyBytes, err := json.Marshal(newPolicy)
	if nil != err {
		return nil, err
	}

	var out *configDto.VariableMutationOutcome

	if nil == policy {
		out, err = srv.bundle.configBundle.VariableService.Create(ctx, configDto.VariableCreateInput{
			BucketId: bucket.Id,
			Name:     in.ApplicationId,
			Value:    string(newPolicyBytes),
			IsLocked: scalar.NilBool(false),
		})

		if nil != err {
			return nil, err
		}
	} else {
		useless, err := policy.EqualTo(newPolicy)
		if nil != err {
			return nil, err
		} else if useless {
			return nil, util.ErrorUselessInput
		}

		out, err = srv.bundle.configBundle.VariableService.Update(ctx, configDto.VariableUpdateInput{
			Id:       policy.Id,
			Version:  in.Version,
			Value:    scalar.NilString(string(newPolicyBytes)),
			IsLocked: scalar.NilBool(false),
		})

		if nil != err {
			return nil, err
		}
	}

	return &model.S3UploadPolicy{
		Id:             out.Variable.Id,
		Version:        out.Variable.Version,
		FileExtensions: newPolicy.FileExtensions,
		RateLimit:      newPolicy.RateLimit,
	}, nil
}

func (srv configService) encrypt(text string) string {
	plaintext := []byte(text)

	block, err := aes.NewCipher([]byte(srv.bundle.cnf.Key))
	if err != nil {
		panic(err)
	}

	cipherText := make([]byte, aes.BlockSize+len(plaintext))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plaintext)

	return base64.URLEncoding.EncodeToString(cipherText)
}

func (srv configService) decrypt(cryptoText string) string {
	cipherText, _ := base64.URLEncoding.DecodeString(cryptoText)
	block, err := aes.NewCipher([]byte(srv.bundle.cnf.Key))
	if err != nil {
		panic(err)
	}

	if len(cipherText) < aes.BlockSize {
		panic("cipherText too short")
	}
	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)

	return string(cipherText)
}

func (srv *configService) client(creds *model.S3Credentials) (*minio.Client, error) {
	endpoint := string(creds.Endpoint)
	endpoint = strings.Replace(endpoint, "http://", "", 1)
	endpoint = strings.Replace(endpoint, "https://", "", 1)

	return minio.New(
		endpoint,
		&minio.Options{
			Creds:     credentials.NewStaticV4(creds.AccessKey, srv.decrypt(creds.SecretKey), ""),
			Secure:    creds.IsSecure,
			Transport: srv.transport,
		},
	)
}
