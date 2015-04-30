package filepicker

import "time"

// MetaTag TODO : (ppknap)
type MetaTag string

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

// MetaOpts TODO : (ppknap)
type MetaOpts struct {
	// Tags TODO : (ppknap)
	Tags []MetaTag

	// Security TODO : (ppknap)
	Security
}

// Metadata TODO : (ppkanp)
type Metadata map[string]interface{}

// Size TODO : (ppknap)
func (md Metadata) Size() (size uint64, ok bool) {
	if val, ok := md[string(TagSize)]; ok {
		return uint64(val.(float64)), ok
	}
	return
}

// Mimetype TODO : (ppknap)
func (md Metadata) Mimetype() (mimetype string, ok bool) {
	if val, ok := md[string(TagMimetype)]; ok {
		return val.(string), ok
	}
	return
}

// Filename TODO : (ppknap)
func (md Metadata) Filename() (filename string, ok bool) {
	if val, ok := md[string(TagFilename)]; ok {
		return val.(string), ok
	}
	return
}

// Width TODO : (ppknap)
func (md Metadata) Width() (width uint64, ok bool) {
	if val, ok := md[string(TagWidth)]; ok && val != nil {
		return uint64(val.(float64)), ok
	}
	return
}

// Height TODO : (ppknap)
func (md Metadata) Height() (height uint64, ok bool) {
	if val, ok := md[string(TagHeight)]; ok && val != nil {
		return uint64(val.(float64)), ok
	}
	return
}

// Uploaded TODO : (ppknap)
func (md Metadata) Uploaded() (uploaded time.Time, ok bool) {
	if val, ok := md[string(TagUploaded)]; ok {
		raw := int64(val.(float64))
		return time.Unix(uraw/1000, uraw%1000), ok
	}
	return
}

// Writeable TODO : (ppknap)
func (md Metadata) Writeable() (writeable, ok bool) {
	if val, ok := md[string(TagWriteable)]; ok {
		return val.(bool), ok
	}
	return
}

// Md5Hash TODO : (ppknap)
func (md Metadata) Md5Hash() (md5hash string, ok bool) {
	if val, ok := md[string(TagMd5Hash)]; ok {
		return val.(string), ok
	}
	return
}

// Location TODO : (ppknap)
func (md Metadata) Location() (location Storage, ok bool) {
	if val, ok := md[string(TagLocation)]; ok {
		return Storage(val.(string)), ok
	}
	return
}

// Path TODO : (ppknap)
func (md Metadata) Path() (path string, ok bool) {
	if val, ok := md[string(TagPath)]; ok {
		return val.(string), ok
	}
	return
}

// Container TODO : (ppknap)
func (md Metadata) Container() (container string, ok bool) {
	if val, ok := md[string(TagContainer)]; ok {
		return val.(string), ok
	}
	return
}
