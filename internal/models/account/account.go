package models

import (
	"os"

	"golang.org/x/crypto/bcrypt"

	db "diploma/go-musthave-diploma-tpl/internal/models/database"
	server "diploma/go-musthave-diploma-tpl/internal/models/server"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
)

type Token struct {
	UserID uint
	jwt.StandardClaims
}

type Account struct {
	ID       uint   `gorm:"primarykey"`
	Login    string `json:"login"`
	Password string `json:"password,omitempty"`
	Token    string `json:"token" sql:"-"`
}

func CreteTable() {
	db.Get().AutoMigrate(&Account{})
}

func (account *Account) Validate() server.Response {

	if len(account.Login) < 6 {
		return server.Message("Login is not valid", 400)
	}
	if len(account.Password) < 6 {
		return server.Message("Password is required", 400)
	}
	existingAccount := &Account{}
	err := db.Get().Table("accounts").Where("login = ?", account.Login).First(existingAccount).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return server.Message("Connection error. Please retry", 502)
	}
	if existingAccount.Login != "" {
		return server.Message("Email address already in use by another user.", 409)
	}
	return server.Message("Requirement passed", 200)
}

func (account *Account) Create() server.Response {

	if resp := account.Validate(); resp.ServerCode != 200 {
		return resp
	}
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	account.Password = string(hashedPassword)
	db.Get().Create(account)
	if account.ID <= 0 {
		return server.Message("Failed to create account, connection error.", 501)
	}
	tokenString := account.passwordHash()
	account.Token = tokenString
	return server.Response{Message: account, ServerCode: 200}
}

func Login(email, password string) server.Response {
	account := &Account{}
	err := db.Get().Table("accounts").Where("login = ?", email).First(account).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return server.Message("Email address not found", 500)
		}
		return server.Message("Connection error. Please retry", 500)
	}
	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return server.Message("Invalid login credentials. Please try again", 401)
	}
	tokenString := account.passwordHash()
	return server.Response{ServerCode: 200, Message: tokenString}
}

func (account *Account) passwordHash() string {
	tk := &Token{UserID: account.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, err := token.SignedString([]byte(os.Getenv("token_password")))
	if err != nil {
		panic(err)
	}
	account.Password = ""
	return tokenString
}
