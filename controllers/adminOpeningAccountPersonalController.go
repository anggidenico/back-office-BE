package controllers

import (
	"mf-bo-api/config"
	"mf-bo-api/lib"
	"mf-bo-api/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

func GetNewOAList(c echo.Context) error {
	errorAuth := initAuthCsKyc()
	if errorAuth != nil {
		return lib.CustomError(http.StatusUnauthorized, "You not allowed to access this page", "You not allowed to access this page")
	}
	var err error
	// var status int
	var responseData []models.OaRequestListResponse

	limitStr := c.QueryParam("limit")
	var limit uint64
	if limitStr != "" {
		limit, err = strconv.ParseUint(limitStr, 10, 64)
		if err == nil {
			if (limit == 0) || (limit > config.LimitQuery) {
				limit = config.LimitQuery
			}
		} else {
			return lib.CustomError(http.StatusBadRequest, "Limit should be number", "Limit should be number")
		}
	} else {
		limit = config.LimitQuery
	}
	// Get parameter page
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

	var getList []models.OaRequestListResponse
	models.GetOpeningAccountIndividuListQuery(&getList, lib.Profile.RoleKey, limit, offset)
	if len(getList) > 0 {
		responseData = getList
	} else {
		return lib.CustomError(http.StatusNotFound, "No Opening Account Data", "No Opening Account Data")
	}

	var response lib.ResponseWithPagination
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Pagination = 0
	response.Data = responseData

}
