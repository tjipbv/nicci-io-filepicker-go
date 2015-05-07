package filepicker

import (
	"encoding/json"
	"net/url"
	"path"
	"time"
)

// MetaTag TODO : (ppknap)
type MetaTag string

// TODO : (ppknap)
const (
	TagSize      = MetaTag("size")
	TagMimetype  = MetaTag("mimetype")
	TagFilename  = MetaTag("filename")
	TagWidth     = MetaTag("width")
	TagHeight    = MetaTag("height")
	TagUploaded  = MetaTag("uploaded")
	TagWriteable = MetaTag("writeable")
	TagMd5Hash   = MetaTag("md5")
	TagLocation  = MetaTag("location")
	TagPath      = MetaTag("path")
	TagContainer = MetaTag("container")
)

// StatOpts TODO : (ppknap)
type StatOpts struct {
	// Tags TODO : (ppknap)
	Tags []MetaTag `json:"tags,omitempty"`

	// Security stores Filepicker.io policy and signature members. If you enable
	// security option in your developer portal, these values must be set in
	// order to perform a valid request call.
	Security
}

// toValues TODO : (ppknap)
func (mo *StatOpts) toValues() url.Values {
	values := toValues(*mo)
	values.Del("tags")
	for _, tag := range mo.Tags {
		values.Add(string(tag), "true")
	}
	return values
}

// Metadata TODO : (ppkanp)
type Metadata map[string]interface{}

// Size returns the size of a stored file in bytes. The second value (ok) is set
// to false if the information is unavailable.
func (md Metadata) Size() (size uint64, ok bool) {
	if val, ok := md[string(TagSize)]; ok {
		return uint64(val.(float64)), ok
	}
	return
}

// Mimetype returns the type of a stored file. The second value (ok) is set to
// false if the information is unavailable.
func (md Metadata) Mimetype() (mimetype string, ok bool) {
	if val, ok := md[string(TagMimetype)]; ok {
		return val.(string), ok
	}
	return
}

// Filename returns the name of a stored file. The second value (ok) is set to
// false if the information is unavailable.
func (md Metadata) Filename() (filename string, ok bool) {
	if val, ok := md[string(TagFilename)]; ok {
		return val.(string), ok
	}
	return
}

// Width returns the width of a stored image. If the file is not an image or the
// information about its size is unavailable, the second value (ok) will be set
// to false.
func (md Metadata) Width() (width uint64, ok bool) {
	if val, ok := md[string(TagWidth)]; ok && val != nil {
		return uint64(val.(float64)), ok
	}
	return
}

// Height returns the height of a stored image. If the file is not an image or
// the information about its size is unavailable, the second value (ok) will be
// set to false.
func (md Metadata) Height() (height uint64, ok bool) {
	if val, ok := md[string(TagHeight)]; ok && val != nil {
		return uint64(val.(float64)), ok
	}
	return
}

// Uploaded returns the upload time of a stored file. The second value (ok) is
// set to false if the information is unavailable.
func (md Metadata) Uploaded() (uploaded time.Time, ok bool) {
	if val, ok := md[string(TagUploaded)]; ok {
		raw := int64(val.(float64))
		return time.Unix(raw/1000, raw%1000), ok
	}
	return
}

// Writeable specifies if the stored file is writeable. The second value (ok) is
// set to false if the information is unavailable.
func (md Metadata) Writeable() (writeable, ok bool) {
	if val, ok := md[string(TagWriteable)]; ok {
		return val.(bool), ok
	}
	return
}

// Md5Hash returns the MD5 hash of the stored file. The second value (ok) is set
// to false if the information is unavailable.
func (md Metadata) Md5Hash() (md5hash string, ok bool) {
	if val, ok := md[string(TagMd5Hash)]; ok {
		return val.(string), ok
	}
	return
}

// Location returns the storage location (S3, etc.) of a stored file. The second
// value (ok) is set to false if the information is unavailable.
func (md Metadata) Location() (location Storage, ok bool) {
	if val, ok := md[string(TagLocation)]; ok {
		return Storage(val.(string)), ok
	}
	return
}

// Path returns the storage path of a stored file. The second value (ok) is set
// to false if the information is unavailable.
func (md Metadata) Path() (path string, ok bool) {
	if val, ok := md[string(TagPath)]; ok {
		return val.(string), ok
	}
	return
}

// Container returns the storage container of a stored file. The second
// value (ok) is set to false if the information is unavailable.
func (md Metadata) Container() (container string, ok bool) {
	if val, ok := md[string(TagContainer)]; ok {
		return val.(string), ok
	}
	return
}

// Stat allows the user to get more detailed metadata about the stored file.
func (c *Client) Stat(src *Blob, opt *StatOpts) (Metadata, error) {
	blobURL, err := url.Parse(src.URL)
	if err != nil {
		return nil, err
	}
	if opt != nil {
		blobURL.RawQuery = opt.toValues().Encode()
	}
	blobURL.Path = path.Join(blobURL.Path, "metadata")
	resp, err := c.do("GET", blobURL.String(), "", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if err := readError(resp); err != nil {
		return nil, err
	}
	md := make(Metadata)
	if err := json.NewDecoder(resp.Body).Decode(&md); err != nil {
		return nil, err
	}
	return md, nil
}
