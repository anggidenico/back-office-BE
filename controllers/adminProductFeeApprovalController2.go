package controllers

import (
	"mf-bo-api/lib"
	"mf-bo-api/models"
	"net/http"

	"github.com/labstack/echo"
)

func ProductFeeApprovalList(c echo.Context) error {

	data := models.GetProductFeeApprovalList()

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = data
	return c.JSON(http.StatusOK, response)
}

func ProductFeeApprovalDetail(c echo.Context) error {

	rec_pk := c.Param("rec_pk")
	if rec_pk == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing: rec_pk")
	}

	data := models.GetProductFeeApprovalDetail(rec_pk)

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = data
	return c.JSON(http.StatusOK, response)
}
