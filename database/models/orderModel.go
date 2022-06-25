package models

import (
	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	UserID     uint    `json:"userid"`
	TotalPrice float32 `json:"totalprice"`
	Items      []Item  `json:"items"`
	Street     string  `json:"street"`
	Nr         string  `json:"nr"`
	Code       string  `json:"code"`
	City       string  `json:"city"`
	Phone      string  `json:"phone"`
	IsPaid     bool    `json:"isPaid"`
	PaymentID  string  `json:"paymentID"`
	IsFinished bool    `json:"isFinished"`
}
