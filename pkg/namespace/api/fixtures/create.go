package fixtures

import (
	"bean/pkg/namespace/model"
	"bean/pkg/namespace/model/dto"
	"bean/pkg/util"
	"bean/pkg/util/api"
)

func NamespaceCreateInputFixture(register bool) dto.NamespaceCreateInput {
	input := dto.NamespaceCreateInput{
		Object: dto.NamespaceCreateInputObject{},
		Context: dto.NamespaceCreateContext{
			UserID: "xxxxxxxx",
		},
	}

	input.Object = dto.NamespaceCreateInputObject{
		Kind:        model.NamespaceKindOrganisation,
		Title:       util.NilString("Home of QA"),
		IsActive:    true,
		Language:    api.LanguageAU,
		DomainNames: nil,
		Features: dto.NamespaceFeaturesInput{
			Register: register,
		},
	}

	input.Object.DomainNames = &dto.DomainNamesInput{}
	input.Object.DomainNames.Primary = &dto.DomainNameInput{
		Verified: util.NilBool(true),
		Value:    util.NilString("http://test.qa"),
		IsActive: util.NilBool(true),
	}

	input.Object.DomainNames.Secondary = []*dto.DomainNameInput{
		{
			Verified: util.NilBool(true),
			Value:    util.NilString("http://beta.test.qa"),
			IsActive: util.NilBool(true),
		},
		{
			Verified: util.NilBool(true),
			Value:    util.NilString("http://rc.test.qa"),
			IsActive: util.NilBool(false),
		},
	}

	return input
}
