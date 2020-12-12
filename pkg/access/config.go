package access

import (
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"

	"bean/components/claim"
	"bean/components/scalar"
	"bean/pkg/util"
)

type (
	AccessConfiguration struct {
		SessionTimeout time.Duration `yaml:"timeout"`
		Jwt            JwtConfig     `yaml:"jwt"`

		mutex      *sync.Mutex
		privateKey interface{}
		publicKey  interface{}
	}

	JwtConfig struct {
		Algorithm  string          `yaml:"algorithm"`
		PrivateKey scalar.FilePath `yaml:"privateKey"`
		PublicKey  scalar.FilePath `yaml:"publicKey"`
		Timeout    time.Duration   `yaml:"timeout"`
	}
)

func (this *AccessConfiguration) init() *AccessConfiguration {
	if nil == this.mutex {
		this.mutex = &sync.Mutex{}
	}

	if time.Duration(0) == this.SessionTimeout {
		this.SessionTimeout, _ = time.ParseDuration("128h")
	}

	// go time to validate configuration
	// …

	return this
}

func (this *AccessConfiguration) signMethod() jwt.SigningMethod {
	switch this.Jwt.Algorithm {
	case "RS512":
		return jwt.SigningMethodRS512

	default:
		panic(util.ErrorToBeImplemented)
	}
}

func (this *AccessConfiguration) GetSignKey() (interface{}, error) {
	if nil == this.privateKey {
		this.mutex.Lock()
		defer this.mutex.Unlock()

		file, err := ioutil.ReadFile(this.Jwt.PrivateKey.String())
		if nil != err {
			return nil, err
		}

		block, _ := pem.Decode(file)
		this.privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)

		if err != nil {
			return nil, err
		}
	}

	return this.privateKey, nil
}

func (this *AccessConfiguration) GetParseKey() (interface{}, error) {
	if nil == this.publicKey {
		this.mutex.Lock()
		defer this.mutex.Unlock()

		pub, err := claim.ParseRsaPublicKeyFromFile(this.Jwt.PublicKey.String())
		if err != nil {
			return nil, err
		} else {
			this.publicKey = pub
		}
	}

	return this.publicKey, nil
}
