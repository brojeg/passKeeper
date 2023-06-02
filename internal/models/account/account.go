package models

import (
	auth "passKeeper/internal/models/auth"
	server "passKeeper/internal/models/server"

	"github.com/jinzhu/gorm"
)

type Account struct {
	ID       uint   `gorm:"primarykey"`
	Login    string `json:"login"`
	Password string `json:"password,omitempty"`
	Token    string `json:"token" sql:"-"`
}

func (account *Account) Validate(dbConn *gorm.DB) server.Response {

	if len(account.Login) < 3 {
		return server.Message("Login is not valid", 400)
	}
	if len(account.Password) < 6 {
		return server.Message("Valid password is required", 400)
	}
	existingAccount := &Account{}
	err := dbConn.Table("accounts").Where("login = ?", account.Login).First(existingAccount).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return server.Message("Connection error. Please retry", 502)
	}
	if existingAccount.Login != "" {
		return server.Message("Email address already in use by another user.", 409)
	}
	return server.Message("Requirement passed", 200)
}

func (account *Account) Create(dbConn *gorm.DB) server.Response {
	if resp := account.Validate(dbConn); resp.ServerCode != 200 {
		return resp
	}
	account.Password = auth.EncryptPassword(account.Password)
	dbConn.Create(account)
	if account.ID == 0 {
		return server.Message("Failed to create account, connection error.", 501)
	}
	account.Token = account.getToken()
	account.Password = ""
	return server.Response{Message: account, ServerCode: 200}
}

func Login(email, password string, dbConn *gorm.DB) server.Response {
	account := &Account{}
	err := dbConn.Table("accounts").Where("login = ?", email).First(account).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return server.Message("Email address not found", 401)
		}
		return server.Message("Connection error. Please retry", 500)
	}

	if !auth.IsPasswordsEqual(account.Password, password) {
		return server.Message("Invalid login credentials. Please try again", 401)
	}
	tokenString := account.getToken()
	return server.Response{ServerCode: 200, Message: tokenString}
}
func (account *Account) getToken() string {

	return auth.GenerateToken(account.ID)
}
