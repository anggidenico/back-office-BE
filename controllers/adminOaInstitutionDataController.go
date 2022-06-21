package controllers

import (
	"database/sql"
	"encoding/json"
	"math"
	"mf-bo-api/config"
	"mf-bo-api/lib"
	"mf-bo-api/models"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
)

func GetListOaInstitusiDataCS(c echo.Context) error {
	if strconv.FormatUint(lib.Profile.RoleKey, 10) != lib.ROLE_CS {
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}
	var oaStatus []string
	oaStatus = append(oaStatus, strconv.FormatUint(uint64(lib.DRAFT), 10))
	fieldNot := ""
	valueNot := ""
	params := make(map[string]string)
	params["o.rec_created_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
	params["o.oa_request_type"] = "127"
	return ListOaInstitusiData(c, oaStatus, fieldNot, valueNot, params)
}

func GetListOaInstitusiApproveCS(c echo.Context) error {
	if strconv.FormatUint(lib.Profile.RoleKey, 10) != lib.ROLE_CS {
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}
	var oaStatus []string
	oaStatus = append(oaStatus, strconv.FormatUint(uint64(lib.OA_ENTRIED), 10))
	fieldNot := "o.rec_modified_by"
	valueNot := strconv.FormatUint(lib.Profile.UserID, 10)
	params := make(map[string]string)
	params["o.oa_request_type"] = "127"
	return ListOaInstitusiData(c, oaStatus, fieldNot, valueNot, params)
}

func GetListOaInstitusiDataBranch(c echo.Context) error {
	if lib.Profile.UserCategoryKey != uint64(lib.USER_CAT_BRANCH) {
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}
	var oaStatus []string
	oaStatus = append(oaStatus, strconv.FormatUint(uint64(lib.DRAFT), 10))
	fieldNot := ""
	valueNot := ""
	params := make(map[string]string)
	params["o.rec_created_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
	params["o.oa_request_type"] = "127"
	if lib.Profile.BranchKey != nil {
		params["b.branch_key"] = strconv.FormatUint(*lib.Profile.BranchKey, 10)
	}
	return ListOaInstitusiData(c, oaStatus, fieldNot, valueNot, params)
}

func GetListOaInstitusiApproveKYC(c echo.Context) error {
	if strconv.FormatUint(lib.Profile.RoleKey, 10) != lib.ROLE_KYC {
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}
	var oaStatus []string
	oaStatus = append(oaStatus, strconv.FormatUint(uint64(lib.CS_APPROVED), 10))
	fieldNot := ""
	valueNot := ""
	params := make(map[string]string)
	params["o.oa_request_type"] = "127"
	return ListOaInstitusiData(c, oaStatus, fieldNot, valueNot, params)
}

func GetListOaInstitusiApproveFundAdmin(c echo.Context) error {
	if strconv.FormatUint(lib.Profile.RoleKey, 10) != lib.ROLE_FUND_ADMIN {
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}
	var oaStatus []string
	oaStatus = append(oaStatus, strconv.FormatUint(uint64(lib.KYC_APPROVED), 10))
	oaStatus = append(oaStatus, strconv.FormatUint(uint64(lib.CUSTOMER_BUILD), 10))
	oaStatus = append(oaStatus, strconv.FormatUint(uint64(lib.SINVEST_DONE), 10))
	fieldNot := ""
	valueNot := ""
	params := make(map[string]string)
	params["o.oa_request_type"] = "127"
	return ListOaInstitusiData(c, oaStatus, fieldNot, valueNot, params)
}

func ListOaInstitusiData(c echo.Context, oaStatus []string, fieldNot string, valueNot string, params map[string]string) error {
	var err error
	var status int

	//Get parameter limit
	limitStr := c.QueryParam("limit")
	var limit uint64
	if limitStr != "" {
		limit, err = strconv.ParseUint(limitStr, 10, 64)
		if err == nil {
			if (limit == 0) || (limit > config.LimitQuery) {
				limit = config.LimitQuery
			}
		} else {
			log.Error("Limit should be number")
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
			log.Error("Page should be number")
			return lib.CustomError(http.StatusBadRequest, "Page should be number", "Page should be number")
		}
	} else {
		page = 1
	}

	var offset uint64
	if page > 1 {
		offset = limit * (page - 1)
	}

	noLimitStr := c.QueryParam("nolimit")
	var noLimit bool
	if noLimitStr != "" {
		noLimit, err = strconv.ParseBool(noLimitStr)
		if err != nil {
			log.Error("Nolimit parameter should be true/false")
			return lib.CustomError(http.StatusBadRequest, "Nolimit parameter should be true/false", "Nolimit parameter should be true/false")
		}
	} else {
		noLimit = false
	}

	items := []string{"oa_request_key", "branch_name", "agent_name", "oa_status", "npwp", "full_name", "tanggal_pendirian", "tempat_pendirian", "nomor_akta", "tanggal_akta", "no_izin_usaha"}

	orderBy := c.QueryParam("order_by")
	if orderBy != "" {
		_, found := lib.Find(items, orderBy)
		if found {
			var ord string
			if orderBy == "branch_name" {
				ord = "b.branch_name"
			} else if orderBy == "oa_request_key" {
				ord = "o.oa_request_key"
			} else if orderBy == "agent_name" {
				ord = "a.agent_name"
			} else if orderBy == "oa_status" {
				ord = "o.oa_status"
			} else if orderBy == "npwp" {
				ord = "oad.tin_number"
			} else if orderBy == "full_name" {
				ord = "oad.full_name"
			} else if orderBy == "tanggal_pendirian" {
				ord = "oad.established_date"
			} else if orderBy == "tempat_pendirian" {
				ord = "oad.established_city"
			} else if orderBy == "nomor_akta" {
				ord = "oad.deed_no"
			} else if orderBy == "tanggal_akta" {
				ord = "oad.deed_date"
			} else if orderBy == "no_izin_usaha" {
				ord = "oad.biz_license"
			}
			params["orderBy"] = ord
			orderType := c.QueryParam("order_type")
			if (orderType == "asc") || (orderType == "ASC") || (orderType == "desc") || (orderType == "DESC") {
				params["orderType"] = orderType
			}
		} else {
			log.Error("Wrong input for parameter order_by")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter order_by", "Wrong input for parameter order_by")
		}
	} else {
		params["orderBy"] = "o.oa_request_key"
		params["orderType"] = "DESC"
	}

	branchKey := c.QueryParam("branch_key")
	if branchKey != "" {
		n, err := strconv.ParseUint(branchKey, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: branch_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: branch_key", "Wrong input for parameter: branch_key")
		}
		params["b.branch_key"] = branchKey
	}

	agentKey := c.QueryParam("agent_key")
	if agentKey != "" {
		n, err := strconv.ParseUint(agentKey, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: agent_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: agent_key", "Wrong input for parameter: agent_key")
		}
		params["a.agent_key"] = agentKey
	}

	searchLike := c.QueryParam("search_like")

	var oaInst []models.AdminListOaInstitutionData
	status, err = models.AdminGetListOaInstitutionData(&oaInst, oaStatus, fieldNot, valueNot, params, limit, offset, noLimit, searchLike)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if len(oaInst) < 1 {
		log.Error("oa institusi not found")
		return lib.CustomError(http.StatusNotFound, "Oa Request Institusi not found", "Oa Request Institusi not found")
	}

	var countData models.OaRequestCountData
	var pagination int
	if limit > 0 {
		status, err = models.AdminCountGetListOaInstitutionData(&countData, oaStatus, fieldNot, valueNot, params, searchLike)
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
		if int(countData.CountData) < int(limit) {
			pagination = 1
		} else {
			calc := math.Ceil(float64(countData.CountData) / float64(limit))
			pagination = int(calc)
		}
	} else {
		pagination = 1
	}

	var response lib.ResponseWithPagination
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Pagination = pagination
	response.Data = oaInst

	return c.JSON(http.StatusOK, response)
}

func GetDetailOaInstitusiAdmin(c echo.Context) error {
	requestKey := c.Param("request_key")
	if requestKey == "" {
		log.Error("Missing required parameter: request_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: request_key", "Missing required parameter: request_key")
	} else {
		n, err := strconv.ParseUint(requestKey, 10, 64)
		if err != nil || n <= 0 {
			log.Error("Wrong input for parameter: request_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: request_key", "Wrong input for parameter: request_key")
		}
	}

	oaDetail, err := DetailInstitution(requestKey)
	if err != nil {
		log.Error("Error get data detail oa institution")
		return lib.CustomError(http.StatusBadRequest, err.Error(), err.Error())
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = oaDetail
	return c.JSON(http.StatusOK, response)
}

func DetailInstitution(oaReqKey string) (models.OaInstitutionDetail, error) {
	var err error
	var status int
	log.Println(status)
	var oaDetail models.OaInstitutionDetail
	decimal.MarshalJSONWithoutQuotes = true

	branchKey := ""

	var userCategory uint64
	userCategory = 3 //Branch
	if lib.Profile.UserCategoryKey == userCategory {
		if lib.Profile.BranchKey != nil {
			branchKey = strconv.FormatUint(*lib.Profile.BranchKey, 10)
		}
	}

	var oareq models.OaRequest
	_, err = models.GetOaRequestInstitution(&oareq, oaReqKey, branchKey)
	if err != nil {
		return oaDetail, err
	}
	var oadata models.OaInstitutionData
	_, err = models.GetOaInstitutionData(&oadata, oaReqKey, "oa_request_key")
	if err != nil {
		return oaDetail, err
	}

	layout := "2006-01-02 15:04:05"
	newLayout := "02 Jan 2006"

	oaDetail.InstitutionDataKey = oadata.InstitutionDataKey
	oaDetail.OaRequestKey = oareq.OaRequestKey
	oaDetail.SalesCode = oareq.SalesCode
	oaDetail.Check1Date = oareq.Check1Date
	if oareq.Check1Date != nil {
		date, err := time.Parse(layout, *oareq.Check1Date)
		if err == nil {
			oke := date.Format(newLayout)
			oaDetail.Check1Date = &oke
		}
	}
	oaDetail.Check1Flag = oareq.Check1Flag
	oaDetail.Check1References = oareq.Check1References
	oaDetail.Check1Notes = oareq.Check1Notes
	oaDetail.Check2Date = oareq.Check2Date
	if oareq.Check2Date != nil {
		date, err := time.Parse(layout, *oareq.Check2Date)
		if err == nil {
			oke := date.Format(newLayout)
			oaDetail.Check2Date = &oke
		}
	}
	oaDetail.Check2Flag = oareq.Check2Flag
	oaDetail.Check2References = oareq.Check2References
	oaDetail.Check2Notes = oareq.Check2Notes
	oaDetail.FullName = oadata.FullName
	oaDetail.ShortName = oadata.ShortName
	oaDetail.TinNumber = oadata.TinNumber
	oaDetail.EstablishedCity = oadata.EstablishedCity
	oaDetail.EstablishedDate = oadata.EstablishedDate
	if oadata.EstablishedDate != nil {
		date, err := time.Parse(layout, *oadata.EstablishedDate)
		if err == nil {
			oke := date.Format(newLayout)
			oaDetail.EstablishedDate = &oke
		}
	}
	oaDetail.DeedNo = oadata.DeedNo
	oaDetail.DeedDate = oadata.DeedDate
	if oadata.DeedDate != nil {
		date, err := time.Parse(layout, *oadata.DeedDate)
		if err == nil {
			oke := date.Format(newLayout)
			oaDetail.DeedDate = &oke
		}
	}
	oaDetail.LsEstablishValidationNo = oadata.LsEstablishValidationNo
	oaDetail.LsEstablishValidationDate = oadata.LsEstablishValidationDate
	if oadata.LsEstablishValidationDate != nil {
		date, err := time.Parse(layout, *oadata.LsEstablishValidationDate)
		if err == nil {
			oke := date.Format(newLayout)
			oaDetail.LsEstablishValidationDate = &oke
		}
	}
	oaDetail.LastChangeAaNo = oadata.LastChangeAaNo
	oaDetail.LastChangeDaDate = oadata.LastChangeAaDate
	if oadata.LastChangeAaDate != nil {
		date, err := time.Parse(layout, *oadata.LastChangeAaDate)
		if err == nil {
			oke := date.Format(newLayout)
			oaDetail.LastChangeDaDate = &oke
		}
	}
	oaDetail.LsLastChangeAaNo = oadata.LsLastChangeAaNo
	oaDetail.LsLastChangeAaDate = oadata.LsLastChangeAaDate
	if oadata.LsLastChangeAaDate != nil {
		date, err := time.Parse(layout, *oadata.LsLastChangeAaDate)
		if err == nil {
			oke := date.Format(newLayout)
			oaDetail.LsLastChangeAaDate = &oke
		}
	}
	oaDetail.ManagementDeedNo = oadata.ManagementDeedNo
	oaDetail.ManagementDeedDate = oadata.ManagementDeedDate
	if oadata.ManagementDeedDate != nil {
		date, err := time.Parse(layout, *oadata.ManagementDeedDate)
		if err == nil {
			oke := date.Format(newLayout)
			oaDetail.ManagementDeedDate = &oke
		}
	}
	oaDetail.LsMgtChangeDeedNo = oadata.LsMgtChangeDeedNo
	oaDetail.LsMgtChangeDeedDate = oadata.LsMgtChangeDeedDate
	if oadata.LsMgtChangeDeedDate != nil {
		date, err := time.Parse(layout, *oadata.LsMgtChangeDeedDate)
		if err == nil {
			oke := date.Format(newLayout)
			oaDetail.LsMgtChangeDeedDate = &oke
		}
	}
	oaDetail.SkdLicenseNo = oadata.SkdLicenseNo
	oaDetail.SkdLicenseDate = oadata.SkdLicenseDate
	if oadata.SkdLicenseDate != nil {
		date, err := time.Parse(layout, *oadata.SkdLicenseDate)
		if err == nil {
			oke := date.Format(newLayout)
			oaDetail.SkdLicenseDate = &oke
		}
	}
	oaDetail.BizLicenseNo = oadata.BizLicenseNo
	oaDetail.BizLicenseDate = oadata.BizLicenseDate
	if oadata.BizLicenseDate != nil {
		date, err := time.Parse(layout, *oadata.BizLicenseDate)
		if err == nil {
			oke := date.Format(newLayout)
			oaDetail.BizLicenseDate = &oke
		}
	}
	oaDetail.NibNo = oadata.NibNo
	oaDetail.NibDate = oadata.NibDate
	if oadata.NibDate != nil {
		date, err := time.Parse(layout, *oadata.NibDate)
		if err == nil {
			oke := date.Format(newLayout)
			oaDetail.NibDate = &oke
		}
	}
	oaDetail.PhoneNo = oadata.PhoneNo
	oaDetail.MobileNo = oadata.MobileNo
	oaDetail.FaxNo = oadata.FaxNo
	oaDetail.EmailAddress = oadata.EmailAddress
	oaDetail.BoName = oadata.BoName
	oaDetail.BoIdnumber = oadata.BoIdnumber
	oaDetail.BoBusiness = oadata.BoBusiness
	oaDetail.BoIdaddress = oadata.BoIdaddress
	oaDetail.BoBusinessAddress = oadata.BoBusinessAddress
	oaDetail.BoAnnuallyIncome = oadata.BoAnnuallyIncome
	oaDetail.AssetY1 = oadata.AssetY1
	oaDetail.AssetY2 = oadata.AssetY2
	oaDetail.AssetY3 = oadata.AssetY3
	oaDetail.OpsProfitY1 = oadata.OpsProfitY1
	oaDetail.OpsProfitY2 = oadata.OpsProfitY2
	oaDetail.OpsProfitY3 = oadata.OpsProfitY3
	oaDetail.InstiRemarks = oadata.InstiRemarks
	oaDetail.InstitutionGroup = oadata.InstitutionGroup
	oaDetail.DocShipmentDate = oadata.DocShipmentDate
	if oadata.DocShipmentDate != nil {
		date, err := time.Parse(layout, *oadata.DocShipmentDate)
		if err == nil {
			oke := date.Format(newLayout)
			oaDetail.DocShipmentDate = &oke
		}
	}
	oaDetail.DocShipmentEmail = oadata.DocShipmentEmail
	oaDetail.DocShipmentNotes = oadata.DocShipmentNotes
	oaDetail.RecCreatedDate = oareq.RecCreatedDate
	if oareq.RecCreatedDate != nil {
		date, err := time.Parse(layout, *oareq.RecCreatedDate)
		if err == nil {
			oke := date.Format(newLayout)
			oaDetail.RecCreatedDate = &oke
		}
	}
	oaDetail.RecCreatedBy = oareq.RecCreatedBy
	oaDetail.RecModifiedDate = oareq.RecModifiedDate
	if oareq.RecModifiedDate != nil {
		date, err := time.Parse(layout, *oareq.RecModifiedDate)
		if err == nil {
			oke := date.Format(newLayout)
			oaDetail.RecModifiedDate = &oke
		}
	}
	oaDetail.RecModifiedBy = oareq.RecModifiedBy
	oaDetail.RecApprovalStatus = oareq.RecApprovalStatus
	oaDetail.RecApprovalStage = oareq.RecApprovalStage
	oaDetail.RecApprovedDate = oareq.RecApprovedDate
	if oareq.RecApprovedDate != nil {
		date, err := time.Parse(layout, *oareq.RecApprovedDate)
		if err == nil {
			oke := date.Format(newLayout)
			oaDetail.RecApprovedDate = &oke
		}
	}
	oaDetail.RecApprovedBy = oareq.RecApprovedBy

	//LOOKUP
	var lookupIds []string

	if oadata.IntitutionType != nil {
		if _, ok := lib.Find(lookupIds, strconv.FormatUint(*oadata.IntitutionType, 10)); !ok {
			lookupIds = append(lookupIds, strconv.FormatUint(*oadata.IntitutionType, 10))
		}
	}

	if oadata.IntitutionClassification != nil {
		if _, ok := lib.Find(lookupIds, strconv.FormatUint(*oadata.IntitutionClassification, 10)); !ok {
			lookupIds = append(lookupIds, strconv.FormatUint(*oadata.IntitutionClassification, 10))
		}
	}

	if oadata.IntitutionCharacteristic != nil {
		if _, ok := lib.Find(lookupIds, strconv.FormatUint(*oadata.IntitutionCharacteristic, 10)); !ok {
			lookupIds = append(lookupIds, strconv.FormatUint(*oadata.IntitutionCharacteristic, 10))
		}
	}

	if oadata.IntitutionBusinessType != nil {
		if _, ok := lib.Find(lookupIds, strconv.FormatUint(*oadata.IntitutionBusinessType, 10)); !ok {
			lookupIds = append(lookupIds, strconv.FormatUint(*oadata.IntitutionBusinessType, 10))
		}
	}

	if oadata.InstiAnnuallyIncome != nil {
		if _, ok := lib.Find(lookupIds, strconv.FormatUint(*oadata.InstiAnnuallyIncome, 10)); !ok {
			lookupIds = append(lookupIds, strconv.FormatUint(*oadata.InstiAnnuallyIncome, 10))
		}
	}

	if oadata.InstiSourceOfIncome != nil {
		if _, ok := lib.Find(lookupIds, strconv.FormatUint(*oadata.InstiSourceOfIncome, 10)); !ok {
			lookupIds = append(lookupIds, strconv.FormatUint(*oadata.InstiSourceOfIncome, 10))
		}
	}

	if oadata.InstiInvestmentPurpose != nil {
		if _, ok := lib.Find(lookupIds, strconv.FormatUint(*oadata.InstiInvestmentPurpose, 10)); !ok {
			lookupIds = append(lookupIds, strconv.FormatUint(*oadata.InstiInvestmentPurpose, 10))
		}
	}

	if oadata.BoRelation != nil {
		if _, ok := lib.Find(lookupIds, strconv.FormatUint(*oadata.BoRelation, 10)); !ok {
			lookupIds = append(lookupIds, strconv.FormatUint(*oadata.BoRelation, 10))
		}
	}

	if oadata.InstiDocShipment != nil {
		if _, ok := lib.Find(lookupIds, strconv.FormatUint(*oadata.InstiDocShipment, 10)); !ok {
			lookupIds = append(lookupIds, strconv.FormatUint(*oadata.InstiDocShipment, 10))
		}
	}

	if oareq.Oastatus != nil {
		if _, ok := lib.Find(lookupIds, strconv.FormatUint(*oareq.Oastatus, 10)); !ok {
			lookupIds = append(lookupIds, strconv.FormatUint(*oareq.Oastatus, 10))
		}
	}

	var lookupInst []models.GenLookup
	if len(lookupIds) > 0 {
		status, err = models.GetGenLookupIn(&lookupInst, lookupIds, "lookup_key")
		if err != nil {
			if err != sql.ErrNoRows {
				log.Error("err get lookup list")
				log.Error(err.Error())
			}
		}
	}

	pData := make(map[uint64]models.GenLookup)
	for _, genLook := range lookupInst {
		pData[genLook.LookupKey] = genLook
	}

	if oadata.IntitutionType != nil {
		if n, ok := pData[*oadata.IntitutionType]; ok {
			var dat models.LookupTrans
			dat.LookupKey = n.LookupKey
			dat.LkpGroupKey = n.LkpGroupKey
			dat.LkpCode = n.LkpCode
			dat.LkpName = n.LkpName
			oaDetail.IntitutionType = &dat
		}
	}

	if oadata.IntitutionClassification != nil {
		if n, ok := pData[*oadata.IntitutionClassification]; ok {
			var dat models.LookupTrans
			dat.LookupKey = n.LookupKey
			dat.LkpGroupKey = n.LkpGroupKey
			dat.LkpCode = n.LkpCode
			dat.LkpName = n.LkpName
			oaDetail.IntitutionClassification = &dat
		}
	}

	if oadata.IntitutionCharacteristic != nil {
		if n, ok := pData[*oadata.IntitutionCharacteristic]; ok {
			var dat models.LookupTrans
			dat.LookupKey = n.LookupKey
			dat.LkpGroupKey = n.LkpGroupKey
			dat.LkpCode = n.LkpCode
			dat.LkpName = n.LkpName
			oaDetail.IntitutionCharacteristic = &dat
		}
	}

	if oadata.IntitutionBusinessType != nil {
		if n, ok := pData[*oadata.IntitutionBusinessType]; ok {
			var dat models.LookupTrans
			dat.LookupKey = n.LookupKey
			dat.LkpGroupKey = n.LkpGroupKey
			dat.LkpCode = n.LkpCode
			dat.LkpName = n.LkpName
			oaDetail.IntitutionBusinessType = &dat
		}
	}

	if oadata.InstiAnnuallyIncome != nil {
		if n, ok := pData[*oadata.InstiAnnuallyIncome]; ok {
			var dat models.LookupTrans
			dat.LookupKey = n.LookupKey
			dat.LkpGroupKey = n.LkpGroupKey
			dat.LkpCode = n.LkpCode
			dat.LkpName = n.LkpName
			oaDetail.InstiAnnuallyIncome = &dat
		}
	}

	if oadata.InstiSourceOfIncome != nil {
		if n, ok := pData[*oadata.InstiSourceOfIncome]; ok {
			var dat models.LookupTrans
			dat.LookupKey = n.LookupKey
			dat.LkpGroupKey = n.LkpGroupKey
			dat.LkpCode = n.LkpCode
			dat.LkpName = n.LkpName
			oaDetail.InstiSourceOfIncome = &dat
		}
	}

	if oadata.InstiInvestmentPurpose != nil {
		if n, ok := pData[*oadata.InstiInvestmentPurpose]; ok {
			var dat models.LookupTrans
			dat.LookupKey = n.LookupKey
			dat.LkpGroupKey = n.LkpGroupKey
			dat.LkpCode = n.LkpCode
			dat.LkpName = n.LkpName
			oaDetail.InstiInvestmentPurpose = &dat
		}
	}

	if oadata.BoRelation != nil {
		if n, ok := pData[*oadata.BoRelation]; ok {
			var dat models.LookupTrans
			dat.LookupKey = n.LookupKey
			dat.LkpGroupKey = n.LkpGroupKey
			dat.LkpCode = n.LkpCode
			dat.LkpName = n.LkpName
			oaDetail.BoRelation = &dat
		}
	}

	if oadata.InstiDocShipment != nil {
		if n, ok := pData[*oadata.InstiDocShipment]; ok {
			var dat models.LookupTrans
			dat.LookupKey = n.LookupKey
			dat.LkpGroupKey = n.LkpGroupKey
			dat.LkpCode = n.LkpCode
			dat.LkpName = n.LkpName
			oaDetail.InstiDocShipment = &dat
		}
	}

	if oareq.Oastatus != nil {
		if n, ok := pData[*oareq.Oastatus]; ok {
			var dat models.LookupTrans
			dat.LookupKey = n.LookupKey
			dat.LkpGroupKey = n.LkpGroupKey
			dat.LkpCode = n.LkpCode
			dat.LkpName = n.LkpName
			oaDetail.OaStatus = &dat
		}
	}

	//ADDRESS
	var postalAddressIds []string
	if oadata.DomicileKey != nil {
		if _, ok := lib.Find(postalAddressIds, strconv.FormatUint(*oadata.DomicileKey, 10)); !ok {
			postalAddressIds = append(postalAddressIds, strconv.FormatUint(*oadata.DomicileKey, 10))
		}
	}
	if oadata.CorrespondenceKey != nil {
		if _, ok := lib.Find(postalAddressIds, strconv.FormatUint(*oadata.CorrespondenceKey, 10)); !ok {
			postalAddressIds = append(postalAddressIds, strconv.FormatUint(*oadata.CorrespondenceKey, 10))
		}
	}
	var address []models.AddressDetail
	if len(postalAddressIds) > 0 {
		status, err = models.GetOaPostalAddressDetailIn(&address, postalAddressIds, "postal_address_key")
		if err != nil {
			if err != sql.ErrNoRows {
				log.Error(err.Error())
				log.Error("ERROR get postal address detail")
			}
		}
	}
	postalData := make(map[uint64]models.AddressDetail)
	for _, posAdd := range address {
		postalData[posAdd.PostalAddressKey] = posAdd
	}

	if len(postalData) > 0 {
		if oadata.DomicileKey != nil {
			if p, ok := postalData[*oadata.DomicileKey]; ok {
				var dat models.AddressDetail
				dat.PostalAddressKey = p.PostalAddressKey
				dat.AddressType = p.AddressType
				dat.AddressTypeName = p.AddressTypeName
				dat.ProvinsiKey = p.ProvinsiKey
				dat.ProvinsiName = p.ProvinsiName
				dat.KabupatenKey = p.KabupatenKey
				dat.KabupatenName = p.KabupatenName
				dat.KecamatanKey = p.KecamatanKey
				dat.KecamatanName = p.KecamatanName
				dat.AddressLine1 = p.AddressLine1
				dat.PostalCode = p.PostalCode

				oaDetail.Domicile = &dat
			}
		}
		if oadata.CorrespondenceKey != nil {
			if p, ok := postalData[*oadata.CorrespondenceKey]; ok {
				var dat models.AddressDetail
				dat.PostalAddressKey = p.PostalAddressKey
				dat.AddressType = p.AddressType
				dat.AddressTypeName = p.AddressTypeName
				dat.ProvinsiKey = p.ProvinsiKey
				dat.ProvinsiName = p.ProvinsiName
				dat.KabupatenKey = p.KabupatenKey
				dat.KabupatenName = p.KabupatenName
				dat.KecamatanKey = p.KecamatanKey
				dat.KecamatanName = p.KecamatanName
				dat.AddressLine1 = p.AddressLine1
				dat.PostalCode = p.PostalCode

				oaDetail.Correspondence = &dat
			}
		}
	}

	//Nationality
	if oadata.Nationality != nil {
		var country models.MsCountry
		_, err = models.GetMsCountry(&country, strconv.FormatUint(*oadata.Nationality, 10))
		if err == nil {
			var dat models.MsCountryList
			dat.CountryKey = country.CountryKey
			dat.CouCode = country.CouCode
			dat.CouName = country.CouName
			oaDetail.Nationality = &dat
		}
	}

	//BRANCH
	if oareq.BranchKey != nil {
		var branch models.MsBranch
		status, err = models.GetMsBranch(&branch, strconv.FormatUint(*oareq.BranchKey, 10))
		if err == nil {
			var b models.MsBranchDropdown
			b.BranchKey = branch.BranchKey
			b.BranchName = branch.BranchName
			oaDetail.Branch = &b
		}

	}

	//AGENT
	if oareq.AgentKey != nil {
		var agent models.MsAgent
		status, err = models.GetMsAgent(&agent, strconv.FormatUint(*oareq.AgentKey, 10))
		if err == nil {
			var a models.MsAgentDropdown
			a.AgentKey = agent.AgentKey
			a.AgentName = agent.AgentName
			oaDetail.Agent = &a
		}

	}

	//BankRequest
	var accBank []models.OaRequestByField
	status, err = models.GetOaRequestBankByField(&accBank, "oa_request_key", strconv.FormatUint(oareq.OaRequestKey, 10))
	if err != nil {
		oaDetail.BankRequest = nil
	} else {
		oaDetail.BankRequest = &accBank
	}

	// RiskProfile
	var oaRiskProfile []models.AdminOaRiskProfile
	status, err = models.AdminGetOaRiskProfile(&oaRiskProfile, strconv.FormatUint(oareq.OaRequestKey, 10))
	if err == nil {
		if len(oaRiskProfile) > 0 {
			oaDetail.RiskProfile = &oaRiskProfile[0]
		}
	} else {
		oaDetail.RiskProfile = nil
	}

	// RiskProfileQuiz
	var oaRiskProfileQuiz []models.AdminOaRiskProfileQuiz
	status, err = models.AdminGetOaRiskProfileQuizByOaRequestKey(&oaRiskProfileQuiz, strconv.FormatUint(oareq.OaRequestKey, 10))
	if err != nil {
		if err != sql.ErrNoRows {
			log.Error(err.Error())
		}
	}
	if len(oaRiskProfileQuiz) > 0 {
		var questionIDs []string
		for _, quiz := range oaRiskProfileQuiz {
			if _, ok := lib.Find(questionIDs, strconv.FormatUint(quiz.QuizQuestionKey, 10)); !ok {
				questionIDs = append(questionIDs, strconv.FormatUint(quiz.QuizQuestionKey, 10))
			}
		}
		var optionDB []models.CmsQuizOptions
		status, err = models.GetCmsQuizOptionsIn(&optionDB, questionIDs, "quiz_question_key")
		if err != nil {
			if err != sql.ErrNoRows {
				log.Error(err.Error())
			}
		}

		optionData := make(map[uint64][]models.CmsQuizOptionsInfo)
		optionUserData := make(map[uint64]models.CmsQuizOptions)
		if len(optionDB) > 0 {
			for _, option := range optionDB {

				optionUserData[option.QuizOptionKey] = option

				var data models.CmsQuizOptionsInfo

				data.QuizOptionKey = option.QuizOptionKey
				if option.QuizOptionLabel != nil {
					data.QuizOptionLabel = *option.QuizOptionLabel
				}
				if option.QuizOptionTitle != nil {
					data.QuizOptionTitle = *option.QuizOptionTitle
				}
				if option.QuizOptionScore != nil {
					data.QuizOptionScore = *option.QuizOptionScore
				}
				if option.QuizOptionDefault != nil {
					data.QuizOptionDefault = *option.QuizOptionDefault
				}

				optionData[option.QuizQuestionKey] = append(optionData[option.QuizQuestionKey], data)
			}
		}

		var riskQuiz []models.RiskProfileQuiz

		for _, oaRiskQuiz := range oaRiskProfileQuiz {
			var risk models.RiskProfileQuiz

			risk.RiskProfileQuizKey = oaRiskQuiz.RiskProfileQuizKey
			if n, ok := optionUserData[oaRiskQuiz.QuizOptionKeyUser]; ok {
				risk.QuizOptionUser.QuizOptionKey = n.QuizOptionKey
				if n.QuizOptionLabel != nil {
					risk.QuizOptionUser.QuizOptionLabel = *n.QuizOptionLabel
				}
				if n.QuizOptionTitle != nil {
					risk.QuizOptionUser.QuizOptionTitle = *n.QuizOptionTitle
				}
				if n.QuizOptionScore != nil {
					risk.QuizOptionUser.QuizOptionScore = *n.QuizOptionScore
				}
				if n.QuizOptionDefault != nil {
					risk.QuizOptionUser.QuizOptionDefault = *n.QuizOptionDefault
				}
			}
			risk.QuizOptionScoreUser = oaRiskQuiz.QuizOptionScoreUser
			risk.QuizQuestionKey = oaRiskQuiz.QuizQuestionKey
			risk.HeaderQuizName = *oaRiskQuiz.HeaderQuizName
			risk.QuizTitle = oaRiskQuiz.QuizTitle

			if opt, ok := optionData[oaRiskQuiz.QuizQuestionKey]; ok {
				risk.Options = opt
			}

			riskQuiz = append(riskQuiz, risk)
		}
		oaDetail.RiskProfileQuiz = &riskQuiz
	}

	// InstitutionDocs
	var instDocs []models.OaInstitutionDocsDetail
	status, err = models.GetOaInstitutionDocsRequest(&instDocs, oaReqKey)
	if err == nil {
		var instDocsResp []models.OaInstitutionDocsDetail
		dir := config.BaseUrl + lib.INST_FILE_PATH + "/"
		for _, dt := range instDocs {
			var data models.OaInstitutionDocsDetail
			data.InstiDocsKey = dt.InstiDocsKey
			data.InstiDocumentType = dt.InstiDocumentType
			data.InstiDocumentTypeName = dt.InstiDocumentTypeName
			data.DocumentFileName = dt.DocumentFileName
			data.InstiDocumentName = dt.InstiDocumentName
			data.InstiDocumentRemarks = dt.InstiDocumentRemarks
			data.Path = dt.Path
			if dt.InstiDocsKey != nil && dt.DocumentFileName != nil {
				pathFile := dir + *dt.DocumentFileName
				data.Path = &pathFile
			}

			instDocsResp = append(instDocsResp, data)
		}
		oaDetail.InstitutionDocs = &instDocsResp
	}

	// InstitutionUserMaker
	var userMaker []models.OaInstitutionUserDetail
	status, err = models.GetOaInstitutionUserRequest(&userMaker, oaReqKey, lib.ROLE_INSTITUTION_MAKER)
	if err == nil {
		oaDetail.InstitutionUserMaker = &userMaker
	}
	// InstitutionUserChecker
	var userChecker []models.OaInstitutionUserDetail
	status, err = models.GetOaInstitutionUserRequest(&userChecker, oaReqKey, lib.ROLE_INSTITUTION_CHECKER)
	if err == nil {
		oaDetail.InstitutionUserChecker = &userChecker
	}
	// InstitutionUserReleaser
	var userReleaser []models.OaInstitutionUserDetail
	status, err = models.GetOaInstitutionUserRequest(&userReleaser, oaReqKey, lib.ROLE_INSTITUTION_RELEASER)
	if err == nil {
		oaDetail.InstitutionUserReleaser = &userReleaser
	}

	// InstitutionSharesHolder
	var sharedHolder []models.OaInstitutionSharesHolderDetail
	status, err = models.GetOaInstitutionSharesHolderRequest(&sharedHolder, oaReqKey)
	if err == nil {
		oaDetail.InstitutionSharesHolder = &sharedHolder
	}

	// oaDetail.InstitutionAuthPerson
	var authPerson []models.OaInstitutionAuthPersonDetail
	status, err = models.GetOaInstitutionAuthPersonRequest(&authPerson, oaReqKey)
	if err == nil {
		oaDetail.InstitutionAuthPerson = &authPerson
	}

	return oaDetail, nil
}

func SaveOaInstitutionData(c echo.Context) error {
	var err error
	var status int
	decimal.MarshalJSONWithoutQuotes = true
	paramsOaRequest := make(map[string]string)
	paramsInstitutionData := make(map[string]string)

	oaRequestKey := c.FormValue("oa_request_key")
	if oaRequestKey != "" {
		n, err := strconv.ParseUint(oaRequestKey, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: oa_request_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: oa_request_key", "Wrong input for parameter: oa_request_key")
		}

		if len(oaRequestKey) > 11 {
			log.Error("Wrong input for parameter: oa_request_key too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: oa_request_key too long, max 11 character", "Missing required parameter: oa_request_key too long, max 11 character")
		}
		paramsOaRequest["oa_request_key"] = oaRequestKey
	}

	if oaRequestKey != "" {
		var oareq models.OaRequest
		_, err = models.GetOaRequestInstitution(&oareq, oaRequestKey, "")
		if err != nil {
			log.Error("OA Request not found.")
			return lib.CustomError(http.StatusBadRequest, "OA Request not found.", "OA Request not found.")
		}

		var userCategory uint64
		userCategory = 3 //Branch
		if lib.Profile.UserCategoryKey == userCategory {
			if oareq.BranchKey != lib.Profile.BranchKey {
				log.Error("User not autorized.")
				return lib.CustomError(http.StatusBadRequest, "User not autorized.", "User not autorized.")
			}
		}

		if strconv.FormatUint(lib.Profile.UserID, 10) != *oareq.RecCreatedBy {
			log.Error("User not autorized.")
			return lib.CustomError(http.StatusBadRequest, "User not autorized.", "User not autorized.")
		}

		if *oareq.Oastatus != uint64(lib.DRAFT) {
			log.Error("User not autorized.")
			return lib.CustomError(http.StatusBadRequest, "User not autorized.", "User not autorized.")
		}
	}

	isSaveDraft := c.FormValue("is_save_draft")
	if isSaveDraft != "" {
		if isSaveDraft != "0" && isSaveDraft != "1" {
			log.Error("Wrong input for parameter: is_save_draft")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: is_save_draft", "Wrong input for parameter: is_save_draft")
		}
	} else {
		log.Error("Missing required parameter: is_save_draft")
		return lib.CustomError(http.StatusBadRequest, "is_save_draft can not be blank", "is_save_draft can not be blank")
	}

	branchKey := c.FormValue("branch_key")
	if branchKey == "" {
		log.Error("Missing required parameter: branch_key")
		return lib.CustomError(http.StatusBadRequest, "branch_key can not be blank", "branch_key can not be blank")
	} else {
		n, err := strconv.ParseUint(branchKey, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: branch_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: branch_key", "Wrong input for parameter: branch_key")
		}

		if len(branchKey) > 11 {
			log.Error("Wrong input for parameter: branch_key too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: branch_key too long, max 11 character", "Missing required parameter: branch_key too long, max 11 character")
		}
	}
	paramsOaRequest["branch_key"] = branchKey

	agentKey := c.FormValue("agent_key")
	if agentKey == "" {
		log.Error("Missing required parameter: agent_key")
		return lib.CustomError(http.StatusBadRequest, "agent_key can not be blank", "agent_key can not be blank")
	} else {
		n, err := strconv.ParseUint(agentKey, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: agent_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: agent_key", "Wrong input for parameter: agent_key")
		}

		if len(agentKey) > 11 {
			log.Error("Wrong input for parameter: agent_key too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: agent_key too long, max 11 character", "Missing required parameter: agent_key too long, max 11 character")
		}
	}
	paramsOaRequest["agent_key"] = agentKey

	nationality := c.FormValue("nationality")
	if nationality != "" {
		n, err := strconv.ParseUint(nationality, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: nationality")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: nationality", "Wrong input for parameter: nationality")
		}

		if len(nationality) > 11 {
			log.Error("Wrong input for parameter: nationality too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: nationality too long, max 11 character", "Missing required parameter: nationality too long, max 11 character")
		}
		paramsInstitutionData["nationality"] = nationality
	}

	fullName := c.FormValue("full_name")
	if fullName == "" {
		log.Error("Missing required parameter: full_name")
		return lib.CustomError(http.StatusBadRequest, "full_name can not be blank", "full_name can not be blank")
	} else {
		if len(fullName) > 150 {
			log.Error("Wrong input for parameter: full_name too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: full_name too long, max 150 character", "Missing required parameter: full_name too long, max 150 character")
		}
	}
	paramsInstitutionData["full_name"] = fullName

	shortName := c.FormValue("short_name")
	if shortName != "" {
		if len(shortName) > 150 {
			log.Error("Wrong input for parameter: short_name too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: short_name too long, max 50 character", "Missing required parameter: short_name too long, max 50 character")
		}
		paramsInstitutionData["short_name"] = shortName
	}

	tinNumber := c.FormValue("tin_number")
	if tinNumber != "" {
		if len(tinNumber) > 50 {
			log.Error("Wrong input for parameter: tin_number too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: tin_number too long, max 50 character", "Missing required parameter: tin_number too long, max 50 character")
		}
		paramsInstitutionData["tin_number"] = tinNumber
	}

	establishedCity := c.FormValue("established_city")
	if establishedCity != "" {
		if len(establishedCity) > 100 {
			log.Error("Wrong input for parameter: established_city too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: established_city too long, max 100 character", "Missing required parameter: established_city too long, max 100 character")
		}
		paramsInstitutionData["established_city"] = establishedCity
	}

	layout := "2006-01-02 15:04:05"

	establishedDate := c.FormValue("established_date")
	if establishedDate != "" {
		establishedDate += " 00:00:00"
		date, _ := time.Parse(layout, establishedDate)
		dateStr := date.Format(layout)
		paramsInstitutionData["established_date"] = dateStr
	}

	deedNo := c.FormValue("deed_no")
	if deedNo != "" {
		if len(deedNo) > 100 {
			log.Error("Wrong input for parameter: deed_no too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: deed_no too long, max 100 character", "Missing required parameter: deed_no too long, max 100 character")
		}
		paramsInstitutionData["deed_no"] = deedNo
	}

	deedDate := c.FormValue("deed_date")
	if deedDate != "" {
		deedDate += " 00:00:00"
		date, _ := time.Parse(layout, deedDate)
		dateStr := date.Format(layout)
		paramsInstitutionData["deed_date"] = dateStr
	}

	lsEstablishValidationNo := c.FormValue("ls_establish_validation_no")
	if lsEstablishValidationNo != "" {
		if len(lsEstablishValidationNo) > 100 {
			log.Error("Wrong input for parameter: ls_establish_validation_no too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: ls_establish_validation_no too long, max 100 character", "Missing required parameter: ls_establish_validation_no too long, max 100 character")
		}
		paramsInstitutionData["ls_establish_validation_no"] = lsEstablishValidationNo
	}

	lsEstablishValidationDate := c.FormValue("ls_establish_validation_date")
	if lsEstablishValidationDate != "" {
		lsEstablishValidationDate += " 00:00:00"
		date, _ := time.Parse(layout, lsEstablishValidationDate)
		dateStr := date.Format(layout)
		paramsInstitutionData["ls_establish_validation_date"] = dateStr
	}

	lastChangeAaNo := c.FormValue("last_change_aa_no")
	if lastChangeAaNo != "" {
		if len(lastChangeAaNo) > 100 {
			log.Error("Wrong input for parameter: last_change_aa_no too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: last_change_aa_no too long, max 100 character", "Missing required parameter: last_change_aa_no too long, max 100 character")
		}
		paramsInstitutionData["last_change_aa_no"] = lastChangeAaNo
	}

	lastChangeAaDate := c.FormValue("last_change_aa_date")
	if lastChangeAaDate != "" {
		lastChangeAaDate += " 00:00:00"
		date, _ := time.Parse(layout, lastChangeAaDate)
		dateStr := date.Format(layout)
		paramsInstitutionData["last_change_aa_date"] = dateStr
	}

	lsLastChangeAaNo := c.FormValue("ls_last_change_aa_no")
	if lsLastChangeAaNo != "" {
		if len(lsLastChangeAaNo) > 100 {
			log.Error("Wrong input for parameter: ls_last_change_aa_no too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: ls_last_change_aa_no too long, max 100 character", "Missing required parameter: ls_last_change_aa_no too long, max 100 character")
		}
		paramsInstitutionData["ls_last_change_aa_no"] = lsLastChangeAaNo
	}
	lsLastChangeAaDate := c.FormValue("ls_last_change_aa_date")
	if lsLastChangeAaDate != "" {
		lsLastChangeAaDate += " 00:00:00"
		date, _ := time.Parse(layout, lsLastChangeAaDate)
		dateStr := date.Format(layout)
		paramsInstitutionData["ls_last_change_aa_date"] = dateStr
	}

	managementDeedNo := c.FormValue("management_deed_no")
	if managementDeedNo != "" {
		if len(managementDeedNo) > 100 {
			log.Error("Wrong input for parameter: management_deed_no too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: management_deed_no too long, max 100 character", "Missing required parameter: management_deed_no too long, max 100 character")
		}
		paramsInstitutionData["management_deed_no"] = managementDeedNo
	}
	managementDeedDate := c.FormValue("management_deed_date")
	if managementDeedDate != "" {
		managementDeedDate += " 00:00:00"
		date, _ := time.Parse(layout, managementDeedDate)
		dateStr := date.Format(layout)
		paramsInstitutionData["management_deed_date"] = dateStr
	}

	lsMgtChangeDeedNo := c.FormValue("ls_mgt_change_deed_no")
	if lsMgtChangeDeedNo != "" {
		if len(lsMgtChangeDeedNo) > 100 {
			log.Error("Wrong input for parameter: ls_mgt_change_deed_no too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: ls_mgt_change_deed_no too long, max 100 character", "Missing required parameter: ls_mgt_change_deed_no too long, max 100 character")
		}
		paramsInstitutionData["ls_mgt_change_deed_no"] = lsMgtChangeDeedNo
	}
	lsMgtChangeDeedDate := c.FormValue("ls_mgt_change_deed_date")
	if lsMgtChangeDeedDate != "" {
		lsMgtChangeDeedDate += " 00:00:00"
		date, _ := time.Parse(layout, lsMgtChangeDeedDate)
		dateStr := date.Format(layout)
		paramsInstitutionData["ls_mgt_change_deed_date"] = dateStr
	}

	skdLicenseNo := c.FormValue("skd_license_no")
	if skdLicenseNo != "" {
		if len(skdLicenseNo) > 100 {
			log.Error("Wrong input for parameter: skd_license_no too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: skd_license_no too long, max 100 character", "Missing required parameter: skd_license_no too long, max 100 character")
		}
		paramsInstitutionData["skd_license_no"] = skdLicenseNo
	}
	skdLicenseDate := c.FormValue("skd_license_date")
	if skdLicenseDate != "" {
		skdLicenseDate += " 00:00:00"
		date, _ := time.Parse(layout, skdLicenseDate)
		dateStr := date.Format(layout)
		paramsInstitutionData["skd_license_date"] = dateStr
	}

	bizLicenseNo := c.FormValue("biz_license_no")
	if bizLicenseNo != "" {
		if len(bizLicenseNo) > 100 {
			log.Error("Wrong input for parameter: biz_license_no too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: biz_license_no too long, max 100 character", "Missing required parameter: biz_license_no too long, max 100 character")
		}
		paramsInstitutionData["biz_license_no"] = bizLicenseNo
	}
	bizLicenseDate := c.FormValue("biz_license_date")
	if bizLicenseDate != "" {
		bizLicenseDate += " 00:00:00"
		date, _ := time.Parse(layout, bizLicenseDate)
		dateStr := date.Format(layout)
		paramsInstitutionData["biz_license_date"] = dateStr
	}

	nibNo := c.FormValue("nib_no")
	if nibNo != "" {
		if len(nibNo) > 100 {
			log.Error("Wrong input for parameter: nib_no too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: nib_no too long, max 100 character", "Missing required parameter: nib_no too long, max 100 character")
		}
		paramsInstitutionData["nib_no"] = nibNo
	}
	nibDate := c.FormValue("nib_date")
	if nibDate != "" {
		nibDate += " 00:00:00"
		date, _ := time.Parse(layout, nibDate)
		dateStr := date.Format(layout)
		paramsInstitutionData["nib_date"] = dateStr
	}

	phoneNo := c.FormValue("phone_no")
	if phoneNo != "" {
		if len(phoneNo) > 50 {
			log.Error("Wrong input for parameter: phone_no too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: phone_no too long, max 50 character", "Missing required parameter: phone_no too long, max 50 character")
		}
		paramsInstitutionData["phone_no"] = phoneNo
	}
	mobileNo := c.FormValue("mobile_no")
	if mobileNo != "" {
		if len(mobileNo) > 50 {
			log.Error("Wrong input for parameter: mobile_no too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: mobile_no too long, max 50 character", "Missing required parameter: mobile_no too long, max 50 character")
		}
		paramsInstitutionData["mobile_no"] = mobileNo
	}
	faxNo := c.FormValue("fax_no")
	if faxNo != "" {
		if len(faxNo) > 50 {
			log.Error("Wrong input for parameter: fax_no too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: fax_no too long, max 50 character", "Missing required parameter: fax_no too long, max 50 character")
		}
		paramsInstitutionData["fax_no"] = faxNo
	}
	emailAddress := c.FormValue("email_address")
	if emailAddress != "" {
		if len(emailAddress) > 100 {
			log.Error("Wrong input for parameter: email_address too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: email_address too long, max 100 character", "Missing required parameter: email_address too long, max 100 character")
		}
		if !lib.IsValidEmail(emailAddress) {
			log.Error("Wrong input for parameter: email_address wrong format email")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: email_address wrong format email", "Wrong input for parameter: email_address wrong format email")
		}
		paramsInstitutionData["email_address"] = emailAddress
	}

	paramsPostalDomicile := make(map[string]string)
	//DOMICILE
	domicilePostalAddressKey := c.FormValue("domicile_postal_address_key")
	if domicilePostalAddressKey != "" {
		n, err := strconv.ParseUint(domicilePostalAddressKey, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: domicile_postal_address_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: domicile_postal_address_key", "Wrong input for parameter: domicile_postal_address_key")
		}
		if len(domicilePostalAddressKey) > 11 {
			log.Error("Wrong input for parameter: domicile_postal_address_key too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: domicile_postal_address_key too long, max 11 character", "Missing required parameter: domicile_postal_address_key too long, max 11 character")
		}
		paramsPostalDomicile["postal_address_key"] = domicilePostalAddressKey
	}
	domicileAddress := c.FormValue("domicile_address")
	if domicileAddress != "" {
		if len(domicileAddress) > 255 {
			log.Error("Wrong input for parameter: domicile_address too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: domicile_address too long, max 255 character", "Missing required parameter: domicile_address too long, max 255 character")
		}
		paramsPostalDomicile["address_line1"] = domicileAddress
	}
	domicileSubdistric := c.FormValue("domicile_subdistric")
	if domicileSubdistric != "" {
		n, err := strconv.ParseUint(domicileSubdistric, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: domicile_subdistric")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: domicile_subdistric", "Wrong input for parameter: domicile_subdistric")
		}
		if len(domicileSubdistric) > 11 {
			log.Error("Wrong input for parameter: domicile_subdistric too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: domicile_subdistric too long, max 11 character", "Missing required parameter: domicile_subdistric too long, max 11 character")
		}
		paramsPostalDomicile["kecamatan_key"] = domicileSubdistric
	}
	domicileCity := c.FormValue("domicile_city")
	if domicileCity != "" {
		n, err := strconv.ParseUint(domicileCity, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: domicile_city")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: domicile_city", "Wrong input for parameter: domicile_city")
		}
		if len(domicileCity) > 11 {
			log.Error("Wrong input for parameter: domicile_city too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: domicile_city too long, max 11 character", "Missing required parameter: domicile_city too long, max 11 character")
		}
		paramsPostalDomicile["kabupaten_key"] = domicileCity
	}
	domicileProvince := c.FormValue("domicile_province")
	log.Println(domicileProvince)

	domicilePostalcode := c.FormValue("domicile_postalcode")
	if domicilePostalcode != "" {
		if len(domicilePostalcode) > 10 {
			log.Error("Wrong input for parameter: domicile_postalcode too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: domicile_postalcode too long, max 10 character", "Missing required parameter: domicile_postalcode too long, max 10 character")
		}
		paramsPostalDomicile["postal_code"] = domicilePostalcode
	}

	paramsPostalCorrespondence := make(map[string]string)
	//CORRESPONDENCE
	correspondencePostalAddressKey := c.FormValue("correspondence_postal_address_key")
	if correspondencePostalAddressKey != "" {
		n, err := strconv.ParseUint(correspondencePostalAddressKey, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: correspondence_postal_address_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: correspondence_postal_address_key", "Wrong input for parameter: correspondence_postal_address_key")
		}
		if len(correspondencePostalAddressKey) > 11 {
			log.Error("Wrong input for parameter: correspondence_postal_address_key too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: correspondence_postal_address_key too long, max 11 character", "Missing required parameter: correspondence_postal_address_key too long, max 11 character")
		}
		paramsPostalCorrespondence["postal_address_key"] = correspondencePostalAddressKey
	}
	correspondenceAddress := c.FormValue("correspondence_address")
	if correspondenceAddress != "" {
		if len(correspondenceAddress) > 255 {
			log.Error("Wrong input for parameter: correspondence_address too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: correspondence_address too long, max 255 character", "Missing required parameter: correspondence_address too long, max 255 character")
		}
		paramsPostalCorrespondence["address_line1"] = correspondenceAddress
	}
	correspondenceSubdistric := c.FormValue("correspondence_subdistric")
	if correspondenceSubdistric != "" {
		n, err := strconv.ParseUint(correspondenceSubdistric, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: correspondence_subdistric")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: correspondence_subdistric", "Wrong input for parameter: correspondence_subdistric")
		}
		if len(correspondenceSubdistric) > 11 {
			log.Error("Wrong input for parameter: correspondence_subdistric too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: correspondence_subdistric too long, max 11 character", "Missing required parameter: correspondence_subdistric too long, max 11 character")
		}
		paramsPostalCorrespondence["kecamatan_key"] = correspondenceSubdistric
	}
	correspondenceCity := c.FormValue("correspondence_city")
	if correspondenceCity != "" {
		n, err := strconv.ParseUint(correspondenceCity, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: correspondence_city")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: correspondence_city", "Wrong input for parameter: correspondence_city")
		}
		if len(correspondenceCity) > 11 {
			log.Error("Wrong input for parameter: correspondence_city too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: correspondence_city too long, max 11 character", "Missing required parameter: correspondence_city too long, max 11 character")
		}
		paramsPostalCorrespondence["kabupaten_key"] = correspondenceCity
	}
	correspondenceProvince := c.FormValue("correspondence_province")
	log.Println(correspondenceProvince)
	correspondencePostalcode := c.FormValue("correspondence_postalcode")
	if correspondencePostalcode != "" {
		if len(correspondencePostalcode) > 10 {
			log.Error("Wrong input for parameter: correspondence_postalcode too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: correspondence_postalcode too long, max 10 character", "Missing required parameter: correspondence_postalcode too long, max 10 character")
		}
		paramsPostalCorrespondence["postal_code"] = correspondencePostalcode
	}

	intitutionType := c.FormValue("intitution_type")
	if intitutionType != "" {
		n, err := strconv.ParseUint(intitutionType, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: intitution_type")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: intitution_type", "Wrong input for parameter: intitution_type")
		}
		if len(intitutionType) > 11 {
			log.Error("Wrong input for parameter: intitution_type too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: intitution_type too long, max 11 character", "Missing required parameter: intitution_type too long, max 11 character")
		}
		paramsInstitutionData["intitution_type"] = intitutionType
	}
	intitutionClassification := c.FormValue("intitution_classification")
	if intitutionClassification != "" {
		n, err := strconv.ParseUint(intitutionClassification, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: intitution_classification")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: intitution_classification", "Wrong input for parameter: intitution_classification")
		}
		if len(intitutionClassification) > 11 {
			log.Error("Wrong input for parameter: intitution_classification too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: intitution_classification too long, max 11 character", "Missing required parameter: intitution_classification too long, max 11 character")
		}
		paramsInstitutionData["intitution_classification"] = intitutionClassification
	}
	intitutionCharacteristic := c.FormValue("intitution_characteristic")
	if intitutionCharacteristic != "" {
		n, err := strconv.ParseUint(intitutionCharacteristic, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: intitution_characteristic")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: intitution_characteristic", "Wrong input for parameter: intitution_characteristic")
		}
		if len(intitutionCharacteristic) > 11 {
			log.Error("Wrong input for parameter: intitution_characteristic too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: intitution_characteristic too long, max 11 character", "Missing required parameter: intitution_characteristic too long, max 11 character")
		}
		paramsInstitutionData["intitution_characteristic"] = intitutionCharacteristic
	}
	intitutionBusinessType := c.FormValue("intitution_business_type")
	if intitutionBusinessType != "" {
		n, err := strconv.ParseUint(intitutionBusinessType, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: intitution_business_type")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: intitution_business_type", "Wrong input for parameter: intitution_business_type")
		}
		if len(intitutionBusinessType) > 11 {
			log.Error("Wrong input for parameter: intitution_business_type too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: intitution_business_type too long, max 11 character", "Missing required parameter: intitution_business_type too long, max 11 character")
		}
		paramsInstitutionData["intitution_business_type"] = intitutionBusinessType
	}

	instiAnnuallyIncome := c.FormValue("insti_annually_income")
	if instiAnnuallyIncome != "" {
		n, err := strconv.ParseUint(instiAnnuallyIncome, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: insti_annually_income")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: insti_annually_income", "Wrong input for parameter: insti_annually_income")
		}
		if len(instiAnnuallyIncome) > 11 {
			log.Error("Wrong input for parameter: insti_annually_income too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: insti_annually_income too long, max 11 character", "Missing required parameter: insti_annually_income too long, max 11 character")
		}
		paramsInstitutionData["insti_annually_income"] = instiAnnuallyIncome
	}
	instiSourceOfIncome := c.FormValue("insti_source_of_income")
	if instiSourceOfIncome != "" {
		n, err := strconv.ParseUint(instiSourceOfIncome, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: insti_source_of_income")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: insti_source_of_income", "Wrong input for parameter: insti_source_of_income")
		}
		if len(instiSourceOfIncome) > 11 {
			log.Error("Wrong input for parameter: insti_source_of_income too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: insti_source_of_income too long, max 11 character", "Missing required parameter: insti_source_of_income too long, max 11 character")
		}
		paramsInstitutionData["insti_source_of_income"] = instiSourceOfIncome
	}
	instiInvestmentPurpose := c.FormValue("insti_investment_purpose")
	if instiInvestmentPurpose != "" {
		n, err := strconv.ParseUint(instiInvestmentPurpose, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: insti_investment_purpose")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: insti_investment_purpose", "Wrong input for parameter: insti_investment_purpose")
		}
		if len(instiInvestmentPurpose) > 11 {
			log.Error("Wrong input for parameter: insti_investment_purpose too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: insti_investment_purpose too long, max 11 character", "Missing required parameter: insti_investment_purpose too long, max 11 character")
		}
		paramsInstitutionData["insti_investment_purpose"] = instiInvestmentPurpose
	}

	//BENEFICIAL OWNER
	boName := c.FormValue("bo_name")
	if boName != "" {
		if len(boName) > 50 {
			log.Error("Wrong input for parameter: bo_name too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: bo_name too long, max 50 character", "Missing required parameter: bo_name too long, max 50 character")
		}
		paramsInstitutionData["bo_name"] = boName
	}
	boIdnumber := c.FormValue("bo_idnumber")
	if boIdnumber != "" {
		if len(boIdnumber) > 30 {
			log.Error("Wrong input for parameter: bo_idnumber too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: bo_idnumber too long, max 30 character", "Missing required parameter: bo_idnumber too long, max 30 character")
		}
		paramsInstitutionData["bo_idnumber"] = boIdnumber
	}
	boBusiness := c.FormValue("bo_business")
	if boBusiness != "" {
		if len(boBusiness) > 30 {
			log.Error("Wrong input for parameter: bo_business too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: bo_business too long, max 30 character", "Missing required parameter: bo_business too long, max 30 character")
		}
		paramsInstitutionData["bo_business"] = boBusiness
	}
	boIdaddress := c.FormValue("bo_idaddress")
	if boIdaddress != "" {
		if len(boIdaddress) > 150 {
			log.Error("Wrong input for parameter: bo_idaddress too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: bo_idaddress too long, max 150 character", "Missing required parameter: bo_idaddress too long, max 150 character")
		}
		paramsInstitutionData["bo_idaddress"] = boIdaddress
	}
	boBusinessAddress := c.FormValue("bo_business_address")
	if boBusinessAddress != "" {
		if len(boBusinessAddress) > 150 {
			log.Error("Wrong input for parameter: bo_business_address too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: bo_business_address too long, max 150 character", "Missing required parameter: bo_business_address too long, max 150 character")
		}
		paramsInstitutionData["bo_business_address"] = boBusinessAddress
	}
	boAnnuallyIncome := c.FormValue("bo_annually_income")
	if boAnnuallyIncome != "" {
		_, err := decimal.NewFromString(boAnnuallyIncome)
		if err != nil {
			log.Error("Wrong input for parameter: bo_annually_income")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: bo_annually_income", "Wrong input for parameter: bo_annually_income")
		}
		paramsInstitutionData["bo_annually_income"] = boAnnuallyIncome
	}
	boRelation := c.FormValue("bo_relation")
	if boRelation != "" {
		n, err := strconv.ParseUint(boRelation, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: bo_relation")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: bo_relation", "Wrong input for parameter: bo_relation")
		}
		if len(boRelation) > 11 {
			log.Error("Wrong input for parameter: bo_relation too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: bo_relation too long, max 11 character", "Missing required parameter: bo_relation too long, max 11 character")
		}
		paramsInstitutionData["bo_relation"] = boRelation
	}

	//SHARES HOLDER //ARRAY
	sharesHolder := c.FormValue("shares_holder")

	saveSharesHolder := false
	var sharesHolderData []interface{}
	if isSaveDraft == "1" {
		if sharesHolder != "" {
			saveSharesHolder = true
			var sharesHolderSlice []interface{}
			err = json.Unmarshal([]byte(sharesHolder), &sharesHolderSlice)
			if err != nil {
				log.Error(err.Error())
				log.Error("Missing required parameter: shares_holder")
				return lib.CustomError(http.StatusBadRequest, err.Error(), "Wrong input for parameter: shares_holder")
			}
			key := 1
			if len(sharesHolderSlice) > 0 {
				for _, val := range sharesHolderSlice {
					sh := make(map[string]interface{})
					valueMap := val.(map[string]interface{})
					if val, ok := valueMap["inst_shares_holder_key"]; ok {
						if val.(string) != "" {
							n, err := strconv.ParseUint(val.(string), 10, 64)
							if err != nil || n == 0 {
								log.Error("Wrong input parameter: shares_holder key : " + strconv.FormatUint(uint64(key), 10) + " inst_shares_holder_key.")
								return lib.CustomError(http.StatusBadRequest,
									"Wrong input parameter: shares_holder key : "+strconv.FormatUint(uint64(key), 10)+" inst_shares_holder_key.",
									"Wrong input parameter: shares_holder key : "+strconv.FormatUint(uint64(key), 10)+" inst_shares_holder_key.")
							}
						}
						sh["inst_shares_holder_key"] = val.(string)
					} else {
						log.Error("Missing required parameter: shares_holder key : " + strconv.FormatUint(uint64(key), 10) + " inst_shares_holder_key tidak ditemukan")
						return lib.CustomError(http.StatusBadRequest,
							"Missing required parameter: shares_holder key : "+strconv.FormatUint(uint64(key), 10)+" inst_shares_holder_key tidak ditemukan",
							"Missing required parameter: shares_holder key : "+strconv.FormatUint(uint64(key), 10)+" inst_shares_holder_key tidak ditemukan")
					}

					if val, ok := valueMap["nationality"]; ok {
						n, err := strconv.ParseUint(val.(string), 10, 64)
						if err != nil || n == 0 {
							log.Error("Wrong input parameter: shares_holder key : " + strconv.FormatUint(uint64(key), 10) + " nationality.")
							return lib.CustomError(http.StatusBadRequest,
								"Wrong input parameter: shares_holder key : "+strconv.FormatUint(uint64(key), 10)+" nationality.",
								"Wrong input parameter: shares_holder key : "+strconv.FormatUint(uint64(key), 10)+" nationality.")
						}
						sh["nationality"] = val.(string)
					} else {
						log.Error("Missing required parameter: shares_holder key : " + strconv.FormatUint(uint64(key), 10) + " nationality tidak ditemukan")
						return lib.CustomError(http.StatusBadRequest,
							"Missing required parameter: shares_holder key : "+strconv.FormatUint(uint64(key), 10)+" nationality tidak ditemukan",
							"Missing required parameter: shares_holder key : "+strconv.FormatUint(uint64(key), 10)+" nationality tidak ditemukan")
					}

					if val, ok := valueMap["holder_full_name"]; ok {
						if len(val.(string)) > 150 {
							log.Error("Missing required parameter: shares_holder key : " + strconv.FormatUint(uint64(key), 10) + " holder_full_name too long, max 150 character.")
							return lib.CustomError(http.StatusBadRequest,
								"Missing required parameter: shares_holder key : "+strconv.FormatUint(uint64(key), 10)+" holder_full_name too long, max 150 character.",
								"Missing required parameter: shares_holder key : "+strconv.FormatUint(uint64(key), 10)+" holder_full_name too long, max 150 character.")
						}
						sh["holder_full_name"] = val.(string)
					} else {
						log.Error("Missing required parameter: shares_holder key : " + strconv.FormatUint(uint64(key), 10) + " holder_full_name tidak ditemukan")
						return lib.CustomError(http.StatusBadRequest,
							"Missing required parameter: shares_holder key : "+strconv.FormatUint(uint64(key), 10)+" holder_full_name tidak ditemukan",
							"Missing required parameter: shares_holder key : "+strconv.FormatUint(uint64(key), 10)+" holder_full_name tidak ditemukan")
					}

					if val, ok := valueMap["idcard_type"]; ok {
						n, err := strconv.ParseUint(val.(string), 10, 64)
						if err != nil || n == 0 {
							log.Error("Wrong input parameter: shares_holder key : " + strconv.FormatUint(uint64(key), 10) + " idcard_type.")
							return lib.CustomError(http.StatusBadRequest,
								"Wrong input parameter: shares_holder key : "+strconv.FormatUint(uint64(key), 10)+" idcard_type.",
								"Wrong input parameter: shares_holder key : "+strconv.FormatUint(uint64(key), 10)+" idcard_type.")
						}
						sh["idcard_type"] = val.(string)
					} else {
						log.Error("Missing required parameter: shares_holder key : " + strconv.FormatUint(uint64(key), 10) + " idcard_type tidak ditemukan")
						return lib.CustomError(http.StatusBadRequest,
							"Missing required parameter: shares_holder key : "+strconv.FormatUint(uint64(key), 10)+" idcard_type tidak ditemukan",
							"Missing required parameter: shares_holder key : "+strconv.FormatUint(uint64(key), 10)+" idcard_type tidak ditemukan")
					}

					if val, ok := valueMap["idcard_no"]; ok {
						if len(val.(string)) > 20 {
							log.Error("Missing required parameter: shares_holder key : " + strconv.FormatUint(uint64(key), 10) + " idcard_no too long, max 20 character.")
							return lib.CustomError(http.StatusBadRequest,
								"Missing required parameter: shares_holder key : "+strconv.FormatUint(uint64(key), 10)+" idcard_no too long, max 20 character.",
								"Missing required parameter: shares_holder key : "+strconv.FormatUint(uint64(key), 10)+" idcard_no too long, max 20 character.")
						}
						sh["idcard_no"] = val.(string)
					} else {
						log.Error("Missing required parameter: shares_holder key : " + strconv.FormatUint(uint64(key), 10) + " idcard_no tidak ditemukan")
						return lib.CustomError(http.StatusBadRequest,
							"Missing required parameter: shares_holder key : "+strconv.FormatUint(uint64(key), 10)+" idcard_no tidak ditemukan",
							"Missing required parameter: shares_holder key : "+strconv.FormatUint(uint64(key), 10)+" idcard_no tidak ditemukan")
					}

					if val, ok := valueMap["holder_dob"]; ok {
						if val.(string) != "" {
							holderDob := val.(string) + " 00:00:00"
							date, _ := time.Parse(layout, holderDob)
							dateStr := date.Format(layout)
							sh["holder_dob"] = dateStr
						} else {
							sh["holder_dob"] = ""
						}
					} else {
						log.Error("Missing required parameter: shares_holder key : " + strconv.FormatUint(uint64(key), 10) + " holder_dob tidak ditemukan")
						return lib.CustomError(http.StatusBadRequest,
							"Missing required parameter: shares_holder key : "+strconv.FormatUint(uint64(key), 10)+" holder_dob tidak ditemukan",
							"Missing required parameter: shares_holder key : "+strconv.FormatUint(uint64(key), 10)+" holder_dob tidak ditemukan")
					}

					if val, ok := valueMap["shares_percent"]; ok {
						_, err := decimal.NewFromString(val.(string))
						if err != nil {
							log.Error("Wrong input parameter: shares_holder key : " + strconv.FormatUint(uint64(key), 10) + " shares_percent")
							return lib.CustomError(http.StatusBadRequest,
								"Wrong input parameter: shares_holder key : "+strconv.FormatUint(uint64(key), 10)+" shares_percent",
								"Wrong input parameter: shares_holder key : "+strconv.FormatUint(uint64(key), 10)+" shares_percent")
						}
						sh["shares_percent"] = val.(string)
					} else {
						log.Error("Missing required parameter: shares_holder key : " + strconv.FormatUint(uint64(key), 10) + " shares_percent tidak ditemukan")
						return lib.CustomError(http.StatusBadRequest,
							"Missing required parameter: shares_holder key : "+strconv.FormatUint(uint64(key), 10)+" shares_percent tidak ditemukan",
							"Missing required parameter: shares_holder key : "+strconv.FormatUint(uint64(key), 10)+" shares_percent tidak ditemukan")
					}

					sharesHolderData = append(sharesHolderData, sh)
					key++
				}
			}
		}
	}

	//AUTH PERSON //ARRAY
	authPerson := c.FormValue("auth_person")

	saveAuthPerson := false
	var authPersonData []interface{}
	if isSaveDraft == "1" {
		if authPerson != "" {
			saveAuthPerson = true
			var authPersonSlice []interface{}
			err = json.Unmarshal([]byte(authPerson), &authPersonSlice)
			if err != nil {
				log.Error(err.Error())
				log.Error("Missing required parameter: auth_person")
				return lib.CustomError(http.StatusBadRequest, err.Error(), "Wrong input for parameter: auth_person")
			}
			key := 1
			if len(authPersonSlice) > 0 {
				for _, val := range authPersonSlice {
					ap := make(map[string]interface{})
					valueMap := val.(map[string]interface{})
					if val, ok := valueMap["insti_auth_person_key"]; ok {
						if val.(string) != "" {
							n, err := strconv.ParseUint(val.(string), 10, 64)
							if err != nil || n == 0 {
								log.Error("Wrong input parameter: shares_holder key : " + strconv.FormatUint(uint64(key), 10) + " insti_auth_person_key.")
								return lib.CustomError(http.StatusBadRequest,
									"Wrong input parameter: shares_holder key : "+strconv.FormatUint(uint64(key), 10)+" insti_auth_person_key.",
									"Wrong input parameter: shares_holder key : "+strconv.FormatUint(uint64(key), 10)+" insti_auth_person_key.")
							}
						}
						ap["insti_auth_person_key"] = val.(string)
					} else {
						log.Error("Missing required parameter: auth_person key : " + strconv.FormatUint(uint64(key), 10) + " insti_auth_person_key tidak ditemukan")
						return lib.CustomError(http.StatusBadRequest,
							"Missing required parameter: auth_person key : "+strconv.FormatUint(uint64(key), 10)+" insti_auth_person_key tidak ditemukan",
							"Missing required parameter: auth_person key : "+strconv.FormatUint(uint64(key), 10)+" insti_auth_person_key tidak ditemukan")
					}

					if val, ok := valueMap["full_name"]; ok {
						if len(val.(string)) > 100 {
							log.Error("Missing required parameter: auth_person key : " + strconv.FormatUint(uint64(key), 10) + " full_name too long, max 100 character.")
							return lib.CustomError(http.StatusBadRequest,
								"Missing required parameter: auth_person key : "+strconv.FormatUint(uint64(key), 10)+" full_name too long, max 100 character.",
								"Missing required parameter: auth_person key : "+strconv.FormatUint(uint64(key), 10)+" full_name too long, max 100 character.")
						}
						ap["full_name"] = val.(string)
					} else {
						log.Error("Missing required parameter: auth_person key : " + strconv.FormatUint(uint64(key), 10) + " full_name tidak ditemukan")
						return lib.CustomError(http.StatusBadRequest,
							"Missing required parameter: auth_person key : "+strconv.FormatUint(uint64(key), 10)+" full_name tidak ditemukan",
							"Missing required parameter: auth_person key : "+strconv.FormatUint(uint64(key), 10)+" full_name tidak ditemukan")
					}

					if val, ok := valueMap["idcard_type"]; ok {
						n, err := strconv.ParseUint(val.(string), 10, 64)
						if err != nil || n == 0 {
							log.Error("Wrong input parameter: auth_person key : " + strconv.FormatUint(uint64(key), 10) + " idcard_type.")
							return lib.CustomError(http.StatusBadRequest,
								"Wrong input parameter: auth_person key : "+strconv.FormatUint(uint64(key), 10)+" idcard_type.",
								"Wrong input parameter: auth_person key : "+strconv.FormatUint(uint64(key), 10)+" idcard_type.")
						}
						ap["idcard_type"] = val.(string)
					} else {
						log.Error("Missing required parameter: auth_person key : " + strconv.FormatUint(uint64(key), 10) + " idcard_type tidak ditemukan")
						return lib.CustomError(http.StatusBadRequest,
							"Missing required parameter: auth_person key : "+strconv.FormatUint(uint64(key), 10)+" idcard_type tidak ditemukan",
							"Missing required parameter: auth_person key : "+strconv.FormatUint(uint64(key), 10)+" idcard_type tidak ditemukan")
					}

					if val, ok := valueMap["idcard_no"]; ok {
						if len(val.(string)) > 20 {
							log.Error("Missing required parameter: auth_person key : " + strconv.FormatUint(uint64(key), 10) + " idcard_no too long, max 20 character.")
							return lib.CustomError(http.StatusBadRequest,
								"Missing required parameter: auth_person key : "+strconv.FormatUint(uint64(key), 10)+" idcard_no too long, max 20 character.",
								"Missing required parameter: auth_person key : "+strconv.FormatUint(uint64(key), 10)+" idcard_no too long, max 20 character.")
						}
						ap["idcard_no"] = val.(string)
					} else {
						log.Error("Missing required parameter: auth_person key : " + strconv.FormatUint(uint64(key), 10) + " idcard_no tidak ditemukan")
						return lib.CustomError(http.StatusBadRequest,
							"Missing required parameter: auth_person key : "+strconv.FormatUint(uint64(key), 10)+" idcard_no tidak ditemukan",
							"Missing required parameter: auth_person key : "+strconv.FormatUint(uint64(key), 10)+" idcard_no tidak ditemukan")
					}

					if val, ok := valueMap["person_dob"]; ok {
						if val.(string) != "" {
							holderDob := val.(string) + " 00:00:00"
							date, _ := time.Parse(layout, holderDob)
							dateStr := date.Format(layout)
							ap["person_dob"] = dateStr
						} else {
							ap["person_dob"] = ""
						}
					} else {
						log.Error("Missing required parameter: auth_person key : " + strconv.FormatUint(uint64(key), 10) + " person_dob tidak ditemukan")
						return lib.CustomError(http.StatusBadRequest,
							"Missing required parameter: auth_person key : "+strconv.FormatUint(uint64(key), 10)+" person_dob tidak ditemukan",
							"Missing required parameter: auth_person key : "+strconv.FormatUint(uint64(key), 10)+" person_dob tidak ditemukan")
					}

					if val, ok := valueMap["nationality"]; ok {
						n, err := strconv.ParseUint(val.(string), 10, 64)
						if err != nil || n == 0 {
							log.Error("Wrong input parameter: auth_person key : " + strconv.FormatUint(uint64(key), 10) + " nationality.")
							return lib.CustomError(http.StatusBadRequest,
								"Wrong input parameter: auth_person key : "+strconv.FormatUint(uint64(key), 10)+" nationality.",
								"Wrong input parameter: auth_person key : "+strconv.FormatUint(uint64(key), 10)+" nationality.")
						}
						ap["nationality"] = val.(string)
					} else {
						log.Error("Missing required parameter: auth_person key : " + strconv.FormatUint(uint64(key), 10) + " nationality tidak ditemukan")
						return lib.CustomError(http.StatusBadRequest,
							"Missing required parameter: auth_person key : "+strconv.FormatUint(uint64(key), 10)+" nationality tidak ditemukan",
							"Missing required parameter: auth_person key : "+strconv.FormatUint(uint64(key), 10)+" nationality tidak ditemukan")
					}

					if val, ok := valueMap["position"]; ok {
						if len(val.(string)) > 50 {
							log.Error("Missing required parameter: auth_person key : " + strconv.FormatUint(uint64(key), 10) + " position too long, max 50 character.")
							return lib.CustomError(http.StatusBadRequest,
								"Missing required parameter: auth_person key : "+strconv.FormatUint(uint64(key), 10)+" position too long, max 50 character.",
								"Missing required parameter: auth_person key : "+strconv.FormatUint(uint64(key), 10)+" position too long, max 50 character.")
						}
						ap["position"] = val.(string)
					} else {
						log.Error("Missing required parameter: auth_person key : " + strconv.FormatUint(uint64(key), 10) + " position tidak ditemukan")
						return lib.CustomError(http.StatusBadRequest,
							"Missing required parameter: auth_person key : "+strconv.FormatUint(uint64(key), 10)+" position tidak ditemukan",
							"Missing required parameter: auth_person key : "+strconv.FormatUint(uint64(key), 10)+" position tidak ditemukan")
					}

					if val, ok := valueMap["phone_no"]; ok {
						if len(val.(string)) > 20 {
							log.Error("Missing required parameter: auth_person key : " + strconv.FormatUint(uint64(key), 10) + " phone_no too long, max 20 character.")
							return lib.CustomError(http.StatusBadRequest,
								"Missing required parameter: auth_person key : "+strconv.FormatUint(uint64(key), 10)+" phone_no too long, max 20 character.",
								"Missing required parameter: auth_person key : "+strconv.FormatUint(uint64(key), 10)+" phone_no too long, max 20 character.")
						}
						ap["phone_no"] = val.(string)
					} else {
						log.Error("Missing required parameter: auth_person key : " + strconv.FormatUint(uint64(key), 10) + " phone_no tidak ditemukan")
						return lib.CustomError(http.StatusBadRequest,
							"Missing required parameter: auth_person key : "+strconv.FormatUint(uint64(key), 10)+" phone_no tidak ditemukan",
							"Missing required parameter: auth_person key : "+strconv.FormatUint(uint64(key), 10)+" phone_no tidak ditemukan")
					}

					if val, ok := valueMap["email_address"]; ok {
						if len(val.(string)) > 50 {
							log.Error("Missing required parameter: auth_person key : " + strconv.FormatUint(uint64(key), 10) + " email_address too long, max 50 character.")
							return lib.CustomError(http.StatusBadRequest,
								"Missing required parameter: auth_person key : "+strconv.FormatUint(uint64(key), 10)+" email_address too long, max 50 character.",
								"Missing required parameter: auth_person key : "+strconv.FormatUint(uint64(key), 10)+" email_address too long, max 50 character.")
						}

						if !lib.IsValidEmail(val.(string)) {
							log.Error("Missing required parameter: auth_person key : " + strconv.FormatUint(uint64(key), 10) + " email_address email_address wrong format email.")
							return lib.CustomError(http.StatusBadRequest,
								"Missing required parameter: auth_person key : "+strconv.FormatUint(uint64(key), 10)+" email_address email_address wrong format email.",
								"Missing required parameter: auth_person key : "+strconv.FormatUint(uint64(key), 10)+" email_address email_address wrong format email.")
						}
						ap["email_address"] = val.(string)
					} else {
						log.Error("Missing required parameter: auth_person key : " + strconv.FormatUint(uint64(key), 10) + " email_address tidak ditemukan")
						return lib.CustomError(http.StatusBadRequest,
							"Missing required parameter: auth_person key : "+strconv.FormatUint(uint64(key), 10)+" email_address tidak ditemukan",
							"Missing required parameter: auth_person key : "+strconv.FormatUint(uint64(key), 10)+" email_address tidak ditemukan")
					}

					authPersonData = append(authPersonData, ap)
					key++
				}
			}
		}
	}

	assetY1 := c.FormValue("asset_y1")
	if assetY1 != "" {
		_, err := decimal.NewFromString(assetY1)
		if err != nil {
			log.Error("Wrong input for parameter: asset_y1")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: asset_y1", "Wrong input for parameter: asset_y1")
		}
		paramsInstitutionData["asset_y1"] = assetY1
	}
	assetY2 := c.FormValue("asset_y2")
	if assetY2 != "" {
		_, err := decimal.NewFromString(assetY2)
		if err != nil {
			log.Error("Wrong input for parameter: asset_y2")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: asset_y2", "Wrong input for parameter: asset_y2")
		}
		paramsInstitutionData["asset_y2"] = assetY2
	}
	assetY3 := c.FormValue("asset_y3")
	if assetY3 != "" {
		_, err := decimal.NewFromString(assetY3)
		if err != nil {
			log.Error("Wrong input for parameter: asset_y3")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: asset_y3", "Wrong input for parameter: asset_y3")
		}
		paramsInstitutionData["asset_y3"] = assetY3
	}
	opsProfitY1 := c.FormValue("ops_profit_y1")
	if opsProfitY1 != "" {
		_, err := decimal.NewFromString(opsProfitY1)
		if err != nil {
			log.Error("Wrong input for parameter: ops_profit_y1")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: ops_profit_y1", "Wrong input for parameter: ops_profit_y1")
		}
		paramsInstitutionData["ops_profit_y1"] = opsProfitY1
	}
	opsProfitY2 := c.FormValue("ops_profit_y2")
	if opsProfitY2 != "" {
		_, err := decimal.NewFromString(opsProfitY2)
		if err != nil {
			log.Error("Wrong input for parameter: ops_profit_y2")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: ops_profit_y2", "Wrong input for parameter: ops_profit_y2")
		}
		paramsInstitutionData["ops_profit_y2"] = opsProfitY2
	}
	opsProfitY3 := c.FormValue("ops_profit_y3")
	if opsProfitY3 != "" {
		_, err := decimal.NewFromString(opsProfitY3)
		if err != nil {
			log.Error("Wrong input for parameter: ops_profit_y3")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: ops_profit_y3", "Wrong input for parameter: ops_profit_y3")
		}
		paramsInstitutionData["ops_profit_y3"] = opsProfitY3
	}

	//REKENING BANK // ARRAY max 3
	bankAccount := c.FormValue("bank_account")

	saveBank := false
	var bankData []interface{}
	if isSaveDraft == "1" {
		if bankAccount != "" {
			saveBank = true
			var bankAccountSlice []interface{}
			err = json.Unmarshal([]byte(bankAccount), &bankAccountSlice)
			if err != nil {
				log.Error(err.Error())
				log.Error("Missing required parameter: bank_account")
				return lib.CustomError(http.StatusBadRequest, err.Error(), "Wrong input for parameter: bank_account")
			}
			if len(bankAccountSlice) > 3 {
				log.Error("Missing required parameter: bank_account")
				return lib.CustomError(http.StatusBadRequest, "bank_account hanya max 3 data.", "bank_account hanya max 3 data.")
			}
			key := 1
			bankIsPriority := 0
			if len(bankAccountSlice) > 0 {
				for _, val := range bankAccountSlice {
					bank := make(map[string]interface{})
					valueMap := val.(map[string]interface{})
					if val, ok := valueMap["req_bankacc_key"]; ok {
						if val.(string) != "" {
							n, err := strconv.ParseUint(val.(string), 10, 64)
							if err != nil || n == 0 {
								log.Error("Wrong input parameter: bank_account key : " + strconv.FormatUint(uint64(key), 10) + " req_bankacc_key.")
								return lib.CustomError(http.StatusBadRequest,
									"Wrong input parameter: bank_account key : "+strconv.FormatUint(uint64(key), 10)+" req_bankacc_key.",
									"Wrong input parameter: bank_account key : "+strconv.FormatUint(uint64(key), 10)+" req_bankacc_key.")
							}
						}
						bank["req_bankacc_key"] = val.(string)
					} else {
						log.Error("Missing required parameter: bank_account key : " + strconv.FormatUint(uint64(key), 10) + " req_bankacc_key tidak ditemukan")
						return lib.CustomError(http.StatusBadRequest,
							"Missing required parameter: bank_account key : "+strconv.FormatUint(uint64(key), 10)+" req_bankacc_key tidak ditemukan",
							"Missing required parameter: bank_account key : "+strconv.FormatUint(uint64(key), 10)+" req_bankacc_key tidak ditemukan")
					}
					if val, ok := valueMap["flag_priority"]; ok {
						bank["flag_priority"] = val.(string)
						if val.(string) == "1" {
							bankIsPriority++
						}
					} else {
						log.Error("Missing required parameter: bank_account key : " + strconv.FormatUint(uint64(key), 10) + " flag_priority tidak ditemukan")
						return lib.CustomError(http.StatusBadRequest,
							"Missing required parameter: bank_account key : "+strconv.FormatUint(uint64(key), 10)+" flag_priority tidak ditemukan",
							"Missing required parameter: bank_account key : "+strconv.FormatUint(uint64(key), 10)+" flag_priority tidak ditemukan")
					}
					if val, ok := valueMap["bank_key"]; ok {
						bank["bank_key"] = val.(string)
						if val.(string) != "" {
							n, err := strconv.ParseUint(val.(string), 10, 64)
							if err != nil || n == 0 {
								log.Error("Wrong input for parameter: bank_account - bank_key")
								return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: bank_account - bank_key", "Wrong input for parameter: bank_account - bank_key")
							}
							if len(val.(string)) > 11 {
								log.Error("Wrong input for parameter: bank_account - bank_key too long")
								return lib.CustomError(http.StatusBadRequest, "Missing required parameter: bank_account - bank_key too long, max 11 character", "Missing required parameter: bank_account - bank_key too long, max 11 character")
							}
						}
					} else {
						log.Error("Missing required parameter: bank_account key : " + strconv.FormatUint(uint64(key), 10) + " bank_key tidak ditemukan")
						return lib.CustomError(http.StatusBadRequest,
							"Missing required parameter: bank_account key : "+strconv.FormatUint(uint64(key), 10)+" bank_key tidak ditemukan",
							"Missing required parameter: bank_account key : "+strconv.FormatUint(uint64(key), 10)+" bank_key tidak ditemukan")
					}
					if val, ok := valueMap["account_no"]; ok {
						bank["account_no"] = val.(string)
						if val.(string) != "" {
							if len(val.(string)) > 30 {
								log.Error("Wrong input for parameter: bank_account - account_no too long")
								return lib.CustomError(http.StatusBadRequest,
									"Missing required parameter: bank_account - account_no too long, max 30 character",
									"Missing required parameter: bank_account - account_no too long, max 30 character")
							}
						}
					} else {
						log.Error("Missing required parameter: bank_account key : " + strconv.FormatUint(uint64(key), 10) + " account_no tidak ditemukan")
						return lib.CustomError(http.StatusBadRequest,
							"Missing required parameter: bank_account key : "+strconv.FormatUint(uint64(key), 10)+" account_no tidak ditemukan",
							"Missing required parameter: bank_account key : "+strconv.FormatUint(uint64(key), 10)+" account_no tidak ditemukan")
					}
					if val, ok := valueMap["account_holder_name"]; ok {
						bank["account_holder_name"] = val.(string)
						if val.(string) != "" {
							if len(val.(string)) > 80 {
								log.Error("Wrong input for parameter: bank_account - account_holder_name too long")
								return lib.CustomError(http.StatusBadRequest,
									"Missing required parameter: bank_account - account_holder_name too long, max 80 character",
									"Missing required parameter: bank_account - account_holder_name too long, max 80 character")
							}
						}
					} else {
						log.Error("Missing required parameter: bank_account key : " + strconv.FormatUint(uint64(key), 10) + " account_holder_name tidak ditemukan")
						return lib.CustomError(http.StatusBadRequest,
							"Missing required parameter: bank_account key : "+strconv.FormatUint(uint64(key), 10)+" account_holder_name tidak ditemukan",
							"Missing required parameter: bank_account key : "+strconv.FormatUint(uint64(key), 10)+" account_holder_name tidak ditemukan")
					}
					if val, ok := valueMap["branch_name"]; ok {
						bank["branch_name"] = val.(string)
						if val.(string) != "" {
							if len(val.(string)) > 80 {
								log.Error("Wrong input for parameter: bank_account - branch_name too long")
								return lib.CustomError(http.StatusBadRequest,
									"Missing required parameter: bank_account - branch_name too long, max 80 character",
									"Missing required parameter: bank_account - branch_name too long, max 80 character")
							}
						}
					} else {
						log.Error("Missing required parameter: bank_account key : " + strconv.FormatUint(uint64(key), 10) + " branch_name tidak ditemukan")
						return lib.CustomError(http.StatusBadRequest,
							"Missing required parameter: bank_account key : "+strconv.FormatUint(uint64(key), 10)+" branch_name tidak ditemukan",
							"Missing required parameter: bank_account key : "+strconv.FormatUint(uint64(key), 10)+" branch_name tidak ditemukan")
					}
					if val, ok := valueMap["currency_key"]; ok {
						bank["currency_key"] = val.(string)
						if val.(string) != "" {
							n, err := strconv.ParseUint(val.(string), 10, 64)
							if err != nil || n == 0 {
								log.Error("Wrong input for parameter: bank_account - currency_key")
								return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: bank_account - currency_key", "Wrong input for parameter: bank_account - currency_key")
							}
							if len(val.(string)) > 11 {
								log.Error("Wrong input for parameter: bank_account - currency_key too long")
								return lib.CustomError(http.StatusBadRequest, "Missing required parameter: bank_account - currency_key too long, max 11 character", "Missing required parameter: bank_account - currency_key too long, max 11 character")
							}
						}
					} else {
						log.Error("Missing required parameter: bank_account key : " + strconv.FormatUint(uint64(key), 10) + " currency_key tidak ditemukan")
						return lib.CustomError(http.StatusBadRequest,
							"Missing required parameter: bank_account key : "+strconv.FormatUint(uint64(key), 10)+" currency_key tidak ditemukan",
							"Missing required parameter: bank_account key : "+strconv.FormatUint(uint64(key), 10)+" currency_key tidak ditemukan")
					}

					bankData = append(bankData, bank)
					key++
				}

				if bankIsPriority > 1 {
					log.Error("Missing required parameter: bank_account priority hanya boleh 1 data")
					return lib.CustomError(http.StatusBadRequest, "bank_account priority hanya boleh 1 data.", "bank_account priority hanya boleh 1 data.")
				}
			}
		}
	}

	if isSaveDraft == "0" { //save with validation
		if tinNumber == "" {
			log.Error("Missing required parameter: tin_number")
			return lib.CustomError(http.StatusBadRequest, "tin_number can not be blank", "tin_number can not be blank")
		}
		if shortName == "" {
			log.Error("Missing required parameter: short_name")
			return lib.CustomError(http.StatusBadRequest, "short_name can not be blank", "short_name can not be blank")
		}

		if establishedCity == "" {
			log.Error("Missing required parameter: established_city")
			return lib.CustomError(http.StatusBadRequest, "established_city can not be blank", "established_city can not be blank")
		}

		if establishedDate == "" {
			log.Error("Missing required parameter: established_date")
			return lib.CustomError(http.StatusBadRequest, "established_date can not be blank", "established_date can not be blank")
		}

		if deedNo == "" {
			log.Error("Missing required parameter: deed_no")
			return lib.CustomError(http.StatusBadRequest, "deed_no can not be blank", "deed_no can not be blank")
		}

		if lsEstablishValidationNo == "" {
			log.Error("Missing required parameter: ls_establish_validation_no")
			return lib.CustomError(http.StatusBadRequest, "ls_establish_validation_no can not be blank", "ls_establish_validation_no can not be blank")
		}

		if lsEstablishValidationDate == "" {
			log.Error("Missing required parameter: ls_establish_validation_date")
			return lib.CustomError(http.StatusBadRequest, "ls_establish_validation_date can not be blank", "ls_establish_validation_date can not be blank")
		}

		if lastChangeAaNo == "" {
			log.Error("Missing required parameter: last_change_aa_no")
			return lib.CustomError(http.StatusBadRequest, "last_change_aa_no can not be blank", "last_change_aa_no can not be blank")
		}

		if lastChangeAaDate == "" {
			log.Error("Missing required parameter: last_change_aa_date")
			return lib.CustomError(http.StatusBadRequest, "last_change_aa_date can not be blank", "last_change_aa_date can not be blank")
		}

		if lsLastChangeAaNo == "" {
			log.Error("Missing required parameter: ls_last_change_aa_no")
			return lib.CustomError(http.StatusBadRequest, "ls_last_change_aa_no can not be blank", "ls_last_change_aa_no can not be blank")
		}
		if lsLastChangeAaDate == "" {
			log.Error("Missing required parameter: ls_last_change_aa_date")
			return lib.CustomError(http.StatusBadRequest, "ls_last_change_aa_date can not be blank", "ls_last_change_aa_date can not be blank")
		}

		if managementDeedNo == "" {
			log.Error("Missing required parameter: management_deed_no")
			return lib.CustomError(http.StatusBadRequest, "management_deed_no can not be blank", "management_deed_no can not be blank")
		}
		if managementDeedDate == "" {
			log.Error("Missing required parameter: management_deed_date")
			return lib.CustomError(http.StatusBadRequest, "management_deed_date can not be blank", "management_deed_date can not be blank")
		}

		if lsMgtChangeDeedNo == "" {
			log.Error("Missing required parameter: ls_mgt_change_deed_no")
			return lib.CustomError(http.StatusBadRequest, "ls_mgt_change_deed_no can not be blank", "ls_mgt_change_deed_no can not be blank")
		}
		if lsMgtChangeDeedDate == "" {
			log.Error("Missing required parameter: ls_mgt_change_deed_date")
			return lib.CustomError(http.StatusBadRequest, "ls_mgt_change_deed_date can not be blank", "ls_mgt_change_deed_date can not be blank")
		}

		if skdLicenseNo == "" {
			log.Error("Missing required parameter: skd_license_no")
			return lib.CustomError(http.StatusBadRequest, "skd_license_no can not be blank", "skd_license_no can not be blank")
		}
		if skdLicenseDate == "" {
			log.Error("Missing required parameter: skd_license_date")
			return lib.CustomError(http.StatusBadRequest, "skd_license_date can not be blank", "skd_license_date can not be blank")
		}

		if bizLicenseNo == "" {
			log.Error("Missing required parameter: biz_license_no")
			return lib.CustomError(http.StatusBadRequest, "biz_license_no can not be blank", "biz_license_no can not be blank")
		}
		if bizLicenseDate == "" {
			log.Error("Missing required parameter: biz_license_date")
			return lib.CustomError(http.StatusBadRequest, "biz_license_date can not be blank", "biz_license_date can not be blank")
		}

		if nibNo == "" {
			log.Error("Missing required parameter: nib_no")
			return lib.CustomError(http.StatusBadRequest, "nib_no can not be blank", "nib_no can not be blank")
		}
		if nibDate == "" {
			log.Error("Missing required parameter: nib_date")
			return lib.CustomError(http.StatusBadRequest, "nib_date can not be blank", "nib_date can not be blank")
		}

		if mobileNo == "" {
			log.Error("Missing required parameter: mobile_no")
			return lib.CustomError(http.StatusBadRequest, "mobile_no can not be blank", "mobile_no can not be blank")
		}

		if faxNo == "" {
			log.Error("Missing required parameter: fax_no")
			return lib.CustomError(http.StatusBadRequest, "fax_no can not be blank", "fax_no can not be blank")
		}
		if emailAddress == "" {
			log.Error("Missing required parameter: email_address")
			return lib.CustomError(http.StatusBadRequest, "email_address can not be blank", "email_address can not be blank")
		}

		//DOMICILE
		if domicileAddress == "" {
			log.Error("Missing required parameter: domicile_address")
			return lib.CustomError(http.StatusBadRequest, "domicile_address can not be blank", "domicile_address can not be blank")
		}
		if domicileSubdistric == "" {
			log.Error("Missing required parameter: domicile_subdistric")
			return lib.CustomError(http.StatusBadRequest, "domicile_subdistric can not be blank", "domicile_subdistric can not be blank")
		}
		if domicileCity == "" {
			log.Error("Missing required parameter: domicile_city")
			return lib.CustomError(http.StatusBadRequest, "domicile_city can not be blank", "domicile_city can not be blank")
		}
		if domicileProvince == "" {
			log.Error("Missing required parameter: domicile_province")
			return lib.CustomError(http.StatusBadRequest, "domicile_province can not be blank", "domicile_province can not be blank")
		}
		if domicilePostalcode == "" {
			log.Error("Missing required parameter: domicile_postalcode")
			return lib.CustomError(http.StatusBadRequest, "domicile_postalcode can not be blank", "domicile_postalcode can not be blank")
		}

		//CORRESPONDENCE
		if correspondenceAddress == "" {
			log.Error("Missing required parameter: correspondence_address")
			return lib.CustomError(http.StatusBadRequest, "correspondence_address can not be blank", "correspondence_address can not be blank")
		}
		if correspondenceSubdistric == "" {
			log.Error("Missing required parameter: correspondence_subdistric")
			return lib.CustomError(http.StatusBadRequest, "correspondence_subdistric can not be blank", "correspondence_subdistric can not be blank")
		}
		if correspondenceCity == "" {
			log.Error("Missing required parameter: correspondence_city")
			return lib.CustomError(http.StatusBadRequest, "correspondence_city can not be blank", "correspondence_city can not be blank")
		}
		if correspondenceProvince == "" {
			log.Error("Missing required parameter: correspondence_province")
			return lib.CustomError(http.StatusBadRequest, "correspondence_province can not be blank", "correspondence_province can not be blank")
		}
		if correspondencePostalcode == "" {
			log.Error("Missing required parameter: correspondence_postalcode")
			return lib.CustomError(http.StatusBadRequest, "correspondence_postalcode can not be blank", "correspondence_postalcode can not be blank")
		}

		if intitutionType == "" {
			log.Error("Missing required parameter: intitution_type")
			return lib.CustomError(http.StatusBadRequest, "intitution_type can not be blank", "intitution_type can not be blank")
		}
		if intitutionClassification == "" {
			log.Error("Missing required parameter: intitution_classification")
			return lib.CustomError(http.StatusBadRequest, "intitution_classification can not be blank", "intitution_classification can not be blank")
		}
		if intitutionCharacteristic == "" {
			log.Error("Missing required parameter: intitution_characteristic")
			return lib.CustomError(http.StatusBadRequest, "intitution_characteristic can not be blank", "intitution_characteristic can not be blank")
		}
		if intitutionBusinessType == "" {
			log.Error("Missing required parameter: intitution_business_type")
			return lib.CustomError(http.StatusBadRequest, "intitution_business_type can not be blank", "intitution_business_type can not be blank")
		}

		if instiAnnuallyIncome == "" {
			log.Error("Missing required parameter: insti_annually_income")
			return lib.CustomError(http.StatusBadRequest, "insti_annually_income can not be blank", "insti_annually_income can not be blank")
		}
		if instiSourceOfIncome == "" {
			log.Error("Missing required parameter: insti_source_of_income")
			return lib.CustomError(http.StatusBadRequest, "insti_source_of_income can not be blank", "insti_source_of_income can not be blank")
		}
		if instiInvestmentPurpose == "" {
			log.Error("Missing required parameter: insti_investment_purpose")
			return lib.CustomError(http.StatusBadRequest, "insti_investment_purpose can not be blank", "insti_investment_purpose can not be blank")
		}

		//BENEFICIAL OWNER
		if boName == "" {
			log.Error("Missing required parameter: bo_name")
			return lib.CustomError(http.StatusBadRequest, "bo_name can not be blank", "bo_name can not be blank")
		}
		if boIdnumber == "" {
			log.Error("Missing required parameter: bo_idnumber")
			return lib.CustomError(http.StatusBadRequest, "bo_idnumber can not be blank", "bo_idnumber can not be blank")
		}
		if boBusiness == "" {
			log.Error("Missing required parameter: bo_business")
			return lib.CustomError(http.StatusBadRequest, "bo_business can not be blank", "bo_business can not be blank")
		}
		if boAnnuallyIncome == "" {
			log.Error("Missing required parameter: bo_annually_income")
			return lib.CustomError(http.StatusBadRequest, "bo_annually_income can not be blank", "bo_annually_income can not be blank")
		}
		if boRelation == "" {
			log.Error("Missing required parameter: bo_relation")
			return lib.CustomError(http.StatusBadRequest, "bo_relation can not be blank", "bo_relation can not be blank")
		}

		//SHARES HOLDER //ARRAY
		if sharesHolder == "" {
			log.Error("Missing required parameter: shares_holder")
			return lib.CustomError(http.StatusBadRequest, "shares_holder can not be blank", "shares_holder can not be blank")
		} else {
			saveSharesHolder = true
			var sharesHolderSlice []interface{}
			err = json.Unmarshal([]byte(sharesHolder), &sharesHolderSlice)
			if err != nil {
				log.Error(err.Error())
				log.Error("Missing required parameter: shares_holder")
				return lib.CustomError(http.StatusBadRequest, err.Error(), "Wrong input for parameter: shares_holder")
			}

			if len(sharesHolderSlice) == 0 {
				log.Error("Missing required parameter: shares_holder")
				return lib.CustomError(http.StatusBadRequest, "shares_holder harus diisi minimal 1 data.", "shares_holder harus diisi minimal 1 data.")
			}

			key := 1
			if len(sharesHolderSlice) > 0 {
				for _, val := range sharesHolderSlice {
					sh := make(map[string]interface{})
					valueMap := val.(map[string]interface{})
					if val, ok := valueMap["inst_shares_holder_key"]; ok {
						if val.(string) != "" {
							n, err := strconv.ParseUint(val.(string), 10, 64)
							if err != nil || n == 0 {
								log.Error("Wrong input parameter: shares_holder key : " + strconv.FormatUint(uint64(key), 10) + " inst_shares_holder_key.")
								return lib.CustomError(http.StatusBadRequest,
									"Wrong input parameter: shares_holder key : "+strconv.FormatUint(uint64(key), 10)+" inst_shares_holder_key.",
									"Wrong input parameter: shares_holder key : "+strconv.FormatUint(uint64(key), 10)+" inst_shares_holder_key.")
							}
						}
						sh["inst_shares_holder_key"] = val.(string)
					} else {
						log.Error("Missing required parameter: shares_holder key : " + strconv.FormatUint(uint64(key), 10) + " inst_shares_holder_key tidak ditemukan")
						return lib.CustomError(http.StatusBadRequest,
							"Missing required parameter: shares_holder key : "+strconv.FormatUint(uint64(key), 10)+" inst_shares_holder_key tidak ditemukan",
							"Missing required parameter: shares_holder key : "+strconv.FormatUint(uint64(key), 10)+" inst_shares_holder_key tidak ditemukan")
					}

					if val, ok := valueMap["nationality"]; ok {
						if val.(string) == "" {
							log.Error("Missing required parameter: shares_holder key : " + strconv.FormatUint(uint64(key), 10) + " nationality can not be blank")
							return lib.CustomError(http.StatusBadRequest,
								"Missing required parameter: shares_holder key : "+strconv.FormatUint(uint64(key), 10)+" nationality can not be blank",
								"Missing required parameter: shares_holder key : "+strconv.FormatUint(uint64(key), 10)+" nationality can not be blank")
						} else {
							n, err := strconv.ParseUint(val.(string), 10, 64)
							if err != nil || n == 0 {
								log.Error("Wrong input parameter: shares_holder key : " + strconv.FormatUint(uint64(key), 10) + " nationality.")
								return lib.CustomError(http.StatusBadRequest,
									"Wrong input parameter: shares_holder key : "+strconv.FormatUint(uint64(key), 10)+" nationality.",
									"Wrong input parameter: shares_holder key : "+strconv.FormatUint(uint64(key), 10)+" nationality.")
							}
							sh["nationality"] = val.(string)
						}
					} else {
						log.Error("Missing required parameter: shares_holder key : " + strconv.FormatUint(uint64(key), 10) + " nationality tidak ditemukan")
						return lib.CustomError(http.StatusBadRequest,
							"Missing required parameter: shares_holder key : "+strconv.FormatUint(uint64(key), 10)+" nationality tidak ditemukan",
							"Missing required parameter: shares_holder key : "+strconv.FormatUint(uint64(key), 10)+" nationality tidak ditemukan")
					}

					if val, ok := valueMap["holder_full_name"]; ok {
						if val.(string) == "" {
							log.Error("Missing required parameter: shares_holder key : " + strconv.FormatUint(uint64(key), 10) + " holder_full_name can not be blank")
							return lib.CustomError(http.StatusBadRequest,
								"Missing required parameter: shares_holder key : "+strconv.FormatUint(uint64(key), 10)+" holder_full_name can not be blank",
								"Missing required parameter: shares_holder key : "+strconv.FormatUint(uint64(key), 10)+" holder_full_name can not be blank")
						} else {
							if len(val.(string)) > 150 {
								log.Error("Missing required parameter: shares_holder key : " + strconv.FormatUint(uint64(key), 10) + " holder_full_name too long, max 150 character.")
								return lib.CustomError(http.StatusBadRequest,
									"Missing required parameter: shares_holder key : "+strconv.FormatUint(uint64(key), 10)+" holder_full_name too long, max 150 character.",
									"Missing required parameter: shares_holder key : "+strconv.FormatUint(uint64(key), 10)+" holder_full_name too long, max 150 character.")
							}
							sh["holder_full_name"] = val.(string)
						}
					} else {
						log.Error("Missing required parameter: shares_holder key : " + strconv.FormatUint(uint64(key), 10) + " holder_full_name tidak ditemukan")
						return lib.CustomError(http.StatusBadRequest,
							"Missing required parameter: shares_holder key : "+strconv.FormatUint(uint64(key), 10)+" holder_full_name tidak ditemukan",
							"Missing required parameter: shares_holder key : "+strconv.FormatUint(uint64(key), 10)+" holder_full_name tidak ditemukan")
					}

					if val, ok := valueMap["idcard_type"]; ok {
						if val.(string) == "" {
							log.Error("Missing required parameter: shares_holder key : " + strconv.FormatUint(uint64(key), 10) + " idcard_type can not be blank")
							return lib.CustomError(http.StatusBadRequest,
								"Missing required parameter: shares_holder key : "+strconv.FormatUint(uint64(key), 10)+" idcard_type can not be blank",
								"Missing required parameter: shares_holder key : "+strconv.FormatUint(uint64(key), 10)+" idcard_type can not be blank")
						} else {
							n, err := strconv.ParseUint(val.(string), 10, 64)
							if err != nil || n == 0 {
								log.Error("Wrong input parameter: shares_holder key : " + strconv.FormatUint(uint64(key), 10) + " idcard_type.")
								return lib.CustomError(http.StatusBadRequest,
									"Wrong input parameter: shares_holder key : "+strconv.FormatUint(uint64(key), 10)+" idcard_type.",
									"Wrong input parameter: shares_holder key : "+strconv.FormatUint(uint64(key), 10)+" idcard_type.")
							}
							sh["idcard_type"] = val.(string)
						}
					} else {
						log.Error("Missing required parameter: shares_holder key : " + strconv.FormatUint(uint64(key), 10) + " idcard_type tidak ditemukan")
						return lib.CustomError(http.StatusBadRequest,
							"Missing required parameter: shares_holder key : "+strconv.FormatUint(uint64(key), 10)+" idcard_type tidak ditemukan",
							"Missing required parameter: shares_holder key : "+strconv.FormatUint(uint64(key), 10)+" idcard_type tidak ditemukan")
					}

					if val, ok := valueMap["idcard_no"]; ok {
						if val.(string) == "" {
							log.Error("Missing required parameter: shares_holder key : " + strconv.FormatUint(uint64(key), 10) + " idcard_no can not be blank")
							return lib.CustomError(http.StatusBadRequest,
								"Missing required parameter: shares_holder key : "+strconv.FormatUint(uint64(key), 10)+" idcard_no can not be blank",
								"Missing required parameter: shares_holder key : "+strconv.FormatUint(uint64(key), 10)+" idcard_no can not be blank")
						} else {
							if len(val.(string)) > 20 {
								log.Error("Missing required parameter: shares_holder key : " + strconv.FormatUint(uint64(key), 10) + " idcard_no too long, max 20 character.")
								return lib.CustomError(http.StatusBadRequest,
									"Missing required parameter: shares_holder key : "+strconv.FormatUint(uint64(key), 10)+" idcard_no too long, max 20 character.",
									"Missing required parameter: shares_holder key : "+strconv.FormatUint(uint64(key), 10)+" idcard_no too long, max 20 character.")
							}
							sh["idcard_no"] = val.(string)
						}
					} else {
						log.Error("Missing required parameter: shares_holder key : " + strconv.FormatUint(uint64(key), 10) + " idcard_no tidak ditemukan")
						return lib.CustomError(http.StatusBadRequest,
							"Missing required parameter: shares_holder key : "+strconv.FormatUint(uint64(key), 10)+" idcard_no tidak ditemukan",
							"Missing required parameter: shares_holder key : "+strconv.FormatUint(uint64(key), 10)+" idcard_no tidak ditemukan")
					}

					if val, ok := valueMap["holder_dob"]; ok {
						if val.(string) == "" {
							log.Error("Missing required parameter: shares_holder key : " + strconv.FormatUint(uint64(key), 10) + " holder_dob can not be blank")
							return lib.CustomError(http.StatusBadRequest,
								"Missing required parameter: shares_holder key : "+strconv.FormatUint(uint64(key), 10)+" holder_dob can not be blank",
								"Missing required parameter: shares_holder key : "+strconv.FormatUint(uint64(key), 10)+" holder_dob can not be blank")
						} else {
							holderDob := val.(string) + " 00:00:00"
							date, _ := time.Parse(layout, holderDob)
							dateStr := date.Format(layout)
							sh["holder_dob"] = dateStr
						}
					} else {
						log.Error("Missing required parameter: shares_holder key : " + strconv.FormatUint(uint64(key), 10) + " holder_dob tidak ditemukan")
						return lib.CustomError(http.StatusBadRequest,
							"Missing required parameter: shares_holder key : "+strconv.FormatUint(uint64(key), 10)+" holder_dob tidak ditemukan",
							"Missing required parameter: shares_holder key : "+strconv.FormatUint(uint64(key), 10)+" holder_dob tidak ditemukan")
					}

					if val, ok := valueMap["shares_percent"]; ok {
						if val.(string) == "" {
							log.Error("Missing required parameter: shares_holder key : " + strconv.FormatUint(uint64(key), 10) + " shares_percent can not be blank")
							return lib.CustomError(http.StatusBadRequest,
								"Missing required parameter: shares_holder key : "+strconv.FormatUint(uint64(key), 10)+" shares_percent can not be blank",
								"Missing required parameter: shares_holder key : "+strconv.FormatUint(uint64(key), 10)+" shares_percent can not be blank")
						} else {
							_, err := decimal.NewFromString(val.(string))
							if err != nil {
								log.Error("Wrong input parameter: shares_holder key : " + strconv.FormatUint(uint64(key), 10) + " shares_percent")
								return lib.CustomError(http.StatusBadRequest,
									"Wrong input parameter: shares_holder key : "+strconv.FormatUint(uint64(key), 10)+" shares_percent",
									"Wrong input parameter: shares_holder key : "+strconv.FormatUint(uint64(key), 10)+" shares_percent")
							}
							sh["shares_percent"] = val.(string)
						}
					} else {
						log.Error("Missing required parameter: shares_holder key : " + strconv.FormatUint(uint64(key), 10) + " shares_percent tidak ditemukan")
						return lib.CustomError(http.StatusBadRequest,
							"Missing required parameter: shares_holder key : "+strconv.FormatUint(uint64(key), 10)+" shares_percent tidak ditemukan",
							"Missing required parameter: shares_holder key : "+strconv.FormatUint(uint64(key), 10)+" shares_percent tidak ditemukan")
					}

					sharesHolderData = append(sharesHolderData, sh)
					key++
				}
			} else {
				log.Error("Missing required parameter: shares_holder harus min 1 data")
				return lib.CustomError(http.StatusBadRequest, "shares_holder harus min 1 data.", "shares_holder harus min 1 data.")
			}
		}

		//AUTH PERSON //ARRAY
		if authPerson == "" {
			log.Error("Missing required parameter: auth_person")
			return lib.CustomError(http.StatusBadRequest, "auth_person can not be blank", "auth_person can not be blank")
		} else {
			if authPerson != "" {
				saveAuthPerson = true
				var authPersonSlice []interface{}
				err = json.Unmarshal([]byte(authPerson), &authPersonSlice)
				if err != nil {
					log.Error(err.Error())
					log.Error("Missing required parameter: auth_person")
					return lib.CustomError(http.StatusBadRequest, err.Error(), "Wrong input for parameter: auth_person")
				}

				if len(authPersonSlice) == 0 {
					log.Error("Missing required parameter: auth_person")
					return lib.CustomError(http.StatusBadRequest, "auth_person harus diisi minimal 1 data.", "auth_person harus diisi minimal 1 data.")
				}

				key := 1
				if len(authPersonSlice) > 0 {
					for _, val := range authPersonSlice {
						ap := make(map[string]interface{})
						valueMap := val.(map[string]interface{})
						if val, ok := valueMap["insti_auth_person_key"]; ok {
							if val.(string) != "" {
								n, err := strconv.ParseUint(val.(string), 10, 64)
								if err != nil || n == 0 {
									log.Error("Wrong input parameter: shares_holder key : " + strconv.FormatUint(uint64(key), 10) + " insti_auth_person_key.")
									return lib.CustomError(http.StatusBadRequest,
										"Wrong input parameter: shares_holder key : "+strconv.FormatUint(uint64(key), 10)+" insti_auth_person_key.",
										"Wrong input parameter: shares_holder key : "+strconv.FormatUint(uint64(key), 10)+" insti_auth_person_key.")
								}
							}
							ap["insti_auth_person_key"] = val.(string)
						} else {
							log.Error("Missing required parameter: auth_person key : " + strconv.FormatUint(uint64(key), 10) + " insti_auth_person_key tidak ditemukan")
							return lib.CustomError(http.StatusBadRequest,
								"Missing required parameter: auth_person key : "+strconv.FormatUint(uint64(key), 10)+" insti_auth_person_key tidak ditemukan",
								"Missing required parameter: auth_person key : "+strconv.FormatUint(uint64(key), 10)+" insti_auth_person_key tidak ditemukan")
						}

						if val, ok := valueMap["full_name"]; ok {
							if val.(string) == "" {
								log.Error("Missing required parameter: shares_holder key : " + strconv.FormatUint(uint64(key), 10) + " shares_percent can not be blank")
								return lib.CustomError(http.StatusBadRequest,
									"Missing required parameter: shares_holder key : "+strconv.FormatUint(uint64(key), 10)+" shares_percent can not be blank",
									"Missing required parameter: shares_holder key : "+strconv.FormatUint(uint64(key), 10)+" shares_percent can not be blank")
							} else {
								if len(val.(string)) > 100 {
									log.Error("Missing required parameter: auth_person key : " + strconv.FormatUint(uint64(key), 10) + " full_name too long, max 100 character.")
									return lib.CustomError(http.StatusBadRequest,
										"Missing required parameter: auth_person key : "+strconv.FormatUint(uint64(key), 10)+" full_name too long, max 100 character.",
										"Missing required parameter: auth_person key : "+strconv.FormatUint(uint64(key), 10)+" full_name too long, max 100 character.")
								}
								ap["full_name"] = val.(string)
							}
						} else {
							log.Error("Missing required parameter: auth_person key : " + strconv.FormatUint(uint64(key), 10) + " full_name tidak ditemukan")
							return lib.CustomError(http.StatusBadRequest,
								"Missing required parameter: auth_person key : "+strconv.FormatUint(uint64(key), 10)+" full_name tidak ditemukan",
								"Missing required parameter: auth_person key : "+strconv.FormatUint(uint64(key), 10)+" full_name tidak ditemukan")
						}

						if val, ok := valueMap["idcard_type"]; ok {
							if val.(string) == "" {
								log.Error("Missing required parameter: shares_holder key : " + strconv.FormatUint(uint64(key), 10) + " shares_percent can not be blank")
								return lib.CustomError(http.StatusBadRequest,
									"Missing required parameter: shares_holder key : "+strconv.FormatUint(uint64(key), 10)+" shares_percent can not be blank",
									"Missing required parameter: shares_holder key : "+strconv.FormatUint(uint64(key), 10)+" shares_percent can not be blank")
							} else {
								n, err := strconv.ParseUint(val.(string), 10, 64)
								if err != nil || n == 0 {
									log.Error("Wrong input parameter: auth_person key : " + strconv.FormatUint(uint64(key), 10) + " idcard_type.")
									return lib.CustomError(http.StatusBadRequest,
										"Wrong input parameter: auth_person key : "+strconv.FormatUint(uint64(key), 10)+" idcard_type.",
										"Wrong input parameter: auth_person key : "+strconv.FormatUint(uint64(key), 10)+" idcard_type.")
								}
								ap["idcard_type"] = val.(string)
							}
						} else {
							log.Error("Missing required parameter: auth_person key : " + strconv.FormatUint(uint64(key), 10) + " idcard_type tidak ditemukan")
							return lib.CustomError(http.StatusBadRequest,
								"Missing required parameter: auth_person key : "+strconv.FormatUint(uint64(key), 10)+" idcard_type tidak ditemukan",
								"Missing required parameter: auth_person key : "+strconv.FormatUint(uint64(key), 10)+" idcard_type tidak ditemukan")
						}

						if val, ok := valueMap["idcard_no"]; ok {
							if val.(string) == "" {
								log.Error("Missing required parameter: shares_holder key : " + strconv.FormatUint(uint64(key), 10) + " shares_percent can not be blank")
								return lib.CustomError(http.StatusBadRequest,
									"Missing required parameter: shares_holder key : "+strconv.FormatUint(uint64(key), 10)+" shares_percent can not be blank",
									"Missing required parameter: shares_holder key : "+strconv.FormatUint(uint64(key), 10)+" shares_percent can not be blank")
							} else {
								if len(val.(string)) > 20 {
									log.Error("Missing required parameter: auth_person key : " + strconv.FormatUint(uint64(key), 10) + " idcard_no too long, max 20 character.")
									return lib.CustomError(http.StatusBadRequest,
										"Missing required parameter: auth_person key : "+strconv.FormatUint(uint64(key), 10)+" idcard_no too long, max 20 character.",
										"Missing required parameter: auth_person key : "+strconv.FormatUint(uint64(key), 10)+" idcard_no too long, max 20 character.")
								}
								ap["idcard_no"] = val.(string)
							}
						} else {
							log.Error("Missing required parameter: auth_person key : " + strconv.FormatUint(uint64(key), 10) + " idcard_no tidak ditemukan")
							return lib.CustomError(http.StatusBadRequest,
								"Missing required parameter: auth_person key : "+strconv.FormatUint(uint64(key), 10)+" idcard_no tidak ditemukan",
								"Missing required parameter: auth_person key : "+strconv.FormatUint(uint64(key), 10)+" idcard_no tidak ditemukan")
						}

						if val, ok := valueMap["nationality"]; ok {
							if val.(string) == "" {
								log.Error("Missing required parameter: shares_holder key : " + strconv.FormatUint(uint64(key), 10) + " shares_percent can not be blank")
								return lib.CustomError(http.StatusBadRequest,
									"Missing required parameter: shares_holder key : "+strconv.FormatUint(uint64(key), 10)+" shares_percent can not be blank",
									"Missing required parameter: shares_holder key : "+strconv.FormatUint(uint64(key), 10)+" shares_percent can not be blank")
							} else {
								n, err := strconv.ParseUint(val.(string), 10, 64)
								if err != nil || n == 0 {
									log.Error("Wrong input parameter: auth_person key : " + strconv.FormatUint(uint64(key), 10) + " nationality.")
									return lib.CustomError(http.StatusBadRequest,
										"Wrong input parameter: auth_person key : "+strconv.FormatUint(uint64(key), 10)+" nationality.",
										"Wrong input parameter: auth_person key : "+strconv.FormatUint(uint64(key), 10)+" nationality.")
								}
								ap["nationality"] = val.(string)
							}
						} else {
							log.Error("Missing required parameter: auth_person key : " + strconv.FormatUint(uint64(key), 10) + " nationality tidak ditemukan")
							return lib.CustomError(http.StatusBadRequest,
								"Missing required parameter: auth_person key : "+strconv.FormatUint(uint64(key), 10)+" nationality tidak ditemukan",
								"Missing required parameter: auth_person key : "+strconv.FormatUint(uint64(key), 10)+" nationality tidak ditemukan")
						}

						if val, ok := valueMap["person_dob"]; ok {
							if val.(string) == "" {
								log.Error("Missing required parameter: shares_holder key : " + strconv.FormatUint(uint64(key), 10) + " shares_percent can not be blank")
								return lib.CustomError(http.StatusBadRequest,
									"Missing required parameter: shares_holder key : "+strconv.FormatUint(uint64(key), 10)+" shares_percent can not be blank",
									"Missing required parameter: shares_holder key : "+strconv.FormatUint(uint64(key), 10)+" shares_percent can not be blank")
							} else {
								if val.(string) != "" {
									personDob := val.(string) + " 00:00:00"
									date, _ := time.Parse(layout, personDob)
									dateStr := date.Format(layout)
									ap["person_dob"] = dateStr
								}
							}
						} else {
							log.Error("Missing required parameter: auth_person key : " + strconv.FormatUint(uint64(key), 10) + " person_dob tidak ditemukan")
							return lib.CustomError(http.StatusBadRequest,
								"Missing required parameter: auth_person key : "+strconv.FormatUint(uint64(key), 10)+" person_dob tidak ditemukan",
								"Missing required parameter: auth_person key : "+strconv.FormatUint(uint64(key), 10)+" person_dob tidak ditemukan")
						}

						if val, ok := valueMap["position"]; ok {
							if val.(string) == "" {
								log.Error("Missing required parameter: shares_holder key : " + strconv.FormatUint(uint64(key), 10) + " shares_percent can not be blank")
								return lib.CustomError(http.StatusBadRequest,
									"Missing required parameter: shares_holder key : "+strconv.FormatUint(uint64(key), 10)+" shares_percent can not be blank",
									"Missing required parameter: shares_holder key : "+strconv.FormatUint(uint64(key), 10)+" shares_percent can not be blank")
							} else {
								if len(val.(string)) > 50 {
									log.Error("Missing required parameter: auth_person key : " + strconv.FormatUint(uint64(key), 10) + " position too long, max 50 character.")
									return lib.CustomError(http.StatusBadRequest,
										"Missing required parameter: auth_person key : "+strconv.FormatUint(uint64(key), 10)+" position too long, max 50 character.",
										"Missing required parameter: auth_person key : "+strconv.FormatUint(uint64(key), 10)+" position too long, max 50 character.")
								}
								ap["position"] = val.(string)
							}
						} else {
							log.Error("Missing required parameter: auth_person key : " + strconv.FormatUint(uint64(key), 10) + " position tidak ditemukan")
							return lib.CustomError(http.StatusBadRequest,
								"Missing required parameter: auth_person key : "+strconv.FormatUint(uint64(key), 10)+" position tidak ditemukan",
								"Missing required parameter: auth_person key : "+strconv.FormatUint(uint64(key), 10)+" position tidak ditemukan")
						}

						if val, ok := valueMap["phone_no"]; ok {
							if val.(string) == "" {
								log.Error("Missing required parameter: shares_holder key : " + strconv.FormatUint(uint64(key), 10) + " shares_percent can not be blank")
								return lib.CustomError(http.StatusBadRequest,
									"Missing required parameter: shares_holder key : "+strconv.FormatUint(uint64(key), 10)+" shares_percent can not be blank",
									"Missing required parameter: shares_holder key : "+strconv.FormatUint(uint64(key), 10)+" shares_percent can not be blank")
							} else {
								if len(val.(string)) > 20 {
									log.Error("Missing required parameter: auth_person key : " + strconv.FormatUint(uint64(key), 10) + " phone_no too long, max 20 character.")
									return lib.CustomError(http.StatusBadRequest,
										"Missing required parameter: auth_person key : "+strconv.FormatUint(uint64(key), 10)+" phone_no too long, max 20 character.",
										"Missing required parameter: auth_person key : "+strconv.FormatUint(uint64(key), 10)+" phone_no too long, max 20 character.")
								}
								ap["phone_no"] = val.(string)
							}
						} else {
							log.Error("Missing required parameter: auth_person key : " + strconv.FormatUint(uint64(key), 10) + " phone_no tidak ditemukan")
							return lib.CustomError(http.StatusBadRequest,
								"Missing required parameter: auth_person key : "+strconv.FormatUint(uint64(key), 10)+" phone_no tidak ditemukan",
								"Missing required parameter: auth_person key : "+strconv.FormatUint(uint64(key), 10)+" phone_no tidak ditemukan")
						}

						if val, ok := valueMap["email_address"]; ok {
							if val.(string) == "" {
								log.Error("Missing required parameter: shares_holder key : " + strconv.FormatUint(uint64(key), 10) + " shares_percent can not be blank")
								return lib.CustomError(http.StatusBadRequest,
									"Missing required parameter: shares_holder key : "+strconv.FormatUint(uint64(key), 10)+" shares_percent can not be blank",
									"Missing required parameter: shares_holder key : "+strconv.FormatUint(uint64(key), 10)+" shares_percent can not be blank")
							} else {
								if len(val.(string)) > 50 {
									log.Error("Missing required parameter: auth_person key : " + strconv.FormatUint(uint64(key), 10) + " email_address too long, max 50 character.")
									return lib.CustomError(http.StatusBadRequest,
										"Missing required parameter: auth_person key : "+strconv.FormatUint(uint64(key), 10)+" email_address too long, max 50 character.",
										"Missing required parameter: auth_person key : "+strconv.FormatUint(uint64(key), 10)+" email_address too long, max 50 character.")
								}
								ap["email_address"] = val.(string)
							}
						} else {
							log.Error("Missing required parameter: auth_person key : " + strconv.FormatUint(uint64(key), 10) + " email_address tidak ditemukan")
							return lib.CustomError(http.StatusBadRequest,
								"Missing required parameter: auth_person key : "+strconv.FormatUint(uint64(key), 10)+" email_address tidak ditemukan",
								"Missing required parameter: auth_person key : "+strconv.FormatUint(uint64(key), 10)+" email_address tidak ditemukan")
						}

						authPersonData = append(authPersonData, ap)
						key++
					}
				}
			}
		}

		if assetY1 == "" {
			log.Error("Missing required parameter: asset_y1")
			return lib.CustomError(http.StatusBadRequest, "asset_y1 can not be blank", "asset_y1 can not be blank")
		}
		if opsProfitY3 == "" {
			log.Error("Missing required parameter: ops_profit_y3")
			return lib.CustomError(http.StatusBadRequest, "ops_profit_y3 can not be blank", "ops_profit_y3 can not be blank")
		}

		//BANK ACCOUNT
		if bankAccount == "" {
			log.Error("Missing required parameter: bank_account")
			return lib.CustomError(http.StatusBadRequest, "bank_account can not be blank", "bank_account can not be blank")
		} else {
			saveBank = true
			var bankAccountSlice []interface{}
			err = json.Unmarshal([]byte(bankAccount), &bankAccountSlice)
			if err != nil {
				log.Error(err.Error())
				log.Error("Missing required parameter: bank_account")
				return lib.CustomError(http.StatusBadRequest, err.Error(), "Wrong input for parameter: bank_account")
			}
			if len(bankAccountSlice) == 0 {
				log.Error("Missing required parameter: bank_account")
				return lib.CustomError(http.StatusBadRequest, "bank_account harus diisi minimal 1 data.", "bank_account harus diisi minimal 1 data.")
			}
			if len(bankAccountSlice) > 3 {
				log.Error("Missing required parameter: bank_account")
				return lib.CustomError(http.StatusBadRequest, "bank_account hanya max 3 data.", "bank_account hanya max 3 data.")
			}
			bankIsPriority := 0
			key := 1
			for _, val := range bankAccountSlice {
				bank := make(map[string]interface{})
				valueMap := val.(map[string]interface{})
				if val, ok := valueMap["req_bankacc_key"]; ok {
					if val.(string) != "" {
						n, err := strconv.ParseUint(val.(string), 10, 64)
						if err != nil || n == 0 {
							log.Error("Wrong input parameter: bank_account key : " + strconv.FormatUint(uint64(key), 10) + " req_bankacc_key.")
							return lib.CustomError(http.StatusBadRequest,
								"Wrong input parameter: bank_account key : "+strconv.FormatUint(uint64(key), 10)+" req_bankacc_key.",
								"Wrong input parameter: bank_account key : "+strconv.FormatUint(uint64(key), 10)+" req_bankacc_key.")
						}
					}
					bank["req_bankacc_key"] = val.(string)
				} else {
					log.Error("Missing required parameter: bank_account key : " + strconv.FormatUint(uint64(key), 10) + " req_bankacc_key tidak ditemukan")
					return lib.CustomError(http.StatusBadRequest,
						"Missing required parameter: bank_account key : "+strconv.FormatUint(uint64(key), 10)+" req_bankacc_key tidak ditemukan",
						"Missing required parameter: bank_account key : "+strconv.FormatUint(uint64(key), 10)+" req_bankacc_key tidak ditemukan")
				}
				if val, ok := valueMap["flag_priority"]; ok {
					if val.(string) == "" {
						log.Error("Missing required parameter: bank_account key : " + strconv.FormatUint(uint64(key), 10) + " flag_priority can not be blank")
						return lib.CustomError(http.StatusBadRequest,
							"Missing required parameter: bank_account key : "+strconv.FormatUint(uint64(key), 10)+" flag_priority can not be blank",
							"Missing required parameter: bank_account key : "+strconv.FormatUint(uint64(key), 10)+" flag_priority can not be blank")
					} else {
						bank["flag_priority"] = val.(string)
						if val.(string) == "1" {
							bankIsPriority++
						}
					}
				} else {
					log.Error("Missing required parameter: bank_account key : " + strconv.FormatUint(uint64(key), 10) + " flag_priority tidak ditemukan")
					return lib.CustomError(http.StatusBadRequest,
						"Missing required parameter: bank_account key : "+strconv.FormatUint(uint64(key), 10)+" flag_priority tidak ditemukan",
						"Missing required parameter: bank_account key : "+strconv.FormatUint(uint64(key), 10)+" flag_priority tidak ditemukan")
				}
				if val, ok := valueMap["bank_key"]; ok {
					if val.(string) == "" {
						log.Error("Missing required parameter: bank_account key : " + strconv.FormatUint(uint64(key), 10) + " bank_key can not be blank")
						return lib.CustomError(http.StatusBadRequest,
							"Missing required parameter: bank_account key : "+strconv.FormatUint(uint64(key), 10)+" bank_key can not be blank",
							"Missing required parameter: bank_account key : "+strconv.FormatUint(uint64(key), 10)+" bank_key can not be blank")
					} else {
						n, err := strconv.ParseUint(val.(string), 10, 64)
						if err != nil || n == 0 {
							log.Error("Wrong input for parameter: bank_account - bank_key")
							return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: bank_account - bank_key", "Wrong input for parameter: bank_account - bank_key")
						}
						if len(val.(string)) > 11 {
							log.Error("Wrong input for parameter: bank_account - bank_key too long")
							return lib.CustomError(http.StatusBadRequest, "Missing required parameter: bank_account - bank_key too long, max 11 character", "Missing required parameter: bank_account - bank_key too long, max 11 character")
						}
						bank["bank_key"] = val.(string)
					}
				} else {
					log.Error("Missing required parameter: bank_account key : " + strconv.FormatUint(uint64(key), 10) + " bank_key tidak ditemukan")
					return lib.CustomError(http.StatusBadRequest,
						"Missing required parameter: bank_account key : "+strconv.FormatUint(uint64(key), 10)+" bank_key tidak ditemukan",
						"Missing required parameter: bank_account key : "+strconv.FormatUint(uint64(key), 10)+" bank_key tidak ditemukan")
				}
				if val, ok := valueMap["account_no"]; ok {
					if val.(string) == "" {
						log.Error("Missing required parameter: bank_account key : " + strconv.FormatUint(uint64(key), 10) + " account_no can not be blank")
						return lib.CustomError(http.StatusBadRequest,
							"Missing required parameter: bank_account key : "+strconv.FormatUint(uint64(key), 10)+" account_no can not be blank",
							"Missing required parameter: bank_account key : "+strconv.FormatUint(uint64(key), 10)+" account_no can not be blank")
					} else {
						if len(val.(string)) > 30 {
							log.Error("Wrong input for parameter: bank_account - account_no too long")
							return lib.CustomError(http.StatusBadRequest,
								"Missing required parameter: bank_account - account_no too long, max 30 character",
								"Missing required parameter: bank_account - account_no too long, max 30 character")
						}
						bank["account_no"] = val.(string)
					}
				} else {
					log.Error("Missing required parameter: bank_account key : " + strconv.FormatUint(uint64(key), 10) + " account_no tidak ditemukan")
					return lib.CustomError(http.StatusBadRequest,
						"Missing required parameter: bank_account key : "+strconv.FormatUint(uint64(key), 10)+" account_no tidak ditemukan",
						"Missing required parameter: bank_account key : "+strconv.FormatUint(uint64(key), 10)+" account_no tidak ditemukan")
				}
				if val, ok := valueMap["account_holder_name"]; ok {
					if val.(string) == "" {
						log.Error("Missing required parameter: bank_account key : " + strconv.FormatUint(uint64(key), 10) + " account_holder_name can not be blank")
						return lib.CustomError(http.StatusBadRequest,
							"Missing required parameter: bank_account key : "+strconv.FormatUint(uint64(key), 10)+" account_holder_name can not be blank",
							"Missing required parameter: bank_account key : "+strconv.FormatUint(uint64(key), 10)+" account_holder_name can not be blank")
					} else {
						if len(val.(string)) > 80 {
							log.Error("Wrong input for parameter: bank_account - account_holder_name too long")
							return lib.CustomError(http.StatusBadRequest,
								"Missing required parameter: bank_account - account_holder_name too long, max 80 character",
								"Missing required parameter: bank_account - account_holder_name too long, max 80 character")
						}
						bank["account_holder_name"] = val.(string)
					}
				} else {
					log.Error("Missing required parameter: bank_account key : " + strconv.FormatUint(uint64(key), 10) + " account_holder_name tidak ditemukan")
					return lib.CustomError(http.StatusBadRequest,
						"Missing required parameter: bank_account key : "+strconv.FormatUint(uint64(key), 10)+" account_holder_name tidak ditemukan",
						"Missing required parameter: bank_account key : "+strconv.FormatUint(uint64(key), 10)+" account_holder_name tidak ditemukan")
				}
				if val, ok := valueMap["branch_name"]; ok {
					if val.(string) == "" {
						log.Error("Missing required parameter: bank_account key : " + strconv.FormatUint(uint64(key), 10) + " branch_name can not be blank")
						return lib.CustomError(http.StatusBadRequest,
							"Missing required parameter: bank_account key : "+strconv.FormatUint(uint64(key), 10)+" branch_name can not be blank",
							"Missing required parameter: bank_account key : "+strconv.FormatUint(uint64(key), 10)+" branch_name can not be blank")
					} else {
						if len(val.(string)) > 80 {
							log.Error("Wrong input for parameter: bank_account - branch_name too long")
							return lib.CustomError(http.StatusBadRequest,
								"Missing required parameter: bank_account - branch_name too long, max 80 character",
								"Missing required parameter: bank_account - branch_name too long, max 80 character")
						}
						bank["branch_name"] = val.(string)
					}
				} else {
					log.Error("Missing required parameter: bank_account key : " + strconv.FormatUint(uint64(key), 10) + " branch_name tidak ditemukan")
					return lib.CustomError(http.StatusBadRequest,
						"Missing required parameter: bank_account key : "+strconv.FormatUint(uint64(key), 10)+" branch_name tidak ditemukan",
						"Missing required parameter: bank_account key : "+strconv.FormatUint(uint64(key), 10)+" branch_name tidak ditemukan")
				}
				if val, ok := valueMap["currency_key"]; ok {
					if val.(string) == "" {
						log.Error("Missing required parameter: bank_account key : " + strconv.FormatUint(uint64(key), 10) + " currency_key can not be blank")
						return lib.CustomError(http.StatusBadRequest,
							"Missing required parameter: bank_account key : "+strconv.FormatUint(uint64(key), 10)+" currency_key can not be blank",
							"Missing required parameter: bank_account key : "+strconv.FormatUint(uint64(key), 10)+" currency_key can not be blank")
					} else {
						n, err := strconv.ParseUint(val.(string), 10, 64)
						if err != nil || n == 0 {
							log.Error("Wrong input for parameter: bank_account - currency_key")
							return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: bank_account - currency_key", "Wrong input for parameter: bank_account - currency_key")
						}
						if len(val.(string)) > 11 {
							log.Error("Wrong input for parameter: bank_account - currency_key too long")
							return lib.CustomError(http.StatusBadRequest, "Missing required parameter: bank_account - currency_key too long, max 11 character", "Missing required parameter: bank_account - currency_key too long, max 11 character")
						}
						bank["currency_key"] = val.(string)
					}
				} else {
					log.Error("Missing required parameter: bank_account key : " + strconv.FormatUint(uint64(key), 10) + " currency_key tidak ditemukan")
					return lib.CustomError(http.StatusBadRequest,
						"Missing required parameter: bank_account key : "+strconv.FormatUint(uint64(key), 10)+" currency_key tidak ditemukan",
						"Missing required parameter: bank_account key : "+strconv.FormatUint(uint64(key), 10)+" currency_key tidak ditemukan")
				}

				bankData = append(bankData, bank)
				key++
			}

			if bankIsPriority != 1 {
				log.Error("Missing required parameter: set 1 data bank_account to priority")
				return lib.CustomError(http.StatusBadRequest,
					"Missing required parameter: set 1 data bank_account to priority",
					"Missing required parameter: set 1 data bank_account to priority")
			}
		}
	}

	// SAVE POSTAL ADDRESS DOMICILE
	error := false
	domicileKey := domicilePostalAddressKey
	if domicilePostalAddressKey != "" {
		//update
		paramsPostalDomicile["rec_modified_date"] = time.Now().Format(layout)
		paramsPostalDomicile["rec_modified_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

		status, err := models.UpdateOaPostalAddress(paramsPostalDomicile)
		if err != nil {
			error = true
			log.Error("Failed update domicile adrress data: " + err.Error())
			return lib.CustomError(status, err.Error(), "failed update data domicile")
		} else {
			error = false
		}
	} else {
		//create
		if domicileAddress != "" || domicileSubdistric != "" || domicileCity != "" ||
			domicileProvince != "" || domicilePostalcode != "" {
			paramsPostalDomicile["rec_order"] = "0"
			paramsPostalDomicile["address_type"] = "18"
			paramsPostalDomicile["rec_status"] = "1"
			paramsPostalDomicile["rec_created_date"] = time.Now().Format(layout)
			paramsPostalDomicile["rec_created_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
			paramsPostalDomicile["rec_modified_date"] = time.Now().Format(layout)
			paramsPostalDomicile["rec_modified_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

			status, err, domicileID := models.CreateOaPostalAddress(paramsPostalDomicile)
			if err != nil {
				error = true
				log.Error("Failed create domicile adrress data: " + err.Error())
				return lib.CustomError(status, err.Error(), "failed input data domicile")
			} else {
				domicileKey = domicileID
				error = false
			}
		}
	}
	// END SAVE POSTAL ADDRESS DOMICILE

	// SAVE POSTAL ADDRESS CORRESPONDENCE
	correspondenceKey := correspondencePostalAddressKey
	if !error {
		if correspondencePostalAddressKey != "" {
			//update
			paramsPostalCorrespondence["rec_modified_date"] = time.Now().Format(layout)
			paramsPostalCorrespondence["rec_modified_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

			status, err := models.UpdateOaPostalAddress(paramsPostalCorrespondence)
			if err != nil {
				error = true
				log.Error("Failed update correspondenc adrress data: " + err.Error())
				return lib.CustomError(status, err.Error(), "failed update data correspondenc")
			} else {
				error = false
			}
		} else {
			//create

			if correspondenceAddress != "" || correspondenceSubdistric != "" || correspondenceCity != "" ||
				correspondenceProvince != "" || correspondencePostalcode != "" {

				paramsPostalCorrespondence["rec_order"] = "0"
				paramsPostalCorrespondence["address_type"] = "19"
				paramsPostalCorrespondence["rec_status"] = "1"
				paramsPostalCorrespondence["rec_created_date"] = time.Now().Format(layout)
				paramsPostalCorrespondence["rec_created_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
				paramsPostalCorrespondence["rec_modified_date"] = time.Now().Format(layout)
				paramsPostalCorrespondence["rec_modified_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

				status, err, correspondenceID := models.CreateOaPostalAddress(paramsPostalCorrespondence)
				if err != nil {
					error = true
					log.Error("Failed create correspondenc adrress data: " + err.Error())
					return lib.CustomError(status, err.Error(), "failed input data correspondenc")
				} else {
					correspondenceKey = correspondenceID
					error = false
				}
			}
		}
	}
	// END SAVE POSTAL ADDRESS CORRESPONDENCE

	// SAVE OA REQUEST
	requestKey := oaRequestKey
	if !error {
		var agent models.MsAgent
		status, err = models.GetMsAgent(&agent, agentKey)
		if err == nil {
			paramsOaRequest["sales_code"] = agent.AgentCode
		}

		if oaRequestKey != "" {
			//update
			paramsOaRequest["oa_entry_end"] = time.Now().Format(layout)
			paramsOaRequest["oa_status"] = strconv.FormatUint(uint64(lib.DRAFT), 10)
			paramsOaRequest["rec_modified_date"] = time.Now().Format(layout)
			paramsOaRequest["rec_modified_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
			status, err := models.UpdateOaRequest(paramsOaRequest)
			if err != nil {
				error = true
				log.Error("Failed create request oa: " + err.Error())
				return lib.CustomError(status, err.Error(), "Failed create request oa")
			} else {
				error = false
			}
		} else {
			//create
			paramsOaRequest["oa_request_type"] = lib.OA_REQ_TYPE_NEW
			paramsOaRequest["oa_entry_start"] = time.Now().Format(layout)
			paramsOaRequest["oa_entry_end"] = time.Now().Format(layout)
			paramsOaRequest["oa_status"] = strconv.FormatUint(uint64(lib.DRAFT), 10)
			paramsOaRequest["rec_order"] = "0"
			paramsOaRequest["rec_status"] = "1"
			paramsOaRequest["rec_created_date"] = time.Now().Format(layout)
			paramsOaRequest["rec_created_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
			paramsOaRequest["rec_modified_date"] = time.Now().Format(layout)
			paramsOaRequest["rec_modified_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
			status, err, requestID := models.CreateOaRequest(paramsOaRequest)
			if err != nil {
				error = true
				log.Error("Failed create request oa: " + err.Error())
				return lib.CustomError(status, err.Error(), "Failed create request oa")
			} else {
				requestKey = requestID
				error = false
			}
		}
	}
	// END SAVE OA REQUEST

	// SAVE OA_INSTITUTION_DATA
	if !error {
		if oaRequestKey != "" {
			//update
			var oadata models.OaInstitutionData
			_, err = models.GetOaInstitutionData(&oadata, requestKey, "oa_request_key")
			if err != nil {
				error = true
				log.Error("Failed get oa institution data: " + err.Error())
				return lib.CustomError(status, err.Error(), "Failed get oa institution data")
			}

			if domicileKey != "" {
				paramsInstitutionData["domicile_key"] = domicileKey
			}
			if correspondenceKey != "" {
				paramsInstitutionData["correspondence_key"] = correspondenceKey
			}
			paramsInstitutionData["institution_data_key"] = strconv.FormatUint(oadata.InstitutionDataKey, 10)
			paramsInstitutionData["rec_modified_date"] = time.Now().Format(layout)
			paramsInstitutionData["rec_modified_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
			status, err = models.UpdateOaInstitutionData(paramsInstitutionData)
			if err != nil {
				error = true
				log.Error("Failed update oa institution data: " + err.Error())
				return lib.CustomError(status, err.Error(), "Failed update oa institution data")
			} else {
				error = false
			}
		} else {
			//create
			paramsInstitutionData["oa_request_key"] = requestKey
			if domicileKey != "" {
				paramsInstitutionData["domicile_key"] = domicileKey
			}
			if correspondenceKey != "" {
				paramsInstitutionData["correspondence_key"] = correspondenceKey
			}

			paramsInstitutionData["rec_order"] = "0"
			paramsInstitutionData["rec_status"] = "1"
			paramsInstitutionData["rec_created_date"] = time.Now().Format(layout)
			paramsInstitutionData["rec_created_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
			paramsInstitutionData["rec_modified_date"] = time.Now().Format(layout)
			paramsInstitutionData["rec_modified_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
			status, err, _ := models.CreateOaInstitutionData(paramsInstitutionData)
			if err != nil {
				error = true
				log.Error("Failed create oa institution data: " + err.Error())
				return lib.CustomError(status, err.Error(), "Failed create oa institution data")
			} else {
				error = false
			}
		}
	}
	// END SAVE OA_INSTITUTION_DATA

	// SAVE MS_BANK_ACCOUNT & OA_REQUEST_BANK_ACCOUNT
	var bankReqKeyNotDelete []string
	if !error {
		if saveBank {
			if len(bankData) > 0 {
				for _, bd := range bankData {
					valueMap := bd.(map[string]interface{})
					if valueMap["req_bankacc_key"].(string) != "" {
						bankReqKeyNotDelete = append(bankReqKeyNotDelete, valueMap["req_bankacc_key"].(string))
						//update
						var bankAccReq models.OaRequestBankAccount
						_, err = models.GetOaRequestBankAccount(&bankAccReq, valueMap["req_bankacc_key"].(string), "req_bankacc_key")
						if err == nil {
							if strconv.FormatUint(bankAccReq.OaRequestKey, 10) == requestKey {
								paramsBank := make(map[string]string)
								paramsBank["bank_account_key"] = strconv.FormatUint(bankAccReq.BankAccountKey, 10)
								paramsBank["bank_key"] = valueMap["bank_key"].(string)
								paramsBank["account_no"] = valueMap["account_no"].(string)
								paramsBank["account_holder_name"] = valueMap["account_holder_name"].(string)
								paramsBank["branch_name"] = valueMap["branch_name"].(string)
								if valueMap["currency_key"].(string) != "" {
									paramsBank["currency_key"] = valueMap["currency_key"].(string)
								}
								paramsBank["rec_modified_date"] = time.Now().Format(layout)
								paramsBank["rec_modified_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
								status, err := models.UpdateMsBankAccount(paramsBank)
								if err != nil {
									error = true
									log.Error("Failed update bank account: " + err.Error())
									return lib.CustomError(status, err.Error(), "Failed update bank account")
								} else {
									error = false
									paramsBankReq := make(map[string]string)
									paramsBankReq["req_bankacc_key"] = valueMap["req_bankacc_key"].(string)
									if valueMap["flag_priority"].(string) == "1" {
										paramsBankReq["flag_priority"] = valueMap["flag_priority"].(string)
									} else {
										paramsBankReq["flag_priority"] = "0"
									}
									paramsBankReq["bank_account_name"] = valueMap["account_holder_name"].(string)
									paramsBankReq["rec_modified_date"] = time.Now().Format(layout)
									paramsBankReq["rec_modified_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
									status, err = models.UpdateOaRequestBankAccount(paramsBankReq)
									if err != nil {
										error = true
										log.Error("Failed update oa bank request: " + err.Error())
										return lib.CustomError(status, err.Error(), "Failed update oa bank request")
									} else {
										error = false
									}
								}

							} else {
								log.Error("Failed update oa_request_bank_account : oa_request_key != oa_request_key bank request")
							}
						} else {
							log.Error("Failed get oa_request_bank_account key: " + valueMap["req_bankacc_key"].(string))
						}
					} else {
						//create
						paramsBank := make(map[string]string)
						paramsBank["bank_key"] = valueMap["bank_key"].(string)
						paramsBank["account_no"] = valueMap["account_no"].(string)
						paramsBank["account_holder_name"] = valueMap["account_holder_name"].(string)
						paramsBank["branch_name"] = valueMap["branch_name"].(string)
						if valueMap["currency_key"].(string) != "" {
							paramsBank["currency_key"] = valueMap["currency_key"].(string)
						}
						paramsBank["bank_account_type"] = "129"
						paramsBank["rec_domain"] = "131"
						paramsBank["rec_order"] = "0"
						paramsBank["rec_status"] = "1"
						paramsBank["rec_created_date"] = time.Now().Format(layout)
						paramsBank["rec_created_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
						paramsBank["rec_modified_date"] = time.Now().Format(layout)
						paramsBank["rec_modified_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
						status, err, bankAccountID := models.CreateMsBankAccount(paramsBank)
						if err != nil {
							error = true
							log.Error("Failed create bank account: " + err.Error())
							return lib.CustomError(status, err.Error(), "Failed create bank account")
						} else {
							error = false
							paramsBankReq := make(map[string]string)
							paramsBankReq["oa_request_key"] = requestKey
							paramsBankReq["bank_account_key"] = bankAccountID
							if valueMap["flag_priority"].(string) == "1" {
								paramsBankReq["flag_priority"] = valueMap["flag_priority"].(string)
							} else {
								paramsBankReq["flag_priority"] = "0"
							}
							paramsBankReq["bank_account_name"] = valueMap["account_holder_name"].(string)
							paramsBankReq["rec_order"] = "0"
							paramsBankReq["rec_status"] = "1"
							paramsBankReq["rec_created_date"] = time.Now().Format(layout)
							paramsBankReq["rec_created_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
							paramsBankReq["rec_modified_date"] = time.Now().Format(layout)
							paramsBankReq["rec_modified_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
							status, err, bankReqID := models.CreateOaRequestBankAccount(paramsBankReq)
							if err != nil {
								error = true
								log.Error("Failed create oa bank request: " + err.Error())
								return lib.CustomError(status, err.Error(), "Failed create oa bank request")
							} else {
								bankReqKeyNotDelete = append(bankReqKeyNotDelete, bankReqID)
								error = false
							}
						}
					}
				}
			}
		}
	}
	// END SAVE MS_BANK_ACCOUNT & OA_REQUEST_BANK_ACCOUNT

	// DELETE OA_REQUEST_BANK_ACCOUNT
	paramsDeleteBankReq := make(map[string]string)
	paramsDeleteBankReq["rec_status"] = "0"
	paramsDeleteBankReq["rec_deleted_date"] = time.Now().Format(layout)
	paramsDeleteBankReq["rec_deleted_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
	status, err = models.DeleteOaRequestBankAccount(paramsDeleteBankReq, bankReqKeyNotDelete, requestKey)
	if err != nil {
		log.Error("Failed delete oa bank request: " + err.Error())
	}
	// END DELETE OA_REQUEST_BANK_ACCOUNT

	// SAVE SHARES_HOLDER
	var sharesHolderKeyNotDelete []string
	if !error {
		if saveSharesHolder {
			if len(sharesHolderData) > 0 {
				for _, shr := range sharesHolderData {
					valueMap := shr.(map[string]interface{})
					if valueMap["inst_shares_holder_key"].(string) != "" {
						sharesHolderKeyNotDelete = append(sharesHolderKeyNotDelete, valueMap["inst_shares_holder_key"].(string))
						//update
						paramsSharesHolder := make(map[string]string)
						paramsSharesHolder["inst_shares_holder_key"] = valueMap["inst_shares_holder_key"].(string)
						paramsSharesHolder["holder_full_name"] = valueMap["holder_full_name"].(string)
						paramsSharesHolder["nationality"] = valueMap["nationality"].(string)
						paramsSharesHolder["idcard_type"] = valueMap["idcard_type"].(string)
						paramsSharesHolder["idcard_no"] = valueMap["idcard_no"].(string)
						paramsSharesHolder["holder_dob"] = valueMap["holder_dob"].(string)
						paramsSharesHolder["shares_percent"] = valueMap["shares_percent"].(string)
						paramsSharesHolder["rec_modified_date"] = time.Now().Format(layout)
						paramsSharesHolder["rec_modified_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
						status, err = models.UpdateOaInstitutionSharesHolder(paramsSharesHolder)
						if err != nil {
							error = true
							log.Error("Failed update shares_holder: " + err.Error())
							return lib.CustomError(status, err.Error(), "Failed update shares_holder")
						} else {
							error = false
						}
					} else {
						//create
						paramsSharesHolder := make(map[string]string)
						paramsSharesHolder["oa_request_key"] = requestKey
						paramsSharesHolder["holder_full_name"] = valueMap["holder_full_name"].(string)
						paramsSharesHolder["nationality"] = valueMap["nationality"].(string)
						paramsSharesHolder["idcard_type"] = valueMap["idcard_type"].(string)
						paramsSharesHolder["idcard_no"] = valueMap["idcard_no"].(string)
						paramsSharesHolder["idcard_never_expired"] = "1"
						paramsSharesHolder["holder_dob"] = valueMap["holder_dob"].(string)
						paramsSharesHolder["shares_percent"] = valueMap["shares_percent"].(string)
						paramsSharesHolder["rec_order"] = "0"
						paramsSharesHolder["rec_status"] = "1"
						paramsSharesHolder["rec_created_date"] = time.Now().Format(layout)
						paramsSharesHolder["rec_created_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
						paramsSharesHolder["rec_modified_date"] = time.Now().Format(layout)
						paramsSharesHolder["rec_modified_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
						status, err, shareHolderID := models.CreateOaInstitutionSharesHolder(paramsSharesHolder)
						if err != nil {
							error = true
							log.Error("Failed create shares_holder: " + err.Error())
							return lib.CustomError(status, err.Error(), "Failed create shares_holder")
						} else {
							sharesHolderKeyNotDelete = append(sharesHolderKeyNotDelete, shareHolderID)
							error = false
						}
					}
				}
			}
		}
	}
	// END SAVE SHARES_HOLDER

	// DELETE SHARES_HOLDER
	paramsDeleteShareHolder := make(map[string]string)
	paramsDeleteShareHolder["rec_status"] = "0"
	paramsDeleteShareHolder["rec_deleted_date"] = time.Now().Format(layout)
	paramsDeleteShareHolder["rec_deleted_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
	status, err = models.DeleteOaInstitutionSharesHolder(paramsDeleteShareHolder, sharesHolderKeyNotDelete, requestKey)
	if err != nil {
		log.Error("Failed delete oa shares holder: " + err.Error())
	}
	// END DELETE SHARES_HOLDER

	// SAVE AUTH_USER
	var authPersonKeyNotDelete []string
	if !error {
		if saveAuthPerson {
			if len(authPersonData) > 0 {
				for _, shr := range authPersonData {
					valueMap := shr.(map[string]interface{})
					if valueMap["insti_auth_person_key"].(string) != "" {
						authPersonKeyNotDelete = append(authPersonKeyNotDelete, valueMap["insti_auth_person_key"].(string))
						//update
						paramsAuthPerson := make(map[string]string)
						paramsAuthPerson["insti_auth_person_key"] = valueMap["insti_auth_person_key"].(string)
						paramsAuthPerson["full_name"] = valueMap["full_name"].(string)
						paramsAuthPerson["nationality"] = valueMap["nationality"].(string)
						paramsAuthPerson["idcard_type"] = valueMap["idcard_type"].(string)
						paramsAuthPerson["idcard_no"] = valueMap["idcard_no"].(string)
						paramsAuthPerson["idcard_never_expired"] = "1"
						paramsAuthPerson["person_dob"] = valueMap["person_dob"].(string)
						paramsAuthPerson["position"] = valueMap["position"].(string)
						paramsAuthPerson["phone_no"] = valueMap["phone_no"].(string)
						paramsAuthPerson["email_address"] = valueMap["email_address"].(string)
						paramsAuthPerson["rec_modified_date"] = time.Now().Format(layout)
						paramsAuthPerson["rec_modified_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
						status, err = models.UpdateOaInstitutionAuthPerson(paramsAuthPerson)
						if err != nil {
							error = true
							log.Error("Failed update auth_person: " + err.Error())
							return lib.CustomError(status, err.Error(), "Failed update auth_person")
						} else {
							error = false
						}
					} else {
						//create
						paramsAuthPerson := make(map[string]string)
						paramsAuthPerson["oa_request_key"] = requestKey
						paramsAuthPerson["full_name"] = valueMap["full_name"].(string)
						paramsAuthPerson["nationality"] = valueMap["nationality"].(string)
						paramsAuthPerson["idcard_type"] = valueMap["idcard_type"].(string)
						paramsAuthPerson["idcard_no"] = valueMap["idcard_no"].(string)
						paramsAuthPerson["idcard_never_expired"] = "1"
						paramsAuthPerson["person_dob"] = valueMap["person_dob"].(string)
						paramsAuthPerson["position"] = valueMap["position"].(string)
						paramsAuthPerson["phone_no"] = valueMap["phone_no"].(string)
						paramsAuthPerson["email_address"] = valueMap["email_address"].(string)
						paramsAuthPerson["rec_order"] = "0"
						paramsAuthPerson["rec_status"] = "1"
						paramsAuthPerson["rec_created_date"] = time.Now().Format(layout)
						paramsAuthPerson["rec_created_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
						paramsAuthPerson["rec_modified_date"] = time.Now().Format(layout)
						paramsAuthPerson["rec_modified_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
						status, err, authPersonID := models.CreateOaInstitutionAuthPerson(paramsAuthPerson)
						if err != nil {
							error = true
							log.Error("Failed create auth_person: " + err.Error())
							return lib.CustomError(status, err.Error(), "Failed create auth_person")
						} else {
							authPersonKeyNotDelete = append(authPersonKeyNotDelete, authPersonID)
							error = false
						}
					}
				}
			}
		}
	}
	// END SAVE AUTH_USER

	// DELETE AUTH_PERSON
	paramsDeleteAuthPerson := make(map[string]string)
	paramsDeleteAuthPerson["rec_status"] = "0"
	paramsDeleteAuthPerson["rec_deleted_date"] = time.Now().Format(layout)
	paramsDeleteAuthPerson["rec_deleted_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
	status, err = models.DeleteOaInstitutionAuthPerson(paramsDeleteAuthPerson, authPersonKeyNotDelete, requestKey)
	if err != nil {
		log.Error("Failed delete oa shares holder: " + err.Error())
	}
	// END DELETE AUTH_PERSON

	responseData := make(map[string]interface{})
	responseData["oa_request_key"] = requestKey

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = responseData

	return c.JSON(http.StatusOK, response)
}

func SaveDocsInstitution(c echo.Context) error {
	var err error
	// var status int
	decimal.MarshalJSONWithoutQuotes = true
	paramsFile := make(map[string]string)
	paramsDocs := make(map[string]string)

	oaRequestKey := c.FormValue("oa_request_key")
	if oaRequestKey != "" {
		n, err := strconv.ParseUint(oaRequestKey, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: oa_request_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: oa_request_key", "Wrong input for parameter: oa_request_key")
		}

		if len(oaRequestKey) > 11 {
			log.Error("Wrong input for parameter: oa_request_key too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: oa_request_key too long, max 11 character", "Missing required parameter: oa_request_key too long, max 11 character")
		}
		paramsDocs["oa_request_key"] = oaRequestKey
	} else {
		log.Error("Missing required parameter: oa_request_key")
		return lib.CustomError(http.StatusBadRequest, "oa_request_key can not be blank", "oa_request_key can not be blank")

	}

	if oaRequestKey != "" {
		var oareq models.OaRequest
		_, err = models.GetOaRequestInstitution(&oareq, oaRequestKey, "")
		if err != nil {
			log.Error("OA Request not found.")
			return lib.CustomError(http.StatusBadRequest, "OA Request not found.", "OA Request not found.")
		}

		var userCategory uint64
		userCategory = 3 //Branch
		if lib.Profile.UserCategoryKey == userCategory {
			if oareq.BranchKey != lib.Profile.BranchKey {
				log.Error("User not autorized.")
				return lib.CustomError(http.StatusBadRequest, "User not autorized.", "User not autorized.")
			}
		}

		if strconv.FormatUint(lib.Profile.UserID, 10) != *oareq.RecCreatedBy {
			log.Error("User not autorized.")
			return lib.CustomError(http.StatusBadRequest, "User not autorized.", "User not autorized.")
		}

		if *oareq.Oastatus != uint64(lib.DRAFT) {
			log.Error("User not autorized.")
			return lib.CustomError(http.StatusBadRequest, "User not autorized.", "User not autorized.")
		}
	}

	instiDocsKey := c.FormValue("insti_docs_key")
	msFileKey := ""
	if instiDocsKey != "" {
		n, err := strconv.ParseUint(instiDocsKey, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: insti_docs_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: insti_docs_key", "Wrong input for parameter: insti_docs_key")
		}
		var insDocs models.OaInstitutionDocs
		_, err = models.GetOaInstitutionDocs(&insDocs, instiDocsKey, "insti_docs_key")
		if err != nil {
			log.Error("Institution Docs not found.")
			return lib.CustomError(http.StatusBadRequest, "Institution Docs not found.", "Institution Docs not found.")
		}

		if strconv.FormatUint(insDocs.OaRequestKey, 10) != oaRequestKey {
			log.Error("User not autorized. OA request berbeda dengan oa_request_key di oa_institution_docs")
			return lib.CustomError(http.StatusBadRequest, "User not autorized.", "User not autorized.")
		}

		paramsDocs["insti_docs_key"] = instiDocsKey
		if insDocs.DocumentFileKey != nil {
			msFileKey = strconv.FormatUint(*insDocs.DocumentFileKey, 10)
		}
	}

	instiDocumentType := c.FormValue("insti_document_type")
	if instiDocumentType == "" {
		log.Error("Missing required parameter: insti_document_type")
		return lib.CustomError(http.StatusBadRequest, "insti_document_type can not be blank", "insti_document_type can not be blank")
	} else {
		n, err := strconv.ParseUint(instiDocumentType, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: insti_document_type")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: insti_document_type", "Wrong input for parameter: insti_document_type")
		}
		paramsDocs["insti_document_type"] = instiDocumentType
	}

	remark := c.FormValue("insti_document_remarks")
	if remark != "" {
		paramsDocs["insti_document_remarks"] = remark
		paramsFile["file_notes"] = remark
	}

	var file *multipart.FileHeader
	filename := ""
	originalFileName := ""
	file, err = c.FormFile("file")
	if file != nil {
		if err != nil {
			return lib.CustomError(http.StatusBadRequest)
		}
		err = os.MkdirAll(config.BasePath+lib.INST_FILE_PATH, 0755)
		if err != nil {
			log.Error(err.Error())
		}
		// Get file extension
		extension := filepath.Ext(file.Filename)
		format := []string{".jpg", ".JPG", ".png", ".PNG", ".jpeg", ".JPEG", ".pdf", ".PDF"}

		_, found := lib.Find(format, extension)
		if !found {
			log.Error("Wrong input for parameter: file. hanya format jpg/png/jpeg/pdf yang diizinkan.")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter:  file. hanya format jpg/png/jpeg/pdf yang diizinkan", "Wrong input for parameter:  file. hanya format jpg/png/jpeg/pdf yang diizinkan")
		}

		if file.Size > int64(lib.MAX_FILE_SIZE) {
			log.Error("Wrong input for parameter: file size max " + lib.MAX_FILE_SIZE_MB + " MB.")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: file size max "+lib.MAX_FILE_SIZE_MB+" MB.", "Wrong input for parameter: file size max "+lib.MAX_FILE_SIZE_MB+" MB.")
		}

		// Generate filename
		originalFileName = file.Filename
		filename = lib.RandStringBytesMaskImprSrc(8) + "-" + strings.ReplaceAll(file.Filename, " ", "-")
		// Upload image and move to proper directory
		err = lib.UploadImage(file, config.BasePath+lib.INST_FILE_PATH+"/"+filename)
		if err != nil {
			log.Println(err)
			return lib.CustomError(http.StatusInternalServerError)
		}
		paramsFile["file_name"] = filename
		paramsFile["file_ext"] = extension
		paramsFile["file_path"] = config.BasePath + lib.INST_FILE_PATH + "/" + filename

		paramsDocs["insti_document_name"] = originalFileName
	} else {
		log.Error("Missing required parameter: file")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: file", "Missing required parameter: file")
	}

	layout := "2006-01-02 15:04:05"
	docsID := instiDocsKey
	if instiDocsKey == "" {
		//create new ms_file
		paramsFile["ref_fk_domain"] = "oa_institution_docs"
		paramsFile["blob_mode"] = "0"
		paramsFile["rec_status"] = "1"
		paramsFile["rec_created_date"] = time.Now().Format(layout)
		paramsFile["rec_created_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
		paramsFile["rec_modified_date"] = time.Now().Format(layout)
		paramsFile["rec_modified_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
		_, err, fileID := models.CreateMsFile(paramsFile)
		if err != nil {
			log.Error("Error create ms_file")
			return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed create File")
		} else {
			//create new oa_institution_docs
			paramsDocs["rec_order"] = "0"
			paramsDocs["rec_status"] = "1"
			paramsDocs["document_file_key"] = fileID
			paramsDocs["rec_created_date"] = time.Now().Format(layout)
			paramsDocs["rec_created_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
			paramsDocs["rec_modified_date"] = time.Now().Format(layout)
			paramsDocs["rec_modified_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
			_, err, docsID = models.CreateOaInstitutionDocs(paramsDocs)
			if err != nil {
				log.Error("Error create oa_institution_docs")
				return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed create Institution Docs")
			}
		}
	} else {
		if msFileKey == "" {
			paramsFile["ref_fk_domain"] = "oa_institution_docs"
			paramsFile["blob_mode"] = "0"
			paramsFile["rec_status"] = "1"
			paramsFile["rec_created_date"] = time.Now().Format(layout)
			paramsFile["rec_created_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
			paramsFile["rec_modified_date"] = time.Now().Format(layout)
			paramsFile["rec_modified_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
			_, err, fileID := models.CreateMsFile(paramsFile)
			if err != nil {
				log.Error("Error create ms_file")
				return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed create File")
			} else {
				msFileKey = fileID
			}

		} else {
			paramsFile["file_key"] = msFileKey
			paramsFile["rec_modified_date"] = time.Now().Format(layout)
			paramsFile["rec_modified_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
			_, err = models.UpdateMsFile(paramsFile)
			if err != nil {
				log.Error("Error update ms_file")
				return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed update File")
			}
		}

		//update new oa_institution_docs
		paramsDocs["document_file_key"] = msFileKey
		paramsDocs["rec_modified_date"] = time.Now().Format(layout)
		paramsDocs["rec_modified_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
		_, err = models.UpdateOaInstitutionDocs(paramsDocs)
		if err != nil {
			log.Error("Error update oa_institution_docs")
			return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed update Institution Docs")
		}
	}

	responseData := make(map[string]interface{})
	responseData["oa_request_key"] = oaRequestKey
	responseData["insti_docs_key"] = docsID
	responseData["document_file_name"] = filename
	responseData["insti_document_name"] = originalFileName
	responseData["insti_document_remarks"] = remark
	responseData["path"] = config.BaseUrl + lib.INST_FILE_PATH + "/" + filename

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = responseData

	return c.JSON(http.StatusOK, response)
}

func DeleteDocsInstitution(c echo.Context) error {
	var err error
	// var status int
	decimal.MarshalJSONWithoutQuotes = true
	paramsDocs := make(map[string]string)

	oaRequestKey := c.FormValue("oa_request_key")
	if oaRequestKey != "" {
		n, err := strconv.ParseUint(oaRequestKey, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: oa_request_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: oa_request_key", "Wrong input for parameter: oa_request_key")
		}

		if len(oaRequestKey) > 11 {
			log.Error("Wrong input for parameter: oa_request_key too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: oa_request_key too long, max 11 character", "Missing required parameter: oa_request_key too long, max 11 character")
		}
	} else {
		log.Error("Missing required parameter: oa_request_key")
		return lib.CustomError(http.StatusBadRequest, "oa_request_key can not be blank", "oa_request_key can not be blank")

	}

	if oaRequestKey != "" {
		var oareq models.OaRequest
		_, err = models.GetOaRequestInstitution(&oareq, oaRequestKey, "")
		if err != nil {
			log.Error("OA Request not found.")
			return lib.CustomError(http.StatusBadRequest, "OA Request not found.", "OA Request not found.")
		}

		var userCategory uint64
		userCategory = 3 //Branch
		if lib.Profile.UserCategoryKey == userCategory {
			if oareq.BranchKey != lib.Profile.BranchKey {
				log.Error("User not autorized.")
				return lib.CustomError(http.StatusBadRequest, "User not autorized.", "User not autorized.")
			}
		}

		if strconv.FormatUint(lib.Profile.UserID, 10) != *oareq.RecCreatedBy {
			log.Error("User not autorized.")
			return lib.CustomError(http.StatusBadRequest, "User not autorized.", "User not autorized.")
		}

		if *oareq.Oastatus != uint64(lib.DRAFT) {
			log.Error("User not autorized.")
			return lib.CustomError(http.StatusBadRequest, "User not autorized.", "User not autorized.")
		}
	}

	instiDocsKey := c.FormValue("insti_docs_key")
	if instiDocsKey != "" {
		n, err := strconv.ParseUint(instiDocsKey, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: insti_docs_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: insti_docs_key", "Wrong input for parameter: insti_docs_key")
		}
		var insDocs models.OaInstitutionDocs
		_, err = models.GetOaInstitutionDocs(&insDocs, instiDocsKey, "insti_docs_key")
		if err != nil {
			log.Error("Institution Docs not found.")
			return lib.CustomError(http.StatusBadRequest, "Institution Docs not found.", "Institution Docs not found.")
		}

		if strconv.FormatUint(insDocs.OaRequestKey, 10) != oaRequestKey {
			log.Error("User not autorized. OA request berbeda dengan oa_request_key di oa_institution_docs")
			return lib.CustomError(http.StatusBadRequest, "User not autorized.", "User not autorized.")
		}

		paramsDocs["insti_docs_key"] = instiDocsKey
	} else {
		log.Error("Missing required parameter: insti_docs_key")
		return lib.CustomError(http.StatusBadRequest, "insti_docs_key can not be blank", "insti_docs_key can not be blank")
	}

	layout := "2006-01-02 15:04:05"

	//update new oa_institution_docs
	paramsDocs["rec_deleted_date"] = time.Now().Format(layout)
	paramsDocs["rec_deleted_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
	paramsDocs["rec_status"] = "0"
	_, err = models.UpdateOaInstitutionDocs(paramsDocs)
	if err != nil {
		log.Error("Error delete oa_institution_docs")
		return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed delete Institution Docs")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = ""

	return c.JSON(http.StatusOK, response)
}

func SaveRiskProfileInstitution(c echo.Context) error {
	var err error
	var status int
	decimal.MarshalJSONWithoutQuotes = true

	oaRequestKey := c.FormValue("oa_request_key")
	if oaRequestKey != "" {
		n, err := strconv.ParseUint(oaRequestKey, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: oa_request_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: oa_request_key", "Wrong input for parameter: oa_request_key")
		}

		if len(oaRequestKey) > 11 {
			log.Error("Wrong input for parameter: oa_request_key too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: oa_request_key too long, max 11 character", "Missing required parameter: oa_request_key too long, max 11 character")
		}
	} else {
		log.Error("Missing required parameter: oa_request_key")
		return lib.CustomError(http.StatusBadRequest, "oa_request_key can not be blank", "oa_request_key can not be blank")
	}

	if oaRequestKey != "" {
		var oareq models.OaRequest
		_, err = models.GetOaRequestInstitution(&oareq, oaRequestKey, "")
		if err != nil {
			log.Error("OA Request not found.")
			return lib.CustomError(http.StatusBadRequest, "OA Request not found.", "OA Request not found.")
		}

		var userCategory uint64
		userCategory = 3 //Branch
		if lib.Profile.UserCategoryKey == userCategory {
			if oareq.BranchKey != lib.Profile.BranchKey {
				log.Error("User not autorized.")
				return lib.CustomError(http.StatusBadRequest, "User not autorized.", "User not autorized.")
			}
		}

		if strconv.FormatUint(lib.Profile.UserID, 10) != *oareq.RecCreatedBy {
			log.Error("User not autorized.")
			return lib.CustomError(http.StatusBadRequest, "User not autorized.", "User not autorized.")
		}

		if *oareq.Oastatus != uint64(lib.DRAFT) {
			log.Error("User not autorized.")
			return lib.CustomError(http.StatusBadRequest, "User not autorized.", "User not autorized.")
		}
	}

	isSaveDraft := c.FormValue("is_save_draft")
	if isSaveDraft != "" {
		if isSaveDraft != "0" && isSaveDraft != "1" {
			log.Error("Wrong input for parameter: is_save_draft")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: is_save_draft", "Wrong input for parameter: is_save_draft")
		}
	} else {
		log.Error("Missing required parameter: is_save_draft")
		return lib.CustomError(http.StatusBadRequest, "is_save_draft can not be blank", "is_save_draft can not be blank")
	}

	riskProfileQuiz := c.FormValue("risk_profile_quiz")

	var riskProfileData []interface{}
	if riskProfileQuiz != "" {
		var riskProfileAccountSlice []interface{}
		err = json.Unmarshal([]byte(riskProfileQuiz), &riskProfileAccountSlice)
		if err != nil {
			log.Error(err.Error())
			log.Error("Missing required parameter: risk_profile_quiz")
			return lib.CustomError(http.StatusBadRequest, err.Error(), "Wrong input for parameter: risk_profile_quiz")
		}
		if len(riskProfileAccountSlice) == 0 {
			log.Error("Missing required parameter: risk_profile_quiz")
			return lib.CustomError(http.StatusBadRequest, err.Error(), "Wrong input for parameter: risk_profile_quiz")
		}
		if len(riskProfileAccountSlice) > 6 {
			log.Error("Missing required parameter: risk_profile_quiz")
			return lib.CustomError(http.StatusBadRequest, "risk_profile_quiz hanya max 6 data.", "risk_profile_quiz hanya max 6 data.")
		}
		if isSaveDraft == "0" {
			if len(riskProfileAccountSlice) != 6 {
				log.Error("Missing required parameter: risk_profile_quiz belum diisi semua")
				return lib.CustomError(http.StatusBadRequest, "risk_profile_quiz harus diisi semua.", "risk_profile_quiz harus diisi semua.")
			}
		}
		key := 1
		if len(riskProfileAccountSlice) > 0 {
			for _, val := range riskProfileAccountSlice {
				risk := make(map[string]interface{})
				valueMap := val.(map[string]interface{})
				if val, ok := valueMap["risk_profile_quiz_key"]; ok {
					if val.(string) != "" {
						n, err := strconv.ParseUint(val.(string), 10, 64)
						if err != nil || n == 0 {
							log.Error("Wrong input parameter: risk_profile_quiz key : " + strconv.FormatUint(uint64(key), 10) + " risk_profile_quiz_key.")
							return lib.CustomError(http.StatusBadRequest,
								"Wrong input parameter: risk_profile_quiz key : "+strconv.FormatUint(uint64(key), 10)+" risk_profile_quiz_key.",
								"Wrong input parameter: risk_profile_quiz key : "+strconv.FormatUint(uint64(key), 10)+" risk_profile_quiz_key.")
						}
					}
					risk["risk_profile_quiz_key"] = val.(string)
				} else {
					log.Error("Missing required parameter: risk_profile_quiz key : " + strconv.FormatUint(uint64(key), 10) + " risk_profile_quiz_key tidak ditemukan")
					return lib.CustomError(http.StatusBadRequest,
						"Missing required parameter: risk_profile_quiz key : "+strconv.FormatUint(uint64(key), 10)+" risk_profile_quiz_key tidak ditemukan",
						"Missing required parameter: risk_profile_quiz key : "+strconv.FormatUint(uint64(key), 10)+" risk_profile_quiz_key tidak ditemukan")
				}

				if val, ok := valueMap["quiz_question_key"]; ok {
					if val.(string) != "" {
						n, err := strconv.ParseUint(val.(string), 10, 64)
						if err != nil || n == 0 {
							log.Error("Wrong input parameter: risk_profile_quiz key : " + strconv.FormatUint(uint64(key), 10) + " quiz_question_key.")
							return lib.CustomError(http.StatusBadRequest,
								"Wrong input parameter: risk_profile_quiz key : "+strconv.FormatUint(uint64(key), 10)+" quiz_question_key.",
								"Wrong input parameter: risk_profile_quiz key : "+strconv.FormatUint(uint64(key), 10)+" quiz_question_key.")
						}
					} else {
						log.Error("Missing required parameter: risk_profile_quiz key : " + strconv.FormatUint(uint64(key), 10) + " quiz_question_key can not be blank")
						return lib.CustomError(http.StatusBadRequest,
							"Missing required parameter: risk_profile_quiz key : "+strconv.FormatUint(uint64(key), 10)+" quiz_question_key can not be blank",
							"Missing required parameter: risk_profile_quiz key : "+strconv.FormatUint(uint64(key), 10)+" quiz_question_key can not be blank")
					}
					risk["quiz_question_key"] = val.(string)
				} else {
					log.Error("Missing required parameter: risk_profile_quiz key : " + strconv.FormatUint(uint64(key), 10) + " quiz_question_key tidak ditemukan")
					return lib.CustomError(http.StatusBadRequest,
						"Missing required parameter: risk_profile_quiz key : "+strconv.FormatUint(uint64(key), 10)+" quiz_question_key tidak ditemukan",
						"Missing required parameter: risk_profile_quiz key : "+strconv.FormatUint(uint64(key), 10)+" quiz_question_key tidak ditemukan")
				}

				if val, ok := valueMap["quiz_option_key"]; ok {
					if val.(string) != "" {
						n, err := strconv.ParseUint(val.(string), 10, 64)
						if err != nil || n == 0 {
							log.Error("Wrong input parameter: risk_profile_quiz key : " + strconv.FormatUint(uint64(key), 10) + " quiz_option_key.")
							return lib.CustomError(http.StatusBadRequest,
								"Wrong input parameter: risk_profile_quiz key : "+strconv.FormatUint(uint64(key), 10)+" quiz_option_key.",
								"Wrong input parameter: risk_profile_quiz key : "+strconv.FormatUint(uint64(key), 10)+" quiz_option_key.")
						}
					} else {
						log.Error("Missing required parameter: risk_profile_quiz key : " + strconv.FormatUint(uint64(key), 10) + " quiz_option_key can not be blank")
						return lib.CustomError(http.StatusBadRequest,
							"Missing required parameter: risk_profile_quiz key : "+strconv.FormatUint(uint64(key), 10)+" quiz_option_key can not be blank",
							"Missing required parameter: risk_profile_quiz key : "+strconv.FormatUint(uint64(key), 10)+" quiz_option_key can not be blank")
					}
					risk["quiz_option_key"] = val.(string)
				} else {
					log.Error("Missing required parameter: risk_profile_quiz key : " + strconv.FormatUint(uint64(key), 10) + " quiz_option_key tidak ditemukan")
					return lib.CustomError(http.StatusBadRequest,
						"Missing required parameter: risk_profile_quiz key : "+strconv.FormatUint(uint64(key), 10)+" quiz_option_key tidak ditemukan",
						"Missing required parameter: risk_profile_quiz key : "+strconv.FormatUint(uint64(key), 10)+" quiz_option_key tidak ditemukan")
				}

				riskProfileData = append(riskProfileData, risk)
				key++
			}
		}
	} else {
		if isSaveDraft == "0" {
			log.Error("Missing required parameter: risk_profile_quiz belum diisi semua")
			return lib.CustomError(http.StatusBadRequest, "risk_profile_quiz harus diisi semua.", "risk_profile_quiz harus diisi semua.")
		} else {
			log.Error("Missing required parameter: risk_profile_quiz min 1")
			return lib.CustomError(http.StatusBadRequest, "risk_profile_quiz min 1 data.", "risk_profile_quiz min 1 data.")
		}
	}

	var riskProfileDelete []string
	layout := "2006-01-02 15:04:05"
	score := 0
	if len(riskProfileData) > 0 {
		paramsOption := make(map[string]string)
		paramsOption["rec_status"] = "1"
		paramsOption["orderBy"] = "quiz_option_key"
		paramsOption["orderType"] = "DESC"
		var options []models.CmsQuizOptions
		_, err = models.GetAllCmsQuizOptions(&options, paramsOption)
		if err != nil {
			if err != sql.ErrNoRows {
				log.Error("err get options")
				log.Error(err.Error())
			}
		}

		pData := make(map[uint64]models.CmsQuizOptions)
		for _, opt := range options {
			pData[opt.QuizOptionKey] = opt
		}

		for _, shr := range riskProfileData {
			valueMap := shr.(map[string]interface{})
			if valueMap["risk_profile_quiz_key"].(string) != "" {

				var riskProfileQuiz models.OaRiskProfileQuiz
				status, err = models.GetOaRiskProfileQuiz(&riskProfileQuiz, valueMap["risk_profile_quiz_key"].(string), "risk_profile_quiz_key")
				if err == nil {
					if strconv.FormatUint(*riskProfileQuiz.OaRequestKey, 10) == oaRequestKey {
						riskProfileDelete = append(riskProfileDelete, valueMap["risk_profile_quiz_key"].(string))
						//update
						risk := make(map[string]string)
						risk["risk_profile_quiz_key"] = valueMap["risk_profile_quiz_key"].(string)
						risk["quiz_question_key"] = valueMap["quiz_question_key"].(string)
						risk["quiz_option_key"] = valueMap["quiz_option_key"].(string)

						num, err := strconv.ParseUint(valueMap["quiz_option_key"].(string), 10, 64)
						if err != nil || num == 0 {
							log.Error("Wrong input for parameter: quiz_option_key")
							return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: quiz_option_key", "Wrong input for parameter: oa_request_key")
						}
						if n, ok := pData[num]; ok {
							risk["quiz_option_score"] = strconv.FormatUint(*n.QuizOptionScore, 10)
							score += int(*n.QuizOptionScore)
						}
						risk["rec_modified_date"] = time.Now().Format(layout)
						risk["rec_modified_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
						status, err := models.UpdateOaRiskProfileQuiz(risk)
						if err != nil {
							log.Error("Failed update oa_risk_profile_quiz: " + err.Error())
							return lib.CustomError(status, err.Error(), "Failed update oa_risk_profile_quiz")
						}
					}
				}
			} else {
				//create
				risk := make(map[string]string)
				risk["oa_request_key"] = oaRequestKey
				risk["quiz_question_key"] = valueMap["quiz_question_key"].(string)
				risk["quiz_option_key"] = valueMap["quiz_option_key"].(string)

				num, err := strconv.ParseUint(valueMap["quiz_option_key"].(string), 10, 64)
				if err != nil || num == 0 {
					log.Error("Wrong input for parameter: quiz_option_key")
					return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: quiz_option_key", "Wrong input for parameter: oa_request_key")
				}
				if n, ok := pData[num]; ok {
					risk["quiz_option_score"] = strconv.FormatUint(*n.QuizOptionScore, 10)
					score += int(*n.QuizOptionScore)
				}
				risk["rec_order"] = "0"
				risk["rec_status"] = "1"
				risk["rec_created_date"] = time.Now().Format(layout)
				risk["rec_created_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
				risk["rec_modified_date"] = time.Now().Format(layout)
				risk["rec_modified_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
				status, err, riskID := models.CreateOaRiskProfileQuiz(risk)
				if err != nil {
					log.Error("Failed create oa_risk_profile_quiz: " + err.Error())
					return lib.CustomError(status, err.Error(), "Failed create oa_risk_profile_quiz")
				} else {
					riskProfileDelete = append(riskProfileDelete, riskID)
				}
			}
		}
	}

	//delete oa_risk_profile_quiz
	paramsDeleteRisk := make(map[string]string)
	paramsDeleteRisk["rec_status"] = "0"
	paramsDeleteRisk["rec_deleted_date"] = time.Now().Format(layout)
	paramsDeleteRisk["rec_deleted_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
	status, err = models.DeleteOaRiskProfileQuiz(paramsDeleteRisk, riskProfileDelete, oaRequestKey)
	if err != nil {
		log.Error("Failed delete oa_risk_profile_quiz: " + err.Error())
	}

	//cek oa_risk_profile
	var riskProfile models.MsRiskProfile
	scoreStr := strconv.FormatUint(uint64(score), 10)
	status, err = models.GetMsRiskProfileScore(&riskProfile, scoreStr)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data risk profile")
	}

	var riskProf models.OaRiskProfile
	_, err = models.GetOaRiskProfile(&riskProf, oaRequestKey, "oa_request_key")
	if err != nil {
		log.Error(err.Error())
		//create new
		paramsOaRiskProfile := make(map[string]string)
		paramsOaRiskProfile["oa_request_key"] = oaRequestKey
		paramsOaRiskProfile["risk_profile_key"] = strconv.FormatUint(riskProfile.RiskProfileKey, 10)
		paramsOaRiskProfile["score_result"] = strconv.FormatUint(uint64(score), 10)
		paramsOaRiskProfile["rec_status"] = "1"
		paramsOaRiskProfile["rec_created_date"] = time.Now().Format(layout)
		paramsOaRiskProfile["rec_created_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

		status, err = models.CreateOaRiskProfile(paramsOaRiskProfile)
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed input data")
		}
	} else {
		//update
		paramsOaRiskProfile := make(map[string]string)
		paramsOaRiskProfile["oa_risk_profile_key"] = strconv.FormatUint(riskProf.OaRiskProfileKey, 10)
		paramsOaRiskProfile["risk_profile_key"] = strconv.FormatUint(riskProfile.RiskProfileKey, 10)
		paramsOaRiskProfile["score_result"] = strconv.FormatUint(uint64(score), 10)
		paramsOaRiskProfile["rec_modified_date"] = time.Now().Format(layout)
		paramsOaRiskProfile["rec_modified_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

		status, err = models.UpdateOaRiskProfile(paramsOaRiskProfile)
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed update data")
		}
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = ""

	return c.JSON(http.StatusOK, response)
}

func SaveInstitutionUser(c echo.Context) error {
	var err error
	// var status int
	decimal.MarshalJSONWithoutQuotes = true

	oaRequestKey := c.FormValue("oa_request_key")
	if oaRequestKey != "" {
		n, err := strconv.ParseUint(oaRequestKey, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: oa_request_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: oa_request_key", "Wrong input for parameter: oa_request_key")
		}

		if len(oaRequestKey) > 11 {
			log.Error("Wrong input for parameter: oa_request_key too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: oa_request_key too long, max 11 character", "Missing required parameter: oa_request_key too long, max 11 character")
		}
	} else {
		log.Error("Missing required parameter: oa_request_key")
		return lib.CustomError(http.StatusBadRequest, "oa_request_key can not be blank", "oa_request_key can not be blank")
	}

	if oaRequestKey != "" {
		var oareq models.OaRequest
		_, err = models.GetOaRequestInstitution(&oareq, oaRequestKey, "")
		if err != nil {
			log.Error("OA Request not found.")
			return lib.CustomError(http.StatusBadRequest, "OA Request not found.", "OA Request not found.")
		}

		var userCategory uint64
		userCategory = 3 //Branch
		if lib.Profile.UserCategoryKey == userCategory {
			if oareq.BranchKey != lib.Profile.BranchKey {
				log.Error("User not autorized.")
				return lib.CustomError(http.StatusBadRequest, "User not autorized.", "User not autorized.")
			}
		}

		if strconv.FormatUint(lib.Profile.UserID, 10) != *oareq.RecCreatedBy {
			log.Error("User not autorized.")
			return lib.CustomError(http.StatusBadRequest, "User not autorized.", "User not autorized.")
		}

		if *oareq.Oastatus != uint64(lib.DRAFT) {
			log.Error("User not autorized.")
			return lib.CustomError(http.StatusBadRequest, "User not autorized.", "User not autorized.")
		}
	}

	isSaveDraft := c.FormValue("is_save_draft")
	if isSaveDraft != "" {
		if isSaveDraft != "0" && isSaveDraft != "1" {
			log.Error("Wrong input for parameter: is_save_draft")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: is_save_draft", "Wrong input for parameter: is_save_draft")
		}
	} else {
		log.Error("Missing required parameter: is_save_draft")
		return lib.CustomError(http.StatusBadRequest, "is_save_draft can not be blank", "is_save_draft can not be blank")
	}

	var institutionUserData []interface{}
	var emailList []string
	var noHPList []string

	// institutionUserMaker := c.FormValue("institution_user_maker")
	// institutionUserChecker := c.FormValue("institution_user_checker")
	// institutionUserReleaser := c.FormValue("institution_user_releaser")

	format := []string{"institution_user_maker", "institution_user_checker", "institution_user_releaser"}

	for _, valData := range format {
		institutionUser := c.FormValue(valData)
		if isSaveDraft == "1" {
			if institutionUser != "" {
				var userSlice []interface{}
				err = json.Unmarshal([]byte(institutionUser), &userSlice)
				if err != nil {
					log.Error(err.Error())
					log.Error("Missing required parameter: " + valData)
					return lib.CustomError(http.StatusBadRequest, err.Error(), "Wrong input for parameter: "+valData)
				}
				if len(userSlice) == 0 {
					log.Error("Missing required parameter: " + valData)
					return lib.CustomError(http.StatusBadRequest, err.Error(), "Wrong input for parameter: "+valData)
				}
				if len(userSlice) > 2 {
					log.Error("Missing required parameter: " + valData)
					return lib.CustomError(http.StatusBadRequest, valData+" hanya max 2 data.", valData+" hanya max 2 data.")
				}
				key := 1
				if len(userSlice) > 0 {
					isPriority := 0
					for _, val := range userSlice {
						user := make(map[string]interface{})
						valueMap := val.(map[string]interface{})
						if valData == "institution_user_maker" {
							user["role_key"] = lib.ROLE_INSTITUTION_MAKER
						}
						if valData == "institution_user_checker" {
							user["role_key"] = lib.ROLE_INSTITUTION_CHECKER
						}
						if valData == "institution_user_releaser" {
							user["role_key"] = lib.ROLE_INSTITUTION_RELEASER
						}
						if val, ok := valueMap["insti_user_key"]; ok {
							if val.(string) != "" {
								n, err := strconv.ParseUint(val.(string), 10, 64)
								if err != nil || n == 0 {
									log.Error("Wrong input parameter: " + val.(string) + " key : " + strconv.FormatUint(uint64(key), 10) + " insti_user_key.")
									return lib.CustomError(http.StatusBadRequest,
										"Wrong input parameter: "+val.(string)+" key : "+strconv.FormatUint(uint64(key), 10)+" insti_user_key.",
										"Wrong input parameter: "+val.(string)+" key : "+strconv.FormatUint(uint64(key), 10)+" insti_user_key.")
								}
								var insUserEks models.OaInstitutionUser
								_, err = models.GetOaInstitutionUser(&insUserEks, val.(string), "insti_user_key")
								if err != nil {
									log.Error("Institution User not found.")
									return lib.CustomError(http.StatusBadRequest, "Institution User not found.", "Institution User not found.")
								}

								if strconv.FormatUint(insUserEks.OaRequestKey, 10) != oaRequestKey {
									log.Error("User not autorized. OA request berbeda dengan oa_request_key di oa_institution_user")
									return lib.CustomError(http.StatusBadRequest, "User not autorized.", "User not autorized.")
								}
							}
							user["insti_user_key"] = val.(string)
						} else {
							log.Error("Missing required parameter: " + valData + " key : " + strconv.FormatUint(uint64(key), 10) + " insti_user_key tidak ditemukan")
							return lib.CustomError(http.StatusBadRequest,
								"Missing required parameter: "+valData+" key : "+strconv.FormatUint(uint64(key), 10)+" insti_user_key tidak ditemukan",
								"Missing required parameter: "+valData+" key : "+strconv.FormatUint(uint64(key), 10)+" insti_user_key tidak ditemukan")
						}

						if val, ok := valueMap["full_name"]; ok {
							if val.(string) != "" {
								if len(val.(string)) > 100 {
									log.Error("Wrong input for parameter: full_name too long")
									return lib.CustomError(http.StatusBadRequest, "Missing required parameter: full_name too long, max 100 character", "Missing required parameter: full_name too long, max 100 character")
								}
							}
							user["full_name"] = val.(string)
						} else {
							log.Error("Missing required parameter: " + valData + " key : " + strconv.FormatUint(uint64(key), 10) + " full_name tidak ditemukan")
							return lib.CustomError(http.StatusBadRequest,
								"Missing required parameter: "+valData+" key : "+strconv.FormatUint(uint64(key), 10)+" full_name tidak ditemukan",
								"Missing required parameter: "+valData+" key : "+strconv.FormatUint(uint64(key), 10)+" full_name tidak ditemukan")
						}

						if val, ok := valueMap["email_address"]; ok {
							if val.(string) != "" {
								if len(val.(string)) > 100 {
									log.Error("Wrong input for parameter: email_address too long")
									return lib.CustomError(http.StatusBadRequest, "Missing required parameter: email_address too long, max 100 character", "Missing required parameter: email_address too long, max 100 character")
								}
								if !lib.IsValidEmail(val.(string)) {
									log.Error("Wrong input for parameter: email_address wrong format email")
									return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: email_address wrong format email", "Wrong input for parameter: email_address wrong format email")
								}
								//cek user login
								var countData models.CountData
								status, err := models.ValidateUniqueData(&countData, "ulogin_email", val.(string), nil)
								if err != nil {
									log.Error(err.Error())
									return lib.CustomError(status, err.Error(), "Failed get data")
								}
								if int(countData.CountData) > int(0) {
									log.Error("Wrong input parameter: " + valData + " key : " + strconv.FormatUint(uint64(key), 10) + " email_address sudah digunakan.")
									return lib.CustomError(http.StatusBadRequest,
										"Wrong input parameter: "+valData+" key : "+strconv.FormatUint(uint64(key), 10)+" email_address sudah digunakan.",
										"Wrong input parameter: "+valData+" key : "+strconv.FormatUint(uint64(key), 10)+" email_address sudah digunakan.")
								}

								//cek oa_institution_user
								var countDataIns models.CountData
								if valueMap["insti_user_key"].(string) != "" {
									ky := valueMap["insti_user_key"].(string)
									status, err = models.ValidateUniqueInstitutionUser(&countDataIns, "email_address", val.(string), &ky)
								} else {
									status, err = models.ValidateUniqueInstitutionUser(&countDataIns, "email_address", val.(string), nil)
								}
								if err != nil {
									log.Error(err.Error())
									return lib.CustomError(status, err.Error(), "Failed get data")
								}
								if int(countDataIns.CountData) > int(0) {
									log.Error("Wrong input parameter: " + valData + " key : " + strconv.FormatUint(uint64(key), 10) + " email_address sudah digunakan.")
									return lib.CustomError(http.StatusBadRequest,
										"Wrong input parameter: "+valData+" key : "+strconv.FormatUint(uint64(key), 10)+" email_address sudah digunakan.",
										"Wrong input parameter: "+valData+" key : "+strconv.FormatUint(uint64(key), 10)+" email_address sudah digunakan.")
								}

								if _, ok := lib.Find(emailList, val.(string)); !ok {
									emailList = append(emailList, val.(string))
								} else {
									log.Error("Wrong input parameter: " + valData + " key : " + strconv.FormatUint(uint64(key), 10) + " email_address sudah digunakan.")
									return lib.CustomError(http.StatusBadRequest,
										"Wrong input parameter: "+valData+" key : "+strconv.FormatUint(uint64(key), 10)+" email_address sudah digunakan.",
										"Wrong input parameter: "+valData+" key : "+strconv.FormatUint(uint64(key), 10)+" email_address sudah digunakan.")
								}
							}
							user["email_address"] = val.(string)
						} else {
							log.Error("Missing required parameter: " + valData + " key : " + strconv.FormatUint(uint64(key), 10) + " email_address tidak ditemukan")
							return lib.CustomError(http.StatusBadRequest,
								"Missing required parameter: "+valData+" key : "+strconv.FormatUint(uint64(key), 10)+" email_address tidak ditemukan",
								"Missing required parameter: "+valData+" key : "+strconv.FormatUint(uint64(key), 10)+" email_address tidak ditemukan")
						}

						if val, ok := valueMap["phone_number"]; ok {
							if val.(string) != "" {
								if len(val.(string)) > 20 {
									log.Error("Wrong input for parameter: phone_number too long")
									return lib.CustomError(http.StatusBadRequest, "Missing required parameter: phone_number too long, max 20 character", "Missing required parameter: phone_number too long, max 20 character")
								}
								//cek user login
								var countData models.CountData
								status, err := models.ValidateUniqueData(&countData, "ulogin_mobileno", val.(string), nil)
								if err != nil {
									log.Error(err.Error())
									return lib.CustomError(status, err.Error(), "Failed get data")
								}
								if int(countData.CountData) > int(0) {
									log.Error("Wrong input parameter: " + valData + " key : " + strconv.FormatUint(uint64(key), 10) + " phone_number sudah digunakan.")
									return lib.CustomError(http.StatusBadRequest,
										"Wrong input parameter: "+valData+" key : "+strconv.FormatUint(uint64(key), 10)+" phone_number sudah digunakan.",
										"Wrong input parameter: "+valData+" key : "+strconv.FormatUint(uint64(key), 10)+" phone_number sudah digunakan.")
								}

								//cek oa_institution_user
								var countDataIns models.CountData
								if valueMap["insti_user_key"].(string) != "" {
									ky := valueMap["insti_user_key"].(string)
									status, err = models.ValidateUniqueInstitutionUser(&countDataIns, "phone_number", val.(string), &ky)
								} else {
									status, err = models.ValidateUniqueInstitutionUser(&countDataIns, "phone_number", val.(string), nil)
								}
								if err != nil {
									log.Error(err.Error())
									return lib.CustomError(status, err.Error(), "Failed get data")
								}
								if int(countDataIns.CountData) > int(0) {
									log.Error("Wrong input parameter: " + valData + " key : " + strconv.FormatUint(uint64(key), 10) + " phone_number sudah digunakan.")
									return lib.CustomError(http.StatusBadRequest,
										"Wrong input parameter: "+valData+" key : "+strconv.FormatUint(uint64(key), 10)+" phone_number sudah digunakan.",
										"Wrong input parameter: "+valData+" key : "+strconv.FormatUint(uint64(key), 10)+" phone_number sudah digunakan.")
								}

								if _, ok := lib.Find(noHPList, val.(string)); !ok {
									noHPList = append(noHPList, val.(string))
								} else {
									log.Error("Wrong input parameter: " + valData + " key : " + strconv.FormatUint(uint64(key), 10) + " phone_number sudah digunakan.")
									return lib.CustomError(http.StatusBadRequest,
										"Wrong input parameter: "+valData+" key : "+strconv.FormatUint(uint64(key), 10)+" phone_number sudah digunakan.",
										"Wrong input parameter: "+valData+" key : "+strconv.FormatUint(uint64(key), 10)+" phone_number sudah digunakan.")
								}
							}
							user["phone_number"] = val.(string)
						} else {
							log.Error("Missing required parameter: " + valData + " key : " + strconv.FormatUint(uint64(key), 10) + " phone_number tidak ditemukan")
							return lib.CustomError(http.StatusBadRequest,
								"Missing required parameter: "+valData+" key : "+strconv.FormatUint(uint64(key), 10)+" phone_number tidak ditemukan",
								"Missing required parameter: "+valData+" key : "+strconv.FormatUint(uint64(key), 10)+" phone_number tidak ditemukan")
						}

						if val, ok := valueMap["user_priority"]; ok {
							if val.(string) != "" {
								if val.(string) != "0" && val.(string) != "1" {
									log.Error("Wrong input parameter: " + valData + " key : " + strconv.FormatUint(uint64(key), 10) + " user_priority.")
									return lib.CustomError(http.StatusBadRequest,
										"Wrong input parameter: "+valData+" key : "+strconv.FormatUint(uint64(key), 10)+" user_priority.",
										"Wrong input parameter: "+valData+" key : "+strconv.FormatUint(uint64(key), 10)+" user_priority.")

								}
								user["user_priority"] = val.(string)
								if val.(string) == "1" {
									isPriority++
								}
							} else {
								user["user_priority"] = "0"
							}
						} else {
							log.Error("Missing required parameter: " + valData + " key : " + strconv.FormatUint(uint64(key), 10) + " user_priority tidak ditemukan")
							return lib.CustomError(http.StatusBadRequest,
								"Missing required parameter: "+valData+" key : "+strconv.FormatUint(uint64(key), 10)+" user_priority tidak ditemukan",
								"Missing required parameter: "+valData+" key : "+strconv.FormatUint(uint64(key), 10)+" user_priority tidak ditemukan")
						}

						institutionUserData = append(institutionUserData, user)
						key++
					}

					if isPriority > 1 {
						log.Error("Missing required parameter: " + valData + " priority hanya boleh 1 data")
						return lib.CustomError(http.StatusBadRequest, valData+" priority hanya boleh 1 data.", valData+" priority hanya boleh 1 data.")
					}
				}
			}
		} else { //save validation
			if institutionUser == "" {
				log.Error("Missing required parameter: " + valData + "")
				return lib.CustomError(http.StatusBadRequest, ""+valData+" can not be blank", ""+valData+" can not be blank")
			} else {
				var userSlice []interface{}
				err = json.Unmarshal([]byte(institutionUser), &userSlice)
				if err != nil {
					log.Error(err.Error())
					log.Error("Missing required parameter: " + valData)
					return lib.CustomError(http.StatusBadRequest, err.Error(), "Wrong input for parameter: "+valData)
				}
				if len(userSlice) == 0 {
					log.Error("Missing required parameter: " + valData)
					return lib.CustomError(http.StatusBadRequest, err.Error(), "Wrong input for parameter: "+valData)
				}
				if len(userSlice) > 2 {
					log.Error("Missing required parameter: " + valData)
					return lib.CustomError(http.StatusBadRequest, valData+" hanya max 2 data.", valData+" hanya max 2 data.")
				}
				key := 1
				if len(userSlice) > 0 {
					isPriority := 0
					for _, val := range userSlice {
						user := make(map[string]interface{})
						valueMap := val.(map[string]interface{})
						if valData == "institution_user_maker" {
							user["role_key"] = lib.ROLE_INSTITUTION_MAKER
						}
						if valData == "institution_user_checker" {
							user["role_key"] = lib.ROLE_INSTITUTION_CHECKER
						}
						if valData == "institution_user_releaser" {
							user["role_key"] = lib.ROLE_INSTITUTION_RELEASER
						}
						if val, ok := valueMap["insti_user_key"]; ok {
							if val.(string) != "" {
								n, err := strconv.ParseUint(val.(string), 10, 64)
								if err != nil || n == 0 {
									log.Error("Wrong input parameter: " + valData + " key : " + strconv.FormatUint(uint64(key), 10) + " insti_user_key.")
									return lib.CustomError(http.StatusBadRequest,
										"Wrong input parameter: "+valData+" key : "+strconv.FormatUint(uint64(key), 10)+" insti_user_key.",
										"Wrong input parameter: "+valData+" key : "+strconv.FormatUint(uint64(key), 10)+" insti_user_key.")
								}
								var insUserEks models.OaInstitutionUser
								_, err = models.GetOaInstitutionUser(&insUserEks, val.(string), "oa_institution_user")
								if err != nil {
									log.Error("Institution User not found.")
									return lib.CustomError(http.StatusBadRequest, "Institution User not found.", "Institution User not found.")
								}

								if strconv.FormatUint(insUserEks.OaRequestKey, 10) != oaRequestKey {
									log.Error("User not autorized. OA request berbeda dengan oa_request_key di oa_institution_user")
									return lib.CustomError(http.StatusBadRequest, "User not autorized.", "User not autorized.")
								}
							}
							user["insti_user_key"] = val.(string)
						} else {
							log.Error("Missing required parameter: " + valData + " key : " + strconv.FormatUint(uint64(key), 10) + " insti_user_key tidak ditemukan")
							return lib.CustomError(http.StatusBadRequest,
								"Missing required parameter: "+valData+" key : "+strconv.FormatUint(uint64(key), 10)+" insti_user_key tidak ditemukan",
								"Missing required parameter: "+valData+" key : "+strconv.FormatUint(uint64(key), 10)+" insti_user_key tidak ditemukan")
						}

						if val, ok := valueMap["full_name"]; ok {
							if val.(string) != "" {
								n, err := strconv.ParseUint(val.(string), 10, 64)
								if err != nil || n == 0 {
									log.Error("Wrong input parameter: " + valData + " key : " + strconv.FormatUint(uint64(key), 10) + " full_name.")
									return lib.CustomError(http.StatusBadRequest,
										"Wrong input parameter: "+valData+" key : "+strconv.FormatUint(uint64(key), 10)+" full_name.",
										"Wrong input parameter: "+valData+" key : "+strconv.FormatUint(uint64(key), 10)+" full_name.")
								}
							} else {
								log.Error("Missing required parameter: " + valData + " key : " + strconv.FormatUint(uint64(key), 10) + " full_name can not be blank")
								return lib.CustomError(http.StatusBadRequest,
									"Missing required parameter: "+valData+" key : "+strconv.FormatUint(uint64(key), 10)+" full_name can not be blank",
									"Missing required parameter: "+valData+" key : "+strconv.FormatUint(uint64(key), 10)+" full_name can not be blank")
							}
							user["full_name"] = val.(string)
						} else {
							log.Error("Missing required parameter: " + valData + " key : " + strconv.FormatUint(uint64(key), 10) + " full_name tidak ditemukan")
							return lib.CustomError(http.StatusBadRequest,
								"Missing required parameter: "+valData+" key : "+strconv.FormatUint(uint64(key), 10)+" full_name tidak ditemukan",
								"Missing required parameter: "+valData+" key : "+strconv.FormatUint(uint64(key), 10)+" full_name tidak ditemukan")
						}

						if val, ok := valueMap["email_address"]; ok {
							if val.(string) != "" {
								if len(val.(string)) > 100 {
									log.Error("Wrong input for parameter: email_address too long")
									return lib.CustomError(http.StatusBadRequest, "Missing required parameter: email_address too long, max 100 character", "Missing required parameter: email_address too long, max 100 character")
								}
								if !lib.IsValidEmail(val.(string)) {
									log.Error("Wrong input for parameter: email_address wrong format email")
									return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: email_address wrong format email", "Wrong input for parameter: email_address wrong format email")
								}
								//cek user login
								var countData models.CountData
								status, err := models.ValidateUniqueData(&countData, "ulogin_email", val.(string), nil)
								if err != nil {
									log.Error(err.Error())
									return lib.CustomError(status, err.Error(), "Failed get data")
								}
								if int(countData.CountData) > int(0) {
									log.Error("Wrong input parameter: " + valData + " key : " + strconv.FormatUint(uint64(key), 10) + " email_address sudah digunakan.")
									return lib.CustomError(http.StatusBadRequest,
										"Wrong input parameter: "+valData+" key : "+strconv.FormatUint(uint64(key), 10)+" email_address sudah digunakan.",
										"Wrong input parameter: "+valData+" key : "+strconv.FormatUint(uint64(key), 10)+" email_address sudah digunakan.")
								}

								//cek oa_institution_user
								var countDataIns models.CountData
								if valueMap["insti_user_key"].(string) != "" {
									ky := valueMap["insti_user_key"].(string)
									status, err = models.ValidateUniqueInstitutionUser(&countDataIns, "email_address", val.(string), &ky)
								} else {
									status, err = models.ValidateUniqueInstitutionUser(&countDataIns, "email_address", val.(string), nil)
								}
								if err != nil {
									log.Error(err.Error())
									return lib.CustomError(status, err.Error(), "Failed get data")
								}
								if int(countDataIns.CountData) > int(0) {
									log.Error("Wrong input parameter: " + valData + " key : " + strconv.FormatUint(uint64(key), 10) + " email_address sudah digunakan.")
									return lib.CustomError(http.StatusBadRequest,
										"Wrong input parameter: "+valData+" key : "+strconv.FormatUint(uint64(key), 10)+" email_address sudah digunakan.",
										"Wrong input parameter: "+valData+" key : "+strconv.FormatUint(uint64(key), 10)+" email_address sudah digunakan.")
								}

								if _, ok := lib.Find(emailList, val.(string)); !ok {
									emailList = append(emailList, val.(string))
								} else {
									log.Error("Wrong input parameter: " + valData + " key : " + strconv.FormatUint(uint64(key), 10) + " email_address sudah digunakan.")
									return lib.CustomError(http.StatusBadRequest,
										"Wrong input parameter: "+valData+" key : "+strconv.FormatUint(uint64(key), 10)+" email_address sudah digunakan.",
										"Wrong input parameter: "+valData+" key : "+strconv.FormatUint(uint64(key), 10)+" email_address sudah digunakan.")
								}
							} else {
								log.Error("Missing required parameter: " + valData + " key : " + strconv.FormatUint(uint64(key), 10) + " email_address can not be blank")
								return lib.CustomError(http.StatusBadRequest,
									"Missing required parameter: "+valData+" key : "+strconv.FormatUint(uint64(key), 10)+" email_address can not be blank",
									"Missing required parameter: "+valData+" key : "+strconv.FormatUint(uint64(key), 10)+" email_address can not be blank")
							}
							user["email_address"] = val.(string)
						} else {
							log.Error("Missing required parameter: " + valData + " key : " + strconv.FormatUint(uint64(key), 10) + " email_address tidak ditemukan")
							return lib.CustomError(http.StatusBadRequest,
								"Missing required parameter: "+valData+" key : "+strconv.FormatUint(uint64(key), 10)+" email_address tidak ditemukan",
								"Missing required parameter: "+valData+" key : "+strconv.FormatUint(uint64(key), 10)+" email_address tidak ditemukan")
						}

						if val, ok := valueMap["phone_number"]; ok {
							if val.(string) != "" {
								if len(val.(string)) > 20 {
									log.Error("Wrong input for parameter: phone_number too long")
									return lib.CustomError(http.StatusBadRequest, "Missing required parameter: phone_number too long, max 20 character", "Missing required parameter: phone_number too long, max 20 character")
								}
								//cek user login
								var countData models.CountData
								status, err := models.ValidateUniqueData(&countData, "ulogin_mobileno", val.(string), nil)
								if err != nil {
									log.Error(err.Error())
									return lib.CustomError(status, err.Error(), "Failed get data")
								}
								if int(countData.CountData) > int(0) {
									log.Error("Wrong input parameter: " + valData + " key : " + strconv.FormatUint(uint64(key), 10) + " phone_number sudah digunakan.")
									return lib.CustomError(http.StatusBadRequest,
										"Wrong input parameter: "+valData+" key : "+strconv.FormatUint(uint64(key), 10)+" phone_number sudah digunakan.",
										"Wrong input parameter: "+valData+" key : "+strconv.FormatUint(uint64(key), 10)+" phone_number sudah digunakan.")
								}

								//cek oa_institution_user
								var countDataIns models.CountData
								if valueMap["insti_user_key"].(string) != "" {
									ky := valueMap["insti_user_key"].(string)
									status, err = models.ValidateUniqueInstitutionUser(&countDataIns, "phone_number", val.(string), &ky)
								} else {
									status, err = models.ValidateUniqueInstitutionUser(&countDataIns, "phone_number", val.(string), nil)
								}
								if err != nil {
									log.Error(err.Error())
									return lib.CustomError(status, err.Error(), "Failed get data")
								}
								if int(countDataIns.CountData) > int(0) {
									log.Error("Wrong input parameter: " + valData + " key : " + strconv.FormatUint(uint64(key), 10) + " phone_number sudah digunakan.")
									return lib.CustomError(http.StatusBadRequest,
										"Wrong input parameter: "+valData+" key : "+strconv.FormatUint(uint64(key), 10)+" phone_number sudah digunakan.",
										"Wrong input parameter: "+valData+" key : "+strconv.FormatUint(uint64(key), 10)+" phone_number sudah digunakan.")
								}

								if _, ok := lib.Find(noHPList, val.(string)); !ok {
									noHPList = append(noHPList, val.(string))
								} else {
									log.Error("Wrong input parameter: " + valData + " key : " + strconv.FormatUint(uint64(key), 10) + " phone_number sudah digunakan.")
									return lib.CustomError(http.StatusBadRequest,
										"Wrong input parameter: "+valData+" key : "+strconv.FormatUint(uint64(key), 10)+" phone_number sudah digunakan.",
										"Wrong input parameter: "+valData+" key : "+strconv.FormatUint(uint64(key), 10)+" phone_number sudah digunakan.")
								}
							} else {
								log.Error("Missing required parameter: " + valData + " key : " + strconv.FormatUint(uint64(key), 10) + " phone_number can not be blank")
								return lib.CustomError(http.StatusBadRequest,
									"Missing required parameter: "+valData+" key : "+strconv.FormatUint(uint64(key), 10)+" phone_number can not be blank",
									"Missing required parameter: "+valData+" key : "+strconv.FormatUint(uint64(key), 10)+" phone_number can not be blank")
							}
							user["phone_number"] = val.(string)
						} else {
							log.Error("Missing required parameter: " + valData + " key : " + strconv.FormatUint(uint64(key), 10) + " phone_number tidak ditemukan")
							return lib.CustomError(http.StatusBadRequest,
								"Missing required parameter: "+valData+" key : "+strconv.FormatUint(uint64(key), 10)+" phone_number tidak ditemukan",
								"Missing required parameter: "+valData+" key : "+strconv.FormatUint(uint64(key), 10)+" phone_number tidak ditemukan")
						}

						if val, ok := valueMap["user_priority"]; ok {
							if val.(string) != "" {
								if val.(string) != "0" && val.(string) != "1" {
									log.Error("Wrong input parameter: " + valData + " key : " + strconv.FormatUint(uint64(key), 10) + " user_priority.")
									return lib.CustomError(http.StatusBadRequest,
										"Wrong input parameter: "+valData+" key : "+strconv.FormatUint(uint64(key), 10)+" user_priority.",
										"Wrong input parameter: "+valData+" key : "+strconv.FormatUint(uint64(key), 10)+" user_priority.")

								}
								user["user_priority"] = val.(string)
								if val.(string) == "1" {
									isPriority++
								}
							} else {
								user["user_priority"] = "0"
							}
						} else {
							log.Error("Missing required parameter: " + valData + " key : " + strconv.FormatUint(uint64(key), 10) + " user_priority tidak ditemukan")
							return lib.CustomError(http.StatusBadRequest,
								"Missing required parameter: "+valData+" key : "+strconv.FormatUint(uint64(key), 10)+" user_priority tidak ditemukan",
								"Missing required parameter: "+valData+" key : "+strconv.FormatUint(uint64(key), 10)+" user_priority tidak ditemukan")
						}

						institutionUserData = append(institutionUserData, user)
						key++
					}

					if isPriority != 1 {
						log.Error("Missing required parameter: " + valData + " pilih 1 priority")
						return lib.CustomError(http.StatusBadRequest, valData+" pilih 1 priority.", valData+" pilih 1 priority.")
					}
				}
			}
		}
	}

	layout := "2006-01-02 15:04:05"
	//save user
	var userNotDelete []string
	if len(institutionUserData) > 0 {
		if len(institutionUserData) > 0 {
			for _, bd := range institutionUserData {
				valueMap := bd.(map[string]interface{})
				if valueMap["role_key"].(string) != "" || valueMap["insti_user_key"].(string) != "" ||
					valueMap["full_name"].(string) != "" || valueMap["email_address"].(string) != "" ||
					valueMap["phone_number"].(string) != "" || valueMap["user_priority"].(string) != "" {
					if valueMap["insti_user_key"].(string) != "" {
						userNotDelete = append(userNotDelete, valueMap["insti_user_key"].(string))
						//update
						var userEx models.OaInstitutionUser
						_, err = models.GetOaInstitutionUser(&userEx, valueMap["insti_user_key"].(string), "insti_user_key")
						if err == nil {
							if strconv.FormatUint(userEx.OaRequestKey, 10) == oaRequestKey {
								paramsUser := make(map[string]string)
								paramsUser["insti_user_key"] = valueMap["insti_user_key"].(string)
								paramsUser["full_name"] = valueMap["full_name"].(string)
								paramsUser["email_address"] = valueMap["email_address"].(string)
								paramsUser["phone_number"] = valueMap["phone_number"].(string)
								paramsUser["user_priority"] = valueMap["user_priority"].(string)
								paramsUser["role_key"] = valueMap["role_key"].(string)
								paramsUser["rec_modified_date"] = time.Now().Format(layout)
								paramsUser["rec_modified_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
								status, err := models.UpdateOaInstitutionUser(paramsUser)
								if err != nil {
									log.Error("Failed update oa_institution_user: " + err.Error())
									return lib.CustomError(status, err.Error(), "Failed update oa_institution_user")
								}
							} else {
								log.Error("Failed update oa_institution_user : oa_request_key != oa_request_key oa_institution_user")
							}
						} else {
							log.Error("Failed get oa_institution_user key: " + valueMap["insti_user_key"].(string))
						}
					} else {
						//create
						paramsUser := make(map[string]string)
						paramsUser["oa_request_key"] = oaRequestKey
						paramsUser["full_name"] = valueMap["full_name"].(string)
						paramsUser["email_address"] = valueMap["email_address"].(string)
						paramsUser["email_verified_flag"] = "0"
						paramsUser["phone_number"] = valueMap["phone_number"].(string)
						paramsUser["phone_verified_flag"] = "0"
						paramsUser["user_priority"] = valueMap["user_priority"].(string)
						paramsUser["role_category_key"] = "1"
						paramsUser["rec_order"] = "0"
						paramsUser["rec_status"] = "1"
						paramsUser["role_key"] = valueMap["role_key"].(string)
						paramsUser["rec_created_date"] = time.Now().Format(layout)
						paramsUser["rec_created_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
						paramsUser["rec_modified_date"] = time.Now().Format(layout)
						paramsUser["rec_modified_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
						status, err, userID := models.CreateOaInstitutionUser(paramsUser)
						if err != nil {
							log.Error("Failed create oa_institution_user: " + err.Error())
							return lib.CustomError(status, err.Error(), "Failed create oa_institution_user")
						} else {
							userNotDelete = append(userNotDelete, userID)
						}
					}
				}
			}
		}
	}
	//delete user
	paramsDeleteuser := make(map[string]string)
	paramsDeleteuser["rec_status"] = "0"
	paramsDeleteuser["rec_deleted_date"] = time.Now().Format(layout)
	paramsDeleteuser["rec_deleted_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
	_, err = models.DeleteOaInstitutionUser(paramsDeleteuser, userNotDelete, oaRequestKey)
	if err != nil {
		log.Error("Failed delete oa_institution_user: " + err.Error())
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = ""

	return c.JSON(http.StatusOK, response)
}

func SaveInstitutionToApprover(c echo.Context) error {
	var err error
	// var status int
	decimal.MarshalJSONWithoutQuotes = true

	oaRequestKey := c.FormValue("oa_request_key")
	if oaRequestKey != "" {
		n, err := strconv.ParseUint(oaRequestKey, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: oa_request_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: oa_request_key", "Wrong input for parameter: oa_request_key")
		}

		if len(oaRequestKey) > 11 {
			log.Error("Wrong input for parameter: oa_request_key too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: oa_request_key too long, max 11 character", "Missing required parameter: oa_request_key too long, max 11 character")
		}
	} else {
		log.Error("Missing required parameter: oa_request_key")
		return lib.CustomError(http.StatusBadRequest, "oa_request_key can not be blank", "oa_request_key can not be blank")
	}

	var oareq models.OaRequest
	if oaRequestKey != "" {
		_, err = models.GetOaRequestInstitution(&oareq, oaRequestKey, "")
		if err != nil {
			log.Error("OA Request not found.")
			return lib.CustomError(http.StatusBadRequest, "OA Request not found.", "OA Request not found.")
		}

		var userCategory uint64
		userCategory = 3 //Branch
		if lib.Profile.UserCategoryKey == userCategory {
			if oareq.BranchKey != lib.Profile.BranchKey {
				log.Error("User not autorized.")
				return lib.CustomError(http.StatusBadRequest, "User not autorized.", "User not autorized.")
			}
		}

		if strconv.FormatUint(lib.Profile.UserID, 10) != *oareq.RecCreatedBy {
			log.Error("User not autorized.")
			return lib.CustomError(http.StatusBadRequest, "User not autorized.", "User not autorized.")
		}

		if *oareq.Oastatus != uint64(lib.DRAFT) {
			log.Error("User not autorized.")
			return lib.CustomError(http.StatusBadRequest, "User not autorized.", "User not autorized.")
		}
	}

	// VALIDASI DATA
	var dataError []interface{}
	var oadata models.OaInstitutionData
	_, err = models.GetOaInstitutionData(&oadata, oaRequestKey, "oa_request_key")
	if err != nil {
		d := make(map[string]interface{})
		d["field"] = "data"
		d["message"] = "OA Data tidak lengkap."
		dataError = append(dataError, d)
	} else {
		if oadata.Nationality == nil {
			d := make(map[string]interface{})
			d["field"] = "nationality"
			d["message"] = "Missing required parameter: nationality"
			dataError = append(dataError, d)
		}
		if oadata.FullName == nil {
			d := make(map[string]interface{})
			d["field"] = "full_name"
			d["message"] = "Missing required parameter: full_name"
			dataError = append(dataError, d)
		}
		if oadata.ShortName == nil {
			d := make(map[string]interface{})
			d["field"] = "short_name"
			d["message"] = "Missing required parameter: short_name"
			dataError = append(dataError, d)
		}
		if oadata.TinNumber == nil {
			d := make(map[string]interface{})
			d["field"] = "tin_number"
			d["message"] = "Missing required parameter: tin_number"
			dataError = append(dataError, d)
		}
		if oadata.EstablishedCity == nil {
			d := make(map[string]interface{})
			d["field"] = "established_city"
			d["message"] = "Missing required parameter: established_city"
			dataError = append(dataError, d)
		}
		if oadata.EstablishedDate == nil {
			d := make(map[string]interface{})
			d["field"] = "established_date"
			d["message"] = "Missing required parameter: established_date"
			dataError = append(dataError, d)
		}
		if oadata.DeedNo == nil {
			d := make(map[string]interface{})
			d["field"] = "deed_no"
			d["message"] = "Missing required parameter: deed_no"
			dataError = append(dataError, d)
		}
		if oadata.LsEstablishValidationNo == nil {
			d := make(map[string]interface{})
			d["field"] = "ls_establish_validation_no"
			d["message"] = "Missing required parameter: ls_establish_validation_no"
			dataError = append(dataError, d)
		}
		if oadata.LsEstablishValidationDate == nil {
			d := make(map[string]interface{})
			d["field"] = "ls_establish_validation_date"
			d["message"] = "Missing required parameter: ls_establish_validation_date"
			dataError = append(dataError, d)
		}
		if oadata.LastChangeAaNo == nil {
			d := make(map[string]interface{})
			d["field"] = "last_change_aa_no"
			d["message"] = "Missing required parameter: last_change_aa_no"
			dataError = append(dataError, d)
		}
		if oadata.LastChangeAaDate == nil {
			d := make(map[string]interface{})
			d["field"] = "last_change_aa_date"
			d["message"] = "Missing required parameter: last_change_aa_date"
			dataError = append(dataError, d)
		}
		if oadata.LsLastChangeAaNo == nil {
			d := make(map[string]interface{})
			d["field"] = "ls_last_change_aa_no"
			d["message"] = "Missing required parameter: ls_last_change_aa_no"
			dataError = append(dataError, d)
		}
		if oadata.LsLastChangeAaDate == nil {
			d := make(map[string]interface{})
			d["field"] = "ls_last_change_aa_date"
			d["message"] = "Missing required parameter: ls_last_change_aa_date"
			dataError = append(dataError, d)
		}
		if oadata.ManagementDeedNo == nil {
			d := make(map[string]interface{})
			d["field"] = "management_deed_no"
			d["message"] = "Missing required parameter: management_deed_no"
			dataError = append(dataError, d)
		}
		if oadata.ManagementDeedDate == nil {
			d := make(map[string]interface{})
			d["field"] = "management_deed_date"
			d["message"] = "Missing required parameter: management_deed_date"
			dataError = append(dataError, d)
		}
		if oadata.LsMgtChangeDeedNo == nil {
			d := make(map[string]interface{})
			d["field"] = "ls_mgt_change_deed_no"
			d["message"] = "Missing required parameter: ls_mgt_change_deed_no"
			dataError = append(dataError, d)
		}
		if oadata.LsMgtChangeDeedDate == nil {
			d := make(map[string]interface{})
			d["field"] = "ls_mgt_change_deed_date"
			d["message"] = "Missing required parameter: ls_mgt_change_deed_date"
			dataError = append(dataError, d)
		}
		if oadata.SkdLicenseNo == nil {
			d := make(map[string]interface{})
			d["field"] = "skd_license_no"
			d["message"] = "Missing required parameter: skd_license_no"
			dataError = append(dataError, d)
		}
		if oadata.SkdLicenseDate == nil {
			d := make(map[string]interface{})
			d["field"] = "skd_license_date"
			d["message"] = "Missing required parameter: skd_license_date"
			dataError = append(dataError, d)
		}
		if oadata.BizLicenseNo == nil {
			d := make(map[string]interface{})
			d["field"] = "biz_license_no"
			d["message"] = "Missing required parameter: biz_license_no"
			dataError = append(dataError, d)
		}
		if oadata.BizLicenseDate == nil {
			d := make(map[string]interface{})
			d["field"] = "biz_license_date"
			d["message"] = "Missing required parameter: biz_license_date"
			dataError = append(dataError, d)
		}
		if oadata.NibNo == nil {
			d := make(map[string]interface{})
			d["field"] = "nib_no"
			d["message"] = "Missing required parameter: nib_no"
			dataError = append(dataError, d)
		}
		if oadata.NibDate == nil {
			d := make(map[string]interface{})
			d["field"] = "nib_date"
			d["message"] = "Missing required parameter: nib_date"
			dataError = append(dataError, d)
		}
		if oadata.MobileNo == nil {
			d := make(map[string]interface{})
			d["field"] = "mobile_no"
			d["message"] = "Missing required parameter: mobile_no"
			dataError = append(dataError, d)
		}
		if oadata.FaxNo == nil {
			d := make(map[string]interface{})
			d["field"] = "fax_no"
			d["message"] = "Missing required parameter: fax_no"
			dataError = append(dataError, d)
		}
		if oadata.EmailAddress == nil {
			d := make(map[string]interface{})
			d["field"] = "email_address"
			d["message"] = "Missing required parameter: email_address"
			dataError = append(dataError, d)
		}
		if oadata.DomicileKey == nil {
			d := make(map[string]interface{})
			d["field"] = "domicile"
			d["message"] = "Missing required parameter: domicile"
			dataError = append(dataError, d)
		} else {
			var dom models.OaPostalAddress
			_, err = models.GetOaPostalAddress(&dom, strconv.FormatUint(*oadata.DomicileKey, 10))
			if err == nil {
				if dom.AddressLine1 == nil {
					d := make(map[string]interface{})
					d["field"] = "domicile_address"
					d["message"] = "Missing required parameter: domicile_address"
					dataError = append(dataError, d)
				}
				if dom.KecamatanKey == nil {
					d := make(map[string]interface{})
					d["field"] = "domicile_subdistric"
					d["message"] = "Missing required parameter: domicile_subdistric"
					dataError = append(dataError, d)
				}
				if dom.KabupatenKey == nil {
					d := make(map[string]interface{})
					d["field"] = "domicile_city"
					d["message"] = "Missing required parameter: domicile_city"
					dataError = append(dataError, d)
				}
				if dom.PostalCode == nil {
					d := make(map[string]interface{})
					d["field"] = "domicile_postalcode"
					d["message"] = "Missing required parameter: domicile_postalcode"
					dataError = append(dataError, d)
				}
			} else {
				d := make(map[string]interface{})
				d["field"] = "domicile"
				d["message"] = "Missing required parameter: domicile"
				dataError = append(dataError, d)
			}
		}
		if oadata.CorrespondenceKey == nil {
			d := make(map[string]interface{})
			d["field"] = "correspondence"
			d["message"] = "Missing required parameter: correspondence"
			dataError = append(dataError, d)
		} else {
			var corr models.OaPostalAddress
			_, err = models.GetOaPostalAddress(&corr, strconv.FormatUint(*oadata.CorrespondenceKey, 10))
			if err == nil {
				if corr.AddressLine1 == nil {
					d := make(map[string]interface{})
					d["field"] = "correspondence_address"
					d["message"] = "Missing required parameter: correspondence_address"
					dataError = append(dataError, d)
				}
				if corr.KecamatanKey == nil {
					d := make(map[string]interface{})
					d["field"] = "correspondence_subdistric"
					d["message"] = "Missing required parameter: correspondence_subdistric"
					dataError = append(dataError, d)
				}
				if corr.KabupatenKey == nil {
					d := make(map[string]interface{})
					d["field"] = "correspondence_city"
					d["message"] = "Missing required parameter: correspondence_city"
					dataError = append(dataError, d)
				}
				if corr.PostalCode == nil {
					d := make(map[string]interface{})
					d["field"] = "correspondence_postalcode"
					d["message"] = "Missing required parameter: correspondence_postalcode"
					dataError = append(dataError, d)
				}
			} else {
				d := make(map[string]interface{})
				d["field"] = "correspondence"
				d["message"] = "Missing required parameter: correspondence"
				dataError = append(dataError, d)
			}
		}
		if oadata.IntitutionType == nil {
			d := make(map[string]interface{})
			d["field"] = "intitution_type"
			d["message"] = "Missing required parameter: intitution_type"
			dataError = append(dataError, d)
		}
		if oadata.IntitutionClassification == nil {
			d := make(map[string]interface{})
			d["field"] = "intitution_classification"
			d["message"] = "Missing required parameter: intitution_classification"
			dataError = append(dataError, d)
		}
		if oadata.IntitutionCharacteristic == nil {
			d := make(map[string]interface{})
			d["field"] = "intitution_characteristic"
			d["message"] = "Missing required parameter: intitution_characteristic"
			dataError = append(dataError, d)
		}
		if oadata.IntitutionBusinessType == nil {
			d := make(map[string]interface{})
			d["field"] = "intitution_business_type"
			d["message"] = "Missing required parameter: intitution_business_type"
			dataError = append(dataError, d)
		}
		if oadata.InstiAnnuallyIncome == nil {
			d := make(map[string]interface{})
			d["field"] = "insti_annually_income"
			d["message"] = "Missing required parameter: insti_annually_income"
			dataError = append(dataError, d)
		}
		if oadata.InstiSourceOfIncome == nil {
			d := make(map[string]interface{})
			d["field"] = "insti_source_of_income"
			d["message"] = "Missing required parameter: insti_source_of_income"
			dataError = append(dataError, d)
		}
		if oadata.InstiInvestmentPurpose == nil {
			d := make(map[string]interface{})
			d["field"] = "insti_investment_purpose"
			d["message"] = "Missing required parameter: insti_investment_purpose"
			dataError = append(dataError, d)
		}
		if oadata.BoName == nil {
			d := make(map[string]interface{})
			d["field"] = "bo_name"
			d["message"] = "Missing required parameter: bo_name"
			dataError = append(dataError, d)
		}
		if oadata.BoIdnumber == nil {
			d := make(map[string]interface{})
			d["field"] = "bo_idnumber"
			d["message"] = "Missing required parameter: bo_idnumber"
			dataError = append(dataError, d)
		}
		if oadata.BoBusiness == nil {
			d := make(map[string]interface{})
			d["field"] = "bo_business"
			d["message"] = "Missing required parameter: bo_business"
			dataError = append(dataError, d)
		}
		if oadata.BoAnnuallyIncome == nil {
			d := make(map[string]interface{})
			d["field"] = "bo_annually_income"
			d["message"] = "Missing required parameter: bo_annually_income"
			dataError = append(dataError, d)
		}
		if oadata.BoRelation == nil {
			d := make(map[string]interface{})
			d["field"] = "bo_relation"
			d["message"] = "Missing required parameter: bo_relation"
			dataError = append(dataError, d)
		}
		if oadata.BoRelation == nil {
			d := make(map[string]interface{})
			d["field"] = "bo_relation"
			d["message"] = "Missing required parameter: bo_relation"
			dataError = append(dataError, d)
		}

		//SHARED HOLDER
		var sharedHolder []models.OaInstitutionSharesHolderDetail
		_, err = models.GetOaInstitutionSharesHolderRequest(&sharedHolder, oaRequestKey)
		if err == nil && len(sharedHolder) > 0 {
			ky := 1
			for _, shar := range sharedHolder {
				if shar.Nationality == nil {
					d := make(map[string]interface{})
					d["field"] = "shares_holder key " + strconv.FormatUint(uint64(ky), 10) + " nationality"
					d["message"] = "Missing required parameter: shares_holder key " + strconv.FormatUint(uint64(ky), 10) + " nationality"
					dataError = append(dataError, d)
				}
				if shar.HolderFullName == nil {
					d := make(map[string]interface{})
					d["field"] = "shares_holder key " + strconv.FormatUint(uint64(ky), 10) + " holder_full_name"
					d["message"] = "Missing required parameter: shares_holder key " + strconv.FormatUint(uint64(ky), 10) + " holder_full_name"
					dataError = append(dataError, d)
				}
				if shar.IdcardType == nil {
					d := make(map[string]interface{})
					d["field"] = "shares_holder key " + strconv.FormatUint(uint64(ky), 10) + " idcard_type"
					d["message"] = "Missing required parameter: shares_holder key " + strconv.FormatUint(uint64(ky), 10) + " idcard_type"
					dataError = append(dataError, d)
				}
				if shar.IdcardNo == nil {
					d := make(map[string]interface{})
					d["field"] = "shares_holder key " + strconv.FormatUint(uint64(ky), 10) + " idcard_no"
					d["message"] = "Missing required parameter: shares_holder key " + strconv.FormatUint(uint64(ky), 10) + " idcard_no"
					dataError = append(dataError, d)
				}
				if shar.HolderDob == nil {
					d := make(map[string]interface{})
					d["field"] = "shares_holder key " + strconv.FormatUint(uint64(ky), 10) + " holder_dob"
					d["message"] = "Missing required parameter: shares_holder key " + strconv.FormatUint(uint64(ky), 10) + " holder_dob"
					dataError = append(dataError, d)
				}
				if shar.SharesPercent == nil {
					d := make(map[string]interface{})
					d["field"] = "shares_holder key " + strconv.FormatUint(uint64(ky), 10) + " shares_percent"
					d["message"] = "Missing required parameter: shares_holder key " + strconv.FormatUint(uint64(ky), 10) + " shares_percent"
					dataError = append(dataError, d)
				}

				ky++
			}
		} else {
			d := make(map[string]interface{})
			d["field"] = "shares_holder"
			d["message"] = "Missing required parameter: shares_holder"
			dataError = append(dataError, d)
		}

		//AUTH PERSON
		var authPerson []models.OaInstitutionAuthPersonDetail
		_, err = models.GetOaInstitutionAuthPersonRequest(&authPerson, oaRequestKey)
		if err == nil && len(authPerson) > 0 {
			ky := 1
			for _, authP := range authPerson {
				if authP.FullName == nil {
					d := make(map[string]interface{})
					d["field"] = "auth_person key " + strconv.FormatUint(uint64(ky), 10) + " full_name"
					d["message"] = "Missing required parameter: auth_person key " + strconv.FormatUint(uint64(ky), 10) + " full_name"
					dataError = append(dataError, d)
				}
				if authP.PersonDob == nil {
					d := make(map[string]interface{})
					d["field"] = "auth_person key " + strconv.FormatUint(uint64(ky), 10) + " person_dob"
					d["message"] = "Missing required parameter: auth_person key " + strconv.FormatUint(uint64(ky), 10) + " person_dob"
					dataError = append(dataError, d)
				}
				if authP.IdcardType == nil {
					d := make(map[string]interface{})
					d["field"] = "auth_person key " + strconv.FormatUint(uint64(ky), 10) + " idcard_type"
					d["message"] = "Missing required parameter: auth_person key " + strconv.FormatUint(uint64(ky), 10) + " idcard_type"
					dataError = append(dataError, d)
				}
				if authP.IdcardNo == nil {
					d := make(map[string]interface{})
					d["field"] = "auth_person key " + strconv.FormatUint(uint64(ky), 10) + " idcard_no"
					d["message"] = "Missing required parameter: auth_person key " + strconv.FormatUint(uint64(ky), 10) + " idcard_no"
					dataError = append(dataError, d)
				}
				if authP.Nationality == nil {
					d := make(map[string]interface{})
					d["field"] = "auth_person key " + strconv.FormatUint(uint64(ky), 10) + " nationality"
					d["message"] = "Missing required parameter: auth_person key " + strconv.FormatUint(uint64(ky), 10) + " nationality"
					dataError = append(dataError, d)
				}
				if authP.Position == nil {
					d := make(map[string]interface{})
					d["field"] = "auth_person key " + strconv.FormatUint(uint64(ky), 10) + " position"
					d["message"] = "Missing required parameter: auth_person key " + strconv.FormatUint(uint64(ky), 10) + " position"
					dataError = append(dataError, d)
				}
				if authP.PhoneNo == nil {
					d := make(map[string]interface{})
					d["field"] = "auth_person key " + strconv.FormatUint(uint64(ky), 10) + " phone_no"
					d["message"] = "Missing required parameter: auth_person key " + strconv.FormatUint(uint64(ky), 10) + " phone_no"
					dataError = append(dataError, d)
				}
				if authP.EmailAddress == nil {
					d := make(map[string]interface{})
					d["field"] = "auth_person key " + strconv.FormatUint(uint64(ky), 10) + " email_address"
					d["message"] = "Missing required parameter: auth_person key " + strconv.FormatUint(uint64(ky), 10) + " email_address"
					dataError = append(dataError, d)
				}

				ky++
			}
		} else {
			d := make(map[string]interface{})
			d["field"] = "auth_person"
			d["message"] = "Missing required parameter: auth_person"
			dataError = append(dataError, d)
		}

		if oadata.AssetY1 == nil {
			d := make(map[string]interface{})
			d["field"] = "asset_y1"
			d["message"] = "Missing required parameter: asset_y1"
			dataError = append(dataError, d)
		}
		if oadata.OpsProfitY3 == nil {
			d := make(map[string]interface{})
			d["field"] = "ops_profit_y3"
			d["message"] = "Missing required parameter: ops_profit_y3"
			dataError = append(dataError, d)
		}

		//BANK ACCOUNT
		var accBank []models.OaRequestByField
		_, err = models.GetOaRequestBankByField(&accBank, "oa_request_key", strconv.FormatUint(oareq.OaRequestKey, 10))
		if err == nil && len(accBank) > 0 {
			ky := 1
			prio := 0
			for _, b := range accBank {
				if b.AccountNo == nil {
					d := make(map[string]interface{})
					d["field"] = "bank_account key " + strconv.FormatUint(uint64(ky), 10) + " account_no"
					d["message"] = "Missing required parameter: bank_account key " + strconv.FormatUint(uint64(ky), 10) + " account_no"
					dataError = append(dataError, d)
				}
				if b.AccountHolderName == nil {
					d := make(map[string]interface{})
					d["field"] = "bank_account key " + strconv.FormatUint(uint64(ky), 10) + " account_holder_name"
					d["message"] = "Missing required parameter: bank_account key " + strconv.FormatUint(uint64(ky), 10) + " account_holder_name"
					dataError = append(dataError, d)
				}
				if b.BranchName == nil {
					d := make(map[string]interface{})
					d["field"] = "bank_account key " + strconv.FormatUint(uint64(ky), 10) + " branch_name"
					d["message"] = "Missing required parameter: bank_account key " + strconv.FormatUint(uint64(ky), 10) + " branch_name"
					dataError = append(dataError, d)
				}
				if b.CurrencyKey == nil {
					d := make(map[string]interface{})
					d["field"] = "bank_account key " + strconv.FormatUint(uint64(ky), 10) + " currency_key"
					d["message"] = "Missing required parameter: bank_account key " + strconv.FormatUint(uint64(ky), 10) + " currency_key"
					dataError = append(dataError, d)
				}

				if b.FlagPriority == 1 {
					prio++
				}

				ky++
			}

			if prio != 1 {
				d := make(map[string]interface{})
				d["field"] = "bank_account field flag_priority"
				d["message"] = "set 1 data bank_account flag_priority"
				dataError = append(dataError, d)
			}
		} else {
			d := make(map[string]interface{})
			d["field"] = "bank_account"
			d["message"] = "Missing required parameter: bank_account"
			dataError = append(dataError, d)
		}
	}

	//DOCS
	var docsError []interface{}
	var instDocs []models.OaInstitutionDocsDetail
	_, err = models.GetOaInstitutionDocsRequest(&instDocs, oaRequestKey)
	if err == nil && len(instDocs) > 0 {
		for _, id := range instDocs {
			if id.InstiDocsKey == nil {
				d := make(map[string]interface{})
				d["field"] = "document"
				d["message"] = "Missing required parameter dokumen: " + *id.InstiDocumentTypeName
				docsError = append(docsError, d)
			}
		}
	} else {
		d := make(map[string]interface{})
		d["field"] = "document"
		d["message"] = "Lengkapi dokumen institusi anda."
		docsError = append(docsError, d)
	}

	//RISK
	var riskError []interface{}
	var oaRiskProfileQuiz []models.AdminOaRiskProfileQuiz
	_, err = models.AdminGetOaRiskProfileQuizByOaRequestKey(&oaRiskProfileQuiz, strconv.FormatUint(oareq.OaRequestKey, 10))
	if err != nil {
		d := make(map[string]interface{})
		d["field"] = "risk_profile"
		d["message"] = "Lengkapi Risk Profile Risiko anda."
		riskError = append(riskError, d)
	} else {
		if len(oaRiskProfileQuiz) != 6 {
			d := make(map[string]interface{})
			d["field"] = "risk_profile"
			d["message"] = "Lengkapi Risk Profile Risiko anda."
			riskError = append(riskError, d)
		}
	}

	//USER
	var userError []interface{}
	// InstitutionUserMaker
	var userMaker []models.OaInstitutionUserDetail
	_, err = models.GetOaInstitutionUserRequest(&userMaker, oaRequestKey, lib.ROLE_INSTITUTION_MAKER)
	if err != nil {
		d := make(map[string]interface{})
		d["field"] = "institution_user_maker"
		d["message"] = "Lengkapi Institution User Maker anda."
		userError = append(userError, d)
	} else {
		if len(userMaker) == 0 {
			d := make(map[string]interface{})
			d["field"] = "institution_user_maker"
			d["message"] = "Lengkapi Institution User Maker anda."
			userError = append(userError, d)
		}
		if len(userMaker) > 2 {
			d := make(map[string]interface{})
			d["field"] = "institution_user_maker"
			d["message"] = "Institution User Maker hanya max. 2 data."
			userError = append(userError, d)
		}
		ky := 1
		prio := 0
		for _, um := range userMaker {
			if um.FullName == nil {
				d := make(map[string]interface{})
				d["field"] = "institution_user_maker key " + strconv.FormatUint(uint64(ky), 10) + " full_name"
				d["message"] = "Missing required parameter: institution_user_maker key " + strconv.FormatUint(uint64(ky), 10) + " full_name"
				userError = append(userError, d)
			}
			if um.EmailAddress == nil {
				d := make(map[string]interface{})
				d["field"] = "institution_user_maker key " + strconv.FormatUint(uint64(ky), 10) + " email_address"
				d["message"] = "Missing required parameter: institution_user_maker key " + strconv.FormatUint(uint64(ky), 10) + " email_address"
				userError = append(userError, d)
			}
			if um.PhoneNumber == nil {
				d := make(map[string]interface{})
				d["field"] = "institution_user_maker key " + strconv.FormatUint(uint64(ky), 10) + " phone_number"
				d["message"] = "Missing required parameter: institution_user_maker key " + strconv.FormatUint(uint64(ky), 10) + " phone_number"
				userError = append(userError, d)
			}
			if um.UserPriority == 1 {
				prio++
			}

			ky++
		}

		if len(userMaker) > 0 {
			if prio != 1 {
				d := make(map[string]interface{})
				d["field"] = "institution_user_maker field flag_priority"
				d["message"] = "set 1 data institution_user_maker flag_priority"
				userError = append(userError, d)
			}
		}
	}

	// InstitutionUserChecker
	var userChecker []models.OaInstitutionUserDetail
	_, err = models.GetOaInstitutionUserRequest(&userChecker, oaRequestKey, lib.ROLE_INSTITUTION_CHECKER)
	if err != nil {
		d := make(map[string]interface{})
		d["field"] = "institution_user_checker"
		d["message"] = "Lengkapi Institution User Maker anda."
		userError = append(userError, d)
	} else {
		if len(userChecker) == 0 {
			d := make(map[string]interface{})
			d["field"] = "institution_user_checker"
			d["message"] = "Lengkapi Institution User Checker anda."
			userError = append(userError, d)
		}
		if len(userChecker) > 2 {
			d := make(map[string]interface{})
			d["field"] = "institution_user_checker"
			d["message"] = "Institution User Checker hanya max. 2 data."
			userError = append(userError, d)
		}
		ky := 1
		prio := 0
		for _, um := range userChecker {
			if um.FullName == nil {
				d := make(map[string]interface{})
				d["field"] = "institution_user_checker key " + strconv.FormatUint(uint64(ky), 10) + " full_name"
				d["message"] = "Missing required parameter: institution_user_checker key " + strconv.FormatUint(uint64(ky), 10) + " full_name"
				userError = append(userError, d)
			}
			if um.EmailAddress == nil {
				d := make(map[string]interface{})
				d["field"] = "institution_user_checker key " + strconv.FormatUint(uint64(ky), 10) + " email_address"
				d["message"] = "Missing required parameter: institution_user_checker key " + strconv.FormatUint(uint64(ky), 10) + " email_address"
				userError = append(userError, d)
			}
			if um.PhoneNumber == nil {
				d := make(map[string]interface{})
				d["field"] = "institution_user_checker key " + strconv.FormatUint(uint64(ky), 10) + " phone_number"
				d["message"] = "Missing required parameter: institution_user_checker key " + strconv.FormatUint(uint64(ky), 10) + " phone_number"
				userError = append(userError, d)
			}
			if um.UserPriority == 1 {
				prio++
			}

			ky++
		}

		if len(userChecker) > 0 {
			if prio != 1 {
				d := make(map[string]interface{})
				d["field"] = "institution_user_checker field flag_priority"
				d["message"] = "set 1 data institution_user_checker flag_priority"
				userError = append(userError, d)
			}
		}
	}

	// InstitutionUserReleaser
	var userReleaser []models.OaInstitutionUserDetail
	_, err = models.GetOaInstitutionUserRequest(&userReleaser, oaRequestKey, lib.ROLE_INSTITUTION_RELEASER)
	if err != nil {
		d := make(map[string]interface{})
		d["field"] = "institution_user_releaser"
		d["message"] = "Lengkapi Institution User Maker anda."
		userError = append(userError, d)
	} else {
		if len(userReleaser) == 0 {
			d := make(map[string]interface{})
			d["field"] = "institution_user_releaser"
			d["message"] = "Lengkapi Institution User Releaser anda."
			userError = append(userError, d)
		}
		if len(userReleaser) > 2 {
			d := make(map[string]interface{})
			d["field"] = "institution_user_releaser"
			d["message"] = "Institution User Releaser hanya max. 2 data."
			userError = append(userError, d)
		}
		ky := 1
		prio := 0
		for _, um := range userReleaser {
			if um.FullName == nil {
				d := make(map[string]interface{})
				d["field"] = "institution_user_releaser key " + strconv.FormatUint(uint64(ky), 10) + " full_name"
				d["message"] = "Missing required parameter: institution_user_releaser key " + strconv.FormatUint(uint64(ky), 10) + " full_name"
				userError = append(userError, d)
			}
			if um.EmailAddress == nil {
				d := make(map[string]interface{})
				d["field"] = "institution_user_releaser key " + strconv.FormatUint(uint64(ky), 10) + " email_address"
				d["message"] = "Missing required parameter: institution_user_releaser key " + strconv.FormatUint(uint64(ky), 10) + " email_address"
				userError = append(userError, d)
			}
			if um.PhoneNumber == nil {
				d := make(map[string]interface{})
				d["field"] = "institution_user_releaser key " + strconv.FormatUint(uint64(ky), 10) + " phone_number"
				d["message"] = "Missing required parameter: institution_user_releaser key " + strconv.FormatUint(uint64(ky), 10) + " phone_number"
				userError = append(userError, d)
			}
			if um.UserPriority == 1 {
				prio++
			}

			ky++
		}
		if len(userReleaser) > 0 {
			if prio != 1 {
				d := make(map[string]interface{})
				d["field"] = "institution_user_releaser field flag_priority"
				d["message"] = "set 1 data institution_user_releaser flag_priority"
				userError = append(userError, d)
			}
		}
	}

	responseData := make(map[string]interface{})
	responseData["data"] = dataError
	responseData["docs"] = docsError
	responseData["risk_profile"] = riskError
	responseData["user"] = userError

	var response lib.Response
	if len(dataError) > 0 || len(docsError) > 0 || len(riskError) > 0 || len(userError) > 0 {
		response.Status.Code = http.StatusBadRequest
		msg := "Lengkapi"

		error := 0
		if len(dataError) > 0 {
			msg += " Data Institusi"
			error++
		}

		if len(docsError) > 0 {
			if error > 0 {
				msg += ", Dokumen Institusi"
			} else {
				msg += " Dokumen Institusi"
			}
			error++
		}

		if len(riskError) > 0 {
			if error > 0 {
				msg += ", Risk Profil Risiko"
			} else {
				msg += " Risk Profil Risiko"
			}
			error++
		}

		if len(userError) > 0 {
			if error > 0 {
				msg += ", User Institusi"
			} else {
				msg += " User Institusi"
			}
			error++
		}

		msg += " anda."

		response.Status.MessageServer = msg
		response.Status.MessageClient = msg
	} else {
		layout := "2006-01-02 15:04:05"

		params := make(map[string]string)
		params["oa_request_key"] = oaRequestKey
		params["oa_entry_end"] = time.Now().Format(layout)
		params["oa_status"] = "258"
		params["rec_modified_date"] = time.Now().Format(layout)
		params["rec_modified_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
		params["rec_attribute_id3"] = c.Request().UserAgent()

		_, err = models.UpdateOaRequest(params)
		if err != nil {
			log.Error("Error update oa request")
			return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed update data")
		}

		//kirim email ke CS Approval
		//kirim email ke Sales
		SentEmailInstitusiOaPengkinianToBackOfficeSales(oareq, oadata, "11", true)

		response.Status.Code = http.StatusOK
		response.Status.MessageServer = "OK"
		response.Status.MessageClient = "OK"
	}

	response.Data = responseData

	return c.JSON(http.StatusOK, response)
}
