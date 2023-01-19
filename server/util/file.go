package util

import (
	"io"
	"os"
	"path"
	"strings"
)

type File interface {
	io.Reader
	io.ReaderAt
	io.Seeker
	io.Closer
}

func SaveFile(file File, filename, uid string) (*os.File, error) {
	dst, err := os.Create(path.Join("uploaded", uid+"_"+filename))
	defer file.Close()
	defer dst.Close()
	if err != nil {
		return nil, err
	}

	// Copy the uploaded file to the created file on the filesystem
	if _, err := io.Copy(dst, file); err != nil {
		return nil, err
	}
	return dst, nil
}

func ExtractFileMeta(fileName string) (string, string) {
	name := strings.Split(fileName, ".")[0]
	ext := strings.Split(fileName, ".")[1]
	return name, ext
}
