package controllers

import (
	"mf-bo-api/lib"
	"mf-bo-api/models"
	"net/http"

	"github.com/labstack/echo"
)

func GetAllocSecController(c echo.Context) error {
	var allocsec []models.AllocSecurity
	status, err := models.GetAllocSecModels(&allocsec)
	if err != nil {
		// log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	var response lib.Response
	if len(allocsec) < 1 {
		response.Status.Code = http.StatusOK
		response.Status.MessageServer = "OK"
		response.Status.MessageClient = "OK"
		response.Data = []models.RiskProfile{}
	} else {
		response.Status.Code = http.StatusOK
		response.Status.MessageServer = "OK"
		response.Status.MessageClient = "OK"
		response.Data = allocsec
	}
	return c.JSON(http.StatusOK, response)
}
