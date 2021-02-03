package models

type TransformData struct {
	FileURL    string      `json:"file"`
	Transforms []Transform `json:"transforms"`
}

type Transform struct {
	Type   TransformAction `json:"type"`
	Params Param           `json:"params"`
}

type TransformAction string

const (
	Crop       TransformAction = "crop"
	Rotate     TransformAction = "rotate"
	RemoveExif TransformAction = "remove exif"
)

type Param struct {
	Degrees float64 `json:"degrees"`
	Width   int     `json:"width"`
	Height  int     `json:"height"`
}
