package missions

import (
	"fmt"

	"github.com/cardboard-citizens/cz-mission-api/internal/models"
)

var (
	missionGroups = make([]map[string]*MissionClass, 0)
)

type MissionClass struct {
	Category    string
	Description string
	Validation  func(*models.Mission, *models.User) (bool, error)
	Logo        string
	Parameters  []string
}

func GetMissionsClasses() map[string]*MissionClass {
	missions := make(map[string]*MissionClass)

	for _, missionGroup := range missionGroups {
		for missionName, mission := range missionGroup {
			missions[fmt.Sprintf("%s-%s", mission.Category, missionName)] = mission
		}
	}

	return missions
}

func GetMissionClassKeys() []string {
	missionClasses := GetMissionsClasses()
	missionClassKeys := make([]string, 0)
	for key := range missionClasses {
		missionClassKeys = append(missionClassKeys, key)
	}
	return missionClassKeys
}
