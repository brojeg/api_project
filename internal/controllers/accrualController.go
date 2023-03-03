package controllers

import (
	"context"
	order "diploma/go-musthave-diploma-tpl/internal/models/order"

	"time"
)

func ApplyAccruals(ctx context.Context, interval, accrualURL string) {

	inter, err := time.ParseDuration(interval)
	if err != nil {
		logger.Errorf("Cannot parse interval value from config. Error is: \n %e", err)
	}
	ticker := time.NewTicker(inter)
	for {

		select {
		case <-ticker.C:
			ordersToProcess := order.GetOrdersToApplyAccrual("NEW")
			for _, order := range ordersToProcess {
				order.ApplyAccrual()
			}

		case <-ctx.Done():
			return

		}
	}

}
