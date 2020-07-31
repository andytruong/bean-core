package space

type Genetic struct {
	Manager ManagerConfig `yaml:"manager"`
}

type ManagerConfig struct {
	MaxNumberOfManager int `yaml:"maxNumberOfManager"`
}
