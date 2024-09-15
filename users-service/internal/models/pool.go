package models

import "gorm.io/gorm"

type Pool struct {
	gorm.Model
	Name  string `json:"name" gorm:"column:name"`
	Tb1ID uint   `json:"tb1_id" gorm:"column:tb1_id"` // Typically user ID
	Tb2ID uint32 `json:"tb2_id" gorm:"column:tb2_id"` // Typically role ID
}
