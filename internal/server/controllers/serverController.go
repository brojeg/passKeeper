package controllers

import (
	"errors"
	"log"
	"net/http"
	db "passKeeper/internal/models/database"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func NewHTTPServer(port string, repo db.DatabaseRepository) {

	router := NewRouter(repo)

	if err := http.ListenAndServeTLS(":443", "server.crt", "server.key", router); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Printf("Cannot start http.ListenAndServe. Error is: /n %e", err)
	} else {
		log.Println("application stopped gracefully")
	}

}

func NewRouter(repo db.DatabaseRepository) *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Recoverer)
	router.Post("/api/account/register", CreateAccount(repo))
	router.Post("/api/account/login", Authenticate(repo))
	router.Group(func(r chi.Router) {
		r.Use(JwtAuthenticationMiddleware)
		r.Get("/api/secret/{id}", GetSecret(repo))
		r.Post("/api/secret", CreateSecret(repo))
		r.Delete("/api/secret/{id}", DeleteSecret(repo))
		r.Get("/api/secrets", GetSecrets(repo))

	})

	return router
}
