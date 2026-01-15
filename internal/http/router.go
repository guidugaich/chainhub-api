package http

import (
	"net/http"

	"chainhub-api/internal/http/handlers"
	authmw "chainhub-api/internal/http/middleware"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
)

func NewRouter(handler *handlers.Handler) http.Handler {
	r := chi.NewRouter()
	r.Use(chimw.RequestID)
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(authmw.CORS("http://localhost:3000"))

	r.Post("/signup", handler.Signup)
	r.Post("/login", handler.Login)
	r.Get("/healthz", handler.Health)
	r.Get("/tree/{username}", handler.GetTreeByUsername)

	r.Route("/links", func(r chi.Router) {
		r.Use(authmw.Auth(handler.Config.JWTSecret))
		r.Get("/", handler.ListLinks)
		r.Post("/", handler.CreateLink)
		r.Put("/{id}", handler.UpdateLink)
		r.Delete("/{id}", handler.DeleteLink)
	})

	return r
}
