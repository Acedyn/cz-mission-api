package utils

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func CheckUserAuth(id, token string) bool {
	authUrl := os.Getenv("AUTH_API_HOST") + "/users/" + id
	client := http.Client{}
	request, _ := http.NewRequest("GET", authUrl, nil)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	response, err := client.Do(request)
	if err != nil {
		return false
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return false
	}

	return !strings.Contains(string(body), "error")
}

func ExtractToken(r *http.Request) string {
	keys := r.URL.Query()
	token := keys.Get("token")
	if token != "" {
		return token
	}
	bearerToken := r.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}
