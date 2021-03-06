package getfile

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"go.uber.org/zap"
)

type FileHandler struct {
	client *http.Client
	logger *zap.Logger
}

func NewFileHandler(client *http.Client, logger *zap.Logger) FileHandler {
	return FileHandler{
		client: client,
		logger: logger,
	}
}

func (f FileHandler) GetFile(fileUrl string) (int64, string, error) {
	handleErr := func(err error) (int64, string, error) {
		return 0, "", fmt.Errorf("get file error: %w", err)
	}
	resp, err := f.client.Get(fileUrl)
	if err != nil {
		handleErr(err)
	}
	defer resp.Body.Close()
	filename, err := f.getfilename(fileUrl)
	if err != nil {
		handleErr(err)
	}
	file, err := f.createFile(filename)
	if err != nil {
		handleErr(err)
	}
	if len(filename) == 0 {
		return 0, "", fmt.Errorf("file name not found. Lenght: %d", len(filename))
	}
	size, err := io.Copy(file, resp.Body)
	if err != nil {
		handleErr(err)
	}
	return size, file.Name(), nil
}

func (f FileHandler) getfilename(fileUrl string) (string, error) {
	fileurl, err := url.Parse(fileUrl)
	if err != nil {
		return "", err
	}
	path := fileurl.Path

	splits := strings.Split(path, "/")
	filename := ""
	for _, v := range splits {
		if contain := strings.Contains(v, ".png"); contain {
			filename = v
		}
		if contain := strings.Contains(v, ".jpg"); contain {
			filename = v
		}

	}
	f.logger.Info("created file", zap.String("filename", filename))
	return filename, nil
}

func (f FileHandler) createFile(filename string) (*os.File, error) {
	file, err := os.Create(filename)
	if err != nil {
		return nil, err
	}
	return file, nil
}
