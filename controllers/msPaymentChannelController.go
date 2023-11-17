package controllers

import (
	"math/big"
	"mf-bo-api/lib"
	"mf-bo-api/models"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo"
)

func GetmsPaymentChannelController(c echo.Context) error {
	var mspayment []models.PaymentChannel
	status, err := models.GetPaymentChannelModels(&mspayment)
	if err != nil {
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = mspayment
	return c.JSON(http.StatusOK, response)
}

func GetMsPaymentDetailController(c echo.Context) error {
	pChannelKey := c.Param("pchannel_key")
	if pChannelKey == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing payment channel key", "Missing payment channel key")
	}
	var detailpayment models.PaymentChannelDetail
	status, err := models.GetDetailPaymentChannelModels(&detailpayment, pChannelKey)
	if err != nil {
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	// log.Println("Not Found")

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = detailpayment
	return c.JSON(http.StatusOK, response)
}

func DeleteMsPaymentChannelController(c echo.Context) error {
	params := make(map[string]string)
	params["rec_status"] = "0"
	params["rec_deleted_date"] = time.Now().Format(lib.TIMESTAMPFORMAT)
	params["rec_deleted_by"] = lib.UserIDStr

	pChannelKey := c.FormValue("pchannel_key")
	if pChannelKey == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing pchannelKey", "Missing pchannelKey")
	}

	status, err := models.DeleteMsPaymentChannel(pChannelKey, params)
	if err != nil {
		return lib.CustomError(status, err.Error(), err.Error())
	}
	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "Berhasil hapus Master Payment Channel!"
	response.Data = ""
	return c.JSON(http.StatusOK, response)
}

