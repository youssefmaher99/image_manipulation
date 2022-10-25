package handlers

import (
	"fmt"
	"log"
	"net/http"
	"path"
	"server/data"
	"server/util"

	"github.com/go-chi/chi/v5"
)

var MyQueue *util.Queue

func CheckStatus(w http.ResponseWriter, r *http.Request) {
	util.EnableCors(&w)
	w.WriteHeader(http.StatusOK)
}

func Upload(w http.ResponseWriter, r *http.Request) {
	util.EnableCors(&w)

	// less than 50 MB max
	r.Body = http.MaxBytesReader(w, r.Body, 1024*1024*5000)

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
	job := util.Job{Uid: r.MultipartForm.Value["uid"][0], Filter: r.MultipartForm.Value["filter"][0]}
	for _, file := range files {

		f, err := file.Open()
		if err != nil {
			w.WriteHeader(400)
			log.Fatal(err)
			return
		}
		defer f.Close()

		// validate image type
		if !util.ValidImageType(file.Header["Content-Type"][0]) {
			w.WriteHeader(400)
			log.Fatal("Invalid data type")
			return
		}

		_, err = util.SaveFile(&f, file.Filename, r.MultipartForm.Value["uid"][0])
		if err != nil {
			log.Fatal(err)
		}

		img := util.Image{Name: file.Filename, Path: path.Join("uploaded", r.MultipartForm.Value["uid"][0]+"_"+file.Filename)}
		job.Images = append(job.Images, img)
	}
	MyQueue.Enqueue(job)
	w.WriteHeader(200)
}

func SessionClosed(w http.ResponseWriter, r *http.Request) {
	util.EnableCors(&w)
	fileName := chi.URLParam(r, "uid")
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
