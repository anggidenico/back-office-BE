package controllers

import (
	"log"
	"mf-bo-api/lib"
	"mf-bo-api/models"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo"
	"github.com/shopspring/decimal"
)

func CreateAllocInstrumentController(c echo.Context) error {
	var err error
	params := make(map[string]string)
	params["rec_created_by"] = lib.UserIDStr
	params["rec_created_date"] = time.Now().Format(lib.TIMESTAMPFORMAT)

	productKey := c.FormValue("product_key")
	if productKey == "" {
		log.Println(err)
		return lib.CustomError(http.StatusBadRequest, "Missing product_key", "Missing product_key")
	}
	periodeKey := c.FormValue("periode_key")
	if periodeKey == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing periode_key", "Missing periode_key")
	}
	instrumentKey := c.FormValue("instrument_key")
	if instrumentKey == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing instrument_key", "Missing instrument_key")
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
		params["rec_order"] = strconv.Itoa(value)
	} else {
		params["rec_order"] = "0"
	}

	// params["rec_order"] = recOrder
	params["product_key"] = productKey
	params["periode_key"] = periodeKey
	params["instrument_key"] = instrumentKey
	params["instrument_value"] = instrumentValue
	params["rec_status"] = "1"

	// Check for duplicate records
	duplicate, key, err := models.CheckDuplicateAllocInstrument(productKey, periodeKey, instrumentKey)
	if err != nil {
		log.Println("Error checking for duplicates:", err)
		return lib.CustomError(http.StatusInternalServerError, "Error checking for duplicates", "Error checking for duplicates")
	}

	log.Println("Duplicate:", duplicate)
	log.Println("Key:", key)

	// Jika duplikasi ditemukan, perbarui data yang sudah ada
	if duplicate {
		log.Println("Duplicate data found.")
		// Cek apakah data yang sudah ada masih aktif atau sudah dihapus
		existingDataStatus, err := models.GetAllocInstrumentStatusByKey(key)
		if err != nil {
			log.Println("Error getting existing data status:", err)
			return lib.CustomError(http.StatusInternalServerError, "Error getting existing data status", "Error getting existing data status")
		}

		// Jika data sudah dihapus (rec_status = 0), perbarui statusnya menjadi aktif (rec_status = 1)
		if existingDataStatus == 0 {
			log.Println("Existing data is deleted. Recreating data.")

			// Set status menjadi aktif (rec_status = 1)
			params["rec_status"] = "1"
			// Update data dengan status baru dan nilai-nilai yang baru
			status, err := models.UpdateAllocInstrument(key, params)
			if err != nil {
				log.Println("Error updating data:", err)
				return lib.CustomError(status, "Error updating data", "Error updating data")
			}
			return c.JSON(http.StatusOK, lib.Response{
				Status: lib.Status{
					Code:          http.StatusOK,
					MessageServer: "OK",
					MessageClient: "OK",
				},
				Data: "Data updated successfully",
			})
		} else {
			// Jika data masih aktif, kembalikan respons kesalahan duplikasi
			log.Println("Existing data is still active. Duplicate data error.")
			return lib.CustomError(http.StatusBadRequest, "Duplicate data. Unable to input data.", "Duplicate data. Unable to input data.")
		}
	} else {
		// Jika tidak ada duplikasi, buat data baru
		status, err := models.CreateAllocInstrument(params)
		if err != nil {
			log.Println("Error create data:", err)
			return lib.CustomError(status, "Duplicate data. Unable to input data.", "Duplicate data. Unable to input data.")
		}
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
	params := make(map[string]string)
	params["rec_modified_by"] = lib.UserIDStr
	params["rec_modified_date"] = time.Now().Format(lib.TIMESTAMPFORMAT)

	allocInstrumentKey := c.FormValue("alloc_instrument_key")
	if allocInstrumentKey != "" {
		return lib.CustomError(http.StatusBadRequest, "Missing alloc_instrument_key", "Missing alloc_instrument_key")
	}
	productKey := c.FormValue("product_key")
	if productKey != "" {
		return lib.CustomError(http.StatusBadRequest, "Missing product_key", "Missing product_key")
	}
	periodeKey := c.FormValue("periode_key")
	if periodeKey != "" {
		return lib.CustomError(http.StatusBadRequest, "Missing periode_key", "Missing periode_key")
	}
	instrumentKey := c.FormValue("instrument_key")
	if instrumentKey != "" {
		return lib.CustomError(http.StatusBadRequest, "Missing instrument_key", "Missing instrument_key")
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
		params["rec_order"] = strconv.Itoa(value)
	} else {
		params["rec_order"] = "0"
	}

	// params["rec_order"] = recOrder
	params["alloc_instrument_key"] = allocInstrumentKey
	params["product_key"] = productKey
	params["periode_key"] = periodeKey
	params["instrument_key"] = instrumentKey
	params["instrument_value"] = instrumentValue
	params["rec_status"] = "1"

	// Check for duplicate records
	duplicate, key, err := models.CheckDuplicateAllocInstrument(productKey, periodeKey, instrumentKey)
	if err != nil {
		log.Println("Error checking for duplicates:", err)
		return lib.CustomError(http.StatusInternalServerError, "Error checking for duplicates", "Error checking for duplicates")
	}
	if duplicate {
		log.Println("Duplicate data found.")
		// Cek apakah data yang sudah ada masih aktif atau sudah dihapus
		_, err := models.GetAllocInstrumentStatusByKey(key)
		if err != nil {
			log.Println("Error getting existing data status:", err)
			return lib.CustomError(http.StatusBadRequest, "Duplicate data. Unable to input data.", "Duplicate data. Unable to input data.")
		}

		if key != allocInstrumentKey {
			return lib.CustomError(http.StatusBadRequest, "Duplicate data", "Duplicate data")
		}

	}
	status, err = models.UpdateAllocInstrument(allocInstrumentKey, params)
	if err != nil {
		return lib.CustomError(status, "Duplicate data. Unable to input data.", "Duplicate data. Unable to input data.")
	}
	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = "Data updated successfully"

	return c.JSON(http.StatusOK, response)
}

func GetAllocInstrumentController(c echo.Context) error {
	decimal.MarshalJSONWithoutQuotes = true
	var instrument []models.AllocInstrument
	status, err := models.GetAllocInstrumentModels(&instrument)
	if err != nil {
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = instrument
	return c.JSON(http.StatusOK, response)
}
