package handlers

import (
	"fmt"
	"log"
	"net/http"
	"server/data"
	"server/util"

	"github.com/go-chi/chi/v5"
)

func CheckStatus(w http.ResponseWriter, r *http.Request) {
	util.EnableCors(&w)
	w.WriteHeader(http.StatusOK)
}

func Upload(w http.ResponseWriter, r *http.Request) {
	util.EnableCors(&w)

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
		if !util.ValidImageType(file.Header["Content-Type"][0]) {
			w.WriteHeader(400)
			return
		}

		// save uploaded files
		img, err := util.SaveFile(&f, file.Filename, r.MultipartForm.Value["uid"][0])
		if err != nil {
			fmt.Println(err)
		}

		// apply filter to images
		err = util.ApplyFilter(img, r.MultipartForm.Value["filter"][0], r.MultipartForm.Value["uid"][0], file.Filename)
		if err != nil {
			w.WriteHeader(int(400))
			return
		}
	}

	err = util.Archive(images, r.MultipartForm.Value["uid"][0])
	if err != nil {
		log.Fatal(err)
	}

	w.WriteHeader(200)
}

func SessionClosed(w http.ResponseWriter, r *http.Request) {
	util.EnableCors(&w)
	fileName := chi.URLParam(r, "uid")
	// fmt.Println(fileName)
	if !data.InMemoryArchives.FileExist(fileName) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	data.InMemoryArchives.Remove(fileName)
	data.RemoveFromDisk(fileName)
}

func CheckFileStatus(w http.ResponseWriter, r *http.Request) {
	util.EnableCors(&w)
	uid := chi.URLParam(r, "uid")

	if data.InMemoryArchives.FileExist(uid) {
		w.WriteHeader(http.StatusOK)
		return
	}

	w.WriteHeader(http.StatusNotFound)
}

func DownloadFile(w http.ResponseWriter, r *http.Request) {
	util.EnableCors(&w)
	fileName := chi.URLParam(r, "uid")
	if !data.InMemoryArchives.FileExist(fileName) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	w.Header().Set("Content-Type", "application/octet-stream")

	archiveExt := util.ExtBasedOnPlatform()
	filePath := "archives/" + fileName + archiveExt
	http.ServeFile(w, r, filePath)
}

// func validImageType(imgType string) bool {
// 	validImgTypes := []string{"image/jpg", "image/jpeg", "image/png"}
// 	for _, validType := range validImgTypes {
// 		if validType == imgType {
// 			return true
// 		}
// 	}
// 	return false
// }

// func saveFile(file *multipart.File, filename, uid string) (*os.File, error) {
// 	dst, err := os.Create(path.Join("uploaded", uid+"_"+filename))
// 	defer (*file).Close()
// 	defer dst.Close()
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Copy the uploaded file to the created file on the filesystem
// 	if _, err := io.Copy(dst, *file); err != nil {
// 		return nil, err
// 	}
// 	return dst, nil
// }

// func applyFilter(img *os.File, filterType, uid, filename string) error {
// 	image, err := imgio.Open(img.Name())
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	switch filterType {
// 	case "gray":
// 		grayFilter(image, filename, path.Ext(img.Name()), uid)
// 		return nil
// 	default:
// 		return errors.New("Invalid filter")
// 	}
// }

// func archive(imageNames []string, uid string) error {
// 	images := make(map[string]string)
// 	for i := 0; i < len(imageNames); i++ {
// 		key := "filtered/" + uid + "_" + "Gray_" + imageNames[i]
// 		images[key] = ""
// 	}
// 	// fmt.Println(images)
// 	archive, err := archiver.FilesFromDisk(nil, images)
// 	if err != nil {
// 		return err
// 	}

// 	archiverExt := extBasedOnPlatform()

// 	out, err := os.Create(path.Join("archives", uid+archiverExt))
// 	if err != nil {
// 		return err
// 	}
// 	defer out.Close()

// 	format := archiver.CompressedArchive{
// 		Archival: archiver.Tar{},
// 	}
// 	err = format.Archive(context.Background(), out, archive)
// 	if err != nil {
// 		return err
// 	}
// 	inMemoryArchives = append(inMemoryArchives, uid)
// 	return nil
// }

// func grayFilter(myimage image.Image, imageName string, ext string, uid string) {
// 	grayImage := effect.Grayscale(myimage)
// 	filename, ext := extractFileMeta(imageName)
// 	image := fmt.Sprintf("%s_Gray_%s.%s", uid, filename, ext)

// 	file, err := os.Create(path.Join("filtered", image))
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	defer file.Close()

// 	err = jpeg.Encode(file, grayImage, nil)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }
