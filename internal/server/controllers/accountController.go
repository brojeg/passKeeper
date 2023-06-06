package controllers

import (
	"encoding/json"
	"net/http"
	acc "passKeeper/internal/models/account"
	db "passKeeper/internal/models/database"
	server "passKeeper/internal/models/server"
)

func CreateAccount(repo db.DatabaseRepository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		account := &acc.Account{}
		err := json.NewDecoder(r.Body).Decode(account)
		if err != nil {
			server.RespondWithMessage(w, 400, "Invalid request")
		}
		resp := repo.CreateAccount(account)
		if resp.ServerCode == 200 {
			w.Header().Add("Authorization", account.Token)
		}
		server.RespondWithMessage(w, resp.ServerCode, resp.Message)
	}
}
