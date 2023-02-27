package models

import (
	"context"
	u "diploma/go-musthave-diploma-tpl/internal/utils"
	log "diploma/go-musthave-diploma-tpl/pkg/logger"

	"go.uber.org/zap"

	"time"
)

var logger *zap.SugaredLogger = log.InitLogger()

type Balance struct {
	ID        uint  `gorm:"primarykey" json:"-"`
	Current   int64 `json:"current"`
	Withdrawn int64 `json:"withdrawn"`
	UserID    uint  `json:"-"`
}

type BalanceHistory struct {
	Order       string    `json:"order"`
	Sum         int64     `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
	UserID      uint      `json:"-"`
}

func GetBalance(id uint) *Balance {

	balance := &Balance{}
	err := GetDB().Table("balances").Where("user_id = ?", id).First(balance).Error
	if err != nil {
		return nil
	}
	return balance
}
func (bh *BalanceHistory) Save() {
	GetDB().Save(bh)

}

func GetBalanceHistory(user uint) []*BalanceHistory {

	history := make([]*BalanceHistory, 0)
	err := GetDB().Table("balance_histories").Where("user_id = ?", user).Order("processed_at	DESC").Find(&history).Error
	if err != nil {
		logger.Error(err)
		return nil
	}

	return history
}

func (balance *Balance) Add(sum int64, user uint) u.Response {
	emptyBalance := Balance{UserID: user}
	if balance == nil {
		balance = &emptyBalance
		balance.Current = sum + balance.Current
		GetDB().Create(balance)
		resp := u.Message(true, "success", 200)
		resp.Message = balance
		return resp
	}
	balance.Current = sum + balance.Current
	GetDB().Save(balance)

	resp := u.Message(true, "success", 200)
	resp.Message = balance
	return resp
}
func (balance *Balance) Withdraw(sum int64) u.Response {

	if balance == nil {
		return u.Message(false, "No active balance avalable", 402)
	}
	if balance.Current < sum {
		return u.Message(false, "Not enough balance points", 402)
	}

	balance.Current = balance.Current - sum
	balance.Withdrawn = balance.Withdrawn + sum
	GetDB().Save(balance)

	resp := u.Message(true, "success", 200)
	resp.Message = balance
	return resp
}

func ApplyAccruals(ctx context.Context, interval string) {
	inter, err := time.ParseDuration(interval)
	if err != nil {
		logger.Errorf("Cannot parse PollInterval value from config. Error is: \n %e", err)
	}
	ticker := time.NewTicker(inter)
	for {

		select {
		case <-ticker.C:
			ordersToProcess := GetOrdersToApplyAccrual("NEW")

			for _, order := range ordersToProcess {
				order.changeOrderStatus("REGISTERED")
				order.changeOrderStatus("PROCESSING")
				data := GetBalance(order.UserID)
				data.Add(order.Accrual, order.UserID)
				order.changeOrderStatus("PROCESSED")

			}

		case <-ctx.Done():
			return

		}
	}

}
