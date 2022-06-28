package controllers

import (
	"backend/database"
	"backend/database/authentication"
	"backend/database/models"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"strconv"
)

const OrderNotFoundException = "Order not found"

func GetOrderController(e *echo.Group) {
	g := e.Group("/order")
	g.GET("", GetOrders)
	g.GET("/:id", GetOrder, middleware.JWTWithConfig(authentication.GetCustomClaimsConfig()))
	g.POST("", PostOrder, middleware.JWTWithConfig(authentication.GetCustomClaimsConfig()))
	g.POST("/:id/:email/:paymentId", CreatePaymentInformation)
	g.DELETE("/:id", DeleteOrder, middleware.JWTWithConfig(authentication.GetCustomClaimsConfig()))
}

func GetOrders(c echo.Context) error {
	var orders []models.Order

	result := database.GetDatabase().Find(&orders)
	if result.Error != nil {
		return c.String(http.StatusNotFound, OrderNotFoundException)
	}

	return c.JSON(http.StatusOK, orders)
}

func GetOrder(c echo.Context) error {
	id := c.Param("id")
	var order models.Order

	result := database.GetDatabase().Find(&order, id)
	if result.Error != nil {
		return c.String(http.StatusNotFound, OrderNotFoundException)
	}

	return c.JSON(http.StatusOK, order)
}

func GetOrderById(id uint) (models.Order, error) {
	var order models.Order

	result := database.GetDatabase().First(&order, "ID = ?", id)
	if result.Error != nil {
		return order, result.Error
	}

	return order, nil
}

func PostOrder(c echo.Context) error {
	order := new(models.Order)

	order.IsPaid = false
	order.IsFinished = false
	err := c.Bind(order)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid body "+err.Error())
	}

	result := database.GetDatabase().Create(order)
	if result.Error != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"code":    403,
			"message": "Database order error " + result.Error.Error(),
		})
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

func CreatePaymentInformation(c echo.Context) error {
	id := c.Param("id")
	email := c.Param("email")
	paymentId := c.Param("paymentId")
	order := new(models.Order)

	result := database.GetDatabase().Find(&order, "ID = ?", id)
	if result == nil {
		return c.JSON(http.StatusNotFound, OrderNotFoundException)
	}

	if email != order.Email {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Not the same email"})
	}

	order.IsPaid = true
	order.PaymentID = paymentId
	order.IsFinished = true
	database.GetDatabase().Save(order)

	return c.JSON(http.StatusOK, order)
}

func DeleteOrder(c echo.Context) error {
	id := c.Param("id")
	var order models.Order

	result := database.GetDatabase().Delete(&order, id)
	if result.Error != nil {
		return c.String(http.StatusNotFound, OrderNotFoundException)
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Order deleted"})
}
