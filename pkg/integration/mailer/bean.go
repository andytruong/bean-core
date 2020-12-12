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

func NewMailerIntegration(genetic *Genetic, logger *zap.Logger) *MailerIntegrationBean {
	this := &MailerIntegrationBean{
		genetic:  genetic,
		logger:   logger,
		resolver: &MailerResolver{},
	}

	return this
}

type MailerIntegrationBean struct {
	genetic  *Genetic
	logger   *zap.Logger
	resolver *MailerResolver
}

func (this MailerIntegrationBean) Resolver() *MailerResolver {
	return this.resolver
}

func (this MailerIntegrationBean) Migrate(tx *gorm.DB, driver string) error {
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

func (this MailerIntegrationBean) Dependencies() []module.Bean {
	panic("implement me")
}

func (this MailerIntegrationBean) Send(message model.Message) error {
	if true {
		return nil
	}

	if this.genetic.Reroute.Enabled {
		// TODO: check matching
		message.Recipient = this.genetic.Reroute.Recipient
	}

	dialer := &gomail.Dialer{} // gomail.NewDialer(host, port, username, password)

	return message.Send(dialer)
}
