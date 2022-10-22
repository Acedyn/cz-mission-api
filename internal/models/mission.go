package models

import (
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Mission struct {
	gorm.Model

	ID               uint32          `gorm:"primary_key"                json:"id"`
	Name             string          `gorm:"size:100"                   json:"name"`
	ShortDescription string          `gorm:"size:255"                   json:"short_description"`
	LongDescription  string          `gorm:"size:255"                   json:"long_description"`
	Category         string          `gorm:"size:255"                   json:"category"`
	Logo             string          `gorm:"size:255"                   json:"logo"`
	Class            string          `gorm:"size:100"                   json:"class"`
	Reward           float64         `gorm:"default:0"                  json:"reward"`
	Canceled         bool            `gorm:"default:false"              json:"canceled"`
	Initialized      bool            `gorm:"default:false"              json:"initialized"`
	CloseAt          time.Time       `                                  json:"close_at"`
	Parameters       datatypes.JSON  `                             json:"parameters"`
	Participations   []Participation `gorm:"foreignKey:mission" json:"participation"`
	CreatedAt        time.Time       `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt        time.Time       `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (mission *Mission) Format() string {
	return fmt.Sprintf("%s#%d", mission.Name, mission.ID)
}

func (mission *Mission) GetParsedParameters() (map[string]string, error) {
	attributeMap := map[string]string{}
	err := json.Unmarshal([]byte(mission.Parameters.String()), &attributeMap)
	if err != nil {
		err = fmt.Errorf("Could not parse mission's parameters\n\t%s", err)
	}

	return attributeMap, err
}

func (mission *Mission) GetParameterValue(key string) string {
	attributeMap, err := mission.GetParsedParameters()
	if err != nil {
		return ""
	}

	value, ok := attributeMap[key]
	if !ok {
		return ""
	}

	return value
}
