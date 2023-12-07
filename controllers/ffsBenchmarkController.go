package controllers

import (
	"database/sql"
	"errors"
	"log"
	"mf-bo-api/lib"
	"mf-bo-api/models"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo"
)

func GetFfsBenchmarkController(c echo.Context) error {
	var benchmark []models.Benchmark
	status, err := models.GetBenchmarkModels(&benchmark)
	if err != nil {
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = benchmark
	return c.JSON(http.StatusOK, response)
}

func GetBenchmarkDetailController(c echo.Context) error {
	benchmarkKey := c.Param("benchmark_key")
	if benchmarkKey == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing benchmark key", "Missing benchmark key")
	}
	var detailbenchmark models.BenchmarkDetail
	status, err := models.GetBenchmarkDetailModels(&detailbenchmark, benchmarkKey)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return lib.CustomError(http.StatusNotFound, "benchmark_key not found", "benchmark_key not found")
		}
		return lib.CustomError(status, err.Error(), err.Error())
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = detailbenchmark
	return c.JSON(http.StatusOK, response)
}

func DeleteBenchmarkController(c echo.Context) error {
	params := make(map[string]string)
	params["rec_status"] = "0"
	params["rec_deleted_date"] = time.Now().Format(lib.TIMESTAMPFORMAT)
	params["rec_deleted_by"] = lib.UserIDStr

	benchmarkKey := c.FormValue("benchmark_key")
	if benchmarkKey == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing benchmarkKey", "Missing benchmarkKey")
	}

	status, err := models.DeleteBenchmark(benchmarkKey, params)
	if err != nil {
		return lib.CustomError(status, err.Error(), err.Error())
	}
	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "Berhasil Menghapus Benchmark!"
	response.Data = ""
	return c.JSON(http.StatusOK, response)
}

