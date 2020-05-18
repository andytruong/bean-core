package namespace

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"bean/pkg/namespace/model/dto"
	"bean/pkg/util"
)

func Test_Create(t *testing.T) {
	ass := assert.New(t)
	input := dto.NamespaceCreateInput{
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
		},
		Context: &dto.NamespaceCreateContext{
			UserID: "xxxxxxxx",
		},
	}

	db := util.MockDatabase()
	module, err := NewNamespaceModule(db, util.MockLogger(), util.MockIdentifier())
	ass.NoError(err)
	util.MockInstall(module, db)

	{
		now := time.Now()
		outcome, err := module.Mutation.NamespaceCreate(context.Background(), input)
		ass.NoError(err)
		ass.Nil(outcome.Errors)
		ass.Equal(input.Object.Title, outcome.Namespace.Title)
		ass.Equal(input.Object.IsActive, outcome.Namespace.IsActive)
		ass.True(outcome.Namespace.CreatedAt.UnixNano() >= now.UnixNano())
		ass.True(outcome.Namespace.UpdatedAt.UnixNano() >= now.UnixNano())
	}
}
