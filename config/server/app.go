package config

import (
	"errors"
	"log"
	"net/http"
	db "passKeeper/internal/models/database"
	"passKeeper/internal/server/controllers"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func (a *App) NewWebServer(conf App) {
	a.repo = conf.repo
	a.Start()
}
func New(sc ServerConfig) *App {
	repo := db.Get(sc.Database)
	return &App{config: sc, repo: repo}
}

func (a App) Start() error {
	a.Server = &http.Server{
		Addr:    a.config.ServerPort,
		Handler: a.NewRouter(),
	}

	if err := http.ListenAndServeTLS(":443", "server.crt", "server.key", a.Server.Handler); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Printf("Cannot start http.ListenAndServe. Error is: /n %e", err)
	} else {
		log.Println("application stopped gracefully")
	}
	return nil
}

func (a App) NewRouter() *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Recoverer)
	router.Post("/api/account/register", a.CreateAccount)
	router.Post("/api/account/login", a.Authenticate)
	router.Group(func(r chi.Router) {
		r.Use(controllers.JwtAuthenticationMiddleware)
		r.Get("/api/secret/{id}", a.GetSecret)
		r.Post("/api/secret", a.CreateSecret)
		r.Delete("/api/secret/{id}", a.DeleteSecret)
		r.Get("/api/secrets", a.GetSecrets)

	})

	return router
}
