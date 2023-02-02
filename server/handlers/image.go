package handlers

import (
	"fmt"
	"net/http"
	"path"
	"server/data"
	"server/logger"
	"server/models"
	"server/notification"
	"server/presist"
	"server/queue"
	"server/util"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

var MyQueue *queue.Queue[models.Job]

func CheckStatus(w http.ResponseWriter, r *http.Request) {
	util.EnableCors(&w)
	w.WriteHeader(http.StatusOK)
}

func Subscribe(w http.ResponseWriter, r *http.Request) {
	util.EnableCors(&w)
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// validate uuid
	_, err := uuid.Parse(chi.URLParam(r, "uid"))
	if err != nil {
		w.WriteHeader(400)
		return
	}

	uuid := chi.URLParam(r, "uid")

	// validate if this uuid is not registered
	if !data.InMemoryUUID.ItemExist(uuid) {
		w.WriteHeader(400)
		return
	}

	// create channel
	notifyChan := make(chan struct{})

	// add channel to chansStore
	notification.NotificationChans.Add(uuid, notifyChan)

	// closing mechanism
	// TODO : to be improved that later
	defer func() {
		notification.NotificationChans.Remove(uuid)
	}()

	// loop (presist connection) and wait for the channel to receive confirmation that file is now ready to be downloaded
	flusher, ok := w.(http.Flusher)
	if !ok {
		logger.MyLog.Fatal("Could not init http flusher")
	}

	for {
		select {
		case <-notification.NotificationChans[uuid]:
			fmt.Fprintf(w, "data: %s\n\n", "1")
			flusher.Flush()
		case <-r.Context().Done():
			fmt.Printf("Client %s disconnected\n", uuid)
			return
		}
	}
}

func Upload(w http.ResponseWriter, r *http.Request) {
	util.EnableCors(&w)

	// less than 50 MB max
	r.Body = http.MaxBytesReader(w, r.Body, 1024*1024*500)

	err := r.ParseMultipartForm(5000)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("File is too large"))
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

	form_uuid := r.MultipartForm.Value["uid"][0]

	// validate uuid
	_, err = uuid.Parse(form_uuid)
	if err != nil {
		w.WriteHeader(400)
		return
	}

	// validate if this uuid is already registered
	if data.InMemoryUUID.ItemExist(form_uuid) {
		w.WriteHeader(400)
		return
	}

	files := r.MultipartForm.File["files"]
	job := models.Job{Uid: form_uuid, Filter: r.MultipartForm.Value["filter"][0]}
	for _, file := range files {

		f, err := file.Open()
		if err != nil {
			w.WriteHeader(400)
			logger.MyLog.Fatal(err)
			return
		}
		defer f.Close()

		// validate image type
		if !util.ValidImageType(file.Header["Content-Type"][0]) {
			w.WriteHeader(400)
			logger.MyLog.Fatal("Invalid data type")
			return
		}

		_, err = util.SaveFile(f, file.Filename, form_uuid)
		if err != nil {
			logger.MyLog.Fatal(err)
		}

		img := models.Image{Name: file.Filename, Path: path.Join("uploaded", form_uuid+"_"+file.Filename)}
		job.Images = append(job.Images, img)
	}

	data.InMemoryUUID.Add(form_uuid, struct{}{})
	presist.AddUUID(form_uuid)

	MyQueue.Enqueue(job)
	presist.AddJob(job)

	w.WriteHeader(200)
}

func SessionClosed(w http.ResponseWriter, r *http.Request) {
	util.EnableCors(&w)
	fileName := chi.URLParam(r, "uid")
	if !data.InMemoryArchives.ItemExist(fileName) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	data.InMemoryArchives.Remove(fileName)
	data.InMemoryUUID.Remove(fileName)

	presist.RemoveArchive(fileName)
	presist.RemoveUUID(fileName)
	presist.DeleteJob(fileName)

	data.RemoveFromDisk(fileName)
}

func CheckFileStatus(w http.ResponseWriter, r *http.Request) {
	util.EnableCors(&w)
	uid := chi.URLParam(r, "uid")

	_, err := uuid.Parse(uid)
	if err != nil {
		w.WriteHeader(400)
		return
	}

	if data.InMemoryArchives.ItemExist(uid) {
		w.WriteHeader(http.StatusOK)
		return
	}

	w.WriteHeader(http.StatusNotFound)
}

func DownloadFile(w http.ResponseWriter, r *http.Request) {
	util.EnableCors(&w)
	fileName := chi.URLParam(r, "uid")
	if !data.InMemoryArchives.ItemExist(fileName) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	w.Header().Set("Content-Type", "application/octet-stream")

	archiveExt := util.ExtBasedOnPlatform()
	filePath := "archives/" + fileName + archiveExt
	http.ServeFile(w, r, filePath)
}
