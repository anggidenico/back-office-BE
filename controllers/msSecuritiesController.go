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

func CreateMsSecuritiesController(c echo.Context) error {
	var err error
	// params := make(map[string]string)
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

	secParentKey := c.FormValue("sec_parent_key")
	if secParentKey != "" {
		value, err := strconv.Atoi(secParentKey)
		if err != nil {
			return lib.CustomError(value, "sec_parent_key must be number", "sec_parent_key must be number")
		}
	}
	params["sec_parent_key"] = secParentKey

	secCategory := c.FormValue("securities_category")
	if secCategory != "" {
		value, err := strconv.Atoi(secCategory)
		if err != nil {
			return lib.CustomError(value, "securities_category must be number", "securities_category must be number")
		}
	}
	params["securities_category"] = secCategory

	secType := c.FormValue("security_type")
	if secType != "" {
		value, err := strconv.Atoi(secType)
		if err != nil {
			return lib.CustomError(value, "security_type must be number", "security_type must be number")
		}
	}
	params["security_type"] = secType

	currencyKey := c.FormValue("currency_key")
	if currencyKey != "" {
		value, err := strconv.Atoi(currencyKey)
		if err != nil {
			return lib.CustomError(value, "currency_key must be number", "currency_key must be number")
		}
	}
	params["currency_key"] = currencyKey

	isinCode := c.FormValue("isin_code")
	params["isin_code"] = isinCode

	sectorKey := c.FormValue("sector_key")
	if sectorKey != "" {
		value, err := strconv.Atoi(sectorKey)
		if err != nil {
			return lib.CustomError(value, "sector_key must be number", "sector_key must be number")
		}
	}
	params["sector_key"] = sectorKey

	secClassification := c.FormValue("sec_classification")
	if secClassification != "" {
		value, err := strconv.Atoi(secClassification)
		if err != nil {
			return lib.CustomError(value, "sec_classification must be number", "sec_classification must be number")
		}
	}
	params["sec_classification"] = secClassification

	secTenorMonth := c.FormValue("sec_tenor_month")
	if secTenorMonth != "" {
		value, err := strconv.Atoi(secTenorMonth)
		if err != nil {
			return lib.CustomError(value, "sec_tenor_month must be number", "sec_tenor_month must be number")
		}
	}
	params["sec_tenor_month"] = secTenorMonth

	securityStatus := c.FormValue("security_status")
	if securityStatus != "" {
		value, err := strconv.Atoi(securityStatus)
		if err != nil {
			return lib.CustomError(value, "security_status must be number", "security_status must be number")
		}
	}
	params["security_status"] = securityStatus

	secShares := c.FormValue("sec_shares")
	if secShares != "" {
		value, err := strconv.Atoi(secShares)
		if err != nil {
			return lib.CustomError(value, "sec_shares must be number", "sec_shares must be number")
		}
	}
	params["sec_shares"] = secShares

	flagSyariahStr := c.FormValue("flag_syariah")
	flagSyariah, err := strconv.ParseBool(flagSyariahStr)
	if err != nil {
		return lib.CustomError(http.StatusBadRequest, "Invalid value for flag_syariah", err.Error())
	}
	if flagSyariah == true {
		params["flag_syariah"] = "1"
	} else {
		params["flag_syariah"] = "0"

	}

	flagIsBreakableStr := c.FormValue("flag_is_breakable")
	flagIsBreakable, err := strconv.ParseBool(flagIsBreakableStr)
	if err != nil {
		return lib.CustomError(http.StatusBadRequest, "Invalid value for flag_is_breakable", err.Error())
	}
	if flagIsBreakable == true {
		params["flag_is_breakable"] = "1"
	} else {
		params["flag_is_breakable"] = "0"
	}

	flaghasCouponStr := c.FormValue("flag_has_coupon")
	flagHasCoupon, err := strconv.ParseBool(flaghasCouponStr)
	if err != nil {
		return lib.CustomError(http.StatusBadRequest, "Invalid value for flag_has_coupon", err.Error())
	}
	if flagHasCoupon == true {
		params["flag_has_coupon"] = "1"
	} else {
		params["flag_has_coupon"] = "0"
	}

	stockMarket := c.FormValue("stock_market")
	if stockMarket != "" {
		value, err := strconv.Atoi(stockMarket)
		if err != nil {
			return lib.CustomError(value, "stock_market must be number", "stock_market must be number")
		}
	}
	params["stock_market"] = stockMarket

	secPaRates := c.FormValue("sec_pa_rates")
	if secPaRates != "" {
		value, err := strconv.Atoi(secPaRates)
		if err != nil {
			return lib.CustomError(value, "sec_pa_rates must be number", "sec_pa_rates must be number")
		}
	}
	params["sec_pa_rates"] = secPaRates

	secPrincipleValue := c.FormValue("sec_principle_value")
	if secPrincipleValue != "" {
		value, err := strconv.Atoi(secPrincipleValue)
		if err != nil {
			return lib.CustomError(value, "sec_principle_value must be number", "sec_principle_value must be number")
		}
	}
	params["sec_principle_value"] = secPrincipleValue

	taxRates := c.FormValue("tax_rates")
	if taxRates != "" {
		value, err := strconv.Atoi(taxRates)
		if err != nil {
			return lib.CustomError(value, "tax_rates must be number", "tax_rates must be number")
		}
	}
	params["tax_rates"] = taxRates

	participantKey := c.FormValue("participant_key")
	if participantKey != "" {
		value, err := strconv.Atoi(participantKey)
		if err != nil {
			return lib.CustomError(value, "participant_key must be number", "participant_key must be number")
		}
	}
	params["participant_key"] = participantKey

	couponType := c.FormValue("coupon_type")
	if couponType != "" {
		value, err := strconv.Atoi(couponType)
		if err != nil {
			return lib.CustomError(value, "coupon_type must be number", "coupon_type must be number")
		}
	}
	params["coupon_type"] = couponType

	recOrder := c.FormValue("rec_order")
	if recOrder != "" {
		if len(recOrder) > 11 {
			return lib.CustomError(http.StatusBadRequest, "rec_order should be exactly 11 characters", "rec_order be exactly 11 characters")
		}
		_, err := strconv.Atoi(recOrder)
		if err != nil {
			return lib.CustomError(http.StatusBadRequest, "rec_order should be a number", "rec_order should be a number")
		}
	}
	params["rec_order"] = recOrder

	today := time.Now()
	pastDue := today.AddDate(1, 0, 0)
	pastDueDate := pastDue.Format(lib.TIMESTAMPFORMAT)

	dateIs := today.AddDate(1, -1, -2)
	dateIssued := dateIs.Format(lib.TIMESTAMPFORMAT)
	params["sec_code"] = secCode
	params["sec_name"] = secName
	params["date_issued"] = dateIssued
	params["date_matured"] = pastDueDate
	params["rec_status"] = "1"

	_, err = models.CreateMsSecurities(params)
	if err != nil {

	}

	// duplicate, key, err := models.CheckDuplicateSecurities(secCode, secName, secType)
	// if err != nil {
	// 	log.Println("Error checking for duplicates:", err)
	// 	return lib.CustomError(http.StatusInternalServerError, "Error checking for duplicates", "Error checking for duplicates")
	// }

	// log.Println("Duplicate:", duplicate)
	// log.Println("Key:", key)

	// // Jika duplikasi ditemukan, perbarui data yang sudah ada
	// if duplicate {
	// 	log.Println("Duplicate data found.")
	// 	// Cek apakah data yang sudah ada masih aktif atau sudah dihapus
	// 	existingDataStatus, err := models.GetSecuritiesStatusByKey(key)
	// 	if err != nil {
	// 		log.Println("Error getting existing data status:", err)
	// 		return lib.CustomError(http.StatusInternalServerError, "Error getting existing data status", "Error getting existing data status")
	// 	}

	// 	// Jika data sudah dihapus (rec_status = 0), perbarui statusnya menjadi aktif (rec_status = 1)
	// 	if existingDataStatus == 0 {
	// 		log.Println("Existing data is deleted. Recreating data.")

	// 		// Set status menjadi aktif (rec_status = 1)
	// 		params["rec_status"] = "1"
	// 		// Update data dengan status baru dan nilai-nilai yang baru
	// 		status, err := models.UpdateMsSecurities(key, params)
	// 		if err != nil {
	// 			log.Println("Error updating data:", err)
	// 			return lib.CustomError(status, "Error updating data", "Error updating data")
	// 		}
	// 		return c.JSON(http.StatusOK, lib.Response{
	// 			Status: lib.Status{
	// 				Code:          http.StatusOK,
	// 				MessageServer: "OK",
	// 				MessageClient: "OK",
	// 			},
	// 			Data: "Data updated successfully",
	// 		})
	// 	} else {
	// 		// Jika data masih aktif, kembalikan respons kesalahan duplikasi
	// 		log.Println("Existing data is still active. Duplicate data error.")
	// 		return lib.CustomError(http.StatusBadRequest, "Duplicate data. Unable to input data.", "Duplicate data. Unable to input data.")
	// 	}
	// } else {
	// 	// Jika tidak ada duplikasi, buat data baru
	// 	status, err := models.CreateMsSecurities(params)
	// 	if err != nil {
	// 		log.Println("Error create data:", err)
	// 		return lib.CustomError(status, "Duplicate data. Unable to input data.", "Duplicate data. Unable to input data.")
	// 	}
	// }

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = "Data Successfully Created"
	return c.JSON(http.StatusOK, response)
}

