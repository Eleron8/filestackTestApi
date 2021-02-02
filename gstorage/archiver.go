package gstorage

import (
	"archive/zip"
	"io"
	"os"
)

func ZipFile(w io.Writer, filename string) error {
	zipWriter := zip.NewWriter(w)
	defer zipWriter.Close()
	if err := AddFileToZip(zipWriter, filename); err != nil {
		return err
	}
	return nil
}

func AddFileToZip(zipWriter *zip.Writer, filename string) error {

	fileToZip, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fileToZip.Close()

	// Get the file information
	info, err := fileToZip.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	header.Name = filename

	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, fileToZip)
	return err
}
