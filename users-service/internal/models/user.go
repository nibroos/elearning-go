package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `json:"name" gorm:"column:name"`
	Username string `json:"username" gorm:"column:username"`
	Email    string `json:"email" gorm:"column:email;unique"`
	Password string `json:"-" gorm:"column:password"`
	Address  string `json:"address" gorm:"column:address"`
	Roles    []Role `json:"roles,omitempty" gorm:"many2many:user_roles"`
}
