package fixtures

import (
	"bean/components/scalar"
	"bean/pkg/namespace/model"
	"bean/pkg/namespace/model/dto"
	"bean/pkg/util/api"
)

func NamespaceCreateInputFixture(register bool) dto.NamespaceCreateInput {
	input := dto.NamespaceCreateInput{
		Object: dto.NamespaceCreateInputObject{},
	}

	input.Object = dto.NamespaceCreateInputObject{
		Kind:        model.NamespaceKindOrganisation,
		Title:       scalar.NilString("Home of QA"),
		IsActive:    true,
		Language:    api.LanguageAU,
		DomainNames: nil,
		Features: dto.NamespaceFeaturesInput{
			Register: register,
		},
	}

	input.Object.DomainNames = &dto.DomainNamesInput{}
	input.Object.DomainNames.Primary = &dto.DomainNameInput{
		Verified: scalar.NilBool(true),
		Value:    scalar.NilString("http://test.qa"),
		IsActive: scalar.NilBool(true),
	}

	input.Object.DomainNames.Secondary = []*dto.DomainNameInput{
		{
			Verified: scalar.NilBool(true),
			Value:    scalar.NilString("http://beta.test.qa"),
			IsActive: scalar.NilBool(true),
		},
		{
			Verified: scalar.NilBool(true),
			Value:    scalar.NilString("http://rc.test.qa"),
			IsActive: scalar.NilBool(false),
		},
	}

	return input
}
