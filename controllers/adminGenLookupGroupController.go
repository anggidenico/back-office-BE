package controllers

import (
	"mf-bo-api/lib"
	"mf-bo-api/models"
	"net/http"

	"github.com/labstack/echo"
)

func AdminGetListDropdownLookupGroup(c echo.Context) error {

	var err error
	var status int

	var lookupgroup []models.ListDropdownLookupGroup

	status, err = models.AdminGetListDropdownLookupGroup(&lookupgroup)

	if err != nil {
		// log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if len(lookupgroup) < 1 {
		// log.Error("Lookup Group Type not found")
		return lib.CustomError(http.StatusNotFound, "Lookup Group Type not found", "Lookup Group Type not found")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = lookupgroup

	return c.JSON(http.StatusOK, response)
}
