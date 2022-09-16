package missions

type MissionClass struct {
	Category   string
	Validation func() bool
}

type MissionInstance struct {
	Data string
	Type MissionClass
}

func GetMissionsClasses() map[string]*MissionClass {
	var missions = make(map[string]*MissionClass)

	var missionGroups = []map[string]*MissionClass{
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
