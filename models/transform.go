package models

type TransformData struct {
	FileURL    string      `json:"file"`
	Transforms []Transform `json:"transforms"`
}

type Transform struct {
	Type   TransformAction `json:"type"`
	Params Param           `json:"params"`
	// IsAdditionalInfo bool            `json:"is_additional_info"`
	// AdditionalInfo   int             `json:"additional_info"`
}

type TransformAction string

const (
	Crop       TransformAction = "crop"
	Rotate     TransformAction = "rotate"
	RemoveExif TransformAction = "remove exif"
)

type Param struct {
	Degrees float64
	Width   int
	Height  int
}
