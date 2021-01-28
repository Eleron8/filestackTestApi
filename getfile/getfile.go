package getfile

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// var httpClient *http.Client

// func init() {
// 	httpClient = &http.Client{
// 		Timeout: time.Second * 10,
// 	}
// }

type FileHandler struct {
	client *http.Client
}

func NewFileHandler(client *http.Client) FileHandler {
	return FileHandler{
		client: client,
	}
}

func (f FileHandler) GetFile(fileUrl string) (int64, error) {
	resp, err := f.client.Get(fileUrl)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	filename, err := f.getfilename(fileUrl)
	if err != nil {
		return 0, err
	}
	file, err := f.createFile(filename)
	if err != nil {
		return 0, err
	}
	if len(filename) == 0 {
		return 0, fmt.Errorf("file name not found. Lenght: %d", len(filename))
	}
	size, err := io.Copy(file, resp.Body)
	if err != nil {
		return 0, err
	}
	return size, nil
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
		} else if contain := strings.Contains(v, ".jpg"); contain {
			filename = v
		} else {
			return "", fmt.Errorf("URL doesn't have image file. Must contain png or jpg")
		}

	}
	// filename := splits[len(splits)-1]
	return filename, nil
}

func (f FileHandler) createFile(filename string) (*os.File, error) {
	file, err := os.Create(filename)
	if err != nil {
		return nil, err
	}
	return file, nil
}
