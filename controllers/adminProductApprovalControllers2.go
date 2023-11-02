package controllers

import (
	"mf-bo-api/lib"
	"mf-bo-api/models"
	"net/http"
	"time"

	"github.com/labstack/echo"
)

func ProductApprovalAction(c echo.Context) error {

	params := make(map[string]string)
	params["rec_created_by"] = lib.UserIDStr
	params["rec_created_date"] = time.Now().Format(lib.TIMESTAMPFORMAT)

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

	err := models.ApprovalAction(params)
	if err != nil {
		return lib.CustomError(http.StatusInternalServerError, err.Error(), err.Error())
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = nil
	return c.JSON(http.StatusOK, response)
}