func GetMsSecuritiesController(c echo.Context) error {
	decimal.MarshalJSONWithoutQuotes = true
	var sec []models.Securities
	status, err := models.GetSecuritiesModels(&sec)
	if err != nil {
		return lib.CustomError(status, err.Error(), err.Error())
	}

	var responseData []models.SecuritiesResponse
	if len(sec) > 0 {
		for _, data := range sec {
			var rData models.SecuritiesResponse
			rData.SecKey = data.SecKey
			rData.CouponName = data.CouponName
			rData.SecParentKey = data.SecParentKey
			rData.SecCode = data.SecCode
			rData.SecName = data.SecName
			rData.CouponType = data.CouponType
			rData.CurrencyCode = data.CurrencyCode
			rData.CurrencyName = data.CurrencyName
			rData.CurrencyKey = data.CurrencyKey
			rData.DateIssued = data.DateIssued
			rData.DateMatured = data.DateMatured
			rData.IsinCode = data.IsinCode
			rData.ParticipantKey = data.ParticipantKey
			rData.ParticipantName = data.ParticipantName
			rData.RecOrder = data.RecOrder
			rData.SecClassification = data.SecClassification
			rData.SecClassificationName = data.SecClassificationName
			rData.SecParates = data.SecParates
			rData.SecPrincipleValue = data.SecPrincipleValue
			rData.SecShares = data.SecShares
			rData.SecTenorMonth = data.SecTenorMonth
			rData.SecuritiesCategory = data.SecuritiesCategory
			rData.SecuritiesCategoryName = data.SecuritiesCategoryName
			rData.SecurityStatus = data.SecurityStatus
			rData.SecurityType = data.SecurityType
			rData.SecurityTypeName = data.SecurityTypeName
			rData.StocKMarketName = data.StocKMarketName
			rData.StockMarket = data.StockMarket
			rData.TaxRates = data.TaxRates

			if len(data.FlagSyariah) > 0 {
				flagSyariahValue := data.FlagSyariah[0] == 1
				rData.FlagSyariah = &flagSyariahValue
			} else {
				rData.FlagSyariah = nil
			}
			if len(data.FlagIsBreakable) > 0 {
				flagBreakableValue := data.FlagSyariah[0] == 1
				rData.FlagIsBreakable = &flagBreakableValue
			} else {
				rData.FlagIsBreakable = nil
			}
			if len(data.FlagHasCoupon) > 0 {
				flagCouponValue := data.FlagHasCoupon[0] == 1
				rData.FlagHasCoupon = &flagCouponValue
			} else {
				rData.FlagIsBreakable = nil
			}
			responseData = append(responseData, rData)
		}
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = responseData
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
	decimal.MarshalJSONWithoutQuotes = true
	secKey := c.Param("sec_key")
	if secKey == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing sec_key", "Missing sec_key")
	}

	var sec models.SecuritiesDetail
	status, err := models.GetMsSecuritiesDetailModels(&sec, secKey)
	if err != nil {
		return lib.CustomError(status, err.Error(), err.Error())
	}

	var responseData models.SecuritiesResponse
	var rData models.SecuritiesResponse

	rData.SecKey = sec.SecKey
	rData.CouponName = sec.CouponName
	rData.SecParentKey = sec.SecParentKey
	rData.SecCode = sec.SecCode
	rData.SecName = sec.SecName
	rData.CouponType = sec.CouponType
	rData.CurrencyCode = sec.CurrencyCode
	rData.CurrencyName = sec.CurrencyName
	rData.CurrencyKey = sec.CurrencyKey
	rData.DateIssued = sec.DateIssued
	rData.DateMatured = sec.DateMatured
	rData.IsinCode = sec.IsinCode
	rData.ParticipantKey = sec.ParticipantKey
	rData.ParticipantName = sec.ParticipantName
	rData.RecOrder = sec.RecOrder
	rData.SecClassification = sec.SecClassification
	rData.SecClassificationName = sec.SecClassificationName
	rData.SecParates = sec.SecParates
	rData.SecPrincipleValue = sec.SecPrincipleValue
	rData.SecShares = sec.SecShares
	rData.SecTenorMonth = sec.SecTenorMonth
	rData.SecuritiesCategory = sec.SecuritiesCategory
	rData.SecuritiesCategoryName = sec.SecuritiesCategoryName
	rData.SecurityStatus = sec.SecurityStatus
	rData.SecurityType = sec.SecurityType
	rData.SecurityTypeName = sec.SecurityTypeName
	rData.SectorKey = sec.SectorKey
	rData.SectorName = sec.SectorName
	rData.SecTenorMonth = sec.SecTenorMonth
	rData.SecurityStatusName = sec.SecurityStatusName
	rData.StocKMarketName = sec.StocKMarketName
	rData.StockMarket = sec.StockMarket
	rData.TaxRates = sec.TaxRates

	if len(sec.FlagSyariah) > 0 {
		flagSyariahValue := sec.FlagSyariah[0] == 1
		rData.FlagSyariah = &flagSyariahValue
	} else {
		rData.FlagSyariah = nil
	}
	if len(sec.FlagIsBreakable) > 0 {
		flagBreakableValue := sec.FlagIsBreakable[0] == 1
		rData.FlagIsBreakable = &flagBreakableValue
	} else {
		rData.FlagIsBreakable = nil
	}
	if len(sec.FlagHasCoupon) > 0 {
		flagCouponValue := sec.FlagHasCoupon[0] == 1
		rData.FlagHasCoupon = &flagCouponValue
	} else {
		rData.FlagHasCoupon = nil
	}

	responseData = rData

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = responseData

	return c.JSON(http.StatusOK, response)
}

