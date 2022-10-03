package missions

import (
	"github.com/cardboard-citizens/cz-mission-api/internal/models"
)

var TwitterMissions = map[string]*MissionClass{
	"like-post": {
		Category: "twitter",
		Validation: func(_ *models.Mission, user *models.User) (bool, error) {
			//return utils.IsUserLikedAPost(user.TwitterUsername, 1568149043733225472), nil
			return true, nil
		},
		Logo:        "https://assets.stickpng.com/thumbs/580b57fcd9996e24bc43c53e.png",
		Description: "Like a twitter post ",
		Parameters:  []string{"link"},
	},
}

func init() {
	missionGroups = append(missionGroups, TwitterMissions)
}
