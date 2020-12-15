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

func NewMailerIntegration(genetic *MailerConfiguration, logger *zap.Logger) *MailerIntegrationBundle {
	this := &MailerIntegrationBundle{
		config: genetic,
		logger: logger,
	}

	return this
}

type MailerIntegrationBundle struct {
	module.AbstractBundle

	config *MailerConfiguration
	logger *zap.Logger
}

func (this MailerIntegrationBundle) Migrate(tx *gorm.DB, driver string) error {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return nil
	}

	runner := migrate.Runner{
		Tx:     tx,
		Logger: this.logger,
		Driver: driver,
		Bean:   "integration.mailer",
		Dir:    path.Dir(filename) + "/model/migration/",
	}

	return runner.Run()
}

func (this MailerIntegrationBundle) Dependencies() []module.Bundle {
	panic("implement me")
}

func (this MailerIntegrationBundle) Send(message model.Message) error {
	if true {
		return nil
	}

	if this.config.Reroute.Enabled {
		// TODO: check matching
		message.Recipient = this.config.Reroute.Recipient
	}

	dialer := &gomail.Dialer{} // gomail.NewDialer(host, port, username, password)

	return message.Send(dialer)
}

func (this *MailerIntegrationBundle) GraphqlResolver() map[string]interface{} {
	return newResoler(this)
}
