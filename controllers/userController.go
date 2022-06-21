package controllers

import (
	"mfbo_api/lib"
	"mfbo_api/models"
	"net/http"

	"github.com/labstack/echo"
)

func GetUserInfo(c echo.Context) error {

	var user models.User

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = user

	return c.JSON(http.StatusOK, response)
}
