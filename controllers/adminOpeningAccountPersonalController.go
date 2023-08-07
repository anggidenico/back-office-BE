package controllers

import (
	"log"
	"mf-bo-api/lib"
	"mf-bo-api/models"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo"
)

func GetNewOAList(c echo.Context) error {
	errorAuth := initAuthCsKyc()
	if errorAuth != nil {
		return lib.CustomError(http.StatusUnauthorized, "You not allowed to access this page", "You not allowed to access this page")
	}
	var err error
	RequestType := uint64(127)

	var responseData []models.PengkinianListResponse

	limitStr := c.QueryParam("limit")
	log.Println(limitStr)
	var limit uint64
	if limitStr != "" {
		limit, err = strconv.ParseUint(limitStr, 10, 64)
		if err != nil {
			return lib.CustomError(http.StatusBadRequest, "Limit should be number", "Limit should be number")
		}
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

	var getList []models.PengkinianListResponse
	pagination := models.GetOARequestIndividuListQuery(&getList, RequestType, lib.Profile.RoleKey, limit, offset)
	if len(getList) > 0 {
		// responseData = getList
		for _, getData := range getList {
			respData := getData
			layout := "02 January 2006 15:04"
			layoutDateBirth := "02 January 2006"

			if getData.DateBirth != nil {
				t1, _ := time.Parse(lib.TIMESTAMPFORMAT, *getData.DateBirth)
				dateBirth := t1.Format(layoutDateBirth)
				respData.DateBirth = &dateBirth
			}

			t2, _ := time.Parse(lib.TIMESTAMPFORMAT, getData.OaDate)
			respData.OaDate = t2.Format(layout)
			responseData = append(responseData, respData)
		}
	} else {
		return lib.CustomError(http.StatusNotFound, "No Data", "No Data")
	}

	var response lib.ResponseWithPagination
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Pagination = pagination
	response.Data = responseData

	return c.JSON(http.StatusOK, response)

}
