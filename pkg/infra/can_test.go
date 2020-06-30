package infra

import (
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func NewMockCan() *Can {
	_ = os.Setenv("DB_MASTER_URL", ":memory:")
	_ = os.Setenv("DB_SLAVE_URL", ":memory:")
	ctn, err := NewCan("../../config.yaml")

	if nil != err {
		panic(err)
	}

	return ctn
}

func Test(t *testing.T) {
	ass := assert.New(t)
	can := NewMockCan()
	id := can.Identifier()
	ass.NotNil(t, id)

	sv, err := can.beans.User()
	ass.NoError(err)
	ass.NotNil(sv)
	ass.Equal("128h0m0s", can.Beans.Access.SessionTimeout.String())
	ass.Equal(100, can.Beans.Namespace.Manager.MaxNumberOfManager)
	ass.Equal("01EBWB516AP6BQD7", can.Beans.Integration.S3.Key)
}

func Test_Request_JWT(t *testing.T) {
	ass := assert.New(t)
	can := NewMockCan()

	r, err := http.NewRequest("GET", "/query", nil)
	ass.NoError(err)

	{
		r.Header.Add("Authorization", "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIwMUVBMzc4RTlFU01STjJaWk5OMURZN04xOCIsImV4cCI6MTU5MTQ0MDgyNywianRpIjoiMDFFQTNROE1WTThDWkY3WFRZTllQWE1OWlYiLCJpYXQiOjE1OTE0NDA1MjcsImlzcyI6ImFjY2VzcyIsInN1YiI6IjAxRUEzNzg0VEU5NEtLVkNUS1I2VkoyR1dDIn0.BO36niHSF3Svzg4oIQ7A8bEScQYrWbvIlBZ5ExakOoEd5CZGuRQbAQRcF0skiqQz8cdVHb3pkcm7LUkJ7zi7WXKdnhd7M-NceGmwQ0XJ9NE9eZvYP5swFxjxVVYxTxjfWQp-5buP3UXkLeL2UhUINsFYJpxQUWKxLG-vdkCzRkcNH8VBkB-XTAfg7lX4ESGObVo-AxyxzrSPj2TWGNHnd5WrB6nFmz_up6vJ89aiizDI7zVnku-lJzPW0AJmiBFAyTD6y9WN0uKdBrGEzJ3wfW8EIadHgqcP7RmCF-XVD4ILIU3nwg-DQ8SQgjBpgcRyPTAawkOsIR6ubfQRS_J21A")
		err := can.beforeServeHTTP(r)
		ass.Contains(err.Error(), "token is expired by")
	}
}
