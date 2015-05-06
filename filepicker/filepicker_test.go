package filepicker_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/filepicker/filepicker-go/filepicker"
)

const (
	FakeApiKey = "0KKK1"
	FakeHandle = "2HHH3"
)

var dummySecurity = filepicker.Security{
	Policy:    "P",
	Signature: "S",
}

type MockedTransport struct {
	Transport http.Transport
}

func (mt *MockedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.URL.Scheme = "http" // Disable SSL
	return mt.Transport.RoundTrip(req)
}

func MockServer(t *testing.T, c *filepicker.Client, h http.HandlerFunc) *httptest.Server {
	mockedServer := httptest.NewServer(http.HandlerFunc(h))
	c.Client.Transport = &MockedTransport{
		Transport: http.Transport{
			Proxy: func(req *http.Request) (*url.URL, error) {
				return url.Parse(mockedServer.URL)
			},
		},
	}
	return mockedServer
}

func ErrorHandler(fperr filepicker.FPError) (
	filepicker.FPError, func(w http.ResponseWriter, req *http.Request)) {
	return fperr, func(w http.ResponseWriter, req *http.Request) {
		http.Error(w, fperr.Error(), int(fperr))
	}
}
