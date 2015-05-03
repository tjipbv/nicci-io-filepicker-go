package filepicker

type FitOption string

const (
	FitClip  = FitOption("clip")
	FitCrop  = FitOption("crop")
	FitScale = FitOption("scale")
	FitMax   = FitOption("max")
)

type AlignOption string

const (
	AlignTop    = AlignOption("top")
	AlignBottom = AlignOption("bottom")
	AlignLeft   = AlignOption("left")
	AlignRight  = AlignOption("right")
	AlignFaces  = AlignOption("faces")
)

type FilterOption string

const (
	FilterBlur    = FilterOption("blur")
	FilterSharpen = FilterOption("sharpen")
)

type ConvertOpt struct {
	Width    int
	Height   int
	Fit      FitOption
	Align    AlignOption
	Crop     []int
	Format   string
	Filter   FilterOption
	Compress bool
	Quality  int
	Rotate   int

	Security
}

// toValues takes all non-zero values from provided ConvertOpt instance and puts
// them to a url.Values object.
func (co *ConvertOpt) toValues() url.Values {
	return toValues(*co)
}
