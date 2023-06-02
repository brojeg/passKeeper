package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	account "passKeeper/internal/models/account"
	auth "passKeeper/internal/models/auth"
	server "passKeeper/internal/models/server"

	"github.com/jinzhu/gorm"
)

func Authenticate(dbCon *gorm.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		acc := &account.Account{}
		err := json.NewDecoder(r.Body).Decode(acc)
		if err != nil || acc.Login == "" || acc.Password == "" {
			server.RespondWithMessage(w, 400, "Invalid request")
		}
		resp := account.Login(acc.Login, acc.Password, dbCon)
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
