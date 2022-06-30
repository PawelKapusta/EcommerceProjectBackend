package main

import (
	"backend/database"
	"backend/database/authentication"
	"backend/database/controllers"
	"context"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
	"log"
	"net/http"
	"os"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	githubConfig := &oauth2.Config{
		ClientID:     os.Getenv("GITHUB_CLIENTID"),
		ClientSecret: os.Getenv("GITHUB_CLIENTKEY"),
		Endpoint:     github.Endpoint,
		Scopes:       []string{"user:email"},
	}

	googleConfig := &oauth2.Config{
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_LINK"),
		ClientID:     os.Getenv("GOOGLE_CLIENTID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENTKEY"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
			"openid",
		},
		Endpoint: google.Endpoint,
	}

	database.Connect()
	database.InitDefaultDatabase()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"https://ecommercef.azurewebsites.net", "http://localhost:3000", "http://localhost:8080"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
	}))

	e.GET("/api/v1", func(c echo.Context) error {
		return c.String(http.StatusOK, "API!")
	})

	g := e.Group("/api/v1")

	controllers.GetCategoryController(g)
	controllers.GetCompanyController(g)
	controllers.GetProductController(g)
	controllers.GetUserController(g)
	controllers.GetOrderController(g)
	controllers.GetOrderProductController(g)
	controllers.GetPaymentController(g)

	e.GET("/api/v1/auth/github", func(c echo.Context) error {
		url := controllers.GetLoginURL(githubConfig)
		return c.JSON(http.StatusOK, map[string]string{"url": url})
	})

	e.GET("/auth/github/callback", func(c echo.Context) error {
		userToken := controllers.GetTokenFromWeb(githubConfig, c.QueryParam("code"))
		userEmail, username, error := controllers.GetUserInfoFromGithub(githubConfig.Client(context.Background(), userToken))
		if error != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "when getting userInfo"})
		}

		foundedUser := controllers.FindUser(userEmail, "github")
		if foundedUser == false {
			controllers.AddUserFromService("GithubName", "GithubSurname", username, "github", userEmail, username+userEmail, *userToken)
		}
		u := controllers.GetUser(userEmail, "github")
		t, err := authentication.CreateToken(u)

		if err != nil {
			return err
		}

		c.Redirect(http.StatusFound, "https://ecommercef.azurewebsites.net/login/auth/github/"+t+"&"+u.Email)
		return c.JSON(http.StatusOK, map[string]string{"token": userToken.AccessToken})
	})

	e.GET("/api/v1/auth/google", func(c echo.Context) error {
		url := controllers.GetLoginURL(googleConfig)
		return c.JSON(http.StatusOK, map[string]string{"url": url})
	})

	e.GET("/auth/google/callback", func(c echo.Context) error {
		userToken := controllers.GetTokenFromWeb(googleConfig, c.QueryParam("code"))
		userInfo, error := controllers.GetUserInfoFromGoogle(googleConfig.Client(context.Background(), userToken))

		if error != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "when getting userInfo"})
		}

		foundedUser := controllers.FindUser(userInfo.Email, "google")
		if foundedUser == false {
			controllers.AddUserFromService(userInfo.GivenName, userInfo.FamilyName, userInfo.Name, "google", userInfo.Email, userInfo.Sub, *userToken)
		}
		user := controllers.GetUser(userInfo.Email, "google")
		t, err := authentication.CreateToken(user)

		if err != nil {
			return err
		}
		c.Redirect(http.StatusFound, "https://ecommercef.azurewebsites.net/login/auth/google/"+t+"&"+user.Email)

		return c.JSON(http.StatusOK, echo.Map{
			"token": t,
			"user":  user,
		})
	})

	e.Logger.Fatal(e.Start(":8080"))
}
