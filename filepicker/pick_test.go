package filepicker_test

import (
	"net/url"
	"testing"

	"github.com/filepicker/filepicker-go/filepicker"
)

func TestPickURL(t *testing.T) {
	const TestURL = "https://www.filepicker.com/image.png"
	tests := []struct {
		Opt *filepicker.PickOpts
		URL string
	}{
		{
			Opt: nil,
			URL: "http://www.filepicker.io/api/pick?key=0KKK1",
		},
		{
			Opt: &filepicker.PickOpts{},
			URL: "http://www.filepicker.io/api/pick?key=0KKK1",
		},
		{
			Opt: &filepicker.PickOpts{
				Security: dummySecurity,
			},
			URL: "http://www.filepicker.io/api/pick?key=0KKK1&policy=P&signature=S",
		},
	}

	var reqURL, reqMethod, reqBody string
	handler := testHandle(&reqURL, &reqMethod, &reqBody)
	client := filepicker.NewClient(FakeApiKey)
	mock := MockServer(t, client, handler)
	defer mock.Close()

	for i, test := range tests {
		blob, err := client.PickURL(TestURL, test.Opt)
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

func TestPickURLError(t *testing.T) {
	fperr, handler := ErrorHandler(dummyErrStr)

	client := filepicker.NewClient(FakeApiKey)
	mock := MockServer(t, client, handler)
	defer mock.Close()

	switch blob, err := client.PickURL("http://www.address.fp", nil); {
	case blob != nil:
		t.Errorf("want blob == nil; got %v", blob)
	case err.Error() != fperr.Error():
		t.Errorf("want error message == %q; got %q", fperr, err)
	}
}
