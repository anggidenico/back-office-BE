package controllers

import (
	"mf-bo-api/lib"
	"mf-bo-api/models"
	"net/http"

	"github.com/labstack/echo"
)

func GetFfsPeriodeController(c echo.Context) error {
	var ffsperiode []models.FfsPeriode
	status, err := models.GetFfsPeriodeModels(&ffsperiode)
	if err != nil {
		return lib.CustomError(status, err.Error(), err.Error())
	}
	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = ffsperiode
	return c.JSON(http.StatusOK, response)
}

func GetFfsPeriodeDetailController(c echo.Context) error {
	periodeKey := c.Param("periode_key")
	if periodeKey == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing periode_key", "Missing periode_key")
	}
	var detperiode models.FfsPeriodeDetail
	status, err := models.GetFfsPeriodeDetailModels(&detperiode, periodeKey)
	if err != nil {
		return lib.CustomError(status, err.Error(), err.Error())
	}
	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = detperiode
	return c.JSON(http.StatusOK, response)
}
