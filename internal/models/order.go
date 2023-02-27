package models

import (
	u "diploma/go-musthave-diploma-tpl/internal/utils"
	"time"

	"github.com/jinzhu/gorm"
)

type RawNumber struct {
	Number int64 `json:"-"`
}

type Order struct {
	//gorm.Model
	ID         uint      `gorm:"primarykey" json:"-"`
	Number     string    `json:"number"`
	Status     string    `json:"status"`
	Accrual    int64     `json:"accrual"`
	UploadedAt time.Time `json:"-"`
	UserID     uint      `json:"-"` //The user that this contact belongs to
}

func (order *Order) Validate() u.Response {

	if order.Number == "" {
		return u.Message(false, "Order number should be on the payload", 500)
	}

	dbOrderExistsForUser := &Order{}
	errorbOrderExistsForUser := GetDB().Table("orders").Where("number = ? AND user_id = ?", order.Number, order.UserID).First(dbOrderExistsForUser).Error
	if errorbOrderExistsForUser != nil && errorbOrderExistsForUser != gorm.ErrRecordNotFound {
		return u.Message(false, "Connection error. Please retry", 500)
	}
	if dbOrderExistsForUser.Number != "" {
		return u.Message(false, "This order already in use by this user.", 200)
	}
	dbOrderExists := &Order{}
	err := GetDB().Table("orders").Where("number = ?", order.Number).First(dbOrderExists).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return u.Message(false, "Connection error. Please retry", 500)
	}
	if dbOrderExists.Number != "" {
		return u.Message(false, "Order already in use by another user.", 409)
	}

	return u.Message(true, "success", 200)
}

func (order *Order) Create() u.Response {

	if resp := order.Validate(); !resp.Status {
		return resp
	}

	GetDB().Create(order)

	resp := u.Message(true, "success", 202)
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

func (order *Order) changeOrderStatus(status string) *Order {

	order.Status = status
	GetDB().Save(order)
	return order

}
