package controllers

import (
	"backend/database"
	"backend/database/models"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func GetUserController(e *echo.Group) {
	g := e.Group("/user")
	g.GET("", GetUsers)
	g.POST("/user", PostUser)
	g.DELETE("/:id", DeleteUser)
}

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

func GetUsers(c echo.Context) error {
	var users []models.User

	result := database.GetDatabase().Find(&users)
	if result.Error != nil {
		return c.String(http.StatusNotFound, "Users not found")
	}

	return c.JSON(http.StatusOK, users)
}

func PostUser(c echo.Context) error {
	var users []models.User

	result := database.GetDatabase().Find(&users)
	if result.Error != nil {
		return c.String(http.StatusNotFound, "Users not found")
	}

	return c.JSON(http.StatusOK, users)
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

func AddUserFromService(name string, surname string, username string, service string, email string, password string, token oauth2.Token) models.User {
	fmt.Println(email, service, token)
	user := new(models.User)
	user.Name = name
	user.Surname = surname
	user.Username = username
	user.Service = service
	user.Email = email
	user.Password = password
	user.Admin = false
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

func GetUserInfoFromGithub(client *http.Client) (string, string, error) {
	response, error := client.Get("https://api.github.com/user/emails")

	if error != nil {
		return "a", "", nil
	}
	defer response.Body.Close()
	data, error := ioutil.ReadAll(response.Body)
	if error != nil {
		return "a", "", nil
	}
	str := string(data)

	email := strings.Replace(strings.Split(strings.Split(strings.Split(str, "}")[0], ",")[0], ":")[1], "\"", "", -1)

	res, err := client.Get("https://api.github.com/user")

	if err != nil {
		return "a", "", nil
	}
	defer res.Body.Close()
	d, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "a", "", nil
	}
	str2 := string(d)
	username := strings.Replace(strings.Split(strings.Split(strings.Split(str2, "}")[0], ",")[0], ":")[1], "\"", "", -1)

	return email, username, error
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

func DeleteUser(c echo.Context) error {
	id := c.Param("id")
	var user models.User

	result := database.GetDatabase().Delete(&user, id)
	if result.Error != nil {
		return c.String(http.StatusNotFound, "User deleted")
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "item deleted"})
}