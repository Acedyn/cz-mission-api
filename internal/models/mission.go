package models

import (
	"errors"
	"fmt"
	"html"
	"strings"
	"time"

	"gorm.io/gorm"

	"golang.org/x/exp/slices"

	"github.com/cardboard-citizens/cz-mission-api/internal/missions"
)

type Mission struct {
	gorm.Model

	ID               uint32    `gorm:"primary_key;auto_increment" json:"id"`
	Name             string    `gorm:"size:100"                   json:"name"`
	ShortDescription string    `gorm:"size:255"                   json:"short_description"`
	LongDescription  string    `gorm:"size:255"                   json:"long_description"`
	Class            string    `gorm:"size:100"                   json:"class"`
	Reward           float64   `gorm:"default:0"                  json:"reward"`
	Canceled         bool      `gorm:"default:false"              json:"canceled"`
	Initialized      bool      `gorm:"default:false"              json:"initialized"`
	CloseAt          time.Time `                                  json:"close_at"`

	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (mission *Mission) Format() string {
	return fmt.Sprintf("%s#%d", mission.Name, mission.ID)
}

func (mission *Mission) GetClassData() *missions.MissionClass {
	return missions.GetMissionsClasses()[mission.Class]
}

func (mission *Mission) Validate() (err error) {
	if mission.Name == "" {
		return errors.New("Invalid mission: No name provided")
	}
	missionClassKeys := missions.GetMissionClassKeys()
	if !slices.Contains(missionClassKeys, mission.Class) {
		return fmt.Errorf(
			"Invalid mission %s: Given class is not part of the available classes (%s)",
			mission.Name,
			mission.Class,
		)
	}

	mission.Name = html.EscapeString(strings.TrimSpace(mission.Name))
	mission.ShortDescription = html.EscapeString(
		strings.TrimSpace(mission.ShortDescription),
	)
	mission.LongDescription = html.EscapeString(
		strings.TrimSpace(mission.LongDescription),
	)
	return err
}
