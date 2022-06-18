package models

import (
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model

	Name        string
	Price       float32
	Description string
	Category    string
	Quantity    int
	ImageUrl    string
}
