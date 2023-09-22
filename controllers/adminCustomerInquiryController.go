package controllers

import (
	"mf-bo-api/lib"
	"mf-bo-api/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

func NewGetCustomerInquiryList(c echo.Context) error {
	var err error
	var responseData []models.CustomerIndividuListResponse
	params := make(map[string]string)

	limitStr := c.QueryParam("limit")
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
	if page > 0 {
		offset = limit * (page - 1)
	}

	cif := c.QueryParam("cif")
	if cif != "" {
		params["cif"] = cif
	}
	fullname := c.QueryParam("full_name")
	if fullname != "" {
		params["full_name"] = fullname
	}
	datebirth := c.QueryParam("date_birth")
	if datebirth != "" {
		params["date_birth"] = datebirth
	}
	idcard_no := c.QueryParam("idcard_no")
	if idcard_no != "" {
		params["idcard_no"] = idcard_no
	}
	mothermaidenname := c.QueryParam("mother_maiden_name")
	if mothermaidenname != "" {
		params["mother_maiden_name"] = mothermaidenname
	}
	branchKey := c.QueryParam("branch_key")
	if branchKey != "" {
		params["branch_key"] = branchKey
	}

	getData, pagination := models.GetCustomerListWithCondition(params, limit, offset)

	if len(getData) > 0 {

		for _, Data := range getData {
			tf := false
			var vData models.CustomerIndividuListResponse
			vData.BranchName = Data.BranchName
			vData.CIF = Data.CIF
			vData.CustomerKey = Data.CustomerKey
			vData.DateBirth = Data.DateBirth
			vData.Email = Data.Email
			vData.FullName = Data.FullName
			vData.IdCardNo = Data.IdCardNo
			vData.MotherMaidenName = Data.MotherMaidenName
			vData.OaRequestKey = Data.OaRequestKey
			vData.PhoneMobile = Data.PhoneMobile
			vData.SID = Data.SID

			if *Data.CIFSuspendFlag == 1 {
				tf = true
				vData.CIFSuspendFlag = &tf
			} else {
				tf = false
				vData.CIFSuspendFlag = &tf
			}

			responseData = append(responseData, vData)
		}
	}

	// valueMap := val.(map[string]interface{})
	// if val, ok := valueMap["customer_key"]; ok {
	// 	dataToAppend["customer_key"] = val.(string)
	// }

	var response lib.ResponseWithPagination
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Pagination = pagination
	response.Data = responseData

	return c.JSON(http.StatusOK, response)
}
