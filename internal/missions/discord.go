package missions

var DiscordMissions = map[string]*MissionClass{
	"react-message": {
		Category: "discord",
		Validation: func() bool {
			return true
		},
		Logo:        "https://assets.stickpng.com/images/5847f2d1cef1014c0b5e4870.png",
		Description: "React to a discord message",
	},
}
