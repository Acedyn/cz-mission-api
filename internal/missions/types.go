package missions

type MissionType struct {
	Category   string
	Validation func() bool
}

type MissionInstance struct {
	Data string
	Type MissionType
}

func GetMissionsTypes() map[string]*MissionType {
	var missions = make(map[string]*MissionType)

	var missionGroups = []map[string]*MissionType{
		TwitterMissions,
	}

	for _, missionGroup := range missionGroups {
		for missionName, Mission := range missionGroup {
			missions[missionName] = Mission
		}
	}

	return missions
}

var Missions map[string]*MissionType = GetMissionsTypes()
