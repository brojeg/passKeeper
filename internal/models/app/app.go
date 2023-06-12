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

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type App struct {
	config        config.Config
	accountRepo   db.AccountRepository
	secretRepo    db.SecretRepository
	migrationRepo db.MigrationRepository
	Server        *http.Server
	JWTConf       auth.JWTSettings
}

func NewApp(config config.Config, accountRepo db.AccountRepository, secretRepo db.SecretRepository, migrationRepo db.MigrationRepository) *App {
	jwt := auth.InitJWTPassword(config.JWTPassword, config.ExpirationTime)
	return &App{config: config, accountRepo: accountRepo, secretRepo: secretRepo, migrationRepo: migrationRepo, JWTConf: jwt}
}

func (a App) CreateTables() {
	a.migrationRepo.AutoMigrate(&acc.Account{}, &sec.Secret{})
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

func (a *App) newRouter() *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Recoverer)

	accountHandler := handlers.NewAccountHandler(a.accountRepo, a.JWTConf)
	secretHandler := handlers.NewSecretHandler(a.secretRepo, a.JWTConf)

	router.Mount("/api/account", accountHandler.Route())
	router.Mount("/api/secret", secretHandler.Route())

	return router
}
