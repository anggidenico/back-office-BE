package controllers

import (
	"database/sql"
	"mf-bo-api/lib"
	"mf-bo-api/models"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo"
)

func GetAllocSecController(c echo.Context) error {
	var allocsec []models.AllocSecurity
	status, err := models.GetAllocSecModels(&allocsec)
	if err != nil {
		// log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	var response lib.Response
	if len(allocsec) < 1 {
		response.Status.Code = http.StatusOK
		response.Status.MessageServer = "OK"
		response.Status.MessageClient = "OK"
		response.Data = []models.RiskProfile{}
	} else {
		response.Status.Code = http.StatusOK
		response.Status.MessageServer = "OK"
		response.Status.MessageClient = "OK"
		response.Data = allocsec
	}
	return c.JSON(http.StatusOK, response)
}
func GetAllocSecDetailController(c echo.Context) error {
	allocSecKey := c.Param("alloc_security_key")
	if allocSecKey == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing alloc_security_key", "Missing alloc_security_key")
	} else {
		_, err := strconv.ParseUint(allocSecKey, 10, 64)
		if err != sql.ErrNoRows {
			// log.Error("Wrong input for parameter: country_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: alloc_security_key", "Wrong input for parameter: alloc_security_key")
		}
	}
	var detailendpoint models.AllocSecDetail
	status, err := models.GetAllocSecDetailModels(&detailendpoint, allocSecKey)
	// log.Println("Not Found")
	if err != nil {
		// log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = detailendpoint
	return c.JSON(http.StatusOK, response)
}
func CreateAllocSecController(c echo.Context) error {
	var err error
	params := make(map[string]string)
	params["rec_created_by"] = lib.UserIDStr
	params["rec_created_date"] = time.Now().Format(lib.TIMESTAMPFORMAT)

	productKey := c.FormValue("product_key")
	if productKey == "" {
		return lib.CustomError(http.StatusBadRequest, "product_key can not be blank", "product_key can not be blank")
	}
	periodeKey := c.FormValue("periode_key")
	if periodeKey == "" {
		return lib.CustomError(http.StatusBadRequest, "periode_key can not be blank", "periode_key can not be blank")
	}
	secKey := c.FormValue("sec_key")
	if secKey == "" {
		return lib.CustomError(http.StatusBadRequest, "sec_key can not be blank", "sec_key can not be blank")
	}
	secValue := c.FormValue("security_value")
	if secValue == "" {
		return lib.CustomError(http.StatusBadRequest, "security_value can not be blank", "security_value can not be blank")
	}
	recOrder := c.FormValue("rec_order")
	if recOrder == "" {
		return lib.CustomError(http.StatusBadRequest, "rec_order can not be blank", "rec_order can not be blank")
	}
	params["product_key"] = productKey
	params["periode_key"] = periodeKey
	params["sec_key"] = secKey
	params["security_value"] = secValue
	params["rec_order"] = recOrder
	params["rec_status"] = "1"

	status, err = models.CreateAllocSec(params)
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
func UpdateAllocSec(c echo.Context) error {
	// var err error
	params := make(map[string]string)
	params["rec_modified_by"] = lib.UserIDStr
	params["rec_modified_date"] = time.Now().Format(lib.TIMESTAMPFORMAT)

	allocSecKey := c.FormValue("alloc_security_key")
	if allocSecKey == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing alloc_security_key", "Missing alloc_security_key")
	}
	params["date_opened"] = allocSecKey
	params["date_closed"] = "2"
	params["rec_status"] = "1"

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = ""

	return c.JSON(http.StatusOK, response)
}
