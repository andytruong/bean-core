package password

import (
	"golang.org/x/crypto/bcrypt"
)

type BCryptAlgorithm struct {
}

func (algo BCryptAlgorithm) Encrypt(rawPassword string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(rawPassword), 3)

	if nil != err {
		return "", err
	}

	return string(hash), nil
}

func (algo BCryptAlgorithm) Name() string {
	return "bcrypt"
}
