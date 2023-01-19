package handlers

import (
	"io/fs"
	"log"
	"os"
	"path"
	"server/logger"
	"server/models"
	"server/util"
	"testing"

	"github.com/google/uuid"
)

var dir = "/home/youssef/Desktop/image data set/PetImages/Cat/"

func TestUpload(t *testing.T) {
	os.Chdir("..")
	files := readNFiles(1000)
	uuid := generateUUID()
	job := models.Job{Uid: generateUUID(), Filter: "gray"}
	for _, file := range files {

		f, err := os.Open(dir + file.Name())
		if err != nil {
			logger.MyLog.Fatal(err)
			return
		}
		defer f.Close()

		_, err = util.SaveFile(f, file.Name(), uuid)
		if err != nil {
			logger.MyLog.Fatal(err)
		}

		img := models.Image{Name: file.Name(), Path: path.Join("uploaded", uuid+"_"+file.Name())}
		job.Images = append(job.Images, img)
	}

	// data.InMemoryUUID.Add(uuid, struct{}{})
	// presist.AddUUID(uuid)

	// MyQueue.Enqueue(job)
	// presist.AddJob(job)
}

func generateUUID() string {
	return uuid.New().String()
}

// read first n files in certain dir
func readNFiles(n int) []fs.DirEntry {
	files, err := os.ReadDir("/home/youssef/Desktop/image data set/PetImages/Cat/")
	if err != nil {
		log.Fatal(err)
	}
	return files[:n]
}

/*

	domain logic and http are tightly coupled


*/
