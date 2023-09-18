package controllers

import (
	"mf-bo-api/lib"
	"mf-bo-api/models"
	"net/http"

	"github.com/labstack/echo"
)

func RevertOAStatus(c echo.Context) error {

	OaRequestyKeys := c.QueryParam("oa_request_key")
	if OaRequestyKeys == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing: oa_request_key", "Missing: oa_request_key")
	}

	err := models.SetRevertOAStatus(OaRequestyKeys)
	if err != nil {
		return lib.CustomError(http.StatusBadRequest, "Cannot Update oa_status", "Cannot Update oa_status")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = nil
	return c.JSON(http.StatusOK, response)
}
