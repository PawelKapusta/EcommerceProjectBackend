package models

import "gorm.io/gorm"

type ProductOrder struct {
	gorm.Model

	OrderId           uint
	ProductId         uint
	QuantityOfProduct int
}
