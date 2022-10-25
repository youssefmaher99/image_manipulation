package main

import (
	"server/handlers"
	"server/router"
	"server/util"

	"log"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

var MyQueue *util.Queue

func main() {

	MyQueue = util.CreateQueue()
	go func(MyQueue *util.Queue) {
		util.SpawnWorkers(MyQueue)
	}(MyQueue)

	r := router.CreateChiRouter(middleware.Logger)
	router.LoadRoutes(r)

	// Package injection
	handlers.MyQueue = MyQueue

	log.Println("Server is connected")
	log.Fatal(http.ListenAndServe(":5000", r))
}
