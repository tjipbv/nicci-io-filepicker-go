package filepicker_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/filepicker/filepicker-go/filepicker"
)

func TestStat(t *testing.T) {
	tests := []struct {
		Opt *filepicker.StatOpts
		URL string
	}{
		{
			Opt: nil,
			URL: "http://www.filepicker.io/api/file/2HHH3/metadata",
		},
		{
			Opt: &filepicker.StatOpts{
				Tags:     nil,
				Security: dummySecurity,
			},
			URL: "http://www.filepicker.io/api/file/2HHH3/metadata?policy=P&signature=S",
		},
		{
			Opt: &filepicker.StatOpts{
				Tags: []filepicker.MetaTag{
					filepicker.TagSize,
					filepicker.TagWidth,
				},
			},
			URL: "http://www.filepicker.io/api/file/2HHH3/metadata?size=true&width=true",
		},
		{
			Opt: &filepicker.StatOpts{
				Tags: []filepicker.MetaTag{
					filepicker.TagLocation,
				},
				Security: dummySecurity,
			},
			URL: "http://www.filepicker.io/api/file/2HHH3/metadata?location=true&policy=P&signature=S",
		},
	}

	var reqURL, reqMethod, reqBody string
	handler := func(w http.ResponseWriter, req *http.Request) {
		body, _ := ioutil.ReadAll(req.Body)
		reqBody = string(body)
		reqURL = req.URL.String()
		reqMethod = req.Method
		w.Write([]byte("{}"))
	}

	blob := filepicker.NewBlob(FakeHandle)
	client := filepicker.NewClient(FakeApiKey)
	mock := MockServer(t, client, handler)
	defer mock.Close()

	for i, test := range tests {
		meta, err := client.Stat(blob, test.Opt)
		if err != nil {
			t.Errorf("want err == nil; got %v (i:%d)", err, i)
		}
		if meta == nil {
			t.Errorf("want meta != nil; got nil (i:%d)", i)
		}
		if test.URL != reqURL {
			t.Errorf("want reqURL == %q; got %q (i:%d)", test.URL, reqURL, i)
		}
		if reqMethod != "GET" {
			t.Errorf("want reqMethod == GET; got %s (i:%d)", reqMethod, i)
		}
		if reqBody != "" {
			t.Errorf("want reqBody == ``; got %q (i:%d)", reqBody, i)
		}
	}
}

