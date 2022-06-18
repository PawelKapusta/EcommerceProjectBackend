package routes

import (
	"backend/database"
	"backend/database/models"
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type UserInfo struct {
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Profile       string `json:"profile"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Gender        string `json:"gender"`
}

func GetUser(email string, service string) models.User {
	var user models.User
	database.GetDatabase().Find(&user, "Email = ? AND Service = ?", email, service)
	return user
}

func FindUser(email string, service string) bool {
	var user models.User
	database.GetDatabase().Find(&user, "Email = ? AND Service = ?", email, service)
	if user.Email == "" {
		return false
	}
	return true
}

func AddUser(email string, service string, token oauth2.Token) models.User {
	user := new(models.User)
	user.Email = email
	user.Service = service
	user.Token = token.AccessToken
	user.GoToken = uuid.NewString()
	database.GetDatabase().Create(user)
	return GetUser(email, service)
}

func GetLoginURL(config *oauth2.Config) string {
	url := config.AuthCodeURL("state")
	return url
}

func GetTokenFromWeb(config *oauth2.Config, code string) *oauth2.Token {
	token, error := config.Exchange(context.Background(), code)

	if error != nil {
		log.Fatal(error)
	}

	return token
}

func GetUserInfoFromGithub(client *http.Client) (string, error) {
	response, error := client.Get("https://api.github.com/user/emails")
	if error != nil {
		return "a", error
	}
	defer response.Body.Close()
	data, error := ioutil.ReadAll(response.Body)
	if error != nil {
		return "a", error
	}
	str := string(data)
	email := strings.Replace(strings.Split(strings.Split(strings.Split(str, "}")[0], ",")[0], ":")[1], "\"", "", -1)
	return email, error
}

func GetUserInfoFromGoogle(client *http.Client) (*UserInfo, error) {
	response, error := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if error != nil {
		return nil, error
	}
	defer response.Body.Close()
	data, error := ioutil.ReadAll(response.Body)
	if error != nil {
		return nil, error
	}
	var result UserInfo
	if error := json.Unmarshal(data, &result); error != nil {
		return nil, error
	}
	return &result, nil
}
