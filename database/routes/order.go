package routes

import (
	"backend/database"
	"backend/database/models"
	"github.com/labstack/echo/v4"
	"net/http"
)

func GetOrders(c echo.Context) error {
	var orders []models.Order
	database.GetDatabase().Find(&orders)
	return c.JSON(http.StatusOK, orders)
}

func GetOrder(c echo.Context) error {
	id := c.Param("id")
	var order models.Order

	result := database.GetDatabase().Find(&order, id)
	if result.Error != nil {
		return c.String(http.StatusNotFound, "Not found order")
	}

	return c.JSON(http.StatusOK, order)
}
