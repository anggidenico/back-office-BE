package controllers

import (
	"database/sql"
	"errors"
	"log"
	"mf-bo-api/lib"
	"mf-bo-api/models"
	"net/http"
	"strconv"
	"time"

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
func DeleteInvestPartnerController(c echo.Context) error {
	params := make(map[string]string)
	dateLayout := "2006-01-02 15:04:05"
	params["rec_status"] = "0"
	params["rec_deleted_date"] = time.Now().Format(dateLayout)
	params["rec_deleted_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

	investPartnerKey := c.FormValue("invest_partner_key")
	if investPartnerKey == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing invest_partner_key", "Missing invest_partner_key")
	}

	status, err := models.DeleteInvestPartnerModels(investPartnerKey, params)
	if err != nil {
		return lib.CustomError(status, err.Error(), err.Error())
	}
	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "Berhasil hapus CMS Invest Partner!"
	response.Data = ""
	return c.JSON(http.StatusOK, response)
}
