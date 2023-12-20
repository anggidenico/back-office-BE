package controllers

import (
	"database/sql"
	"errors"
	"mf-bo-api/lib"
	"mf-bo-api/models"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/shopspring/decimal"
)

func GetAutoInvestRegisController(c echo.Context) error {
	decimal.MarshalJSONWithoutQuotes = true
	var value []models.AutoInvestRegistration
	status, err := models.GetAutoInvestRegistrationModels(&value)
	if err != nil {
		// log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	var response lib.Response
	if len(value) < 1 {
		response.Status.Code = http.StatusOK
		response.Status.MessageServer = "OK"
		response.Status.MessageClient = "OK"
		response.Data = []models.RiskProfile{}
	} else {
		response.Status.Code = http.StatusOK
		response.Status.MessageServer = "OK"
		response.Status.MessageClient = "OK"
		response.Data = value
	}
	return c.JSON(http.StatusOK, response)
}

func GetAutoInvestRegisDetailController(c echo.Context) error {
	autoInvestKey := c.Param("autoinvest_key")
	if autoInvestKey == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing autoinvest_key", "Missing autoinvest_key")
	}
	var invest models.AutoInvestRegistration
	status, err := models.GetAutoInvestRegisDetailModels(&invest, autoInvestKey)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return lib.CustomError(http.StatusNotFound, "autoinvest_key not found", "autoinvest_key not found")
		}
		return lib.CustomError(status, err.Error(), err.Error())
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = invest
	return c.JSON(http.StatusOK, response)
}

func DeleteAutoInvestRegisController(c echo.Context) error {
	params := make(map[string]string)
	params["rec_status"] = "0"
	params["rec_deleted_date"] = time.Now().Format(lib.TIMESTAMPFORMAT)
	params["rec_deleted_by"] = lib.UserIDStr

	autoInvestKey := c.FormValue("auto_invest_key")
	if autoInvestKey == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing auto_invest_key", "Missing sec_key")
	}

	status, err := models.DeleteAutoInvestRegis(autoInvestKey, params)
	if err != nil {
		return lib.CustomError(status, err.Error(), err.Error())
	}
	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "Berhasil hapus Message Posting!"
	response.Data = ""
	return c.JSON(http.StatusOK, response)
}
