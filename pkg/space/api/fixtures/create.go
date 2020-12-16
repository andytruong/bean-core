package fixtures

import (
	"bean/components/scalar"
	"bean/pkg/infra/api"
	"bean/pkg/space/model"
	"bean/pkg/space/model/dto"
)

func SpaceCreateInputFixture(register bool) dto.SpaceCreateInput {
	input := dto.SpaceCreateInput{
		Object: dto.SpaceCreateInputObject{},
	}

	input.Object = dto.SpaceCreateInputObject{
		Kind:        model.SpaceKindOrganisation,
		Title:       scalar.NilString("Home of QA"),
		IsActive:    true,
		Language:    api.LanguageAU,
		DomainNames: nil,
		Features: dto.SpaceFeaturesInput{
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
