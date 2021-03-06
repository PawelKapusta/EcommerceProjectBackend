package controllers

import (
	"backend/database"
	"backend/database/authentication"
	"backend/database/models"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

const CategoryNotFoundException = "Category not found"

func GetCategoryController(e *echo.Group) {
	g := e.Group("/category")
	g.GET("", GetCategories)
	g.GET("/:id", GetCategory)
	g.POST("", CreateCategory, middleware.JWTWithConfig(authentication.GetCustomClaimsConfig()))
	g.PUT("/:id", UpdateCategory, middleware.JWTWithConfig(authentication.GetCustomClaimsConfig()))
	g.DELETE("/:id", DeleteCategory, middleware.JWTWithConfig(authentication.GetCustomClaimsConfig()))
}

func GetCategories(c echo.Context) error { // c.Request().Host
	var categories []models.Category
	result := database.GetDatabase().Find(&categories)
	if result.Error != nil {
		return c.String(http.StatusNotFound, CategoryNotFoundException)
	}

	return c.JSON(http.StatusOK, categories)
}

func GetCategory(c echo.Context) error {
	id := c.Param("id")

	var category models.Category

	result := database.GetDatabase().Find(&category, id)
	if result.Error != nil {
		return c.String(http.StatusNotFound, CategoryNotFoundException)
	}

	return c.JSON(http.StatusOK, category)
}

func CreateCategory(c echo.Context) error {
	category := new(models.Category)

	err := c.Bind(category)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid body "+err.Error())
	}

	result := database.GetDatabase().Create(category)
	if result.Error != nil {
		return c.String(http.StatusBadRequest, "Invalid "+result.Error.Error())
	}

	return c.JSON(http.StatusOK, category)
}

func UpdateCategory(c echo.Context) error {
	id := c.Param("id")
	var category models.Category
	result := database.GetDatabase().Find(&category, id)
	if result.Error != nil {
		return c.String(http.StatusNotFound, CategoryNotFoundException)
	}

	values := new(models.Category)
	err := c.Bind(values)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid body "+err.Error())
	}

	category.Name = values.Name
	category.Description = values.Description
	database.GetDatabase().Save(&category)

	return c.JSON(http.StatusOK, category)
}

func DeleteCategory(c echo.Context) error {
	id := c.Param("id")
	var category models.Category

	result := database.GetDatabase().Delete(&category, id)
	if result.Error != nil {
		return c.String(http.StatusNotFound, "Category not found")
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Category deleted"})
}
