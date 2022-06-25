package controllers

import (
	"backend/database"
	"backend/database/models"
	"github.com/labstack/echo/v4"
	"net/http"
)

func GetOrderProductController(e *echo.Group) {
	g := e.Group("/orderproduct")
	g.GET("", GetOrdersProducts)
	g.GET("/:orderId", GetOrderProducts)
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

func PostOrderProduct(name string, price float32, description string,
	categoryId uint, companyId uint, imageUrl string, orderId string,
	quantity int) models.OrderProduct {
	product := new(models.OrderProduct)
	product.Name = name
	product.Price = price
	product.Description = description
	product.CategoryID = categoryId
	product.CompanyID = companyId
	product.ImageUrl = imageUrl
	product.OrderID = orderId
	product.Quantity = quantity
	database.GetDatabase().Create(product)

	return *product
}
