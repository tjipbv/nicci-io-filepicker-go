package filepicker_test

import (
	"net/url"
	"os"
	"testing"

	"github.com/filepicker/filepicker-go/filepicker"
)

func TestWrite(t *testing.T) {
	tests := []struct {
		Opt *filepicker.WriteOpts
		URL string
	}{
		{
			Opt: nil,
			URL: "http://www.filepicker.io/api/file/2HHH3",
		},
		{
			Opt: &filepicker.WriteOpts{
				Base64Decode: true,
			},
			URL: "http://www.filepicker.io/api/file/2HHH3?base64decode=true",
		},
		{
			Opt: &filepicker.WriteOpts{
				Base64Decode: true,
				Security:     dummySecurity,
			},
			URL: "http://www.filepicker.io/api/file/2HHH3?base64decode=true&policy=P&signature=S",
		},
	}

	filename := tempFile(t)
	defer os.Remove(filename)

	var reqURL, reqMethod, reqBody string
	handler := testHandle(&reqURL, &reqMethod, &reqBody)

	blob := filepicker.NewBlob(FakeHandle)
	client := filepicker.NewClient(FakeApiKey)
	mock := MockServer(t, client, handler)
	defer mock.Close()

	for i, test := range tests {
		blob, err := client.Write(blob, filename, test.Opt)
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

func TestWriteError(t *testing.T) {
	fperr, handler := ErrorHandler(dummyErrStr)

	client := filepicker.NewClient(FakeApiKey)
	mock := MockServer(t, client, handler)
	defer mock.Close()

	blob := filepicker.NewBlob("XYZ")
	filename := tempFile(t)
	defer os.Remove(filename)

	switch blob, err := client.Write(blob, filename, nil); {
	case blob != nil:
		t.Errorf("want blob == nil; got %v", blob)
	case err.Error() != fperr.Error():
		t.Errorf("want error message == %q; got %q", fperr, err)
	}
}

func TestWriteErrorNoFile(t *testing.T) {
	blob := filepicker.NewBlob("XYZ")
	client := filepicker.NewClient(FakeApiKey)
	switch blob, err := client.Write(blob, "unknown.unknown.file", nil); {
	case blob != nil:
		t.Errorf("want blob == nil; got %v", blob)
	case err == nil:
		t.Error("want err != nil; got nil")
	}
}

func TestWriteURL(t *testing.T) {
	const TestURL = "https://www.filepicker.com/image.png"
	tests := []struct {
		Opt *filepicker.WriteOpts
		URL string
	}{
		{
			Opt: nil,
			URL: "http://www.filepicker.io/api/file/2HHH3",
		},
		{
			Opt: &filepicker.WriteOpts{
				Base64Decode: true,
			},
			URL: "http://www.filepicker.io/api/file/2HHH3?base64decode=true",
		},
		{
			Opt: &filepicker.WriteOpts{
				Base64Decode: true,
				Security:     dummySecurity,
			},
			URL: "http://www.filepicker.io/api/file/2HHH3?base64decode=true&policy=P&signature=S",
		},
	}

	var reqURL, reqMethod, reqBody string
	handler := testHandle(&reqURL, &reqMethod, &reqBody)

	blob := filepicker.NewBlob(FakeHandle)
	client := filepicker.NewClient(FakeApiKey)
	mock := MockServer(t, client, handler)
	defer mock.Close()

	for i, test := range tests {
		blob, err := client.WriteURL(blob, TestURL, test.Opt)
		if err != nil {
			t.Errorf("want err == nil; got %v (i:%d)", err, i)
		}
		if blob == nil {
			t.Errorf("want blob != nil; got nil (i:%d)", i)
		}
		if test.URL != reqURL {
			t.Errorf("want reqURL == %q; got %q (i:%d)", test.URL, reqURL, i)
		}
		if reqMethod != "POST" {
			t.Errorf("want reqMethod == POST; got %s (i:%d)", reqMethod, i)
		}
		if TestUrlEsc := "url=" + url.QueryEscape(TestURL); reqBody != TestUrlEsc {
			t.Errorf("want reqBody == %q; got %q (i:%d)", TestUrlEsc, reqBody, i)
		}
	}
}

func TestWriteURLError(t *testing.T) {
	fperr, handler := ErrorHandler(dummyErrStr)

	blob := filepicker.NewBlob(FakeHandle)
	client := filepicker.NewClient(FakeApiKey)
	mock := MockServer(t, client, handler)
	defer mock.Close()

	switch blob, err := client.WriteURL(blob, "http://www.address.fp", nil); {
	case blob != nil:
		t.Errorf("want blob == nil; got %v", blob)
	case err.Error() != fperr.Error():
		t.Errorf("want error message == %q; got %q", fperr, err)
	}
}
