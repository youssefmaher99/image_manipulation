package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"server/data"
	"server/handlers"
	"server/logger"
	"server/models"
	"server/presist"
	"server/queue"
	"server/router"
	"server/worker"

	"github.com/go-chi/chi/v5/middleware"
)

var MyQueue *queue.Queue[models.Job]
var InMemoryArchives data.InMemory = make(data.InMemory)
var InMemoryUUID data.InMemory = make(data.InMemory)

func main() {
	MyQueue = queue.CreateQueue[models.Job]()

	// Inject global queue in package
	handlers.MyQueue = MyQueue

	// Inject UUID and archives in their package
	data.InMemoryArchives = InMemoryArchives
	data.InMemoryUUID = InMemoryUUID

	presist.Builder(MyQueue)
	presist.InitiateDestroyerWorker()
	data.RemoveDeadRefs()

	go func(MyQueue *queue.Queue[models.Job]) {
		worker.SpawnWorkers(MyQueue)
	}(MyQueue)

	//Graceful shutdown cleaning dirs
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	go func() {
		// blocking
		<-sigChan

		// clean up
		cleanDirs(map[string]string{"uploaded": "jpg", "filtered": "jpg", "archives": "gz"})
		// os.Exit(0)
	}()

	r := router.CreateChiRouter(
		middleware.Logger,
		middleware.Heartbeat("/health"),
	)
	router.LoadRoutes(r)

	logger.MyLog.Println("SERVER CONNECTED")
	logger.MyLog.Fatal(http.ListenAndServe(":5000", r))
}

func cleanDirs(dirs map[string]string) {
	for dir, ext := range dirs {
		cmd := exec.Command("bash", "-c", fmt.Sprintf("rm %s/*.%s", dir, ext))
		err := cmd.Run()
		// if dirs are empty err != nil so i continue to check the other dirs
		if err != nil {
			continue
		}
	}
	os.Exit(0)
}
