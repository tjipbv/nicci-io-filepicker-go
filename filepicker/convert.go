package filepicker

import (
	"net/url"
	"path"
	"strings"
)

// FitOption specifies how to resize the image.
type FitOption string

// TODO : (ppknap)
const (
	FitClip  = FitOption("clip")
	FitCrop  = FitOption("crop")
	FitScale = FitOption("scale")
	FitMax   = FitOption("max")
)

// AlignOption defines how the image is aligned when resizing and using the
// "fit" parameter.
type AlignOption string

// TODO : (ppknap)
const (
	AlignTop    = AlignOption("top")
	AlignBottom = AlignOption("bottom")
	AlignLeft   = AlignOption("left")
	AlignRight  = AlignOption("right")
	AlignFaces  = AlignOption("faces")
)

// ConvertOpts structure allows the user to set conversion and security options.
type ConvertOpts struct {
	// Width of the inputted image, in pixels. This property is ignored when the
	// file is not an image.
	Width int `json:"width,omitempty"`

	// Height of the inputted image, in pixels. This property is ignored when
	// the file is not an image.
	Height int `json:"height,omitempty"`

	// Fit specifies how to resize the image.
	Fit FitOption `json:"fit,omitempty"`

	// Align determines how the image is aligned when resizing and using the
	// "Fit" parameter. Defaults to cropping to the center of the image.
	Align AlignOption `json:"align,omitempty"`

	// Format TODO : (ppknap)
	Format string `json:"format,omitempty"`

	// Compress property works only for jpeg and png files. It specifies whether
	// image should be compressed.
	Compress bool `json:"compress,omitempty"`

	// Quality specifies the quality of the resultant image. It is ignored when
	// the file is not of jpeg type.
	Quality int8 `json:"quality,omitempty"`

	// Filename specifies the name of the stored file. If this variable is
	// empty, filepicker service will choose the label automatically.
	Filename string `json:"filename,omitempty"`

	// Location contains the name of file storage service which will be used to
	// store a file. If this field is not set, filepicker client will use Simple
	// Storage Service (S3).
	Location Storage `json:"storeLocation,omitempty"`

	// Path to store the file at within the specified file store. If the
	// provided path ends in a '/', it will be treated as a folder.
	Path string `json:"storePath,omitempty"`

	// Container or a bucket in the specified file store where the file should
	// end up. If this parameter is omitted, the file is stored in the default
	// container specified in the user's developer portal.
	Container string `json:"storeContainer,omitempty"`

	// Access allows to use direct links to underlying file store service.
	Access string `json:"storeAccess,omitempty"`

	// Security stores Filepicker.io policy and signature members. If you enable
	// security option in your developer portal, these values must be set in
	// order to perform a valid request call.
	Security
}

// toValues takes all non-zero values from provided ConvertOpt instance and puts
// them to url.Values object.
func (co *ConvertOpts) toValues() url.Values {
	return toValues(*co)
}

// ConvertAndStore TODO : (ppknap)
func (c *Client) ConvertAndStore(src *Blob, opt *ConvertOpts) (*Blob, error) {
	const content = "application/x-www-form-urlencoded"
	blobURL, err := url.Parse(src.URL)
	if err != nil {
		return nil, err
	}
	blobURL.Path = path.Join(blobURL.Path, "convert")
	if opt == nil {
		panic("filepicker: convert options pointer cannot be set to nil")
	}
	values := opt.toValues()
	values.Set("key", c.apiKey)
	return storeRes(c.do("POST", blobURL.String(), content, strings.NewReader(values.Encode())))
}
