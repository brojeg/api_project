package main

import (
	"bytes"
	"diploma/go-musthave-diploma-tpl/internal/controllers"
	balance "diploma/go-musthave-diploma-tpl/internal/models/balance"
	db "diploma/go-musthave-diploma-tpl/internal/models/database"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type BalanceTest struct {
	suite.Suite
	router *mux.Router
}

func (suite *BalanceTest) SetupTest() {

	suite.router = controllers.NewRouter()
	code, token := Login(TestData.Login, TestData.Password, suite.router)
	if code == 200 {
		TestData.Token = token
	}

}

func (suite *BalanceTest) TestBalance() {
	type testCase struct {
		name           string
		token          string
		expectedResult int
	}

	testCases := []testCase{
		{
			name:           "// Test case 1: Verify that the GET /api/user/balance endpoint returns the expected 200 code for the user",
			token:          TestData.Token,
			expectedResult: 200,
		}, {
			name:           "// Test case 2: Verify that the GET /api/user/balance endpoint returns the expected 401 code for the wrong Authorization header",
			token:          "some_wrong_token",
			expectedResult: 401,
		},
		{
			name:           "// Test case 3: Verify that the GET /api/user/balance endpoint returns the expected 401 code for the missing Authorization header",
			expectedResult: 401,
		},
	}

	for _, tc := range testCases {
		resposeCode := GetBalance(tc.token, suite.router)
		assert.Equal(suite.T(), tc.expectedResult, resposeCode, fmt.Sprintf("%v, want %v, got %v", tc.name, tc.expectedResult, resposeCode))

	}
}
func (suite *BalanceTest) TestWithdraw() {
	type testCase struct {
		name           string
		token          string
		orderID        int
		Sum            int
		expectedResult int
	}

	testCases := []testCase{
		{
			name:           "// Test case 1: Verify that the POST /api/user/balance/withdraw endpoint returns the expected 402 code for the existing user and not enough points",
			token:          TestData.Token,
			orderID:        TestData.OrderID,
			Sum:            1,
			expectedResult: 402,
		},
		{
			name:           "// Test case 2: Verify that the POST /api/user/balance/withdraw endpoint returns the expected 422 code for the existing user and not Luhn order id",
			token:          TestData.Token,
			orderID:        1000,
			Sum:            100,
			expectedResult: 422,
		},
		{
			name:           "// Test case 3: Verify that the POST /api/user/balance/withdraw endpoint returns the expected 401 code for the wrong Authorization header",
			token:          "some_wrong_token",
			expectedResult: 401,
		},
		{
			name:           "// Test case 5: Verify that thePOST /api/user/balance/withdraw endpoint returns the expected 401 code for the missing Authorization header",
			expectedResult: 401,
		},
	}

	for _, tc := range testCases {
		resposeCode := WithdrawFromBalance(tc.token, tc.orderID, tc.Sum, suite.router)
		assert.Equal(suite.T(), tc.expectedResult, resposeCode, fmt.Sprintf("%v, want %v, got %v", tc.name, tc.expectedResult, resposeCode))

	}
}

func GetBalance(token string, router *mux.Router) int {

	req, err := http.NewRequest("GET", "/api/user/balance", nil)
	req.Header.Add("Authorization", token)
	if err != nil {
		panic(err)
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	return rr.Code
}

func WithdrawFromBalance(token string, orderID, sum int, router *mux.Router) int {
	type WithdrawRequest struct {
		Order string `json:"order"`
		Sum   int    `json:"sum"`
	}
	order := strconv.Itoa(orderID)
	requestBody := &WithdrawRequest{
		Order: order,
		Sum:   sum,
	}
	requestBodyJSON, errMarshal := json.Marshal(requestBody)
	if errMarshal != nil {
		panic(errMarshal)
	}
	req, err := http.NewRequest("POST", "/api/user/balance/withdraw", bytes.NewBuffer(requestBodyJSON))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Add("Authorization", token)
	if err != nil {
		panic(err)
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	return rr.Code
}
func DeletBalance(id uint) {

	db.Get().Where("user_id = ?", id).Delete(&balance.Balance{})
}
