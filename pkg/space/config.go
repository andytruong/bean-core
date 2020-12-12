package space

type SpaceConfiguration struct {
	Manager ManagerConfig `yaml:"manager"`
}

type ManagerConfig struct {
	MaxNumberOfManager int `yaml:"maxNumberOfManager"`
}
