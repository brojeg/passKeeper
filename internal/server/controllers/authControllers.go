package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	account "passKeeper/internal/models/account"
	auth "passKeeper/internal/models/auth"
	db "passKeeper/internal/models/database"
	server "passKeeper/internal/models/server"
)

func Authenticate(repo db.DatabaseRepository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		acc := &account.Account{}
		err := json.NewDecoder(r.Body).Decode(acc)
		if err != nil || acc.Login == "" || acc.Password == "" {
			server.RespondWithMessage(w, 400, "Invalid request")
		}
		resp := repo.LoginAccount(acc.Login, acc.Password)
		w.Header().Add("Authorization", resp.Message.(string))
		server.RespondWithMessage(w, resp.ServerCode, resp.Message)
	}
}

var JwtAuthenticationMiddleware = func(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := auth.ValidateToken(r)
		if resp.ServerCode != 200 {
			server.RespondWithMessage(w, resp.ServerCode, resp.Message)
			return
		}

		ctx := context.WithValue(r.Context(), auth.ContextUserKey, resp.Message.(uint))
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})

}
