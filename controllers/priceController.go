package controllers

import (
	"database/sql"
	"errors"
	"log"
	"mf-bo-api/lib"
	"mf-bo-api/models"
	"net/http"

	"github.com/labstack/echo"
)

func GetPriceListController(c echo.Context) error {
	var price []models.PriceList
	status, err := models.GetPriceListModels(&price)
	if err != nil {
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = price
	log.Printf("Response Data: %+v\n", response.Data)
	return c.JSON(http.StatusOK, response)
}

func GetPriceDetailController(c echo.Context) error {
	priceKey := c.Param("price_key")
	if priceKey == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing price_key", "Missing price_key")
	}
	var price models.PriceList
	status, err := models.GetPriceDetailModels(&price, priceKey)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return lib.CustomError(http.StatusNotFound, "price_key not found", "price_key not found")
		}
		return lib.CustomError(status, err.Error(), err.Error())
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = price
	return c.JSON(http.StatusOK, response)
}
