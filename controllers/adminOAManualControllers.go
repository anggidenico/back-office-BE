package controllers

import (
	"mf-bo-api/lib"
	"mf-bo-api/models"
	"net/http"

	"github.com/labstack/echo"
)

func ManualOaRequestCreate(c echo.Context) error {
	responseData := make(map[string]interface{})

	step := c.Param("step")
	if step == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing: step", "Missing: step")
	}

	if step == "1" {
		err, oa_request_key := SaveStep1(c)
		responseData["oa_request_key"] = oa_request_key
		if err != nil {
			return lib.CustomError(http.StatusInternalServerError, err.Error(), err.Error())
		}
	} else if step == "2" {
		err, oa_request_key := SaveStep2(c)
		responseData["oa_request_key"] = oa_request_key
		if err != nil {
			return lib.CustomError(http.StatusInternalServerError, err.Error(), err.Error())
		}
	} else if step == "3" {
		err, oa_request_key := SaveStep3(c)
		responseData["oa_request_key"] = oa_request_key
		if err != nil {
			return lib.CustomError(http.StatusInternalServerError, err.Error(), err.Error())
		}
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = responseData

	return c.JSON(http.StatusOK, response)
}

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
