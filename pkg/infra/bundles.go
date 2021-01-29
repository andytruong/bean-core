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
	bundleList struct {
		container *Container

		user   *user.Bundle
		space  *space.Bundle
		config *config.Bundle
		access *access.Bundle
		s3     *s3.Bundle
		mailer *mailer.Bundle
		app    *app.Bundle
	}
)
