package filepicker

import (
	"net/http"
	"net/url"
)

// RemoveOpts structure allows the user to set additional options when removing
// the data.
type RemoveOpts struct {
	// Security stores Filepicker.io policy and signature members. If you enable
	// security option in your developer portal, these values must be set in
	// order to perform a valid request call.
	Security
}

// toValues takes all non-zero values from provided RemoveOpts entity and puts
// them to a url.Values object.
func (ro *RemoveOpts) toValues() url.Values {
	return toValues(*ro)
}

// Remove is used to delete a file from Filepicker.io and any underlying storage.
func (c *Client) Remove(src *Blob, opt *RemoveOpts) (err error) {
	blobUrl, err := url.Parse(src.Url)
	if err != nil {
		return
	}
	values := url.Values{}
	if opt != nil {
		values = opt.toValues()
	}
	values.Set("key", c.apiKey)
	blobUrl.RawQuery = values.Encode()
	req, err := http.NewRequest("DELETE", blobUrl.String(), nil)
	if err != nil {
		return
	}
	resp, err := c.Client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if invalidResCode(resp.StatusCode) {
		return FPError(resp.StatusCode)
	}
	return
}
