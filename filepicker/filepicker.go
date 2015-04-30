package filepicker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path"
)

// FilepickerURL is a link to the filepicker.io service.
const FilepickerURL = "https://www.filepicker.io/"

// apiURL is a URL representation of FilepickerURL address.
var apiURL *url.URL

func init() {
	var err error
	if apiURL, err = url.Parse(FilepickerURL); err != nil {
		panic("filepicker: invalid filepicker address " + FilepickerURL)
	}
}

// Storage represents cloud storage services supported by filepicker client.
type Storage string

const (
	S3        = Storage("S3")        // Amazon Simple Storage Service.
	Azure     = Storage("azure")     // Microsoft Azure storage.
	Dropbox   = Storage("dropbox")   // Dropbox folder.
	Rackspace = Storage("rackspace") // Rackspace Cloud files.
)

// Blob TODO : (ppknap)
type Blob struct {
	Url       string  `json:"url,omitempty"`
	Filename  string  `json:"filename,omitempty"`
	Mimetype  string  `json:"type,omitempty"`
	Size      uint64  `json:"size,omitempty"`
	Key       string  `json:"key,omitempty"`
	Container Storage `json:"container,omitempty"`
	Writeable bool    `json:"isWriteable,omitempty"`
	Path      string  `json:"path,omitempty"`
}

// NewBlob TODO : (ppknap)
func NewBlob(handle string) (blob Blob) {
	return newBlob(handle, Security{})
}

// NewBlobSecurity TODO : (ppknap)
func NewBlobSecurity(handle string, security Security) (blob Blob) {
	return newBlob(handle, security)
}

// NewBlob TODO : (ppknap)
func newBlob(handle string, security Security) (blob Blob) {
	blobUrl := url.URL{
		Scheme:   apiURL.Scheme,
		Host:     apiURL.Host,
		Path:     path.Join("api", "file", handle),
		RawQuery: security.toValues().Encode(),
	}
	return Blob{Url: blobUrl.String()}
}

// StoreOpts structure allows user to configure how to store the data.
type StoreOpts struct {
	// Filename specifies the name of the stored file. If this variable is
	// empty, filepicker's server will choose the label automatically.
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
	// TODO : make type
	Access string `json:"access,omitempty"`

	// TODO : (ppknap)
	Security
}

// toValues takes all non-zero values from provided StoreOpts entity and puts
// them to a url.Values object.
func (so *StoreOpts) toValues() url.Values {
	return toValues(*so)
}

// Client TODO : (ppknap)
type Client struct {
	apiKey  string
	storage Storage
	Client  *http.Client
}

// NewClient TODO : (ppknap)
func NewClient(apiKey string) *Client {
	return newClient(apiKey, S3)
}

// NewClientStorage TODO : (ppknap)
func NewClientStorage(apiKey string, storage Storage) *Client {
	return newClient(apiKey, storage)
}

// newClient TODO : (ppknap)
func newClient(apiKey string, storage Storage) *Client {
	return &Client{
		apiKey:  apiKey,
		storage: storage,
		Client:  &http.Client{},
	}
}

// StoreURL TODO : (ppknap)
// TODO : mv url storeable(?)
func (c *Client) StoreURL(dataUrl string, opt StoreOpts) (blob Blob, err error) {
	values := url.Values{}
	values.Set("url", dataUrl)
	return storeRes(c.Client.PostForm(c.newStoreURL(&opt).String(), values))
}

// Store TODO : (ppknap)
// TODO : mv path storeable(?)
func (c *Client) Store(name string, opt StoreOpts) (blob Blob, err error) {
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
	return storeRes(c.Client.Post(c.newStoreURL(&opt).String(), content, buff))
}

// storeRes handles client response errors and if there are none, the function
// reads response's Body and unmarshals it into Blob object.
func storeRes(resp *http.Response, respErr error) (blob Blob, err error) {
	switch {
	case respErr != nil:
		return blob, err
	case invalidResCode(resp.StatusCode):
		return blob, FPError(resp.StatusCode)
	}
	defer resp.Body.Close()
	return blob, json.NewDecoder(resp.Body).Decode(&blob)
}

// invalidResCode returns true when response code is not valid.
func invalidResCode(code int) bool {
	return code != http.StatusOK
}

func (c *Client) newStoreURL(opt *StoreOpts) *url.URL {
	storage := c.storage
	if opt.Location != "" {
		storage = opt.Location
	}
	vals := opt.toValues()
	vals.Set("key", c.apiKey)
	return &url.URL{
		Scheme:   apiURL.Scheme,
		Host:     apiURL.Host,
		Path:     path.Join("api", "store", string(storage)),
		RawQuery: vals.Encode(),
	}
}

// toValues takes all non-zero values from provided interface and puts them to
// a url.Values object.
func toValues(val interface{}) url.Values {
	data, err := json.Marshal(val)
	if err != nil {
		panic("filepicker: invalid field " + err.Error())
	}
	temp := make(map[string]interface{})
	json.Unmarshal(data, &temp)
	values := url.Values{}
	for k, v := range temp {
		values.Set(k, fmt.Sprint(v))
	}
	return values
}
