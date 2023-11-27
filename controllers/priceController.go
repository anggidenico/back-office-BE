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

func CreatePriceController(c echo.Context) error {
	var err error
	params := make(map[string]string)
	params["rec_created_by"] = lib.UserIDStr
	params["rec_created_date"] = time.Now().Format(lib.TIMESTAMPFORMAT)

	benchmarkKey := c.FormValue("benchmark_key")
	if benchmarkKey == "" {
		return lib.CustomError(http.StatusBadRequest, "benchmark_key can not be blank", "benchmark_key can not be blank")
	}
	priceType := c.FormValue("price_type")
	if priceType == "" {
		return lib.CustomError(http.StatusBadRequest, "price_type can not be blank", "price_type can not be blank")
	}
	priceDate := c.FormValue("price_date")
	if priceDate == "" {
		return lib.CustomError(http.StatusBadRequest, "price_date can not be blank", "price_date can not be blank")
	}
	expectedDateFormat := "2006-01-02"
	_, err = time.Parse(expectedDateFormat, priceDate)
	if err != nil {
		return lib.CustomError(http.StatusBadRequest, "price_date should be a valid date in the format "+expectedDateFormat, "price_date should be a valid date in the format "+expectedDateFormat)
	}
	priceValue := c.FormValue("price_value")
	if priceValue == "" {
		return lib.CustomError(http.StatusBadRequest, "price_value can not be blank", "price_value can not be blank")
	}
	priceRemarks := c.FormValue("price_remarks")
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
	} else {
		params["rec_order"] = "0"
	}
	params["benchmark_key"] = benchmarkKey
	params["price_type"] = priceType
	params["price_date"] = priceDate
	params["price_value"] = priceValue
	params["price_remarks"] = priceRemarks
	params["rec_status"] = "1"

	// Check for duplicate records
	duplicate, key, err := models.CheckDuplicateFfsPrice(params["benchmark_key"], params["price_type"], params["price_date"])
	if err != nil {
		log.Println("Error checking for duplicates:", err)
		return lib.CustomError(http.StatusInternalServerError, "Error checking for duplicates", "Error checking for duplicates")
	}
	log.Println("Duplicate:", duplicate)
	log.Println("Key:", key)
	// Jika duplikasi ditemukan, perbarui data yang sudah ada
	if duplicate {
		status, err := models.UpdatePrice(key, params)
		if err != nil {
			log.Println("Failed to update data:", err)
			return lib.CustomError(status, "Failed to update data", "Failed to update data")
		}
		return c.JSON(http.StatusOK, lib.Response{
			Status: lib.Status{
				Code:          http.StatusOK,
				MessageServer: "OK",
				MessageClient: "OK",
			},
			Data: "Data updated successfully",
		})
	}

	// Jika tidak ada duplikasi, buat data baru
	status, err = models.CreateFfsPrice(params)
	if err != nil {
		return lib.CustomError(status, "Failed input data", "Failed input data")
	}

	return c.JSON(http.StatusOK, lib.Response{
		Status: lib.Status{
			Code:          http.StatusOK,
			MessageServer: "OK",
			MessageClient: "OK",
		},
		Data: "Data created successfully",
	})
}
func UpdatePriceController(c echo.Context) error {
	var err error
	params := make(map[string]string)
	params["rec_modified_by"] = lib.UserIDStr
	params["rec_modified_date"] = time.Now().Format(lib.TIMESTAMPFORMAT)

	priceKey := c.FormValue("price_key")
	if priceKey == "" {
		return lib.CustomError(http.StatusBadRequest, "price_key can not be blank", "price_key can not be blank")
	}
	benchmarkKey := c.FormValue("benchmark_key")
	if benchmarkKey == "" {
		return lib.CustomError(http.StatusBadRequest, "benchmark_key can not be blank", "benchmark_key can not be blank")
	}
	priceType := c.FormValue("price_type")
	if priceType == "" {
		return lib.CustomError(http.StatusBadRequest, "price_type can not be blank", "price_type can not be blank")
	}
	priceDate := c.FormValue("price_date")
	if priceDate == "" {
		return lib.CustomError(http.StatusBadRequest, "price_date can not be blank", "price_date can not be blank")
	}
	expectedDateFormat := "2006-01-02"
	_, err = time.Parse(expectedDateFormat, priceDate)
	if err != nil {
		return lib.CustomError(http.StatusBadRequest, "price_date should be a valid date in the format "+expectedDateFormat, "price_date should be a valid date in the format "+expectedDateFormat)
	}
	priceValue := c.FormValue("price_value")
	if priceValue == "" {
		return lib.CustomError(http.StatusBadRequest, "price_value can not be blank", "price_value can not be blank")
	}
	priceRemarks := c.FormValue("price_remarks")
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
	} else {
		params["rec_order"] = "0"
	}
	params["price_key"] = priceKey
	params["benchmark_key"] = benchmarkKey
	params["price_type"] = priceType
	params["price_date"] = priceDate
	params["price_value"] = priceValue
	params["price_remarks"] = priceRemarks
	params["rec_status"] = "1"

	status, err = models.UpdatePrice(priceKey, params)
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
