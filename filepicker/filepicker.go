package filepicker

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"
)

func init() {
	var err error
	if apiURL, err = url.Parse(FilepickerURL); err != nil {
		panic("filepicker: invalid filepicker.io service address " + FilepickerURL)
	}
}

// UserAgentID TODO : (ppknap)
const UserAgentID = "filepicker-go 0.1"

// FilepickerURL is a link to filepicker.io service.
const FilepickerURL = "https://www.filepicker.io/"

// apiURL is a url.URL type representation of FilepickerURL address.
var apiURL *url.URL

// Storage represents cloud storage services supported by filepicker.io client.
type Storage string

// TODO : (ppknap)
const (
	S3        = Storage("S3")        // Amazon S3 bucket.
	Azure     = Storage("azure")     // Azure blob storage container.
	Dropbox   = Storage("dropbox")   // Dropbox folder.
	Rackspace = Storage("rackspace") // Rackspace cloud files container.
)

// Blob contains information about the stored file.
type Blob struct {
	// URL points to where the file is stored.
	URL string `json:"url,omitempty"`

	// Filename is the name of the file, if available.
	Filename string `json:"filename,omitempty"`

	// Mimetype is the mimetype of the file, if available.
	Mimetype string `json:"type,omitempty"`

	// Size is the size of the file in bytes. When this value is not available,
	// the user can get it by using Client.Stat method.
	Size uint64 `json:"size,omitempty"`

	// Key shows where in the file storage the data was put.
	Key string `json:"key,omitempty"`

	// Container points to the storage in which the file was put.
	Container Storage `json:"container,omitempty"`

	// Writeable specifies whether the underlying file is writeable.
	Writeable bool `json:"isWriteable,omitempty"`

	// Path indicates Blob's position in the hierarchy of files uploaded when
	// {folders:true} is set.
	Path string `json:"path,omitempty"`
}

// NewBlob creates a new Blob object from a given file handle.
func NewBlob(handle string) *Blob {
	blobURL := url.URL{
		Scheme: apiURL.Scheme,
		Host:   apiURL.Host,
		Path:   path.Join("api", "file", handle),
	}
	return &Blob{URL: blobURL.String()}
}

// Handle returns the unique identifier of the file. Its value is used by
// filepicker service to locate the data.
func (b *Blob) Handle() string {
	blobURL, err := url.Parse(b.URL)
	if err != nil {
		return ""
	}
	return path.Base(blobURL.Path)
}

// fperror represents an error that can be returned from filepicker.io service.
type fperror struct {
	Code    int
	Message string
}

// Error satisfies builtin.error interface. It prints an error string with
// the reason of failure.
func (e fperror) Error() string {
	return fmt.Sprintf("filepicker: %d - %s", e.Code, e.Message)
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

func (c *Client) do(method, urlStr, bodyType string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, urlStr, body)
	if err != nil {
		return nil, err
	}
	if bodyType != "" {
		req.Header.Set("Content-Type", bodyType)
	}
	req.Header.Set("User-Agent", UserAgentID)
	return c.Client.Do(req)
}

// toValues takes all non-zero values from provided interface and puts them to
// url.Values object.
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

func readError(resp *http.Response) error {
	if resp.StatusCode == http.StatusOK {
		return nil
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return fperror{
		Code:    resp.StatusCode,
		Message: strings.TrimSpace(string(bytes)),
	}
}
