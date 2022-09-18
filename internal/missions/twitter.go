package missions

var TwitterMissions = map[string]*MissionClass{
	"like-post": {
		Category: "twitter",
		Validation: func() bool {
			return true
		},
		Logo:        "https://assets.stickpng.com/thumbs/580b57fcd9996e24bc43c53e.png",
		Description: "Like a twitter post ",
	},
}
