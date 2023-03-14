package main

import (
	"bytes"
	"fmt"
	"math/rand"

	"diploma/go-musthave-diploma-tpl/config"
	"diploma/go-musthave-diploma-tpl/internal/controllers"
	account "diploma/go-musthave-diploma-tpl/internal/models/account"
	auth "diploma/go-musthave-diploma-tpl/internal/models/auth"
	db "diploma/go-musthave-diploma-tpl/internal/models/database"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AccountTest struct {
	suite.Suite
	router *chi.Mux
}
type CreateUserRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (suite *AccountTest) SetupTest() {
	// Initialize the router
	suite.router = controllers.NewRouter()

}

func (suite *AccountTest) TearDownSuite() {

	//DeleteAccount(TestUser.Login)

}

func (suite *AccountTest) TestCreateAccount() {
	type testCase struct {
		name           string
		login          string
		password       string
		expectedResult int
	}

	testCases := []testCase{
		{
			name:           "// Test case 1: Verify that the /api/user/register endpoint returns the expected 200 code for non-existing user.",
			login:          TestData.Login,
			password:       TestData.Password,
			expectedResult: 200,
		},
		{
			name:           "// Test case 2:Verify that the /api/user/register endpoint returns the expected 409 code for the existing user.",
			login:          TestData.Login,
			password:       TestData.Password,
			expectedResult: 409,
		},
		{
			name:           "// Test case 3: Verify that the /api/user/register endpoint returns the expected 400 code for the account with missing Login field.",
			password:       TestData.Password,
			expectedResult: 400,
		},
		{
			name:           "// Test case 4: Verify that the /api/user/register endpoint returns the expected 400 code for the account with missing Password field.",
			login:          TestData.Login,
			expectedResult: 400,
		},
		{
			name:           "// Test case 5: Verify that the /api/user/register endpoint returns the expected 400 code for the empty request.",
			expectedResult: 400,
		},
		{
			name:           "// Test case 6: Verify that the /api/user/register endpoint returns the expected 400 code for the short Login field.",
			login:          "l",
			password:       TestData.Password,
			expectedResult: 400,
		},
		{
			name:           "// Test case 7: Verify that the /api/user/register endpoint returns the expected 400 code for the short Password field.",
			login:          TestData.Login,
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
		assert.Equal(suite.T(), tc.expectedResult, rr.Code, fmt.Sprintf("%v, want %v, got %v", tc.name, tc.expectedResult, rr.Code))

	}
}

func (suite *AccountTest) TestLoginAccount() {
	type testCase struct {
		name           string
		login          string
		password       string
		expectedResult int
	}

	testCases := []testCase{
		{
			name:           "// Test case 1: Verify that the /api/user/login endpoint returns the expected 200 code for the existing user.",
			login:          TestData.Login,
			password:       TestData.Password,
			expectedResult: 200,
		},
		{
			name:           "// Test case 2: Verify that the /api/user/login endpoint returns the expected 401 code for the non-existing user.",
			login:          randString(5),
			password:       TestData.Password,
			expectedResult: 401,
		},
		{
			name:           "// Test case 3: Verify that the /api/user/login endpoint returns the expected 401 code for the existing user but wrong password",
			login:          TestData.Login,
			password:       randString(17),
			expectedResult: 401,
		},
		{
			name:           "// Test case 4: Verify that the /api/user/login endpoint returns the expected 400 code for the missing Login fieled",
			password:       randString(17),
			expectedResult: 400,
		},
		{
			name:           "// Test case 4: Verify that the /api/user/login endpoint returns the expected 400 code for the missing Password fieled",
			login:          TestData.Login,
			expectedResult: 400,
		},
		{
			name:           "// Test case 4: Verify that the /api/user/login endpoint returns the expected 400 code for empty request",
			expectedResult: 400,
		},
	}

	for _, tc := range testCases {
		resposeCode, _ := Login(tc.login, tc.password, suite.router)
		assert.Equal(suite.T(), tc.expectedResult, resposeCode, fmt.Sprintf("%v, want %v, got %v", tc.name, tc.expectedResult, resposeCode))

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

	userID := GetUserIDFromToken(TestData.Token)
	DeletBalance(userID)
}

func Login(login, password string, router *chi.Mux) (int, string) {
	requestBody := &CreateUserRequest{
		Login:    login,
		Password: password,
	}
	requestBodyJSON, _ := json.Marshal(requestBody)
	req, err := http.NewRequest("POST", "/api/user/login", bytes.NewBuffer(requestBodyJSON))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	if err != nil {
		panic(err)
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	accessToken := rr.Header().Get("Authorization")

	return rr.Code, accessToken
}

func GetUserIDFromToken(token string) uint {
	tk := &auth.Token{StandardClaims: jwt.StandardClaims{ExpiresAt: int64(config.Param.ExpirationTime)}}
	_, err := jwt.ParseWithClaims(token, tk, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Param.JWTPassword), nil
	})
	if err != nil {
		panic(err)
	}
	return tk.UserID
}
