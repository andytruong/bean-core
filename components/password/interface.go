package password

import (
	"fmt"
)

type PasswordAlgorithm interface {
	Name() string
	Encrypt(rawPassword string) (string, error)
}

func Get(algo string) (PasswordAlgorithm, error) {
	switch algo {
	case "bcrypt":
		return &BCryptAlgorithm{}, nil

	default:
		return nil, fmt.Errorf("algorithm not found: %s", algo)
	}
}
