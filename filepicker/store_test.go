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

const storeFileContent = "STORETEST"

func tempFile(t *testing.T) (name string) {
	file, err := ioutil.TempFile("", "FP")
	if err != nil {
		t.Fatalf("want err == nil; got %v", err)
	}
	defer file.Close()
	if _, err := file.WriteString(storeFileContent); err != nil {
		t.Fatalf("want err == nil; got %v", err)
	}
	return file.Name()
}

func testHandle(reqUrl, reqMethod, reqBody *string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		body, _ := ioutil.ReadAll(req.Body)
		*reqBody = string(body)
		*reqUrl = req.URL.String()
		*reqMethod = req.Method
		blob := &filepicker.Blob{}
		data, err := json.Marshal(blob)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if _, err := w.Write(data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func TestStore(t *testing.T) {
	tests := []struct {
		Opt *filepicker.StoreOpts
		Url string
	}{
		{
			Opt: nil,
			Url: "http://www.filepicker.io/api/store/S3?key=0KKK1",
		},
		{
			Opt: &filepicker.StoreOpts{
				Location: filepicker.Azure,
			},
			Url: "http://www.filepicker.io/api/store/azure?key=0KKK1&location=azure",
		},
		{
			Opt: &filepicker.StoreOpts{
				Filename: "file.txt",
			},
			Url: "http://www.filepicker.io/api/store/S3?filename=file.txt&key=0KKK1",
		},
		{
			Opt: &filepicker.StoreOpts{
				Path:         "path",
				Base64Decode: true,
			},
			Url: "http://www.filepicker.io/api/store/S3?base64decode=true&key=0KKK1&path=path",
		},
	}

	filename := tempFile(t)
	defer os.Remove(filename)

	var reqUrl, reqMethod, reqBody string
	handler := testHandle(&reqUrl, &reqMethod, &reqBody)
	client := filepicker.NewClient(FakeApiKey)
	mock := MockServer(t, client, handler)
	defer mock.Close()

	for _, test := range tests {
		blob, err := client.Store(filename, test.Opt)
		if err != nil {
			t.Errorf("want err == nil; got %v", err)
		}
		if blob == nil {
			t.Error("want blob != nil; got nil")
		}
		if test.Url != reqUrl {
			t.Errorf("want test.Url == reqUrl; got %q != %q", test.Url, reqUrl)
		}
		if reqMethod != "POST" {
			t.Errorf("want reqMethod == POST; got %s", reqMethod)
		}
	}
}

func TestStoreError(t *testing.T) {
	fperr, handler := ErrorHandler(filepicker.ErrFileStoreUnreachable)

	client := filepicker.NewClient(FakeApiKey)
	mock := MockServer(t, client, handler)
	defer mock.Close()

	filename := tempFile(t)
	defer os.Remove(filename)

	switch blob, err := client.Store(filename, nil); {
	case blob != nil:
		t.Errorf("want blob == nil; got %v", blob)
	case err != fperr:
		t.Errorf("want err == fperr(%v); got %v", fperr, err)
	}
}

func TestStoreErrorNoFile(t *testing.T) {
	client := filepicker.NewClient(FakeApiKey)
	switch blob, err := client.Store("unknown.unknown.file", nil); {
	case blob != nil:
		t.Errorf("want blob == nil; got %v", blob)
	case err == nil:
		t.Error("want err != nil; got nil")
	}
}

func TestStoreUrl(t *testing.T) {
	const TestUrl = "https://www.filepicker.com/image.png"
	tests := []struct {
		Opt *filepicker.StoreOpts
		Url string
	}{
		{
			Opt: nil,
			Url: "http://www.filepicker.io/api/store/S3?key=0KKK1",
		},
		{
			Opt: &filepicker.StoreOpts{
				Location: filepicker.Azure,
			},
			Url: "http://www.filepicker.io/api/store/azure?key=0KKK1&location=azure",
		},
		{
			Opt: &filepicker.StoreOpts{
				Access: "public",
			},
			Url: "http://www.filepicker.io/api/store/S3?access=public&key=0KKK1",
		},
		{
			Opt: &filepicker.StoreOpts{
				Path:         "path",
				Base64Decode: true,
			},
			Url: "http://www.filepicker.io/api/store/S3?base64decode=true&key=0KKK1&path=path",
		},
	}

	var reqUrl, reqMethod, reqBody string
	handler := testHandle(&reqUrl, &reqMethod, &reqBody)
	client := filepicker.NewClient(FakeApiKey)
	mock := MockServer(t, client, handler)
	defer mock.Close()

	for _, test := range tests {
		blob, err := client.StoreURL(TestUrl, test.Opt)
		if err != nil {
			t.Errorf("want err == nil; got %v", err)
		}
		if test.Url != reqUrl {
			t.Errorf("want test.Url == reqUrl; got %q != %q", test.Url, reqUrl)
		}
		if reqMethod != "POST" {
			t.Errorf("want reqMethod == POST; got %s", reqMethod)
		}
		if TestUrlEsc := "url=" + url.QueryEscape(TestUrl); reqBody != TestUrlEsc {
			t.Errorf("want reqBody == TestUrlEsc; got %q != %q", reqBody, TestUrlEsc)
		}
		if blob == nil {
			t.Error("want blob != nil; got nil")
		}
	}
}

func TestStoreURLError(t *testing.T) {
	fperr, handler := ErrorHandler(filepicker.ErrRemoteUrlUnreachable)

	client := filepicker.NewClient(FakeApiKey)
	mock := MockServer(t, client, handler)
	defer mock.Close()

	switch blob, err := client.StoreURL("http://www.address.fp", nil); {
	case blob != nil:
		t.Errorf("want blob == nil; got %v", blob)
	case err != fperr:
		t.Errorf("want err == fperr(%v); got %v", fperr, err)
	}
}
