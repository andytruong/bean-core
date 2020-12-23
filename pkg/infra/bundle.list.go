package infra

import (
	"bean/components/module"
)

// TODO: Generate this code
func (list *bundles) List() []module.Bundle {
	userBundle, _ := list.User()
	spaceBundle, _ := list.Space()
	accessBundle, _ := list.Access()
	s3Bundle, _ := list.S3()
	mailerIntegrationBundle, _ := list.Mailer()

	return []module.Bundle{userBundle, spaceBundle, accessBundle, s3Bundle, mailerIntegrationBundle}
}
