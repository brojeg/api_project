package main

import (
	"bytes"
	"math/rand"

	"diploma/go-musthave-diploma-tpl/config"
	"diploma/go-musthave-diploma-tpl/internal/controllers"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	account "diploma/go-musthave-diploma-tpl/internal/models/account"
	db "diploma/go-musthave-diploma-tpl/internal/models/database"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var testLogin string
var testPassword string

type APISuite struct {
	suite.Suite
	router *mux.Router
}
type CreateUserRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (suite *APISuite) SetupTest() {
	// Initialize the router
	suite.router = controllers.NewRouter()

}

func (suite *APISuite) TestCreateAccount() {
	// Create a request to the route
	requestBody := &CreateUserRequest{
		Login:    testLogin,
		Password: testPassword,
	}
	requestBodyJSON, _ := json.Marshal(requestBody)
	req, err := http.NewRequest("POST", "/api/user/register", bytes.NewBuffer(requestBodyJSON))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	if err != nil {
		panic(err)
	}

	// Create a response recorder to capture the response
	rr := httptest.NewRecorder()

	// Call the handler function with the request and response recorder
	suite.router.ServeHTTP(rr, req)

	// Check the response status code
	assert.Equal(suite.T(), http.StatusOK, rr.Code)

}

func (suite *APISuite) TestCreateExingingAccount() {
	// Create a request to the route
	requestBody := &CreateUserRequest{
		Login:    testLogin,
		Password: testPassword,
	}
	requestBodyJSON, _ := json.Marshal(requestBody)
	req, err := http.NewRequest("POST", "/api/user/register", bytes.NewBuffer(requestBodyJSON))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	if err != nil {
		panic(err)
	}

	// Create a response recorder to capture the response
	rr := httptest.NewRecorder()

	// Call the handler function with the request and response recorder
	suite.router.ServeHTTP(rr, req)

	// Check the response status code
	assert.Equal(suite.T(), http.StatusConflict, rr.Code)

}
func (suite *APISuite) TearDownSuite() {

	Delete(testLogin)

}

func TestAPISuite(t *testing.T) {
	config.Init()
	testLogin = randString(4)
	testPassword = randString(16)
	suite.Run(t, new(APISuite))
}

func randString(length int) string {
	rand.Seed(time.Now().Unix())

	ranStr := make([]byte, length)

	// Generating Random string
	for i := 0; i < length; i++ {
		ranStr[i] = byte(65 + rand.Intn(26))
	}

	// Converting byte slice to string
	str := string(ranStr)

	return str
}

func Delete(login string) {

	db.Get().Where("login = ?", login).Delete(&account.Account{})
}
