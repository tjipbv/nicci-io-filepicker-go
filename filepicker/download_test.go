package filepicker_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/filepicker/filepicker-go/filepicker"
)

const downloadFileContent = "DOWNLOADTEST"

func TestDownloadTo(t *testing.T) {
	tests := []struct {
		Opt *filepicker.DownloadOpts
		URL string
	}{
		{
			Opt: nil,
			URL: "http://www.filepicker.io/api/file/2HHH3",
		},
		{
			Opt: &filepicker.DownloadOpts{
				Base64Decode: true,
			},
			URL: "http://www.filepicker.io/api/file/2HHH3?base64decode=true",
		},
		{
			Opt: &filepicker.DownloadOpts{
				Base64Decode: true,
				Security:     dummySecurity,
			},
			URL: "http://www.filepicker.io/api/file/2HHH3?base64decode=true&policy=P&signature=S",
		},
	}

	var reqURL, reqMethod, reqBody string
	handler := func(w http.ResponseWriter, req *http.Request) {
		body, _ := ioutil.ReadAll(req.Body)
		reqBody = string(body)
		reqURL = req.URL.String()
		reqMethod = req.Method
		if _, err := w.Write([]byte(downloadFileContent)); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	blob := filepicker.NewBlob(FakeHandle)
	client := filepicker.NewClient(FakeApiKey)
	mock := MockServer(t, client, handler)
	defer mock.Close()

	for i, test := range tests {
		var buff bytes.Buffer
		byteRead, err := client.DownloadTo(blob, test.Opt, &buff)
		if err != nil {
			t.Errorf("want err == nil; got %v (i:%d)", err, i)
		}
		if l := int64(len(downloadFileContent)); l != byteRead {
			t.Errorf("want byteRead == %d; got %d (i:%d)", l, byteRead, i)
		}
		if test.URL != reqURL {
			t.Errorf("want reqURL == %q; got %q (i:%d)", test.URL, reqURL, i)
		}
		if reqMethod != "GET" {
			t.Errorf("want reqMethod == GET; got %q (i:%d)", reqMethod, i)
		}
		content := string(buff.Bytes())
		if content != downloadFileContent {
			t.Errorf("want content == %q; got %q (i:%d)", downloadFileContent, content, i)
		}
	}
}

func TestDownloadToError(t *testing.T) {
	fperr, handler := ErrorHandler(dummyErrStr)

	blob := filepicker.NewBlob(FakeHandle)
	client := filepicker.NewClient(FakeApiKey)
	mock := MockServer(t, client, handler)
	defer mock.Close()

	var buff bytes.Buffer
	switch byteRead, err := client.DownloadTo(blob, nil, &buff); {
	case byteRead != 0:
		t.Errorf("want byteRead == 0; got %d", byteRead)
	case err.Error() != fperr.Error():
		t.Errorf("want error message == %q; got %q", fperr, err)
	}
}

func TestDownloadToFile(t *testing.T) {
	var testCounter int
	tests := []struct {
		Opt     *filepicker.DownloadOpts
		URL     string
		XName   string
		FileDir string
		Path    string
	}{
		{
			Opt:     nil,
			URL:     "http://www.filepicker.io/api/file/2HHH3",
			XName:   "document.txt",
			FileDir: ".",
			Path:    "./document.txt",
		},
		{
			Opt:     nil,
			URL:     "http://www.filepicker.io/api/file/2HHH3",
			XName:   "",
			FileDir: "something.txt",
			Path:    "./something.txt",
		},
		{
			Opt:     nil,
			URL:     "http://www.filepicker.io/api/file/2HHH3",
			XName:   "document.txt",
			FileDir: "./doc.txt",
			Path:    "./doc.txt",
		},
		{
			Opt: &filepicker.DownloadOpts{
				Base64Decode: true,
			},
			URL:     "http://www.filepicker.io/api/file/2HHH3?base64decode=true",
			XName:   "dc.txt",
			FileDir: ".",
			Path:    "./dc.txt",
		},
		{
			Opt: &filepicker.DownloadOpts{
				Base64Decode: true,
				Security:     dummySecurity,
			},
			URL:     "http://www.filepicker.io/api/file/2HHH3?base64decode=true&policy=P&signature=S",
			XName:   "dc.txt",
			FileDir: ".",
			Path:    "./dc.txt",
		},
	}

	var reqURL, reqMethod, reqBody string
	handler := func(w http.ResponseWriter, req *http.Request) {
		body, _ := ioutil.ReadAll(req.Body)
		reqBody = string(body)
		reqURL = req.URL.String()
		reqMethod = req.Method
		w.Header().Set("X-File-Name", tests[testCounter].XName)
		testCounter++
		if _, err := w.Write([]byte(downloadFileContent)); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	blob := filepicker.NewBlob(FakeHandle)
	client := filepicker.NewClient(FakeApiKey)
	mock := MockServer(t, client, handler)
	defer mock.Close()

	for i, test := range tests {
		if err := client.DownloadToFile(blob, test.Opt, test.FileDir); err != nil {
			t.Errorf("want err == nil; got %v (i:%d)", err, i)
		}
		if test.URL != reqURL {
			t.Errorf("want reqURL == %q; got %q (i:%d)", test.URL, reqURL, i)
		}
		if reqMethod != "GET" {
			t.Errorf("want reqMethod == GET; got %s (i:%d)", reqMethod, i)
		}
		b, err := ioutil.ReadFile(test.Path)
		if err != nil {
			t.Errorf("want err == nil; got %v (i:%d)", err, i)
		}
		if content := string(b); content != downloadFileContent {
			t.Errorf("want content == %q; got %q (i:%d)", downloadFileContent, content, i)
		}
		if err := os.Remove(test.Path); err != nil {
			t.Errorf("want err == nil; got %v (i:%d)", err, i)
		}
	}
}

func TestDownloadToFileError(t *testing.T) {
	fperr, handler := ErrorHandler(dummyErrStr)

	blob := filepicker.NewBlob(FakeHandle)
	client := filepicker.NewClientStorage(FakeApiKey, filepicker.Azure)
	mock := MockServer(t, client, handler)
	defer mock.Close()

	if err := client.DownloadToFile(blob, nil, "."); err.Error() != fperr.Error() {
		t.Errorf("want error message == %q; got %q", fperr, err)
	}
}
