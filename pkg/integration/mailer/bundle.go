package mailer

import (
	"path"
	"runtime"

	"go.uber.org/zap"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"

	"bean/components/module"
	"bean/components/module/migrate"
	"bean/pkg/integration/mailer/model"
)

func NewMailerBundle(genetic *MailerConfiguration, logger *zap.Logger) *MailerBundle {
	this := &MailerBundle{
		config: genetic,
		logger: logger,
	}

	this.resolvers = newResoler(this)

	return this
}

type MailerBundle struct {
	module.AbstractBundle

	config    *MailerConfiguration
	logger    *zap.Logger
	resolvers map[string]interface{}
}

func (bundle MailerBundle) Migrate(tx *gorm.DB, driver string) error {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return nil
	}

	runner := migrate.Runner{
		Tx:     tx,
		Logger: bundle.logger,
		Driver: driver,
		Bean:   "integration.mailer",
		Dir:    path.Dir(filename) + "/model/migration/",
	}

	return runner.Run()
}

func (bundle MailerBundle) Dependencies() []module.Bundle {
	panic("implement me")
}

func (bundle MailerBundle) Send(message model.Message) error {
	if true {
		return nil
	}

	if bundle.config.Reroute.Enabled {
		// TODO: check matching
		message.Recipient = bundle.config.Reroute.Recipient
	}

	dialer := &gomail.Dialer{} // gomail.NewDialer(host, port, username, password)

	return message.Send(dialer)
}

func (bundle *MailerBundle) GraphqlResolver() map[string]interface{} {
	return bundle.resolvers
}
