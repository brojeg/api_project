package models

import (
	u "diploma/go-musthave-diploma-tpl/internal/utils"
	"os"

	"golang.org/x/crypto/bcrypt"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
)

/*
JWT claims struct
*/
type Token struct {
	UserID uint
	jwt.StandardClaims
}

// a struct to rep user account
type Account struct {
	ID       uint   `gorm:"primarykey"`
	Login    string `json:"login"`
	Password string `json:"password,omitempty"`
	Token    string `json:"token" sql:"-"`
}

// Validate incoming user details...
func (account *Account) Validate() u.Response {

	if len(account.Password) < 6 {
		return u.Message(false, "Password is required", 400)
	}

	//Email must be unique
	temp := &Account{}

	//check for errors and duplicate emails
	err := GetDB().Table("accounts").Where("login = ?", account.Login).First(temp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return u.Message(false, "Connection error. Please retry", 502)
	}
	if temp.Login != "" {
		return u.Message(false, "Email address already in use by another user.", 409)
	}

	return u.Message(true, "Requirement passed", 200)
}

func (account *Account) Create() (u.Response, string, uint) {

	if resp := account.Validate(); !resp.Status {
		return resp, "", 0
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	account.Password = string(hashedPassword)

	GetDB().Create(account)

	if account.ID <= 0 {
		return u.Message(false, "Failed to create account, connection error.", 501), "", 0
	}

	//Create new JWT token for the newly registered account
	tk := &Token{UserID: account.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, err := token.SignedString([]byte(os.Getenv("token_password")))
	if err != nil {
		panic(err)
	}
	account.Token = tokenString

	account.Password = "" //delete password

	response := u.Message(true, "Account has been created", 200)
	response.Message = account
	return response, account.Token, account.ID
}

func Login(email, password string) (u.Response, string) {

	account := &Account{}

	err := GetDB().Table("accounts").Where("login = ?", email).First(account).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.Message(false, "Email address not found", 500), ""
		}
		return u.Message(false, "Connection error. Please retry", 500), ""
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		return u.Message(false, "Invalid login credentials. Please try again", 401), ""
	}
	//Worked! Logged In
	account.Password = ""

	//Create JWT token
	tk := &Token{UserID: account.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	account.Token = tokenString //Store the token in the response

	resp := u.Message(true, "Logged In", 200)
	resp.Message = account
	return resp, tokenString
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
