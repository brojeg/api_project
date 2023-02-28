package controllers

import (
	"context"
	"diploma/go-musthave-diploma-tpl/internal/models"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func RequestAccrual(endpont, orderid string) *models.Accrual {
	accrual := &models.Accrual{}
	URL := endpont + "/api/orders/" + orderid
	resp, err := http.Get(URL)
	if err != nil {
		logger.Error(err)
	}
	errDecode := json.NewDecoder(resp.Body).Decode(accrual)
	if errDecode != nil {
		logger.Error(err)
	}
	// body, error := io.ReadAll(resp.Body)
	// if error != nil {
	// 	fmt.Println(error)
	// }
	// // close response body
	resp.Body.Close()

	// // print response body
	// fmt.Println(string(body))

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
				// var balance *models.Balance
				var accrualForOrder *models.Accrual
				// order.changeOrderStatus("REGISTERED")
				// order.changeOrderStatus("PROCESSING")
				order := models.GetOrderByNumber(order.Number)
				if order != nil {
					// balance = models.GetBalance(order.UserID)
					accrualForOrder = RequestAccrual(accrualURL, order.Number)
				}
				fmt.Println(accrualForOrder)
				if order != nil && accrualForOrder != nil {
					order.Accrual = accrualForOrder.Accrual
					order.Status = accrualForOrder.Status
					// balance.Add(order.Accrual, order.UserID)
					order.Save()
					// balance.Save()
				}
				// if accrualForOrder != nil {

				// }

				// order.changeOrderStatus("PROCESSED")

			}

		case <-ctx.Done():
			return

		}
	}

}
