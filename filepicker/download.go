package filepicker

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

// DownloadOpts structure defines a set of additional options that may be
// required to successfully download the stored data.
type DownloadOpts struct {
	// Base64Decode indicates whether the data should be first decoded from
	// base64 before being written to the file.
	Base64Decode bool `json:"base64decode,omitempty"`

	// Security stores Filepicker.io policy and signature members. If you enable
	// security option in your developer portal, these values must be set in
	// order to perform a valid request call.
	Security
}

// toValues takes all non-zero values from provided DownloadOpts entity and puts
// them to url.Values object.
func (do *DownloadOpts) toValues() url.Values {
	return toValues(*do)
}

// DownloadTo TODO : (ppknap)
func (c *Client) DownloadTo(src *Blob, opt *DownloadOpts, dst io.Writer) (int64, error) {
	resp, err := c.download(src, opt)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	return io.Copy(dst, resp.Body)
}

// DownloadToFile TODO : (ppknap)
func (c *Client) DownloadToFile(src *Blob, opt *DownloadOpts, filedir string) error {
	resp, err := c.download(src, opt)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	directory, filename := filepath.Split(filedir)
	if filename == "" || filename == "." {
		if filename = resp.Header.Get("X-File-Name"); filename == "" {
			return fmt.Errorf("filepicker: invalid file name (handle %q)",
				src.Handle())
		}
	}
	file, err := os.Create(filepath.Clean(filepath.Join(directory, filename)))
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(file, resp.Body)
	return err
}

func (c *Client) download(src *Blob, opt *DownloadOpts) (resp *http.Response, err error) {
	blobURL, err := url.Parse(src.URL)
	if err != nil {
		return
	}
	if opt != nil {
		blobURL.RawQuery = opt.toValues().Encode()
	}
	if resp, err = c.do("GET", blobURL.String(), "", nil); err != nil {
		return
	}
	if err = readError(resp); err != nil {
		resp.Body.Close()
	}
	return
}
