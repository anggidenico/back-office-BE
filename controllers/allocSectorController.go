package controllers

import (
	"database/sql"
	"errors"
	"log"
	"mf-bo-api/lib"
	"mf-bo-api/models"
	"net/http"
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
	params := make(map[string]string)
	params["rec_created_by"] = lib.UserIDStr
	params["rec_created_date"] = time.Now().Format(lib.TIMESTAMPFORMAT)

	productKey := c.FormValue("product_key")
	if productKey == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing product_key", "Missing product_key")
	}
	params["product_key"] = productKey

	params["rec_status"] = "1"

	// Check for duplicate records
	duplicate, key, err := models.CheckDuplicateFfsBenchmark(params["benchmark_code"], params["benchmark_name"])
	if err != nil {
		log.Println("Error checking for duplicates:", err)
		return lib.CustomError(http.StatusInternalServerError, "Error checking for duplicates", "Error checking for duplicates")
	}
	log.Println("Duplicate:", duplicate)
	log.Println("Key:", key)
	// Jika duplikasi ditemukan, perbarui data yang sudah ada
	if duplicate {
		status, err := models.UpdateBenchmark(key, params)
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
	status, err = models.CreateFfsBenchmark(params)
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
