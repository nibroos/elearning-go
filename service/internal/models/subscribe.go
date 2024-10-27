package models

import (
	"gorm.io/gorm"
)

type Subscribe struct {
	gorm.Model
	ID          uint   `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Name        string `json:"name" gorm:"column:name"`
	Description string `json:"desc" gorm:"column:description"`
	CreatedByID *uint  `json:"created_by_id" gorm:"column:created_by_id"`
	UpdatedByID *uint  `json:"updated_by_id" gorm:"column:updated_by_id"`
}
