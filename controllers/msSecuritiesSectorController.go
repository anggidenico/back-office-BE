package controllers

import (
	"database/sql"
	"errors"
	"log"
	"mf-bo-api/lib"
	"mf-bo-api/models"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo"
)

func GetSecuritiesSectorController(c echo.Context) error {
	var value []models.SecuritiesSector
	status, err := models.GetSecuritiesSectorModels(&value)
	if err != nil {
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = value
	log.Printf("Response Data: %+v\n", response.Data)
	return c.JSON(http.StatusOK, response)
}

func GetSecSectorDetailController(c echo.Context) error {
	sectorKey := c.Param("sector_key")
	if sectorKey == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing sector_key", "Missing sector_key")
	}
	var price models.SecuritiesSector
	status, err := models.GetSecuritiesSectorDetailModels(&price, sectorKey)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return lib.CustomError(http.StatusNotFound, "sector_key not found", "sector_key not found")
		}
		return lib.CustomError(status, err.Error(), err.Error())
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = price
	return c.JSON(http.StatusOK, response)
}

func CreateSecuritiesSectorController(c echo.Context) error {
	var err error
	params := make(map[string]string)
	params["rec_created_by"] = lib.UserIDStr
	params["rec_created_date"] = time.Now().Format(lib.TIMESTAMPFORMAT)

	sectorCode := c.FormValue("sector_code")
	if sectorCode != "" {
		if len(sectorCode) > 20 {
			return lib.CustomError(http.StatusBadRequest, "sector_code must be <= 20 characters", "sector_code must be <= 20 characters")
		}
		sectorCode = strings.ToUpper(sectorCode)
	}
	sectorName := c.FormValue("sector_name")
	if sectorName != "" {
		if len(sectorName) > 150 {
			return lib.CustomError(http.StatusBadRequest, "sector_name must be <= 150 characters", "sector_name must be <= 150 characters")
		}
		sectorName = strings.ToUpper(sectorName)
	} else {
		return lib.CustomError(http.StatusBadRequest, "sector_name can not be blank", "sector_name can not be blank")
	}
	sectorDesc := c.FormValue("sector_description")
	if sectorDesc != "" {
		if len(sectorDesc) > 255 {
			return lib.CustomError(http.StatusBadRequest, "sector_description must be <= 255 characters", "sector_description must be <= 255 characters")
		}
		sectorDesc = strings.ToUpper(sectorDesc)
	}
	SectorParentKey := c.FormValue("sector_parent_key")
	if SectorParentKey != "" {
		if len(SectorParentKey) > 11 {
			return lib.CustomError(http.StatusBadRequest, "sector_parent_key must be <= 11 characters", "sector_parent_key must be <= 11 characters")
		}
		_, err = strconv.Atoi(SectorParentKey)
		if err != nil {
			return lib.CustomError(http.StatusBadRequest, "sector_parent_key must be number", "sector_parent_key must be number")
		}
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
	params["sector_parent_key"] = SectorParentKey
	params["sector_code"] = sectorCode
	params["sector_name"] = sectorName
	params["sector_description"] = sectorDesc
	params["rec_status"] = "1"

	// Check for duplicate records
	duplicate, key, err := models.CheckDuplicateSecuritiesSector(sectorCode)
	if err != nil {
		log.Println("Error checking for duplicates:", err)
		return lib.CustomError(http.StatusInternalServerError, "Error checking for duplicates", "Error checking for duplicates")
	}
	log.Println("Duplicate:", duplicate)
	log.Println("Key:", key)
	// Jika duplikasi ditemukan, perbarui data yang sudah ada
	if duplicate {
		status, err := models.UpdateSecuritiesSector(key, params)
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
	status, err = models.CreateSecuritiesSector(params)
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
func UpdateSecuritiesSectorController(c echo.Context) error {
	var err error
	params := make(map[string]string)
	params["rec_modified_by"] = lib.UserIDStr
	params["rec_modified_date"] = time.Now().Format(lib.TIMESTAMPFORMAT)

	sectorKey := c.FormValue("sector_key")
	if sectorKey == "" {
		return lib.CustomError(http.StatusBadRequest, "sector_key can not be blank", "sector_key can not be blank")
	}
	sectorCode := c.FormValue("sector_code")
	if sectorCode != "" {
		if len(sectorCode) > 20 {
			return lib.CustomError(http.StatusBadRequest, "sector_code must be <= 20 characters", "sector_code must be <= 20 characters")
		}
		sectorCode = strings.ToUpper(sectorCode)
	}

	sectorName := c.FormValue("sector_name")
	if sectorName != "" {
		if len(sectorName) > 150 {
			return lib.CustomError(http.StatusBadRequest, "sector_name must be <= 150 characters", "sector_name must be <= 150 characters")
		}
		sectorName = strings.ToUpper(sectorName)
	} else {
		return lib.CustomError(http.StatusBadRequest, "sector_name can not be blank", "sector_name can not be blank")
	}

	sectorParentKey := c.FormValue("sector_parent_key")
	// var sectorParentKeyInt int
	if sectorParentKey != "" {
		if len(sectorParentKey) > 11 {
			return lib.CustomError(http.StatusBadRequest, "sector_parent_key must be <= 11 characters", "sector_parent_key must be <= 11 characters")
		}
		_, err := strconv.Atoi(sectorParentKey)
		if err != nil {
			return lib.CustomError(http.StatusBadRequest, "sector_parent_key must be a number", "sector_parent_key must be a number")
		}

		// Set nilai yang vali
	} else {
		// Jika sector_parent_key kosong, set ke nilai yang sesuai dengan definisi null pada tipe int
		sectorParentKey = "0" // Atau sesuaikan dengan nilai default yang sesuai
	}
	sectorDesc := c.FormValue("sector_description")
	if sectorDesc != "" {
		if len(sectorDesc) > 255 {
			return lib.CustomError(http.StatusBadRequest, "sector_description must be <= 255 characters", "sector_description must be <= 255 characters")
		}
		sectorDesc = strings.ToUpper(sectorDesc)
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
	params["sector_parent_key"] = sectorParentKey
	params["sector_code"] = sectorCode
	params["sector_name"] = sectorName
	params["sector_description"] = sectorDesc
	params["rec_status"] = "1"

	status, err = models.UpdateSecuritiesSector(sectorKey, params)
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

func DeleteSecuritiesSectorController(c echo.Context) error {
	params := make(map[string]string)
	dateLayout := "2006-01-02 15:04:05"
	params["rec_status"] = "0"
	params["rec_deleted_date"] = time.Now().Format(dateLayout)
	params["rec_deleted_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

	sectorKey := c.FormValue("sector_key")
	if sectorKey == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing sector_key", "Missing sector_key")
	}

	status, err := models.DeleteSecuritiesSectorModels(sectorKey, params)
	if err != nil {
		return lib.CustomError(status, err.Error(), err.Error())
	}
	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = "Successfully deleted Master Securities Sector"
	return c.JSON(http.StatusOK, response)
}
