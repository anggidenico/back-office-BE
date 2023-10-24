package controllers

import (
	"mf-bo-api/lib"
	"mf-bo-api/models"
	"net/http"
	"time"

	"github.com/labstack/echo"
)

func ProductPaymentChannelList(c echo.Context) error {

	productKey := c.Param("product_key")
	if productKey == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing product_key")
	}

	data := models.GetProductPaymentChannels(productKey)

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = data

	return c.JSON(http.StatusOK, response)
}

func CreateProductPaymentChannels(c echo.Context) error {

	params := make(map[string]string)
	params["rec_created_by"] = lib.UserIDStr
	params["rec_created_date"] = time.Now().Format(lib.TIMESTAMPFORMAT)
	params["rec_status"] = "1"

	productKey := c.FormValue("product_key")
	if productKey == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing product_key", "Missing product_key")
	}
	params["product_key"] = productKey

	pchannelKey := c.FormValue("pchannel_key")
	if pchannelKey == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing pchannel_key", "Missing pchannel_key")
	}
	params["pchannel_key"] = pchannelKey

	cotTransaction := c.FormValue("cot_transaction")
	if cotTransaction == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing cot_transaction", "Missing cot_transaction")
	}
	params["cot_transaction"] = cotTransaction

	cotSettlement := c.FormValue("cot_settlement")
	if cotSettlement == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing cot_settlement", "Missing cot_settlement")
	}
	params["cot_settlement"] = cotSettlement

	err := models.CreateProductPaymentChannels(params)
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

func UpdateProductPaymentChannels(c echo.Context) error {

	params := make(map[string]string)
	params["rec_modified_by"] = lib.UserIDStr
	params["rec_modified_date"] = time.Now().Format(lib.TIMESTAMPFORMAT)
	params["rec_status"] = "1"

	prodChannelKey := c.FormValue("prod_channel_key")
	if prodChannelKey == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing prod_channel_key", "Missing prod_channel_key")
	}

	productKey := c.FormValue("product_key")
	if productKey == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing product_key", "Missing product_key")
	}
	params["product_key"] = productKey

	pchannelKey := c.FormValue("pchannel_key")
	if pchannelKey == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing pchannel_key", "Missing pchannel_key")
	}
	params["pchannel_key"] = pchannelKey

	cotTransaction := c.FormValue("cot_transaction")
	if cotTransaction == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing cot_transaction", "Missing cot_transaction")
	}
	params["cot_transaction"] = cotTransaction

	cotSettlement := c.FormValue("cot_settlement")
	if cotSettlement == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing cot_settlement", "Missing cot_settlement")
	}
	params["cot_settlement"] = cotSettlement

	err := models.UpdateProductPaymentChannels(params)
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

func DeleteProductPaymentChannels(c echo.Context) error {

	prodChannelKey := c.FormValue("prod_channel_key")
	if prodChannelKey == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing prod_channel_key", "Missing prod_channel_key")
	}

	err := models.DeleteProductPaymentChannels(prodChannelKey)
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
