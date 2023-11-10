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

func CreateMsSecuritiesController(c echo.Context) error {
	var err error
	params := make(map[string]string)
	params["rec_created_by"] = lib.UserIDStr
	params["rec_created_date"] = time.Now().Format(lib.TIMESTAMPFORMAT)

	secCode := c.FormValue("sec_code")
	if secCode == "" {
		return lib.CustomError(http.StatusBadRequest, "sec_code can not be blank", "sec_code can not be blank")
	}
	secName := c.FormValue("sec_name")
	if secName == "" {
		return lib.CustomError(http.StatusBadRequest, "sec_name can not be blank", "sec_name can not be blank")
	}
	// Memeriksa apakah sec_name sudah ada dalam peta (map) validasi
	// if _, exists := validatedSecNames[secName]; exists {
	// 	return lib.CustomError(http.StatusBadRequest, "sec_name already exists", "sec_name already exists")
	// }

	// // ... (lanjutkan dengan kode validasi lainnya)

	// // Setelah berhasil memvalidasi, tambahkan sec_name ke dalam peta validasi
	// validatedSecNames[secName] = true
	secCategory := c.FormValue("securities_category")
	if secCategory == "" {
		return lib.CustomError(http.StatusBadRequest, "securities_category can not be blank", "securities_category can not be blank")
	}
	secType := c.FormValue("security_type")
	if secType == "" {
		return lib.CustomError(http.StatusBadRequest, "security_type can not be blank", "security_type can not be blank")
	}
	currencyKey := c.FormValue("currency_key")
	if currencyKey == "" {
		return lib.CustomError(http.StatusBadRequest, "currency_key can not be blank", "currency_key can not be blank")
	}
	secStatus := c.FormValue("security_status")
	if secStatus == "" {
		return lib.CustomError(http.StatusBadRequest, "security_status can not be blank", "security_status can not be blank")
	}
	isinCode := c.FormValue("isin_code")
	if isinCode == "" {
		return lib.CustomError(http.StatusBadRequest, "isin_code can not be blank", "isin_code can not be blank")
	}
	secClassification := c.FormValue("sec_classification")
	if secClassification == "" {
		return lib.CustomError(http.StatusBadRequest, "sec_classification can not be blank", "sec_classification can not be blank")
	}
	today := time.Now()
	pastDue := today.AddDate(1, 0, 0)
	pastDueDate := pastDue.Format(lib.TIMESTAMPFORMAT)

	dateIs := today.AddDate(1, -1, -2)
	dateIssued := dateIs.Format(lib.TIMESTAMPFORMAT)

	params["sec_code"] = secCode
	params["sec_name"] = secName
	params["securities_category"] = secCategory
	params["security_type"] = secType
	params["currency_key"] = currencyKey
	params["security_status"] = secStatus
	params["isin_code"] = isinCode
	params["sec_classification"] = secClassification
	params["date_issued"] = dateIssued
	params["date_matured"] = pastDueDate
	params["rec_status"] = "1"

	status, err = models.CreateMsSecurities(params)
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
func GetMsSecuritiesController(c echo.Context) error {
	var sec []models.Securities
	status, err := models.GetSecuritiesModels(&sec)
	if err != nil {
		return lib.CustomError(status, err.Error(), err.Error())
	}
	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = sec
	return c.JSON(http.StatusOK, response)
}
func DeleteMsSecuritiesController(c echo.Context) error {
	params := make(map[string]string)
	params["rec_status"] = "0"
	params["rec_deleted_date"] = time.Now().Format(lib.TIMESTAMPFORMAT)
	params["rec_deleted_by"] = lib.UserIDStr

	secKey := c.FormValue("sec_key")
	if secKey == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing sec_key", "Missing sec_key")
	}

	status, err := models.DeleteMsSecurities(secKey, params)
	if err != nil {
		return lib.CustomError(status, err.Error(), err.Error())
	}
	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "Berhasil hapus Master Securities!"
	response.Data = ""
	return c.JSON(http.StatusOK, response)
}
func GetMsSecuritiesDetailController(c echo.Context) error {
	secKey := c.Param("sec_key")
	if secKey == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing sec_key", "Missing sec_key")
	} else {
		_, err := strconv.ParseUint(secKey, 10, 64)
		if err != sql.ErrNoRows {
			// log.Error("Wrong input for parameter: country_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: sec_key", "Wrong input for parameter: sec_key")
		}
	}
	var detailmssec models.SecuritiesDetail
	status, err := models.GetMsSecuritiesDetailModels(&detailmssec, secKey)
	if err != nil {
		return lib.CustomError(status, err.Error(), "Failed input data")
	}
	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = detailmssec
	return c.JSON(http.StatusOK, response)
}
func UpdateMsSecuritiesController(c echo.Context) error {
	var err error
	params := make(map[string]string)
	params["rec_modified_by"] = lib.UserIDStr
	params["rec_modified_date"] = time.Now().Format(lib.TIMESTAMPFORMAT)

	secKey := c.FormValue("sec_key")
	if secKey == "" {
		return lib.CustomError(http.StatusBadRequest, "sec_key can not be blank", "sec_key can not be blank")
	}
	secCode := c.FormValue("sec_code")
	if secCode == "" {
		return lib.CustomError(http.StatusBadRequest, "sec_code can not be blank", "sec_code can not be blank")
	}
	secName := c.FormValue("sec_name")
	if secName == "" {
		return lib.CustomError(http.StatusBadRequest, "sec_name can not be blank", "sec_name can not be blank")
	}
	secCategory := c.FormValue("securities_category")
	if secCategory == "" {
		return lib.CustomError(http.StatusBadRequest, "securities_category can not be blank", "securities_category can not be blank")
	}
	secType := c.FormValue("security_type")
	if secType == "" {
		return lib.CustomError(http.StatusBadRequest, "security_type can not be blank", "security_type can not be blank")
	}
	currencyKey := c.FormValue("currency_key")
	if currencyKey == "" {
		return lib.CustomError(http.StatusBadRequest, "currency_key can not be blank", "currency_key can not be blank")
	}
	secStatus := c.FormValue("security_status")
	if secStatus == "" {
		return lib.CustomError(http.StatusBadRequest, "security_status can not be blank", "security_status can not be blank")
	}
	isinCode := c.FormValue("isin_code")
	if isinCode == "" {
		return lib.CustomError(http.StatusBadRequest, "isin_code can not be blank", "isin_code can not be blank")
	}
	secClassification := c.FormValue("sec_classification")
	if secClassification == "" {
		return lib.CustomError(http.StatusBadRequest, "sec_classification can not be blank", "sec_classification can not be blank")
	}
	today := time.Now()
	pastDue := today.AddDate(1, 0, 0)
	pastDueDate := pastDue.Format(lib.TIMESTAMPFORMAT)

	dateIs := today.AddDate(1, -1, -2)
	dateIssued := dateIs.Format(lib.TIMESTAMPFORMAT)
	params["sec_code"] = secCode
	params["sec_name"] = secName
	params["securities_category"] = secCategory
	params["security_type"] = secType
	params["currency_key"] = currencyKey
	params["security_status"] = secStatus
	params["isin_code"] = isinCode
	params["sec_classification"] = secClassification
	params["date_issued"] = dateIssued
	params["date_matured"] = pastDueDate
	params["rec_status"] = "1"

	status, err = models.UpdateMsSecurities(secKey, params)
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
