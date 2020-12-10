package mailer

type Genetic struct {
	Reroute struct {
		Enabled   bool              `yaml:"enabled"`
		Recipient string            `yaml:"recipient"`
		Filters   map[string]string `yaml:"filters"`
	} `yaml:"reroute"`

	Attachment struct {
		SizeLimit        uint64   `yaml:"sizeLimit"`
		SizeLimitEach    uint64   `yaml:"sizeLimitEach"`
		AllowContentType []string `yaml:"allowContentType"`
	}
}
