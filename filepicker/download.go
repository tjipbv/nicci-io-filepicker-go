package filepicker

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

// DownloadOpts TODO
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
// them to a url.Values object.
func (do *DownloadOpts) toValues() url.Values {
	return toValues(*do)
}

// DownloadTo TODO : (ppknap)
func (c *Client) DownloadTo(src *Blob, opt *DownloadOpts, dst io.Writer) (
	written int64, err error) {
	resp, err := c.download(src.Url, opt)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	return io.Copy(dst, resp.Body)
}

// DownloadToFile TODO : (ppknap)
func (c *Client) DownloadToFile(src *Blob, opt *DownloadOpts, filedir string) (err error) {
	resp, err := c.download(src.Url, opt)
	if err != nil {
		return
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
		return
	}
	defer file.Close()
	_, err = io.Copy(file, resp.Body)
	return
}

func (c *Client) download(src *Blob, opt *DownloadOpts) (resp *http.Response, err error) {
	blobUrl, err := url.Parse(src.Url)
	if err != nil {
		return
	}
	if opt != nil {
		blobUrl.RawQuery = opt.toValues().Encode()
	}
	if resp, err = c.Client.Get(blobUrl.String()); err != nil {
		return
	}
	if invalidResCode(resp.StatusCode) {
		resp.Body.Close()
		return nil, FPError(resp.StatusCode)
	}
	return
}
