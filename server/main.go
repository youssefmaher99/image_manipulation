package main

import (
	"server/router"

	"log"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := router.CreateChiRouter(middleware.Logger)
	router.LoadRoutes(r)

	err := http.ListenAndServe(":5000", r)
	if err != nil {
		log.Fatal(err)
	}
}
