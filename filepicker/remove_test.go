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
		URL string
	}{
		{
			Opt: nil,
			URL: "http://www.filepicker.io/api/file/2HHH3?key=0KKK1",
		},
		{
			Opt: &filepicker.RemoveOpts{},
			URL: "http://www.filepicker.io/api/file/2HHH3?key=0KKK1",
		},
		{
			Opt: &filepicker.RemoveOpts{
				Security: dummySecurity,
			},
			URL: "http://www.filepicker.io/api/file/2HHH3?key=0KKK1&policy=P&signature=S",
		},
	}

	var reqURL, reqMethod, reqBody string
	handler := func(w http.ResponseWriter, req *http.Request) {
		body, _ := ioutil.ReadAll(req.Body)
		reqBody = string(body)
		reqURL = req.URL.String()
		reqMethod = req.Method
	}

	blob := filepicker.NewBlob(FakeHandle)
	client := filepicker.NewClient(FakeApiKey)
	mock := MockServer(t, client, handler)
	defer mock.Close()

	for i, test := range tests {
		if err := client.Remove(blob, test.Opt); err != nil {
			t.Errorf("want err == nil; got %v (i:%d)", err, i)
		}
		if test.URL != reqURL {
			t.Errorf("want reqURL == %q; got %q (i:%d)", test.URL, reqURL, i)
		}
		if reqMethod != "DELETE" {
			t.Errorf("want reqMethod == DELETE; got %s (i:%d)", reqMethod, i)
		}
		if reqBody != "" {
			t.Errorf("want reqBody == ``; got %q (i:%d)", reqBody, i)
		}
	}
}

func TestRemoveError(t *testing.T) {
	fperr, handler := ErrorHandler(dummyErrStr)

	blob := filepicker.NewBlob(FakeHandle)
	client := filepicker.NewClient(FakeApiKey)

	mock := MockServer(t, client, handler)
	defer mock.Close()

	if err := client.Remove(blob, nil); err.Error() != fperr.Error() {
		t.Errorf("want error message == %q; got %q", fperr, err)
	}
}
