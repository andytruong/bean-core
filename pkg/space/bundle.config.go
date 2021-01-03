package space

type Config struct {
	Manager ManagerConfig `yaml:"manager"`
}

type ManagerConfig struct {
	MaxNumberOfManager int `yaml:"maxNumberOfManager"`
}
