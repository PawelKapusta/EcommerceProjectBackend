package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email   string
	Service string
	Token   string
	GoToken string
}
