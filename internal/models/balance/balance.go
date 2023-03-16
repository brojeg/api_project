package models

import (
	db "diploma/go-musthave-diploma-tpl/internal/models/database"
	server "diploma/go-musthave-diploma-tpl/internal/models/server"
)

type Balance struct {
	ID        uint    `gorm:"primarykey" json:"-"`
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
	UserID    uint    `json:"-"`
}

func CreateTable() {
	db.Get().AutoMigrate(&Balance{})
}

func Create(id uint) {
	balance := &Balance{UserID: id}
	balance.Save()
}

func Get(id uint) *Balance {
	balance := &Balance{}
	err := db.Get().Table("balances").Where("user_id = ?", id).First(balance).Error
	if err != nil {
		return nil
	}
	return balance
}

func (balance *Balance) Save() {
	db.Get().Save(balance)

}

func (balance *Balance) Add(sum float64, user uint) server.Response {
	balance.Current = sum + balance.Current
	db.Get().Save(balance)
	resp := server.Message("success", 200)
	resp.Message = balance
	return resp
}
func (balance *Balance) Withdraw(sum float64) server.Response {

	if balance == nil {
		return server.Message("No active balance avalable", 402)
	}
	if balance.Current < sum {
		return server.Message("Not enough balance points", 402)
	}

	balance.Current = balance.Current - sum
	balance.Withdrawn = balance.Withdrawn + sum
	db.Get().Save(balance)

	resp := server.Message("success", 200)
	resp.Message = balance
	return resp
}
