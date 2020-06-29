package s3

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"

	"gorm.io/gorm"

	"bean/pkg/integration/s3/model"
	"bean/pkg/integration/s3/model/dto"
	"bean/pkg/util"
	"bean/pkg/util/connect"
)

type coreCredentials struct {
	bean *S3IntegrationBean
}

func (this *coreCredentials) onAppCreate(tx *gorm.DB, app *model.Application, in dto.S3ApplicationCredentialsCreateInput) error {
	cre := model.Credentials{
		ID:            this.bean.id.MustULID(),
		ApplicationId: app.ID,
		Endpoint:      in.Endpoint,
		AccessKey:     in.AccessKey,
		SecretKey:     this.encrypt(in.SecretKey),
		IsSecure:      in.IsSecure,
	}

	return tx.Table(connect.TableIntegrationS3Credentials).Create(&cre).Error
}

func (this *coreCredentials) onAppUpdate(tx *gorm.DB, app *model.Application, in *dto.S3ApplicationCredentialsUpdateInput) error {
	if nil == in {
		return nil
	}

	// load
	cre := &model.Credentials{}
	err := tx.Table(connect.TableIntegrationS3Credentials).Where("application_id = ?", app.ID).First(&cre).Error

	if nil == err {
		// if found -> update
		changed := false

		if nil != in.Endpoint {
			changed = true
			cre.Endpoint = *in.Endpoint
		}

		if nil != in.IsSecure {
			changed = true
			cre.IsSecure = *in.IsSecure
		}

		if nil != in.AccessKey || nil != in.SecretKey {
			if nil != in.AccessKey {
				changed = true
				cre.AccessKey = *in.AccessKey
			}

			if nil != in.SecretKey {
				changed = true
				cre.SecretKey = this.encrypt(*in.SecretKey)
			}
		}

		if changed {
			return tx.Table(connect.TableIntegrationS3Credentials).Save(&cre).Error
		}

		return util.ErrorUselessInput
	} else {
		if gorm.ErrRecordNotFound != err {
			return err
		}

		if nil != in.Endpoint && nil != in.AccessKey && nil != in.SecretKey {
			// if not found -> create
			cre = &model.Credentials{
				ID:            this.bean.id.MustULID(),
				ApplicationId: app.ID,
				Endpoint:      *in.Endpoint,
				AccessKey:     *in.AccessKey,
				SecretKey:     this.encrypt(*in.SecretKey),
				IsSecure:      false,
			}

			err := tx.Table(connect.TableIntegrationS3Credentials).Create(&cre).Error
			if nil != err {
				return err
			}
		}
	}

	return nil
}

func (this coreCredentials) encrypt(text string) string {
	plaintext := []byte(text)

	block, err := aes.NewCipher(this.bean.genetic.Key)
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

func (this coreCredentials) decrypt(cryptoText string) string {
	cipherText, _ := base64.URLEncoding.DecodeString(cryptoText)
	block, err := aes.NewCipher(this.bean.genetic.Key)
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

	return fmt.Sprintf("%s", cipherText)
}
