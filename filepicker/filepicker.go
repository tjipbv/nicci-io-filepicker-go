package filepicker

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

func init() {
	var err error
	if apiURL, err = url.Parse(FilepickerURL); err != nil {
		panic("filepicker: invalid filepicker.io service address " + FilepickerURL)
	}
}

// FilepickerURL is a link to filepicker.io service.
const FilepickerURL = "https://www.filepicker.io/"

// apiURL is a url.URL type representation of FilepickerURL address.
var apiURL *url.URL

// Storage represents cloud storage services supported by filepicker.io client.
type Storage string

const (
	S3        = Storage("S3")        // Amazon S3 bucket.
	Azure     = Storage("azure")     // Azure blob storage container.
	Dropbox   = Storage("dropbox")   // Dropbox folder.
	Rackspace = Storage("rackspace") // Rackspace cloud files container.
)

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

// invalidResCode returns true when response code is not valid.
func invalidResCode(code int) bool {
	return code != http.StatusOK
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
