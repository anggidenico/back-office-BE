package controllers

import (
	"encoding/json"
	"mf-bo-api/lib"
	"mf-bo-api/models"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo"
)

func ProductFeeCreateRequest(c echo.Context) error {
	var err error
	params := make(map[string]string)

	productkey := c.FormValue("product_key")
	if productkey == "" {
		// log.Error("Missing required parameter: product_key cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: product_key cann't be blank", "Missing required parameter: product_key cann't be blank")
	}
	strproductkey, err := strconv.ParseUint(productkey, 10, 64)
	if err == nil && strproductkey > 0 {
		params["product_key"] = productkey
	} else {
		// log.Error("Wrong input for parameter: product_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: product_key", "Missing required parameter: product_key")
	}

	//fee_type
	feetype := c.FormValue("fee_type")
	if feetype != "" {
		strfeetype, err := strconv.ParseUint(feetype, 10, 64)
		if err == nil && strfeetype > 0 {
			params["fee_type"] = feetype
		} else {
			// log.Error("Wrong input for parameter: fee_type")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: fee_type", "Missing required parameter: fee_type")
		}
	}

	//fee_code
	feecode := c.FormValue("fee_code")
	if feecode != "" {
		params["fee_code"] = feecode
	}

	//flag_show_ontnc
	flagshowontnc := c.FormValue("flag_show_ontnc")
	var flagshowontncBool bool
	if flagshowontnc != "" {
		flagshowontncBool, err = strconv.ParseBool(flagshowontnc)
		if err != nil {
			// log.Error("flag_show_ontnc parameter should be true/false")
			return lib.CustomError(http.StatusBadRequest, "flag_show_ontnc parameter should be true/false", "flag_show_ontnc parameter should be true/false")
		}
		if flagshowontncBool == true {
			params["flag_show_ontnc"] = "1"
		} else {
			params["flag_show_ontnc"] = "0"
		}
	} else {
		params["flag_show_ontnc"] = "0"
	}

	//fee_annotation
	feeannotation := c.FormValue("fee_annotation")
	if feeannotation != "" {
		params["fee_annotation"] = feeannotation
	}

	//fee_desc
	feedesc := c.FormValue("fee_desc")
	if feedesc != "" {
		params["fee_desc"] = feedesc
	}

	//fee_date_start
	feedatestart := c.FormValue("fee_date_start")
	if feedatestart != "" {
		params["fee_date_start"] = feedatestart + " 00:00:00"
	}

	//fee_date_thru
	feedatethru := c.FormValue("fee_date_thru")
	if feedatethru != "" {
		params["fee_date_thru"] = feedatethru + " 00:00:00"
	}

	//fee_nominal_type
	feenominaltype := c.FormValue("fee_nominal_type")
	if feenominaltype != "" {
		strfeenominaltype, err := strconv.ParseUint(feenominaltype, 10, 64)
		if err == nil && strfeenominaltype > 0 {
			params["fee_nominal_type"] = feenominaltype
		} else {
			// log.Error("Wrong input for parameter: fee_nominal_type")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: fee_nominal_type", "Missing required parameter: fee_nominal_type")
		}
	}

	//enabled_min_amount
	enabledminamount := c.FormValue("enabled_min_amount")
	var enabledminamountBool bool
	if enabledminamount != "" {
		enabledminamountBool, err = strconv.ParseBool(enabledminamount)
		if err != nil {
			// log.Error("enabled_min_amount parameter should be true/false")
			return lib.CustomError(http.StatusBadRequest, "enabled_min_amount parameter should be true/false", "enabled_min_amount parameter should be true/false")
		}
		if enabledminamountBool == true {
			params["enabled_min_amount"] = "1"
		} else {
			params["enabled_min_amount"] = "0"
		}
	} else {
		// log.Error("enabled_min_amount parameter should be true/false")
		return lib.CustomError(http.StatusBadRequest, "enabled_min_amount parameter should be true/false", "enabled_min_amount parameter should be true/false")
	}

	//fee_min_amount
	feeminamount := c.FormValue("fee_min_amount")
	if feeminamount == "" {
		if enabledminamountBool == true {
			// log.Error("Missing required parameter: fee_min_amount cann't be blank")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: fee_min_amount cann't be blank", "Missing required parameter: fee_min_amount cann't be blank")
		}
	} else {
		feeminamountFloat, err := strconv.ParseFloat(feeminamount, 64)
		if err == nil {
			if feeminamountFloat < 0 {
				// log.Error("Wrong input for parameter: fee_min_amount cann't negatif")
				return lib.CustomError(http.StatusBadRequest, "Missing required parameter: fee_min_amount must cann't negatif", "Missing required parameter: fee_min_amount cann't negatif")
			}
			params["fee_min_amount"] = feeminamount
		} else {
			// log.Error("Wrong input for parameter: fee_min_amount number")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: fee_min_amount must number", "Missing required parameter: fee_min_amount number")
		}
	}

	//enabled_max_amount
	enabledmaxamount := c.FormValue("enabled_max_amount")
	var enabledmaxamountBool bool
	if enabledmaxamount != "" {
		enabledmaxamountBool, err = strconv.ParseBool(enabledmaxamount)
		if err != nil {
			// log.Error("enabled_max_amount parameter should be true/false")
			return lib.CustomError(http.StatusBadRequest, "enabled_max_amount parameter should be true/false", "enabled_max_amount parameter should be true/false")
		}
		if enabledmaxamountBool == true {
			params["enabled_max_amount"] = "1"
		} else {
			params["enabled_max_amount"] = "0"
		}
	} else {
		// log.Error("enabled_max_amount parameter should be true/false")
		return lib.CustomError(http.StatusBadRequest, "enabled_max_amount parameter should be true/false", "enabled_max_amount parameter should be true/false")
	}

	//fee_max_amount
	feemaxamount := c.FormValue("fee_max_amount")
	if feemaxamount == "" {
		if enabledmaxamountBool == true {
			// log.Error("Missing required parameter: fee_max_amount cann't be blank")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: fee_max_amount cann't be blank", "Missing required parameter: fee_max_amount cann't be blank")
		}
	} else {
		feemaxamountFloat, err := strconv.ParseFloat(feemaxamount, 64)
		if err == nil {
			if feemaxamountFloat < 0 {
				// log.Error("Wrong input for parameter: fee_max_amount cann't negatif")
				return lib.CustomError(http.StatusBadRequest, "Missing required parameter: fee_max_amount must cann't negatif", "Missing required parameter: fee_max_amount cann't negatif")
			}
			params["fee_max_amount"] = feemaxamount
		} else {
			// log.Error("Wrong input for parameter: fee_max_amount number")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: fee_max_amount must number", "Missing required parameter: fee_max_amount number")
		}
	}

	//fee_calc_method
	feecalcmethod := c.FormValue("fee_calc_method")
	if feecalcmethod != "" {
		strfeecalcmethod, err := strconv.ParseUint(feecalcmethod, 10, 64)
		if err == nil && strfeecalcmethod > 0 {
			params["fee_calc_method"] = feecalcmethod
		} else {
			// log.Error("Wrong input for parameter: fee_calc_method")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: fee_calc_method", "Missing required parameter: fee_calc_method")
		}
	}

	//calculation_baseon
	calculationbaseon := c.FormValue("calculation_baseon")
	if calculationbaseon != "" {
		strcalculationbaseon, err := strconv.ParseUint(calculationbaseon, 10, 64)
		if err == nil && strcalculationbaseon > 0 {
			params["calculation_baseon"] = calculationbaseon
		} else {
			// log.Error("Wrong input for parameter: calculation_baseon")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: calculation_baseon", "Missing required parameter: calculation_baseon")
		}
	}

	//period_hold
	periodhold := c.FormValue("period_hold")
	if periodhold == "" {
		// log.Error("Missing required parameter: period_hold cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: period_hold cann't be blank", "Missing required parameter: period_hold cann't be blank")
	}
	strperiodhold, err := strconv.ParseUint(periodhold, 10, 64)
	if err == nil && strperiodhold > 0 {
		params["period_hold"] = periodhold
	} else {
		// log.Error("Wrong input for parameter: period_hold")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: period_hold", "Missing required parameter: period_hold")
	}

	//days_inyear
	daysinyear := c.FormValue("days_inyear")
	if daysinyear != "" {
		strdaysinyear, err := strconv.ParseUint(daysinyear, 10, 64)
		if err == nil && strdaysinyear > 0 {
			params["days_inyear"] = daysinyear
		} else {
			// log.Error("Wrong input for parameter: days_inyear")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: days_inyear", "Missing required parameter: days_inyear")
		}
	}

	productFeeItems := c.FormValue("product_fee_items")
	if productFeeItems == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: product_fee_items", "Missing required parameter: product_fee_items")
	}

	params["rec_status"] = "1"
	params["rec_created_date"] = time.Now().Format(lib.TIMESTAMPFORMAT)
	params["rec_created_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
	params["rec_action"] = "CREATE"

	var feeItems []models.FeeItemData
	err = json.Unmarshal([]byte(productFeeItems), &feeItems)
	if err != nil {
		return lib.CustomError(http.StatusInternalServerError, err.Error(), err.Error())
	}

	status, err = models.ProductFeeCreateRequest(params, feeItems)
	if err != nil {
		return lib.CustomError(http.StatusInternalServerError, err.Error(), err.Error())
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = nil
	return c.JSON(http.StatusOK, response)
}
