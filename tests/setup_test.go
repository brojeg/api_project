package main

import (
	"diploma/go-musthave-diploma-tpl/config"
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestAPISuite(t *testing.T) {
	config.Init()
	suite.Run(t, new(APISuite))
}
