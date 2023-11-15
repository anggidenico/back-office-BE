package controllers

import (
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
		_, err := strconv.Atoi(minNominalTrx)
		if err != nil {
			return lib.CustomError(http.StatusBadRequest, "min_nominal_trx harus berupa angka", "min_nominal_trx harus berupa angka")
		}
		params["min_nominal_trx"] = minNominalTrx
	} else {
		params["min_nominal_trx"] = "0"
	}

	feeValue := c.FormValue("fee_value")
	if feeValue == "" {
		return lib.CustomError(http.StatusBadRequest, "fee_value can not be blank", "fee_value can not be blank")
	}
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
	if settlePaymentMethod == "" {
		return lib.CustomError(http.StatusBadRequest, "settle_payment_method can not be blank", "settle_payment_method can not be blank")
	}
	valueType := c.FormValue("value_type")
	if valueType == "" {
		return lib.CustomError(http.StatusBadRequest, "settle_channel can not be blank", "settle_channel can not be blank")
	}
	feeMinValue := c.FormValue("fee_min_value")

	feeMaxValue := c.FormValue("fee_max_value")

	fixedDmrFee := c.FormValue("fixed_dmr_fee")

	fixedAmountFee := c.FormValue("fixed_amount_fee")

	pgTnc := c.FormValue("pg_tnc")
	if pgTnc != "" {
		if len(pgTnc) > 5000 {
			return lib.CustomError(http.StatusBadRequest, "pg_tnc harus kurang dari 255 karakter", "pg_tnc harus kurang dari 255 karakter")
		}
	}

	pgRemarks := c.FormValue("pg_remarks")

	paymentLoginUrl := c.FormValue("payment_login_url")

	paymentEntryUrl := c.FormValue("payment_entry_url")

	paymentErrorUrl := c.FormValue("payment_error_url")

	paymentSuccessUrl := c.FormValue("payment_success_url")

	pgPrefix := c.FormValue("pg_prefix")

	picName := c.FormValue("pic_name")

	picPhoneNo := c.FormValue("pic_phone_no")

	picEmailAddress := c.FormValue("pic_email_address")

	recOrder := c.FormValue("rec_order")
	if recOrder == "" {
		return lib.CustomError(http.StatusBadRequest, "rec_order can not be blank", "rec_order can not be blank")
	}
	params["pchannel_code"] = pChannelCode
	params["pchannel_name"] = pchannelName
	// params["min_nominal_trx"] = minNominalTrx
	params["pchannel_name"] = pchannelName
	params["fee_value"] = feeValue
	params["has_min_max"] = hasMinMax
	// params["settle_channel"] = SettleChannel
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
	if pChannelCode == "" {
		return lib.CustomError(http.StatusBadRequest, "pchannel_code can not be blank", "pchannel_code can not be blank")
	}
	pchannelName := c.FormValue("pchannel_name")
	if pchannelName == "" {
		return lib.CustomError(http.StatusBadRequest, "pchannel_name can not be blank", "pchannel_name can not be blank")
	}
	minNominalTrx := c.FormValue("min_nominal_trx")
	if minNominalTrx == "" {
		return lib.CustomError(http.StatusBadRequest, "min_nominal_trx can not be blank", "min_nominal_trx can not be blank")
	}
	feeValue := c.FormValue("fee_value")
	if feeValue == "" {
		return lib.CustomError(http.StatusBadRequest, "fee_value can not be blank", "fee_value can not be blank")
	}
	hasMinMax := c.FormValue("has_min_max")
	if hasMinMax == "" {
		return lib.CustomError(http.StatusBadRequest, "has_min_max can not be blank", "has_min_max can not be blank")
	}
	settleChannel := c.FormValue("settle_channel")
	if settleChannel == "" {
		return lib.CustomError(http.StatusBadRequest, "settle_channel can not be blank", "settle_channel can not be blank")
	}
	settlePaymentMethod := c.FormValue("settle_payment_method")
	if settlePaymentMethod == "" {
		return lib.CustomError(http.StatusBadRequest, "settle_payment_method can not be blank", "settle_payment_method can not be blank")
	}
	valueType := c.FormValue("value_type")
	if valueType == "" {
		return lib.CustomError(http.StatusBadRequest, "settle_channel can not be blank", "settle_channel can not be blank")
	}
	feeMinValue := c.FormValue("fee_min_value")
	if feeMinValue == "" {
		return lib.CustomError(http.StatusBadRequest, "fee_min_value can not be blank", "fee_min_value can not be blank")
	}
	feeMaxValue := c.FormValue("fee_max_value")
	if feeMaxValue == "" {
		return lib.CustomError(http.StatusBadRequest, "fee_max_value can not be blank", "fee_max_value can not be blank")
	}
	fixedDmrFee := c.FormValue("fixed_dmr_fee")
	if fixedDmrFee == "" {
		return lib.CustomError(http.StatusBadRequest, "fixed_dmr_fee can not be blank", "fixed_dmr_fee can not be blank")
	}
	fixedAmountFee := c.FormValue("fixed_amount_fee")
	if fixedAmountFee == "" {
		return lib.CustomError(http.StatusBadRequest, "fee_max_value can not be blank", "fee_max_value can not be blank")
	}
	pgTnc := c.FormValue("pg_tnc")
	if pgTnc == "" {
		return lib.CustomError(http.StatusBadRequest, "pg_tnc can not be blank", "pg_tnc can not be blank")
	}
	pgRemarks := c.FormValue("pg_remarks")
	if pgRemarks == "" {
		return lib.CustomError(http.StatusBadRequest, "pg_remarks can not be blank", "pg_remarks can not be blank")
	}
	paymentLoginUrl := c.FormValue("payment_login_url")
	if paymentLoginUrl == "" {
		return lib.CustomError(http.StatusBadRequest, "payment_login_url can not be blank", "payment_login_url can not be blank")
	}
	paymentEntryUrl := c.FormValue("payment_entry_url")
	if paymentEntryUrl == "" {
		return lib.CustomError(http.StatusBadRequest, "payment_entry_url can not be blank", "payment_entry_url can not be blank")
	}
	paymentErrorUrl := c.FormValue("payment_error_url")
	if paymentErrorUrl == "" {
		return lib.CustomError(http.StatusBadRequest, "payment_error_url can not be blank", "payment_error_url can not be blank")
	}
	paymentSuccessUrl := c.FormValue("payment_success_url")
	if paymentSuccessUrl == "" {
		return lib.CustomError(http.StatusBadRequest, "payment_success_url can not be blank", "payment_success_url can not be blank")
	}
	pgPrefix := c.FormValue("pg_prefix")
	if pgPrefix == "" {
		return lib.CustomError(http.StatusBadRequest, "pg_prefix can not be blank", "pg_prefix can not be blank")
	}
	picName := c.FormValue("pic_name")
	if picName == "" {
		return lib.CustomError(http.StatusBadRequest, "pic_name can not be blank", "pic_name can not be blank")
	}
	picPhoneNo := c.FormValue("pic_phone_no")
	if picPhoneNo == "" {
		return lib.CustomError(http.StatusBadRequest, "pic_phone_no can not be blank", "pic_phone_no can not be blank")
	}
	picEmailAddress := c.FormValue("pic_email_address")
	if picEmailAddress == "" {
		return lib.CustomError(http.StatusBadRequest, "pic_email_address can not be blank", "pic_email_address can not be blank")
	}
	recOrder := c.FormValue("rec_order")
	if recOrder == "" {
		return lib.CustomError(http.StatusBadRequest, "rec_order can not be blank", "rec_order can not be blank")
	}
	params["pchannel_key"] = pChannelKey
	params["pchannel_code"] = pChannelCode
	params["pchannel_name"] = pchannelName
	params["min_nominal_trx"] = minNominalTrx
	params["pchannel_name"] = pchannelName
	params["fee_value"] = feeValue
	params["has_min_max"] = hasMinMax
	params["settle_channel"] = settleChannel
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
