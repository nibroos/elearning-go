package models

import (
	"gorm.io/gorm"
)

type Module struct {
	gorm.Model
	ID          uint   `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	ClassID     uint   `json:"class_id" gorm:"column:class_id"`
	Name        string `json:"name" gorm:"column:name"`
	Description string `json:"desc" gorm:"column:description"`
	LogoURL     string `json:"logo_url" gorm:"column:logo_url"`
	VideoURL    string `json:"video_url" gorm:"column:video_url"`
	CreatedByID *uint  `json:"created_by_id" gorm:"column:created_by_id"`
	UpdatedByID *uint  `json:"updated_by_id" gorm:"column:updated_by_id"`
}
