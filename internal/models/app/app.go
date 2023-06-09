package models

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	config "passKeeper/config/server"
	"passKeeper/internal/handlers"
	acc "passKeeper/internal/models/account"
	auth "passKeeper/internal/models/auth"
	db "passKeeper/internal/models/database"
	sec "passKeeper/internal/models/secret"
	"passKeeper/internal/server/controllers"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type App struct {
	config  config.Config
	repo    db.DatabaseRepository
	Server  *http.Server
	JWTConf auth.JWTSettings
}

func NewApp(config config.Config, repo db.DatabaseRepository) *App {
	jwt := auth.InitJWTPassword(config.JWTPassword, config.ExpirationTime)
	return &App{config: config, repo: repo, JWTConf: jwt}
}

func (a App) CreateTables() {
	a.repo.AutoMigrate(&acc.Account{}, &sec.Secret{})
}

func (a *App) StartWebServer() error {
	if a.config.TLSCertFile == "" || a.config.TLSKeyFile == "" || a.config.ServerPort == "" {
		return fmt.Errorf("server configuration is not complete")
	}
	a.Server = &http.Server{
		Addr:    a.config.ServerPort,
		Handler: a.newRouter(),
	}

	if err := http.ListenAndServeTLS(":443", a.config.TLSCertFile, a.config.TLSKeyFile, a.Server.Handler); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Printf("Cannot start http.ListenAndServe. Error is: /n %e", err)
	} else {
		log.Println("application stopped gracefully")
	}
	return nil
}

func (a App) newRouter() *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Recoverer)

	accountHandler := handlers.NewAccountHandler(a.repo, a.JWTConf)
	secretHandler := handlers.NewSecretHandler(a.repo)

	if err := http.ListenAndServeTLS(":443", a.config.TLSCertFile, a.config.TLSKeyFile, a.Server.Handler); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Printf("Cannot start http.ListenAndServe. Error is: /n %e", err)
	} else {
		log.Println("application stopped gracefully")
	}
	router.Post("/api/account/register", accountHandler.CreateAccount)
	router.Post("/api/account/login", accountHandler.Authenticate)
	router.Group(func(r chi.Router) {
		r.Use(controllers.JwtAuthenticationMiddleware(a.JWTConf))
		r.Get("/api/secret/{id}", secretHandler.GetSecret)
		r.Post("/api/secret", secretHandler.CreateSecret)
		r.Delete("/api/secret/{id}", secretHandler.DeleteSecret)
		r.Get("/api/secrets", secretHandler.GetSecrets)

	})

	return router
}
