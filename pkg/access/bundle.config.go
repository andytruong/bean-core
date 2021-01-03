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
	"bean/components/util"
)

type (
	Config struct {
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

func (cnf *Config) init() *Config {
	if nil == cnf.mutex {
		cnf.mutex = &sync.Mutex{}
	}

	if time.Duration(0) == cnf.SessionTimeout {
		cnf.SessionTimeout, _ = time.ParseDuration("128h")
	}

	// go time to validate configuration
	// …

	return cnf
}

func (cnf *Config) signMethod() jwt.SigningMethod {
	switch cnf.Jwt.Algorithm {
	case "RS512":
		return jwt.SigningMethodRS512

	default:
		panic(util.ErrorToBeImplemented)
	}
}

func (cnf *Config) GetSignKey() (interface{}, error) {
	if nil == cnf.privateKey {
		cnf.mutex.Lock()
		defer cnf.mutex.Unlock()

		file, err := ioutil.ReadFile(cnf.Jwt.PrivateKey.String())
		if nil != err {
			return nil, err
		}

		block, _ := pem.Decode(file)
		cnf.privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)

		if err != nil {
			return nil, err
		}
	}

	return cnf.privateKey, nil
}

func (cnf *Config) GetParseKey() (interface{}, error) {
	if nil == cnf.publicKey {
		cnf.mutex.Lock()
		defer cnf.mutex.Unlock()

		pub, err := claim.ParseRsaPublicKeyFromFile(cnf.Jwt.PublicKey.String())
		if err != nil {
			return nil, err
		} else {
			cnf.publicKey = pub
		}
	}

	return cnf.publicKey, nil
}
