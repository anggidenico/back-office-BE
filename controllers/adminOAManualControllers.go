package controllers

import (
	"mf-bo-api/lib"
	"mf-bo-api/models"
	"net/http"

	"github.com/labstack/echo"
)

func GetCountryList(c echo.Context) error {
	paramSearch := make(map[string]string)

	country_name := c.QueryParam("country_name")
	if country_name != "" {
		paramSearch["country_name"] = country_name
	}

	result := models.GetCountryList(paramSearch)

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = result

	return c.JSON(http.StatusOK, response)
}

func GetCityList(c echo.Context) error {
	paramSearch := make(map[string]string)

	city_level := c.Param("city_level")
	if city_level != "" {
		paramSearch["city_level"] = city_level
	}

	parent_key := c.QueryParam("parent_key")
	if parent_key != "" {
		paramSearch["parent_key"] = parent_key
	}

	city_name := c.QueryParam("city_name")
	if city_name != "" {
		paramSearch["city_name"] = city_name
	}

	result := models.GetCityList(paramSearch)

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = result

	return c.JSON(http.StatusOK, response)
}
