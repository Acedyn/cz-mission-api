package missions

var TwitterMissions = map[string]*MissionType{
	"lile-post": {
		Category: "twitter",
		Validation: func() bool {
			return true
		},
	},
}
