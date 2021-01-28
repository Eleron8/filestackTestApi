package models

type TransformData struct {
	File       string      `json:"file"`
	Transforms []Transform `json:"transforms"`
}

type Transform struct {
	Action           TransformAction `json:"action"`
	IsAdditionalInfo bool            `json:"is_additional_info"`
	AdditionalInfo   int             `json:"additional_info"`
}

type TransformAction string

const (
	Crop       TransformAction = "crop"
	Rotate     TransformAction = "rotate"
	RemoveExif TransformAction = "remove exif"
)
