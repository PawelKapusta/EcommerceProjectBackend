package main

import (
	"backend/database"
	"backend/database/routes"
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

	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAccessControlAllowOrigin, echo.HeaderAccessControlAllowCredentials},
		AllowMethods: []string{http.MethodGet, http.MethodPost},
	}))

	e.GET("/api/v1", func(c echo.Context) error {
		return c.String(http.StatusOK, "API!")
	})

	e.GET("/api/v1/products", routes.GetProducts)
	e.GET("/api/v1/products/:id", routes.GetProduct)
	e.POST("/api/v1/products", routes.PostProduct)
	e.GET("/api/v1/orders", routes.GetOrders)
	e.GET("/api/v1/orders/:id", routes.GetOrder)

	e.GET("/api/v1/auth/github", func(c echo.Context) error {
		url := routes.GetLoginURL(githubConfig)
		return c.JSON(http.StatusOK, map[string]string{"url": url})
	})

	e.GET("/auth/github/callback", func(c echo.Context) error {
		userToken := routes.GetTokenFromWeb(githubConfig, c.QueryParam("code"))
		userEmail, error := routes.GetUserInfoFromGithub(githubConfig.Client(context.Background(), userToken))
		if error != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "when getting userInfo"})
		}

		foundedUser := routes.FindUser(userEmail, "github")
		if foundedUser == false {
			routes.AddUser(userEmail, "github", *userToken)
		}
		u := routes.GetUser(userEmail, "github")

		c.Redirect(http.StatusFound, "http://localhost:3000/login/auth/github/"+u.GoToken+"&"+u.Email)
		return c.JSON(http.StatusOK, map[string]string{"token": userToken.AccessToken})
	})

	e.GET("/api/v1/auth/google", func(c echo.Context) error {
		url := routes.GetLoginURL(googleConfig)
		return c.JSON(http.StatusOK, map[string]string{"url": url})
	})

	e.GET("/auth/google/callback", func(c echo.Context) error {
		userToken := routes.GetTokenFromWeb(googleConfig, c.QueryParam("code"))
		userInfo, error := routes.GetUserInfoFromGoogle(googleConfig.Client(context.Background(), userToken))

		if error != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "when getting userInfo"})
		}

		foundedUser := routes.FindUser(userInfo.Email, "google")
		if foundedUser == false {
			routes.AddUser(userInfo.Email, "google", *userToken)
		}
		user := routes.GetUser(userInfo.Email, "google")

		c.Redirect(http.StatusFound, "http://localhost:3000/login/auth/google/"+user.GoToken+"&"+user.Email)

		return c.JSON(http.StatusOK, map[string]string{"token": userToken.AccessToken})
	})

	e.Logger.Fatal(e.Start(":8080"))
}
