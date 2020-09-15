package mailer

type Genetic struct {
	ConnectionUrl string `yaml:"connectionUrl"`
	Sender        struct {
		Name    string `yaml:"name"`
		Address string `yaml:"address"`
	} `yaml:"sender"`
	
	Reroute struct {
		Enabled   bool              `yaml:"enabled"`
		Recipient string            `yaml:"recipient"`
		Filters   map[string]string `yaml:"filters"`
	} `yaml:"reroute"`
	
	Attachment struct {
		SizeLimit        uint64   `yaml:"sizeLimit"`
		AllowContentType []string `yaml:"allowContentType"`
	}
}
