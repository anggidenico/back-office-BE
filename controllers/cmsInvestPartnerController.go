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

func GetInvestPurposeController(c echo.Context) error {
	var invest []models.InvestPurpose
	status, err := models.GetInvestPurposeModels(&invest)
	if err != nil {
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = invest
	log.Printf("Response Data: %+v\n", response.Data)
	return c.JSON(http.StatusOK, response)
}

func GetInvestPartnerController(c echo.Context) error {
	var invest []models.InvestPartner
	status, err := models.GetInvestPartnerModels(&invest)
	if err != nil {
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = invest
	log.Printf("Response Data: %+v\n", response.Data)
	return c.JSON(http.StatusOK, response)
}

func GetInvestPartnerDetailController(c echo.Context) error {
	partnerKey := c.Param("invest_partner_key")
	if partnerKey == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing invest_partner_key", "Missing invest_partner_key")
	}
	var invest models.InvestPartner
	status, err := models.GetInvestPartnerDetailModels(&invest, partnerKey)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return lib.CustomError(http.StatusNotFound, "invest_partner_key not found", "invest_partner_key not found")
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
