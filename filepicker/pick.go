package filepicker

import (
	"net/url"
	"path"
)

// PickOpts structure allows the user to configure security options when picking a file.
type PickOpts struct {
	// Security stores Filepicker.io policy and signature members. If you enable
	// security option in your developer portal, these values must be set in
	// order to perform a valid request call.
	Security
}

// toValues takes all non-zero values from provided PickOpts entity and puts
// them to url.Values object.
func (po *PickOpts) toValues() url.Values {
	return toValues(*po)
}

// PickURL creates a symlink to the underlaying file. Thus, if the user deletes
// the file from its storage, the blob object returned from this call will be
// invalid.
func (c *Client) PickURL(dataURL string, opt *PickOpts) (*Blob, error) {
	return c.storeURL(dataURL, func() string {
		return c.toPickURL(opt).String()
	})
}

func (c *Client) toPickURL(opt *PickOpts) *url.URL {
	values := url.Values{}
	if opt != nil {
		values = opt.toValues()
	}
	values.Set("key", c.apiKey)
	return &url.URL{
		Scheme:   apiURL.Scheme,
		Host:     apiURL.Host,
		Path:     path.Join("api", "pick"),
		RawQuery: values.Encode(),
	}
}
