package controllers

import (
	"log"
	"mf-bo-api/lib"
	"mf-bo-api/models"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo"
)

func CreateAllocInstrumentController(c echo.Context) error {
	var err error
	params := make(map[string]interface{})
	params["rec_created_by"] = lib.UserIDStr
	params["rec_created_date"] = time.Now().Format(lib.TIMESTAMPFORMAT)

	productKey, err := strconv.ParseInt(c.FormValue("product_key"), 10, 64)
	if err != nil {
		return lib.CustomError(http.StatusBadRequest, "Invalid product_key", "Invalid product_key")
	}
	periodeKey, err := strconv.ParseInt(c.FormValue("periode_key"), 10, 64)
	if err != nil {
		return lib.CustomError(http.StatusBadRequest, "Invalid periode_key", "Invalid periode_key")
	}
	instrumentKey, err := strconv.ParseInt(c.FormValue("instrument_key"), 10, 64)
	if err != nil {
		return lib.CustomError(http.StatusBadRequest, "Invalid instrument_key", "Invalid instrument_key")
	}
	instrumentValue := c.FormValue("instrument_value")
	if instrumentValue == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing instrument_value", "Missing instrument_value")
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
		params["rec_order"] = value
	} else {
		params["rec_order"] = 0
	}

	// params["rec_order"] = recOrder
	params["product_key"] = productKey
	params["periode_key"] = periodeKey
	params["instrument_key"] = instrumentKey
	params["instrument_value"] = instrumentValue
	params["rec_status"] = "1"

	// Check for duplicate records
	if productKey != 0 && periodeKey != 0 && instrumentKey != 0 {
		duplicate, key, err := models.CheckDuplicateAllocInstrument(periodeKey, productKey, instrumentKey)
		if err != nil {
			log.Println("Error checking for duplicates:", err)
			return lib.CustomError(http.StatusInternalServerError, "Error checking for duplicates", "Error checking for duplicates")
		}
		log.Println("Duplicate:", duplicate)
		log.Println("Key:", key)
		// Jika duplikasi ditemukan, perbarui data yang sudah ada
		if duplicate {
			status, err := models.UpdateAllocInstrument(key, params)
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
	}
	// Jika tidak ada duplikasi, buat data baru
	status, err = models.CreateAllocInstrument(params)
	if err != nil {
		log.Println(err)
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

func UpdateAllocInstrumentController(c echo.Context) error {
	var err error
	params := make(map[string]interface{})
	params["rec_modified_by"] = lib.UserIDStr
	params["rec_modified_date"] = time.Now().Format(lib.TIMESTAMPFORMAT)

	allocInstrumentKey, err := strconv.ParseInt(c.FormValue("alloc_instrument_key"), 10, 64)
	if err != nil {
		return lib.CustomError(http.StatusBadRequest, "Invalid alloc_instrument_key", "Invalid alloc_instrument_key")
	}
	productKey, err := strconv.ParseInt(c.FormValue("product_key"), 10, 64)
	if err != nil {
		return lib.CustomError(http.StatusBadRequest, "Invalid product_key", "Invalid product_key")
	}
	periodeKey, err := strconv.ParseInt(c.FormValue("periode_key"), 10, 64)
	if err != nil {
		return lib.CustomError(http.StatusBadRequest, "Invalid periode_key", "Invalid periode_key")
	}
	instrumentKey, err := strconv.ParseInt(c.FormValue("instrument_key"), 10, 64)
	if err != nil {
		return lib.CustomError(http.StatusBadRequest, "Invalid instrument_key", "Invalid instrument_key")
	}
	instrumentValue := c.FormValue("instrument_value")
	if instrumentValue == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing instrument_value", "Missing instrument_value")
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
		params["rec_order"] = value
	} else {
		params["rec_order"] = 0
	}

	// params["rec_order"] = recOrder
	params["alloc_instrument_key"] = allocInstrumentKey
	params["product_key"] = productKey
	params["periode_key"] = periodeKey
	params["instrument_key"] = instrumentKey
	params["instrument_value"] = instrumentValue
	params["rec_status"] = "1"

	// Check for duplicate records
	if productKey != 0 && periodeKey != 0 && instrumentKey != 0 {
		duplicate, key, err := models.CheckDuplicateAllocInstrument(periodeKey, productKey, instrumentKey)
		if err != nil {
			log.Println("Error checking for duplicates:", err)
			return lib.CustomError(http.StatusInternalServerError, "Error checking for duplicates", "Error checking for duplicates")
		}
		log.Println("Duplicate:", duplicate)
		log.Println("Key:", key)
		// Jika duplikasi ditemukan, perbarui data yang sudah ada
		if duplicate {
			status, err := models.UpdateAllocInstrument(key, params)
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
	}
	return c.JSON(http.StatusOK, lib.Response{
		Status: lib.Status{
			Code:          http.StatusOK,
			MessageServer: "OK",
			MessageClient: "OK",
		},
		Data: "No action taken",
	})
}
