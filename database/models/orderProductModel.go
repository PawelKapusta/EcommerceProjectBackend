package models

import (
	"gorm.io/gorm"
)

type OrderProduct struct {
	gorm.Model
	Name        string  `json:"name"`
	Price       float32 `json:"price"`
	Description string  `json:"description"`
	CategoryID  uint    `json:"categoryid"`
	CompanyID   uint    `json:"companyid"`
	ImageUrl    string  `json:"imageUrl"`
	OrderID     string  `json:"order_id"`
	Quantity    int     `json:"quantity"`
}
