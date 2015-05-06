package filepicker_test

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/filepicker/filepicker-go/filepicker"
)

func TestRemove(t *testing.T) {
	tests := []struct {
		Opt *filepicker.RemoveOpts
		Url string
	}{
		{
			Opt: nil,
			Url: "http://www.filepicker.io/api/file/2HHH3?key=0KKK1",
		},
		{
			Opt: &filepicker.RemoveOpts{},
			Url: "http://www.filepicker.io/api/file/2HHH3?key=0KKK1",
		},
		{
			Opt: &filepicker.RemoveOpts{
				Security: dummySecurity,
			},
			Url: "http://www.filepicker.io/api/file/2HHH3?key=0KKK1&policy=P&signature=S",
		},
	}

	var reqUrl, reqMethod, reqBody string
	handler := func(w http.ResponseWriter, req *http.Request) {
		body, _ := ioutil.ReadAll(req.Body)
		reqBody = string(body)
		reqUrl = req.URL.String()
		reqMethod = req.Method
	}

	blob := filepicker.NewBlob(FakeHandle)
	client := filepicker.NewClient(FakeApiKey)
	mock := MockServer(t, client, handler)
	defer mock.Close()

	for _, test := range tests {
		if err := client.Remove(blob, test.Opt); err != nil {
			t.Errorf("want err == nil; got %v", err)
		}
		if test.Url != reqUrl {
			t.Errorf("want test.Url == reqUrl; got %q != %q", test.Url, reqUrl)
		}
		if reqMethod != "DELETE" {
			t.Errorf("want reqMethod == DELETE; got %s", reqMethod)
		}
		if reqBody != "" {
			t.Errorf("want reqBody == ``; got %q", reqBody)
		}
	}
}

func TestRemoveError(t *testing.T) {
	fperr, handler := ErrorHandler(filepicker.ErrRmFileCannotBeFound)

	blob := filepicker.NewBlob(FakeHandle)
	client := filepicker.NewClient(FakeApiKey)

	mock := MockServer(t, client, handler)
	defer mock.Close()

	if err := client.Remove(blob, nil); err != fperr {
		t.Errorf("want err == fperr(%v), got %v", fperr, err)
	}
}
