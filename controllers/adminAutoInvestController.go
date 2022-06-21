package controllers

import (
	"bytes"
	"database/sql"
	"html/template"
	"math"
	"mf-bo-api/config"
	"mf-bo-api/lib"
	"mf-bo-api/models"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo"
	"github.com/leekchan/accounting"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

func ListAdminAutoInvest(c echo.Context) error {

	var err error
	var status int
	decimal.MarshalJSONWithoutQuotes = true

	//Get parameter limit
	limitStr := c.QueryParam("limit")
	var limit uint64
	if limitStr != "" {
		limit, err = strconv.ParseUint(limitStr, 10, 64)
		if err == nil {
			if (limit == 0) || (limit > config.LimitQuery) {
				limit = config.LimitQuery
			}
		} else {
			log.Error("Limit should be number")
			return lib.CustomError(http.StatusBadRequest, "Limit should be number", "Limit should be number")
		}
	} else {
		limit = config.LimitQuery
	}
	// Get parameter page
	pageStr := c.QueryParam("page")
	var page uint64
	if pageStr != "" {
		page, err = strconv.ParseUint(pageStr, 10, 64)
		if err == nil {
			if page == 0 {
				page = 1
			}
		} else {
			log.Error("Page should be number")
			return lib.CustomError(http.StatusBadRequest, "Page should be number", "Page should be number")
		}
	} else {
		page = 1
	}
	var offset uint64
	if page > 1 {
		offset = limit * (page - 1)
	}

	noLimitStr := c.QueryParam("nolimit")
	var noLimit bool
	if noLimitStr != "" {
		noLimit, err = strconv.ParseBool(noLimitStr)
		if err != nil {
			log.Error("Nolimit parameter should be true/false")
			return lib.CustomError(http.StatusBadRequest, "Nolimit parameter should be true/false", "Nolimit parameter should be true/false")
		}
	} else {
		noLimit = false
	}

	items := []string{"full_name", "product_name", "invest_amount", "invest_date_execute", "date_thru", "date_last_generate", "settle_channel", "settle_payment_method", "bank_name", "account_no"}

	params := make(map[string]string)
	orderBy := c.QueryParam("order_by")
	if orderBy != "" {
		_, found := lib.Find(items, orderBy)
		if found {
			var ord string
			if orderBy == "full_name" {
				ord = "c.full_name"
			} else if orderBy == "product_name" {
				ord = "p.product_name_alt"
			} else if orderBy == "invest_amount" {
				ord = "a.invest_amount"
			} else if orderBy == "invest_date_execute" {
				ord = "a.invest_date_execute"
			} else if orderBy == "date_thru" {
				ord = "a.date_thru"
			} else if orderBy == "date_last_generate" {
				ord = "a.date_last_generate"
			} else if orderBy == "settle_channel" {
				ord = "setchannel.lkp_name"
			} else if orderBy == "settle_payment_method" {
				ord = "setpayment.lkp_name"
			} else if orderBy == "bank_name" {
				ord = "b.bank_name"
			} else if orderBy == "account_no" {
				ord = "ba.account_no"
			} else {
				ord = "a." + orderBy
			}

			params["orderBy"] = ord
			orderType := c.QueryParam("order_type")
			if (orderType == "asc") || (orderType == "ASC") || (orderType == "desc") || (orderType == "DESC") {
				params["orderType"] = orderType
			}
		} else {
			log.Error("Wrong input for parameter order_by")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter order_by", "Wrong input for parameter order_by")
		}
	} else {
		params["orderBy"] = "a.date_last_generate"
		params["orderType"] = "ASC"
	}

	productKey := c.QueryParam("product_key")
	if productKey != "" {
		n, err := strconv.ParseUint(productKey, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: product_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: product_key", "Wrong input for parameter: product_key")
		}

		if len(productKey) > 11 {
			log.Error("Wrong input for parameter: product_key too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: product_key too long, max 11 character", "Missing required parameter: product_key too long, max 11 character")
		}
		params["ta.product_key"] = productKey
	}

	customerKey := c.QueryParam("customer_key")
	if customerKey != "" {
		n, err := strconv.ParseUint(customerKey, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: customer_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: customer_key", "Wrong input for parameter: customer_key")
		}

		if len(customerKey) > 11 {
			log.Error("Wrong input for parameter: customer_key too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: customer_key too long, max 11 character", "Missing required parameter: customer_key too long, max 11 character")
		}
		params["ta.customer_key"] = customerKey
	}

	var trAutoinvest []models.AdminListAutoInvestRegistration
	status, err = models.GetAdminListAutoInvestRegistration(&trAutoinvest, params, limit, offset, noLimit)

	if err != nil {
		if err != sql.ErrNoRows {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
	}
	if len(trAutoinvest) < 1 {
		log.Error("Data not found")
		return lib.CustomError(http.StatusNotFound, "Data not found", "Data not found")
	}

	var countData models.CountData
	var pagination int
	if limit > 0 {
		status, err = models.GetAdminCountListAutoInvestRegistration(&countData, params)
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
		if int(countData.CountData) < int(limit) {
			pagination = 1
		} else {
			calc := math.Ceil(float64(countData.CountData) / float64(limit))
			pagination = int(calc)
		}
	} else {
		pagination = 1
	}

	var response lib.ResponseWithPagination
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Pagination = pagination
	response.Data = trAutoinvest

	return c.JSON(http.StatusOK, response)
}

func AdminCreateTrAutoInvest(c echo.Context) error {
	var err error
	var status int

	params := make(map[string]string)
	paramsAcc := make(map[string]string)

	productKey := c.FormValue("product_key")
	if productKey != "" {
		n, err := strconv.ParseUint(productKey, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: product_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: product_key", "Wrong input for parameter: product_key")
		}

		if len(productKey) > 11 {
			log.Error("Wrong input for parameter: product_key too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: product_key too long, max 11 character", "Missing required parameter: product_key too long, max 11 character")
		}
		paramsAcc["product_key"] = productKey
	} else {
		log.Error("Missing required parameter: product_key")
		return lib.CustomError(http.StatusBadRequest, "product_key can not be blank", "product_key can not be blank")
	}

	var product models.MsProduct
	status, err = models.GetMsProduct(&product, productKey)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(http.StatusBadRequest, err.Error(), "Product tidak ditemukan")
	}

	if product.FlagEnabled != 1 || product.FlagSubscription != 1 {
		log.Error("Tidak dapat melakukan autoinvest pada produk ini.")
		return lib.CustomError(http.StatusBadRequest, "Tidak dapat melakukan autoinvest pada produk ini.", "Tidak dapat melakukan autoinvest pada produk ini.")
	}

	customerKey := c.FormValue("customer_key")
	if customerKey != "" {
		n, err := strconv.ParseUint(customerKey, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: customer_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: customer_key", "Wrong input for parameter: customer_key")
		}

		if len(customerKey) > 11 {
			log.Error("Wrong input for parameter: customer_key too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: customer_key too long, max 11 character", "Missing required parameter: customer_key too long, max 11 character")
		}
		paramsAcc["customer_key"] = customerKey
	} else {
		log.Error("Missing required parameter: customer_key")
		return lib.CustomError(http.StatusBadRequest, "customer_key can not be blank", "customer_key can not be blank")
	}

	var customer models.MsCustomer
	status, err = models.GetMsCustomer(&customer, customerKey)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
	} else {
		if customer.CifSuspendFlag == 1 {
			log.Error("Account customer tersuspend. CIF Suspend")
			return lib.CustomError(http.StatusBadRequest, "Customer Account Suspended", "Customer Account Suspended")

		}
	}

	investAmount := c.FormValue("invest_amount")
	if investAmount != "" {
		value, err := decimal.NewFromString(investAmount)
		if err != nil {
			log.Error("Wrong input for parameter: invest_amount")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: invest_amount", "Wrong input for parameter: invest_amount")
		}
		if value.Cmp(product.MinSubAmount) == -1 {
			log.Error("invest_amount < minimum sub")
			return lib.CustomError(http.StatusBadRequest, "invest_amount < minum sub", "Minumum invest_amount untuk product ini adalah: "+product.MinSubAmount.String())
		}
		if investAmount == "0" {
			log.Error("Wrong input for parameter: invest_amount")
			return lib.CustomError(http.StatusBadRequest, "invest_amount harus lebih dari 0", "invest_amount harus lebih dari 0")
		}
		params["invest_amount"] = investAmount
	} else {
		log.Error("Missing required parameter: invest_amount")
		return lib.CustomError(http.StatusBadRequest, "invest_amount can not be blank", "invest_amount can not be blank")
	}

	investFeeRates := c.FormValue("invest_fee_rates")
	if investFeeRates != "" {
		_, err := decimal.NewFromString(investFeeRates)
		if err != nil {
			log.Error("Wrong input for parameter: invest_fee_rates")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: invest_fee_rates", "Wrong input for parameter: invest_fee_rates")
		}
		params["invest_fee_rates"] = investFeeRates
	} else {
		log.Error("Missing required parameter: invest_fee_rates")
		return lib.CustomError(http.StatusBadRequest, "invest_fee_rates can not be blank", "invest_fee_rates can not be blank")
	}

	investFeeAmount := c.FormValue("invest_fee_amount")
	if investFeeAmount != "" {
		_, err := decimal.NewFromString(investFeeAmount)
		if err != nil {
			log.Error("Wrong input for parameter: invest_fee_amount")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: invest_fee_amount", "Wrong input for parameter: invest_fee_amount")
		}
		params["invest_fee_amount"] = investFeeAmount
	} else {
		log.Error("Missing required parameter: invest_fee_amount")
		return lib.CustomError(http.StatusBadRequest, "invest_fee_amount can not be blank", "invest_fee_amount can not be blank")
	}

	investFeeCharges := c.FormValue("invest_fee_charges")
	if investFeeCharges != "" {
		_, err := decimal.NewFromString(investFeeCharges)
		if err != nil {
			log.Error("Wrong input for parameter: invest_fee_charges")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: invest_fee_charges", "Wrong input for parameter: invest_fee_charges")
		}
		params["invest_fee_charges"] = investFeeCharges
	} else {
		log.Error("Missing required parameter: invest_fee_charges")
		return lib.CustomError(http.StatusBadRequest, "invest_fee_charges can not be blank", "invest_fee_charges can not be blank")
	}

	investDateExecute := c.FormValue("invest_date_execute")
	if investDateExecute != "" {
		n, err := strconv.ParseUint(investDateExecute, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: invest_date_execute")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: invest_date_execute", "Wrong input for parameter: invest_date_execute")
		}
		if n > 28 {
			log.Error("invest_date_execute max tanggal 28")
			return lib.CustomError(http.StatusBadRequest, "invest_date_execute max tanggal 28", "invest_date_execute max tanggal 28")
		}
		params["invest_date_execute"] = investDateExecute
	} else {
		log.Error("Missing required parameter: invest_date_execute")
		return lib.CustomError(http.StatusBadRequest, "invest_date_execute can not be blank", "invest_date_execute can not be blank")
	}

	dateThru := c.FormValue("date_thru")
	if dateThru != "" {
		params["date_thru"] = dateThru
	} else {
		log.Error("Missing required parameter: date_thru")
		return lib.CustomError(http.StatusBadRequest, "date_thru can not be blank", "date_thru can not be blank")
	}

	settleChannel := c.FormValue("settle_channel")
	if settleChannel != "" {
		n, err := strconv.ParseUint(settleChannel, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: settle_channel")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: settle_channel", "Wrong input for parameter: settle_channel")
		}

		params["settle_channel"] = settleChannel
	} else {
		log.Error("Missing required parameter: settle_channel")
		return lib.CustomError(http.StatusBadRequest, "settle_channel can not be blank", "settle_channel can not be blank")
	}

	settlePaymentMethod := c.FormValue("settle_payment_method")
	if settlePaymentMethod != "" {
		n, err := strconv.ParseUint(settlePaymentMethod, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: settle_payment_method")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: settle_payment_method", "Wrong input for parameter: settle_payment_method")
		}

		params["settle_payment_method"] = settlePaymentMethod
	} else {
		log.Error("Missing required parameter: settle_payment_method")
		return lib.CustomError(http.StatusBadRequest, "settle_payment_method can not be blank", "settle_payment_method can not be blank")
	}

	prodBankAccKey := c.FormValue("prod_bank_acc_key")
	if prodBankAccKey != "" {
		n, err := strconv.ParseUint(prodBankAccKey, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: prod_bank_acc_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: prod_bank_acc_key", "Wrong input for parameter: prod_bank_acc_key")
		}

	} else {
		log.Error("Missing required parameter: prod_bank_acc_key")
		return lib.CustomError(http.StatusBadRequest, "prod_bank_acc_key can not be blank", "prod_bank_acc_key can not be blank")
	}

	investRemarks := c.FormValue("invest_remarks")
	if investRemarks != "" {
		if len(investRemarks) > 140 {
			log.Error("Wrong input for parameter: invest_remarks too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: invest_remarks too long, max 140 character", "Missing required parameter: invest_remarks too long, max 140 character")
		}
		params["invest_remarks"] = investRemarks
	}

	dateLayout := "2006-01-02 15:04:05"
	var accKey string
	paramsAcc["rec_status"] = "1"
	var trAccountDB []models.TrAccount
	status, err = models.GetAllTrAccount(&trAccountDB, paramsAcc)
	if len(trAccountDB) > 0 {
		accKey = strconv.FormatUint(trAccountDB[0].AccKey, 10)
		if trAccountDB[0].SubSuspendFlag != nil && *trAccountDB[0].SubSuspendFlag == 1 {
			log.Error("Account suspended for this product")
			return lib.CustomError(status, "Account suspended for this product", "Account suspended for this product")
		}
	} else {
		paramsAcc["acc_status"] = "204"
		paramsAcc["rec_created_date"] = time.Now().Format(dateLayout)
		paramsAcc["rec_created_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
		status, err, accKey = models.CreateTrAccount(paramsAcc)
		if err != nil {
			log.Error("Failed create account product data: " + err.Error())
			return lib.CustomError(status, err.Error(), "failed input data")
		}
	}

	var countData models.CountData
	status, err = models.AdminValidateAccAndiInvestDateExecute(&countData, accKey, investDateExecute, "")
	if err != nil {
		log.Error(err.Error())
	}
	if int(countData.CountData) > int(0) {
		log.Error("Ada data product dengan execute date date yang sama")
		return lib.CustomError(status, "Terdapat data product dengan execute date date yang sama", "Terdapat data product dengan execute date date yang sama")
	}

	params["acc_key"] = accKey
	params["attempt_count"] = "0"
	params["date_start"] = time.Now().Format(dateLayout)
	params["rec_created_date"] = time.Now().Format(dateLayout)
	params["rec_created_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
	params["rec_status"] = "1"

	var productBankAccount models.MsProductBankAccount
	status, err = models.GetMsProductBankAccount(&productBankAccount, prodBankAccKey)
	if err != nil {
		log.Error("Failed get product bank account: " + err.Error())
	} else {
		params["bank_account_key"] = strconv.FormatUint(*productBankAccount.BankAccountKey, 10)
		var bankAccount models.MsBankAccount
		status, err = models.GetBankAccount(&bankAccount, strconv.FormatUint(*productBankAccount.BankAccountKey, 10))
		if err != nil {
			log.Error("Failed get bank account: " + err.Error())
		} else {
			params["bank_key"] = strconv.FormatUint(bankAccount.BankKey, 10)
		}
	}

	_, err, _ = models.CreateTrAutoinvestRegistration(params)
	if err != nil {
		log.Error("Failed create tr_autoinvest_registration: " + err.Error())
		return lib.CustomError(status, err.Error(), "failed input data")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = ""
	return c.JSON(http.StatusOK, response)
}

func AdminUpdateTrAutoInvest(c echo.Context) error {
	var err error
	var status int

	params := make(map[string]string)
	paramsAcc := make(map[string]string)

	autoinvestKey := c.FormValue("autoinvest_key")
	if autoinvestKey != "" {
		n, err := strconv.ParseUint(autoinvestKey, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: autoinvest_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: autoinvest_key", "Wrong input for parameter: autoinvest_key")
		}

		if len(autoinvestKey) > 11 {
			log.Error("Wrong input for parameter: autoinvest_key too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: autoinvest_key too long, max 11 character", "Missing required parameter: autoinvest_key too long, max 11 character")
		}
		params["autoinvest_key"] = autoinvestKey
	} else {
		log.Error("Missing required parameter: autoinvest_key")
		return lib.CustomError(http.StatusBadRequest, "autoinvest_key can not be blank", "autoinvest_key can not be blank")
	}

	var invest models.TrAutoinvestRegistration
	_, err = models.GetTrAutoinvestRegistration(&invest, autoinvestKey)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(http.StatusBadRequest, "failed get data", "failed get data")
	}

	productKey := c.FormValue("product_key")
	if productKey != "" {
		n, err := strconv.ParseUint(productKey, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: product_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: product_key", "Wrong input for parameter: product_key")
		}

		if len(productKey) > 11 {
			log.Error("Wrong input for parameter: product_key too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: product_key too long, max 11 character", "Missing required parameter: product_key too long, max 11 character")
		}
		paramsAcc["product_key"] = productKey
	} else {
		log.Error("Missing required parameter: product_key")
		return lib.CustomError(http.StatusBadRequest, "product_key can not be blank", "product_key can not be blank")
	}

	var product models.MsProduct
	status, err = models.GetMsProduct(&product, productKey)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(http.StatusBadRequest, err.Error(), "Product tidak ditemukan")
	}

	if product.FlagEnabled != 1 || product.FlagSubscription != 1 {
		log.Error("Tidak dapat melakukan autoinvest pada produk ini.")
		return lib.CustomError(http.StatusBadRequest, "Tidak dapat melakukan autoinvest pada produk ini.", "Tidak dapat melakukan autoinvest pada produk ini.")
	}

	customerKey := c.FormValue("customer_key")
	if customerKey != "" {
		n, err := strconv.ParseUint(customerKey, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: customer_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: customer_key", "Wrong input for parameter: customer_key")
		}

		if len(customerKey) > 11 {
			log.Error("Wrong input for parameter: customer_key too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: customer_key too long, max 11 character", "Missing required parameter: customer_key too long, max 11 character")
		}
		paramsAcc["customer_key"] = customerKey
	} else {
		log.Error("Missing required parameter: customer_key")
		return lib.CustomError(http.StatusBadRequest, "customer_key can not be blank", "customer_key can not be blank")
	}

	var customer models.MsCustomer
	status, err = models.GetMsCustomer(&customer, customerKey)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
	} else {
		if customer.CifSuspendFlag == 1 {
			log.Error("Account customer tersuspend. CIF Suspend")
			return lib.CustomError(http.StatusBadRequest, "Customer Account Suspended", "Customer Account Suspended")

		}
	}

	investAmount := c.FormValue("invest_amount")
	if investAmount != "" {
		value, err := decimal.NewFromString(investAmount)
		if err != nil {
			log.Error("Wrong input for parameter: invest_amount")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: invest_amount", "Wrong input for parameter: invest_amount")
		}
		if value.Cmp(product.MinSubAmount) == -1 {
			log.Error("invest_amount < minimum sub")
			return lib.CustomError(http.StatusBadRequest, "invest_amount < minum sub", "Minumum invest_amount untuk product ini adalah: "+product.MinSubAmount.String())
		}
		if investAmount == "0" {
			log.Error("Wrong input for parameter: invest_amount")
			return lib.CustomError(http.StatusBadRequest, "invest_amount harus lebih dari 0", "invest_amount harus lebih dari 0")
		}
		params["invest_amount"] = investAmount
	} else {
		log.Error("Missing required parameter: invest_amount")
		return lib.CustomError(http.StatusBadRequest, "invest_amount can not be blank", "invest_amount can not be blank")
	}

	investFeeRates := c.FormValue("invest_fee_rates")
	if investFeeRates != "" {
		_, err := decimal.NewFromString(investFeeRates)
		if err != nil {
			log.Error("Wrong input for parameter: invest_fee_rates")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: invest_fee_rates", "Wrong input for parameter: invest_fee_rates")
		}
		params["invest_fee_rates"] = investFeeRates
	} else {
		log.Error("Missing required parameter: invest_fee_rates")
		return lib.CustomError(http.StatusBadRequest, "invest_fee_rates can not be blank", "invest_fee_rates can not be blank")
	}

	investFeeAmount := c.FormValue("invest_fee_amount")
	if investFeeAmount != "" {
		_, err := decimal.NewFromString(investFeeAmount)
		if err != nil {
			log.Error("Wrong input for parameter: invest_fee_amount")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: invest_fee_amount", "Wrong input for parameter: invest_fee_amount")
		}
		params["invest_fee_amount"] = investFeeAmount
	} else {
		log.Error("Missing required parameter: invest_fee_amount")
		return lib.CustomError(http.StatusBadRequest, "invest_fee_amount can not be blank", "invest_fee_amount can not be blank")
	}

	investFeeCharges := c.FormValue("invest_fee_charges")
	if investFeeCharges != "" {
		_, err := decimal.NewFromString(investFeeCharges)
		if err != nil {
			log.Error("Wrong input for parameter: invest_fee_charges")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: invest_fee_charges", "Wrong input for parameter: invest_fee_charges")
		}
		params["invest_fee_charges"] = investFeeCharges
	} else {
		log.Error("Missing required parameter: invest_fee_charges")
		return lib.CustomError(http.StatusBadRequest, "invest_fee_charges can not be blank", "invest_fee_charges can not be blank")
	}

	investDateExecute := c.FormValue("invest_date_execute")
	if investDateExecute != "" {
		n, err := strconv.ParseUint(investDateExecute, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: invest_date_execute")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: invest_date_execute", "Wrong input for parameter: invest_date_execute")
		}
		if n > 28 {
			log.Error("invest_date_execute max tanggal 28")
			return lib.CustomError(http.StatusBadRequest, "invest_date_execute max tanggal 28", "invest_date_execute max tanggal 28")
		}
		params["invest_date_execute"] = investDateExecute
	} else {
		log.Error("Missing required parameter: invest_date_execute")
		return lib.CustomError(http.StatusBadRequest, "invest_date_execute can not be blank", "invest_date_execute can not be blank")
	}

	dateThru := c.FormValue("date_thru")
	if dateThru != "" {
		params["date_thru"] = dateThru
	} else {
		log.Error("Missing required parameter: date_thru")
		return lib.CustomError(http.StatusBadRequest, "date_thru can not be blank", "date_thru can not be blank")
	}

	settleChannel := c.FormValue("settle_channel")
	if settleChannel != "" {
		n, err := strconv.ParseUint(settleChannel, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: settle_channel")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: settle_channel", "Wrong input for parameter: settle_channel")
		}

		params["settle_channel"] = settleChannel
	} else {
		log.Error("Missing required parameter: settle_channel")
		return lib.CustomError(http.StatusBadRequest, "settle_channel can not be blank", "settle_channel can not be blank")
	}

	settlePaymentMethod := c.FormValue("settle_payment_method")
	if settlePaymentMethod != "" {
		n, err := strconv.ParseUint(settlePaymentMethod, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: settle_payment_method")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: settle_payment_method", "Wrong input for parameter: settle_payment_method")
		}

		params["settle_payment_method"] = settlePaymentMethod
	} else {
		log.Error("Missing required parameter: settle_payment_method")
		return lib.CustomError(http.StatusBadRequest, "settle_payment_method can not be blank", "settle_payment_method can not be blank")
	}

	prodBankAccKey := c.FormValue("prod_bank_acc_key")
	if prodBankAccKey != "" {
		n, err := strconv.ParseUint(prodBankAccKey, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: prod_bank_acc_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: prod_bank_acc_key", "Wrong input for parameter: prod_bank_acc_key")
		}

	} else {
		log.Error("Missing required parameter: prod_bank_acc_key")
		return lib.CustomError(http.StatusBadRequest, "prod_bank_acc_key can not be blank", "prod_bank_acc_key can not be blank")
	}

	investRemarks := c.FormValue("invest_remarks")
	if investRemarks != "" {
		if len(investRemarks) > 140 {
			log.Error("Wrong input for parameter: invest_remarks too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: invest_remarks too long, max 140 character", "Missing required parameter: invest_remarks too long, max 140 character")
		}
		params["invest_remarks"] = investRemarks
	}

	dateLayout := "2006-01-02 15:04:05"
	var accKey string
	paramsAcc["rec_status"] = "1"
	var trAccountDB []models.TrAccount
	status, err = models.GetAllTrAccount(&trAccountDB, paramsAcc)
	if len(trAccountDB) > 0 {
		accKey = strconv.FormatUint(trAccountDB[0].AccKey, 10)
		if trAccountDB[0].SubSuspendFlag != nil && *trAccountDB[0].SubSuspendFlag == 1 {
			log.Error("Account suspended for this product")
			return lib.CustomError(status, "Account suspended for this product", "Account suspended for this product")
		}
	} else {
		paramsAcc["acc_status"] = "204"
		paramsAcc["rec_created_date"] = time.Now().Format(dateLayout)
		paramsAcc["rec_created_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
		status, err, accKey = models.CreateTrAccount(paramsAcc)
		if err != nil {
			log.Error("Failed create account product data: " + err.Error())
			return lib.CustomError(status, err.Error(), "failed input data")
		}
	}

	var countData models.CountData
	status, err = models.AdminValidateAccAndiInvestDateExecute(&countData, accKey, investDateExecute, autoinvestKey)
	if err != nil {
		log.Error(err.Error())
	}
	if int(countData.CountData) > int(0) {
		log.Error("Ada data product dengan execute date date yang sama")
		return lib.CustomError(status, "Terdapat data product dengan execute date date yang sama", "Terdapat data product dengan execute date date yang sama")
	}

	params["acc_key"] = accKey
	params["rec_modified_date"] = time.Now().Format(dateLayout)
	params["rec_modified_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

	var productBankAccount models.MsProductBankAccount
	status, err = models.GetMsProductBankAccount(&productBankAccount, prodBankAccKey)
	if err != nil {
		log.Error("Failed get product bank account: " + err.Error())
	} else {
		params["bank_account_key"] = strconv.FormatUint(*productBankAccount.BankAccountKey, 10)
		var bankAccount models.MsBankAccount
		status, err = models.GetBankAccount(&bankAccount, strconv.FormatUint(*productBankAccount.BankAccountKey, 10))
		if err != nil {
			log.Error("Failed get bank account: " + err.Error())
		} else {
			params["bank_key"] = strconv.FormatUint(bankAccount.BankKey, 10)
		}
	}

	_, err = models.UpdateTrAutoinvestRegistration(params)
	if err != nil {
		log.Error("Failed update tr_autoinvest_registration: " + err.Error())
		return lib.CustomError(status, err.Error(), "failed update data")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = ""
	return c.JSON(http.StatusOK, response)
}

func AdminDeleteTrAutoInvest(c echo.Context) error {
	var err error
	var status int
	params := make(map[string]string)

	autoinvestKey := c.FormValue("autoinvest_key")
	if autoinvestKey != "" {
		n, err := strconv.ParseUint(autoinvestKey, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: autoinvest_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: autoinvest_key", "Wrong input for parameter: autoinvest_key")
		}

		if len(autoinvestKey) > 11 {
			log.Error("Wrong input for parameter: autoinvest_key too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: autoinvest_key too long, max 11 character", "Missing required parameter: autoinvest_key too long, max 11 character")
		}
		params["autoinvest_key"] = autoinvestKey
	} else {
		log.Error("Missing required parameter: autoinvest_key")
		return lib.CustomError(http.StatusBadRequest, "autoinvest_key can not be blank", "autoinvest_key can not be blank")
	}

	var invest models.TrAutoinvestRegistration
	_, err = models.GetTrAutoinvestRegistration(&invest, autoinvestKey)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(http.StatusBadRequest, "failed get data", "failed get data")
	}

	dateLayout := "2006-01-02 15:04:05"
	params["rec_status"] = "0"
	params["rec_deleted_date"] = time.Now().Format(dateLayout)
	params["rec_deleted_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

	_, err = models.UpdateTrAutoinvestRegistration(params)
	if err != nil {
		log.Error("Failed delete tr_autoinvest_registration: " + err.Error())
		return lib.CustomError(status, err.Error(), "failed delete data")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = nil
	return c.JSON(http.StatusOK, response)

}

func AdminDetailTrAutoInvest(c echo.Context) error {
	var err error

	autoinvestKey := c.Param("autoinvest_key")
	if autoinvestKey != "" {
		n, err := strconv.ParseUint(autoinvestKey, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: autoinvest_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: autoinvest_key", "Wrong input for parameter: autoinvest_key")
		}

		if len(autoinvestKey) > 11 {
			log.Error("Wrong input for parameter: autoinvest_key too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: autoinvest_key too long, max 11 character", "Missing required parameter: autoinvest_key too long, max 11 character")
		}
	} else {
		log.Error("Missing required parameter: autoinvest_key")
		return lib.CustomError(http.StatusBadRequest, "autoinvest_key can not be blank", "autoinvest_key can not be blank")
	}

	var invest models.DetailAutoinvestRegistration
	_, err = models.AdminGetDetailTrAutoinvestRegistration(&invest, autoinvestKey)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(http.StatusBadRequest, "failed get data", "failed get data")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = invest
	return c.JSON(http.StatusOK, response)

}

func AdminGenerateTransactionFromTrAutoInvest(c echo.Context) error {
	var err error
	var status int

	params := make(map[string]string)

	autoinvestKey := c.FormValue("autoinvest_key")
	if autoinvestKey != "" {
		n, err := strconv.ParseUint(autoinvestKey, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: autoinvest_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: autoinvest_key", "Wrong input for parameter: autoinvest_key")
		}

		if len(autoinvestKey) > 11 {
			log.Error("Wrong input for parameter: autoinvest_key too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: autoinvest_key too long, max 11 character", "Missing required parameter: autoinvest_key too long, max 11 character")
		}
	} else {
		log.Error("Missing required parameter: autoinvest_key")
		return lib.CustomError(http.StatusBadRequest, "autoinvest_key can not be blank", "autoinvest_key can not be blank")
	}

	var invest models.DetailAutoinvestRegistration
	_, err = models.AdminGetDetailTrAutoinvestRegistration(&invest, autoinvestKey)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(http.StatusBadRequest, "failed get data", "failed get data")
	}

	if invest.SubSuspendFlag != nil && *invest.SubSuspendFlag == uint8(1) {
		log.Error("Account Customer Product not allowed subscription")
		return lib.CustomError(http.StatusBadRequest, "Account Customer Product Not Allowed Subscription", "Account Customer Product Not Allowed Subscription")
	}

	productKey := strconv.FormatUint(invest.ProductKey, 10)

	var product models.MsProduct
	status, err = models.GetMsProduct(&product, productKey)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(http.StatusBadRequest, err.Error(), "Product tidak ditemukan")
	}

	if product.FlagSubscription != uint8(1) {
		log.Error("Product not allowed subscription")
		return lib.CustomError(http.StatusBadRequest, "Product Not Allowed Subscription", "Product Not Allowed Subscription")
	}

	customerKey := strconv.FormatUint(invest.CustomerKey, 10)
	investAmount := invest.InvestAmount
	settleChannel := strconv.FormatUint(invest.SettleChannel, 10)
	settlePaymentMethod := strconv.FormatUint(invest.SettlePaymentMethod, 10)
	prodBankAccKey := strconv.FormatUint(*invest.ProdBankaccKey, 10)
	accKey := invest.AccKey

	transRemarks := c.FormValue("trans_remarks")
	if transRemarks != "" {
		if len(transRemarks) > 140 {
			log.Error("Wrong input for parameter: trans_remarks too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: trans_remarks too long, max 140 character", "Missing required parameter: trans_remarks too long, max 140 character")
		}
	}

	nol := decimal.NewFromInt(0)
	feePercent := invest.InvestFeeRates
	feeAmount := invest.InvestFeeAmount
	chargesFeeAmount := invest.InvestFeeCharges

	var fee models.ProductFeeValueSubscription
	status, err = models.GetProductFeeValueSubscription(&fee, productKey)
	if err == nil {
		if fee.FeeNominalType == 192 { //percent
			feePercent = fee.FeeValue
			feeAmount = feePercent.Div(decimal.NewFromInt(100)).Mul(invest.InvestAmount)
		} else {
			feePercent = nol
			feeAmount = fee.FeeValue
		}
	}

	totalAmount := invest.InvestAmount.Add(feeAmount).Add(chargesFeeAmount).Truncate(0)

	dateLayout := "2006-01-02 15:04:05"

	agentKey := "1"
	var agent models.CustomerAgent
	status, err = models.GetCustomerLastAgent(&agent, customerKey)
	if err == nil {
		agentKey = strconv.FormatUint(agent.AgentKey, 10)
	}

	branchKey := "1"
	var agentBranch models.MsAgentBranch
	status, err = models.GetLastBranchAgent(&agentBranch, agentKey)
	if err == nil {
		branchKey = strconv.FormatUint(agentBranch.BranchKey, 10)
	}

	//cek tr_account_agent / save
	paramsAccAgent := make(map[string]string)
	paramsAccAgent["acc_key"] = strconv.FormatUint(accKey, 10)
	paramsAccAgent["agent_key"] = agentKey
	paramsAccAgent["rec_status"] = "1"

	var acaKey string
	var accountAgentDB []models.TrAccountAgent
	status, err = models.GetAllTrAccountAgent(&accountAgentDB, paramsAccAgent)
	if len(accountAgentDB) > 0 {
		acaKey = strconv.FormatUint(accountAgentDB[0].AcaKey, 10)
	} else {
		paramsCreateAccAgent := make(map[string]string)
		paramsCreateAccAgent["acc_key"] = strconv.FormatUint(accKey, 10)
		paramsCreateAccAgent["eff_date"] = time.Now().Format(dateLayout)
		paramsCreateAccAgent["rec_created_date"] = time.Now().Format(dateLayout)
		paramsCreateAccAgent["rec_created_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
		paramsCreateAccAgent["agent_key"] = agentKey
		paramsCreateAccAgent["branch_key"] = branchKey
		paramsCreateAccAgent["rec_status"] = "1"
		status, err, acaKey = models.CreateTrAccountAgent(paramsCreateAccAgent)
		if err != nil {
			log.Error("Failed create account agent data: " + err.Error())
			return lib.CustomError(status, err.Error(), "failed input data")
		}
	}

	paramsTrans := make(map[string]string)
	paramsTrans["branch_key"] = branchKey
	paramsTrans["agent_key"] = agentKey
	paramsTrans["customer_key"] = customerKey
	paramsTrans["product_key"] = productKey
	paramsTrans["trans_status_key"] = "2"
	paramsTrans["trans_date"] = time.Now().Format(dateLayout)
	paramsTrans["trans_type_key"] = "1"
	paramsTrans["trx_code"] = "137"

	layoutISO := "2006-01-02"
	navDate := time.Now().Format(layoutISO)

	paramHoliday := make(map[string]string)
	paramHoliday["holiday_date"] = navDate

	var holiday []models.MsHoliday
	status, err = models.GetAllMsHoliday(&holiday, paramHoliday)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
	}

	if len(holiday) > 0 {
		var tglBursa models.TanggalBursa
		status, err = models.GetTanggalBursa(&tglBursa, navDate, "1")
		if err == nil {
			navDate = tglBursa.TanggalBursa
		}
	} else {
		t, _ := time.Parse(layoutISO, navDate)
		strDate := t.Format(layoutISO)
		w, _ := time.Parse(layoutISO, strDate)
		w = time.Date(w.Year(), w.Month(), w.Day(), 0, 0, 0, 0, time.UTC)
		cek := lib.IsWeekend(w)
		if cek {
			var tglBursa models.TanggalBursa
			status, err = models.GetTanggalBursa(&tglBursa, navDate, "1")
			if err == nil {
				navDate = tglBursa.TanggalBursa
			}
		} else {
			now := time.Now()
			if (now.Hour() >= 12 && now.Minute() > 0) || now.Hour() > 12 {
				var tglBursa models.TanggalBursa
				status, err = models.GetTanggalBursa(&tglBursa, navDate, "1")
				if err == nil {
					navDate = tglBursa.TanggalBursa
				}
			}
		}
	}

	paramsTrans["nav_date"] = navDate
	paramsTrans["entry_mode"] = "140"
	paramsTrans["aca_key"] = acaKey
	paramsTrans["trans_calc_method"] = "289"
	paramsTrans["trans_amount"] = investAmount.Truncate(2).String()
	paramsTrans["trans_unit"] = "0"

	paramsTr := make(map[string]string)
	paramsTr["customer_key"] = customerKey
	paramsTr["product_key"] = productKey
	paramsTr["trans_type_key"] = "1"
	var transactionDB []models.TrTransaction
	status, err = models.GetAllTrTransaction(&transactionDB, paramsTr)
	if err != nil {
		paramsTrans["flag_newsub"] = "1"
	} else {
		paramsTrans["flag_newsub"] = "0"
	}

	paramsTrans["trans_fee_percent"] = feePercent.String()
	paramsTrans["trans_fee_amount"] = feeAmount.String()
	paramsTrans["charges_fee_amount"] = chargesFeeAmount.String()
	paramsTrans["services_fee_amount"] = "0"
	paramsTrans["total_amount"] = totalAmount.String()
	paramsTrans["trans_remarks"] = transRemarks
	risk := "0"

	var riskProfil models.RiskProfilCustomer
	status, err = models.GetRiskProfilCustomer(&riskProfil, customerKey)
	if err != nil {
		if product.RiskProfileKey != nil {
			if riskProfil.RiskProfileKey < *product.RiskProfileKey {
				risk = "1"
			}
		}
	}
	paramsTrans["risk_waiver"] = risk

	paramsTrans["trans_source"] = "142"
	paramsTrans["payment_method"] = settlePaymentMethod
	paramsTrans["rec_order"] = "0"
	paramsTrans["rec_status"] = "1"
	paramsTrans["rec_created_date"] = time.Now().Format(dateLayout)
	paramsTrans["rec_created_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
	paramsTrans["rec_modified_date"] = time.Now().Format(dateLayout)
	paramsTrans["rec_modified_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

	//create tr_transaction
	status, err, transactionID := models.CreateTrTransaction(paramsTrans)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed input data")
	}

	//create tr_transaction_bank_account
	bankAccKeySource := "1"
	paramsTransactionBankAccount := make(map[string]string)
	paramsTransactionBankAccount["transaction_key"] = transactionID
	paramsTransactionBankAccount["prod_bankacc_key"] = prodBankAccKey
	var customerBankDB []models.MsCustomerBankAccount
	paramCustomerBank := make(map[string]string)
	paramCustomerBank["customer_key"] = customerKey
	paramCustomerBank["rec_status"] = "1"
	paramCustomerBank["orderBy"] = "flag_priority"
	paramCustomerBank["orderType"] = "DESC"
	status, err = models.GetAllMsCustomerBankAccount(&customerBankDB, paramCustomerBank)
	if err != nil {
		log.Error(err.Error())
		paramsTransactionBankAccount["cust_bankacc_key"] = "1"
	} else {
		bankAccKeySource = strconv.FormatUint(customerBankDB[0].BankAccountKey, 10)
		paramsTransactionBankAccount["cust_bankacc_key"] = strconv.FormatUint(customerBankDB[0].CustBankaccKey, 10)
	}
	paramsTransactionBankAccount["rec_status"] = "1"
	paramsTransactionBankAccount["rec_created_date"] = time.Now().Format(dateLayout)
	paramsTransactionBankAccount["rec_created_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
	status, err = models.CreateTrTransactionBankAccount(paramsTransactionBankAccount)
	if err != nil {
		log.Error(err.Error())
	}
	//create tr_transaction_settlement
	settlementParams := make(map[string]string)
	settlementParams["transaction_key"] = transactionID
	settlementParams["settle_purposed"] = "297"
	settlementParams["settle_date"] = navDate
	settlementParams["settle_nominal"] = totalAmount.String()
	settlementParams["client_subaccount_no"] = ""
	settlementParams["settled_status"] = "243"
	settlementParams["source_bank_account_key"] = bankAccKeySource
	settlementParams["target_bank_account_key"] = prodBankAccKey
	settlementParams["settle_channel"] = settleChannel
	settlementParams["settle_payment_method"] = "10" //Manual
	settlementParams["rec_status"] = "1"

	var tglBursaExpired models.TanggalBursa
	status, err = models.GetTanggalBursa(&tglBursaExpired, time.Now().Format(dateLayout), "1")
	if err == nil {
		settlementParams["expired_date"] = tglBursaExpired.TanggalBursa
	}

	_, err, _ = models.CreateTrTransactionSettlement(settlementParams)
	if err != nil {
		log.Error(err.Error())
	}

	// update tr_autoinvest_registration
	attempt := invest.AttempCount + 1
	params["autoinvest_key"] = autoinvestKey
	params["attempt_count"] = strconv.FormatUint(attempt, 10)
	params["date_last_generate"] = time.Now().Format(dateLayout)
	params["rec_modified_date"] = time.Now().Format(dateLayout)
	params["rec_modified_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
	_, err = models.UpdateTrAutoinvestRegistration(params)
	if err != nil {
		log.Error("Failed update tr_autoinvest_registration: " + err.Error())
		return lib.CustomError(status, err.Error(), "failed update data")
	}

	//create tr_transaction_settlement
	paramsAutoExe := make(map[string]string)
	paramsAutoExe["autoinvest_key"] = autoinvestKey
	paramsAutoExe["transaction_key"] = transactionID
	paramsAutoExe["execute_date"] = time.Now().Format(dateLayout)
	paramsAutoExe["invest_remarks"] = transRemarks
	paramsAutoExe["rec_created_date"] = time.Now().Format(dateLayout)
	paramsAutoExe["rec_created_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
	paramsAutoExe["rec_status"] = "1"

	_, err = models.CreateTrAutoinvestExecution(paramsAutoExe)
	if err != nil {
		log.Error(err.Error())
	}

	//create message
	var body string
	subject := "Segera Upload Bukti Transfer Kamu"
	if params["flag_newsub"] == "1" {
		body = "Segera Upload Bukti Transfer Kamu agar Subscribe kamu dapat segera kami proses."
	} else {
		body = "Segera Upload Bukti Transfer Kamu agar Top Up kamu dapat segera kami proses."
	}

	var userData models.ScUserLogin
	status, err = models.GetScUserLoginByCustomerKey(&userData, customerKey)
	if err != nil {
		return lib.CustomError(status)
	}

	customerUserLoginKey := strconv.FormatUint(userData.UserLoginKey, 10)
	paramsUserMessage := make(map[string]string)
	paramsUserMessage["umessage_type"] = "245"
	paramsUserMessage["umessage_recipient_key"] = customerUserLoginKey
	paramsUserMessage["umessage_receipt_date"] = time.Now().Format(dateLayout)
	paramsUserMessage["flag_read"] = "0"
	paramsUserMessage["umessage_sent_date"] = time.Now().Format(dateLayout)
	paramsUserMessage["flag_sent"] = "1"

	paramsUserMessage["umessage_subject"] = subject
	paramsUserMessage["umessage_body"] = body

	paramsUserMessage["umessage_category"] = "248"
	paramsUserMessage["flag_archieved"] = "0"
	paramsUserMessage["archieved_date"] = time.Now().Format(dateLayout)
	paramsUserMessage["rec_status"] = "1"
	paramsUserMessage["rec_created_date"] = time.Now().Format(dateLayout)
	paramsUserMessage["rec_created_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

	status, err = models.CreateScUserMessage(paramsUserMessage)
	if err != nil {
		log.Error("Error create user message")
	} else {
		log.Error("Sukses insert user message")
	}
	lib.CreateNotifCustomerFromAdminByCustomerId(customerKey, subject, body, "TRANSACTION")

	paramsEmail := make(map[string]string)
	paramsEmail["currency"] = strconv.FormatUint(*product.CurrencyKey, 10)
	paramsEmail["trans_fee_amount"] = paramsTrans["trans_fee_amount"]
	paramsEmail["trans_amount"] = paramsTrans["trans_amount"]
	paramsEmail["trans_date"] = time.Now().Format(dateLayout)
	paramsEmail["product_name"] = invest.ProductName
	paramsEmail["customer_key"] = customerKey
	paramsEmail["email_customer"] = userData.UloginEmail
	err = mailSubscriptionAutoInvest(paramsEmail)

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = ""
	return c.JSON(http.StatusOK, response)
}

func mailSubscriptionAutoInvest(params map[string]string) error {
	var err error
	var mailTemp, subject string
	decimal.MarshalJSONWithoutQuotes = true
	ac0 := accounting.Accounting{Symbol: "", Precision: 0, Thousand: ".", Decimal: ","}
	mailParam := make(map[string]string)
	if params["currency"] == "1" {
		mailParam["Symbol"] = "Rp. "
	} else if params["currency"] == "2" {
		mailParam["Symbol"] = "$"
	}
	val, _ := decimal.NewFromString(params["trans_fee_amount"])
	mailParam["Fee"] = ac0.FormatMoneyDecimal(val.Truncate(0))
	mailTemp = "index-subscription-uncomplete.html"
	subject = "Ayo Upload Bukti Transfer Kamu"
	val, _ = decimal.NewFromString(params["trans_amount"])
	mailParam["Amount"] = ac0.FormatMoneyDecimal(val)
	var customer models.MsCustomer
	_, err = models.GetMsCustomer(&customer, params["customer_key"])
	if err != nil {
		log.Error("Failed send mail: " + err.Error())
		return err
	}
	mailParam["Name"] = customer.FullName
	mailParam["Cif"] = customer.UnitHolderIDno
	layout := "2006-01-02 15:04:05"
	dateLayout := "02 Jan 2006"
	date, _ := time.Parse(layout, params["trans_date"])
	mailParam["Date"] = date.Format(dateLayout)
	hr, min, _ := date.Clock()
	mailParam["Time"] = strconv.Itoa(hr) + "." + strconv.Itoa(min) + " WIB"

	mailParam["ProductName"] = params["product_name"]
	mailParam["ProductIn"] = params["product_name"]

	mailParam["FileUrl"] = config.FileUrl + "/images/mail"

	t := template.New(mailTemp)

	t, err = t.ParseFiles(config.BasePath + "/mail/" + mailTemp)
	if err != nil {
		log.Error("Failed send mail: " + err.Error())
		return err
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, mailParam); err != nil {
		log.Error("Failed send mail: " + err.Error())
		return err
	}

	result := tpl.String()

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", config.EmailFrom)
	mailer.SetHeader("To", params["email_customer"])
	mailer.SetHeader("Subject", "[MotionFunds] "+subject)
	mailer.SetBody("text/html", result)

	err = lib.SendEmail(mailer)
	if err != nil {
		log.Error("Failed send mail: " + err.Error())
		return err
	}

	// dialer := gomail.NewDialer(
	// 	config.EmailSMTPHost,
	// 	int(config.EmailSMTPPort),
	// 	config.EmailFrom,
	// 	config.EmailFromPassword,
	// )
	// dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// err = dialer.DialAndSend(mailer)
	// if err != nil {
	// 	log.Error("Failed send mail: " + err.Error())
	// 	return err
	// }
	log.Info("Email sent")
	return nil
}
