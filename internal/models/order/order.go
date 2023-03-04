package models

import (
	"time"

	"github.com/jinzhu/gorm"

	accrual "diploma/go-musthave-diploma-tpl/internal/models/accrual"
	balance "diploma/go-musthave-diploma-tpl/internal/models/balance"
	db "diploma/go-musthave-diploma-tpl/internal/models/database"
	server "diploma/go-musthave-diploma-tpl/internal/models/server"
	log "diploma/go-musthave-diploma-tpl/pkg/logger"

	"go.uber.org/zap"
)

var logger *zap.SugaredLogger = log.Init()
var accrualURL string

type RawNumber struct {
	Number int64 `json:"-"`
}

type Order struct {
	ID         uint      `gorm:"primarykey" json:"-"`
	Number     string    `json:"number"`
	Status     string    `json:"status"`
	Accrual    float64   `json:"accrual"`
	UploadedAt time.Time `json:"-"`
	UserID     uint      `json:"-"`
}

func CreteTable() {
	db.Get().AutoMigrate(&Order{})
}
func InitAccrualURL(config string) {
	accrualURL = config
}
func (newOrder *Order) Validate() server.Response {

	if newOrder.Number == "" {
		return server.Message("Order number should be on the payload", 500)
	}

	dbOrderExistsForUser := &Order{}
	errorbOrderExistsForUser := db.Get().Table("orders").Where("number = ? AND user_id = ?", newOrder.Number, newOrder.UserID).First(dbOrderExistsForUser).Error
	if errorbOrderExistsForUser != nil && errorbOrderExistsForUser != gorm.ErrRecordNotFound {
		return server.Message("Connection error. Please retry", 500)
	}
	if dbOrderExistsForUser.Number != "" {
		return server.Message("This order already in use by this user.", 200)
	}
	dbOrderExists := &Order{}
	errOrderExists := db.Get().Table("orders").Where("number = ?", newOrder.Number).First(dbOrderExists).Error
	if errOrderExists != nil && errOrderExists != gorm.ErrRecordNotFound {
		return server.Message("Connection error. Please retry", 500)
	}
	if dbOrderExists.Number != "" {
		return server.Message("Order already in use by another user.", 409)
	}

	return server.Response{}
}

func (newOrder *Order) Create() server.Response {
	if resp := newOrder.Validate(); resp.ServerCode != 0 {
		return resp
	}
	db.Get().Create(newOrder)
	resp := server.Response{Message: newOrder, ServerCode: 200}
	return resp
}
func (newOrder *Order) ApplyAccrual() {
	var currentBalance *balance.Balance
	var accrualForOrder *accrual.Accrual
	if newOrder != nil {
		currentBalance = balance.Get(newOrder.UserID)
		accrualForOrder = accrual.RequestAccrual(accrualURL, newOrder.Number)
	}
	if accrualForOrder != nil {
		newOrder.Accrual = accrualForOrder.Accrual
		newOrder.Status = accrualForOrder.Status
		currentBalance.Add(newOrder.Accrual, newOrder.UserID)
		newOrder.Save()
		currentBalance.Save()
	}
}

func GetOrderByUser(id uint) *Order {

	order := &Order{}
	err := db.Get().Table("orders").Where("id = ?", id).First(order).Error
	if err != nil {
		return nil
	}
	return order
}
func GetOrderByNumber(number string) *Order {

	order := &Order{}
	err := db.Get().Table("orders").Where("number = ?", number).First(order).Error
	if err != nil {
		return nil
	}
	return order
}

func GetOrders(user uint) []*Order {

	orders := make([]*Order, 0)
	err := db.Get().Table("orders").Where("user_id = ?", user).Find(&orders).Error
	if err != nil {
		logger.Error(err)
		return nil
	}

	return orders
}

func GetOrdersToApplyAccrual(status string) []*Order {

	orders := make([]*Order, 0)
	err := db.Get().Table("orders").Where("status = ?", status).Find(&orders).Error
	if err != nil {
		logger.Error(err)
		return nil
	}

	return orders
}

func (newOrder *Order) Save() {
	db.Get().Save(newOrder)
}
