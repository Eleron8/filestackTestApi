package usecase

import (
	"archive/zip"
	"fmt"
	"io"
	"regexp"

	"sync"

	"github.com/Eleron8/filestackTestApi/gstorage"
	"github.com/Eleron8/filestackTestApi/imageprocess"
	"github.com/Eleron8/filestackTestApi/models"
	"go.uber.org/zap"
)

type FileGet interface {
	GetFile(fileUrl string) (int64, string, error)
}

type Storage interface {
	GetWriter() io.Writer
}
type Usecase struct {
	FileGet       FileGet
	Storage       Storage
	logger        *zap.Logger
	FolderName    string
	MaxGoroutines int
}

func NewUsecase(fileget FileGet, storage Storage, folderName string, maxgoroutines int, logger *zap.Logger) *Usecase {
	return &Usecase{
		FileGet:       fileget,
		Storage:       storage,
		FolderName:    folderName,
		MaxGoroutines: maxgoroutines,
		logger:        logger,
	}
}

func (u *Usecase) FileFlow(data models.TransformData) (io.Reader, error) {
	handleErr := func(err error) (io.Reader, error) {
		return nil, fmt.Errorf("file flow file's url %s: %w", data.FileURL, err)
	}
	var wg = sync.WaitGroup{}
	_, file, err := u.FileGet.GetFile(data.FileURL)
	if err != nil {
		u.logger.Info("get file failed", zap.Error(err))
		return handleErr(err)
	}
	imgProc, format, err := imageprocess.OpenImage(file, data.Transforms, u.logger)
	if err != nil {
		return handleErr(err)
	}
	_, radius := imgProc.GetPixels()
	formatPart := "." + format
	reg := regexp.MustCompile(formatPart)
	names := reg.Split(file, 2)
	filename := fmt.Sprintf("%s/%s", u.FolderName, names[0])
	guard := make(chan struct{}, u.MaxGoroutines)
	rc, wc := io.Pipe()
	zipWriter := zip.NewWriter(wc)
	defer zipWriter.Close()
	for _, v := range data.Transforms {
		guard <- struct{}{}
		wg.Add(1)
		go func(d models.Transform, zipwriter *zip.Writer) (io.Reader, error) {
			name, err := imgProc.ImageTransform(d, format, filename, radius, formatPart)
			if err != nil {
				return handleErr(err)
			}
			if err := gstorage.AddFileToZip(zipwriter, "createdImages/"+name); err != nil {
				return handleErr(err)
			}
			<-guard
			wg.Done()
			return nil, nil
		}(v, zipWriter)
	}

	// for _, v := range data.Transforms {
	// 	if err := imgProc.ImageTransform(v, format, filename, radius, formatPart); err != nil {
	// 		return handleErr(err)
	// 	}
	// }

	// if err := imgProc.ImageTransforms(format, filename); err != nil {
	// 	return handleErr(err)
	// }

	return rc, nil
}
