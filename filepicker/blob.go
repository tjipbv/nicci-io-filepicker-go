package filepicker

import (
	"net/url"
	"path"
)

// Blob contains information about the stored file.
type Blob struct {
	// Url points to where the file is stored.
	Url string `json:"url,omitempty"`

	// Filename is the name of the file, if available.
	Filename string `json:"filename,omitempty"`

	// Mimetype is the mimetype of the file, if available.
	Mimetype string `json:"type,omitempty"`

	// Size is the size of the file in bytes. When this value is not available,
	// you can get it by using Client.Stat method.
	Size uint64 `json:"size,omitempty"`

	// Key shows where in the file storage the file was put.
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
	blobUrl := url.URL{
		Scheme: apiURL.Scheme,
		Host:   apiURL.Host,
		Path:   path.Join("api", "file", handle),
	}
	return &Blob{Url: blobUrl.String()}
}

// Handle returns the unique identifier of the file. Its value is used by
// filepicker service to locate the data.
func (b *Blob) Handle() string {
	blobUrl, err := url.Parse(b.Url)
	if err != nil {
		return ""
	}
	return path.Base(blobUrl.Path)
}
