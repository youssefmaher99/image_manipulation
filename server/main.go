package main

import (
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
)

func greet(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	w.Write([]byte("hello world"))
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

	files := r.MultipartForm.File["files"]
	for _, file := range files {
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
		err = applyFilter(img, r.MultipartForm.Value["filter"][0])
		if err != nil {
			w.WriteHeader(int(400))
			return
		}
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

func applyFilter(img *os.File, filterType string) error {
	image, err := imgio.Open(img.Name())
	if err != nil {
		log.Fatal(err)
	}

	switch filterType {
	case "gray":
		grayFilter(image, img.Name(), path.Ext(img.Name()))
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

func grayFilter(myimage image.Image, imageName string, ext string) {
	// size := myimage.Bounds().Size()
	// rect := image.Rect(0, 0, size.X, size.Y)
	// wImg := image.NewRGBA(rect)

	// for i := 0; i < width; i++ {
	// 	for j := 0; j < height; j++ {
	// 		pix := myimage.At(i, j)
	// 		orgColor := color.RGBAModel.Convert(pix).(color.RGBA)
	// 		r := float64(orgColor.R) * 0.92126
	// 		g := float64(orgColor.G) * 0.97152
	// 		b := float64(orgColor.B) * 0.90722

	// 		grey := uint8((r + g + b) / 3)
	// 		c := color.RGBA{
	// 			R: grey, G: grey, B: grey, A: orgColor.A,
	// 		}
	// 		wImg.Set(i, j, c)
	// 	}
	// }
	grayImage := effect.Grayscale(myimage)
	filename, ext := extractFileMeta(imageName)
	image := fmt.Sprintf("Gray_%s.%s", filename, ext)
	file, err := os.Create(path.Join("files", image))
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	err = jpeg.Encode(file, grayImage, nil)
}

func extractFileMeta(fileName string) (string, string) {
	nameAndExt := strings.Split(fileName, "/")[1]
	name := strings.Split(nameAndExt, ".")[0]
	ext := strings.Split(nameAndExt, ".")[1]
	return name, ext
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

	// catFile, err := os.Open("cat1.jpeg")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer catFile.Close()

	// ext := path.Ext(catFile.Name())
	// fmt.Println(ext)

	// cat, _, err := image.Decode(catFile)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// img, err := imgio.Open()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// blackAndWhiteFilter(img, imgName, imgExtension)
}
