package models

import (
	log "diploma/go-musthave-diploma-tpl/pkg/logger"

	"go.uber.org/zap"

	"time"
)

var logger *zap.SugaredLogger = log.Init()

type Balance struct {
	ID        uint    `gorm:"primarykey" json:"-"`
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
	UserID    uint    `json:"-"`
}

type BalanceHistory struct {
	Order       string    `json:"order"`
	Sum         float64   `json:"sum"`
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
func (balance *Balance) Save() {
	GetDB().Save(balance)

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

func (balance *Balance) Add(sum float64, user uint) Response {
	emptyBalance := Balance{UserID: user}
	if balance == nil {
		balance = &emptyBalance
		balance.Current = sum + balance.Current
		GetDB().Create(balance)
		resp := Message("success", 200)
		resp.Message = balance
		return resp
	}
	balance.Current = sum + balance.Current
	GetDB().Save(balance)

	resp := Message("success", 200)
	resp.Message = balance
	return resp
}
func (balance *Balance) Withdraw(sum float64) Response {

	if balance == nil {
		return Message("No active balance avalable", 402)
	}
	if balance.Current < sum {
		return Message("Not enough balance points", 402)
	}

	balance.Current = balance.Current - sum
	balance.Withdrawn = balance.Withdrawn + sum
	GetDB().Save(balance)

	resp := Message("success", 200)
	resp.Message = balance
	return resp
}
