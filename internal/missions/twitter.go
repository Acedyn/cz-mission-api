package missions

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cardboard-citizens/cz-mission-api/internal/models"
)

var TwitterMissions = map[string]*MissionClass{
	"like-post": {
		Category: "twitter",
		Validation: func(_ *models.Mission, user *models.User) (bool, error) {
			// twitterUsername, err := GetTwitterNickname(user.ID)
			// if err != nil {
			// 	return false, nil
			// }

			// return utils.IsUserLikedAPost(twitterUsername, <TWITTER_POSTID_HERE>), nil
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

func GetTwitterNickname(userid uint32) (string, error) {
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("http://localhost:8081/users/%d", userid),
		nil,
	)

	//os.Setenv("AUTHAPI_TOKEN", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE2NjQ5MzIyNTcsInVzZXJfaWQiOjF9.Frw_R-UyrxZg1xemMCnF5NElY7OQqvZlF2cDdTHIrDc")

	//req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE2NjQ5MzIyNTcsInVzZXJfaWQiOjF9.Frw_R-UyrxZg1xemMCnF5NElY7OQqvZlF2cDdTHIrDc"))
	if err != nil {
		return "", err
	}
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	var userData = struct {
		TwitterUsername string `json:"twitter_username"`
	}{}

	err = json.NewDecoder(resp.Body).Decode(&userData)

	if err != nil {
		return "", err
	}

	return userData.TwitterUsername, nil
}
