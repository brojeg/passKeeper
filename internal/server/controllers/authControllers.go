package controllers

import (
	"context"
	"net/http"
	auth "passKeeper/internal/models/auth"
	server "passKeeper/internal/models/server"
)

func JwtAuthenticationMiddleware(jwtSettings auth.JWTSettings) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			resp := auth.ValidateToken(r, jwtSettings)
			if resp.ServerCode != 200 {
				server.RespondWithMessage(w, resp.ServerCode, resp.Message)
				return
			}

			ctx := context.WithValue(r.Context(), auth.ContextUserKey, resp.Message.(uint))
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}
