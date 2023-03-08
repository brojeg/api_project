package models

import (
	db "diploma/go-musthave-diploma-tpl/internal/models/database"
	log "diploma/go-musthave-diploma-tpl/pkg/logger"
	"time"

	"go.uber.org/zap"
)

var logger *zap.SugaredLogger = log.Init()

func CreateTable() {
	db.Get().AutoMigrate(&BalanceHistory{})
}

type BalanceHistory struct {
	Order       string    `json:"order"`
	Sum         float64   `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
	UserID      uint      `json:"-"`
}

func (bh *BalanceHistory) Save() {
	db.Get().Save(bh)

}

func GetBalanceHistory(user uint) []*BalanceHistory {

	history := make([]*BalanceHistory, 0)
	err := db.Get().Table("balance_histories").Where("user_id = ?", user).Order("processed_at DESC").Find(&history).Error
	if err != nil {
		logger.Error(err)
		return nil
	}

	return history
}
