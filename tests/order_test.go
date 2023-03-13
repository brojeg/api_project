package main

import (
	"bytes"
	"diploma/go-musthave-diploma-tpl/internal/controllers"
	db "diploma/go-musthave-diploma-tpl/internal/models/database"
	order "diploma/go-musthave-diploma-tpl/internal/models/order"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type OrderTest struct {
	suite.Suite
	router *mux.Router
}

func (suite *OrderTest) SetupTest() {

	suite.router = controllers.NewRouter()
	code, token := Login(TestData.Login, TestData.Password, suite.router)
	if code == 200 {
		TestData.Token = token
	}

}

func (suite *OrderTest) TearDownSuite() {
	// stringOrederID := strconv.Itoa(TestData.OrderID)
	// DeleteOrder(stringOrederID)
}

func (suite *OrderTest) TestCreateOrder() {
	type testCase struct {
		name             string
		token            string
		orderID          int
		wrongOrderIdType string
		expectedResult   int
	}

	testCases := []testCase{
		{
			name:           "// Test case 1: Verify that the /api/user/orders endpoint returns the expected 202 code for the non-exiting order.",
			token:          TestData.Token,
			orderID:        TestData.OrderID,
			expectedResult: 202,
		},
		{
			name:           "// Test case 2: Verify that the /api/user/orders endpoint returns the expected 200 code for the already uploaded orderID by this user",
			token:          TestData.Token,
			orderID:        TestData.OrderID,
			expectedResult: 200,
		},
		{
			name:           "// Test case 3: Verify that the /api/user/orders endpoint returns the expected 422 code for the non-Luhna number",
			token:          TestData.Token,
			orderID:        1,
			expectedResult: 422,
		},
		{
			name:             "// Test case 4: Verify that the /api/user/orders endpoint returns the expected 422 code for the wrong nubmer format",
			token:            TestData.Token,
			wrongOrderIdType: "test",
			expectedResult:   422,
		},
		{
			name:           "// Test case 5: Verify that the /api/user/orders endpoint returns the expected 401 code for the missing Authorization header",
			expectedResult: 401,
		},
	}

	for _, tc := range testCases {
		resposeCode := CreateOrder(tc.token, tc.wrongOrderIdType, tc.orderID, suite.router)
		assert.Equal(suite.T(), tc.expectedResult, resposeCode, fmt.Sprintf("%v, want %v, got %v", tc.name, tc.expectedResult, resposeCode))

	}
}

func (suite *OrderTest) TestGetOrder() {
	type testCase struct {
		name           string
		token          string
		expectedResult int
	}

	testCases := []testCase{
		{
			name:           "// Test case 1: Verify that the GET /api/user/orders endpoint returns the expected 200 code for the exiting user",
			token:          TestData.Token,
			expectedResult: 200,
		},
		{
			name:           "// Test case 2: Verify that the GET /api/user/orders endpoint returns the expected 401 code for the missing Authorization header",
			expectedResult: 401,
		},
		{
			name:           "// Test case 3: Verify that the GET /api/user/orders endpoint returns the expected 401 code for the non-valid Authorization header",
			token:          "some_token_string",
			expectedResult: 401,
		},
		{
			name:           "// Test case 4: Verify that the GET /api/user/orders endpoint returns the expected 401 code for the expired Authorization header",
			token:          "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOjEzNSwiZXhwIjoxNjc4NTUzMzcyfQ.InFcXDMjFUW7hX7xjjnXctHpV5MTGHj9G3ah0GznGNU",
			expectedResult: 401,
		},
	}

	for _, tc := range testCases {
		resposeCode := GetOrder(tc.token, suite.router)
		assert.Equal(suite.T(), tc.expectedResult, resposeCode, fmt.Sprintf("%v, want %v, got %v", tc.name, tc.expectedResult, resposeCode))

	}
	stringOrederID := strconv.Itoa(TestData.OrderID)
	DeleteOrder(stringOrederID)
	resp := GetOrder(TestData.Token, suite.router)
	assert.Equal(suite.T(), 204, resp, fmt.Sprintf("%v, want %v, got %v", "// Test case 5: Verify that the GET /api/user/orders endpoint returns the expected 204 code for the existing user with no orders", 204, resp))

}

func CreateOrder(token, badOrderid string, orderID int, router *mux.Router) int {
	requestBody := strconv.Itoa(orderID)
	if requestBody == "0" && badOrderid != "" {
		requestBody = badOrderid
	}
	req, err := http.NewRequest("POST", "/api/user/orders", bytes.NewBufferString(requestBody))
	req.Header.Add("Authorization", token)
	if err != nil {
		panic(err)
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	return rr.Code
}

func GetOrder(token string, router *mux.Router) int {

	req, err := http.NewRequest("GET", "/api/user/orders", nil)
	req.Header.Add("Authorization", token)
	if err != nil {
		panic(err)
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	return rr.Code
}

func DeleteOrder(orderID string) {

	db.Get().Where("number = ?", orderID).Delete(&order.Order{})
}
