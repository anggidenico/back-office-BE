package controllers

import (
	"database/sql"
	"errors"
	"log"
	"mf-bo-api/config"
	"mf-bo-api/lib"
	"mf-bo-api/models"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/labstack/echo"
)

func GetInvestPurposeController(c echo.Context) error {
	var invest []models.InvestPurpose
	status, err := models.GetInvestPurposeModels(&invest)
	if err != nil {
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = invest
	// log.Printf("Response Data: %+v\n", response.Data)
	return c.JSON(http.StatusOK, response)
}

func GetInvestPartnerController(c echo.Context) error {
	var invest []models.InvestPartner
	status, err := models.GetInvestPartnerModels(&invest)
	if err != nil {
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = invest
	// log.Printf("Response Data: %+v\n", response.Data)
	return c.JSON(http.StatusOK, response)
}

func GetInvestPartnerDetailController(c echo.Context) error {
	partnerKey := c.Param("invest_partner_key")
	if partnerKey == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing invest_partner_key", "Missing invest_partner_key")
	}
	var invest models.InvestPartner
	status, err := models.GetInvestPartnerDetailModels(&invest, partnerKey)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return lib.CustomError(http.StatusNotFound, "invest_partner_key not found", "invest_partner_key not found")
		}
		return lib.CustomError(status, err.Error(), err.Error())
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = invest
	return c.JSON(http.StatusOK, response)
}
func DeleteInvestPartnerController(c echo.Context) error {
	params := make(map[string]string)
	dateLayout := "2006-01-02 15:04:05"
	params["rec_status"] = "0"
	params["rec_deleted_date"] = time.Now().Format(dateLayout)
	params["rec_deleted_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

	investPartnerKey := c.FormValue("invest_partner_key")
	if investPartnerKey == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing invest_partner_key", "Missing invest_partner_key")
	}

	status, err := models.DeleteInvestPartnerModels(investPartnerKey, params)
	if err != nil {
		return lib.CustomError(status, err.Error(), err.Error())
	}
	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "Berhasil hapus CMS Invest Partner!"
	response.Data = ""
	return c.JSON(http.StatusOK, response)
}

