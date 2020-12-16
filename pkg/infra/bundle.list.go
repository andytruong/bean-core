package infra

import (
	"bean/components/module"
)

// TODO: Generate this code
func (this *bundles) List() []module.Bundle {
	userBundle, _ := this.User()
	spaceBundle, _ := this.Space()
	accessBundle, _ := this.Access()
	s3Bundle, _ := this.S3()
	mailerIntegrationBundle, _ := this.Mailer()

	return []module.Bundle{userBundle, spaceBundle, accessBundle, s3Bundle, mailerIntegrationBundle}
}
