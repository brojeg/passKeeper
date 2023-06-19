package models

import (
	auth "passKeeper/internal/models/auth"
)

type Account struct {
	ID       uint   `gorm:"primarykey"`
	Login    string `json:"login"`
	Password string `json:"password,omitempty"`
	Token    string `json:"token" sql:"-"`
}

func (account *Account) GetToken(jwtSettings auth.JWTSettings) string {
	return auth.GenerateToken(account.ID, jwtSettings)
}
