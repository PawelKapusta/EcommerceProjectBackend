package models

import "gorm.io/gorm"

type Item struct {
	gorm.Model
	ProductID uint `json:productid`
	Quantity  int  `json:"quantity"`
	OrderID   uint `json:orderid`
}
