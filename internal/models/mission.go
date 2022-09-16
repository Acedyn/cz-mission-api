package models

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"html"
	"strings"
	"time"

	"golang.org/x/exp/slices"

	"github.com/cardboard-citizens/cz-goodboard-api/internal/missions"
)

type Mission struct {
	gorm.Model

	ID               uint32    `gorm:"primary_key;auto_increment" json:"id"`
	Name             string    `gorm:"size:100" json:"name"`
	ShortDescription string    `gorm:"size:255" json:"short_description"`
	LongDescription  string    `gorm:"size:255" json:"long_description"`
	Class            string    `gorm:"size:100" json:"class"`
	Reward           float64   `gorm:"default:0" json:"reward"`
	CloseAt          time.Time `json:"close_at"`

	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (mission *Mission) Setup() (err error) {
	if mission.Name == "" {
		return errors.New("Could not setup new mission: No name provided")
	}
	missionClassKeys := missions.GetMissionClassKeys()
	if !slices.Contains(missionClassKeys, mission.Class) {
		return fmt.Errorf("Could not setup mission %s: Invalid class (%s)", mission.Name, mission.Class)
	}

	mission.ID = 0
	mission.Name = html.EscapeString(strings.TrimSpace(mission.Name))
	mission.ShortDescription = html.EscapeString(strings.TrimSpace(mission.ShortDescription))
	mission.LongDescription = html.EscapeString(strings.TrimSpace(mission.LongDescription))
	mission.Reward = 0
	mission.CloseAt = time.Now()
	mission.CreatedAt = time.Now()
	mission.UpdatedAt = time.Now()

	return err
}
