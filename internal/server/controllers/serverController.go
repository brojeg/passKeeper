package controllers

import (
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/jinzhu/gorm"
)

func NewHTTPServer(port string, dbConn *gorm.DB) {

	router := NewRouter(dbConn)

	if err := http.ListenAndServe(port, router); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Printf("Cannot start http.ListenAndServe. Error is: /n %e", err)
	} else {
		log.Println("application stopped gracefully")
	}

}

func NewRouter(dbConn *gorm.DB) *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Recoverer)
	router.Post("/api/account/register", CreateAccount(dbConn))
	router.Post("/api/account/login", Authenticate(dbConn))
	router.Group(func(r chi.Router) {
		r.Use(JwtAuthenticationMiddleware)
		r.Get("/api/secret/{id}", GetSecret(dbConn))
		r.Post("/api/secret", CreateSecret(dbConn))
		r.Delete("/api/secret/{id}", DeleteSecret(dbConn))
		r.Get("/api/secrets", GetSecrets(dbConn))

	})

	return router
}
