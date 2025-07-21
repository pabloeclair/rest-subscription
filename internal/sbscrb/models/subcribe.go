package models

import (
	"time"

	"gorm.io/gorm"
)

type SubscribeDto struct {
	gorm.Model
	Id          uint       `gorm:"primaryKey" json:"id"`
	ServiceName string     `gorm:"size:255;not null;default:null" json:"service_name"`
	Price       int        `gorm:"not null;default:null" json:"price"`
	UserId      string     `gorm:"not null;default:null" json:"user_id"`
	StartDate   time.Time  `gorm:"not null;default:null" json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
}
