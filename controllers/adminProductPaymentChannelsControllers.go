package controllers

import (
	"mf-bo-api/lib"
	"mf-bo-api/models"
	"net/http"

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

	productKey := c.FormValue("product_key")
	if productKey == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing product_key", "Missing product_key")
	}

	pchannelKey := c.FormValue("pchannel_key")
	if pchannelKey == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing pchannel_key", "Missing pchannel_key")
	}

	cotTransaction := c.FormValue("cot_transaction")
	if cotTransaction == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing cot_transaction", "Missing cot_transaction")
	}

	cotSettlement := c.FormValue("cot_settlement")
	if cotSettlement == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing cot_settlement", "Missing cot_settlement")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = nil

	return c.JSON(http.StatusOK, response)
}