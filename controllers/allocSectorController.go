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
	"github.com/shopspring/decimal"
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
	// log.Printf("Response Data: %+v\n", response.Data)
	return c.JSON(http.StatusOK, response)
}
func GetSectorSecuController(c echo.Context) error {
	decimal.MarshalJSONWithoutQuotes = true
	var sector []models.SectorKey
	status, err := models.GetSectorSecuModels(&sector)
	if err != nil {
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = sector
	// log.Printf("Response Data: %+v\n", response.Data)
	return c.JSON(http.StatusOK, response)
}
func GetAllocSectorDetailController(c echo.Context) error {
	decimal.MarshalJSONWithoutQuotes = true
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

	allocSectorKey := c.FormValue("alloc_sector_key")
	if allocSectorKey == "" {
		return lib.CustomError(http.StatusBadRequest, "Invalid product_key", "Invalid product_key")
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

	duplicate, key, err := models.CheckDuplicateAllocSector(periodeKey, productKey, sectorKey)
	if err != nil {
		log.Println("Error checking for duplicates:", err)
		return lib.CustomError(http.StatusInternalServerError, "Error checking for duplicates", "Error checking for duplicates")
	}
	if duplicate {
		log.Println("Duplicate data found.")
		// Cek apakah data yang sudah ada masih aktif atau sudah dihapus
		_, err := models.GetAllocSectorStatusByKey(key)
		if err != nil {
			log.Println("Error getting existing data status:", err)
			return lib.CustomError(http.StatusBadRequest, "Duplicate data. Unable to input data.", "Duplicate data. Unable to input data.")
		}
		// if existingDataStatus != 0 {
		// 	log.Println("Existing DATA")
		// 	return lib.CustomError(http.StatusBadRequest, "Duplicate data. Unable to input data.", "Duplicate data. Unable to input data.")
		// }
		if key != allocSectorKey {
			return lib.CustomError(http.StatusBadRequest, "Duplicate data", "Duplicate data")
		}

	}
	status, err = models.UpdateAllocSector(allocSectorKey, params)
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

func DeleteAllocSectorController(c echo.Context) error {
	params := make(map[string]string)
	params["rec_status"] = "0"
	params["rec_deleted_date"] = time.Now().Format(lib.TIMESTAMPFORMAT)
	params["rec_deleted_by"] = lib.UserIDStr

	allocSecKey := c.FormValue("alloc_sector_key")
	if allocSecKey == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing alloc_sector_key", "Missing alloc_sector_key")
	}
	params["alloc_sector_key"] = allocSecKey

	status, err := models.DeleteAllocSector(allocSecKey, params)
	if err != nil {
		return lib.CustomError(status, err.Error(), err.Error())
	}
	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "Berhasil hapus AllocSector!"
	response.Data = ""
	return c.JSON(http.StatusOK, response)
}
