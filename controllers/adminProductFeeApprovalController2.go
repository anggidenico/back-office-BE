package controllers

import (
	"mf-bo-api/lib"
	"mf-bo-api/models"
	"net/http"
	"time"

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

func ProductFeeApprovalAction(c echo.Context) error {
	params := make(map[string]string)
	params["rec_by"] = lib.UserIDStr
	params["rec_date"] = time.Now().Format(lib.TIMESTAMPFORMAT)

	RecPK := c.FormValue("rec_pk")
	if RecPK == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing: rec_pk", "Missing: rec_pk")
	}
	params["rec_pk"] = RecPK

	Approval := c.FormValue("approval")
	if Approval == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing: approval", "Missing: approval")
	}
	if Approval == "true" {
		params["approval"] = "1"
	} else {
		params["approval"] = "0"
	}

	Reason := c.FormValue("reason")
	if Approval == "false" && Reason == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing: reason", "Missing: reason")
	}
	params["reason"] = Reason

	err := models.ProductFeeApprovalAction(params)
	if err != nil {
		return lib.CustomError(http.StatusInternalServerError, err.Error())
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = nil
	return c.JSON(http.StatusOK, response)
}
