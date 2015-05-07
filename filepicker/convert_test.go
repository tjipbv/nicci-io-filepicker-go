package filepicker_test

import (
	"testing"

	"github.com/filepicker/filepicker-go/filepicker"
)

func TestConvertAndStore(t *testing.T) {
	tests := []struct {
		Opt  *filepicker.ConvertOpts
		URL  string
		Body string
	}{
		{
			Opt: &filepicker.ConvertOpts{
				Width:    100,
				Location: filepicker.Azure,
			},
			URL:  "http://www.filepicker.io/api/file/2HHH3/convert",
			Body: "key=0KKK1&storeLocation=azure&width=100",
		},
		{
			Opt: &filepicker.ConvertOpts{
				Width:  150,
				Height: 200,
				Fit:    filepicker.FitScale,
			},
			URL:  "http://www.filepicker.io/api/file/2HHH3/convert",
			Body: "fit=scale&height=200&key=0KKK1&width=150",
		},
		{
			Opt: &filepicker.ConvertOpts{
				Align:   filepicker.AlignTop,
				Quality: 34,
			},
			URL:  "http://www.filepicker.io/api/file/2HHH3/convert",
			Body: "align=top&key=0KKK1&quality=34",
		},
		{
			Opt: &filepicker.ConvertOpts{
				Width:    100,
				Security: dummySecurity,
			},
			URL:  "http://www.filepicker.io/api/file/2HHH3/convert",
			Body: "key=0KKK1&policy=P&signature=S&width=100",
		},
	}

	var reqURL, reqMethod, reqBody string
	handler := testHandle(&reqURL, &reqMethod, &reqBody)

	blob := filepicker.NewBlob(FakeHandle)
	client := filepicker.NewClient(FakeApiKey)
	mock := MockServer(t, client, handler)
	defer mock.Close()

	for i, test := range tests {
		blob, err := client.ConvertAndStore(blob, test.Opt)
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
		if reqBody != test.Body {
			t.Errorf("want reqBody == %q; got %q (i:%d)", test.Body, reqBody, i)
		}
	}
}

func TestConvertAndStoreError(t *testing.T) {
	fperr, handler := ErrorHandler(dummyErrStr)

	blob := filepicker.NewBlob(FakeHandle)
	opts := &filepicker.ConvertOpts{Width: 300}
	client := filepicker.NewClient(FakeApiKey)
	mock := MockServer(t, client, handler)
	defer mock.Close()

	switch blob, err := client.ConvertAndStore(blob, opts); {
	case blob != nil:
		t.Errorf("want blob == nil; got %v", blob)
	case err.Error() != fperr.Error():
		t.Errorf("want error message == %q; got %q", fperr, err)
	}
}
