package controllers

import (
	"context"
	"diploma/go-musthave-diploma-tpl/internal/models"
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
			ordersToProcess := models.GetOrdersToApplyAccrual("NEW")
			for _, order := range ordersToProcess {
				order.ApplyAccrual()
			}

		case <-ctx.Done():
			return

		}
	}

}
