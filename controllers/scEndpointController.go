package controllers

import (
	"log"
	"mf-bo-api/lib"
	"mf-bo-api/models"
	"net/http"

	"github.com/labstack/echo"
)

func GetEndpointscController(c echo.Context) error {

	result := models.GetEndpointscModels()
	log.Println("Not Found")

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = result
	return c.JSON(http.StatusOK, response)
}

func GetEndpointDetailController(c echo.Context) error {
	endpointKey := c.Param("endpoint_key")
	if endpointKey == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing question key", "Missing question key")
	}
	result := models.GetDetailEndpointModels(endpointKey)
	// log.Println("Not Found")

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = result
	return c.JSON(http.StatusOK, response)
}
