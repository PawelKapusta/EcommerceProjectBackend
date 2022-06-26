package controllers

import (
	"backend/database"
	"backend/database/authentication"
	"backend/database/models"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

func GetOrderProductController(e *echo.Group) {
	g := e.Group("/orderproduct")
	g.GET("", GetOrdersProducts, middleware.JWTWithConfig(authentication.GetCustomClaimsConfig()))
	g.GET("/:orderId", GetOrderProducts, middleware.JWTWithConfig(authentication.GetCustomClaimsConfig()))
}

func GetOrdersProducts(c echo.Context) error {
	var orderProduct []models.OrderProduct

	result := database.GetDatabase().Find(&orderProduct)
	if result.Error != nil {
		return c.String(http.StatusNotFound, "Products not found")
	}

	return c.JSON(http.StatusOK, orderProduct)

}

func GetOrderProducts(c echo.Context) error {
	orderId := c.Param("orderId")
	var orderProducts []models.OrderProduct

	result := database.GetDatabase().Find(&orderProducts, "order_id = ?", orderId)

	if result.Error != nil {
		return c.String(http.StatusNotFound, "Order Products not found")
	}

	return c.JSON(http.StatusOK, orderProducts)
}

func PostOrderProduct(parameters models.OrderProduct) models.OrderProduct {
	product := new(models.OrderProduct)
	product.Name = parameters.Name
	product.Price = parameters.Price
	product.Description = parameters.Description
	product.CategoryID = parameters.CategoryID
	product.CompanyID = parameters.CompanyID
	product.ImageUrl = parameters.ImageUrl
	product.OrderID = parameters.OrderID
	product.Quantity = parameters.Quantity
	database.GetDatabase().Create(product)

	return *product
}
