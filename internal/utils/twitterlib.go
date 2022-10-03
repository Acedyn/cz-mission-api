package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type TwitterError struct {
	Value        string `json:"value"`
	Detail       string `json:"detail"`
	Title        string `json:"title"`
	ResourceType string `json:"resource_type"`
	Parameter    string `json:"parameter"`
	ResourceId   string `json:"resource_id"`
	Type         string `json:"type"`
}

type TwitterData struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
}

type TwitterResponse struct {
	Errors []TwitterError
	Data   TwitterData
}

type TwitterUser struct {
	ID         uint64 `json:"id"`
	IDStr      string `json:"id_str"`
	Name       string `json:"name"`
	ScreenName string `json:"screen_name"`
}

type TwitterPostResponse struct {
	CreatedAt string      `json:"created_at"`
	ID        uint64      `json:"id"`
	User      TwitterUser `json:"user"`
}

type TwitterMeta struct {
	ReslutCount uint64 `json:"result_count"`
	NewestID    string `json:"newest_id"`
	OldestID    string `json:"oldest_id"`
}
type TwitterTweet struct {
	ID               string         `json:"id"`
	CreatedAt        string         `json:"created_at"`
	Text             string         `json:"text"`
	Type             string         `json:"type"`
	ReferencedTweets []TwitterTweet `json:"referenced_tweets"`
}

type TwitterRetweetsResponse struct {
	Data []TwitterTweet `json:"data"`
	Meta TwitterMeta    `json:"meta"`
}

type TwitterTweetResponse struct {
	Data TwitterTweet `json:"data"`
}

func GetTwitterIDFromUsername(username string) string {
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("https://api.twitter.com/2/users/by/username/%s", username),
		nil,
	)

	if err != nil {
		return ""
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("TWITTER_TOKEN")))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ""
	}

	var response TwitterResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return ""
	}

	return response.Data.ID
}

func GetTwitterTweetData(postID uint64) *TwitterTweetResponse {
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("https://api.twitter.com/2/tweets/%d?tweet.fields=created_at", postID),
		nil,
	)

	if err != nil {
		return nil
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("TWITTER_TOKEN")))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil
	}

	var response TwitterTweetResponse
	err = json.NewDecoder(resp.Body).Decode(&response)

	if err != nil {
		return nil
	}
	return &response
}

func GetTwitterUserData(username string) bool {
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("https://api.twitter.com/2/users/by/username/%s", username),
		nil,
	)

	if err != nil {
		return false
	}
	defer req.Body.Close()

	q := req.URL.Query()
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("TWITTER_TOKEN")))

	q.Add("tweet.fields", "created_at")
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return false
	}

	var response TwitterResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return false
	}

	return len(response.Errors) == 0
}

func IsUserLikedAPost(username string, postID uint64) bool {
	os.Setenv("TWITTER_TOKEN", "AAAAAAAAAAAAAAAAAAAAAAhbgwEAAAAAzLyI2%2FSavMB57c57BF8H3btj3JE%3DyH3d1fys9WgmFI3BuBkUuPqvTkYMPTgMkKWSRfaOkbOd5XSjBH")

	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("https://api.twitter.com/1.1/favorites/list.json?screen_name=%s&since_id=%d&max_id=%d&", username, postID-1, postID+1),
		nil,
	)

	if err != nil {
		return false
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("TWITTER_TOKEN")))

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	var posts []TwitterPostResponse
	err = json.NewDecoder(resp.Body).Decode(&posts)
	if err != nil {
		return false
	}

	return len(posts) == 1 && posts[0].ID == postID
}

func IsUserRetweetAPost(username string, postID uint64) bool {

	tweeData := GetTwitterTweetData(postID)
	if tweeData.Data.ID == "" {
		return false
	}

	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("https://api.twitter.com/2/users/%s/tweets?expansions=referenced_tweets.id&exclude=replies&start_time=%s", GetTwitterIDFromUsername(username), tweeData.Data.CreatedAt),
		nil,
	)

	if err != nil {
		return false
	}

	req.URL.Query()
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("TWITTER_TOKEN")))

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		return false
	}
	var retweet TwitterRetweetsResponse
	err = json.NewDecoder(resp.Body).Decode(&retweet)
	if err != nil {
		return false
	}

	sizeRetweeData := len(retweet.Data)
	for i := 0; i < sizeRetweeData; i++ {
		sizeReferencesTweetsData := len(retweet.Data[i].ReferencedTweets)
		for j := 0; j < sizeReferencesTweetsData; j++ {
			if retweet.Data[i].ReferencedTweets[j].ID == fmt.Sprintf("%d", postID) {
				return true
			}
		}
	}
	return false
}
