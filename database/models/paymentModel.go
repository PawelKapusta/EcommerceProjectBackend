package models

import "gorm.io/gorm"

type Payment struct {
	gorm.Model
	PaymentType string `json:"paymenttype" validate:"required"`
	OrderID     uint   `json:"orderid" validate:"required"`
}
