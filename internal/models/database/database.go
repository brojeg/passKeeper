package models

import (
	"errors"
	"log"

	acc "passKeeper/internal/models/account"
	auth "passKeeper/internal/models/auth"
	sec "passKeeper/internal/models/secret"
	server "passKeeper/internal/models/server"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func ConnectDB(connStr string) *gorm.DB {
	conn, err := gorm.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error is %e \n Connection string is %s", err, connStr)
	}
	return conn
}

func GetAccountRepo(db *gorm.DB) AccountRepository {
	return &GormRepository{db: db}
}

func GetSecretRepo(db *gorm.DB) SecretRepository {
	return &GormRepository{db: db}
}

func GetMigrationRepo(db *gorm.DB) MigrationRepository {
	return &GormRepository{db: db}
}

type AccountRepository interface {
	CreateAccount(account *acc.Account, jwtSettings auth.JWTSettings) server.Response
	ValidateAccount(account *acc.Account) server.Response
	LoginAccount(email, password string, jwtSettings auth.JWTSettings) server.Response
}

type SecretRepository interface {
	GetSecretByID(secretID uint) (*sec.Secret, error)
	SaveSecret(s *sec.Secret) (*sec.Secret, error)
	GetSecretsForUser(userID uint) ([]sec.Secret, error)
	DeleteSecret(s *sec.Secret) error
}

type MigrationRepository interface {
	AutoMigrate(models ...interface{}) error
}

type GormRepository struct {
	db *gorm.DB
}

func (g *GormRepository) AutoMigrate(models ...interface{}) error {
	g.db = g.db.AutoMigrate(models...)
	if g.db.Error != nil {
		return g.db.Error
	}
	return nil
}
func (g *GormRepository) LoginAccount(email, password string, jwtSettings auth.JWTSettings) server.Response {
	account := &acc.Account{}
	err := g.db.Table("accounts").Where("login = ?", email).First(account).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return server.Message("Email address not found", 401)
		}
		return server.Message("Connection error. Please retry", 500)
	}

	if !auth.IsPasswordsEqual(account.Password, password) {
		return server.Message("Invalid login credentials. Please try again", 401)
	}
	tokenString := account.GetToken(jwtSettings)
	return server.Response{ServerCode: 200, Message: tokenString}
}

func (g *GormRepository) CreateAccount(account *acc.Account, jwtSettings auth.JWTSettings) server.Response {
	if resp := g.ValidateAccount(account); resp.ServerCode != 200 {
		return resp
	}
	account.Password = auth.EncryptPassword(account.Password)
	g.db.Create(account)
	if account.ID == 0 {
		return server.Message("Failed to create account, connection error.", 501)
	}
	account.Token = account.GetToken(jwtSettings)
	account.Password = ""
	return server.Response{Message: account, ServerCode: 200}
}
func (g *GormRepository) ValidateAccount(account *acc.Account) server.Response {
	if len(account.Login) < 3 {
		return server.Message("Login is not valid", 400)
	}
	if len(account.Password) < 6 {
		return server.Message("Valid password is required", 400)
	}
	existingAccount := &acc.Account{}
	err := g.db.Table("accounts").Where("login = ?", account.Login).First(existingAccount).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return server.Message("Connection error. Please retry", 502)
	}
	if existingAccount.Login != "" {
		return server.Message("Email address already in use by another user.", 409)
	}
	return server.Message("Requirement passed", 200)
}

func (g *GormRepository) DeleteSecret(s *sec.Secret) error {
	sec, err := g.GetSecretByID(s.ID)
	if err != nil {
		return err
	}
	if sec.UserID == s.UserID {
		result := g.db.Delete(s)
		if result.Error != nil {
			return result.Error
		}
	}
	return nil
}
func (g *GormRepository) GetSecretByID(secretID uint) (*sec.Secret, error) {
	secret := sec.Secret{}
	err := g.db.Table("secrets").Where("ID = ?", secretID).Find(&secret).Error
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &secret, nil
}

func (g *GormRepository) SaveSecret(s *sec.Secret) (*sec.Secret, error) {
	result := g.db.Save(s)
	if result.Error != nil || s.ID == 0 {
		return nil, errors.New("failed to save secret, connection error")
	}
	return s, nil
}
func (g *GormRepository) GetSecretsForUser(userID uint) ([]sec.Secret, error) {
	var secrets []sec.Secret
	result := g.db.Table("secrets").Where("User_ID = ?", userID).Find(&secrets)
	if result.Error != nil {
		return nil, result.Error
	}
	return secrets, nil
}
