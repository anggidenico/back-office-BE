package controllers

import (
	"mf-bo-api/config"
	"mf-bo-api/lib"
	"mf-bo-api/models"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo"
)

func GetPengkinianPersonalDataList(c echo.Context) error {
	errorAuth := initAuthCsKyc()
	if errorAuth != nil {
		return lib.CustomError(http.StatusUnauthorized, "You not allowed to access this page", "You not allowed to access this page")
	}
	var err error
	// var status int
	var responseData []models.OaRequestListResponse

	RequestType := uint64(296)

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
	if page > 1 {
		offset = limit * (page - 1)
	}

	var getList []models.OaRequestListResponse
	pagination := models.GetOARequestIndividuListQuery(&getList, RequestType, lib.Profile.RoleKey, limit, offset)
	if len(getList) > 0 {
		// responseData = getList
		for _, getData := range getList {
			respData := getData

			layout := "02 Jan 2006 15:04"
			layoutDateBirth := "02 Jan 2006"

			t1, _ := time.Parse(lib.TIMESTAMPFORMAT, getData.DateBirth)
			respData.DateBirth = t1.Format(layoutDateBirth)

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

func GetPengkinianPersonalDataDetails(c echo.Context) error {

	oaRequestKey := c.Param("key")
	if oaRequestKey == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing: oaRequestKey")
	}

	var responseData models.PengkinianPersonalDataResponse

	var getPersonalData models.PengkinianPersonalDataModels
	err := models.GetPersonalDataOnlyQuery(&getPersonalData, oaRequestKey)
	if err != nil {
		return lib.CustomError(http.StatusInternalServerError, err.Error(), err.Error())
	}

	var getBankAccount []models.OaRequestBankDetails
	err = models.GetPengkinianBankAccount(&getBankAccount, oaRequestKey)
	if err != nil {
		return lib.CustomError(http.StatusInternalServerError, err.Error(), err.Error())
	}
	responseData.BankAccount = &getBankAccount

	responseData.Agent = getPersonalData.Agent
	responseData.AnnualIncome = getPersonalData.AnnualIncome
	responseData.BankAccount = &getBankAccount
	responseData.BeneficialFullName = getPersonalData.BeneficialFullName
	responseData.BeneficialRelation = getPersonalData.BeneficialRelation
	responseData.Branch = getPersonalData.Branch
	responseData.DateBirth = getPersonalData.DateBirth
	responseData.DomicileAddress = getPersonalData.DomicileAddress
	responseData.DomicileCity = getPersonalData.DomicileCity
	responseData.DomicilePostalCode = getPersonalData.DomicilePostalCode
	responseData.DomicileProvince = getPersonalData.DomicileProvince
	responseData.Education = getPersonalData.Education
	responseData.EmergencyFullName = getPersonalData.EmergencyFullName
	responseData.EmergencyPhoneNo = getPersonalData.EmergencyPhoneNo
	responseData.EmergencyRelation = getPersonalData.EmergencyRelation
	responseData.FullName = getPersonalData.FullName
	responseData.Gender = getPersonalData.Gender
	responseData.IdCardAddress = getPersonalData.IdCardAddress
	responseData.IdCardCity = getPersonalData.IdCardCity
	responseData.IdCardPostalCode = getPersonalData.IdCardPostalCode
	responseData.IdCardProvince = getPersonalData.IdCardProvince
	responseData.IdCardType = getPersonalData.IdCardType
	responseData.InvesmentObjectives = getPersonalData.InvesmentObjectives
	responseData.MaritalStatus = getPersonalData.MaritalStatus
	responseData.MotherMaidenName = getPersonalData.MotherMaidenName
	responseData.Nationality = getPersonalData.Nationality
	responseData.OaEntryEnd = getPersonalData.OaEntryEnd
	responseData.OaEntryStart = getPersonalData.OaEntryStart
	responseData.OaRequestKey = getPersonalData.OaRequestKey
	responseData.OaRequestType = getPersonalData.OaRequestType
	responseData.OaRiskLevel = getPersonalData.OaRiskLevel
	responseData.OaStatus = getPersonalData.OaStatus
	responseData.OccupAddress = getPersonalData.OccupAddress
	responseData.OccupBusinessFields = getPersonalData.OccupBusinessFields
	responseData.OccupCompany = getPersonalData.OccupCompany
	responseData.OccupJob = getPersonalData.OccupJob
	responseData.OccupPosition = getPersonalData.OccupJob
	responseData.PepName = getPersonalData.PepName
	responseData.PepPosition = getPersonalData.PepPosition
	responseData.PepStatus = getPersonalData.PepStatus
	responseData.PhoneHome = getPersonalData.PhoneHome
	responseData.PhoneMobile = getPersonalData.PhoneMobile
	responseData.PlaceBirth = getPersonalData.PlaceBirth
	responseData.RelationBusinessFields = getPersonalData.RelationBusinessFields
	responseData.RelationFullName = getPersonalData.RelationFullName
	responseData.RelationOccupation = getPersonalData.RelationOccupation
	responseData.RelationType = getPersonalData.RelationType
	responseData.Religion = getPersonalData.Religion
	responseData.SalesCode = getPersonalData.SalesCode
	responseData.SiteReferrer = getPersonalData.SiteReferrer
	responseData.SourceOfFund = getPersonalData.SourceOfFund

	i := config.ImageUrl + "/images/user/" + strconv.FormatUint(getPersonalData.UserLoginKey, 10) + "/" + *getPersonalData.PicKtp
	responseData.PicKtp = &i
	i2 := config.ImageUrl + "/images/user/" + strconv.FormatUint(getPersonalData.UserLoginKey, 10) + "/" + *getPersonalData.PicSelfieKtp
	responseData.PicSelfieKtp = &i2

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = responseData

	return c.JSON(http.StatusOK, response)
}
