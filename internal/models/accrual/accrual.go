package models

import (
	log "diploma/go-musthave-diploma-tpl/pkg/logger"
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

var logger *zap.SugaredLogger = log.Init()

type Accrual struct {
	Order   string  `json:"order"`
	Status  string  `json:"status,omitempty"`
	Accrual float64 `json:"accrual,omitempty"`
}

func RequestAccrual(endpont, orderid string) *Accrual {
	accrual := &Accrual{}
	URL := endpont + "/api/orders/" + orderid
	resp, err := http.Get(URL)
	if err != nil {
		logger.Error(err)
		return nil
	}
	errDecode := json.NewDecoder(resp.Body).Decode(accrual)
	if errDecode != nil {
		logger.Error(err)
		return nil
	}
	resp.Body.Close()
	return accrual
}
