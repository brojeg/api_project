package controllers

import (
	"context"
	"diploma/go-musthave-diploma-tpl/internal/models"
	"encoding/json"
	"net/http"
	"time"
)

func RequestAccrual(endpont, orderid string) *models.Accrual {
	accrual := &models.Accrual{}
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

func ApplyAccruals(ctx context.Context, interval, accrualURL string) {

	inter, err := time.ParseDuration(interval)
	if err != nil {
		logger.Errorf("Cannot parse interval value from config. Error is: \n %e", err)
	}
	ticker := time.NewTicker(inter)
	for {

		select {
		case <-ticker.C:
			ordersToProcess := models.GetOrdersToApplyAccrual("NEW")
			for _, order := range ordersToProcess {
				var balance *models.Balance
				var accrualForOrder *models.Accrual
				order := models.GetOrderByNumber(order.Number)
				if order != nil {
					balance = models.GetBalance(order.UserID)
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

		case <-ctx.Done():
			return

		}
	}

}
