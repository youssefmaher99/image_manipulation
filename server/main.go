package main

import (
	"server/router"

	"log"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

// var inMemoryArchives []string

// func checkStatus(w http.ResponseWriter, r *http.Request) {
// 	enableCors(&w)
// 	w.WriteHeader(http.StatusOK)
// }

// func sessionClosed(w http.ResponseWriter, r *http.Request) {
// 	enableCors(&w)
// 	fileName := chi.URLParam(r, "uid")
// 	// fmt.Println(fileName)
// 	if !fileExist(fileName) {
// 		w.WriteHeader(http.StatusNotFound)
// 		return
// 	}

// 	removeFromInMemoryArchives(fileName)
// 	removeFromDisk(fileName)
// }

// func removeFromDisk(fileName string) {

// 	// find matching files in uploded directory
// 	uplodedFiles, err := findMatchingFiles("uploaded", fileName)
// 	handleErr(err)

// 	// remove matching files in uploded directory
// 	err = removeMatchingFiles(uplodedFiles)
// 	handleErr(err)

// 	// find matching files in filtered directory
// 	filteredFiles, err := findMatchingFiles("filtered", fileName)
// 	handleErr(err)

// 	// remove matching files in filtered directory
// 	err = removeMatchingFiles(filteredFiles)
// 	handleErr(err)

// 	// find matching files in archives directory
// 	archivesFiles, err := findMatchingFiles("archives", fileName)
// 	handleErr(err)

// 	// remove matching files in archives directory
// 	err = removeMatchingFiles(archivesFiles)
// 	handleErr(err)

// }

// func handleErr(err error) {
// 	if err != nil {
// 		log.Println(err)
// 	}
// }

// func findMatchingFiles(dir string, fileName string) ([]string, error) {
// 	files, err := filepath.Glob(dir + "/" + fileName + "*")
// 	if err != nil {
// 		return nil, err
// 	}
// 	return files, nil
// }

// func removeMatchingFiles(files []string) error {
// 	for _, f := range files {
// 		if err := os.Remove(f); err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

// func removeFromInMemoryArchives(fileName string) {
// 	fmt.Printf("Before : %v\n", inMemoryArchives)
// 	for i := 0; i < len(inMemoryArchives); i++ {
// 		if inMemoryArchives[i] == fileName {
// 			inMemoryArchives = append(inMemoryArchives[:i], inMemoryArchives[i+1:]...)
// 			break
// 		}
// 	}
// 	fmt.Printf("After : %v\n", inMemoryArchives)

// }

// func downloadFile(w http.ResponseWriter, r *http.Request) {
// 	enableCors(&w)
// 	fileName := chi.URLParam(r, "uid")
// 	if !fileExist(fileName) {
// 		w.WriteHeader(http.StatusNotFound)
// 		return
// 	}

// 	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
// 	w.Header().Set("Content-Type", "application/octet-stream")

// 	archiveExt := extBasedOnPlatform()
// 	filePath := "archives/" + fileName + archiveExt
// 	http.ServeFile(w, r, filePath)
// }

// func extBasedOnPlatform() string {
// 	if runtime.GOOS == "linux" {
// 		return ".tar.gz"
// 	} else {
// 		return ".zip"
// 	}
// }

// // func checkFileStatus(w http.ResponseWriter, r *http.Request) {
// // 	enableCors(&w)
// // 	uid := chi.URLParam(r, "uid")
// // 	if fileExist(uid) {
// // 		w.WriteHeader(http.StatusOK)
// // 		return
// // 	}

// // 	w.WriteHeader(http.StatusNotFound)
// // }

// // func fileExistInMemory(fileName string) bool {
// // 	for i := 0; i < len(inMemoryArchives); i++ {
// // 		if inMemoryArchives[i] == fileName {
// // 			return true
// // 		}
// // 	}
// // 	return false
// // }

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

// func validImageType(imgType string) bool {
// 	validImgTypes := []string{"image/jpg", "image/jpeg", "image/png"}
// 	for _, validType := range validImgTypes {
// 		if validType == imgType {
// 			return true
// 		}
// 	}
// 	return false
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

// func extractFileMeta(fileName string) (string, string) {
// 	name := strings.Split(fileName, ".")[0]
// 	ext := strings.Split(fileName, ".")[1]
// 	return name, ext
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

func main() {
	r := router.CreateChiRouter(middleware.Logger)
	router.LoadRoutes(r)

	err := http.ListenAndServe(":5000", r)
	if err != nil {
		log.Fatal(err)
	}
}
