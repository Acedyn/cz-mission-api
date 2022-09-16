package missions

var TwitterMissions = map[string]*MissionClass{
	"like-post": {
		Category: "twitter",
		Validation: func() bool {
			return true
		},
	},
}
