package controllers

import (
	"mf-bo-api/lib"
	"mf-bo-api/models"
	"net/http"
	"strconv"

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
