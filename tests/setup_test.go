package main

import (
	"diploma/go-musthave-diploma-tpl/config"
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestAPISuite(t *testing.T) {
	config.Init()
	initTestData()
	suite.Run(t, new(AccountTest))
	suite.Run(t, new(OrderTest))
	DeleteAccount(TestData.Login)
}

var TestData Data

type Data struct {
	Login    string
	Password string
	Token    string
	OrderID  int
}

func initTestData() {
	TestData = Data{Login: randString(4), Password: randString(16), OrderID: 512556643322227}
}
