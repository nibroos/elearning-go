package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID        uint    `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Name      string  `json:"name" gorm:"column:name"`
	Username  *string `json:"username" gorm:"column:username;unique"`
	Email     string  `json:"email" gorm:"column:email;unique"`
	Password  string  `json:"-" gorm:"column:password"`
	Address   *string `json:"address" gorm:"column:address"`
	Roles     []Role  `json:"roles,omitempty" gorm:"many2many:user_roles"`
	CreatedAt string  `json:"created_at" gorm:"column:created_at"`
}
