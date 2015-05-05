package filepicker_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/filepicker/filepicker-go/filepicker"
)

func TestWrite(t *testing.T) {
	var testCounter int
	tests := []struct {
		Src *filepicker.Blob
		Opt *filepicker.WriteOpts
		Url string
	}{
		{
			Src: filepicker.NewBlob(FakeHandle),
			Opt: nil,
			Url: `http://www.filepicker.io/api/file/2HHH3`,
		},
		{
			Src: filepicker.NewBlob(FakeHandle),
			Opt: &filepicker.WriteOpts{
				Base64Decode: true,
			},
			Url: `http://www.filepicker.io/api/file/2HHH3?base64decode=true`,
		},
		{
			Src: filepicker.NewBlob(FakeHandle),
			Opt: &filepicker.WriteOpts{
				Base64Decode: true,
				Security:     dummySecurity,
			},
			Url: `http://www.filepicker.io/api/file/2HHH3?base64decode=true&policy=P&signature=S`,
		},
	}

	filename := tempFile(t)
	defer os.Remove(filename)

	var reqUrl, reqMethod, reqBody string
	writeHandle := func(w http.ResponseWriter, req *http.Request) {
		body, _ := ioutil.ReadAll(req.Body)
		reqBody = string(body)
		reqUrl = req.URL.String()
		reqMethod = req.Method
		blob := &filepicker.Blob{}
		data, err := json.Marshal(blob)
		testCounter++
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if _, err := w.Write(data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	client := filepicker.NewClient(FakeApiKey)
	mock := MockServer(t, client, writeHandle)
	defer mock.Close()

	for _, test := range tests {
		blob, err := client.Write(test.Src, filename, test.Opt)
		if err != nil {
			t.Errorf(`want err == nil, got %v`, err)
		}
		if blob == nil {
			t.Error(`want blob != nil, got nil`)
		}
		if test.Url != reqUrl {
			t.Errorf(`want test.Url == reqUrl; got %q != %q`, test.Url, reqUrl)
		}
		if reqMethod != `POST` {
			t.Errorf(`want reqMethod == POST; got %s`, reqMethod)
		}
	}
}

func TestWriteError(t *testing.T) {
	fperr, handle := ErrorHandler(filepicker.ErrCannotFindWriteBlob)

	client := filepicker.NewClient(FakeApiKey)
	mock := MockServer(t, client, handle)
	defer mock.Close()

	blob := filepicker.NewBlob("XYZ")
	filename := tempFile(t)
	defer os.Remove(filename)

	switch blob, err := client.Write(blob, filename, nil); {
	case blob != nil:
		t.Errorf("want blob == nil, got %v", blob)
	case err != fperr:
		t.Errorf("want err == fperr(%v), got %v", fperr, err)
	}
}

func TestWriteErrorNoFile(t *testing.T) {
	blob := filepicker.NewBlob("XYZ")
	client := filepicker.NewClient(FakeApiKey)
	switch blob, err := client.Write(blob, "unknown.unknown.file", nil); {
	case blob != nil:
		t.Error("want blob == nil, got %v", blob)
	case err == nil:
		t.Error("want err != nil, got nil")
	}
}

func TestWriteUrl(t *testing.T) {
	const TestUrl = `https://www.filepicker.com/image.png`
	var testCounter int
	tests := []struct {
		Src *filepicker.Blob
		Opt *filepicker.WriteOpts
		Url string
	}{
		{
			Src: filepicker.NewBlob(FakeHandle),
			Opt: nil,
			Url: `http://www.filepicker.io/api/file/2HHH3`,
		},
		{
			Src: filepicker.NewBlob(FakeHandle),
			Opt: &filepicker.WriteOpts{
				Base64Decode: true,
			},
			Url: `http://www.filepicker.io/api/file/2HHH3?base64decode=true`,
		},
		{
			Src: filepicker.NewBlob(FakeHandle),
			Opt: &filepicker.WriteOpts{
				Base64Decode: true,
				Security:     dummySecurity,
			},
			Url: `http://www.filepicker.io/api/file/2HHH3?base64decode=true&policy=P&signature=S`,
		},
	}

	var reqUrl, reqMethod, reqBody string
	writeHandle := func(w http.ResponseWriter, req *http.Request) {
		body, _ := ioutil.ReadAll(req.Body)
		reqBody = string(body)
		reqUrl = req.URL.String()
		reqMethod = req.Method
		blob := &filepicker.Blob{}
		data, err := json.Marshal(blob)
		testCounter++
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if _, err := w.Write(data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	client := filepicker.NewClient(FakeApiKey)
	mock := MockServer(t, client, writeHandle)
	defer mock.Close()

	for _, test := range tests {
		blob, err := client.WriteURL(test.Src, TestUrl, test.Opt)
		if err != nil {
			t.Errorf(`want err == nil, got %v`, err)
		}
		if test.Url != reqUrl {
			t.Errorf(`want test.Url == reqUrl; got %q != %q`, test.Url, reqUrl)
		}
		if reqMethod != `POST` {
			t.Errorf(`want reqMethod == POST; got %s`, reqMethod)
		}
		if TestUrlEsc := "url=" + url.QueryEscape(TestUrl); reqBody != TestUrlEsc {
			t.Errorf(`want reqBody == TestUrlEsc; got %q != %q`, reqBody, TestUrlEsc)
		}
		if blob == nil {
			t.Error(`want blob != nil`)
		}
	}
}

func TestWriteURLError(t *testing.T) {
	fperr, handle := ErrorHandler(filepicker.ErrWriteUrlUnreachable)

	blob := filepicker.NewBlob(FakeHandle)
	client := filepicker.NewClient(FakeApiKey)
	mock := MockServer(t, client, handle)
	defer mock.Close()

	switch blob, err := client.WriteURL(blob, "http://www.address.fp", nil); {
	case blob != nil:
		t.Errorf("want blob == nil, got %v", blob)
	case err != fperr:
		t.Errorf("want err == fperr(%v), got %v", fperr, err)
	}
}
