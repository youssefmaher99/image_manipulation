package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"server/handlers"
	"server/queue"
	"server/router"
	"server/worker"

	"log"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

var MyQueue *queue.Queue

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	MyQueue = queue.CreateQueue()

	go func(MyQueue *queue.Queue) {
		worker.SpawnWorkers(MyQueue)
	}(MyQueue)

	// Inject global queue in package
	handlers.MyQueue = MyQueue

	//Graceful shutdown cleaning dirs
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	go func() {
		// blocking
		<-sigChan

		// clean up
		cleanDirs(map[string]string{"uploaded": "jpg", "filtered": "jpg", "archives": "gz"})
	}()

	r := router.CreateChiRouter(middleware.Logger)
	router.LoadRoutes(r)

	log.Println("Server is connected")
	log.Fatal(http.ListenAndServe(":5000", r))
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
