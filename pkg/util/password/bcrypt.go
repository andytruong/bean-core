package password

import (
	"golang.org/x/crypto/bcrypt"
)

type BCryptAlgorithm struct {
}

func (this BCryptAlgorithm) Encrypt(rawPassword string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(rawPassword), 3)

	if nil != err {
		return "", err
	}

	return string(hash), nil
}

func (this BCryptAlgorithm) Name() string {
	return "bcrypt"
}
