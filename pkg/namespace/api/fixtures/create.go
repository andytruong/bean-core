package fixtures

import (
	"bean/pkg/namespace/model/dto"
	"bean/pkg/util"
)

func NamespaceCreateInputFixture(register bool) dto.NamespaceCreateInput {
	return dto.NamespaceCreateInput{
		Object: &dto.NamespaceCreateInputObject{
			Title:    util.NilString("Home of QA"),
			IsActive: true,
			DomainNames: &dto.DomainNamesInput{
				Primary: &dto.DomainNameInput{
					Verified: util.NilBool(true),
					Value:    util.NilString("http://test.qa"),
					IsActive: util.NilBool(true),
				},
				Secondary: []*dto.DomainNameInput{
					{
						Verified: util.NilBool(true),
						Value:    util.NilString("http://beta.test.qa"),
						IsActive: util.NilBool(true),
					},
					{
						Verified: util.NilBool(true),
						Value:    util.NilString("http://rc.test.qa"),
						IsActive: util.NilBool(true),
					},
				},
			},
			Features: dto.NamespaceFeaturesInput{
				Register: register,
			},
		},
		Context: &dto.NamespaceCreateContext{
			UserID: "xxxxxxxx",
		},
	}
}
