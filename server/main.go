package main

import (
	"context"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/jpeg"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/anthonynsimon/bild/effect"
	"github.com/anthonynsimon/bild/imgio"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/mholt/archiver/v4"
)

func greet(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	w.Write([]byte("hello world"))
}

func downloadFile(w http.ResponseWriter, r *http.Request) {

}

func upload(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	// less than 5 MB max
	r.Body = http.MaxBytesReader(w, r.Body, 1024*1024*5)

	err := r.ParseMultipartForm(5000)
	if err != nil {
		w.WriteHeader(400)
		fmt.Println(err)
		return
	}

	// check if no filter option in request body
	if len(r.MultipartForm.Value["filter"]) <= 0 {
		w.WriteHeader(400)
		return
	}

	// check if no sessionId in request body
	if len(r.MultipartForm.Value["uid"]) <= 0 {
		w.WriteHeader(400)
		return
	}

	files := r.MultipartForm.File["files"]
	var images []string
	for _, file := range files {
		images = append(images, file.Filename)
		f, err := file.Open()
		if err != nil {
			log.Fatal(err)
			break
		}
		defer f.Close()

		// validate image type
		if !validImageType(file.Header["Content-Type"][0]) {
			w.WriteHeader(400)
			return
		}

		// save uploaded files
		img, err := saveFile(&f, file.Filename)
		if err != nil {
			fmt.Println(err)
		}

		// apply filter to images
		err = applyFilter(img, r.MultipartForm.Value["filter"][0], r.MultipartForm.Value["uid"][0])
		if err != nil {
			w.WriteHeader(int(400))
			return
		}
	}

	err = archive(images, r.MultipartForm.Value["uid"][0])
	if err != nil {
		log.Fatal(err)
	}

	w.WriteHeader(200)
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func saveFile(file *multipart.File, filename string) (*os.File, error) {
	dst, err := os.Create(path.Join("files", filename))
	defer (*file).Close()
	defer dst.Close()
	if err != nil {
		return nil, err
	}

	// Copy the uploaded file to the created file on the filesystem
	if _, err := io.Copy(dst, *file); err != nil {
		return nil, err
	}
	return dst, nil
}

func applyFilter(img *os.File, filterType string, uid string) error {
	image, err := imgio.Open(img.Name())
	if err != nil {
		log.Fatal(err)
	}

	switch filterType {
	case "gray":
		grayFilter(image, img.Name(), path.Ext(img.Name()), uid)
		return nil
	default:
		return errors.New("Invalid filter")
	}
}

func validImageType(imgType string) bool {
	validImgTypes := []string{"image/jpg", "image/jpeg", "image/png"}
	for _, validType := range validImgTypes {
		if validType == imgType {
			return true
		}
	}
	return false
}

func grayFilter(myimage image.Image, imageName string, ext string, uid string) {
	grayImage := effect.Grayscale(myimage)
	filename, ext := extractFileMeta(imageName)
	image := fmt.Sprintf("%s_Gray_%s.%s", uid, filename, ext)
	// Check if directory does not exist
	if _, err := os.Stat("files/cat"); os.IsNotExist(err) {
		fmt.Println(os.IsNotExist(err))
		// Create directory
		if err := os.Mkdir("cat", os.ModePerm); err != nil {
			log.Fatal(err)
		}
	}

	file, err := os.Create(path.Join("files", "cat", image))
	if err != nil {
		fmt.Println(err)
	}

	defer file.Close()

	err = jpeg.Encode(file, grayImage, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func extractFileMeta(fileName string) (string, string) {
	nameAndExt := strings.Split(fileName, "/")[1]
	name := strings.Split(nameAndExt, ".")[0]
	ext := strings.Split(nameAndExt, ".")[1]
	return name, ext
}

func archive(imageNames []string, uid string) error {
	images := make(map[string]string)
	for i := 0; i < len(imageNames); i++ {
		key := "files/cat/" + uid + "_" + "Gray_" + imageNames[i]
		fmt.Println(key)
		images[key] = ""
	}
	archive, err := archiver.FilesFromDisk(nil, images)
	if err != nil {
		return err
	}

	out, err := os.Create(strings.Join([]string{"Gray", ".zip"}, ""))
	if err != nil {
		return err
	}
	defer out.Close()

	format := archiver.CompressedArchive{
		Archival: archiver.Zip{},
	}
	err = format.Archive(context.Background(), out, archive)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/test", greet)
	r.Post("/upload", upload)

	err := http.ListenAndServe(":5000", r)
	if err != nil {
		log.Fatal(err)
	}

}
