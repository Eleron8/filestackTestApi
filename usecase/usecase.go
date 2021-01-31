package usecase

import (
	"fmt"

	"github.com/Eleron8/filestackTestApi/imageprocess"
	"github.com/Eleron8/filestackTestApi/models"
)

type FileGet interface {
	GetFile(fileUrl string) (int64, string, error)
}

type Usecase struct {
	FileGet FileGet
}

func NewUsecase(fileget FileGet) Usecase {
	return Usecase{
		FileGet: fileget,
	}
}

func (u Usecase) FileFlow(data models.TransformData) error {
	handleErr := func(err error) error {
		return fmt.Errorf("file flow file's url %s: %w", data.FileURL, err)
	}
	_, file, err := u.FileGet.GetFile(data.FileURL)
	if err != nil {
		return handleErr(err)
	}
	imgProc, err := imageprocess.OpenImage(file, data.Transforms)
	if err != nil {
		return handleErr(err)
	}
	if err := imgProc.ImageTransform("format string", file); err != nil {
		return handleErr(err)
	}

	return nil
}
