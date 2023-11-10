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

func CreateFfsBenchmarkProductController(c echo.Context) error {
	var err error
	params := make(map[string]string)
	params["rec_created_by"] = lib.UserIDStr
	params["rec_created_date"] = time.Now().Format(lib.TIMESTAMPFORMAT)

	productKey := c.FormValue("product_key")
	if productKey == "" {
		return lib.CustomError(http.StatusBadRequest, "product_key can not be blank", "product_key can not be blank")
	}
	benchmarkRatio := c.FormValue("benchmark_ratio")
	if benchmarkRatio == "" {
		return lib.CustomError(http.StatusBadRequest, "benchmark_ratio can not be blank", "benchmark_ratio can not be blank")
	}
	params["product_key"] = productKey
	params["benchmark_ratio"] = benchmarkRatio
	params["rec_status"] = "1"

	status, err = models.CreateFfsProductBenchmark(params)
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
func GetFfsBenchmarkProductController(c echo.Context) error {
	var benchmarkprod []models.BenchmarkProduct
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
	} else {
		_, err := strconv.ParseUint(benchProdKey, 10, 64)
		if err != sql.ErrNoRows {
			// log.Error("Wrong input for parameter: country_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: bench_prod_key", "Wrong input for parameter: bench_prod_key")
		}
	}
	var detailbenchmarkprod models.BenchmarkProdDetail
	status, err := models.GetBenchmarkProductDetailModels(&detailbenchmarkprod, benchProdKey)
	if err != nil {
		return lib.CustomError(status, err.Error(), "Failed input data")
	}
	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = detailbenchmarkprod
	return c.JSON(http.StatusOK, response)
}
