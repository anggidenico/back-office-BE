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

func GetFfsBenchmarkProductController(c echo.Context) error {
	var benchmarkprod []models.BenchmarkProduct
	log.Println(benchmarkprod)
	status, err := models.GetBenchmarkProductModels(&benchmarkprod)
	if err != nil {
		return lib.CustomError(status, err.Error(), "Failed input data")
	}
	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = benchmarkprod
	return c.JSON(http.StatusOK, response)
}
func GetBenchmarkProdDetailController(c echo.Context) error {
	benchProdKey := c.Param("bench_prod_key")
	if benchProdKey == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing benchmark_product_key", "Missing benchmark_product_key")
	}

	var detailbenchmarkprod models.BenchmarkProdDetail
	status, err := models.GetBenchmarkProductDetailModels(&detailbenchmarkprod, benchProdKey)
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
	response.Data = detailbenchmarkprod
	return c.JSON(http.StatusOK, response)
}
func DeleteBenchmarkProdController(c echo.Context) error {
	params := make(map[string]string)
	params["rec_status"] = "0"
	params["rec_deleted_date"] = time.Now().Format(lib.TIMESTAMPFORMAT)
	params["rec_deleted_by"] = lib.UserIDStr

	benchProdKey := c.FormValue("bench_prod_key")
	if benchProdKey == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing bench_prod_key", "Missing bench_prod_key")
	}

	status, err := models.DeleteBenchmarkProduct(benchProdKey, params)
	if err != nil {
		return lib.CustomError(status, err.Error(), err.Error())
	}
	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "Berhasil Menghapus Benchmark Product!"
	response.Data = ""
	return c.JSON(http.StatusOK, response)
}
func CreateBenchProdController(c echo.Context) error {
	var err error
	params := make(map[string]interface{})
	params["rec_created_by"] = lib.UserIDStr
	params["rec_created_date"] = time.Now().Format(lib.TIMESTAMPFORMAT)

	productKey := c.FormValue("product_key")
	if productKey == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing product_key", "Missing product_key")
	}
	benchmarkKey := c.FormValue("benchmark_key")
	if benchmarkKey == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing benchmark_key", "Missing benchmark_key")
	}
	benchmarkRatio := c.FormValue("benchmark_ratio")
	if benchmarkRatio == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing benchmark_ratio", "Missing benchmark_ratio")
	}
	benchmarkRemarks := c.FormValue("benchmark_remarks")
	if benchmarkRemarks != "" {
		if len(benchmarkRemarks) > 255 {
			return lib.CustomError(http.StatusBadRequest, "benchmark_remarks must be <= 225 characters", "benchmark_remarks must be <= 225 characters")
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
		params["rec_order"] = value
	} else {
		params["rec_order"] = nil
	}

	// params["rec_order"] = recOrder
	// params["bench_prod_key"] = benchProdKey
	params["product_key"] = productKey
	params["benchmark_key"] = benchmarkKey
	params["benchmark_ratio"] = benchmarkRatio
	params["benchmark_remarks"] = benchmarkRemarks
	params["rec_status"] = "1"

	// Check for duplicate records
	duplicate, key, err := models.CheckDuplicateBenchmarkProd(productKey, benchmarkKey)
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
		existingDataStatus, err := models.GetBenchmarkProdStatusByKey(key)
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
			status, err := models.UpdateBenchmarkProd(key, params)
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
		status, err := models.CreateBenchmarkProd(params)
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

func UpdateBenchmarkProdController(c echo.Context) error {
	var err error
	params := make(map[string]interface{})
	params["rec_modified_by"] = lib.UserIDStr
	params["rec_modified_date"] = time.Now().Format(lib.TIMESTAMPFORMAT)

	benchProdKey := c.FormValue("bench_prod_key")
	if benchProdKey == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing bench_prod_key", "Missing bench_prod_key")
	}
	productKey := c.FormValue("product_key")
	if productKey == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing product_key", "Missing product_key")
	}
	benchmarkKey := c.FormValue("benchmark_key")
	if benchmarkKey == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing benchmark_key", "Missing benchmark_key")
	}
	benchmarkRatio, err := strconv.ParseInt(c.FormValue("benchmark_ratio"), 10, 64)
	if err != nil {
		return lib.CustomError(http.StatusBadRequest, "Invalid benchmark_ratio", "Invalid benchmark_ratio")
	}
	benchmarkRemarks := c.FormValue("benchmark_remarks")
	if benchmarkRemarks == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing benchmark_remarks", "Missing benchmark_remarks")
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
	params["bench_prod_key"] = benchProdKey
	params["product_key"] = productKey
	params["benchmark_key"] = benchmarkKey
	params["benchmark_ratio"] = benchmarkRatio
	params["benchmark_remarks"] = benchmarkRemarks
	params["rec_status"] = "1"

	// Check for duplicate records
	duplicate, key, err := models.CheckDuplicateBenchmarkProd(productKey, benchmarkKey)
	if err != nil {
		log.Println("Error checking for duplicates:", err)
		return lib.CustomError(http.StatusInternalServerError, "Error checking for duplicates", "Error checking for duplicates")
	}
	if duplicate {
		log.Println("Duplicate data found.")
		// Cek apakah data yang sudah ada masih aktif atau sudah dihapus
		_, err := models.GetBenchmarkProdStatusByKey(key)
		if err != nil {
			log.Println("Error getting existing data status:", err)
			return lib.CustomError(http.StatusBadRequest, "Duplicate data. Unable to input data.", "Duplicate data. Unable to input data.")
		}
		// if existingDataStatus != 0 {
		// 	log.Println("Existing DATA")
		// 	return lib.CustomError(http.StatusBadRequest, "Duplicate data. Unable to input data.", "Duplicate data. Unable to input data.")
		// }
		if key != benchProdKey {
			return lib.CustomError(http.StatusBadRequest, "Duplicate data", "Duplicate data")
		}

	}
	status, err = models.UpdateBenchmarkProd(benchProdKey, params)
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
