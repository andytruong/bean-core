package mailer

type MailerConfiguration struct {
	Account struct {
		MaxPerSpace uint `yaml:"maxNumber"`
	} `yaml:"account"`
	
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
