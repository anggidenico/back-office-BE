package controllers

import (
	"database/sql"
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
	} else {
		_, err := strconv.ParseUint(periodeKey, 10, 64)
		if err != sql.ErrNoRows {
			// log.Error("Wrong input for parameter: country_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: periode_key", "Wrong input for parameter: periode_key")
		}
	}
	var detperiode models.FfsPeriodeDetail
	status, err := models.GetFfsPeriodeDetailModels(&detperiode, periodeKey)
	if err != nil {
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

	periodeDate := c.FormValue("periode_date")
	if periodeDate == "" {
		return lib.CustomError(http.StatusBadRequest, "periode_date can not be blank", "periode_date can not be blank")
	}
	periodeName := c.FormValue("periode_name")
	if periodeName == "" {
		return lib.CustomError(http.StatusBadRequest, "periode_name can not be blank", "periode_name can not be blank")
	}
	dateOpened := c.FormValue("date_opened")
	if dateOpened == "" {
		return lib.CustomError(http.StatusBadRequest, "date_opened can not be blank", "date_opened can not be blank")
	}
	dateClosed := c.FormValue("date_closed")
	if dateClosed == "" {
		return lib.CustomError(http.StatusBadRequest, "date_closed can not be blank", "date_closed can not be blank")
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

	status, err = models.CreateFfsPeriode(params)
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
	periodeName := c.FormValue("periode_name")
	if periodeName == "" {
		return lib.CustomError(http.StatusBadRequest, "periode_name can not be blank", "periode_name can not be blank")
	}
	dateOpened := c.FormValue("date_opened")
	if dateOpened == "" {
		return lib.CustomError(http.StatusBadRequest, "date_opened can not be blank", "date_opened can not be blank")
	}
	dateClosed := c.FormValue("date_closed")
	if dateClosed == "" {
		return lib.CustomError(http.StatusBadRequest, "date_closed can not be blank", "date_closed can not be blank")
	}
	status, err = models.UpdateFfsPeriode(periodeKey, params)
	if err != nil {
		return lib.CustomError(status, err.Error(), "Failed input data")
	}
	params["periode_key"] = periodeKey
	params["periode_date"] = periodeDate
	params["periode_name"] = periodeName
	params["date_opened"] = dateOpened
	params["date_closed"] = dateClosed
	params["rec_status"] = "1"

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = ""

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
