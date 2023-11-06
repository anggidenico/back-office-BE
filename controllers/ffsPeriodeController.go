package controllers

import (
	"mf-bo-api/lib"
	"mf-bo-api/models"
	"net/http"

	"github.com/labstack/echo"
)

func GetFfsPeriodeController(c echo.Context) error {

	result := models.GetFfsPeriodeModels()

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = result
	return c.JSON(http.StatusOK, response)
}

func GetFfsPeriodeDetailController(c echo.Context) error {
	periodeKey := c.Param("periode_key")
	if periodeKey == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing peiode_key", "Missing peiode_key")
	}
	result := models.GetFfsPeriodeDetailModels(periodeKey)

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = result
	return c.JSON(http.StatusOK, response)
}
