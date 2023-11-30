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

func GetSecuritiesSectorController(c echo.Context) error {
	var value []models.SecuritiesSector
	status, err := models.GetSecuritiesSectorModels(&value)
	if err != nil {
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = value
	log.Printf("Response Data: %+v\n", response.Data)
	return c.JSON(http.StatusOK, response)
}

func GetSecSectorDetailController(c echo.Context) error {
	sectorKey := c.Param("sector_key")
	if sectorKey == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing sector_key", "Missing sector_key")
	}
	var price models.SecuritiesSector
	status, err := models.GetSecuritiesSectorDetailModels(&price, sectorKey)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return lib.CustomError(http.StatusNotFound, "sector_key not found", "sector_key not found")
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