func TestStatMetadata(t *testing.T) {
	tests := []struct {
		Res   map[filepicker.MetaTag]interface{}
		Call  interface{}
		Value interface{}
		Ok    bool
	}{
		{
			Res:   map[filepicker.MetaTag]interface{}{filepicker.TagSize: 100},
			Call:  filepicker.Metadata.Size,
			Value: uint64(100),
			Ok:    true,
		},
		{
			Res:   map[filepicker.MetaTag]interface{}{},
			Call:  filepicker.Metadata.Size,
			Value: uint64(0),
			Ok:    false,
		},
		{
			Res:   map[filepicker.MetaTag]interface{}{filepicker.TagMimetype: "text/plain"},
			Call:  filepicker.Metadata.Mimetype,
			Value: "text/plain",
			Ok:    true,
		},
		{
			Res:   map[filepicker.MetaTag]interface{}{},
			Call:  filepicker.Metadata.Mimetype,
			Value: "",
			Ok:    false,
		},
		{
			Res:   map[filepicker.MetaTag]interface{}{filepicker.TagFilename: "example.txt"},
			Call:  filepicker.Metadata.Filename,
			Value: "example.txt",
			Ok:    true,
		},
		{
			Res:   map[filepicker.MetaTag]interface{}{},
			Call:  filepicker.Metadata.Filename,
			Value: "",
			Ok:    false,
		},
		{
			Res:   map[filepicker.MetaTag]interface{}{filepicker.TagWidth: 800},
			Call:  filepicker.Metadata.Width,
			Value: uint64(800),
			Ok:    true,
		},
		{
			Res:   map[filepicker.MetaTag]interface{}{},
			Call:  filepicker.Metadata.Width,
			Value: uint64(0),
			Ok:    false,
		},
		{
			Res:   map[filepicker.MetaTag]interface{}{filepicker.TagHeight: 600},
			Call:  filepicker.Metadata.Height,
			Value: uint64(600),
			Ok:    true,
		},
		{
			Res:   map[filepicker.MetaTag]interface{}{},
			Call:  filepicker.Metadata.Height,
			Value: uint64(0),
			Ok:    false,
		},
		{
			Res:   map[filepicker.MetaTag]interface{}{filepicker.TagUploaded: 1257894000000.0},
			Call:  filepicker.Metadata.Uploaded,
			Value: time.Unix(1257894000, 0),
			Ok:    true,
		},
		{
			Res:   map[filepicker.MetaTag]interface{}{},
			Call:  filepicker.Metadata.Uploaded,
			Value: time.Time{},
			Ok:    false,
		},
		{
			Res:   map[filepicker.MetaTag]interface{}{filepicker.TagWriteable: true},
			Call:  filepicker.Metadata.Writeable,
			Value: true,
			Ok:    true,
		},
		{
			Res:   map[filepicker.MetaTag]interface{}{},
			Call:  filepicker.Metadata.Writeable,
			Value: false,
			Ok:    false,
		},
		{
			Res:   map[filepicker.MetaTag]interface{}{filepicker.TagMd5Hash: "f31dbf9b885e315d98e136f1db0daf52"},
			Call:  filepicker.Metadata.Md5Hash,
			Value: "f31dbf9b885e315d98e136f1db0daf52",
			Ok:    true,
		},
		{
			Res:   map[filepicker.MetaTag]interface{}{},
			Call:  filepicker.Metadata.Md5Hash,
			Value: "",
			Ok:    false,
		},
		{
			Res:   map[filepicker.MetaTag]interface{}{filepicker.TagLocation: "S3"},
			Call:  filepicker.Metadata.Location,
			Value: filepicker.S3,
			Ok:    true,
		},
		{
			Res:   map[filepicker.MetaTag]interface{}{},
			Call:  filepicker.Metadata.Location,
			Value: filepicker.Storage(""),
			Ok:    false,
		},
		{
			Res:   map[filepicker.MetaTag]interface{}{filepicker.TagPath: "example.txt"},
			Call:  filepicker.Metadata.Path,
			Value: "example.txt",
			Ok:    true,
		},
		{
			Res:   map[filepicker.MetaTag]interface{}{},
			Call:  filepicker.Metadata.Path,
			Value: "",
			Ok:    false,
		},
		{
			Res:   map[filepicker.MetaTag]interface{}{filepicker.TagContainer: "container"},
			Call:  filepicker.Metadata.Container,
			Value: "container",
			Ok:    true,
		},
		{
			Res:   map[filepicker.MetaTag]interface{}{},
			Call:  filepicker.Metadata.Container,
			Value: "",
			Ok:    false,
		},
	}

	for i, test := range tests {
		b, err := json.Marshal(test.Res)
		if err != nil {
			t.Errorf("want err == nil; got %v (i:%d)", err, i)
		}
		metaraw := make(map[string]interface{})
		if err := json.Unmarshal(b, &metaraw); err != nil {
			t.Errorf("want err == nil; got %v (i:%d)", err, i)
		}
		metaval := reflect.ValueOf(filepicker.Metadata(metaraw))
		res := reflect.ValueOf(test.Call).Call([]reflect.Value{metaval})
		if l := len(res); l != 2 {
			t.Errorf("want len(res) == 2; got %d (i:%d)", l, i)
		}
		if !reflect.DeepEqual(test.Value, res[0].Interface()) {
			t.Errorf("values `%v` and `%v` are not equal (i:%d)", test.Value, res[0].Interface(), i)
		}
		if !reflect.DeepEqual(test.Ok, res[1].Interface()) {
			t.Errorf("boolean values `%v` and `%v` are not equal (i:%d)", test.Ok, res[1].Interface(), i)
		}
	}
}

func TestStatError(t *testing.T) {
	fperr, handler := ErrorHandler(dummyErrStr)

	blob := filepicker.NewBlob(FakeHandle)
	client := filepicker.NewClient(FakeApiKey)
	mock := MockServer(t, client, handler)
	defer mock.Close()

	switch meta, err := client.Stat(blob, nil); {
	case meta != nil:
		t.Errorf("want meta == nil; got %v", meta)
	case err.Error() != fperr.Error():
		t.Errorf("want error message == %q; got %q", fperr, err)
	}
}
