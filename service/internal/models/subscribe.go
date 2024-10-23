package models

import (
	"gorm.io/gorm"
)

type Subscribe struct {
	gorm.Model
	ID          uint   `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Name        string `json:"name" gorm:"column:name"`
	Description string `json:"desc" gorm:"column:description"`
}
