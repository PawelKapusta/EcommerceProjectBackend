package controllers

import (
	"backend/database"
	"backend/database/authentication"
	"backend/database/models"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

const CompanyNotFoundException = "Company not found"

func GetCompanyController(e *echo.Group) {
	g := e.Group("/company")
	g.GET("", GetCompanies)
	g.GET("/:id", GetCompany)
	g.POST("", PostCompany, middleware.JWTWithConfig(authentication.GetCustomClaimsConfig()))
	g.PUT("/:id", UpdateCompany, middleware.JWTWithConfig(authentication.GetCustomClaimsConfig()))
	g.DELETE("/:id", DeleteCompany, middleware.JWTWithConfig(authentication.GetCustomClaimsConfig()))
}

func GetCompanies(c echo.Context) error {
	var companies []models.Company

	result := database.GetDatabase().Find(&companies)
	if result.Error != nil {
		return c.String(http.StatusNotFound, CompanyNotFoundException)
	}

	return c.JSON(http.StatusOK, companies)
}

func GetCompany(c echo.Context) error {
	id := c.Param("id")
	var company models.Company

	result := database.GetDatabase().Find(&company, id)
	if result.Error != nil {
		return c.String(http.StatusNotFound, CompanyNotFoundException)
	}

	return c.JSON(http.StatusOK, company)
}

func PostCompany(c echo.Context) error {
	company := new(models.Company)
	err := c.Bind(company)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid body "+err.Error())
	}
	result := database.GetDatabase().Create(company)
	if result.Error != nil {
		return c.String(http.StatusBadRequest, "Database error "+result.Error.Error())
	}

	return c.JSON(http.StatusOK, company)
}

func UpdateCompany(c echo.Context) error {
	id := c.Param("id")
	var company models.Company
	result := database.GetDatabase().Find(&company, id)
	if result.Error != nil {
		return c.String(http.StatusNotFound, CompanyNotFoundException)
	}

	values := new(models.Company)
	err := c.Bind(values)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid body "+err.Error())
	}

	company.Name = values.Name
	company.Description = values.Description
	database.GetDatabase().Save(&company)

	return c.JSON(http.StatusOK, company)
}

func DeleteCompany(c echo.Context) error {
	id := c.Param("id")
	var company models.Company

	result := database.GetDatabase().Delete(&company, id)
	if result.Error != nil {
		return c.String(http.StatusNotFound, "Company not found")
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Company deleted"})
}
