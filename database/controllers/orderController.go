package controllers

import (
	"backend/database"
	"backend/database/models"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func GetOrderController(e *echo.Group) {
	g := e.Group("/order")
	g.GET("", GetOrders)
	g.GET("/:id", GetOrder)
	g.POST("", PostOrder)
}

func GetOrders(c echo.Context) error {
	var orders []models.Order

	result := database.GetDatabase().Find(&orders)
	if result.Error != nil {
		return c.String(http.StatusNotFound, "Order not found")
	}

	return c.JSON(http.StatusOK, orders)
}

func GetOrder(c echo.Context) error {
	id := c.Param("id")
	var order models.Order

	result := database.GetDatabase().Find(&order, id)
	if result.Error != nil {
		return c.String(http.StatusNotFound, "Order not found")
	}

	return c.JSON(http.StatusOK, order)
}

func PostOrder(c echo.Context) error {
	order := new(models.Order)

	order.IsPaid = false
	err := c.Bind(order)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid body "+err.Error())
	}

	result := database.GetDatabase().Create(order)
	if result.Error != nil {
		return c.String(http.StatusBadRequest, "Invalid "+result.Error.Error())
	}

	allItems := order.Items

	for i := 0; i < len(allItems); i++ {
		product, err := GetProductByID(allItems[i].ProductID)
		if err != nil {
			return c.String(http.StatusBadRequest, "Invalid body "+err.Error())
		}
		var parameters models.OrderProduct
		parameters.Name = product.Name
		parameters.Price = product.Price
		parameters.Description = product.Description
		parameters.CategoryID = product.CategoryID
		parameters.CompanyID = product.CompanyID
		parameters.ImageUrl = product.ImageUrl
		parameters.OrderID = strconv.FormatUint(uint64(order.ID), 10)
		parameters.Quantity = allItems[i].Quantity
		PostOrderProduct(parameters)

	}

	return c.JSON(http.StatusOK, order)
}
