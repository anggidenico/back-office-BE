package controllers

import (
	"mf-bo-api/lib"
	"mf-bo-api/models"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo"
)

func GetUserBackOfficeList(c echo.Context) error {
	var err error
	params := make(map[string]string)

	limitStr := c.QueryParam("limit")
	var limit uint64
	if limitStr != "" {
		limit, err = strconv.ParseUint(limitStr, 10, 64)
		if err != nil {
			return lib.CustomError(http.StatusBadRequest, "Limit should be number", "Limit should be number")
		}
	}
	pageStr := c.QueryParam("page")
	var page uint64
	if pageStr != "" {
		page, err = strconv.ParseUint(pageStr, 10, 64)
		if err == nil {
			if page == 0 {
				page = 1
			}
		} else {
			return lib.CustomError(http.StatusBadRequest, "Page should be number", "Page should be number")
		}
	} else {
		page = 1
	}
	var offset uint64
	if page > 1 {
		offset = limit * (page - 1)
	}

	email := c.QueryParam("email")
	if email != "" {
		params["email"] = email
	}

	var responseData []models.UserCategoryBackOfficeList
	var pagination int
	responseData, pagination = models.GetUserCategoryBackOfficeList(params, limit, offset)

	var response lib.ResponseWithPagination
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Pagination = pagination
	response.Data = responseData

	return c.JSON(http.StatusOK, response)
}

func GetUserCustomerList(c echo.Context) error {
	var err error
	params := make(map[string]string)

	limitStr := c.QueryParam("limit")
	var limit uint64
	if limitStr != "" {
		limit, err = strconv.ParseUint(limitStr, 10, 64)
		if err != nil {
			return lib.CustomError(http.StatusBadRequest, "Limit should be number", "Limit should be number")
		}
	}
	pageStr := c.QueryParam("page")
	var page uint64
	if pageStr != "" {
		page, err = strconv.ParseUint(pageStr, 10, 64)
		if err == nil {
			if page == 0 {
				page = 1
			}
		} else {
			return lib.CustomError(http.StatusBadRequest, "Page should be number", "Page should be number")
		}
	} else {
		page = 1
	}
	var offset uint64
	if page > 1 {
		offset = limit * (page - 1)
	}

	email := c.QueryParam("email")
	if email != "" {
		params["email"] = email
	}
	phone_mobile := c.QueryParam("phone_mobile")
	if phone_mobile != "" {
		params["phone_mobile"] = phone_mobile
	}
	full_name := c.QueryParam("full_name")
	if full_name != "" {
		params["full_name"] = full_name
	}
	user_category := c.QueryParam("user_category")
	if user_category != "" {
		params["user_category"] = user_category
	}

	var responseData []models.UserCategoryCustomerList
	var pagination int
	responseData, pagination = models.GetUserCategoryCustomerList(params, limit, offset)

	var response lib.ResponseWithPagination
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Pagination = pagination
	response.Data = responseData

	return c.JSON(http.StatusOK, response)
}

func UserCustomerEdit(c echo.Context) error {
	var err error
	UpdtScUser := make(map[string]string)

	userLoginKey := c.FormValue("user_login_key")
	if userLoginKey == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing user_login_key", "Missing user_login_key")
	} else {
		UpdtScUser["user_login_key"] = userLoginKey
	}

	enabled := c.FormValue("enabled")
	if enabled != "" {
		if enabled == "true" {
			UpdtScUser["ulogin_enabled"] = "1"
		} else if enabled == "false" {
			UpdtScUser["ulogin_enabled"] = "0"
		}
	}

	locked := c.FormValue("locked")
	if locked != "" {
		if locked == "true" {
			UpdtScUser["ulogin_locked"] = "1"
		} else if locked == "false" {
			UpdtScUser["ulogin_locked"] = "0"
			UpdtScUser["ulogin_failed_count"] = "0"
		}
	}

	verifEmail := c.FormValue("verified_email")
	if verifEmail != "" {
		if verifEmail == "true" {
			UpdtScUser["verified_email"] = "1"
		} else if verifEmail == "false" {
			UpdtScUser["verified_email"] = "0"
		}
	}

	verifMobileno := c.FormValue("verified_mobileno")
	if verifMobileno != "" {
		if verifMobileno == "true" {
			UpdtScUser["verified_email"] = "1"
			UpdtScUser["last_verified_mobileno"] = time.Now().Format(lib.TIMESTAMPFORMAT)
		} else if verifMobileno == "false" {
			UpdtScUser["verified_email"] = "0"
		}
	}

	UpdtScUser["rec_modified_date"] = time.Now().Format(lib.TIMESTAMPFORMAT)
	UpdtScUser["rec_modified_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
	_, err = models.UpdateScUserLogin(UpdtScUser)
	if err != nil {
		return lib.CustomError(http.StatusBadRequest, err.Error(), "Failed update")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "Berhasil Update Data Customer"
	response.Data = nil

	return c.JSON(http.StatusOK, response)
}
