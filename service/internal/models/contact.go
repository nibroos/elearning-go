package models

import (
	"time"

	"gorm.io/gorm"
)

type Contact struct {
	gorm.Model
	ID            uint       `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	TypeContactID uint       `json:"type_contact_id" gorm:"column:type_contact_id"`
	UserID        *uint      `json:"user_id" gorm:"column:user_id"`
	RefNum        string     `json:"ref_num" gorm:"column:ref_num"`
	Status        uint       `json:"status" gorm:"column:status"`
	OptionsJSON   *string    `json:"options_json" gorm:"column:options_json"`
	CreatedAt     *time.Time `json:"created_at" gorm:"column:created_at"`
	DeletedAt     *time.Time `json:"deleted_at" gorm:"column:deleted_at"`
}
