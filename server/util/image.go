package util

import (
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"path"
	"server/logger"

	"github.com/anthonynsimon/bild/effect"
	"github.com/anthonynsimon/bild/imgio"
)

type Image struct {
	Name string
	Path string
}

type Job struct {
	Filter string
	Images []Image
	Uid    string
}

func ValidImageType(imgType string) bool {
	validImgTypes := []string{"image/jpg", "image/jpeg", "image/png"}
	for _, validType := range validImgTypes {
		if validType == imgType {
			return true
		}
	}
	return false
}

func ApplyFilter(img *os.File, filterType, uid, filename string) error {
	image, err := imgio.Open(img.Name())
	if err != nil {
		logger.MyLog.Fatal(err)
	}
	switch filterType {
	case "gray":
		grayFilter(image, filename, uid)
		return nil
	default:
		return errors.New("invalid filter")
	}
}

func grayFilter(myimage image.Image, imageName string, uid string) {
	grayImage := effect.Grayscale(myimage)
	filename, ext := ExtractFileMeta(imageName)
	image := fmt.Sprintf("%s_Gray_%s.%s", uid, filename, ext)

	file, err := os.Create(path.Join("filtered", image))
	if err != nil {
		logger.MyLog.Fatal(err)
	}

	defer file.Close()

	err = jpeg.Encode(file, grayImage, nil)
	if err != nil {
		logger.MyLog.Fatal(err)
	}
}
