package util

import "net/http"

type MockRoundTrip struct {
	Callback func(*http.Request) (*http.Response, error)
}

func (this MockRoundTrip) RoundTrip(request *http.Request) (*http.Response, error) {
	return this.Callback(request)
}
