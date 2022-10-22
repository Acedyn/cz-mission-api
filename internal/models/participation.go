package models

import (
	"time"

	"gorm.io/gorm"
)

type Participation struct {
	gorm.Model
	ID uint32 `gorm:"primary_key" json:"id"`

	Users    []*User `gorm:"primary_key;many2many:user_participations;" json:"users"`
	Mission  Mission `gorm:"primary_key;foreignKey:id" json:"mission"`
	Progress float64 `gorm:"default:0" json:"progress"`
	Name     string  `gorm:"default:unnamed" json:"name"`

	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}
