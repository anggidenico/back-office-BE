package controllers

import (
	"log"
	"mf-bo-api/lib"
	"mf-bo-api/models"
	"net/http"
	"strconv"
	"strings"
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

func GetOaRequestListKycApproved(c echo.Context) error {

	var responseData []models.OaRequestListKYCApprove
	result := models.GetOaRequestKYCApproveListQuery()

	if len(result) > 0 {
		for _, rData := range result {
			datas := rData

			if rData.DateBirth != nil {
				layoutDateBirth := "02 Jan 2006"
				t1, _ := time.Parse(lib.TIMESTAMPFORMAT, *rData.DateBirth)
				dateBirth := t1.Format(layoutDateBirth)
				datas.DateBirth = &dateBirth
			}

			layout := "02 Jan 2006 15:04"
			t2, _ := time.Parse(lib.TIMESTAMPFORMAT, *rData.OaDate)
			*datas.OaDate = t2.Format(layout)

			responseData = append(responseData, datas)
		}
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = responseData

	return c.JSON(http.StatusOK, response)
}

func DownloadOaRequestTextFile(c echo.Context) error {
	// var err error
	var responseData []models.OaRequestCsvFormatFiksTxt

	errorAuth := initAuthFundAdmin()
	if errorAuth != nil {
		// log.Error("User Autorizer")
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}

	OaRequestyKeys := c.QueryParam("oa_request_key")
	if OaRequestyKeys == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing: oa_request_key", "Missing: oa_request_key")
	}

	GetOaRequest := models.GetOaRequestKYCApproveListByRequestKey(OaRequestyKeys)

	txtHeader := "Type|SA Code|SID|First Name|Middle Name| Last Name|Country of Nationality|ID no |ID Expiration Date| NPWP no|NPWP Registration Date|Country of Birth|Place of Birth|Date of Birth |Gender|Educational Background|Mothers Maiden Name|Religion|Occupation|Income Level (IDR)|Marital Status|Spouses Name|Investors Risk Profile|Investment Objective|Source of Fund|Asset Owner|KTP Address|KTP City Code|KTP Postal Code|Correspondence Address|Correspondence City Code|Correspondence City Name|Correspondence Postal Code|Country of Correspondence|Domicile Address|Domicile City Code|Domicile City Name|Domicile Postal Code|Country of Domicile|Home Phone|Mobile Phone|Facsimile|Email|Statement Type|FATCA (Status)|TIN / Foreign TIN|TIN / Foreign TIN Issuance Country|REDM Payment Bank BIC Code 1|REDM Payment Bank BI Member Code 1|REDM Payment Bank Name 1| REDM Payment Bank Country 1| REDM Payment Bank Branch 1|REDM Payment A/C CCY 1|REDM Payment A/C No. 1|REDM Payment A/C Name 1|REDM Payment Bank BIC Code 2|REDM Payment Bank BI Member Code 2|REDM Payment Bank Name 2|REDM Payment Bank Country 2|REDM Payment Bank Branch 2|REDM Payment A/C CCY 2|REDM Payment A/C No. 2|REDM Payment A/C Name 2|REDM Payment Bank BIC Code 3|REDM Payment Bank BI Member Code 3| REDM Payment Bank Name 3|REDM Payment Bank Country 3|REDM Payment Bank Branch 3|REDM Payment A/C CCY 3|REDM Payment A/C No. 3|REDM Payment A/C Name 3|Client Code"
	var dataRow models.OaRequestCsvFormatFiksTxt
	dataRow.DataRow = txtHeader
	responseData = append(responseData, dataRow)

	if len(GetOaRequest) > 0 {
		for _, oarData := range GetOaRequest {

			// JIKA OA_FIRST / BARU PERTAMA KALI REGISTRASI
			if *oarData.OaRequestTypeInt == uint64(lib.OA_REQ_TYPE_NEW_INT) {
				var data models.OaRequestCsvFormat
				personalData := GetThePersonalDataDetails(strconv.FormatUint(oarData.OaRequestKey, 10))

				sliceName := strings.Fields(*personalData.FullName)
				firstName := ""
				middleName := ""
				lastName := ""

				if len(sliceName) > 0 {
					if len(sliceName) == 1 {
						firstName = sliceName[0]
						lastName = sliceName[0]
					}
					if len(sliceName) == 2 {
						firstName = sliceName[0]
						lastName = sliceName[1]
					}
					if len(sliceName) > 2 {
						ln := len(sliceName)
						firstName = sliceName[0]
						middleName = sliceName[1]
						lastName = strings.Join(sliceName[2:ln], " ")
						// lastName = lastName
					}
				}

				newLayout := "20060102"
				date1, _ := time.Parse(lib.TIMESTAMPFORMAT, *oarData.DateBirth)

				if *personalData.RelationTypeKey == uint64(87) || *personalData.RelationTypeKey == uint64(96) {
					data.SpouseName = *personalData.RelationFullName
				} else {
					data.SpouseName = ""
				}

				QuizResult := models.GetRiskProfileQuizResult(strconv.FormatUint(oarData.OaRequestKey, 10))

				data.Type = "1"
				data.SACode = lib.SA_CODE_MAM
				data.SID = ""
				data.FirstName = firstName
				data.MiddleName = middleName
				data.LastName = lastName
				data.CountryOfNationality = *personalData.CountryCode
				data.IDno = *personalData.IdCardNo
				data.IDExpirationDate = ""
				data.NpwpNo = ""
				data.NpwpRegistrationDate = ""
				data.CountryOfBirth = *personalData.CountryCode
				data.PlaceOfBirth = *personalData.PlaceBirth
				data.DateOfBirth = date1.Format(newLayout)
				data.Gender = *personalData.Gender
				data.EducationalBackground = strconv.FormatUint(*personalData.EducationKey, 10)
				data.MotherMaidenName = *personalData.MotherMaidenName
				data.Religion = strconv.FormatUint(*personalData.ReligionKey, 10)
				data.Occupation = strconv.FormatUint(*personalData.OccupJobKey, 10)
				data.IncomeLevel = strconv.FormatUint(*personalData.AnnualIncomeKey, 10)
				data.MaritalStatus = strconv.FormatUint(*personalData.MaritalStatusKey, 10)
				data.InvestorRiskProfile = strconv.FormatUint(QuizResult.RiskProfileKey, 10)
				data.InvestmentObjective = strconv.FormatUint(*personalData.InvesmentObjectivesKey, 10)
				data.SourceOfFund = strconv.FormatUint(*personalData.SourceOfFundkey, 10)
				data.AssetOwner = "1"
				data.KTPAddress = *personalData.IdCardAddress
				data.KTPCityCode = *personalData.IdCardCityCode
				data.KTPPostalCode = *personalData.IdCardPostalCode
				data.CorrespondenceAddress = ""
				data.CorrespondenceCityCode = ""
				data.CorrespondenceCityName = ""
				data.CorrespondencePostalCode = ""
				data.CountryOfCorrespondence = ""
				data.DomicileAddress = *personalData.DomicileAddress
				data.DomicileCityCode = *personalData.DomicileCityCode
				data.DomicileCityName = *personalData.DomicileCity
				data.DomicilePostalCode = *personalData.DomicilePostalCode
				data.CountryOfDomicile = *personalData.CountryCode
				data.HomePhone = *personalData.PhoneHome
				data.MobilePhone = *personalData.PhoneMobile
				data.Facsimile = ""
				data.Email = *personalData.EmailAddress

			}
		}
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = responseData

	return c.JSON(http.StatusOK, response)
}
