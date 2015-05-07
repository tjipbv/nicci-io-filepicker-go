package filepicker

import "net/url"

// RemoveOpts structure allows the user to set additional options when removing
// the data.
type RemoveOpts struct {
	// Security stores Filepicker.io policy and signature members. If you enable
	// security option in your developer portal, these values must be set in
	// order to perform a valid request call.
	Security
}

// toValues takes all non-zero values from provided RemoveOpts entity and puts
// them to url.Values object.
func (ro *RemoveOpts) toValues() url.Values {
	return toValues(*ro)
}

// Remove is used to delete a file from Filepicker.io and any underlying storage.
func (c *Client) Remove(src *Blob, opt *RemoveOpts) error {
	blobURL, err := url.Parse(src.URL)
	if err != nil {
		return err
	}
	values := url.Values{}
	if opt != nil {
		values = opt.toValues()
	}
	values.Set("key", c.apiKey)
	blobURL.RawQuery = values.Encode()
	resp, err := c.do("DELETE", blobURL.String(), "", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return readError(resp)
}
