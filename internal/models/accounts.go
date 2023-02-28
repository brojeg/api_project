package models

import (
	u "diploma/go-musthave-diploma-tpl/internal/utils"
	"os"

	"golang.org/x/crypto/bcrypt"

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

func (account *Account) Validate() u.Response {

	if len(account.Password) < 6 {
		return u.Message("Password is required", 400)
	}
	temp := &Account{}
	err := GetDB().Table("accounts").Where("login = ?", account.Login).First(temp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return u.Message("Connection error. Please retry", 502)
	}
	if temp.Login != "" {
		return u.Message("Email address already in use by another user.", 409)
	}
	return u.Message("Requirement passed", 200)
}

func (account *Account) Create() u.Response {

	if resp := account.Validate(); resp.ServerCode != 200 {
		return resp
	}
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	account.Password = string(hashedPassword)
	GetDB().Create(account)
	if account.ID <= 0 {
		return u.Message("Failed to create account, connection error.", 501)
	}
	tk := &Token{UserID: account.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, err := token.SignedString([]byte(os.Getenv("token_password")))
	if err != nil {
		panic(err)
	}
	account.Token = tokenString
	account.Password = ""
	response := u.Response{Message: account, ServerCode: 200}
	return response
}

func Login(email, password string) u.Response {
	account := &Account{}
	err := GetDB().Table("accounts").Where("login = ?", email).First(account).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.Message("Email address not found", 500)
		}
		return u.Message("Connection error. Please retry", 500)
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return u.Message("Invalid login credentials. Please try again", 401)
	}
	account.Password = ""
	tk := &Token{UserID: account.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	// account.Token = tokenString
	resp := u.Response{ServerCode: 200, Message: tokenString}

	return resp
}

func GetUser(u uint) *Account {

	acc := &Account{}
	GetDB().Table("accounts").Where("id = ?", u).First(acc)
	if acc.Login == "" { //User not found!
		return nil
	}

	acc.Password = ""
	return acc
}
