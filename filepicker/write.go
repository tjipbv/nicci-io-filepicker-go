package filepicker

import "net/url"

// WriteOpts structure defines a set of additional options that may be required
// to successfully rewrite the contents of the stored file.
type WriteOpts struct {
	// Base64Decode indicates whether the data should be first decoded from
	// base64 before being written to the file.
	Base64Decode bool `json:"base64decode,omitempty"`

	// Security stores Filepicker.io policy and signature members. If you enable
	// security option in your developer portal, these values must be set in
	// order to perform a valid request call.
	Security
}

// toValues takes all non-zero values from provided WriteOpts entity and puts
// them to url.Values object.
func (wo *WriteOpts) toValues() url.Values {
	return toValues(*wo)
}

// Write TODO : (ppknap)
func (c *Client) Write(src *Blob, name string, opt *WriteOpts) (*Blob, error) {
	return c.store(name, func() string {
		return c.toWriteURL(src, opt).String()
	})
}

// WriteURL TODO : (ppknap)
func (c *Client) WriteURL(src *Blob, dataURL string, opt *WriteOpts) (*Blob, error) {
	return c.storeURL(dataURL, func() string {
		return c.toWriteURL(src, opt).String()
	})
}

func (c *Client) toWriteURL(src *Blob, opt *WriteOpts) *url.URL {
	blobURL, err := url.Parse(src.URL)
	if err != nil {
		return &url.URL{}
	}
	values := url.Values{}
	if opt != nil {
		values = opt.toValues()
	}
	blobURL.RawQuery = values.Encode()
	return blobURL
}
