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

func GetAllocSectorController(c echo.Context) error {
	var sector []models.AllocSector
	status, err := models.GetAllocSectorModels(&sector)
	if err != nil {
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = sector
	log.Printf("Response Data: %+v\n", response.Data)
	return c.JSON(http.StatusOK, response)
}

func GetAllocSectorDetailController(c echo.Context) error {
	allocSectorKey := c.Param("alloc_sector_key")
	if allocSectorKey == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing alloc_sector_key", "Missing alloc_sector_key")
	}
	var alloc models.AllocSectorDetail
	status, err := models.GetAllocSectorDetailModels(&alloc, allocSectorKey)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return lib.CustomError(http.StatusNotFound, "alloc_sector_key not found", "alloc_sector_key not found")
		}
		return lib.CustomError(status, err.Error(), err.Error())
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = alloc
	return c.JSON(http.StatusOK, response)
}

func CreateAllocSectorController(c echo.Context) error {
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
	sectorKey, err := strconv.ParseInt(c.FormValue("sector_key"), 10, 64)
	if err != nil {
		return lib.CustomError(http.StatusBadRequest, "Invalid sector_key", "Invalid sector_key")
	}
	sectorValue := c.FormValue("sector_value")
	if sectorValue == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing sector_value", "Missing sector_value")
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
	params["sector_key"] = sectorKey
	params["sector_value"] = sectorValue
	params["rec_status"] = "1"

	// Check for duplicate records
	if productKey != 0 && periodeKey != 0 && sectorKey != 0 {
		duplicate, key, err := models.CheckDuplicateAllocSector(periodeKey, productKey, sectorKey)
		if err != nil {
			log.Println("Error checking for duplicates:", err)
			return lib.CustomError(http.StatusInternalServerError, "Error checking for duplicates", "Error checking for duplicates")
		}
		log.Println("Duplicate:", duplicate)
		log.Println("Key:", key)
		// Jika duplikasi ditemukan, perbarui data yang sudah ada
		if duplicate {
			status, err := models.UpdateAllocSector(key, params)
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
	status, err = models.CreateAllocSector(params)
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

func UpdateAllocSectorController(c echo.Context) error {
	var err error
	params := make(map[string]interface{})
	params["rec_modified_by"] = lib.UserIDStr
	params["rec_modified_date"] = time.Now().Format(lib.TIMESTAMPFORMAT)

	allocSectorKey, err := strconv.ParseInt(c.FormValue("alloc_sector_key"), 10, 64)
	if err != nil {
		return lib.CustomError(http.StatusBadRequest, "Invalid alloc_sector_key", "Invalid alloc_sector_key")
	}
	productKey, err := strconv.ParseInt(c.FormValue("product_key"), 10, 64)
	if err != nil {
		return lib.CustomError(http.StatusBadRequest, "Invalid product_key", "Invalid product_key")
	}
	periodeKey, err := strconv.ParseInt(c.FormValue("periode_key"), 10, 64)
	if err != nil {
		return lib.CustomError(http.StatusBadRequest, "Invalid periode_key", "Invalid periode_key")
	}
	sectorKey, err := strconv.ParseInt(c.FormValue("sector_key"), 10, 64)
	if err != nil {
		return lib.CustomError(http.StatusBadRequest, "Invalid sector_key", "Invalid sector_key")
	}
	sectorValue := c.FormValue("sector_value")
	if sectorValue == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing sector_value", "Missing sector_value")
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
	params["alloc_sector_key"] = allocSectorKey
	params["product_key"] = productKey
	params["periode_key"] = periodeKey
	params["sector_key"] = sectorKey
	params["sector_value"] = sectorValue
	params["rec_status"] = "1"

	// Check for duplicate records
	if productKey != 0 && periodeKey != 0 && sectorKey != 0 {
		duplicate, key, err := models.CheckDuplicateAllocSector(periodeKey, productKey, sectorKey)
		if err != nil {
			log.Println("Error checking for duplicates:", err)
			return lib.CustomError(http.StatusInternalServerError, "Error checking for duplicates", "Error checking for duplicates")
		}
		log.Println("Duplicate:", duplicate)
		log.Println("Key:", key)
		// Jika duplikasi ditemukan, perbarui data yang sudah ada
		if duplicate {
			status, err := models.UpdateAllocSector(key, params)
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
	status, err = models.CreateAllocSector(params)
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
