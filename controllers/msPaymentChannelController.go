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

func CreateMsPaymentChannelController(c echo.Context) error {
	var err error
	params := make(map[string]string)
	params["rec_created_by"] = lib.UserIDStr
	params["rec_created_date"] = time.Now().Format(lib.TIMESTAMPFORMAT)

	pChannelCode := c.FormValue("pchannel_code")
	if pChannelCode == "" {
		return lib.CustomError(http.StatusBadRequest, "pchannel_code can not be blank", "pchannel_code can not be blank")
	}
	pchannelName := c.FormValue("pchannel_name")
	if pchannelName == "" {
		return lib.CustomError(http.StatusBadRequest, "pchannel_name can not be blank", "pchannel_name can not be blank")
	}
	minNominalTrx := c.FormValue("min_nominal_trx")
	if minNominalTrx == "" {
		return lib.CustomError(http.StatusBadRequest, "min_nominal_trx can not be blank", "min_nominal_trx can not be blank")
	}
	feeValue := c.FormValue("fee_value")
	if feeValue == "" {
		return lib.CustomError(http.StatusBadRequest, "fee_value can not be blank", "fee_value can not be blank")
	}
	hasMinMax := c.FormValue("has_min_max")
	if hasMinMax == "" {
		return lib.CustomError(http.StatusBadRequest, "has_min_max can not be blank", "has_min_max can not be blank")
	}
	settleChannel := c.FormValue("settle_channel")
	if settleChannel == "" {
		return lib.CustomError(http.StatusBadRequest, "settle_channel can not be blank", "settle_channel can not be blank")
	}
	settlePaymentMethod := c.FormValue("settle_payment_method")
	if settlePaymentMethod == "" {
		return lib.CustomError(http.StatusBadRequest, "settle_payment_method can not be blank", "settle_payment_method can not be blank")
	}
	valueType := c.FormValue("value_type")
	if valueType == "" {
		return lib.CustomError(http.StatusBadRequest, "settle_channel can not be blank", "settle_channel can not be blank")
	}

	params["pchannel_code"] = pChannelCode
	params["pchannel_name"] = pchannelName
	params["min_nominal_trx"] = minNominalTrx
	params["pchannel_name"] = pchannelName
	params["fee_value"] = feeValue
	params["has_min_max"] = hasMinMax
	params["settle_channel"] = settleChannel
	params["settle_payment_method"] = settlePaymentMethod
	params["value_type"] = valueType
	params["rec_status"] = "1"
	status, err = models.CreateMsPaymentChannel(params)
	if err != nil {
		return lib.CustomError(status, err.Error(), "Failed input data")
	}
	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = ""

	return c.JSON(http.StatusOK, response)
}

func UpdateMsPaymentChannelController(c echo.Context) error {
	var err error
	params := make(map[string]string)
	params["rec_modified_by"] = lib.UserIDStr
	params["rec_modified_date"] = time.Now().Format(lib.TIMESTAMPFORMAT)

	pChannelKey := c.FormValue("pchannel_key")
	if pChannelKey == "" {
		return lib.CustomError(http.StatusBadRequest, "pchannel_key can not be blank", "pchannel_key can not be blank")
	}
	pChannelCode := c.FormValue("pchannel_code")
	if pChannelCode == "" {
		return lib.CustomError(http.StatusBadRequest, "pchannel_code can not be blank", "pchannel_code can not be blank")
	}
	pchannelName := c.FormValue("pchannel_name")
	if pchannelName == "" {
		return lib.CustomError(http.StatusBadRequest, "pchannel_name can not be blank", "pchannel_name can not be blank")
	}
	minNominalTrx := c.FormValue("min_nominal_trx")
	if minNominalTrx == "" {
		return lib.CustomError(http.StatusBadRequest, "min_nominal_trx can not be blank", "min_nominal_trx can not be blank")
	}
	feeValue := c.FormValue("fee_value")
	if feeValue == "" {
		return lib.CustomError(http.StatusBadRequest, "fee_value can not be blank", "fee_value can not be blank")
	}
	hasMinMax := c.FormValue("has_min_max")
	if hasMinMax == "" {
		return lib.CustomError(http.StatusBadRequest, "has_min_max can not be blank", "has_min_max can not be blank")
	}
	settleChannel := c.FormValue("settle_channel")
	if settleChannel == "" {
		return lib.CustomError(http.StatusBadRequest, "settle_channel can not be blank", "settle_channel can not be blank")
	}
	settlePaymentMethod := c.FormValue("settle_payment_method")
	if settlePaymentMethod == "" {
		return lib.CustomError(http.StatusBadRequest, "settle_payment_method can not be blank", "settle_payment_method can not be blank")
	}
	valueType := c.FormValue("value_type")
	if valueType == "" {
		return lib.CustomError(http.StatusBadRequest, "settle_channel can not be blank", "settle_channel can not be blank")
	}
	params["pchannel_key"] = pChannelKey
	params["pchannel_code"] = pChannelCode
	params["pchannel_name"] = pchannelName
	params["min_nominal_trx"] = minNominalTrx
	params["pchannel_name"] = pchannelName
	params["fee_value"] = feeValue
	params["has_min_max"] = hasMinMax
	params["settle_channel"] = settleChannel
	params["settle_payment_method"] = settlePaymentMethod
	params["value_type"] = valueType
	params["rec_status"] = "1"

	err = models.UpdateMsPaymentChannel(pChannelKey, params)
	if err != nil {
		return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed input data")
	}
	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = ""

	return c.JSON(http.StatusOK, response)
}
