package infra

import (
	"bean/components/module"
	"bean/pkg/access"
	"bean/pkg/integration/mailer"
	"bean/pkg/integration/s3"
	"bean/pkg/space"
	"bean/pkg/user"
)

type (
	beans struct {
		can    *Can
		user   *user.UserBean
		space  *space.SpaceBean
		access *access.AccessBean
		s3     *s3.S3IntegrationBean
		mailer *mailer.MailerIntegrationBean
	}
)

func (this *beans) List() []module.Bean {
	mUser, _ := this.User()
	mSpace, _ := this.Space()
	mAccess, _ := this.Access()
	mMailer, _ := this.Mailer()
	
	return []module.Bean{mUser, mSpace, mAccess, mMailer}
}

func (this *beans) User() (*user.UserBean, error) {
	var err error
	
	if nil == this.user {
		db, err := this.can.dbs.master()
		if nil != err {
			return nil, err
		}
		
		this.user = user.NewUserBean(
			db,
			this.can.logger,
			this.can.Identifier(),
		)
	}
	
	return this.user, err
}

func (this *beans) Space() (*space.SpaceBean, error) {
	var err error
	
	if nil == this.space {
		db, err := this.can.dbs.master()
		if nil != err {
			return nil, err
		}
		
		mUser, err := this.User()
		if nil != err {
			return nil, err
		}
		
		this.space = space.NewSpaceBean(
			db,
			this.can.logger,
			this.can.Identifier(),
			mUser,
			this.can.Beans.Space,
		)
	}
	
	return this.space, err
}

func (this *beans) Access() (*access.AccessBean, error) {
	if nil == this.access {
		db, err := this.can.dbs.master()
		if nil != err {
			return nil, err
		}
		
		mUser, err := this.User()
		if nil != err {
			return nil, err
		}
		
		mSpace, err := this.Space()
		if nil != err {
			return nil, err
		}
		
		this.access = access.NewAccessBean(
			db,
			this.can.Identifier(),
			this.can.logger,
			mUser,
			mSpace,
			this.can.Beans.Access,
		)
	}
	
	return this.access, nil
}

func (this *beans) Mailer() (*mailer.MailerIntegrationBean, error) {
	if nil == this.mailer {
		this.mailer = mailer.NewMailerIntegration(this.can.Beans.Integration.Mailer, this.can.logger)
	}
	
	return this.mailer, nil
}

func (this *beans) S3() (*s3.S3IntegrationBean, error) {
	if nil == this.s3 {
		db, err := this.can.dbs.master()
		if nil != err {
			return nil, err
		}
		
		this.s3 = s3.NewS3Integration(db, this.can.Identifier(), this.can.logger, this.can.Beans.Integration.S3)
	}
	
	return this.s3, nil
}
