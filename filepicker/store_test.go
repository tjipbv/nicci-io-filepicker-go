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

func testHandle(reqURL, reqMethod, reqBody *string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		body, _ := ioutil.ReadAll(req.Body)
		*reqBody = string(body)
		*reqURL = req.URL.String()
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
		URL string
	}{
		{
			Opt: nil,
			URL: "http://www.filepicker.io/api/store/S3?key=0KKK1",
		},
		{
			Opt: &filepicker.StoreOpts{
				Location: filepicker.Azure,
			},
			URL: "http://www.filepicker.io/api/store/azure?key=0KKK1&location=azure",
		},
		{
			Opt: &filepicker.StoreOpts{
				Filename: "file.txt",
			},
			URL: "http://www.filepicker.io/api/store/S3?filename=file.txt&key=0KKK1",
		},
		{
			Opt: &filepicker.StoreOpts{
				Path:         "path",
				Base64Decode: true,
			},
			URL: "http://www.filepicker.io/api/store/S3?base64decode=true&key=0KKK1&path=path",
		},
	}

	filename := tempFile(t)
	defer os.Remove(filename)

	var reqURL, reqMethod, reqBody string
	handler := testHandle(&reqURL, &reqMethod, &reqBody)
	client := filepicker.NewClient(FakeApiKey)
	mock := MockServer(t, client, handler)
	defer mock.Close()

	for i, test := range tests {
		blob, err := client.Store(filename, test.Opt)
		if err != nil {
			t.Errorf("want err == nil; got %v (i:%d)", err, i)
		}
		if blob == nil {
			t.Errorf("want blob != nil; got nil (i:%d)", i)
		}
		if reqMethod != "POST" {
			t.Errorf("want reqMethod == POST; got %s (i:%d)", reqMethod, i)
		}
		if test.URL != reqURL {
			t.Errorf("want reqURL == %q; got %q (i:%d)", test.URL, reqURL, i)
		}
	}
}

func TestStoreError(t *testing.T) {
	fperr, handler := ErrorHandler(dummyErrStr)

	client := filepicker.NewClient(FakeApiKey)
	mock := MockServer(t, client, handler)
	defer mock.Close()

	filename := tempFile(t)
	defer os.Remove(filename)

	switch blob, err := client.Store(filename, nil); {
	case blob != nil:
		t.Errorf("want blob == nil; got %v", blob)
	case err.Error() != fperr.Error():
		t.Errorf("want error message == %q; got %q", fperr, err)
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

func TestStoreURL(t *testing.T) {
	const TestURL = "https://www.filepicker.com/image.png"
	tests := []struct {
		Opt *filepicker.StoreOpts
		URL string
	}{
		{
			Opt: nil,
			URL: "http://www.filepicker.io/api/store/S3?key=0KKK1",
		},
		{
			Opt: &filepicker.StoreOpts{
				Location: filepicker.Azure,
			},
			URL: "http://www.filepicker.io/api/store/azure?key=0KKK1&location=azure",
		},
		{
			Opt: &filepicker.StoreOpts{
				Access: "public",
			},
			URL: "http://www.filepicker.io/api/store/S3?access=public&key=0KKK1",
		},
		{
			Opt: &filepicker.StoreOpts{
				Path:         "path",
				Base64Decode: true,
			},
			URL: "http://www.filepicker.io/api/store/S3?base64decode=true&key=0KKK1&path=path",
		},
	}

	var reqURL, reqMethod, reqBody string
	handler := testHandle(&reqURL, &reqMethod, &reqBody)
	client := filepicker.NewClient(FakeApiKey)
	mock := MockServer(t, client, handler)
	defer mock.Close()

	for i, test := range tests {
		blob, err := client.StoreURL(TestURL, test.Opt)
		if err != nil {
			t.Errorf("want err == nil; got %v (i:%d)", err, i)
		}
		if test.URL != reqURL {
			t.Errorf("want reqURL == %q; got %q (i:%d)", test.URL, reqURL, i)
		}
		if reqMethod != "POST" {
			t.Errorf("want reqMethod == POST; got %s (i:%d)", reqMethod, i)
		}
		if TestURLEsc := "url=" + url.QueryEscape(TestURL); reqBody != TestURLEsc {
			t.Errorf("want reqBody == %q; got %q (i:%d)", TestURLEsc, reqBody, i)
		}
		if blob == nil {
			t.Errorf("want blob != nil; got nil (i:%d)", i)
		}
	}
}

func TestStoreURLError(t *testing.T) {
	fperr, handler := ErrorHandler(dummyErrStr)

	client := filepicker.NewClient(FakeApiKey)
	mock := MockServer(t, client, handler)
	defer mock.Close()

	switch blob, err := client.StoreURL("http://www.address.fp", nil); {
	case blob != nil:
		t.Errorf("want blob == nil; got %v", blob)
	case err.Error() != fperr.Error():
		t.Errorf("want error message == %q; got %q", fperr, err)
	}
}
