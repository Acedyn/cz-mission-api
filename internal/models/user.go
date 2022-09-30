package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	ID             uint32           `gorm:"primary_key;auto_increment"      json:"id"`
	Points         float64          `gorm:"default:0"                       json:"points"`
	Participations []*Participation `gorm:"many2many:user_participations"   json:"participations"`
}
