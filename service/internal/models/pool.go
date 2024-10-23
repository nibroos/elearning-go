package models

import "gorm.io/gorm"

type Pool struct {
	gorm.Model
	Group1ID uint32 `json:"group1_id" gorm:"column:group1_id"`
	Group2ID uint32 `json:"group2_id" gorm:"column:group2_id"`
	Mv1ID    uint32 `json:"mv1_id" gorm:"column:mv1_id"` // Typically user ID
	Mv2ID    uint32 `json:"mv2_id" gorm:"column:mv2_id"` // Typically role ID
}
