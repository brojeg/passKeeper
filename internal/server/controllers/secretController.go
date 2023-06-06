package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	auth "passKeeper/internal/models/auth"
	db "passKeeper/internal/models/database"
	sec "passKeeper/internal/models/secret"
	server "passKeeper/internal/models/server"
	"strconv"

	"github.com/go-chi/chi"
)

func CreateSecret(repo db.DatabaseRepository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := auth.GetUserFromContext(r.Context())
		if !ok {
			server.RespondWithMessage(w, 500, "Could not get user from context")
			return
		}

		var req sec.SecretRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			server.RespondWithMessage(w, 400, "Invalid request")
			return
		}

		value, err := sec.GetSecretFromRequest(req, user)
		if err != nil {
			server.RespondWithMessage(w, 500, "Could not create secret from request")
			return
		}

		secret, err := sec.NewSecret(user, req.Type, value, req.Meta)
		if err != nil {
			server.RespondWithMessage(w, 500, "Could not create secret")
			return
		}

		if req.ID != 0 {
			secret.ID = req.ID
		}

		savedSecret, err := repo.SaveSecret(&secret)
		if err != nil {
			log.Printf("cannot create secret - %s", secret.SecretType)
			server.RespondWithMessage(w, 500, "Could not save secret")
			return
		}

		server.RespondWithMessage(w, 200, savedSecret)
	}
}
func DeleteSecret(repo db.DatabaseRepository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := auth.GetUserFromContext(r.Context())
		if !ok {
			server.RespondWithMessage(w, 500, "Could not get user from context")
			return
		}
		id := chi.URLParam(r, "id")
		if id == "" {
			server.RespondWithMessage(w, 400, "Bad request. Id Is missing.")
		}
		i, err := strconv.Atoi(id)
		if err != nil {
			server.RespondWithMessage(w, 400, "Bad request.")
		}
		var secretToDelete sec.Secret
		secretToDelete.ID = uint(i)
		secretToDelete.UserID = user

		err = repo.DeleteSecret(&secretToDelete)
		if err != nil {
			server.RespondWithMessage(w, 500, "Could not delete secret")
			return
		}

		server.RespondWithMessage(w, 200, nil)
	}
}

func GetSecret(repo db.DatabaseRepository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		user, ok := auth.GetUserFromContext(r.Context())
		if !ok {
			server.RespondWithMessage(w, 500, "Could not get user from context")
		}
		id := chi.URLParam(r, "id")
		if id == "" {
			server.RespondWithMessage(w, 400, "Bad request. Id Is missing.")
		}
		i, err := strconv.Atoi(id)
		if err != nil {
			server.RespondWithMessage(w, 400, "Bad request.")
		}

		data, err := repo.GetSecretByID(uint(i))
		if err != nil {
			server.RespondWithMessage(w, 500, "Could not get Secret")
		}
		if data != nil {
			if data.UserID == user {

				server.RespondWithMessage(w, 200, data)
			}
		}
	}
}

func GetSecrets(repo db.DatabaseRepository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := auth.GetUserFromContext(r.Context())
		if !ok {
			server.RespondWithMessage(w, 500, "Could not get user from context")
			return
		}
		secrets, err := repo.GetSecretsForUser(user)
		if err != nil {
			server.RespondWithMessage(w, 500, "Could not get secrets")
			return
		}
		resp := server.Response{Message: secrets, ServerCode: 200}
		server.RespondWithMessage(w, resp.ServerCode, resp.Message)
	}
}