func UpdateMsSecuritiesController(c echo.Context) error {
	var err error
	params := make(map[string]interface{})
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

	secParentKey := c.FormValue("sec_parent_key")
	if secParentKey != "" {
		value, err := strconv.Atoi(secParentKey)
		if err != nil {
			return lib.CustomError(value, "sec_parent_key must be number", "sec_parent_key must be number")
		}
		params["sec_parent_key"] = secParentKey
	} else {
		params["sec_parent_key"] = nil
	}

	secCategory := c.FormValue("securities_category")
	if secCategory != "" {
		value, err := strconv.Atoi(secCategory)
		if err != nil {
			return lib.CustomError(value, "securities_category must be number", "securities_category must be number")
		}
		params["securities_category"] = secCategory
	} else {
		params["securities_category"] = nil
	}
	secType := c.FormValue("security_type")
	if secType != "" {
		value, err := strconv.Atoi(secType)
		if err != nil {
			return lib.CustomError(value, "security_type must be number", "security_type must be number")
		}
		params["security_type"] = secType
	} else {
		params["security_type"] = nil
	}
	currencyKey := c.FormValue("currency_key")
	if currencyKey != "" {
		value, err := strconv.Atoi(currencyKey)
		if err != nil {
			return lib.CustomError(value, "currency_key must be number", "currency_key must be number")
		}
		params["currency_key"] = currencyKey
	} else {
		params["currency_key"] = nil
	}

	isinCode := c.FormValue("isin_code")
	params["isin_code"] = isinCode

	sectorKey := c.FormValue("sector_key")
	if sectorKey != "" {
		value, err := strconv.Atoi(sectorKey)
		if err != nil {
			return lib.CustomError(value, "sector_key must be number", "sector_key must be number")
		}
		params["sector_key"] = sectorKey
	} else {
		params["sector_key"] = nil
	}
	secClassification := c.FormValue("sec_classification")
	if secClassification != "" {
		value, err := strconv.Atoi(secClassification)
		if err != nil {
			return lib.CustomError(value, "sec_classification must be number", "sec_classification must be number")
		}
		params["sec_classification"] = secClassification
	} else {
		params["sec_classification"] = nil
	}
	secTenorMonth := c.FormValue("sec_tenor_month")
	if secTenorMonth != "" {
		value, err := strconv.Atoi(secTenorMonth)
		if err != nil {
			return lib.CustomError(value, "sec_tenor_month must be number", "sec_tenor_month must be number")
		}
		params["sec_tenor_month"] = secTenorMonth
	} else {
		params["sec_tenor_month"] = nil
	}

	securityStatus := c.FormValue("security_status")
	if securityStatus != "" {
		value, err := strconv.Atoi(securityStatus)
		if err != nil {
			return lib.CustomError(value, "security_status must be number", "security_status must be number")
		}
		params["security_status"] = securityStatus
	} else {
		params["security_status"] = nil
	}

	secShares := c.FormValue("sec_shares")
	if secShares != "" {
		value, err := strconv.Atoi(secShares)
		if err != nil {
			return lib.CustomError(value, "sec_shares must be number", "sec_shares must be number")
		}
		params["sec_shares"] = secShares
	} else {
		params["sec_shares"] = nil
	}

	flagSyariahStr := c.FormValue("flag_syariah")
	flagSyariah, err := strconv.ParseBool(flagSyariahStr)
	if err != nil {
		return lib.CustomError(http.StatusBadRequest, "Invalid value for flag_syariah", err.Error())
	}
	params["flag_syariah"] = flagSyariah

	flagIsBreakableStr := c.FormValue("flag_is_breakable")
	flagIsBreakable, err := strconv.ParseBool(flagIsBreakableStr)
	if err != nil {
		return lib.CustomError(http.StatusBadRequest, "Invalid value for flag_is_breakable", err.Error())
	}
	params["flag_is_breakable"] = flagIsBreakable

	flaghasCouponStr := c.FormValue("flag_has_coupon")
	flagHasCoupon, err := strconv.ParseBool(flaghasCouponStr)
	if err != nil {
		return lib.CustomError(http.StatusBadRequest, "Invalid value for flag_has_coupon", err.Error())
	}
	params["flag_has_coupon"] = flagHasCoupon

	stockMarket := c.FormValue("stock_market")
	if stockMarket != "" {
		value, err := strconv.Atoi(stockMarket)
		if err != nil {
			return lib.CustomError(value, "stock_market must be number", "stock_market must be number")
		}
		params["stock_market"] = stockMarket
	} else {
		params["stock_market"] = nil
	}

	secPaRates := c.FormValue("sec_pa_rates")
	if secPaRates != "" {
		value, err := strconv.Atoi(secPaRates)
		if err != nil {
			return lib.CustomError(value, "sec_pa_rates must be number", "sec_pa_rates must be number")
		}
		params["sec_pa_rates"] = secPaRates
	} else {
		params["sec_pa_rates"] = nil
	}
	secPrincipleValue := c.FormValue("sec_principle_value")
	if secPrincipleValue != "" {
		value, err := strconv.Atoi(secPrincipleValue)
		if err != nil {
			return lib.CustomError(value, "sec_principle_value must be number", "sec_principle_value must be number")
		}
		params["sec_principle_value"] = secPrincipleValue
	} else {
		params["sec_principle_value"] = nil
	}

	taxRates := c.FormValue("tax_rates")
	if taxRates != "" {
		value, err := strconv.Atoi(taxRates)
		if err != nil {
			return lib.CustomError(value, "tax_rates must be number", "tax_rates must be number")
		}
		params["tax_rates"] = taxRates
	} else {
		params["tax_rates"] = nil
	}

	participantKey := c.FormValue("participant_key")
	if participantKey != "" {
		value, err := strconv.Atoi(participantKey)
		if err != nil {
			return lib.CustomError(value, "participant_key must be number", "participant_key must be number")
		}
		params["participant_key"] = participantKey
	} else {
		params["participant_key"] = nil
	}

	couponType := c.FormValue("coupon_type")
	if couponType != "" {
		value, err := strconv.Atoi(couponType)
		if err != nil {
			return lib.CustomError(value, "coupon_type must be number", "coupon_type must be number")
		}
		params["coupon_type"] = couponType
	} else {
		params["coupon_type"] = nil
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

	today := time.Now()
	pastDue := today.AddDate(1, 0, 0)
	pastDueDate := pastDue.Format(lib.TIMESTAMPFORMAT)

	dateIs := today.AddDate(1, -1, -2)
	dateIssued := dateIs.Format(lib.TIMESTAMPFORMAT)
	params["sec_code"] = secCode
	params["sec_name"] = secName
	params["date_issued"] = dateIssued
	params["date_matured"] = pastDueDate
	params["rec_status"] = "1"

	duplicate, key, err := models.CheckDuplicateSecurities(secCode, secName, secType)
	if err != nil {
		log.Println("Error checking for duplicates:", err)
		return lib.CustomError(http.StatusInternalServerError, "Error checking for duplicates", "Error checking for duplicates")
	}
	if duplicate {
		log.Println("Duplicate data found.")
		// Cek apakah data yang sudah ada masih aktif atau sudah dihapus
		_, err := models.GetSecuritiesStatusByKey(key)
		if err != nil {
			log.Println("Error getting existing data status:", err)
			return lib.CustomError(http.StatusBadRequest, "Duplicate data. Unable to input data.", "Duplicate data. Unable to input data.")
		}

		if key != secKey {
			return lib.CustomError(http.StatusBadRequest, "Duplicate data", "Duplicate data")
		}

	}
	status, err = models.UpdateMsSecurities(secKey, params)
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

func GetParticipantKeyController(c echo.Context) error {
	var value []models.ParticipantList
	status, err := models.GetParticipantListModels(&value)
	if err != nil {
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = value
	// log.Printf("Response Data: %+v\n", response.Data)
	return c.JSON(http.StatusOK, response)
}
