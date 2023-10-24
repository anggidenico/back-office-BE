package controllers

import (
	"mf-bo-api/lib"
	"mf-bo-api/models"
	"net/http"
	"time"

	"github.com/labstack/echo"
)

func GetmsPaymentChannelController(c echo.Context) error {

	result := models.GetPaymentChannelModels()

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = result
	return c.JSON(http.StatusOK, response)
}

func GetMsPaymentDetailController(c echo.Context) error {
	pChannelKey := c.Param("pchannel_key")
	if pChannelKey == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing payment channel key", "Missing payment channel key")
	}
	result := models.GetDetailPaymentChannelModels(pChannelKey)
	// log.Println("Not Found")

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = result
	return c.JSON(http.StatusOK, response)
}

func DeleteMsPaymentChannelController(c echo.Context) error {
	params := make(map[string]string)
	params["rec_status"] = "0"
	params["rec_deleted_date"] = time.Now().Format(lib.TIMESTAMPFORMAT)
	params["rec_deleted_by"] = lib.UserIDStr

	pChannelKey := c.FormValue("pchannel_key")
	if pChannelKey == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing pchannelKey", "Missing pchannelKey")
	}

	err := models.DeleteMsPaymentChannel(pChannelKey, params)
	if err != nil {
		return lib.CustomError(http.StatusInternalServerError, err.Error(), err.Error())
	}
	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "Berhasil hapus Master Payment Channel!"
	response.Data = ""
	return c.JSON(http.StatusOK, response)
}
