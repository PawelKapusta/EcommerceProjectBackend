package controllers

import (
	"backend/database"
	"backend/database/authentication"
	"backend/database/models"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

const ProductNotFoundException = "Product not found"

func GetProductController(e *echo.Group) {
	g := e.Group("/product")
	g.GET("", GetProducts)
	g.GET("/:id", GetProduct)
	g.POST("", PostProduct, middleware.JWTWithConfig(authentication.GetCustomClaimsConfig()))
	g.PUT("/:id", UpdateProduct, middleware.JWTWithConfig(authentication.GetCustomClaimsConfig()))
	g.DELETE("/:id", DeleteProduct, middleware.JWTWithConfig(authentication.GetCustomClaimsConfig()))
}

func GetProducts(c echo.Context) error {
	var products []models.Product

	result := database.GetDatabase().Find(&products)
	if result.Error != nil {
		return c.String(http.StatusNotFound, ProductNotFoundException)
	}

	return c.JSON(http.StatusOK, products)
}

func GetProduct(c echo.Context) error {
	id := c.Param("id")
	var product models.Product

	result := database.GetDatabase().Find(&product, id)
	if result.Error != nil {
		return c.String(http.StatusNotFound, ProductNotFoundException)
	}

	return c.JSON(http.StatusOK, product)
}

func GetProductByID(id uint) (models.Product, error) {
	var product models.Product

	result := database.GetDatabase().First(&product, id)
	if result.Error != nil {
		return product, result.Error
	}

	return product, nil
}

func PostProduct(c echo.Context) error {
	product := new(models.Product)

	err := c.Bind(product)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid body "+err.Error())
	}

	result := database.GetDatabase().Create(product)
	if result.Error != nil {
		return c.String(http.StatusBadRequest, "Database error "+result.Error.Error())
	}

	return c.JSON(http.StatusOK, product)
}

func UpdateProduct(c echo.Context) error {
	id := c.Param("id")
	var product models.Product
	result := database.GetDatabase().Find(&product, id)
	if result.Error != nil {
		return c.String(http.StatusNotFound, ProductNotFoundException)
	}

	values := new(models.Product)
	err := c.Bind(values)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid body "+err.Error())
	}

	product.Name = values.Name
	product.Price = values.Price
	product.CategoryID = values.CategoryID
	product.CompanyID = values.CompanyID
	product.Description = values.Description
	database.GetDatabase().Save(&product)

	return c.JSON(http.StatusOK, product)
}

func DeleteProduct(c echo.Context) error {
	id := c.Param("id")
	var product models.Product

	result := database.GetDatabase().Delete(&product, id)
	if result.Error != nil {
		return c.String(http.StatusNotFound, "Product not found")
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Product deleted"})
}