func CreateFfsBenchmarkController(c echo.Context) error {
	var err error
	params := make(map[string]string)
	params["rec_created_by"] = lib.UserIDStr
	params["rec_created_date"] = time.Now().Format(lib.TIMESTAMPFORMAT)

	fundTypeKey := c.FormValue("fund_type_key")
	if fundTypeKey != "" {
		_, err := strconv.Atoi(fundTypeKey)
		if err != nil {
			return lib.CustomError(http.StatusBadRequest, "fund_type_key should be a number", "fund_type_key should be a number")
		}
		if len(fundTypeKey) > 11 {
			return lib.CustomError(http.StatusBadRequest, "fund_type_key <= 11 digits", "fund_type_key <= 11 digits")
		}
	} else {
		fundTypeKey = "0"
	}
	benchmarkCode := c.FormValue("benchmark_code")
	if benchmarkCode != "" {
		if len(benchmarkCode) > 50 {
			return lib.CustomError(http.StatusBadRequest, "benchmark_code must be <= 50 characters", "benchmark_code must be <= 50 characters")
		}
		benchmarkCode = strings.ToUpper(benchmarkCode)
	}

	benchmarkName := c.FormValue("benchmark_name")
	if benchmarkName != "" {
		if len(benchmarkName) > 150 {
			return lib.CustomError(http.StatusBadRequest, "benchmark_name must be <= 150 characters", "benchmark_name must be <= 150 characters")
		}
	}
	benchmarkShortName := c.FormValue("benchmark_short_name")
	if benchmarkShortName != "" {
		if len(benchmarkShortName) > 70 {
			return lib.CustomError(http.StatusBadRequest, "benchmark_short_name must be <= 70 characters", "benchmark_short_name must be <= 70 characters")
		}
	}
	recAttributeID1 := c.FormValue("rec_attribute_id1")
	if recAttributeID1 != "" {
		if len(recAttributeID1) > 50 {
			return lib.CustomError(http.StatusBadRequest, "rec_attribute_id1 should be exactly 50 characters", "rec_attribute_id1 should be exactly 50 characters")
		}
		params["rec_attribute_id1"] = recAttributeID1
	} else {
		params["rec_attribute_id1"] = "NULL"
	}
	params["fund_type_key"] = fundTypeKey
	params["benchmark_code"] = benchmarkCode
	params["benchmark_name"] = benchmarkName
	params["benchmark_short_name"] = benchmarkShortName
	params["rec_status"] = "1"

	// Check for duplicate records
	duplicate, key, err := models.CheckDuplicateBenchmark(benchmarkCode, benchmarkName)
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
		existingDataStatus, err := models.GetBenchmarkStatusByKey(key)
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
			status, err := models.UpdateBenchmark(key, params)
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
		status, err := models.CreateBenchmark(params)
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
func UpdateFfsBenchmarkController(c echo.Context) error {
	var err error
	params := make(map[string]string)
	params["rec_modified_by"] = lib.UserIDStr
	params["rec_modified_date"] = time.Now().Format(lib.TIMESTAMPFORMAT)

	benchmarkKey := c.FormValue("benchmark_key")
	if benchmarkKey == "" {
		return lib.CustomError(http.StatusBadRequest, "benchmark_key can not be blank", "benchmark_key can not be blank")
	}
	fundTypeKey := c.FormValue("fund_type_key")
	if fundTypeKey != "" {
		_, err := strconv.Atoi(fundTypeKey)
		if err != nil {
			return lib.CustomError(http.StatusBadRequest, "fund_type_key should be number", "fund_type_key should be number")
		}
		if len(fundTypeKey) > 11 {
			return lib.CustomError(http.StatusBadRequest, "fund_type_key <= 11 digits", "fund_type_key <= 11 digits")
		}
	} else {
		fundTypeKey = "0"
	}
	benchmarkCode := c.FormValue("benchmark_code")
	if benchmarkCode != "" {
		if len(benchmarkCode) > 50 {
			return lib.CustomError(http.StatusBadRequest, "fund_type_key must be <= 50 characters", "fund_type_key must be <= 50 characters")
		}
		benchmarkCode = strings.ToUpper(benchmarkCode)
	}
	benchmarkName := c.FormValue("benchmark_name")
	if benchmarkName != "" {
		if len(benchmarkName) > 150 {
			return lib.CustomError(http.StatusBadRequest, "benchmark_name  must be <= 150 character", "benchmark_name must be <= 150 characters")
		}
	}
	benchmarkShortName := c.FormValue("benchmark_short_name")
	if benchmarkShortName != "" {
		if len(benchmarkShortName) > 70 {
			return lib.CustomError(http.StatusBadRequest, "benchmark_short_name must be <= 70 characters", "benchmark_short_name must be <= 70 characters")
		}
	}
	recAttributeID1 := c.FormValue("rec_attribute_id1")
	if recAttributeID1 != "" {
		if len(recAttributeID1) > 50 {
			return lib.CustomError(http.StatusBadRequest, "rec_attribute_id1 must be <= 50 characters", "rec_attribute_id1 must be <= 70 characters")
		}
	}
	recAttributeID2 := c.FormValue("rec_attribute_id2")
	if recAttributeID2 != "" {
		if len(recAttributeID2) > 50 {
			return lib.CustomError(http.StatusBadRequest, "rec_attribute_id1 must be <= 50 characters", "rec_attribute_id1 must be <= 70 characters")
		}
	}
	recAttributeID3 := c.FormValue("rec_attribute_id3")
	if recAttributeID3 != "" {
		if len(recAttributeID2) > 50 {
			return lib.CustomError(http.StatusBadRequest, "rec_attribute_id1 must be <= 50 characters", "rec_attribute_id1 must be <= 70 characters")
		}
	}
	params["fund_type_key"] = fundTypeKey
	params["benchmark_code"] = benchmarkCode
	params["benchmark_name"] = benchmarkName
	params["benchmark_short_name"] = benchmarkShortName
	params["rec_attribute_id1"] = recAttributeID1
	params["rec_attribute_id2"] = recAttributeID2
	params["rec_attribute_id3"] = recAttributeID3
	params["rec_status"] = "1"

	duplicate, key, err := models.CheckDuplicateBenchmark(benchmarkCode, benchmarkName)
	if err != nil {
		log.Println("Error checking for duplicates:", err)
		return lib.CustomError(http.StatusInternalServerError, "Error checking for duplicates", "Error checking for duplicates")
	}
	if duplicate {
		log.Println("Duplicate data found.")
		// Cek apakah data yang sudah ada masih aktif atau sudah dihapus
		existingDataStatus, err := models.GetBenchmarkStatusByKey(key)
		if err != nil {
			log.Println("Error getting existing data status:", err)
			return lib.CustomError(http.StatusBadRequest, "Duplicate data. Unable to input data.", "Duplicate data. Unable to input data.")
		}
		if existingDataStatus != 0 {
			log.Println("Existing DATA")
			return lib.CustomError(http.StatusBadRequest, "Duplicate data. Unable to input data.", "Duplicate data. Unable to input data.")
		}
	}
	status, err = models.UpdateBenchmark(benchmarkKey, params)
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
