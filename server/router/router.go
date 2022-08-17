package router

import (
	"net/http"
	"server/handlers"

	"github.com/go-chi/chi/v5"
)

func CreateChiRouter(middlewares ...func(http.Handler) http.Handler) *chi.Mux {
	router := chi.NewRouter()
	for _, middleware := range middlewares {
		router.Use(middleware)
	}
	return router
}

func LoadRoutes(r *chi.Mux) {
	r.Route("/", func(r chi.Router) {
		r.Get("/test", handlers.CheckStatus)
		r.Get("/check/{uid}", handlers.CheckFileStatus)
		r.Get("/clear/{uid}", handlers.SessionClosed)
		r.Get("/download/{uid}", handlers.DownloadFile)

		r.Post("/upload", handlers.Upload)
	})
}
