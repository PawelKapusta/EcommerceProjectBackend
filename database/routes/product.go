package routes

import (
	"backend/database"
	"backend/database/models"
	"github.com/labstack/echo/v4"
	"net/http"
)

func GetProducts(c echo.Context) error {
	var products []models.Product
	database.GetDatabase().Find(&products)
	return c.JSON(http.StatusOK, products)
}

func GetProduct(c echo.Context) error {
	id := c.Param("id")
	var product models.Product

	result := database.GetDatabase().Find(&product, id)
	if result.Error != nil {
		return c.String(http.StatusNotFound, "Not found product")
	}

	return c.JSON(http.StatusOK, product)
}

func PostProduct(c echo.Context) error {
	product := new(models.Product)

	err := c.Bind(product)
	if err != nil {
		return c.String(http.StatusBadRequest, "Error in body"+err.Error())
	}

	result := database.GetDatabase().Create(product)
	if result.Error != nil {
		return c.String(http.StatusBadRequest, "Error when creating in Database... "+result.Error.Error())
	}

	return c.JSON(http.StatusOK, product)
}

//e.GET("/api/v1/products", func(c echo.Context) error {
//
//}

//
//e.GET("/api/v1/products/:id", func(c echo.Context) error {

//})
//
//e.POST("/api/v1/products", func(c echo.Context) error {

//})
