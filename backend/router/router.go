package router

import (
	"net/http"
	"strings"

	"learnlang-backend/handlers"
	"learnlang-backend/utils"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func NewRouter() http.Handler {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	// CORS for frontend (Next.js dev on :3000 and production domain)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "https://localhost:3000", "https://learnlang.app", "https://www.learnlang.app"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	// Static uploads (served from local disk)
	fs := http.StripPrefix("/files/", http.FileServer(http.Dir(utils.UploadDir())))
	// fs := http.FileServer(http.Dir(utils.UploadDir()))
	r.Handle("/files/*", intercept(fs))

	// Routes
	r.Route("/api", func(r chi.Router) {
		r.Get("/languages", handlers.GetLanguagesHandler)

		r.Get("/packs", handlers.GetPacksHandler)
		r.Post("/packs", handlers.CreatePackHandler)
		r.Get("/packs/{id}", handlers.GetPackByIDHandler)

		r.Post("/vocabs", handlers.CreateVocabHandler)

		r.Get("/flashcards", handlers.GetFlashcardsHandler)
	})

	return r
}

// to prevent directory listing
func intercept(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}
