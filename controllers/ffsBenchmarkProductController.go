package controllers

import (
	"mf-bo-api/lib"
	"mf-bo-api/models"
	"net/http"
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

	result := models.GetBenchmarkProductModels()

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = result
	return c.JSON(http.StatusOK, response)
}
