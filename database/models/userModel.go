package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Username string `json:"username"`
	Service  string `json:"service" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password"`
	Admin    bool   `json:"admin"`
	Token    string `json:"token"`
	GoToken  string `json:"gotoken"`
}
