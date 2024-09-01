package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `json:"name" gorm:"column:name"`
	Email    string `json:"email" gorm:"column:email;unique"`
	Password string `json:"-" gorm:"column:password"`
	Address  string `json:"address"`
	Roles    []Role `json:"roles,omitempty" gorm:"many2many:user_roles"`
}
