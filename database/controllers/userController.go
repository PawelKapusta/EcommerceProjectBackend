package controllers

import (
	"backend/database"
	"backend/database/authentication"
	"backend/database/models"
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/oauth2"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func GetUserController(e *echo.Group) {
	g := e.Group("/user")
	g.GET("", GetUsers)
	g.POST("/login", login)
	g.POST("/register", CreateUser)
	g.DELETE("/:id", DeleteUser, middleware.JWTWithConfig(authentication.GetCustomClaimsConfig()))
}

func login(c echo.Context) error {
	var user models.User
	username := c.FormValue("username")
	password := c.FormValue("password")

	result := database.GetDatabase().First(&user, "Username = ? AND Password = ?", username, password)
	if result.Error != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"code":    404,
			"message": "Database error " + result.Error.Error(),
		})
	} else {
		t, err := authentication.CreateToken(user)

		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, echo.Map{
			"token": t,
			"user":  user,
		})
	}
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

func CreateUser(c echo.Context) error {
	user := new(models.User)

	err := c.Bind(user)

	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid body "+err.Error())
	}

	errors := database.GetDatabase().Where("Username = ? AND Email = ?", user.Username, user.Email).Take(&user)

	if errors.RowsAffected == 0 {
		res := database.GetDatabase().Create(user)
		if res.Error != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"code":    403,
				"message": "Database error " + res.Error.Error(),
			})
		}
		t, err := authentication.CreateToken(*user)

		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, echo.Map{
			"token": t,
			"user":  user,
		})
	}

	if errors.RowsAffected > 0 {
		return c.String(http.StatusNotFound,
			"Database error "+errors.Error.Error(),
		)
	}
	return c.JSON(http.StatusOK, user)
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

	return c.JSON(http.StatusOK, map[string]string{"message": "User deleted"})
}
