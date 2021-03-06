package connect

import "net/http"

type MockRoundTrip struct {
	Callback func(*http.Request) (*http.Response, error)
}

func (mock MockRoundTrip) RoundTrip(request *http.Request) (*http.Response, error) {
	return mock.Callback(request)
}