func CreateMsPaymentChannelController(c echo.Context) error {
	var err error
	params := make(map[string]string)
	params["rec_created_by"] = lib.UserIDStr
	params["rec_created_date"] = time.Now().Format(lib.TIMESTAMPFORMAT)

	pChannelCode := c.FormValue("pchannel_code")
	if pChannelCode != "" {
		if len(pChannelCode) > 50 {
			return lib.CustomError(http.StatusBadRequest, "pchannel_code harus kurang dari 255 karakter", "pchannel_code harus kurang dari 255 karakter")
		}
	}
	pchannelName := c.FormValue("pchannel_name")
	if pchannelName != "" {
		if len(pchannelName) > 150 {
			return lib.CustomError(http.StatusBadRequest, "pchannel_name should be exactly 150 characters", "pchannel_name should be exactly 150 characters")
		}
	}
	minNominalTrx := c.FormValue("min_nominal_trx")
	if minNominalTrx != "" {
		value, success := new(big.Int).SetString(minNominalTrx, 10)
		if !success {
			return lib.CustomError(http.StatusBadRequest, "min_nominal_trx must be a numeric value", "min_nominal_trx must be a numeric value")
		}
		if value.BitLen() > 18*3 { // 3 bits per digit to account for decimal places
			return lib.CustomError(http.StatusBadRequest, "min_nominal_trx should not exceed 18 digits", "min_nominal_trx should not exceed 18 digits")
		}
		params["min_nominal_trx"] = minNominalTrx
	}

	// feeValue := c.FormValue("fee_value")
	// if feeValue != "" {
	// 	// Cek apakah fee_value adalah numeric
	// 	if len(feeValue) > 18 {
	// 		return lib.CustomError(http.StatusBadRequest, "kepanjangan yang diinput", "kepanjangan yang diinput")
	// 	}
	// 	_, err := strconv.ParseFloat(feeValue, 64)
	// 	if err != nil {
	// 		return lib.CustomError(http.StatusBadRequest, "fee_value must be a numeric value", "fee_value must be a numeric value")
	// 	}
	// } else {
	// 	if feeValue == "" {
	// 		feeValue = "0"
	// 	}
	// }
	// params["fee_value"] = feeValue

	feeValue := c.FormValue("fee_value")
	valueType := c.FormValue("value_type")

	// Jika value_type adalah 316, cek fee_value
	// if valueType == "316" {
	// 	if feeValue == "" {
	// 		feeValue = "0"
	// 	} else {
	// if len(feeValue) > 18 {
	// 	return lib.CustomError(http.StatusBadRequest, "kepanjangan yang diinput", "kepanjangan yang diinput")
	// }
	if feeValue != "" {
		_, err := strconv.ParseFloat(feeValue, 64)
		if err != nil {
			return lib.CustomError(http.StatusBadRequest, "fee_value must be a numeric value", "fee_value must be a numeric value")
		}
	}
	// } else {
	// 	_, err := strconv.Atoi(feeValue)
	// 	if err != nil {
	// 		return lib.CustomError(http.StatusBadRequest, "fee_value harus harus angka", "fee_value harus harus angka")
	// 	}
	// }
	// Set nilai fee_value ke dalam params
	params["fee_value"] = feeValue

	// valueType := c.FormValue("value_type")
	// if valueType != "" {
	// 	if len(valueType) > 11 {
	// 		return lib.CustomError(http.StatusBadRequest, "value_type should be exactly 11 characters", "value_type be exactly 11 characters")
	// 	}
	// 	value, err := strconv.Atoi(valueType)
	// 	if err != nil {
	// 		return lib.CustomError(http.StatusBadRequest, "value_type should be a number", "value_type should be a number")
	// 	}
	// 	if value == 315 {

	// 	}
	// } else {
	// 	return lib.CustomError(http.StatusBadRequest, "value_type can not be blank", "value_type can not be blank")
	// }

	hasMinMax := c.FormValue("has_min_max")
	if hasMinMax == "" {
		return lib.CustomError(http.StatusBadRequest, "has_min_max can not be blank", "has_min_max can not be blank")
	}
	if hasMinMax != "" {
		feeMinValue := c.FormValue("fee_min_value")
		if feeMinValue == "0" {
			return lib.CustomError(http.StatusBadRequest, "fee_min_value cannot be empty", "fee_min_value cannot be empty")
		}
		feeMaxValue := c.FormValue("fee_max_value")
		if feeMaxValue == "" {
			return lib.CustomError(http.StatusBadRequest, "fee_max_value cannot be empty", "fee_max_value cannot be empty")
		}
	}
	settleChannelInput := c.FormValue("settle_channel")
	if settleChannelInput != "" {
		if len(settleChannelInput) > 11 {
			return lib.CustomError(http.StatusBadRequest, "settle_channel should be exactly 11 characters", "settle_channel should be exactly 11 characters")
		}
		// Validasi bahwa settleChannelInput adalah bilangan bulat
		settleChannel, err := strconv.Atoi(settleChannelInput)
		if err != nil {
			return lib.CustomError(http.StatusBadRequest, "settle_channel should be a number", "settle_channel should be a number")
		}
		params["settle_channel"] = strconv.Itoa(settleChannel)
	} else {
		return lib.CustomError(http.StatusBadRequest, "settle_channel can not be blank", "settle_channel can not be blank")
	}

	settlePaymentMethod := c.FormValue("settle_payment_method")
	if settlePaymentMethod != "" {
		if len(settlePaymentMethod) > 11 {
			return lib.CustomError(http.StatusBadRequest, "settle_payment_method should be exactly 11 characters", "settle_payment_method be exactly 11 characters")
		}
		settlepaymethod, err := strconv.Atoi(settlePaymentMethod)
		if err != nil {
			return lib.CustomError(http.StatusBadRequest, "settle_payment_method should be a number", "settle_payment_method should be a number")
		}
		params["settle_payment_method"] = strconv.Itoa(settlepaymethod)
	} else {
		return lib.CustomError(http.StatusBadRequest, "settle_payment_method can not be blank", "settle_payment_method can not be blank")
	}

	feeMinValue := c.FormValue("fee_min_value") // tanya eka

	feeMaxValue := c.FormValue("fee_max_value") // tanya eka

	fixedDmrFee := c.FormValue("fixed_dmr_fee")
	if fixedDmrFee != "" {
		value, success := new(big.Int).SetString(fixedDmrFee, 10)
		if !success {
			return lib.CustomError(http.StatusBadRequest, "fixed_dmr_fee must be a numeric value", "fixed_dmr_fee must be a numeric value")
		}
		if value.BitLen() > 18*3 { // 3 bits per digit to account for decimal places
			return lib.CustomError(http.StatusBadRequest, "fixed_dmr_fee should not exceed 18 digits", "fixed_dmr_fee should not exceed 18 digits")
		}
	} else {
		return lib.CustomError(http.StatusBadRequest, "fixed_dmr_fee can not be blank", "fixed_dmr_fee can not be blank")
	}

	fixedAmountFee := c.FormValue("fixed_amount_fee")
	if fixedAmountFee != "" {
		value, success := new(big.Int).SetString(fixedAmountFee, 10)
		if !success {
			return lib.CustomError(http.StatusBadRequest, "fixed_amount_fee must be a numeric value", "fixed_amount_fee must be a numeric value")
		}
		if value.BitLen() > 18*3 { // 3 bits per digit to account for decimal places
			return lib.CustomError(http.StatusBadRequest, "fixed_amount_fee should not exceed 18 digits", "fixed_amount_fee should not exceed 18 digits")
		}
	} else {
		return lib.CustomError(http.StatusBadRequest, "fixed_amount_fee can not be blank", "fixed_amount_fee can not be blank")
	}
	pgTnc := c.FormValue("pg_tnc")

	pgRemarks := c.FormValue("pg_remarks")

	paymentLoginUrl := c.FormValue("payment_login_url")
	if len(paymentLoginUrl) > 255 {
		return lib.CustomError(http.StatusBadRequest, "fixed_amount_fee can not be blank", "fixed_amount_fee can not be blank")
	}
	paymentEntryUrl := c.FormValue("payment_entry_url")
	if len(paymentEntryUrl) > 255 {
		return lib.CustomError(http.StatusBadRequest, "payment_entry_url can not be blank", "payment_entry_url can not be blank")
	}
	paymentErrorUrl := c.FormValue("payment_error_url")
	if len(paymentErrorUrl) > 255 {
		return lib.CustomError(http.StatusBadRequest, "payment_error_url can not be blank", "payment_error_url can not be blank")
	}
	paymentSuccessUrl := c.FormValue("payment_success_url")
	if len(paymentSuccessUrl) > 255 {
		return lib.CustomError(http.StatusBadRequest, "payment_success_url can not be blank", "payment_success_url can not be blank")
	}
	pgPrefix := c.FormValue("pg_prefix")
	if len(pgPrefix) > 150 {
		return lib.CustomError(http.StatusBadRequest, "pg_prefix can not be blank", "pg_prefix can not be blank")
	}
	picName := c.FormValue("pic_name")
	if len(picName) > 150 {
		return lib.CustomError(http.StatusBadRequest, "pic_name can not be blank", "pic_name can not be blank")
	}
	picPhoneNo := c.FormValue("pic_phone_no")
	if len(picPhoneNo) > 150 {
		return lib.CustomError(http.StatusBadRequest, "pic_phone_no can not be blank", "pic_phone_no can not be blank")
	}
	picEmailAddress := c.FormValue("pic_email_address")
	if len(picEmailAddress) > 150 {
		return lib.CustomError(http.StatusBadRequest, "pic_email_address can not be blank", "pic_email_address can not be blank")
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
	}
	params["rec_order"] = recOrder

	params["pchannel_code"] = pChannelCode
	params["pchannel_name"] = pchannelName
	params["pchannel_name"] = pchannelName
	params["fee_value"] = feeValue
	params["has_min_max"] = hasMinMax
	params["settle_payment_method"] = settlePaymentMethod
	params["value_type"] = valueType
	params["fee_min_value"] = feeMinValue
	params["fee_max_value"] = feeMaxValue
	params["fixed_dmr_fee"] = fixedDmrFee
	params["fixed_amount_fee"] = fixedAmountFee
	params["pg_tnc"] = pgTnc
	params["pg_remarks"] = pgRemarks
	params["payment_login_url"] = paymentLoginUrl
	params["payment_entry_url"] = paymentEntryUrl
	params["payment_error_url"] = paymentErrorUrl
	params["payment_success_url"] = paymentSuccessUrl
	params["pg_prefix"] = pgPrefix
	params["pic_name"] = picName
	params["pic_phone_no"] = picPhoneNo
	params["pic_email_address"] = picEmailAddress
	params["rec_status"] = "1"
	status, err = models.CreateMsPaymentChannel(params)
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

func UpdateMsPaymentChannelController(c echo.Context) error {
	var err error
	params := make(map[string]string)
	params["rec_modified_by"] = lib.UserIDStr
	params["rec_modified_date"] = time.Now().Format(lib.TIMESTAMPFORMAT)

	pChannelKey := c.FormValue("pchannel_key")
	if pChannelKey == "" {
		return lib.CustomError(http.StatusBadRequest, "pchannel_key can not be blank", "pchannel_key can not be blank")
	}
	pChannelCode := c.FormValue("pchannel_code")
	if pChannelCode != "" {
		if len(pChannelCode) > 50 {
			return lib.CustomError(http.StatusBadRequest, "pchannel_code harus kurang dari 255 karakter", "pchannel_code harus kurang dari 255 karakter")
		}
	}
	pchannelName := c.FormValue("pchannel_name")
	if pchannelName != "" {
		if len(pchannelName) > 150 {
			return lib.CustomError(http.StatusBadRequest, "pchannel_name should be exactly 150 characters", "pchannel_name should be exactly 150 characters")
		}
	}
	minNominalTrx := c.FormValue("min_nominal_trx")
	if minNominalTrx != "" {
		value, success := new(big.Int).SetString(minNominalTrx, 10)
		if !success {
			return lib.CustomError(http.StatusBadRequest, "min_nominal_trx must be a numeric value", "min_nominal_trx must be a numeric value")
		}
		if value.BitLen() > 18*3 { // 3 bits per digit to account for decimal places
			return lib.CustomError(http.StatusBadRequest, "min_nominal_trx should not exceed 18 digits", "min_nominal_trx should not exceed 18 digits")
		}
		params["min_nominal_trx"] = minNominalTrx
	}

	// feeValue := c.FormValue("fee_value")
	// if feeValue != "" {
	// 	// Cek apakah fee_value adalah numeric
	// 	if len(feeValue) > 18 {
	// 		return lib.CustomError(http.StatusBadRequest, "kepanjangan yang diinput", "kepanjangan yang diinput")
	// 	}
	// 	_, err := strconv.ParseFloat(feeValue, 64)
	// 	if err != nil {
	// 		return lib.CustomError(http.StatusBadRequest, "fee_value must be a numeric value", "fee_value must be a numeric value")
	// 	}
	// } else {
	// 	if feeValue == "" {
	// 		feeValue = "0"
	// 	}
	// }
	// params["fee_value"] = feeValue

	feeValue := c.FormValue("fee_value")
	valueType := c.FormValue("value_type")

	// Jika value_type adalah 316, cek fee_value
	if valueType == "316" {
		if feeValue == "" {
			feeValue = "0"
		} else {
			if len(feeValue) > 18 {
				return lib.CustomError(http.StatusBadRequest, "kepanjangan yang diinput", "kepanjangan yang diinput")
			}
			_, err := strconv.ParseFloat(feeValue, 64)
			if err != nil {
				return lib.CustomError(http.StatusBadRequest, "fee_value must be a numeric value", "fee_value must be a numeric value")
			}
		}
	} else {
		_, err := strconv.Atoi(feeValue)
		if err != nil {
			return lib.CustomError(http.StatusBadRequest, "fee_value harus harus angka", "fee_value harus harus angka")
		}
	}
	// Set nilai fee_value ke dalam params
	params["fee_value"] = feeValue

	// valueType := c.FormValue("value_type")
	// if valueType != "" {
	// 	if len(valueType) > 11 {
	// 		return lib.CustomError(http.StatusBadRequest, "value_type should be exactly 11 characters", "value_type be exactly 11 characters")
	// 	}
	// 	value, err := strconv.Atoi(valueType)
	// 	if err != nil {
	// 		return lib.CustomError(http.StatusBadRequest, "value_type should be a number", "value_type should be a number")
	// 	}
	// 	if value == 315 {

	// 	}
	// } else {
	// 	return lib.CustomError(http.StatusBadRequest, "value_type can not be blank", "value_type can not be blank")
	// }

	hasMinMax := c.FormValue("has_min_max")
	if hasMinMax == "" {
		return lib.CustomError(http.StatusBadRequest, "has_min_max can not be blank", "has_min_max can not be blank")
	}
	if hasMinMax != "" {
		feeMinValue := c.FormValue("fee_min_value")
		if feeMinValue == "0" {
			return lib.CustomError(http.StatusBadRequest, "fee_min_value cannot be empty", "fee_min_value cannot be empty")
		}
		feeMaxValue := c.FormValue("fee_max_value")
		if feeMaxValue == "" {
			return lib.CustomError(http.StatusBadRequest, "fee_max_value cannot be empty", "fee_max_value cannot be empty")
		}
	}
	settleChannelInput := c.FormValue("settle_channel")
	if settleChannelInput != "" {
		if len(settleChannelInput) > 11 {
			return lib.CustomError(http.StatusBadRequest, "settle_channel should be exactly 11 characters", "settle_channel should be exactly 11 characters")
		}
		// Validasi bahwa settleChannelInput adalah bilangan bulat
		settleChannel, err := strconv.Atoi(settleChannelInput)
		if err != nil {
			return lib.CustomError(http.StatusBadRequest, "settle_channel should be a number", "settle_channel should be a number")
		}
		params["settle_channel"] = strconv.Itoa(settleChannel)
	} else {
		return lib.CustomError(http.StatusBadRequest, "settle_channel can not be blank", "settle_channel can not be blank")
	}

	settlePaymentMethod := c.FormValue("settle_payment_method")
	if settlePaymentMethod != "" {
		if len(settlePaymentMethod) > 11 {
			return lib.CustomError(http.StatusBadRequest, "settle_payment_method should be exactly 11 characters", "settle_payment_method be exactly 11 characters")
		}
		settlepaymethod, err := strconv.Atoi(settlePaymentMethod)
		if err != nil {
			return lib.CustomError(http.StatusBadRequest, "settle_payment_method should be a number", "settle_payment_method should be a number")
		}
		params["settle_payment_method"] = strconv.Itoa(settlepaymethod)
	} else {
		return lib.CustomError(http.StatusBadRequest, "settle_payment_method can not be blank", "settle_payment_method can not be blank")
	}

	feeMinValue := c.FormValue("fee_min_value") // tanya eka

	feeMaxValue := c.FormValue("fee_max_value") // tanya eka

	fixedDmrFee := c.FormValue("fixed_dmr_fee")
	if fixedDmrFee != "" {
		value, success := new(big.Int).SetString(fixedDmrFee, 10)
		if !success {
			return lib.CustomError(http.StatusBadRequest, "fixed_dmr_fee must be a numeric value", "fixed_dmr_fee must be a numeric value")
		}
		if value.BitLen() > 18*3 { // 3 bits per digit to account for decimal places
			return lib.CustomError(http.StatusBadRequest, "fixed_dmr_fee should not exceed 18 digits", "fixed_dmr_fee should not exceed 18 digits")
		}
	} else {
		return lib.CustomError(http.StatusBadRequest, "fixed_dmr_fee can not be blank", "fixed_dmr_fee can not be blank")
	}

	fixedAmountFee := c.FormValue("fixed_amount_fee")
	if fixedAmountFee != "" {
		value, success := new(big.Int).SetString(fixedAmountFee, 10)
		if !success {
			return lib.CustomError(http.StatusBadRequest, "fixed_amount_fee must be a numeric value", "fixed_amount_fee must be a numeric value")
		}
		if value.BitLen() > 18*3 { // 3 bits per digit to account for decimal places
			return lib.CustomError(http.StatusBadRequest, "fixed_amount_fee should not exceed 18 digits", "fixed_amount_fee should not exceed 18 digits")
		}
	} else {
		return lib.CustomError(http.StatusBadRequest, "fixed_amount_fee can not be blank", "fixed_amount_fee can not be blank")
	}
	pgTnc := c.FormValue("pg_tnc")

	pgRemarks := c.FormValue("pg_remarks")

	paymentLoginUrl := c.FormValue("payment_login_url")
	if len(paymentLoginUrl) > 255 {
		return lib.CustomError(http.StatusBadRequest, "fixed_amount_fee can not be blank", "fixed_amount_fee can not be blank")
	}
	paymentEntryUrl := c.FormValue("payment_entry_url")
	if len(paymentEntryUrl) > 255 {
		return lib.CustomError(http.StatusBadRequest, "payment_entry_url can not be blank", "payment_entry_url can not be blank")
	}
	paymentErrorUrl := c.FormValue("payment_error_url")
	if len(paymentErrorUrl) > 255 {
		return lib.CustomError(http.StatusBadRequest, "payment_error_url can not be blank", "payment_error_url can not be blank")
	}
	paymentSuccessUrl := c.FormValue("payment_success_url")
	if len(paymentSuccessUrl) > 255 {
		return lib.CustomError(http.StatusBadRequest, "payment_success_url can not be blank", "payment_success_url can not be blank")
	}
	pgPrefix := c.FormValue("pg_prefix")
	if len(pgPrefix) > 150 {
		return lib.CustomError(http.StatusBadRequest, "pg_prefix can not be blank", "pg_prefix can not be blank")
	}
	picName := c.FormValue("pic_name")
	if len(picName) > 150 {
		return lib.CustomError(http.StatusBadRequest, "pic_name can not be blank", "pic_name can not be blank")
	}
	picPhoneNo := c.FormValue("pic_phone_no")
	if len(picPhoneNo) > 150 {
		return lib.CustomError(http.StatusBadRequest, "pic_phone_no can not be blank", "pic_phone_no can not be blank")
	}
	picEmailAddress := c.FormValue("pic_email_address")
	if len(picEmailAddress) > 150 {
		return lib.CustomError(http.StatusBadRequest, "pic_email_address can not be blank", "pic_email_address can not be blank")
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
	}
	params["rec_order"] = recOrder

	params["pchannel_code"] = pChannelCode
	params["pchannel_name"] = pchannelName
	params["pchannel_name"] = pchannelName
	params["fee_value"] = feeValue
	params["has_min_max"] = hasMinMax
	params["settle_payment_method"] = settlePaymentMethod
	params["value_type"] = valueType
	params["fee_min_value"] = feeMinValue
	params["fee_max_value"] = feeMaxValue
	params["fixed_dmr_fee"] = fixedDmrFee
	params["fixed_amount_fee"] = fixedAmountFee
	params["pg_tnc"] = pgTnc
	params["pg_remarks"] = pgRemarks
	params["payment_login_url"] = paymentLoginUrl
	params["payment_entry_url"] = paymentEntryUrl
	params["payment_error_url"] = paymentErrorUrl
	params["payment_success_url"] = paymentSuccessUrl
	params["pg_prefix"] = pgPrefix
	params["pic_name"] = picName
	params["pic_phone_no"] = picPhoneNo
	params["pic_email_address"] = picEmailAddress
	params["rec_status"] = "1"

	status, err = models.UpdateMsPaymentChannel(pChannelKey, params)
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