func CreateInvestPartnerController(c echo.Context) error {
	var err error
	params := make(map[string]string)
	params["rec_created_by"] = lib.UserIDStr
	params["rec_created_date"] = time.Now().Format(lib.TIMESTAMPFORMAT)

	investPurposeKey := c.FormValue("invest_purpose_key")

	partnerCode := c.FormValue("partner_code")
	if partnerCode != "" {
		if len(partnerCode) > 30 {
			return lib.CustomError(http.StatusBadRequest, "partner_code must be <= 30 characters", "partner_code must be <= 30 characters")
		}
	} else {
		return lib.CustomError(http.StatusBadRequest, "partner_code cant be blank", "partner_code cant be blank")
	}
	partnerDesc := c.FormValue("partner_desc")
	if partnerDesc != "" {
		if len(partnerDesc) > 150 {
			return lib.CustomError(http.StatusBadRequest, "partner_desc must be <= 150 characters", "partner_desc must be <= 150 characters")
		}
	}
	params["partner_desc"] = partnerDesc
	partnerBusinessName := c.FormValue("partner_business_name")
	if partnerBusinessName != "" {
		if len(partnerBusinessName) > 50 {
			return lib.CustomError(http.StatusBadRequest, "partner_business_name must be <= 50 characters", "partner_business_name must be <= 50 characters")
		}
	}
	partnerPicName := c.FormValue("partner_picname")
	if partnerPicName != "" {
		if len(partnerPicName) > 150 {
			return lib.CustomError(http.StatusBadRequest, "partner_picname must be <= 150 characters", "partner_picname must be <= 150 characters")
		}
	}
	partnerMobileNo := c.FormValue("partner_mobileno")
	partnerOfficeNo := c.FormValue("partner_officeno")
	partnerEmail := c.FormValue("partner_email")
	if partnerEmail != "" {
		if len(partnerEmail) > 100 {
			return lib.CustomError(http.StatusBadRequest, "partner_email must be <= 100 characters", "partner_email must be <= 100 characters")
		}
	}
	partnerCity := c.FormValue("partner_city")
	if partnerCity != "" {
		if len(partnerCity) > 100 {
			return lib.CustomError(http.StatusBadRequest, "partner_city must be <= 100 characters", "partner_city must be <= 100 characters")
		}
	}
	partnerAddress := c.FormValue("partner_address")
	if partnerAddress != "" {
		if len(partnerAddress) > 200 {
			return lib.CustomError(http.StatusBadRequest, "partner_address must be <= 200 characters", "partner_address must be <= 200 characters")
		}
	}
	partnerUrl := c.FormValue("partner_url")
	if partnerUrl != "" {
		if len(partnerUrl) > 255 {
			return lib.CustomError(http.StatusBadRequest, "partner_url must be <= 200 characters", "partner_url must be <= 200 characters")
		}
	}
	partnerDateStarted := c.FormValue("partner_date_started")
	if partnerDateStarted != "" {
		expectedDateFormat := "2006-01-02"
		_, err = time.Parse(expectedDateFormat, partnerDateStarted)
		if err != nil {
			return lib.CustomError(http.StatusBadRequest, "partner_date_started should be a valid date in the format "+expectedDateFormat, "partner_date_started should be a valid date in the format "+expectedDateFormat)
		}
	} else {
		params["partner_date_started"] = "0000-00-00"
	}
	partnerDateExpired := c.FormValue("partner_date_expired")
	if partnerDateExpired != "" {
		expectedDateFormat := "2006-01-02"
		_, err = time.Parse(expectedDateFormat, partnerDateExpired)
		if err != nil {
			return lib.CustomError(http.StatusBadRequest, "partner_date_expired should be a valid date in the format "+expectedDateFormat, "partner_date_expired should be a valid date in the format "+expectedDateFormat)
		}
	} else {
		params["partner_date_expired"] = "0000-00-00"
	}
	partnerBannerHits := c.FormValue("partner_banner_hits")
	if partnerBannerHits != "" {
		if len(partnerBannerHits) > 11 {
			return lib.CustomError(http.StatusBadRequest, "partner_banner_hits must be <= 11 characters", "partner_banner_hits must be <= 11 characters")
		}
		value, err := strconv.Atoi(partnerBannerHits)
		if err != nil {
			return lib.CustomError(http.StatusBadRequest, "partner_banner_hits should be number", "partner_banner_hits should be number")
		}
		params["partner_banner_hits"] = strconv.Itoa(value)
	} else {
		params["partner_banner_hits"] = "0"
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
	var fileUpload *multipart.FileHeader
	_, err = c.FormFile("rec_image1")
	if fileUpload != nil {
		err = os.MkdirAll(config.BasePathImage+"/images/user/"+strconv.FormatUint(lib.Profile.UserID, 10), 0755)
		if err != nil {
			// log.Println(err.Error())
		} else {
			if fileUpload.Size > int64(lib.MAX_FILE_SIZE) {
				msgs := "Maximum size is " + lib.MAX_FILE_SIZE_MB + " MB"
				return lib.CustomError(http.StatusBadRequest, msgs, msgs)
			}
			// Get file extension
			extension := filepath.Ext(fileUpload.Filename)
			// Generate filename
			filename := "upload_image" + strconv.FormatUint(lib.Profile.UserID, 10) + "_" + lib.RandStringBytesMaskImprSrc(26)
			// log.Println("Generate filename:", filename)
			targetDir := config.BasePathImage + "/images/user/" + strconv.FormatUint(lib.Profile.UserID, 10) + "/" + filename + extension
			// Upload image and move to proper directory
			err = lib.UploadImage(fileUpload, targetDir)
			if err != nil {
				// log.Println(err)
				return lib.CustomError(http.StatusInternalServerError)
			}
			params["rec_image1"] = filename + extension
		}
	} else {
		return lib.CustomError(http.StatusBadRequest, "Foto KTP tidak boleh kosong", "Foto KTP tidak boleh kosong")
	}
	params["invest_purpose_key"] = investPurposeKey
	params["partner_code"] = partnerCode
	params["partner_business_name"] = partnerBusinessName
	params["partner_picname"] = partnerPicName
	params["partner_mobileno"] = partnerMobileNo
	params["partner_officeno"] = partnerOfficeNo
	params["partner_email"] = partnerEmail
	params["partner_city"] = partnerCity
	params["partner_address"] = partnerAddress
	params["partner_url"] = partnerUrl
	params["rec_status"] = "1"

	// Check for duplicate records
	duplicate, key, err := models.CheckDuplicateInvestPartner(params["partner_code"])
	if err != nil {
		log.Println("Error checking for duplicates:", err)
		return lib.CustomError(http.StatusInternalServerError, "Error checking for duplicates", "Error checking for duplicates")
	}
	log.Println("Duplicate:", duplicate)
	log.Println("Key:", key)
	// Jika duplikasi ditemukan, perbarui data yang sudah ada
	if duplicate {
		status, err := models.UpdateInvestPartner(key, params)
		if err != nil {
			log.Println("Failed to update data:", err)
			return lib.CustomError(status, "Failed to update data", "Failed to update data")
		}
		return c.JSON(http.StatusOK, lib.Response{
			Status: lib.Status{
				Code:          http.StatusOK,
				MessageServer: "OK",
				MessageClient: "OK",
			},
			Data: "Data updated successfully",
		})
	}

	// Jika tidak ada duplikasi, buat data baru
	status, err = models.CreateInvestPartner(params)
	log.Println("Error create data:", err)
	if err != nil {
		return lib.CustomError(status, "Failed input data", "Failed input data")
	}

	return c.JSON(http.StatusOK, lib.Response{
		Status: lib.Status{
			Code:          http.StatusOK,
			MessageServer: "OK",
			MessageClient: "OK",
		},
		Data: "Data created successfully",
	})
}
func UpdateInvestPartnerController(c echo.Context) error {
	var err error
	params := make(map[string]string)
	params["rec_modified_by"] = lib.UserIDStr
	params["rec_modified_date"] = time.Now().Format(lib.TIMESTAMPFORMAT)

	investPartnerKey := c.FormValue("invest_partner_key")
	if investPartnerKey == "" {
		return lib.CustomError(http.StatusBadRequest, "invest_partner_key can not be blank", "invest_partner_key can not be blank")
	}
	investPurposeKey := c.FormValue("invest_purpose_key")

	partnerCode := c.FormValue("partner_code")
	if partnerCode != "" {
		if len(partnerCode) > 30 {
			return lib.CustomError(http.StatusBadRequest, "partner_code must be <= 30 characters", "partner_code must be <= 30 characters")
		}
	} else {
		return lib.CustomError(http.StatusBadRequest, "partner_code cant be blank", "partner_code cant be blank")
	}

	partnerBusinessName := c.FormValue("partner_business_name")
	if partnerBusinessName != "" {
		if len(partnerBusinessName) > 50 {
			return lib.CustomError(http.StatusBadRequest, "partner_business_name must be <= 50 characters", "partner_business_name must be <= 50 characters")
		}
	}
	partnerDesc := c.FormValue("partner_desc")
	if partnerDesc != "" {
		if len(partnerDesc) > 150 {
			return lib.CustomError(http.StatusBadRequest, "partner_desc must be <= 150 characters", "partner_desc must be <= 150 characters")
		}
	}
	params["partner_desc"] = partnerDesc
	partnerPicName := c.FormValue("partner_picname")
	if partnerPicName != "" {
		if len(partnerPicName) > 150 {
			return lib.CustomError(http.StatusBadRequest, "partner_picname must be <= 150 characters", "partner_picname must be <= 150 characters")
		}
	}
	partnerMobileNo := c.FormValue("partner_mobileno")
	partnerOfficeNo := c.FormValue("partner_officeno")
	partnerEmail := c.FormValue("partner_email")
	if partnerEmail != "" {
		if len(partnerEmail) > 100 {
			return lib.CustomError(http.StatusBadRequest, "partner_email must be <= 100 characters", "partner_email must be <= 100 characters")
		}
	}
	partnerCity := c.FormValue("partner_city")
	if partnerCity != "" {
		if len(partnerCity) > 100 {
			return lib.CustomError(http.StatusBadRequest, "partner_city must be <= 100 characters", "partner_city must be <= 100 characters")
		}
	}
	partnerAddress := c.FormValue("partner_address")
	if partnerAddress != "" {
		if len(partnerAddress) > 200 {
			return lib.CustomError(http.StatusBadRequest, "partner_address must be <= 200 characters", "partner_address must be <= 200 characters")
		}
	}
	partnerUrl := c.FormValue("partner_url")
	if partnerUrl != "" {
		if len(partnerUrl) > 255 {
			return lib.CustomError(http.StatusBadRequest, "partner_url must be <= 200 characters", "partner_url must be <= 200 characters")
		}
	}
	partnerDateStarted := c.FormValue("partner_date_started")
	if partnerDateStarted != "" {
		expectedDateFormat := "2006-01-02"
		_, err = time.Parse(expectedDateFormat, partnerDateStarted)
		if err != nil {
			return lib.CustomError(http.StatusBadRequest, "partner_date_started should be a valid date in the format "+expectedDateFormat, "partner_date_started should be a valid date in the format "+expectedDateFormat)
		}
	} else {
		params["partner_date_started"] = "0000-00-00"
	}
	partnerDateExpired := c.FormValue("partner_date_expired")
	if partnerDateExpired != "" {
		expectedDateFormat := "2006-01-02"
		_, err = time.Parse(expectedDateFormat, partnerDateExpired)
		if err != nil {
			return lib.CustomError(http.StatusBadRequest, "partner_date_expired should be a valid date in the format "+expectedDateFormat, "partner_date_expired should be a valid date in the format "+expectedDateFormat)
		}
	} else {
		params["partner_date_expired"] = "0000-00-00"
	}
	partnerBannerHits := c.FormValue("partner_banner_hits")
	if partnerBannerHits != "" {
		if len(partnerBannerHits) > 11 {
			return lib.CustomError(http.StatusBadRequest, "partner_banner_hits must be <= 11 characters", "partner_banner_hits must be <= 11 characters")
		}
		value, err := strconv.Atoi(partnerBannerHits)
		if err != nil {
			return lib.CustomError(http.StatusBadRequest, "partner_banner_hits should be number", "partner_banner_hits should be number")
		}
		params["partner_banner_hits"] = strconv.Itoa(value)
	} else {
		params["partner_banner_hits"] = "0"
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
	params["invest_purpose_key"] = investPurposeKey
	params["partner_code"] = partnerCode
	params["partner_business_name"] = partnerBusinessName
	params["partner_picname"] = partnerPicName
	params["partner_mobileno"] = partnerMobileNo
	params["partner_officeno"] = partnerOfficeNo
	params["partner_email"] = partnerEmail
	params["partner_city"] = partnerCity
	params["partner_address"] = partnerAddress
	params["partner_url"] = partnerUrl
	params["rec_status"] = "1"

	status, err = models.UpdateInvestPartner(investPartnerKey, params)
	if err != nil {
		return lib.CustomError(status, err.Error(), "Failed input data")
	}
	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = "Data updated successfully"

	return c.JSON(http.StatusOK, response)
}
