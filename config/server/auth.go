package config

import (
	"encoding/json"
	"net/http"
	acc "passKeeper/internal/models/account"
	server "passKeeper/internal/models/server"
)

func (a *App) Authenticate(w http.ResponseWriter, r *http.Request) {

	acc := &acc.Account{}
	err := json.NewDecoder(r.Body).Decode(acc)
	if err != nil || acc.Login == "" || acc.Password == "" {
		server.RespondWithMessage(w, 400, "Invalid request")
	}
	resp := a.repo.LoginAccount(acc.Login, acc.Password)
	w.Header().Add("Authorization", resp.Message.(string))
	server.RespondWithMessage(w, resp.ServerCode, resp.Message)

}
