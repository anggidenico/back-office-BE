package controllers

import (
	"mf-bo-api/lib"
	"mf-bo-api/models"
	"net/http"
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
		return lib.CustomError(status, err.Error(), "Failed get data")
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
	if fundTypeKey == "" {
		return lib.CustomError(http.StatusBadRequest, "fund_type_key can not be blank", "fund_type_key can not be blank")
	}
	benchmarkCode := c.FormValue("benchmark_code")
	if benchmarkCode == "" {
		return lib.CustomError(http.StatusBadRequest, "benchmark_code can not be blank", "benchmark_code can not be blank")
	}
	benchmarkName := c.FormValue("benchmark_name")
	if benchmarkName == "" {
		return lib.CustomError(http.StatusBadRequest, "benchmark_name can not be blank", "benchmark_name can not be blank")
	}
	benchmarkShortName := c.FormValue("benchmark_short_name")
	if benchmarkShortName == "" {
		return lib.CustomError(http.StatusBadRequest, "benchmark_short_name can not be blank", "benchmark_short_name can not be blank")
	}
	recAttributeID1 := c.FormValue("rec_attribute_id1")
	if recAttributeID1 == "" {
		return lib.CustomError(http.StatusBadRequest, "rec_attribute_id1 can not be blank", "rec_attribute_id1 can not be blank")
	}

	params["fund_type_key"] = fundTypeKey
	params["benchmark_code"] = benchmarkCode
	params["benchmark_name"] = benchmarkName
	params["benchmark_short_name"] = benchmarkShortName
	params["rec_attribute_id1"] = recAttributeID1
	params["rec_status"] = "1"

	status, err = models.CreateFfsBenchmark(params)
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
	if fundTypeKey == "" {
		return lib.CustomError(http.StatusBadRequest, "fund_type_key can not be blank", "fund_type_key can not be blank")
	}
	benchmarkCode := c.FormValue("benchmark_code")
	if benchmarkCode == "" {
		return lib.CustomError(http.StatusBadRequest, "benchmark_code can not be blank", "benchmark_code can not be blank")
	}
	benchmarkName := c.FormValue("benchmark_name")
	if benchmarkName == "" {
		return lib.CustomError(http.StatusBadRequest, "benchmark_name can not be blank", "benchmark_name can not be blank")
	}
	benchmarkShortName := c.FormValue("benchmark_short_name")
	if benchmarkShortName == "" {
		return lib.CustomError(http.StatusBadRequest, "benchmark_short_name can not be blank", "benchmark_short_name can not be blank")
	}
	// recAttributeID1 := c.FormValue("rec_attribute_id1")
	// if recAttributeID1 == "" {
	// 	return lib.CustomError(http.StatusBadRequest, "rec_attribute_id1 can not be blank", "rec_attribute_id1 can not be blank")
	// }
	// params["benchmark_key"] = benchmarkKey
	params["fund_type_key"] = fundTypeKey
	params["benchmark_code"] = benchmarkCode
	params["benchmark_name"] = benchmarkName
	params["benchmark_short_name"] = benchmarkShortName
	// params["rec_attribute_id1"] = recAttributeID1
	params["rec_status"] = "1"

	status, err = models.UpdateFfsBenchmark(benchmarkKey, params)
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
func UpdateBenchmarkProdController(c echo.Context) error {
	var err error
	params := make(map[string]string)
	params["rec_modified_by"] = lib.UserIDStr
	params["rec_modified_date"] = time.Now().Format(lib.TIMESTAMPFORMAT)

	benchProdKey := c.FormValue("bench_prod_key")
	if benchProdKey == "" {
		return lib.CustomError(http.StatusBadRequest, "bench_prod_key can not be blank", "brench_prod_key can not be blank")
	}
	productKey := c.FormValue("product_key")
	if productKey == "" {
		return lib.CustomError(http.StatusBadRequest, "product_key can not be blank", "product_key can not be blank")
	}
	benchmarkRatio := c.FormValue("benchmark_ratio")
	if benchmarkRatio == "" {
		return lib.CustomError(http.StatusBadRequest, "benchmark_ratio can not be blank", "benchmark_ratio can not be blank")
	}
	params["bench_prod_key"] = benchProdKey
	params["product_key"] = productKey
	params["benchmark_ratio"] = benchmarkRatio
	params["rec_status"] = "1"

	status, err = models.UpdateBenchmarkProduct(benchProdKey, params)
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
