package missions

var DiscordMissions = map[string]*MissionClass{
	"react-message": {
		Category: "discord",
		Validation: func() bool {
			return true
		},
	},
}
