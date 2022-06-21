package controllers

import (
	"mf-bo-api/lib"
	"net/http"

	"github.com/labstack/echo"
)

func HelloGreeting(c echo.Context) error {

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = "Hello from MF-BO-API.... it works ! ;)"

	return c.JSON(http.StatusOK, response)

}
