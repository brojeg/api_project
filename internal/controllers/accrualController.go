package controllers

import (
	"context"
	"diploma/go-musthave-diploma-tpl/internal/models"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func RequestAccrual(endpont, orderid string) *models.Accrual {
	accrual := &models.Accrual{}
	resp, err := http.Get("endpont" + "/api/orders/" + orderid)
	if err != nil {
		log.Fatalln(err)
	}
	errDecode := json.NewDecoder(resp.Body).Decode(accrual)
	if errDecode != nil {
		log.Fatalln(err)
	}
	return accrual
}

func ApplyAccruals(ctx context.Context, interval, accrualURL string) {

	inter, err := time.ParseDuration(interval)
	if err != nil {
		logger.Errorf("Cannot parse PollInterval value from config. Error is: \n %e", err)
	}
	ticker := time.NewTicker(inter)
	for {

		select {
		case <-ticker.C:
			ordersToProcess := models.GetOrdersToApplyAccrual("NEW")

			for _, order := range ordersToProcess {
				// order.changeOrderStatus("REGISTERED")
				// order.changeOrderStatus("PROCESSING")
				order := models.GetOrderByNumber(order.Number)
				balance := models.GetBalance(order.UserID)
				accrual := RequestAccrual(accrualURL, order.Number)
				if accrual != nil {
					order.Accrual = float64(accrual.Accrual)
					order.Status = accrual.Status
					balance.Add(order.Accrual, order.UserID)
					order.Save()
					balance.Save()
				}

				// order.changeOrderStatus("PROCESSED")

			}

		case <-ctx.Done():
			return

		}
	}

}
