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
	var responseData []models.PengkinianListResponse
	RequestType := uint64(296) // PENGKINIAN PERSONAL DATA

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

	var getList []models.PengkinianListResponse
	pagination := models.GetOARequestIndividuListQuery(&getList, RequestType, lib.Profile.RoleKey, limit, offset)
	if len(getList) > 0 {
		// responseData = getList
		for _, getData := range getList {
			respData := getData

			layout := "02 Jan 2006 15:04"
			layoutDateBirth := "02 Jan 2006"

			if getData.DateBirth != nil {
				t1, _ := time.Parse(lib.TIMESTAMPFORMAT, *getData.DateBirth)
				dateBirth := t1.Format(layoutDateBirth)
				respData.DateBirth = &dateBirth
			}

			t2, _ := time.Parse(lib.TIMESTAMPFORMAT, getData.OaDate)
			respData.OaDate = t2.Format(layout)

			responseData = append(responseData, respData)
		}
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
	var responseData models.PengkinianPersonalDataCompareResponse

	responseData.New = GetThePersonalDataDetails(oaRequestKey)
	OldKey := models.GetLastActiveOaKeyByNewOaKey(oaRequestKey)
	if OldKey != nil {
		if *OldKey != "" {
			responseData.Old = GetThePersonalDataDetails(*OldKey)
		}
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = responseData

	return c.JSON(http.StatusOK, response)
}

func GetThePersonalDataDetails(OaRequestKey string) models.PengkinianPersonalDataResponse {
	var OaData models.PengkinianPersonalDataResponse

	var getPersonalData models.PengkinianPersonalDataModels
	err := models.GetPersonalDataOnlyQuery(&getPersonalData, OaRequestKey)
	if err != nil {
	}

	var getBankAccount []models.OaRequestBankDetails
	err = models.GetOaRequestBankAccountNew(&getBankAccount, OaRequestKey)
	if err != nil {
	}
	OaData.BankAccount = &getBankAccount

	layout := "02 Jan 2006 15:04"
	layout2 := "02 Jan 2006"
	var DateBirth string
	var EntryStart string
	var EntryEnd string

	if getPersonalData.DateBirth != nil {
		ts, _ := time.Parse(lib.TIMESTAMPFORMAT, *getPersonalData.DateBirth)
		DateBirth = ts.Format(layout2)
	}
	if getPersonalData.OaEntryStart != nil {
		t1, _ := time.Parse(lib.TIMESTAMPFORMAT, *getPersonalData.OaEntryStart)
		EntryStart = t1.Format(layout)
	}
	if getPersonalData.OaEntryEnd != nil {
		t2, _ := time.Parse(lib.TIMESTAMPFORMAT, *getPersonalData.OaEntryEnd)
		EntryEnd = t2.Format(layout)
	}

	OaData.Agent = getPersonalData.Agent
	OaData.AnnualIncome = getPersonalData.AnnualIncome
	OaData.AnnualIncomeKey = getPersonalData.AnnualIncomeKey
	OaData.BeneficialFullName = getPersonalData.BeneficialFullName
	OaData.BeneficialRelation = getPersonalData.BeneficialRelation
	OaData.Branch = getPersonalData.Branch
	OaData.DateBirth = &DateBirth
	OaData.DomicileAddress = getPersonalData.DomicileAddress
	OaData.DomicileCity = getPersonalData.DomicileCity
	OaData.DomicileCityCode = getPersonalData.DomicileCityCode
	OaData.DomicilePostalCode = getPersonalData.DomicilePostalCode
	OaData.DomicileProvince = getPersonalData.DomicileProvince
	OaData.Education = getPersonalData.Education
	OaData.EducationKey = getPersonalData.EducationKey
	OaData.EmergencyFullName = getPersonalData.EmergencyFullName
	OaData.EmergencyPhoneNo = getPersonalData.EmergencyPhoneNo
	OaData.EmergencyRelation = getPersonalData.EmergencyRelation
	OaData.EmailAddress = getPersonalData.EmailAddress
	OaData.FullName = getPersonalData.FullName
	OaData.Gender = getPersonalData.Gender
	OaData.IdCardNo = getPersonalData.IdCardNo
	OaData.IdCardAddress = getPersonalData.IdCardAddress
	OaData.IdCardCity = getPersonalData.IdCardCity
	OaData.IdCardCityCode = getPersonalData.IdCardCityCode
	OaData.IdCardPostalCode = getPersonalData.IdCardPostalCode
	OaData.IdCardProvince = getPersonalData.IdCardProvince
	OaData.IdCardType = getPersonalData.IdCardType
	OaData.InvesmentObjectives = getPersonalData.InvesmentObjectives
	OaData.MaritalStatus = getPersonalData.MaritalStatus
	OaData.MaritalStatusKey = getPersonalData.MaritalStatusKey
	OaData.MotherMaidenName = getPersonalData.MotherMaidenName
	OaData.Nationality = getPersonalData.Nationality
	OaData.CountryCode = getPersonalData.CountryCode
	OaData.OaEntryEnd = &EntryEnd
	OaData.OaEntryStart = &EntryStart
	OaData.OaRequestKey = getPersonalData.OaRequestKey
	OaData.OaRequestType = getPersonalData.OaRequestType
	OaData.OaRiskLevel = getPersonalData.OaRiskLevel
	OaData.OaStatus = getPersonalData.OaStatus
	OaData.OccupAddress = getPersonalData.OccupAddress
	OaData.OccupBusinessFields = getPersonalData.OccupBusinessFields
	OaData.OccupCompany = getPersonalData.OccupCompany
	OaData.OccupJob = getPersonalData.OccupJob
	OaData.OccupJobKey = getPersonalData.OccupJobKey
	OaData.OccupPosition = getPersonalData.OccupPosition
	OaData.PepName = getPersonalData.PepName
	OaData.PepPosition = getPersonalData.PepPosition
	OaData.PepStatus = getPersonalData.PepStatus
	OaData.PhoneHome = getPersonalData.PhoneHome
	OaData.PhoneMobile = getPersonalData.PhoneMobile
	OaData.PlaceBirth = getPersonalData.PlaceBirth
	OaData.RelationBusinessFields = getPersonalData.RelationBusinessFields
	OaData.RelationFullName = getPersonalData.RelationFullName
	OaData.RelationOccupation = getPersonalData.RelationOccupation
	OaData.RelationType = getPersonalData.RelationType
	OaData.RelationTypeKey = getPersonalData.RelationTypeKey
	OaData.Religion = getPersonalData.Religion
	OaData.ReligionKey = getPersonalData.ReligionKey
	OaData.SalesCode = getPersonalData.SalesCode
	OaData.SiteReferrer = getPersonalData.SiteReferrer
	OaData.SourceOfFund = getPersonalData.SourceOfFund
	OaData.OaSource = getPersonalData.OaSource

	PersonalDataKey := strconv.FormatUint(getPersonalData.PersonalDataKey, 10)

	if getPersonalData.ReligionKey != nil && *getPersonalData.ReligionKey == uint64(26) { // agama
		OaData.Religion = models.GetUdfValueData("1", PersonalDataKey)
	}

	if getPersonalData.OccupJobKey != nil && *getPersonalData.OccupJobKey == uint64(35) { // pekerjaan
		OaData.OccupJob = models.GetUdfValueData("2", PersonalDataKey)
	}

	if getPersonalData.EducationKey != nil && *getPersonalData.EducationKey == uint64(43) { // pendidikan
		OaData.Education = models.GetUdfValueData("3", PersonalDataKey)
	}

	if getPersonalData.OccupBusinessFieldsKey != nil && *getPersonalData.OccupBusinessFieldsKey == uint64(60) { // bisnis field
		OaData.OccupBusinessFields = models.GetUdfValueData("4", PersonalDataKey)
	}

	if getPersonalData.SourceOfFundKey != nil && *getPersonalData.SourceOfFundKey == uint64(76) { // sumber uang
		OaData.SourceOfFund = models.GetUdfValueData("5", PersonalDataKey)
	}

	if getPersonalData.InvesmentObjectivesKey != nil && *getPersonalData.InvesmentObjectivesKey == uint64(82) { // tujuan investasi
		OaData.InvesmentObjectives = models.GetUdfValueData("6", PersonalDataKey)
	}

	if getPersonalData.BeneficialRelationKey != nil && *getPersonalData.BeneficialRelationKey == uint64(106) { //
		OaData.BeneficialRelation = models.GetUdfValueData("7", PersonalDataKey)
	}

	if getPersonalData.RelationOccupationKey != nil && *getPersonalData.RelationOccupationKey == uint64(35) { //
		OaData.RelationOccupation = models.GetUdfValueData("8", PersonalDataKey)
	}

	if getPersonalData.RelationBusinessFieldsKey != nil && *getPersonalData.RelationBusinessFieldsKey == uint64(60) { //
		OaData.RelationBusinessFields = models.GetUdfValueData("9", PersonalDataKey)
	}

	if getPersonalData.OccupPositionKey != nil && *getPersonalData.OccupPositionKey == uint64(295) { //
		OaData.OccupPosition = models.GetUdfValueData("10", PersonalDataKey)
	}

	if getPersonalData.PicKtp != nil {
		i := config.ImageUrl + "/images/user/" + strconv.FormatUint(getPersonalData.UserLoginKey, 10) + "/" + *getPersonalData.PicKtp
		OaData.PicKtp = &i
	}
	if getPersonalData.PicSelfieKtp != nil {
		i2 := config.ImageUrl + "/images/user/" + strconv.FormatUint(getPersonalData.UserLoginKey, 10) + "/" + *getPersonalData.PicSelfieKtp
		OaData.PicSelfieKtp = &i2
	}
	if getPersonalData.SignatureImage != nil {
		i3 := config.ImageUrl + "/images/user/" + strconv.FormatUint(getPersonalData.UserLoginKey, 10) + "/" + *getPersonalData.SignatureImage
		OaData.SignatureImage = &i3
	}

	return OaData
}
