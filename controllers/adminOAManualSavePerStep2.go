package controllers

import (
	"encoding/json"
	"fmt"
	"mf-bo-api/config"
	"mf-bo-api/lib"
	"mf-bo-api/models"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo"
)

func SaveStep4(c echo.Context) (error, int64) {
	var OaRequestKey int64

	paramsOaRequest := make(map[string]string)
	paramsOaRequest["oa_step"] = "4"

	oa_request_key := c.FormValue("oa_request_key")
	if oa_request_key == "" {
		return fmt.Errorf("Missing: oa_request_key"), OaRequestKey
	}
	paramsOaRequest["oa_request_key"] = oa_request_key

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

			err, OaRequestKey = models.CreateOaBankAccount(paramsOaRequest, params)
			if err != nil {
				return err, OaRequestKey
			}
		}
	}

	return nil, OaRequestKey
}

func SaveStep5(c echo.Context) (error, int64) {
	var OaRequestKey int64
	paramsOaRequest := make(map[string]string)
	paramsOaRequest["oa_step"] = "5"
	paramsOaRequest["rec_status"] = "1"
	paramsOaRequest["rec_created_by"] = lib.UserIDStr
	paramsOaRequest["rec_created_date"] = time.Now().Format(lib.TIMESTAMPFORMAT)

	oa_request_key := c.FormValue("oa_request_key")
	if oa_request_key == "" {
		return fmt.Errorf("Missing: oa_request_key"), OaRequestKey
	}
	paramsOaRequest["oa_request_key"] = oa_request_key

	quiz_answers := c.FormValue("quiz_answers")
	if quiz_answers == "" {
		return fmt.Errorf("Missing: quiz_answers"), OaRequestKey
	}

	var quiz_array []models.OaQuizAnswer
	err := json.Unmarshal([]byte(quiz_answers), &quiz_array)
	if err != nil {
		return err, OaRequestKey
	}

	if len(quiz_array) == 0 {
		return fmt.Errorf("Missing: quiz"), OaRequestKey
	}

	err, OaRequestKey = models.CreateOrUpdateOaRiskProfileQuiz(paramsOaRequest, quiz_array)
	if err != nil {
		return err, OaRequestKey
	}

	return nil, OaRequestKey
}

func SaveStep6(c echo.Context) (error, int64) {
	var OaRequestKey int64

	paramsOaRequest := make(map[string]string)
	paramsOaRequest["oa_step"] = "6"

	oa_request_key := c.FormValue("oa_request_key")
	if oa_request_key == "" {
		return fmt.Errorf("Missing: oa_request_key"), OaRequestKey
	}
	paramsOaRequest["oa_request_key"] = oa_request_key

	getParamsData := models.GetOptionByLookupGroupKey("105")
	if len(getParamsData) > 0 {
		for _, data := range getParamsData {
			file_upload, err := c.FormFile("file_upload_" + strconv.FormatUint(data.Key, 10))
			if err != nil {
				return err, OaRequestKey
			}

			if file_upload != nil {
				err = os.MkdirAll(config.BasePathImage+"/images/oa_manual/"+oa_request_key, 0755)
				if err != nil {
					return err, OaRequestKey
				}

				extension := filepath.Ext(file_upload.Filename)
				filename := strings.ReplaceAll(data.Value, " ", "_") + "_" + lib.RandStringBytesMaskImprSrc(26)
				targetDir := config.BasePathImage + "/images/oa_manual/" + oa_request_key + "/" + filename + extension
				err = lib.UploadImage(file_upload, targetDir)
				if err != nil {
					return err, OaRequestKey
				}
				// paramsOaPersonalData["pic_selfie_ktp"] = filename + extension

				createFile := make(map[string]string)
				createFile["ref_fk_key"] = oa_request_key
				createFile["ref_fk_domain"] = "oa_request"
				createFile["rec_status"] = "1"
				createFile["rec_created_date"] = time.Now().Format(lib.TIMESTAMPFORMAT)
				createFile["rec_created_by"] = lib.UserIDStr
				createFile["file_name"] = filename
				createFile["file_ext"] = extension
				createFile["file_path"] = "/images/oa_manual/" + oa_request_key + "/" + filename + extension
				createFile["file_url"] = config.ImageUrl + "/images/oa_manual/" + oa_request_key + "/" + filename + extension

				err, _ = models.CreateOrUpdateFileOaManual(paramsOaRequest, createFile)
				if err != nil {
					return err, OaRequestKey
				}
			}
		}
	}

	OaRequestKey, _ = strconv.ParseInt(paramsOaRequest["oa_request_key"], 10, 64)
	return nil, OaRequestKey
}
