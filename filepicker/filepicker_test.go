package filepicker_test

import (
	"errors"
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

var dummyErrStr = "dummy error"

func (mt *MockedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.URL.Scheme = "http" // Disable SSL
	if userAgentID := req.Header.Get("User-Agent"); userAgentID != filepicker.UserAgentID {
		panic("filepicker: invalid User-Agent header field: " + userAgentID)
	}
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

func ErrorHandler(errstr string) (error, func(w http.ResponseWriter, req *http.Request)) {
	return errors.New("filepicker: 404 - " + dummyErrStr), func(w http.ResponseWriter, req *http.Request) {
		http.Error(w, errstr, 404)
	}
}

func TestBlobHandle(t *testing.T) {
	if blob := filepicker.NewBlob(FakeHandle); blob.Handle() != FakeHandle {
		t.Errorf("want blob.Handle() == %q; got %q", FakeHandle, blob.Handle())
	}
}
