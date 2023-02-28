package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

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

func (order *Order) Validate() Response {

	if order.Number == "" {
		return Message("Order number should be on the payload", 500)
	}

	dbOrderExistsForUser := &Order{}
	errorbOrderExistsForUser := GetDB().Table("orders").Where("number = ? AND user_id = ?", order.Number, order.UserID).First(dbOrderExistsForUser).Error
	if errorbOrderExistsForUser != nil && errorbOrderExistsForUser != gorm.ErrRecordNotFound {
		return Message("Connection error. Please retry", 500)
	}
	if dbOrderExistsForUser.Number != "" {
		return Message("This order already in use by this user.", 200)
	}
	dbOrderExists := &Order{}
	errOrderExists := GetDB().Table("orders").Where("number = ?", order.Number).First(dbOrderExists).Error
	if errOrderExists != nil && errOrderExists != gorm.ErrRecordNotFound {
		return Message("Connection error. Please retry", 500)
	}
	if dbOrderExists.Number != "" {
		return Message("Order already in use by another user.", 409)
	}

	return Response{}
}

func (order *Order) Create() Response {

	if resp := order.Validate(); resp.ServerCode != 0 {
		return resp
	}

	GetDB().Create(order)

	resp := Message("success", 202)
	resp.Message = order
	return resp
}

func GetOrder(id uint) *Order {

	order := &Order{}
	err := GetDB().Table("orders").Where("id = ?", id).First(order).Error
	if err != nil {
		return nil
	}
	return order
}
func GetOrderByNumber(number string) *Order {

	order := &Order{}
	err := GetDB().Table("orders").Where("number = ?", number).First(order).Error
	if err != nil {
		return nil
	}
	return order
}

func GetOrders(user uint) []*Order {

	orders := make([]*Order, 0)
	err := GetDB().Table("orders").Where("user_id = ?", user).Find(&orders).Error
	if err != nil {
		logger.Error(err)
		return nil
	}

	return orders
}

func GetOrdersToApplyAccrual(status string) []*Order {

	orders := make([]*Order, 0)
	err := GetDB().Table("orders").Where("status = ?", status).Find(&orders).Error
	if err != nil {
		logger.Error(err)
		return nil
	}

	return orders
}

func (order *Order) Save() {

	GetDB().Save(order)

}
