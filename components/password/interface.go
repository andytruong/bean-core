package password

type PasswordAlgorithm interface {
	Name() string
	Encrypt(rawPassword string) (string, error)
}

func New() (PasswordAlgorithm, error) {
	return &BCryptAlgorithm{}, nil
}
