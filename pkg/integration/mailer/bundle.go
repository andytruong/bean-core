package mailer

import (
	"context"
	"path"
	"runtime"

	"go.uber.org/zap"
	"gopkg.in/gomail.v2"

	"bean/components/connect"
	"bean/components/module"
	"bean/pkg/integration/mailer/model"
)

func NewMailerBundle(genetic *Config, logger *zap.Logger) *Bundle {
	this := &Bundle{
		config: genetic,
		logger: logger,
	}

	this.resolvers = newResoler(this)

	return this
}

type Bundle struct {
	module.AbstractBundle

	config    *Config
	logger    *zap.Logger
	resolvers map[string]interface{}
}

func (Bundle) Name() string {
	return "Mailer"
}

func (bundle Bundle) Migrate(ctx context.Context, driver string) error {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return nil
	}

	runner := connect.Runner{
		Logger: bundle.logger,
		Driver: driver,
		Bundle: "integration.mailer",
		Dir:    path.Dir(filename) + "/model/migration/",
	}

	return runner.Run(ctx)
}

func (bundle Bundle) Dependencies() []module.Bundle {
	panic("implement me")
}

func (bundle Bundle) Send(message model.Message) error {
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

func (bundle *Bundle) GraphqlResolver() map[string]interface{} {
	return bundle.resolvers
}
