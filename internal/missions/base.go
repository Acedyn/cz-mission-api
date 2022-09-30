package missions

import (
	"github.com/cardboard-citizens/cz-mission-api/internal/models"
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

	missionGroups := []map[string]*MissionClass{
		TwitterMissions,
		DiscordMissions,
	}

	for _, missionGroup := range missionGroups {
		for missionName, Mission := range missionGroup {
			missions[missionName] = Mission
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
