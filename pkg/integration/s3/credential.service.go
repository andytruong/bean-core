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
)

type credentialService struct {
	bundle    *S3Bundle
	transport http.RoundTripper
}

func (srv *credentialService) load(ctx context.Context, appId string) (*model.S3Credentials, error) {
	var (
		err      error
		bucket   *configModel.ConfigBucket
		variable *configModel.ConfigVariable
		cre      = &model.S3Credentials{}
	)

	// load current bucket
	bucket, err = srv.bundle.configBundle.BucketService.Load(ctx, configDto.BucketKey{Slug: credentialsConfigSlug})
	if nil != err {
		return nil, err
	}

	// load current variable
	variable, err = srv.bundle.configBundle.VariableService.Load(ctx, configDto.VariableKey{BucketId: bucket.Id, Name: appId})
	if nil != err {
		return nil, err
	}

	err = json.Unmarshal([]byte(variable.Value), cre)
	if nil != err {
		return nil, err
	}

	cre.Id = variable.Id
	cre.Version = variable.Version

	return cre, nil
}

func (srv *credentialService) save(ctx context.Context, in dto.S3CredentialsInput) (*model.S3Credentials, error) {
	var (
		err    error
		bucket *configModel.ConfigBucket
	)

	// load current bucket
	bucket, err = srv.bundle.configBundle.BucketService.Load(ctx, configDto.BucketKey{Slug: credentialsConfigSlug})
	if nil != err {
		return nil, errors.Wrap(err, "bucket load error")
	}

	cre, err := srv.load(ctx, in.ApplicationId)
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

func (srv credentialService) encrypt(text string) string {
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

func (srv credentialService) decrypt(cryptoText string) string {
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

func (srv *credentialService) client(creds *model.S3Credentials) (*minio.Client, error) {
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
