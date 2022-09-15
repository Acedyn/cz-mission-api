package missions

var DiscordMissions = map[string]*MissionType{
	"react-message": {
		Category: "discord",
		Validation: func() bool {
			return true
		},
	},
}
