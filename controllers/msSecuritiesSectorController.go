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
	// log.Printf("Response Data: %+v\n", response.Data)
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
	secParKey := c.FormValue("sector_parent_key")
	if secParKey != "" {
		if err != nil {
			return lib.CustomError(http.StatusBadRequest, "tolong dipilih gan", "tolong pilih gan")
		}
		params["sector_parent_key"] = secParKey
	}
	// else {
	// 	params["sector_parent_key"] = "NULL" // Set ke string "NULL" untuk kasus ini
	// }
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
	// params["sector_parent_key"] = SectorParentKey
	params["sector_code"] = sectorCode
	params["sector_name"] = sectorName
	params["sector_description"] = sectorDesc
	params["rec_status"] = "1"

	duplicate, key, err := models.CheckDuplicateSecuritiesSector(sectorCode, sectorName)
	if err != nil {
		log.Println("Error checking for duplicates:", err)
		return lib.CustomError(http.StatusInternalServerError, "Error checking for duplicates", "Error checking for duplicates")
	}

	log.Println("Duplicate:", duplicate)
	log.Println("Key:", key)

	// Jika duplikasi ditemukan, perbarui data yang sudah ada
	if duplicate {
		log.Println("Duplicate data found.")
		// Cek apakah data yang sudah ada masih aktif atau sudah dihapus
		existingDataStatus, err := models.GetSecuritiesSectorStatusByKey(key)
		if err != nil {
			log.Println("Error getting existing data status:", err)
			return lib.CustomError(http.StatusInternalServerError, "Error getting existing data status", "Error getting existing data status")
		}

		// Jika data sudah dihapus (rec_status = 0), perbarui statusnya menjadi aktif (rec_status = 1)
		if existingDataStatus == 0 {
			log.Println("Existing data is deleted. Recreating data.")

			// Set status menjadi aktif (rec_status = 1)
			params["rec_status"] = "1"
			// Update data dengan status baru dan nilai-nilai yang baru
			status, err := models.UpdateSecuritiesSector(key, params)
			if err != nil {
				log.Println("Error updating data:", err)
				return lib.CustomError(status, "Error updating data", "Error updating data")
			}
			return c.JSON(http.StatusOK, lib.Response{
				Status: lib.Status{
					Code:          http.StatusOK,
					MessageServer: "OK",
					MessageClient: "OK",
				},
				Data: "Data updated successfully",
			})
		} else {
			// Jika data masih aktif, kembalikan respons kesalahan duplikasi
			log.Println("Existing data is still active. Duplicate data error.")
			return lib.CustomError(http.StatusBadRequest, "Duplicate data. Unable to input data.", "Duplicate data. Unable to input data.")
		}
	} else {
		// Jika tidak ada duplikasi, buat data baru
		status, err := models.CreateSecuritiesSector(params)
		if err != nil {
			log.Println("Error create data:", err)
			return lib.CustomError(status, "Duplicate data. Unable to input data.", "Duplicate data. Unable to input data.")
		}
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

	secParKey := c.FormValue("sector_parent_key")
	if secParKey != "" {
		if err != nil {
			return lib.CustomError(http.StatusBadRequest, "tolong dipilih gan", "tolong pilih gan")
		}
		params["sector_parent_key"] = secParKey
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

	params["sector_code"] = sectorCode
	params["sector_name"] = sectorName
	params["sector_description"] = sectorDesc
	params["rec_status"] = "1"

	duplicate, key, err := models.CheckDuplicateSecuritiesSector(sectorCode, sectorName)
	if err != nil {
		log.Println("Error checking for duplicates:", err)
		return lib.CustomError(http.StatusInternalServerError, "Error checking for duplicates", "Error checking for duplicates")
	}
	if duplicate {
		log.Println("Duplicate data found.")
		// Cek apakah data yang sudah ada masih aktif atau sudah dihapus
		_, err := models.GetSecuritiesSectorStatusByKey(key)
		if err != nil {
			log.Println("Error getting existing data status:", err)
			return lib.CustomError(http.StatusBadRequest, "Duplicate data. Unable to input data.", "Duplicate data. Unable to input data.")
		}
		if key != sectorKey {
			return lib.CustomError(http.StatusBadRequest, "Duplicate data", "Duplicate data")
		}
	}
	status, err = models.UpdateSecuritiesSector(sectorKey, params)
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
