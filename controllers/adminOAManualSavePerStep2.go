package controllers

import (
	"encoding/json"
	"fmt"
	"mf-bo-api/lib"
	"mf-bo-api/models"
	"strconv"
	"time"

	"github.com/labstack/echo"
)

func SaveStep4(c echo.Context) (error, int64) {
	var OaRequestKey int64

	oa_request_key := c.FormValue("oa_request_key")
	if oa_request_key == "" {
		return fmt.Errorf("Missing: bank_accounts"), OaRequestKey

	}

	bank_accounts := c.FormValue("bank_accounts")
	if bank_accounts == "" {
		return fmt.Errorf("Missing: bank_accounts"), OaRequestKey
	}

	var BankAccountsModel []models.OaRequestBankAccountDetails
	err := json.Unmarshal([]byte(bank_accounts), &BankAccountsModel)
	if err != nil {
		return err, OaRequestKey
	}

	if len(BankAccountsModel) == 0 {
		return fmt.Errorf("Missing: bank_accounts"), OaRequestKey
	} else {
		params := make(map[string]string)
		params["oa_request_key"] = oa_request_key
		params["rec_status"] = "1"
		params["rec_created_by"] = lib.UserIDStr
		params["rec_created_date"] = time.Now().Format(lib.TIMESTAMPFORMAT)

		for _, data := range BankAccountsModel {
			params["account_holder_name"] = *data.BankAccountName
			params["account_no"] = *data.BankAccountNo
			// params["bank_account_key"] = *data.BankAccountKey
			params["bank_key"] = strconv.FormatUint(*data.BankKey, 10)
			params["branch_name"] = *data.BankBranchName
			params["currency_key"] = strconv.FormatUint(*data.CurrencyKey, 10)
			params["flag_priority"] = strconv.FormatUint(*data.FlagPriority, 10)

			err, OaRequestKey = models.CreateOaBankAccount(params)
			if err != nil {
				return err, OaRequestKey
			}
		}
	}

	return nil, OaRequestKey
}
