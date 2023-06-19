package handlers

import (
	"encoding/json"
	"net/http"
	acc "passKeeper/internal/models/account"
	auth "passKeeper/internal/models/auth"
	db "passKeeper/internal/models/database"
	server "passKeeper/internal/models/server"

	"github.com/go-chi/chi"
)

type accountHandler struct {
	Repo        db.AccountRepository
	jwtSettings auth.JWTSettings
}

func NewAccountHandler(repo db.AccountRepository, jwtSettings auth.JWTSettings) *accountHandler {
	return &accountHandler{Repo: repo, jwtSettings: jwtSettings}
}

func (ah *accountHandler) Route() *chi.Mux {
	router := chi.NewRouter()
	router.Post("/register", ah.CreateAccount)
	router.Post("/login", ah.Authenticate)
	return router
}

func (ah *accountHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	account := &acc.Account{}
	err := json.NewDecoder(r.Body).Decode(account)
	if err != nil {
		server.RespondWithMessage(w, 400, "Invalid request")
	}
	resp := ah.Repo.CreateAccount(account, ah.jwtSettings)
	if resp.ServerCode == 200 {
		w.Header().Add("Authorization", account.Token)
	}
	server.RespondWithMessage(w, resp.ServerCode, resp.Message)
}
func (ah *accountHandler) Authenticate(w http.ResponseWriter, r *http.Request) {

	acc := &acc.Account{}
	err := json.NewDecoder(r.Body).Decode(acc)
	if err != nil || acc.Login == "" || acc.Password == "" {
		server.RespondWithMessage(w, 400, "Invalid request")
	}
	resp := ah.Repo.LoginAccount(acc.Login, acc.Password, ah.jwtSettings)
	w.Header().Add("Authorization", resp.Message.(string))
	server.RespondWithMessage(w, resp.ServerCode, resp.Message)

}
