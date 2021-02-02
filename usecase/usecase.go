package usecase

import (
	"archive/zip"
	"fmt"
	"io"
	"regexp"
	"time"

	"sync"

	"github.com/Eleron8/filestackTestApi/gstorage"
	"github.com/Eleron8/filestackTestApi/imageprocess"
	"github.com/Eleron8/filestackTestApi/models"
	"go.uber.org/zap"
)

type FileGet interface {
	GetFile(fileUrl string) (int64, string, error)
}

// type Storage interface {
// 	GetWriter() io.Writer
// }
type Usecase struct {
	FileGet FileGet
	// Storage       Storage
	logger        *zap.Logger
	FolderName    string
	MaxGoroutines int
}

func NewUsecase(fileget FileGet, folderName string, maxgoroutines int, logger *zap.Logger) *Usecase {
	return &Usecase{
		FileGet: fileget,
		// Storage:       storage,
		FolderName:    folderName,
		MaxGoroutines: maxgoroutines,
		logger:        logger,
	}
}

func (u *Usecase) FileFlow(data models.TransformData, wr io.Writer) error {
	handleErr := func(err error) error {
		return fmt.Errorf("file flow file's url %s: %w", data.FileURL, err)
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
	errChan1 := make(chan error)
	errChan2 := make(chan error)

	zipWriter := zip.NewWriter(wr)
	defer zipWriter.Close()
	for _, v := range data.Transforms {
		guard <- struct{}{}
		wg.Add(1)
		go func(d models.Transform) {
			defer func() {
				wg.Done()
			}()
			name, err := imgProc.ImageTransform(d, format, filename, radius, formatPart)
			if err != nil {

				errChan1 <- err
				return
			}

			u.logger.Info("add to zip file starts", zap.String("filename", name))
			if err := gstorage.AddFileToZip(zipWriter, name); err != nil {
				u.logger.Info("add file to zip Failed", zap.String("filename", name), zap.Error(err))
				errChan1 <- err

				return
			}
			<-guard
			u.logger.Info("image transform", zap.String("filename", name))

		}(v)
		time.Sleep(8 * time.Second)

	}
	wg.Wait()

	select {
	case err := <-errChan1:
		return handleErr(err)
	case err := <-errChan2:
		return handleErr(err)
	default:
		return nil

	}

}

// absPath, _ := filepath.Abs("../createdImages/" + name)
// if err := gstorage.AddFileToZip(zipwriter, absPath); err != nil {
// 	u.logger.Info("add file to zip", zap.String("filename", name), zap.Error(err))
// 	errChan <- err
// 	return
// }
// close(errChan1)
// close(filenames)

// select {
// case err := <-errChan1:
// 	return handleErr(err)
// case err := <-errChan2:
// 	return handleErr(err)
// }
