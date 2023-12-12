package controllers

import (
	"database/sql"
	"errors"
	"log"
	"mf-bo-api/lib"
	"mf-bo-api/models"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo"
)

func GetFfsPeriodeController(c echo.Context) error {
	var ffsperiode []models.FfsPeriode
	status, err := models.GetFfsPeriodeModels(&ffsperiode)
	if err != nil {
		return lib.CustomError(status, err.Error(), err.Error())
	}
	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = ffsperiode
	return c.JSON(http.StatusOK, response)
}

func GetFfsPeriodeDetailController(c echo.Context) error {
	periodeKey := c.Param("periode_key")
	if periodeKey == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing periode_key", "Missing periode_key")
	}
	var detperiode models.FfsPeriodeDetail
	status, err := models.GetFfsPeriodeDetailModels(&detperiode, periodeKey)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return lib.CustomError(http.StatusNotFound, "Periode key not found", "Periode key not found")
		}
		return lib.CustomError(status, err.Error(), err.Error())
	}
	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = detperiode
	return c.JSON(http.StatusOK, response)
}
func CreateFfsPeriodeController(c echo.Context) error {
	var err error
	params := make(map[string]string)
	params["rec_created_by"] = lib.UserIDStr
	params["rec_created_date"] = time.Now().Format(lib.TIMESTAMPFORMAT)

	periodeDate := c.FormValue("periode_date") //validate date
	if periodeDate == "" {
		return lib.CustomError(http.StatusBadRequest, "periode_date can not be blank", "periode_date can not be blank")
	}
	expectedDateFormat := "2006-01-02"
	_, err = time.Parse(expectedDateFormat, periodeDate)
	if err != nil {
		return lib.CustomError(http.StatusBadRequest, "periode_date should be a valid date in the format "+expectedDateFormat, "periode_date should be a valid date in the format "+expectedDateFormat)
	}
	periodeName := c.FormValue("periode_name")
	if periodeName == "" {
		return lib.CustomError(http.StatusBadRequest, "periode_name can not be blank", "periode_name can not be blank")
	}
	dateOpened := c.FormValue("date_opened")
	if dateOpened == "" {
		return lib.CustomError(http.StatusBadRequest, "date_opened can not be blank", "date_opened can not be blank")
	}
	_, err = time.Parse(expectedDateFormat, dateOpened)
	if err != nil {
		return lib.CustomError(http.StatusBadRequest, "date_opened should be a valid date in the format "+expectedDateFormat, "date_opened should be a valid date in the format "+expectedDateFormat)
	}
	dateClosed := c.FormValue("date_closed")
	if dateClosed == "" {
		return lib.CustomError(http.StatusBadRequest, "date_closed can not be blank", "date_closed can not be blank")
	}
	_, err = time.Parse(expectedDateFormat, dateClosed)
	if err != nil {
		return lib.CustomError(http.StatusBadRequest, "date_closed should be a valid date in the format "+expectedDateFormat, "date_closed should be a valid date in the format "+expectedDateFormat)
	}
	recOrder := c.FormValue("rec_order")
	if recOrder != "" {
		_, err := strconv.ParseUint(recOrder, 10, 64)
		if err != nil {
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: rec_order", "Wrong input for parameter: rec_order")
		}
		params["rec_order"] = recOrder
	}
	params["periode_date"] = periodeDate
	params["periode_name"] = periodeName
	params["date_opened"] = dateOpened
	params["date_closed"] = dateClosed
	params["rec_status"] = "1"

	duplicate, key, err := models.CheckDuplicatePeriode(periodeName)
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
		existingDataStatus, err := models.GetPeriodeStatusByKey(key)
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
			status, err := models.UpdateFfsPeriode(key, params)
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
		status, err := models.CreatePeriode(params)
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
func UpdateFfsPeriode(c echo.Context) error {
	var err error
	params := make(map[string]string)
	params["rec_modified_by"] = lib.UserIDStr
	params["rec_modified_date"] = time.Now().Format(lib.TIMESTAMPFORMAT)

	periodeKey := c.FormValue("periode_key")
	if periodeKey == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing risk_profile_key", "Missing risk_profile_key")
	}

	periodeDate := c.FormValue("periode_date")
	if periodeDate == "" {
		return lib.CustomError(http.StatusBadRequest, "periode_date can not be blank", "periode_date can not be blank")
	}
	expectedDateFormat := "2006-01-02"
	_, err = time.Parse(expectedDateFormat, periodeDate)
	if err != nil {
		return lib.CustomError(http.StatusBadRequest, "periode_date should be a valid date in the format "+expectedDateFormat, "periode_date should be a valid date in the format "+expectedDateFormat)
	}
	periodeName := c.FormValue("periode_name")
	if periodeName == "" {
		return lib.CustomError(http.StatusBadRequest, "periode_name can not be blank", "periode_name can not be blank")
	}
	dateOpened := c.FormValue("date_opened")
	if dateOpened == "" {
		return lib.CustomError(http.StatusBadRequest, "date_opened can not be blank", "date_opened can not be blank")
	}
	_, err = time.Parse(expectedDateFormat, dateOpened)
	if err != nil {
		return lib.CustomError(http.StatusBadRequest, "date_opened should be a valid date in the format "+expectedDateFormat, "date_opened should be a valid date in the format "+expectedDateFormat)
	}
	dateClosed := c.FormValue("date_closed")
	if dateClosed == "" {
		return lib.CustomError(http.StatusBadRequest, "date_closed can not be blank", "date_closed can not be blank")
	}
	_, err = time.Parse(expectedDateFormat, dateClosed)
	if err != nil {
		return lib.CustomError(http.StatusBadRequest, "date_closed should be a valid date in the format "+expectedDateFormat, "date_closed should be a valid date in the format "+expectedDateFormat)
	}

	params["periode_key"] = periodeKey
	params["periode_date"] = periodeDate
	params["periode_name"] = periodeName
	params["date_opened"] = dateOpened
	params["date_closed"] = dateClosed
	params["rec_status"] = "1"

	duplicate, key, err := models.CheckDuplicatePeriode(periodeName)
	if err != nil {
		log.Println("Error checking for duplicates:", err)
		return lib.CustomError(http.StatusInternalServerError, "Error checking for duplicates", "Error checking for duplicates")
	}
	if duplicate {
		log.Println("Duplicate data found.")
		// Cek apakah data yang sudah ada masih aktif atau sudah dihapus
		_, err := models.GetPeriodeStatusByKey(key)
		if err != nil {
			log.Println("Error getting existing data status:", err)
			return lib.CustomError(http.StatusBadRequest, "Duplicate data. Unable to input data.", "Duplicate data. Unable to input data.")
		}

		if key != periodeKey {
			return lib.CustomError(http.StatusBadRequest, "Duplicate data", "Duplicate data")
		}

	}
	status, err = models.UpdateFfsPeriode(periodeKey, params)
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

func DeleteFfsPeriode(c echo.Context) error {
	params := make(map[string]string)
	dateLayout := "2006-01-02 15:04:05"
	params["rec_status"] = "0"
	params["rec_deleted_date"] = time.Now().Format(dateLayout)
	params["rec_deleted_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

	periodeKey := c.FormValue("periode_key")
	if periodeKey == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing periode_key", "Missing periode_key")
	}

	status, err := models.DeleteFfsPeriode(periodeKey, params)
	if err != nil {
		return lib.CustomError(status, err.Error(), err.Error())
	}
	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "Berhasil hapus FFS Periode"
	response.Data = ""
	return c.JSON(http.StatusOK, response)
}
