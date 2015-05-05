package filepicker_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"testing"

	"github.com/filepicker/filepicker-go/filepicker"
)

const tmpFileContent = "STORETEST"

func tempFile(t *testing.T) (name string) {
	file, err := ioutil.TempFile("", "FP")
	if err != nil {
		t.Fatalf("want err == nil; got %v", err)
	}
	defer file.Close()
	if _, err := file.WriteString(tmpFileContent); err != nil {
		t.Fatalf("want err == nil; got %v", err)
	}
	return file.Name()
}

func TestStore(t *testing.T) {
	var testCounter int
	tests := []struct {
		Opt  *filepicker.StoreOpts
		Url  string
		Blob *filepicker.Blob
	}{
		{
			Opt:  nil,
			Url:  `http://www.filepicker.io/api/store/S3?key=0KKK1`,
			Blob: &filepicker.Blob{},
		},
		{
			Opt: &filepicker.StoreOpts{
				Location: filepicker.Azure,
			},
			Url:  `http://www.filepicker.io/api/store/azure?key=0KKK1&location=azure`,
			Blob: &filepicker.Blob{},
		},
		{
			Opt: &filepicker.StoreOpts{
				Filename: "file.txt",
			},
			Url: `http://www.filepicker.io/api/store/S3?filename=file.txt&key=0KKK1`,
			Blob: &filepicker.Blob{
				Filename: "file.txt",
			},
		},
		{
			Opt: &filepicker.StoreOpts{
				Path:         "path",
				Base64Decode: true,
			},
			Url:  `http://www.filepicker.io/api/store/S3?base64decode=true&key=0KKK1&path=path`,
			Blob: &filepicker.Blob{},
		},
	}

	filename := tempFile(t)
	defer os.Remove(filename)

	var reqUrl, reqMethod, reqBody string
	storeHandle := func(w http.ResponseWriter, req *http.Request) {
		body, _ := ioutil.ReadAll(req.Body)
		reqBody = string(body)
		reqUrl = req.URL.String()
		reqMethod = req.Method
		data, err := json.Marshal(tests[testCounter].Blob)
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
	mock := MockServer(t, client, storeHandle)
	defer mock.Close()

	for _, test := range tests {
		blob, err := client.Store(filename, test.Opt)
		if err != nil {
			t.Errorf(`want err == nil, got %v`, err)
		}
		if test.Url != reqUrl {
			t.Errorf(`want test.Url == reqUrl; got %q != %q`, test.Url, reqUrl)
		}
		if reqMethod != `POST` {
			t.Errorf(`want reqMethod == POST; got %s`, reqMethod)
		}
		if !reflect.DeepEqual(*blob, *test.Blob) {
			t.Errorf(`want *test.blob(%v) == *blob(%v)`, *test.Blob, *blob)
		}
	}
}

func TestStoreError(t *testing.T) {
	fperr, handle := ErrorHandler(filepicker.ErrFileStoreUnreachable)

	client := filepicker.NewClient(FakeApiKey)
	mock := MockServer(t, client, handle)
	defer mock.Close()

	filename := tempFile(t)
	defer os.Remove(filename)

	switch blob, err := client.Store(filename, nil); {
	case blob != nil:
		t.Error("want blob == nil, got %v", blob)
	case err != fperr:
		t.Errorf("want err == fperr(%v), got %v", fperr, err)
	}
}

func TestStoreErrorNoFile(t *testing.T) {
	client := filepicker.NewClient(FakeApiKey)
	switch blob, err := client.Store("unknown.unknown.file", nil); {
	case blob != nil:
		t.Error("want blob == nil, got %v", blob)
	case err == nil:
		t.Error("want err != nil, got nil")
	}
}

func TestStoreUrl(t *testing.T) {
	const TestUrl = `https://www.filepicker.com/image.png`
	var testCounter int
	tests := []struct {
		Opt *filepicker.StoreOpts
		Url string
	}{
		{
			Opt: nil,
			Url: `http://www.filepicker.io/api/store/S3?key=0KKK1`,
		},
		{
			Opt: &filepicker.StoreOpts{
				Location: filepicker.Azure,
			},
			Url: `http://www.filepicker.io/api/store/azure?key=0KKK1&location=azure`,
		},
		{
			Opt: &filepicker.StoreOpts{
				Access: "public",
			},
			Url: `http://www.filepicker.io/api/store/S3?access=public&key=0KKK1`,
		},
		{
			Opt: &filepicker.StoreOpts{
				Path:         "path",
				Base64Decode: true,
			},
			Url: `http://www.filepicker.io/api/store/S3?base64decode=true&key=0KKK1&path=path`,
		},
	}

	var reqUrl, reqMethod, reqBody string
	storeHandle := func(w http.ResponseWriter, req *http.Request) {
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
	mock := MockServer(t, client, storeHandle)
	defer mock.Close()

	for _, test := range tests {
		blob, err := client.StoreURL(TestUrl, test.Opt)
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

func TestStoreURLError(t *testing.T) {
	fperr, handle := ErrorHandler(filepicker.ErrRemoteUrlUnreachable)

	client := filepicker.NewClient(FakeApiKey)
	mock := MockServer(t, client, handle)
	defer mock.Close()

	switch blob, err := client.StoreURL("some_link", nil); {
	case blob != nil:
		t.Error("want blob == nil, got %v", blob)
	case err != fperr:
		t.Errorf("want err == fperr(%v), got %v", fperr, err)
	}
}
