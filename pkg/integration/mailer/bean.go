package mailer

import (
	"net/url"
	
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
	
	"bean/components/module"
	"bean/pkg/integration/mailer/model"
)

func NewMailerIntegration(genetic *Genetic) *MailerIntegrationBean {
	con, err := url.Parse(genetic.ConnectionUrl)
	if nil != err {
		panic(err)
	}
	
	username := con.Query().Get("username")
	if "" == username {
		panic("incorrect mailer.connectionUrl.username")
	}
	
	password := con.Query().Get("password")
	if "" == password {
		panic("incorrect mailer.connectionUrl.password")
	}
	
	this := &MailerIntegrationBean{
		genetic: genetic,
		dialer:  gomail.NewDialer(con.Host, 587, username, password),
	}
	
	return this
}

type MailerIntegrationBean struct {
	genetic *Genetic
	dialer  *gomail.Dialer
}

func (this MailerIntegrationBean) Migrate(tx *gorm.DB, driver string) error {
	panic("implement me")
}

func (this MailerIntegrationBean) Dependencies() []module.Bean {
	panic("implement me")
}

func (this MailerIntegrationBean) Send(message model.Message) error {
	if this.genetic.Reroute.Enabled {
		// TODO: check matching
		message.Recipient = this.genetic.Reroute.Recipient
	}
	
	return message.Send(this.dialer)
}
