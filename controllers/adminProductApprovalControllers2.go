package controllers

import (
	"mf-bo-api/lib"
	"net/http"

	"github.com/labstack/echo"
)

func ProductApprovalAction(c echo.Context) error {

	

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = nil
	return c.JSON(http.StatusOK, response)
}
