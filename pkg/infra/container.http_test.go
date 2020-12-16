package infra

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func Test_Authorization(t *testing.T) {
	t.Parallel()
	ass := assert.New(t)
	can := NewMockCan()
	router := can.HttpRouter(mux.NewRouter())

	t.Run("ping without auth", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/query", strings.NewReader(`{ "query": "{ __schema { queryType { kind } } }" }`))
		r.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, r)

		ass.Equal(http.StatusOK, w.Code)
		ass.Contains(w.Body.String(), `{"data":{"__schema":{"queryType":{"kind":"OBJECT"}}}}`)
	})

	t.Run("ping with invalid auth", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/query", strings.NewReader(`{ "query": "{ "query": "{ __schema { queryType { kind } } }" }" }`))
		r.Header.Set("Content-Type", "application/json")
		r.Header.Add("Authorization", "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIwMUVBMzc4RTlFU01STjJaWk5OMURZN04xOCIsImV4cCI6MTU5MTQ0MDgyNywianRpIjoiMDFFQTNROE1WTThDWkY3WFRZTllQWE1OWlYiLCJpYXQiOjE1OTE0NDA1MjcsImlzcyI6ImFjY2VzcyIsInN1YiI6IjAxRUEzNzg0VEU5NEtLVkNUS1I2VkoyR1dDIn0.BO36niHSF3Svzg4oIQ7A8bEScQYrWbvIlBZ5ExakOoEd5CZGuRQbAQRcF0skiqQz8cdVHb3pkcm7LUkJ7zi7WXKdnhd7M-NceGmwQ0XJ9NE9eZvYP5swFxjxVVYxTxjfWQp-5buP3UXkLeL2UhUINsFYJpxQUWKxLG-vdkCzRkcNH8VBkB-XTAfg7lX4ESGObVo-AxyxzrSPj2TWGNHnd5WrB6nFmz_up6vJ89aiizDI7zVnku-lJzPW0AJmiBFAyTD6y9WN0uKdBrGEzJ3wfW8EIadHgqcP7RmCF-XVD4ILIU3nwg-DQ8SQgjBpgcRyPTAawkOsIR6ubfQRS_J21A")
		router.ServeHTTP(w, r)

		ass.Equal(http.StatusForbidden, w.Code)
		ass.Contains(w.Body.String(), `{"errors":[{"message":"token is expired by`)
	})

	t.Run("ping with valid auth", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/query", strings.NewReader(`{ "query": "{ __schema { queryType { kind } } }" }`))
		r.Header.Set("Content-Type", "application/json")
		r.Header.Add("Authorization", func() string {
			claims := jwt.StandardClaims{Id: "1"}
			access, _ := can.bundles.Access()
			token, _ := access.JwtService.Sign(claims)

			return "Bearer " + token
		}())
		router.ServeHTTP(w, r)

		ass.Equal(http.StatusOK, w.Code)
		ass.Contains(w.Body.String(), `{"data":{"__schema":{"queryType":{"kind":"OBJECT"}}}}`)
	})
}
