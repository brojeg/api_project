package main

import (
	"bytes"
	"math/rand"

	"diploma/go-musthave-diploma-tpl/internal/controllers"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"time"

	account "diploma/go-musthave-diploma-tpl/internal/models/account"
	db "diploma/go-musthave-diploma-tpl/internal/models/database"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var testLogin string = randString(4)
var testPassword string = randString(16)

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

func (suite *APISuite) TearDownSuite() {

	DeleteAccount(testLogin)

}

func (suite *APISuite) TestCreateAccount() {
	type testCase struct {
		name           string
		login          string
		password       string
		expectedResult int
	}

	testCases := []testCase{
		{
			name:           "// Test case 1: Verify that the /api/user/register endpoint returns the expected 200 code for non-existing user.",
			login:          testLogin,
			password:       testPassword,
			expectedResult: 200,
		},
		{
			name:           "// Test case 2:Verify that the /api/user/register endpoint returns the expected 409 code for the existing user.",
			login:          testLogin,
			password:       testPassword,
			expectedResult: 409,
		},
		{
			name:           "// Test case 2: Verify that the /api/user/register endpoint returns the expected 400 code for the account with missing Login field.",
			password:       testPassword,
			expectedResult: 400,
		},
		{
			name:           "// Test case 3: Verify that the /api/user/register endpoint returns the expected 400 code for the account with missing Password field.",
			login:          testLogin,
			expectedResult: 400,
		},
		{
			name:           "// Test case 4: Verify that the /api/user/register endpoint returns the expected 400 code for the empty request.",
			expectedResult: 400,
		},
		{
			name:           "// Test case 5: Verify that the /api/user/register endpoint returns the expected 400 code for the short Login field.",
			login:          "l",
			password:       testPassword,
			expectedResult: 400,
		},
		{
			name:           "// Test case 6: Verify that the /api/user/register endpoint returns the expected 400 code for the short Password field.",
			login:          testLogin,
			password:       "p",
			expectedResult: 400,
		},
	}

	for _, tc := range testCases {
		requestBody := &CreateUserRequest{
			Login:    tc.login,
			Password: tc.password,
		}
		requestBodyJSON, _ := json.Marshal(requestBody)
		req, err := http.NewRequest("POST", "/api/user/register", bytes.NewBuffer(requestBodyJSON))
		req.Header.Set("Content-Type", "application/json; charset=UTF-8")
		if err != nil {
			panic(err)
		}
		rr := httptest.NewRecorder()
		suite.router.ServeHTTP(rr, req)
		assert.Equal(suite.T(), tc.expectedResult, rr.Code)

	}
}

func (suite *APISuite) TestLoginAccount() {
	type testCase struct {
		name           string
		login          string
		password       string
		expectedResult int
	}

	testCases := []testCase{
		{
			name:           "// Test case 1: Verify that the /api/user/login endpoint returns the expected 200 code for the existing user.",
			login:          testLogin,
			password:       testPassword,
			expectedResult: 200,
		},
		{
			name:           "// Test case 2: Verify that the /api/user/login endpoint returns the expected 401 code for the non-existing user.",
			login:          randString(4),
			password:       testPassword,
			expectedResult: 401,
		},
		{
			name:           "// Test case 3: Verify that the /api/user/login endpoint returns the expected 401 code for the existing user but wrong password",
			login:          testLogin,
			password:       randString(16),
			expectedResult: 401,
		},
	}

	for _, tc := range testCases {
		requestBody := &CreateUserRequest{
			Login:    tc.login,
			Password: tc.password,
		}
		requestBodyJSON, _ := json.Marshal(requestBody)
		req, err := http.NewRequest("POST", "/api/user/login", bytes.NewBuffer(requestBodyJSON))
		req.Header.Set("Content-Type", "application/json; charset=UTF-8")
		if err != nil {
			panic(err)
		}
		rr := httptest.NewRecorder()
		suite.router.ServeHTTP(rr, req)
		assert.Equal(suite.T(), tc.expectedResult, rr.Code)

	}
}

func randString(length int) string {
	rand.Seed(time.Now().Unix())

	ranStr := make([]byte, length)
	for i := 0; i < length; i++ {
		ranStr[i] = byte(65 + rand.Intn(26))
	}
	str := string(ranStr)

	return str
}

func DeleteAccount(login string) {

	db.Get().Where("login = ?", login).Delete(&account.Account{})
}
