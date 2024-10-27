package models

import (
	"gorm.io/gorm"
)

type Class struct {
	gorm.Model
	ID          uint    `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Name        string  `json:"name" gorm:"column:name"`
	Description string  `json:"desc" gorm:"column:description"`
	BannerURL   *string `json:"banner_url" gorm:"column:banner_url"`
	LogoURL     *string `json:"logo_url" gorm:"column:logo_url"`
	VideoURL    *string `json:"video_url" gorm:"column:video_url"`
	SubscribeID uint    `json:"subscribe_id" gorm:"column:subscribe_id"`
	InchargeID  uint    `json:"incharge_id" gorm:"column:incharge_id"`
	CreatedByID *uint   `json:"created_by_id" gorm:"column:created_by_id"`
	UpdatedByID *uint   `json:"updated_by_id" gorm:"column:updated_by_id"`
}
