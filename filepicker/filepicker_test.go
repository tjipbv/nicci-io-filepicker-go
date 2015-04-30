package filepicker_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/filepicker/filepicker-go/filepicker"
)

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
				req.URL.Scheme = "http"
				fmt.Println(*req)
				fmt.Println(url.Parse(mockedServer.URL))
				return url.Parse(mockedServer.URL)
			},
		},
	}
	return mockedServer
}

func TestTmp(t *testing.T) {
	printerHandle := func(w http.ResponseWriter, req *http.Request) {
		fmt.Println(*req)
		fmt.Println("================")
		byteData, _ := ioutil.ReadAll(req.Body)
		fmt.Println(string(byteData))
	}

	client := filepicker.NewClient("AAA")
	mock := MockServer(t, client, printerHandle)
	defer mock.Close()

	client.StoreURL("https://url.com/sth", filepicker.StoreOpts{})
}
