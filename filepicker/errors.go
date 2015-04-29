package filepicker

/// INFO : (ppknap) this file will be rewritten.

type FPError int

func (e FPError) Error() string {
	return "TODO"
}

const (
	ErrBadParameters             = 400
	ErrInvalidRequest            = 403
	ErrDOMFilesUnsupported       = 111 // Store
	ErrFileNotFound              = 115 // Store
	ErrGeneralRead               = 118 // Store
	ErrWriteBlobUnreachable      = 121
	ErrWriteURLUnreachable       = 122
	ErrStoreFileUnreachable      = 151 // Store URL/Store
	ErrStoreURLUnreachable       = 152 // Store URL
	ErrStatFileCannotBeFound     = 161 // Stat
	ErrStatCannotFetchMetadata   = 162 // Stat
	ErrRmFileCannotBeFound       = 171 // remove
	ErrRmContentStoreUnreachable = 172 // remove

)

var errors = map[int]string{
	ErrBadParameters:             "filepicker: bad parameters were passed to the server",
	ErrInvalidRequest:            "filepicker: invalid request",
	ErrDOMFilesUnsupported:       "filepicker: DOM file objects are not supported",
	ErrFileNotFound:              "filepicker: file not found",
	ErrGeneralRead:               "filepicker: general read error",
	ErrWriteBlobUnreachable:      "filepicker: the Blob to write to could not be found",
	ErrWriteURLUnreachable:       "filepicker: the remote URL is unreachable",
	ErrStoreFileUnreachable:      "filepicker: the file store is unreachable",
	ErrStoreURLUnreachable:       "filepicker: the remote URL is unreachable",
	ErrStatFileCannotBeFound:     "filepicker: the file cannot be found",
	ErrStatCannotFetchMetadata:   "filepicker: cannot fetch file metadata",
	ErrRmFileCannotBeFound:       "filepicker: cannot find the requested file",
	ErrRmContentStoreUnreachable: "filepicker: the underlying content store is unreachable",
}
