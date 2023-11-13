package controllers

import (
	"mf-bo-api/lib"
	"mf-bo-api/models"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo"
)

func CreateProductBankRequest(c echo.Context) error {
	paramsProductBankAccount := make(map[string]string)
	paramsBankAcc := make(map[string]string)

	//product_key
	productkey := c.FormValue("product_key")
	if productkey == "" {
		// log.Error("Missing required parameter: product_key cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: product_key cann't be blank", "Missing required parameter: product_key cann't be blank")
	}
	strproductkey, err := strconv.ParseUint(productkey, 10, 64)
	if err == nil && strproductkey > 0 {
		paramsProductBankAccount["product_key"] = productkey
	} else {
		// log.Error("Wrong input for parameter: product_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: product_key", "Missing required parameter: product_key")
	}

	//bank_key
	bankkey := c.FormValue("bank_key")
	if bankkey == "" {
		// log.Error("Missing required parameter: bank_key cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: bank_key cann't be blank", "Missing required parameter: bank_key cann't be blank")
	}
	strbankkey, err := strconv.ParseUint(bankkey, 10, 64)
	if err == nil && strbankkey > 0 {
		paramsBankAcc["bank_key"] = bankkey
	} else {
		// log.Error("Wrong input for parameter: bank_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: bank_key", "Missing required parameter: bank_key")
	}

	//account_no
	accountno := c.FormValue("account_no")
	if accountno == "" {
		// log.Error("Missing required parameter: account_no cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: account_no cann't be blank", "Missing required parameter: account_no cann't be blank")
	}
	paramsBankAcc["account_no"] = accountno

	//account_holder_name
	accountholdername := c.FormValue("account_holder_name")
	if accountholdername == "" {
		// log.Error("Missing required parameter: account_holder_name cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: account_holder_name cann't be blank", "Missing required parameter: account_holder_name cann't be blank")
	}
	paramsBankAcc["account_holder_name"] = accountholdername

	//branch_name
	branchname := c.FormValue("branch_name")
	if branchname != "" {
		paramsBankAcc["branch_name"] = branchname
	}

	//currency_key
	currencykey := c.FormValue("currency_key")
	if currencykey == "" {
		// log.Error("Missing required parameter: currency_key cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: currency_key cann't be blank", "Missing required parameter: currency_key cann't be blank")
	}
	strcurrencykey, err := strconv.ParseUint(currencykey, 10, 64)
	if err == nil && strcurrencykey > 0 {
		paramsBankAcc["currency_key"] = currencykey
	} else {
		// log.Error("Wrong input for parameter: currency_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: currency_key", "Missing required parameter: currency_key")
	}

	//bank_account_type
	bankaccounttype := c.FormValue("bank_account_type")
	if bankaccounttype == "" {
		// log.Error("Missing required parameter: bank_account_type cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: bank_account_type cann't be blank", "Missing required parameter: bank_account_type cann't be blank")
	}
	strbankaccounttype, err := strconv.ParseUint(bankaccounttype, 10, 64)
	if err == nil && strbankaccounttype > 0 {
		paramsBankAcc["bank_account_type"] = bankaccounttype
	} else {
		// log.Error("Wrong input for parameter: bank_account_type")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: bank_account_type", "Missing required parameter: bank_account_type")
	}

	paramsBankAcc["rec_domain"] = "132"

	//swift_code
	swiftcode := c.FormValue("swift_code")
	if swiftcode != "" {
		paramsBankAcc["swift_code"] = swiftcode
	}

	//bank_account_name
	bankaccountname := c.FormValue("bank_account_name")
	if bankaccountname == "" {
		// log.Error("Missing required parameter: bank_account_name cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: bank_account_name cann't be blank", "Missing required parameter: bank_account_name cann't be blank")
	}
	paramsProductBankAccount["bank_account_name"] = bankaccountname

	//bank_account_type
	bankaccountpurpose := c.FormValue("bank_account_purpose")
	if bankaccountpurpose == "" {
		// log.Error("Missing required parameter: bank_account_purpose cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: bank_account_purpose cann't be blank", "Missing required parameter: bank_account_purpose cann't be blank")
	}
	strbankaccountpurpose, err := strconv.ParseUint(bankaccountpurpose, 10, 64)
	if err == nil && strbankaccountpurpose > 0 {
		paramsProductBankAccount["bank_account_purpose"] = bankaccountpurpose
	} else {
		// log.Error("Wrong input for parameter: bank_account_purpose")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: bank_account_purpose", "Missing required parameter: bank_account_purpose")
	}

	paramsBankAcc["rec_status"] = "1"
	paramsBankAcc["rec_created_date"] = time.Now().Format(lib.TIMESTAMPFORMAT)
	paramsBankAcc["rec_created_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

	paramsProductBankAccount["rec_status"] = "1"
	paramsProductBankAccount["rec_action"] = "CREATE"
	paramsProductBankAccount["rec_created_date"] = time.Now().Format(lib.TIMESTAMPFORMAT)
	paramsProductBankAccount["rec_created_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

	err = models.CreateRequestProductBankAccount(paramsBankAcc, paramsProductBankAccount)
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
