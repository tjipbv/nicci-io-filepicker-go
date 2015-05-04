package filepicker

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path"
)

// StoreOpts structure allows the user to configure how to store the data.
type StoreOpts struct {
	// Filename specifies the name of the stored file. If this variable is
	// empty, filepicker service will choose the label automatically.
	Filename string `json:"filename,omitempty"`

	// Mimetype specifies the type of the stored file.
	Mimetype string `json:"mimetype,omitempty"`

	// Location contains the name of file storage service which will be used to
	// store a file. If this field is not set, filepicker client will use Simple
	// Storage Service (S3).
	Location Storage `json:"location,omitempty"`

	// Path to store the file at within the specified file store. If the
	// provided path ends in a '/', it will be treated as a folder.
	Path string `json:"path,omitempty"`

	// Container or a bucket in the specified file store where the file should
	// end up. If this parameter is omitted, the file is stored in the default
	// container specified in the user's developer portal.
	Container string `json:"container,omitempty"`

	// Base64Decode indicates whether the data should be first decoded from
	// base64 before being written to the file.
	Base64Decode bool `json:"base64decode,omitempty"`

	// Access allows to use direct links to underlying file store service.
	Access string `json:"access,omitempty"`

	// Security stores Filepicker.io policy and signature members. If you enable
	// security option in your developer portal, these values must be set in
	// order to perform a valid request call.
	Security
}

// toValues takes all non-zero values from provided StoreOpts entity and puts
// them to a url.Values object.
func (so *StoreOpts) toValues() url.Values {
	return toValues(*so)
}

// Store opens the named file and sends it content to client's storage bucket.
// If there is no error, this function returns a blob object that contains
// information about the stored file.
//
// StoreOpt defines how filepicker.io will store the data. If a nil pointer is
// provided, this function will use default storage options.
func (c *Client) Store(name string, opt *StoreOpts) (blob *Blob, err error) {
	buff := &bytes.Buffer{}
	wr := multipart.NewWriter(buff)
	file, err := os.Open(name)
	if err != nil {
		return
	}
	defer file.Close()
	mimewr, err := wr.CreateFormFile("fileUpload", name)
	if err != nil {
		return
	}
	if _, err = io.Copy(mimewr, file); err != nil {
		return
	}
	content := wr.FormDataContentType()
	wr.Close()
	return storeRes(c.Client.Post(c.toStoreURL(opt).String(), content, buff))
}

// StoreURL takes a URL that points to the data to store and sends them directly
// to client's storage bucket. If the call succeeds, this function will return a
// blob object that contains information about the stored file.
//
// StoreOpt defines how filepicker.io will store the data. If a nil pointer is
// provided, this function will use default storage options.
func (c *Client) StoreURL(dataUrl string, opt *StoreOpts) (blob *Blob, err error) {
	values := url.Values{}
	values.Set("url", dataUrl)
	return storeRes(c.Client.PostForm(c.toStoreURL(opt).String(), values))
}

// storeRes handles client response error and, if there is none, this function
// reads response's Body and unmarshals it into Blob object.
func storeRes(resp *http.Response, respErr error) (blob *Blob, err error) {
	if respErr != nil {
		return nil, respErr
	}
	defer resp.Body.Close()
	if invalidResCode(resp.StatusCode) {
		return nil, FPError(resp.StatusCode)
	}
	blob = &Blob{}
	return blob, json.NewDecoder(resp.Body).Decode(blob)
}

func (c *Client) toStoreURL(opt *StoreOpts) *url.URL {
	storage := c.storage
	values := url.Values{}
	if opt != nil {
		values = opt.toValues()
		if opt.Location != "" {
			storage = opt.Location
		}
	}
	values.Set("key", c.apiKey)
	return &url.URL{
		Scheme:   apiURL.Scheme,
		Host:     apiURL.Host,
		Path:     path.Join("api", "store", string(storage)),
		RawQuery: values.Encode(),
	}
}
