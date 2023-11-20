package controllers

import (
	"database/sql"
	"errors"
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
	}
	var detail models.AllocSecDetail
	status, err := models.GetAllocSecDetailModels(&detail, allocSecKey)
	// log.Println("Not Found")
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return lib.CustomError(http.StatusNotFound, "alloc_sec_key not found", "alloc_sec_key not found")
		}
		return lib.CustomError(status, err.Error(), err.Error())
	}
	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = detail
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
func UpdateAllocSecController(c echo.Context) error {
	var err error
	params := make(map[string]string)
	params["rec_modified_by"] = lib.UserIDStr
	params["rec_modified_date"] = time.Now().Format(lib.TIMESTAMPFORMAT)

	allocSecKey := c.FormValue("alloc_security_key")
	if allocSecKey == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing alloc_security_key", "Missing alloc_security_key")
	}
	productKey := c.FormValue("product_key")
	if productKey != "" {
		if len(productKey) > 11 {
			return lib.CustomError(http.StatusBadRequest, "product_key must be <= 11 characters", "product_key must be <= 11 characters")
		}
		_, err := strconv.Atoi(productKey)
		if err != nil {
			return lib.CustomError(http.StatusBadRequest, "product key must be number", "product key must be number")
		}

	} else {
		return lib.CustomError(http.StatusBadRequest, "Missing alloc_security_key", "Missing alloc_security_key")
	}
	periodeKey := c.FormValue("periode_key")
	if periodeKey != "" {
		if len(periodeKey) > 11 {
			return lib.CustomError(http.StatusBadRequest, "periode_key must be <= 11 characters", "periode_key must be <= 11 characters")
		}
		_, err := strconv.Atoi(productKey)
		if err != nil {
			return lib.CustomError(http.StatusBadRequest, "periode_key must be number", "periode_key must be number")
		}

	} else {
		return lib.CustomError(http.StatusBadRequest, "Missing periode_key", "Missing periode_key")
	}
	secKey := c.FormValue("sec_key")
	if secKey != "" {
		if len(secKey) > 11 {
			return lib.CustomError(http.StatusBadRequest, "sec_key must be <= 11 characters", "sec_key must be <= 11 characters")
		}
		_, err := strconv.Atoi(secKey)
		if err != nil {
			return lib.CustomError(http.StatusBadRequest, "sec_key must be number", "sec_key must be number")
		}

	} else {
		return lib.CustomError(http.StatusBadRequest, "Missing sec_key", "Missing sec_key")
	}
	secValue := c.FormValue("security_value")
	if secValue != "" {
		if len(secValue) > 11 {
			return lib.CustomError(http.StatusBadRequest, "security_value must be <= 11 characters", "security_value must be <= 11 characters")
		}
		_, err := strconv.Atoi(secValue)
		if err != nil {
			return lib.CustomError(http.StatusBadRequest, "security_value must be number", "security_value must be number")
		}
	} else {
		if secValue == "" {
			return lib.CustomError(http.StatusBadRequest, "Missing security_value", "Missing security_value")
		}
	}
	recOrder := c.FormValue("rec_order")
	if recOrder != "" {
		if len(recOrder) > 11 {
			return lib.CustomError(http.StatusBadRequest, "rec_order should be exactly 11 characters", "rec_order be exactly 11 characters")
		}
		value, err := strconv.Atoi(recOrder)
		if err != nil {
			return lib.CustomError(http.StatusBadRequest, "rec_order should be a number", "rec_order should be a number")
		}
		params["rec_order"] = strconv.Itoa(value)
	}
	params["rec_order"] = recOrder
	params["product_key"] = productKey
	params["periode_key"] = periodeKey
	params["sec_key"] = secKey
	params["security_value"] = secValue
	params["rec_status"] = "1"

	status, err = models.UpdateAllocSec(allocSecKey, params)
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

func DeleteAllocSecController(c echo.Context) error {
	params := make(map[string]string)
	params["rec_status"] = "0"
	params["rec_deleted_date"] = time.Now().Format(lib.TIMESTAMPFORMAT)
	params["rec_deleted_by"] = lib.UserIDStr

	allocSecKey := c.FormValue("alloc_security_key")
	if allocSecKey == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing alloc_security_key", "Missing alloc_security_key")
	}

	status, err := models.DeleteAllocSec(allocSecKey, params)
	if err != nil {
		return lib.CustomError(status, err.Error(), err.Error())
	}
	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "Berhasil hapus Portfolio Instrument!"
	response.Data = ""
	return c.JSON(http.StatusOK, response)
}
