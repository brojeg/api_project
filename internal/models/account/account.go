package models

import (
	auth "diploma/go-musthave-diploma-tpl/internal/models/auth"
	db "diploma/go-musthave-diploma-tpl/internal/models/database"
	server "diploma/go-musthave-diploma-tpl/internal/models/server"

	"github.com/jinzhu/gorm"
)

type Account struct {
	ID       uint   `gorm:"primarykey"`
	Login    string `json:"login"`
	Password string `json:"password,omitempty"`
	Token    string `json:"token" sql:"-"`
}

func CreateTable() {
	db.Get().AutoMigrate(&Account{})
}

func (account *Account) Validate() server.Response {

	if len(account.Login) < 3 {
		return server.Message("Login is not valid", 400)
	}
	if len(account.Password) < 6 {
		return server.Message("Valid password is required", 400)
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
	account.Password = auth.EncryptPassword(account.Password)
	db.Get().Create(account)
	if account.ID == 0 {
		return server.Message("Failed to create account, connection error.", 501)
	}
	account.Token = account.getToken()
	account.Password = ""
	return server.Response{Message: account, ServerCode: 200}
}

func Login(email, password string) server.Response {
	account := &Account{}
	err := db.Get().Table("accounts").Where("login = ?", email).First(account).Error
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

func RefreshToken(user_id uint) server.Response {
	user := &Account{ID: user_id}
	err := db.Get().Table("accounts").Where("ID = ?", user_id).First(user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return server.Message("Email address not found", 401)
		}
		return server.Message("Connection error. Please retry", 500)
	}

	return server.Response{ServerCode: 200, Message: user.getToken()}
}
