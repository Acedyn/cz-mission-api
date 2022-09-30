package missions

import (
	"github.com/cardboard-citizens/cz-mission-api/internal/models"
)

var DiscordMissions = map[string]*MissionClass{
	"react-message": {
		Category: "discord",
		Validation: func(_ *models.Mission, _ *models.User) (bool, error) {
			return true, nil
		},
		Logo:        "https://assets.stickpng.com/images/5847f2d1cef1014c0b5e4870.png",
		Description: "React to a discord message",
		Parameters:  []string{"link"},
	},
}
