package connect

type DatabaseConfig struct {
	Driver string `yaml:"driver"`
	Url    string `yaml:"url"`
}
