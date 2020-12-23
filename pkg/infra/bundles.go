package infra

import (
	"bean/pkg/access"
	"bean/pkg/app"
	"bean/pkg/config"
	"bean/pkg/integration/mailer"
	"bean/pkg/integration/s3"
	"bean/pkg/space"
	"bean/pkg/user"
)

type (
	bundles struct {
		container *Container

		user   *user.UserBundle
		space  *space.SpaceBundle
		config *config.ConfigBundle
		access *access.AccessBundle
		s3     *s3.S3Bundle
		mailer *mailer.MailerBundle
		app    *app.AppBundle
	}
)
