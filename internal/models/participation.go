package models

import (
	"gorm.io/gorm"
	"time"
)

type Participation struct {
	gorm.Model

	Users    []*User `gorm:"primary_key;many2many:user_participations;" json:"users"`
	Mission  Mission `gorm:"primary_key;foreignKey:ID" json:"mission"`
	Progress float64 `gorm:"default:0" json:"progress"`
	Name     string  `gorm:"default:unnamed" json:"name"`

	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}
