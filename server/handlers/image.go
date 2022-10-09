package handlers

import (
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"server/data"
	"server/util"
	"sync"

	"github.com/go-chi/chi/v5"
)

func CheckStatus(w http.ResponseWriter, r *http.Request) {
	util.EnableCors(&w)
	w.WriteHeader(http.StatusOK)
}

func Upload(w http.ResponseWriter, r *http.Request) {
	util.EnableCors(&w)

	// less than 5 MB max
	r.Body = http.MaxBytesReader(w, r.Body, 1024*1024*500)

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
	wg := sync.WaitGroup{}
	for _, file := range files {
		wg.Add(1)
		go func(file *multipart.FileHeader) {
			images = append(images, file.Filename)
			f, err := file.Open()
			if err != nil {
				log.Fatal(err)
				// break
				return
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
			wg.Done()
		}(file)
	}
	wg.Wait()
	fmt.Println(images)
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
