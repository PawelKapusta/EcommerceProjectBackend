package models

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	OrderPrice    float32
	PaymentMethod string
}
