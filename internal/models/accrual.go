package models

import (
	"encoding/json"
	"net/http"
)

var accrualURL string

type Accrual struct {
	Order   string  `json:"order"`
	Status  string  `json:"status,omitempty"`
	Accrual float64 `json:"accrual,omitempty"`
}

func InitAccrualURL(config string) {
	accrualURL = config
}

func (order *Order) ApplyAccrual() {
	var balance *Balance
	var accrualForOrder *Accrual
	if order != nil {
		balance = GetBalance(order.UserID)
		accrualForOrder = RequestAccrual(accrualURL, order.Number)
	}
	if accrualForOrder != nil {
		order.Accrual = accrualForOrder.Accrual
		order.Status = accrualForOrder.Status
		balance.Add(order.Accrual, order.UserID)
		order.Save()
		balance.Save()
	}
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
