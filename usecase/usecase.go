package usecase

import (
	"archive/zip"
	"fmt"
	"io"
	"regexp"

	"github.com/Eleron8/filestackTestApi/gstorage"
	"github.com/Eleron8/filestackTestApi/imageprocess"
	"github.com/Eleron8/filestackTestApi/models"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
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
	g := new(errgroup.Group)
	zipWriter := zip.NewWriter(wr)
	defer zipWriter.Close()
	filenames := make([]string, len(data.Transforms))
	for i := 0; i < len(data.Transforms); i++ {
		guard <- struct{}{}
		n := i
		g.Go(func() error {

			name, err := imgProc.ImageTransform(data.Transforms[n], format, filename, radius, formatPart)
			if err != nil {
				u.logger.Info("image transform failed", zap.String("filename", name), zap.String("action", string(data.Transforms[n].Type)))
				return err

			}
			filenames[n] = name

			<-guard
			u.logger.Info("image transform", zap.String("filename", name))
			return err
		})

	}
	if err := g.Wait(); err != nil {
		return handleErr(err)
	}
	for _, n := range filenames {
		if err := gstorage.AddFileToZip(zipWriter, n); err != nil {
			u.logger.Info("add file to zip Failed", zap.String("filename", n), zap.Error(err))
			return err
		}
	}
	if err := gstorage.RemoveContents(u.FolderName); err != nil {
		return handleErr(err)
	}
	return nil

}
