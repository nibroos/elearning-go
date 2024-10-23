package models

import "gorm.io/gorm"

type Permission struct {
	gorm.Model
	Name string `json:"name" gorm:"column:name"`
}
