package controllers

import (
	"bytes"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"html/template"
	"math"
	"mf-bo-api/config"
	"mf-bo-api/db"
	"mf-bo-api/lib"
	"mf-bo-api/models"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/badoux/checkmail"
	"github.com/labstack/echo"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"

	wkhtml "github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

func GetListCustomerIndividuInquiry(c echo.Context) error {

	var err error
	var status int
	decimal.MarshalJSONWithoutQuotes = true

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

	items := []string{"oa_request_key", "cif", "full_name", "sid", "date_birth", "customer_key", "phone_mobile", "cif_suspend_flag", "mother_maiden_name", "ktp"}

	params := make(map[string]string)
	orderBy := c.QueryParam("order_by")
	if orderBy != "" {
		_, found := lib.Find(items, orderBy)
		if found {
			params["orderBy"] = orderBy
			orderType := c.QueryParam("order_type")
			if (orderType == "asc") || (orderType == "ASC") || (orderType == "desc") || (orderType == "DESC") {
				params["orderType"] = orderType
			}
		} else {
			log.Error("Wrong input for parameter order_by")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter order_by", "Wrong input for parameter order_by")
		}
	} else {
		params["orderBy"] = "oa_request_key"
		params["orderType"] = "DESC"
	}

	params["c.investor_type"] = "263"

	paramsLike := make(map[string]string)

	cif := c.QueryParam("cif")
	if cif != "" {
		paramsLike["c.unit_holder_idno"] = cif
	}
	fullname := c.QueryParam("full_name")
	if fullname != "" {
		paramsLike["pd.full_name"] = fullname
	}
	datebirth := c.QueryParam("date_birth")
	if datebirth != "" {
		paramsLike["pd.date_birth"] = datebirth
	}
	ktp := c.QueryParam("ktp")
	if ktp != "" {
		paramsLike["pd.idcard_no"] = ktp
	}
	mothermaidenname := c.QueryParam("mother_maiden_name")
	if mothermaidenname != "" {
		paramsLike["pd.mother_maiden_name"] = mothermaidenname
	}
	branchKey := c.QueryParam("branch_key")
	if branchKey != "" {
		params["r.branch_key"] = branchKey
	}

	//if user category  = 3 -> user branch, 2 = user HO
	var userCategory uint64
	userCategory = 3
	if lib.Profile.UserCategoryKey == userCategory {
		log.Println(lib.Profile)
		if lib.Profile.BranchKey != nil {
			strBranchKey := strconv.FormatUint(*lib.Profile.BranchKey, 10)
			params["r.branch_key"] = strBranchKey
		} else {
			log.Error("User Branch haven't Branch")
			return lib.CustomError(http.StatusBadRequest, "Wrong User Branch haven't Branch", "Wrong User Branch haven't Branch")
		}
	}

	var customers []models.CustomerIndividuInquiry

	status, err = models.AdminGetAllCustomerIndividuInquery(&customers, limit, offset, params, paramsLike, noLimit)

	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if len(customers) < 1 {
		log.Error("Customer not found")
		return lib.CustomError(http.StatusNotFound, "Customer not found", "Customer not found")
	}

	var countData models.CountData
	var pagination int
	if limit > 0 {
		status, err = models.CountAdminGetAllCustomerIndividuInquery(&countData, params, paramsLike)
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
	response.Data = customers

	return c.JSON(http.StatusOK, response)
}

func GetListCustomerInstitutionInquiry(c echo.Context) error {

	var err error
	var status int
	decimal.MarshalJSONWithoutQuotes = true

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

	items := []string{"customer_key", "full_name", "sid", "npwp", "cif", "institution", "cif_suspend_flag"}

	params := make(map[string]string)
	orderBy := c.QueryParam("order_by")
	if orderBy != "" {
		_, found := lib.Find(items, orderBy)
		if found {

			var orderByJoin string
			orderByJoin = "c.customer_key"
			if orderBy == "cif" {
				orderByJoin = "c.unit_holder_idno"
			} else if orderBy == "full_name" {
				orderByJoin = "c.full_name"
			} else if orderBy == "sid" {
				orderByJoin = "c.sid"
			} else if orderBy == "institution" {
				orderByJoin = "pd.insti_full_name"
			} else if orderBy == "cif_suspend_flag" {
				orderByJoin = "c.cif_suspend_flag"
			} else if orderBy == "npwp" {
				orderByJoin = "pd.npwp_no"
			} else if orderBy == "ktp" {
				orderByJoin = "pd.idcard_no"
			}

			params["orderBy"] = orderByJoin
			orderType := c.QueryParam("order_type")
			if (orderType == "asc") || (orderType == "ASC") || (orderType == "desc") || (orderType == "DESC") {
				params["orderType"] = orderType
			}
		} else {
			log.Error("Wrong input for parameter order_by")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter order_by", "Wrong input for parameter order_by")
		}
	} else {
		params["orderBy"] = "c.customer_key"
		params["orderType"] = "DESC"
	}

	params["c.investor_type"] = "264"

	paramsLike := make(map[string]string)

	cif := c.QueryParam("cif")
	if cif != "" {
		paramsLike["c.unit_holder_idno"] = cif
	}
	fullname := c.QueryParam("full_name")
	if fullname != "" {
		paramsLike["c.full_name"] = fullname
	}
	npwp := c.QueryParam("npwp")
	if npwp != "" {
		paramsLike["pd.npwp_no"] = npwp
	}

	//if user category  = 3 -> user branch, 2 = user HO
	var userCategory uint64
	userCategory = 3
	if lib.Profile.UserCategoryKey == userCategory {
		log.Println(lib.Profile)
		if lib.Profile.BranchKey != nil {
			strBranchKey := strconv.FormatUint(*lib.Profile.BranchKey, 10)
			params["c.openacc_branch_key"] = strBranchKey
		} else {
			log.Error("User Branch haven't Branch")
			return lib.CustomError(http.StatusBadRequest, "Wrong User Branch haven't Branch", "Wrong User Branch haven't Branch")
		}
	}

	var customers []models.CustomerInstituionInquiry

	status, err = models.AdminGetAllCustomerInstitutionInquery(&customers, limit, offset, params, paramsLike, noLimit)

	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if len(customers) < 1 {
		log.Error("Customer not found")
		return lib.CustomError(http.StatusNotFound, "Customer not found", "Customer not found")
	}

	var countData models.CountData
	var pagination int
	if limit > 0 {
		status, err = models.CountAdminGetAllCustomerInstitutionInquery(&countData, params, paramsLike)
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
	response.Data = customers

	return c.JSON(http.StatusOK, response)
}

func GetDetailCustomerIndividu(c echo.Context) error {
	var err error

	keyStr := c.Param("key") //oa_request_key
	if keyStr == "" {
		log.Error("Missing required parameter: key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: key", "Missing required parameter: key")
	}
	key, _ := strconv.ParseUint(keyStr, 10, 64)
	if key == 0 {
		return lib.CustomError(http.StatusNotFound)
	}

	var customer models.CustomerIndividuInquiry
	_, err = models.AdminGetHeaderCustomerIndividu(&customer, keyStr)
	if err != nil {
		return lib.CustomError(http.StatusNotFound)
	}

	//if user category  = 3 -> user branch, 2 = user HO
	var userCategory uint64
	userCategory = 3
	if lib.Profile.UserCategoryKey == userCategory {
		if customer.BranchKey == nil {
			log.Error("User Branch null, not match with customer branch")
			return lib.CustomError(http.StatusNotFound)
		} else {
			strCusBranch := strconv.FormatUint(*customer.BranchKey, 10)
			strUserBranch := strconv.FormatUint(*lib.Profile.BranchKey, 10)
			if strCusBranch != strUserBranch {
				log.Error("User Branch not match with customer branch")
				return lib.CustomError(http.StatusNotFound)
			}
		}
	}

	var oaCustomer []models.OaCustomer
	if customer.CustomerKey != nil {
		_, err = models.AdminGetAllOaByCustomerKey(&oaCustomer, strconv.FormatUint(*customer.CustomerKey, 10))
		if err != nil {
			if err != sql.ErrNoRows {
				return lib.CustomError(http.StatusNotFound)
			}
		}
	} else {
		_, err = models.AdminGetAllOaByOaKey(&oaCustomer, strconv.FormatUint(customer.OaRequestKey, 10))
		if err != nil {
			if err != sql.ErrNoRows {
				return lib.CustomError(http.StatusNotFound)
			}
		}

	}

	var responseData models.DetailCustomerIndividuInquiry
	responseData.Header = customer
	responseData.PersonalData = &oaCustomer

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = responseData

	return c.JSON(http.StatusOK, response)
}

func GetDetailCustomerInstitution(c echo.Context) error {
	var err error

	keyStr := c.Param("key")
	if keyStr == "" {
		log.Error("Missing required parameter: key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: key", "Missing required parameter: key")
	}
	key, _ := strconv.ParseUint(keyStr, 10, 64)
	if key == 0 {
		return lib.CustomError(http.StatusNotFound)
	}

	var customer models.CustomerInstituionInquiry
	_, err = models.AdminGetHeaderCustomerInstitution(&customer, keyStr)
	if err != nil {
		return lib.CustomError(http.StatusNotFound)
	}

	//if user category  = 3 -> user branch, 2 = user HO
	var userCategory uint64
	userCategory = 3
	if lib.Profile.UserCategoryKey == userCategory {
		if customer.BranchKey == nil {
			log.Error("User Branch null, not match with customer branch")
			return lib.CustomError(http.StatusNotFound)
		} else {
			strCusBranch := strconv.FormatUint(*customer.BranchKey, 10)
			strUserBranch := strconv.FormatUint(*lib.Profile.BranchKey, 10)
			if strCusBranch != strUserBranch {
				log.Error("User Branch not match with customer branch")
				return lib.CustomError(http.StatusNotFound)
			}
		}
	}

	var oaCustomer []models.OaCustomer
	_, err = models.AdminGetAllOaByCustomerKey(&oaCustomer, keyStr)
	if err != nil {
		if err != sql.ErrNoRows {
			return lib.CustomError(http.StatusNotFound)
		}
	}

	var responseData models.DetailCustomerInstitutionInquiry
	responseData.Header = customer
	responseData.PersonalData = &oaCustomer

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = responseData

	return c.JSON(http.StatusOK, response)
}

func GetDetailCustomerInquiry(c echo.Context) error {
	var err error

	keyStr := c.Param("key")
	if keyStr == "" {
		log.Error("Missing required parameter: key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: key", "Missing required parameter: key")
	}
	key, _ := strconv.ParseUint(keyStr, 10, 64)
	if key == 0 {
		return lib.CustomError(http.StatusNotFound)
	}

	var customer models.DetailCustomerInquiry
	_, err = models.AdminGetHeaderDetailCustomer(&customer, keyStr)
	if err != nil {
		return lib.CustomError(http.StatusNotFound)
	}

	var customerHeader models.DetailHeaderCustomerInquiry

	customerHeader.InvestorType = customer.InvestorType
	customerHeader.CustomerKey = customer.CustomerKey
	customerHeader.Cif = customer.Cif
	customerHeader.FullName = customer.FullName
	customerHeader.SidNo = customer.SidNo
	customerHeader.CifSuspendFlag = customer.CifSuspendFlag

	if customer.InvestorType == "263" {
		customerHeader.DateBirth = customer.DateBirth
		customerHeader.IDcardNo = customer.IDcardNo
		customerHeader.PhoneMobile = customer.PhoneMobile
		customerHeader.MotherMaidenName = customer.MotherMaidenName
	} else {
		customerHeader.Npwp = customer.Npwp
		customerHeader.Institusion = customer.Institusion
	}

	var oaCustomer []models.OaCustomer
	_, err = models.AdminGetAllOaByCustomerKey(&oaCustomer, keyStr)
	if err != nil {
		if err != sql.ErrNoRows {
			return lib.CustomError(http.StatusNotFound)
		}
	}

	var responseData models.DetailCustomerInquiryResponse
	responseData.Header = customerHeader
	responseData.PersonalData = &oaCustomer

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = responseData

	return c.JSON(http.StatusOK, response)
}

func DetailPersonalDataCustomerIndividu(c echo.Context) error {
	var err error
	var status int
	decimal.MarshalJSONWithoutQuotes = true
	//Get parameter limit
	keyStr := c.Param("key")
	key, _ := strconv.ParseUint(keyStr, 10, 64)
	if key == 0 {
		return lib.CustomError(http.StatusNotFound)
	}

	var oareq models.OaRequest
	status, err = models.GetOaRequest(&oareq, keyStr)
	if err != nil {
		return lib.CustomError(status)
	}

	// if oareq.CustomerKey == nil {
	// 	log.Println("data belum jadi customer")
	// 	return lib.CustomError(http.StatusNotFound)
	// }

	var responseData models.DetailPersonalDataCustomerIndividu

	layout := "2006-01-02 15:04:05"
	newLayout := "02 Jan 2006"

	responseData.OaRequestKey = oareq.OaRequestKey
	date, _ := time.Parse(layout, oareq.OaEntryStart)
	responseData.OaEntryStart = date.Format(newLayout)
	date, _ = time.Parse(layout, oareq.OaEntryEnd)
	responseData.OaEntryEnd = date.Format(newLayout)
	if oareq.SalesCode != nil {
		responseData.SalesCode = *oareq.SalesCode
	} else {
		responseData.SalesCode = ""
	}

	var oaRequestLookupIds []string

	if oareq.OaRequestType != nil {
		if _, ok := lib.Find(oaRequestLookupIds, strconv.FormatUint(*oareq.OaRequestType, 10)); !ok {
			oaRequestLookupIds = append(oaRequestLookupIds, strconv.FormatUint(*oareq.OaRequestType, 10))
		}
	}
	if oareq.Oastatus != nil {
		if _, ok := lib.Find(oaRequestLookupIds, strconv.FormatUint(*oareq.Oastatus, 10)); !ok {
			oaRequestLookupIds = append(oaRequestLookupIds, strconv.FormatUint(*oareq.Oastatus, 10))
		}
	}
	if oareq.OaRiskLevel != nil {
		if _, ok := lib.Find(oaRequestLookupIds, strconv.FormatUint(*oareq.OaRiskLevel, 10)); !ok {
			oaRequestLookupIds = append(oaRequestLookupIds, strconv.FormatUint(*oareq.OaRiskLevel, 10))
		}
	}

	//gen lookup oa request
	var lookupOaReq []models.GenLookup
	if len(oaRequestLookupIds) > 0 {
		status, err = models.GetGenLookupIn(&lookupOaReq, oaRequestLookupIds, "lookup_key")
		if err != nil {
			if err != sql.ErrNoRows {
				log.Error(err.Error())
				return lib.CustomError(status, err.Error(), "Failed get data")
			}
		}
	}

	gData := make(map[uint64]models.GenLookup)
	for _, gen := range lookupOaReq {
		gData[gen.LookupKey] = gen
	}

	if oareq.OaRequestType != nil {
		if n, ok := gData[*oareq.OaRequestType]; ok {
			responseData.OaRequestType = n.LkpName
		}
	}

	if oareq.OaRiskLevel != nil {
		if n, ok := gData[*oareq.OaRiskLevel]; ok {
			responseData.OaRiskLevel = n.LkpName
		}
	}

	if oareq.Oastatus != nil {
		if n, ok := gData[*oareq.Oastatus]; ok {
			responseData.Oastatus = *n.LkpName
		}
	}

	//check personal data by oa request key
	var oapersonal models.OaPersonalData
	strKey := strconv.FormatUint(oareq.OaRequestKey, 10)
	status, err = models.GetOaPersonalDataByOaRequestKey(&oapersonal, strKey)
	if err != nil {
		if err != sql.ErrNoRows {
			return lib.CustomError(status)
		}
	} else {
		responseData.FullName = oapersonal.FullName
		responseData.IDCardNo = oapersonal.IDcardNo
		date, _ = time.Parse(layout, oapersonal.DateBirth)
		responseData.DateBirth = date.Format(newLayout)
		responseData.PhoneNumber = oapersonal.PhoneMobile
		responseData.EmailAddress = oapersonal.EmailAddress
		responseData.PlaceBirth = oapersonal.PlaceBirth
		responseData.PhoneHome = oapersonal.PhoneHome
		responseData.RelationFullName = oapersonal.RelationFullName

		dir := config.ImageUrl + "/images/user/" + strconv.FormatUint(*oareq.UserLoginKey, 10) + "/"

		if oapersonal.PicKtp != nil && *oapersonal.PicKtp != "" {
			path := dir + *oapersonal.PicKtp
			responseData.PicKtp = &path
		}

		if oapersonal.PicSelfie != nil && *oapersonal.PicSelfie != "" {
			path := dir + *oapersonal.PicSelfie
			responseData.PicSelfie = &path
		}

		if oapersonal.PicSelfieKtp != nil && *oapersonal.PicSelfieKtp != "" {
			path := dir + *oapersonal.PicSelfieKtp
			responseData.PicSelfieKtp = &path
		}

		if oapersonal.RecImage1 != nil && *oapersonal.RecImage1 != "" {
			path := dir + "/signature/" + *oapersonal.RecImage1
			responseData.Signature = &path
		}

		responseData.OccupCompany = oapersonal.OccupCompany
		responseData.OccupPhone = oapersonal.OccupPhone
		responseData.OccupWebURL = oapersonal.OccupWebUrl
		responseData.MotherMaidenName = oapersonal.MotherMaidenName
		responseData.BeneficialFullName = oapersonal.BeneficialFullName

		//mapping gen lookup
		var personalDataLookupIds []string
		if oapersonal.Gender != nil {
			if _, ok := lib.Find(personalDataLookupIds, strconv.FormatUint(*oapersonal.Gender, 10)); !ok {
				personalDataLookupIds = append(personalDataLookupIds, strconv.FormatUint(*oapersonal.Gender, 10))
			}
		}
		if _, ok := lib.Find(personalDataLookupIds, strconv.FormatUint(oapersonal.IDcardType, 10)); !ok {
			personalDataLookupIds = append(personalDataLookupIds, strconv.FormatUint(oapersonal.IDcardType, 10))
		}
		if oapersonal.PepStatus != nil {
			if _, ok := lib.Find(personalDataLookupIds, strconv.FormatUint(*oapersonal.PepStatus, 10)); !ok {
				personalDataLookupIds = append(personalDataLookupIds, strconv.FormatUint(*oapersonal.PepStatus, 10))
			}
		}
		if oapersonal.MaritalStatus != nil {
			if _, ok := lib.Find(personalDataLookupIds, strconv.FormatUint(*oapersonal.MaritalStatus, 10)); !ok {
				personalDataLookupIds = append(personalDataLookupIds, strconv.FormatUint(*oapersonal.MaritalStatus, 10))
			}
		}
		if oapersonal.Religion != nil {
			if _, ok := lib.Find(personalDataLookupIds, strconv.FormatUint(*oapersonal.Religion, 10)); !ok {
				personalDataLookupIds = append(personalDataLookupIds, strconv.FormatUint(*oapersonal.Religion, 10))
			}
		}
		if oapersonal.Education != nil {
			if _, ok := lib.Find(personalDataLookupIds, strconv.FormatUint(*oapersonal.Education, 10)); !ok {
				personalDataLookupIds = append(personalDataLookupIds, strconv.FormatUint(*oapersonal.Education, 10))
			}
		}
		if oapersonal.OccupJob != nil {
			if _, ok := lib.Find(personalDataLookupIds, strconv.FormatUint(*oapersonal.OccupJob, 10)); !ok {
				personalDataLookupIds = append(personalDataLookupIds, strconv.FormatUint(*oapersonal.OccupJob, 10))
			}
		}
		if oapersonal.OccupPosition != nil {
			if _, ok := lib.Find(personalDataLookupIds, strconv.FormatUint(*oapersonal.OccupPosition, 10)); !ok {
				personalDataLookupIds = append(personalDataLookupIds, strconv.FormatUint(*oapersonal.OccupPosition, 10))
			}
		}
		if oapersonal.AnnualIncome != nil {
			if _, ok := lib.Find(personalDataLookupIds, strconv.FormatUint(*oapersonal.AnnualIncome, 10)); !ok {
				personalDataLookupIds = append(personalDataLookupIds, strconv.FormatUint(*oapersonal.AnnualIncome, 10))
			}
		}
		if oapersonal.SourceofFund != nil {
			if _, ok := lib.Find(personalDataLookupIds, strconv.FormatUint(*oapersonal.SourceofFund, 10)); !ok {
				personalDataLookupIds = append(personalDataLookupIds, strconv.FormatUint(*oapersonal.SourceofFund, 10))
			}
		}
		if oapersonal.InvesmentObjectives != nil {
			if _, ok := lib.Find(personalDataLookupIds, strconv.FormatUint(*oapersonal.InvesmentObjectives, 10)); !ok {
				personalDataLookupIds = append(personalDataLookupIds, strconv.FormatUint(*oapersonal.InvesmentObjectives, 10))
			}
		}
		if oapersonal.Correspondence != nil {
			if _, ok := lib.Find(personalDataLookupIds, strconv.FormatUint(*oapersonal.Correspondence, 10)); !ok {
				personalDataLookupIds = append(personalDataLookupIds, strconv.FormatUint(*oapersonal.Correspondence, 10))
			}
		}
		if oapersonal.BeneficialRelation != nil {
			if _, ok := lib.Find(personalDataLookupIds, strconv.FormatUint(*oapersonal.BeneficialRelation, 10)); !ok {
				personalDataLookupIds = append(personalDataLookupIds, strconv.FormatUint(*oapersonal.BeneficialRelation, 10))
			}
		}
		if oapersonal.EmergencyRelation != nil {
			if _, ok := lib.Find(personalDataLookupIds, strconv.FormatUint(*oapersonal.EmergencyRelation, 10)); !ok {
				personalDataLookupIds = append(personalDataLookupIds, strconv.FormatUint(*oapersonal.EmergencyRelation, 10))
			}
		}
		if oapersonal.RelationType != nil {
			if _, ok := lib.Find(personalDataLookupIds, strconv.FormatUint(*oapersonal.RelationType, 10)); !ok {
				personalDataLookupIds = append(personalDataLookupIds, strconv.FormatUint(*oapersonal.RelationType, 10))
			}
		}
		if oapersonal.RelationOccupation != nil {
			if _, ok := lib.Find(personalDataLookupIds, strconv.FormatUint(*oapersonal.RelationOccupation, 10)); !ok {
				personalDataLookupIds = append(personalDataLookupIds, strconv.FormatUint(*oapersonal.RelationOccupation, 10))
			}
		}
		if oapersonal.RelationBusinessFields != nil {
			if _, ok := lib.Find(personalDataLookupIds, strconv.FormatUint(*oapersonal.RelationBusinessFields, 10)); !ok {
				personalDataLookupIds = append(personalDataLookupIds, strconv.FormatUint(*oapersonal.RelationBusinessFields, 10))
			}
		}
		if oapersonal.OccupBusinessFields != nil {
			if _, ok := lib.Find(personalDataLookupIds, strconv.FormatUint(*oapersonal.OccupBusinessFields, 10)); !ok {
				personalDataLookupIds = append(personalDataLookupIds, strconv.FormatUint(*oapersonal.OccupBusinessFields, 10))
			}
		}
		//gen lookup personal data
		var lookupPersonData []models.GenLookup
		if len(personalDataLookupIds) > 0 {
			status, err = models.GetGenLookupIn(&lookupPersonData, personalDataLookupIds, "lookup_key")
			if err != nil {
				if err != sql.ErrNoRows {
					log.Error(err.Error())
					return lib.CustomError(status, err.Error(), "Failed get data")
				}
			}
		}

		pData := make(map[uint64]models.GenLookup)
		for _, genLook := range lookupPersonData {
			pData[genLook.LookupKey] = genLook
		}

		if oapersonal.Gender != nil {
			if n, ok := pData[*oapersonal.Gender]; ok {
				responseData.Gender = n.LkpName
			}
		}

		if oapersonal.PepStatus != nil {
			if n, ok := pData[*oapersonal.PepStatus]; ok {
				responseData.PepStatus = n.LkpName
			}
		}
		if oapersonal.MaritalStatus != nil {
			if n, ok := pData[*oapersonal.MaritalStatus]; ok {
				responseData.MaritalStatus = n.LkpName
			}
		}
		if oapersonal.Religion != nil {
			if n, ok := pData[*oapersonal.Religion]; ok {
				responseData.Religion = n.LkpName
			}
		}

		var country models.MsCountry

		strCountry := strconv.FormatUint(oapersonal.Nationality, 10)
		status, err = models.GetMsCountry(&country, strCountry)
		if err != nil {
			if err != sql.ErrNoRows {
				log.Error("Error Personal Data not Found")
				return lib.CustomError(status, err.Error(), "Personal data not found")
			}
		} else {
			responseData.Nationality = &country.CouName
		}
		if n, ok := pData[oapersonal.IDcardType]; ok {
			responseData.IDCardType = n.LkpName
		}

		if oapersonal.Education != nil {
			if n, ok := pData[*oapersonal.Education]; ok {
				responseData.Education = n.LkpName
			}
		}
		if oapersonal.OccupJob != nil {
			if n, ok := pData[*oapersonal.OccupJob]; ok {
				responseData.OccupJob = n.LkpName
			}
		}
		if oapersonal.OccupPosition != nil {
			if n, ok := pData[*oapersonal.OccupPosition]; ok {
				responseData.OccupPosition = n.LkpName
			}
		}
		if oapersonal.AnnualIncome != nil {
			if n, ok := pData[*oapersonal.AnnualIncome]; ok {
				responseData.AnnualIncome = n.LkpName
			}
		}
		if oapersonal.SourceofFund != nil {
			if n, ok := pData[*oapersonal.SourceofFund]; ok {
				responseData.SourceofFund = n.LkpName
			}
		}
		if oapersonal.InvesmentObjectives != nil {
			if n, ok := pData[*oapersonal.InvesmentObjectives]; ok {
				responseData.InvesmentObjectives = n.LkpName
			}
		}
		if oapersonal.Correspondence != nil {
			if n, ok := pData[*oapersonal.Correspondence]; ok {
				responseData.Correspondence = n.LkpName
			}
		}
		if oapersonal.BeneficialRelation != nil {
			if n, ok := pData[*oapersonal.BeneficialRelation]; ok {
				responseData.BeneficialRelation = n.LkpName
			}
		}
		if oapersonal.OccupBusinessFields != nil {
			if n, ok := pData[*oapersonal.OccupBusinessFields]; ok {
				responseData.OccupBusinessFields = n.LkpName
			}
		}

		//mapping idcard address &  domicile
		var postalAddressIds []string
		if oapersonal.IDcardAddressKey != nil {
			if _, ok := lib.Find(postalAddressIds, strconv.FormatUint(*oapersonal.IDcardAddressKey, 10)); !ok {
				postalAddressIds = append(postalAddressIds, strconv.FormatUint(*oapersonal.IDcardAddressKey, 10))
			}
		}
		if oapersonal.DomicileAddressKey != nil {
			if _, ok := lib.Find(postalAddressIds, strconv.FormatUint(*oapersonal.DomicileAddressKey, 10)); !ok {
				postalAddressIds = append(postalAddressIds, strconv.FormatUint(*oapersonal.DomicileAddressKey, 10))
			}
		}
		if oapersonal.OccupAddressKey != nil {
			if _, ok := lib.Find(postalAddressIds, strconv.FormatUint(*oapersonal.OccupAddressKey, 10)); !ok {
				postalAddressIds = append(postalAddressIds, strconv.FormatUint(*oapersonal.OccupAddressKey, 10))
			}
		}
		var oaPstalAddressList []models.OaPostalAddress
		if len(postalAddressIds) > 0 {
			status, err = models.GetOaPostalAddressIn(&oaPstalAddressList, postalAddressIds, "postal_address_key")
			if err != nil {
				if err != sql.ErrNoRows {
					log.Error(err.Error())
					return lib.CustomError(status, err.Error(), "Failed get data")
				}
			}
		}

		postalData := make(map[uint64]models.OaPostalAddress)
		for _, posAdd := range oaPstalAddressList {
			postalData[posAdd.PostalAddressKey] = posAdd
		}

		if len(postalData) > 0 {
			if oapersonal.IDcardAddressKey != nil {
				if p, ok := postalData[*oapersonal.IDcardAddressKey]; ok {
					responseData.IDcardAddress.Address = p.AddressLine1
					responseData.IDcardAddress.PostalCode = p.PostalCode

					var cityIds []string
					if p.KabupatenKey != nil {
						if _, ok := lib.Find(cityIds, strconv.FormatUint(*p.KabupatenKey, 10)); !ok {
							cityIds = append(cityIds, strconv.FormatUint(*p.KabupatenKey, 10))
						}
					}
					if p.KecamatanKey != nil {
						if _, ok := lib.Find(cityIds, strconv.FormatUint(*p.KecamatanKey, 10)); !ok {
							cityIds = append(cityIds, strconv.FormatUint(*p.KecamatanKey, 10))
						}
					}

					var cityList []models.MsCity
					if len(cityIds) > 0 {
						status, err = models.GetMsCityIn(&cityList, cityIds, "city_key")
						if err != nil {
							if err != sql.ErrNoRows {
								log.Error(err.Error())
								return lib.CustomError(status, err.Error(), "Failed get data")
							}
						}
					}
					cityData := make(map[uint64]models.MsCity)
					for _, city := range cityList {
						cityData[city.CityKey] = city
					}
					if c, ok := cityData[*p.KabupatenKey]; ok {
						var prov models.MsCity
						status, err = models.GetMsCity(&prov, strconv.FormatUint(*c.ParentKey, 10))
						if err == nil {
							responseData.IDcardAddress.Provinsi = &prov.CityName
						}
						responseData.IDcardAddress.Kabupaten = &c.CityName
					}
					if c, ok := cityData[*p.KecamatanKey]; ok {
						responseData.IDcardAddress.Kecamatan = &c.CityName
					}
				}
			}
			if oapersonal.DomicileAddressKey != nil {
				if p, ok := postalData[*oapersonal.DomicileAddressKey]; ok {
					responseData.DomicileAddress.Address = p.AddressLine1
					responseData.DomicileAddress.PostalCode = p.PostalCode

					var cityIds []string
					if p.KabupatenKey != nil {
						if _, ok := lib.Find(cityIds, strconv.FormatUint(*p.KabupatenKey, 10)); !ok {
							cityIds = append(cityIds, strconv.FormatUint(*p.KabupatenKey, 10))
						}
					}
					if p.KecamatanKey != nil {
						if _, ok := lib.Find(cityIds, strconv.FormatUint(*p.KecamatanKey, 10)); !ok {
							cityIds = append(cityIds, strconv.FormatUint(*p.KecamatanKey, 10))
						}
					}

					var cityList []models.MsCity
					if len(cityIds) > 0 {
						status, err = models.GetMsCityIn(&cityList, cityIds, "city_key")
						if err != nil {
							if err != sql.ErrNoRows {
								log.Error(err.Error())
								return lib.CustomError(status, err.Error(), "Failed get data")
							}
						}
					}
					cityData := make(map[uint64]models.MsCity)
					for _, city := range cityList {
						cityData[city.CityKey] = city
					}
					if p.KabupatenKey != nil {
						if c, ok := cityData[*p.KabupatenKey]; ok {
							var prov models.MsCity
							status, err = models.GetMsCity(&prov, strconv.FormatUint(*c.ParentKey, 10))
							if err == nil {
								responseData.DomicileAddress.Provinsi = &prov.CityName
							}
							responseData.DomicileAddress.Kabupaten = &c.CityName
						}
					}
					if p.KecamatanKey != nil {
						if c, ok := cityData[*p.KecamatanKey]; ok {
							responseData.DomicileAddress.Kecamatan = &c.CityName
						}
					}
				}
			}
			if oapersonal.OccupAddressKey != nil {
				if p, ok := postalData[*oapersonal.OccupAddressKey]; ok {
					responseData.OccupAddressKey.Address = p.AddressLine1
					responseData.OccupAddressKey.PostalCode = p.PostalCode

					var cityIds []string
					if p.KabupatenKey != nil {
						if _, ok := lib.Find(cityIds, strconv.FormatUint(*p.KabupatenKey, 10)); !ok {
							cityIds = append(cityIds, strconv.FormatUint(*p.KabupatenKey, 10))
						}
					}
					if p.KecamatanKey != nil {
						if _, ok := lib.Find(cityIds, strconv.FormatUint(*p.KecamatanKey, 10)); !ok {
							cityIds = append(cityIds, strconv.FormatUint(*p.KecamatanKey, 10))
						}
					}

					var cityList []models.MsCity
					if len(cityIds) > 0 {
						status, err = models.GetMsCityIn(&cityList, cityIds, "city_key")
						if err != nil {
							if err != sql.ErrNoRows {
								log.Error(err.Error())
								return lib.CustomError(status, err.Error(), "Failed get data")
							}
						}
					}
					cityData := make(map[uint64]models.MsCity)
					for _, city := range cityList {
						cityData[city.CityKey] = city
					}
					if p.KabupatenKey != nil {
						if c, ok := cityData[*p.KabupatenKey]; ok {
							responseData.DomicileAddress.Kabupaten = &c.CityName
						}
					}
					if p.KecamatanKey != nil {
						if c, ok := cityData[*p.KecamatanKey]; ok {
							responseData.DomicileAddress.Kecamatan = &c.CityName
						}
					}
				}
			}
		}

		//mapping relation
		if oapersonal.RelationType != nil {
			if n, ok := pData[*oapersonal.RelationType]; ok {
				responseData.Relation.RelationType = n.LkpName
			}
		}
		responseData.Relation.RelationFullName = oapersonal.RelationFullName
		if oapersonal.RelationOccupation != nil {
			if n, ok := pData[*oapersonal.RelationOccupation]; ok {
				responseData.Relation.RelationOccupation = n.LkpName
			}
		}
		if oapersonal.RelationBusinessFields != nil {
			if n, ok := pData[*oapersonal.RelationBusinessFields]; ok {
				responseData.Relation.RelationBusinessFields = n.LkpName
			}
		}

		//mapping emergency
		responseData.Emergency.EmergencyFullName = oapersonal.EmergencyFullName
		if oapersonal.EmergencyRelation != nil {
			if n, ok := pData[*oapersonal.EmergencyRelation]; ok {
				responseData.Emergency.EmergencyRelation = n.LkpName
			}
		}
		responseData.Emergency.EmergencyPhoneNo = oapersonal.EmergencyPhoneNo

		var oaRiskProfile []models.AdminOaRiskProfile
		status, err = models.AdminGetOaRiskProfile(&oaRiskProfile, strKey)
		if err != nil {
			if err != sql.ErrNoRows {
				log.Error(err.Error())
				return lib.CustomError(status, err.Error(), "Failed get data")
			}
		}
		responseData.RiskProfile = oaRiskProfile

		//mapping oa risk profile quiz
		var oaRiskProfileQuiz []models.AdminOaRiskProfileQuiz
		status, err = models.AdminGetOaRiskProfileQuizByOaRequestKey(&oaRiskProfileQuiz, strKey)
		if err != nil {
			if err != sql.ErrNoRows {
				log.Error(err.Error())
				return lib.CustomError(status, err.Error(), "Failed get data")
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
					return lib.CustomError(status, err.Error(), "Failed get data")
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
			responseData.RiskProfileQuiz = riskQuiz
		}

		//add response field Sinvest
		if oareq.CustomerKey != nil {
			var customer models.MsCustomer
			strCustomerKey := strconv.FormatUint(*oareq.CustomerKey, 10)
			status, err = models.GetMsCustomer(&customer, strCustomerKey)
			if err != nil {
				if err != sql.ErrNoRows {
					log.Error(err.Error())
					return lib.CustomError(status, err.Error(), "Failed get data")
				}
			}

			responseData.FirstName = customer.FirstName
			responseData.MiddleName = customer.MiddleName
			responseData.LastName = customer.LastName
			responseData.ClientCode = customer.ClientCode
			responseData.TinNumber = customer.TinNumber

			if customer.TinIssuanceDate != nil {
				layout := "2006-01-02 15:04:05"
				newLayout := "02 Jan 2006"
				date, _ := time.Parse(layout, *customer.TinIssuanceDate)
				oke := date.Format(newLayout)
				responseData.TinIssuanceDate = &oke
			}

			if customer.FatcaStatus != nil {
				var fatca models.GenLookup
				strLookKey := strconv.FormatUint(*customer.FatcaStatus, 10)
				status, err = models.GetGenLookup(&fatca, strLookKey)
				if err != nil {
					if err != sql.ErrNoRows {
						log.Error(err.Error())
						return lib.CustomError(status, err.Error(), "Failed get data")
					}
				}
				responseData.FatcaStatus = fatca.LkpName
			}

			if customer.TinIssuanceCountry != nil {
				var country models.MsCountry
				strCountryKey := strconv.FormatUint(*customer.TinIssuanceCountry, 10)
				status, err = models.GetMsCountry(&country, strCountryKey)
				if err != nil {
					if err != sql.ErrNoRows {
						log.Error(err.Error())
						return lib.CustomError(status, err.Error(), "Failed get data")
					}
				}
				responseData.TinIssuanceCountry = &country.CouName
			}
		} else {
			sliceName := strings.Fields(oapersonal.FullName)
			if len(sliceName) > 0 {
				if len(sliceName) == 1 {
					responseData.FirstName = &sliceName[0]
					responseData.LastName = &sliceName[0]
				}
				if len(sliceName) == 2 {
					responseData.FirstName = &sliceName[0]
					responseData.LastName = &sliceName[1]
				}
				if len(sliceName) > 2 {
					ln := len(sliceName)
					responseData.FirstName = &sliceName[0]
					responseData.LastName = &sliceName[1]
					lastName := strings.Join(sliceName[2:ln], " ")
					responseData.LastName = &lastName
				}
			}
		}
	}

	//mapping user approval
	var userApprovalIds []string
	if oareq.Check1References != nil {

		usercs, _ := strconv.ParseUint(*oareq.Check1References, 10, 64)
		if usercs > 0 {
			if _, ok := lib.Find(userApprovalIds, strconv.FormatUint(usercs, 10)); !ok {
				userApprovalIds = append(userApprovalIds, strconv.FormatUint(usercs, 10))
			}
		}
	}
	if oareq.Check2References != nil {
		userkyc, _ := strconv.ParseUint(*oareq.Check2References, 10, 64)
		if userkyc > 0 {
			if _, ok := lib.Find(userApprovalIds, strconv.FormatUint(userkyc, 10)); !ok {
				userApprovalIds = append(userApprovalIds, strconv.FormatUint(userkyc, 10))
			}
		}
	}

	//gen lookup personal data
	var userappr []models.ScUserLogin
	if len(userApprovalIds) > 0 {
		status, err = models.GetScUserLoginIn(&userappr, userApprovalIds, "user_login_key")
		if err != nil {
			if err != sql.ErrNoRows {
				log.Error(err.Error())
				return lib.CustomError(status, err.Error(), "Failed get data")
			}
		}
	}

	usrData := make(map[uint64]models.ScUserLogin)
	for _, usr := range userappr {
		usrData[usr.UserLoginKey] = usr
	}

	//set approv cs
	if oareq.Check1References != nil {
		usercs, _ := strconv.ParseUint(*oareq.Check1References, 10, 64)
		if usercs > 0 {
			if n, ok := usrData[usercs]; ok {
				var approvecs models.ApprovalData
				approvecs.ApproveUser = n.UloginFullName
				if oareq.Check1Date != nil {
					date, _ = time.Parse(layout, *oareq.Check1Date)
					oke := date.Format(newLayout)
					approvecs.ApproveDate = &oke
				}
				approvecs.ApproveNotes = oareq.Check1Notes
				if oareq.Check1Flag != nil {
					if *oareq.Check1Flag == 1 {
						approvecs.ApproveStatus = "Approved"
					} else {
						approvecs.ApproveStatus = "Rejected"
					}
				} else {
					approvecs.ApproveStatus = "-"
				}

				responseData.ApproveCS = &approvecs
			}
		}
	}

	//set approv kyc
	if oareq.Check2References != nil {
		userkyc, _ := strconv.ParseUint(*oareq.Check2References, 10, 64)
		if userkyc > 0 {
			if n, ok := usrData[userkyc]; ok {
				var approvekyc models.ApprovalData
				approvekyc.ApproveUser = n.UloginFullName
				if oareq.Check1Date != nil {
					date, _ = time.Parse(layout, *oareq.Check2Date)
					oke := date.Format(newLayout)
					approvekyc.ApproveDate = &oke
				}
				approvekyc.ApproveNotes = oareq.Check2Notes
				if oareq.Check2Flag != nil {
					if *oareq.Check2Flag == 1 {
						approvekyc.ApproveStatus = "Approved"
					} else {
						approvekyc.ApproveStatus = "Rejected"
					}
				} else {
					approvekyc.ApproveStatus = "-"
				}

				responseData.ApproveKYC = &approvekyc
			}
		}
	}

	//set bank_request
	var accBank []models.OaRequestByField
	status, err = models.GetOaRequestBankByField(&accBank, "oa_request_key", strconv.FormatUint(oareq.OaRequestKey, 10))
	if err != nil {
		responseData.BankRequest = nil
	} else {
		responseData.BankRequest = &accBank
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = responseData

	return c.JSON(http.StatusOK, response)
}

func AdminCreateCustomerIndividu(c echo.Context) error {
	var err error
	var status int
	var bindVarz [][]string

	paramsOaPersonalData := make(map[string]string)

	branchkey := c.FormValue("branch_key")
	if branchkey == "" {
		log.Error("Missing required parameter: branch_key")
		return lib.CustomError(http.StatusBadRequest, "branch_key can not be blank", "branch_key can not be blank")
	} else {
		n, err := strconv.ParseUint(branchkey, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: branch_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: branch_key", "Wrong input for parameter: branch_key")
		}
		if len(branchkey) > 11 {
			log.Error("Wrong input for parameter: branch_key too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: branch_key too long, max 11 character", "Missing required parameter: branch_key too long, max 11 character")
		}
	}

	agentkey := c.FormValue("agent_key")
	if agentkey == "" {
		log.Error("Missing required parameter: agent_key")
		return lib.CustomError(http.StatusBadRequest, "agent_key can not be blank", "agent_key can not be blank")
	} else {
		n, err := strconv.ParseUint(agentkey, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: agent_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: agent_key", "Wrong input for parameter: agent_key")
		}
		if len(agentkey) > 11 {
			log.Error("Wrong input for parameter: agent_key too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: agent_key too long, max 11 character", "Missing required parameter: agent_key too long, max 11 character")
		}
	}

	//SC_USER_LOGIN
	// Check parameters
	email := c.FormValue("email")
	if email == "" {
		log.Error("Missing required parameter: email")
		return lib.CustomError(http.StatusBadRequest, "email can not be blank", "email can not be blank")
	} else {
		if len(email) > 50 {
			log.Error("Wrong input for parameter: email too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: email too long, max 50 character", "Missing required parameter: email too long, max 50 character")
		}
	}
	uloginenabled := c.FormValue("ulogin_enabled")
	if uloginenabled == "" {
		log.Error("Missing required parameter: ulogin_enabled")
		return lib.CustomError(http.StatusBadRequest, "ulogin_enabled can not be blank", "ulogin_enabled can not be blank")
	} else {
		n, err := strconv.ParseUint(uloginenabled, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: ulogin_enabled")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: ulogin_enabled", "Wrong input for parameter: ulogin_enabled")
		}
		if len(agentkey) > 11 {
			log.Error("Wrong input for parameter: ulogin_enabled too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: ulogin_enabled too long, max 11 character", "Missing required parameter: ulogin_enabled too long, max 11 character")
		}
	}

	phone := c.FormValue("phone")
	if phone == "" {
		log.Error("Missing required parameter: phone")
		return lib.CustomError(http.StatusBadRequest, "phone can not be blank", "phone can not be blank")
	} else {
		if len(phone) > 20 {
			log.Error("Wrong input for parameter: phone too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: phone too long, max 20 character", "Missing required parameter: phone too long, max 20 character")
		}
	}
	password := c.FormValue("password")
	if password == "" {
		log.Error("Missing required parameter: password")
		return lib.CustomError(http.StatusBadRequest, "password can not be blank", "password can not be blank")
	} else {
		if len(password) > 50 {
			log.Error("Wrong input for parameter: password too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: password too long, max 50 character", "Missing required parameter: password too long, max 50 character")
		}
	}
	confpassword := c.FormValue("conf_password")
	if confpassword == "" {
		log.Error("Missing required parameter: conf_password")
		return lib.CustomError(http.StatusBadRequest, "conf_password can not be blank", "conf_password can not be blank")
	} else {
		if len(confpassword) > 50 {
			log.Error("Wrong input for parameter: conf_password too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: conf_password too long, max 50 character", "Missing required parameter: conf_password too long, max 50 character")
		}
	}
	// Validate email
	err = checkmail.ValidateFormat(email)
	if err != nil {
		log.Error("Email format is not valid")
		return lib.CustomError(http.StatusBadRequest, "Email format is not valid", "Email format is not valid")
	}
	var user []models.ScUserLogin
	paramsCekUserLogin := make(map[string]string)
	paramsCekUserLogin["ulogin_email"] = email
	status, err = models.GetAllScUserLogin(&user, 0, 0, paramsCekUserLogin, true)
	if err != nil {
		log.Error("Error get email " + email)
		return lib.CustomError(status, err.Error(), "Error get email")
	}
	if len(user) > 0 {
		log.Error("Email " + email + " already registered")
		return lib.CustomError(http.StatusBadRequest, "Email "+email+" already registered", "Email kamu sudah terdaftar.\nSilakan masukkan email lainnya atau hubungi Customer.")
	}

	// Validate password
	length, number, upper, special := verifyPassword(password)
	if length == false || number == false || upper == false || special == false {
		log.Error("Password does meet the criteria")
		return lib.CustomError(http.StatusBadRequest, "Password does meet the criteria", "Your password need at least 8 character length, has lower and upper case letter, has numeric letter, and has special character")
	}

	if strings.Contains(password, confpassword) == false {
		log.Error("Missing required parameter: conf_password must equal with password")
		return lib.CustomError(http.StatusBadRequest, "conf_password must equal with password", "conf_password must equal with password")
	}

	//INFORMASI NASABAH
	fullname := c.FormValue("full_name")
	if fullname == "" {
		log.Error("Missing required parameter: full_name")
		return lib.CustomError(http.StatusBadRequest, "full_name can not be blank", "full_name can not be blank")
	} else {
		if len(fullname) > 50 {
			log.Error("Wrong input for parameter: full_name too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: full_name too long, max 50 character", "Missing required parameter: full_name too long, max 50 character")
		}
	}

	nationality := c.FormValue("nationality")
	if nationality == "" {
		log.Error("Missing required parameter: nationality")
		return lib.CustomError(http.StatusBadRequest, "nationality can not be blank", "nationality can not be blank")
	} else {
		n, err := strconv.ParseUint(nationality, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: nationality")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: nationality", "Wrong input for parameter: nationality")
		}
		if len(nationality) > 11 {
			log.Error("Wrong input for parameter: nationality too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: nationality too long, max 11 character", "Missing required parameter: nationality too long, max 11 character")
		}
	}

	idcardNumber := c.FormValue("idcard_no")
	if idcardNumber == "" {
		log.Error("Missing required parameter: idcard_no")
		return lib.CustomError(http.StatusBadRequest, "idcard_no can not be blank", "idcard_no can not be blank")
	} else {
		if len(idcardNumber) > 20 {
			log.Error("Wrong input for parameter: idcard_no too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: idcard_no too long, max 20 character", "Missing required parameter: idcard_no too long, max 20 character")
		}
	}

	gender := c.FormValue("gender")
	if gender == "" {
		log.Error("Missing required parameter: gender")
		return lib.CustomError(http.StatusBadRequest, "gender can not be blank", "gender can not be blank")
	} else {
		n, err := strconv.ParseUint(gender, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: gender")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: gender", "Wrong input for parameter: gender")
		}
		if len(gender) > 11 {
			log.Error("Wrong input for parameter: gender too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: gender too long, max 11 character", "Missing required parameter: gender too long, max 11 character")
		}
	}

	salesCode := c.FormValue("sales_code")
	if salesCode == "" {
		salesCode = "1"
		// log.Error("Missing required parameter: sales_code")
		// return lib.CustomError(http.StatusBadRequest, "sales_code can not be blank", "sales_code can not be blank")
	} else {
		if len(salesCode) > 20 {
			log.Error("Wrong input for parameter: sales_code too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: sales_code too long, max 20 character", "Missing required parameter: sales_code too long, max 20 character")
		}
	}

	placeBirth := c.FormValue("place_birth")
	if placeBirth == "" {
		log.Error("Missing required parameter: place_birth")
		return lib.CustomError(http.StatusBadRequest, "place_birth can not be blank", "place_birth can not be blank")
	} else {
		if len(placeBirth) > 50 {
			log.Error("Wrong input for parameter: place_birth too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: place_birth too long, max 50 character", "Missing required parameter: place_birth too long, max 50 character")
		}
	}

	dateBirth := c.FormValue("date_birth")
	if dateBirth == "" {
		log.Error("Missing required parameter: date_birth")
		return lib.CustomError(http.StatusBadRequest, "date_birth can not be blank", "date_birth can not be blank")
	} else {
		if len(dateBirth) > 50 {
			log.Error("Wrong input for parameter: date_birth too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: date_birth too long, max 50 character", "Missing required parameter: date_birth too long, max 50 character")
		}
	}

	maritalStatus := c.FormValue("marital_status")
	if maritalStatus == "" {
		log.Error("Missing required parameter: marital_status")
		return lib.CustomError(http.StatusBadRequest, "marital_status can not be blank", "marital_status can not be blank")
	} else {
		n, err := strconv.ParseUint(maritalStatus, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: marital_status")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: marital_status", "Wrong input for parameter: marital_status")
		}
		if len(maritalStatus) > 11 {
			log.Error("Wrong input for parameter: marital_status too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: marital_status too long, max 11 character", "Missing required parameter: marital_status too long, max 11 character")
		}
	}

	addressid := c.FormValue("address_id")
	if addressid == "" {
		log.Error("Missing required parameter: address_id")
		return lib.CustomError(http.StatusBadRequest, "address_id can not be blank", "address_id can not be blank")
	} else {
		if len(addressid) > 200 {
			log.Error("Wrong input for parameter: address_id too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: address_id too long, max 200 character", "Missing required parameter: address_id too long, max 200 character")
		}
	}

	kabupatenid := c.FormValue("kabupaten_id")
	if kabupatenid != "" {
		n, err := strconv.ParseUint(kabupatenid, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: kabupaten_id")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: kabupaten_id", "Wrong input for parameter: kabupaten_id")
		}
		if len(kabupatenid) > 11 {
			log.Error("Wrong input for parameter: kabupaten_id too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: kabupaten_id too long, max 11 character", "Missing required parameter: kabupaten_id too long, max 11 character")
		}
	} else {
		log.Error("Missing required parameter: kabupaten_id")
		return lib.CustomError(http.StatusBadRequest, "kabupaten_id can not be blank", "kabupaten_id can not be blank")
	}

	kecamatanid := c.FormValue("kecamatan_id")
	if kecamatanid != "" {
		n, err := strconv.ParseUint(kecamatanid, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: kecamatan_id")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: kecamatan_id", "Wrong input for parameter: kecamatan_id")
		}
		if len(kecamatanid) > 11 {
			log.Error("Wrong input for parameter: kecamatan_id too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: kecamatan_id too long, max 11 character", "Missing required parameter: kecamatan_id too long, max 11 character")
		}
	} else {
		log.Error("Missing required parameter: kecamatan_id")
		return lib.CustomError(http.StatusBadRequest, "kecamatan_id can not be blank", "kecamatan_id can not be blank")
	}

	postalid := c.FormValue("postal_id")
	if postalid != "" {
		ps, err := strconv.ParseUint(postalid, 10, 64)
		if err != nil || ps == 0 {
			log.Error("Wrong input for parameter: postal_id")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: postal_id", "Wrong input for parameter: postal_id")
		}
		if len(postalid) > 10 {
			log.Error("Wrong input for parameter: postal_id too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: postal_id too long, max 10 character", "Missing required parameter: postal_id too long, max 10 character")
		}
	} else {
		log.Error("Missing required parameter: postal_id")
		return lib.CustomError(http.StatusBadRequest, "postal_id can not be blank", "postal_id can not be blank")
	}

	addressdomicile := c.FormValue("address_domicile")
	if addressdomicile == "" {
		log.Error("Missing required parameter: address_domicile")
		return lib.CustomError(http.StatusBadRequest, "address_domicile can not be blank", "address_domicile can not be blank")
	} else {
		if len(addressdomicile) > 200 {
			log.Error("Wrong input for parameter: address_domicile too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: address_domicile too long, max 200 character", "Missing required parameter: address_domicile too long, max 200 character")
		}
	}

	kabupatendomicile := c.FormValue("kabupaten_domicile")
	if kabupatendomicile != "" {
		n, err := strconv.ParseUint(kabupatendomicile, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: kabupaten_domicile")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: kabupaten_domicile", "Wrong input for parameter: kabupaten_domicile")
		}
		if len(kabupatendomicile) > 11 {
			log.Error("Wrong input for parameter: kabupaten_domicile too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: kabupaten_domicile too long, max 11 character", "Missing required parameter: kabupaten_domicile too long, max 11 character")
		}
	} else {
		log.Error("Missing required parameter: kabupaten_domicile")
		return lib.CustomError(http.StatusBadRequest, "kabupaten_domicile can not be blank", "kabupaten_domicile can not be blank")
	}

	kecamatandomicile := c.FormValue("kecamatan_domicile")
	if kecamatandomicile != "" {
		n, err := strconv.ParseUint(kecamatandomicile, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: kecamatan_domicile")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: kecamatan_domicile", "Wrong input for parameter: kecamatan_domicile")
		}
		if len(kecamatandomicile) > 11 {
			log.Error("Wrong input for parameter: kecamatan_domicile too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: kecamatan_domicile too long, max 11 character", "Missing required parameter: kecamatan_domicile too long, max 11 character")
		}
	} else {
		log.Error("Missing required parameter: kecamatan_domicile")
		return lib.CustomError(http.StatusBadRequest, "kecamatan_domicile can not be blank", "kecamatan_domicile can not be blank")
	}

	postaldomicile := c.FormValue("postal_domicile")
	if postaldomicile != "" {
		ps, err := strconv.ParseUint(postaldomicile, 10, 64)
		if err != nil || ps == 0 {
			log.Error("Wrong input for parameter: postal_domicile")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: postal_domicile", "Wrong input for parameter: postal_domicile")
		}
		if len(postaldomicile) > 10 {
			log.Error("Wrong input for parameter: postal_domicile too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: postal_domicile too long, max 10 character", "Missing required parameter: postal_domicile too long, max 10 character")
		}
	} else {
		log.Error("Missing required parameter: postal_domicile")
		return lib.CustomError(http.StatusBadRequest, "postal_domicile can not be blank", "postal_domicile can not be blank")
	}

	phoneHome := c.FormValue("phone_home")
	if phoneHome == "" {
		log.Error("Missing required parameter: phone_home")
		return lib.CustomError(http.StatusBadRequest, "phone_home can not be blank", "phone_home can not be blank")
	} else {
		if len(phoneHome) > 20 {
			log.Error("Wrong input for parameter: phone_home too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: phone_home too long, max 20 character", "Missing required parameter: phone_home too long, max 20 character")
		}
	}

	religion := c.FormValue("religion")
	if religion == "" {
		log.Error("Missing required parameter: religion")
		return lib.CustomError(http.StatusBadRequest, "religion can not be blank", "religion can not be blank")
	} else {
		ps, err := strconv.ParseUint(religion, 10, 64)
		if err != nil || ps == 0 {
			log.Error("Wrong input for parameter: religion")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: religion", "Wrong input for parameter: religion")
		}
		if len(religion) > 11 {
			log.Error("Wrong input for parameter: religion too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: religion too long, max 11 character", "Missing required parameter: religion too long, max 11 character")
		}
	}

	religionOther := c.FormValue("religion_other")
	if religionOther != "" {
		var rowz []string
		rowz = append(rowz, "1")
		rowz = append(rowz, "0")
		rowz = append(rowz, religionOther)
		bindVarz = append(bindVarz, rowz)
	}

	education := c.FormValue("education")
	if education == "" {
		log.Error("Missing required parameter: education")
		return lib.CustomError(http.StatusBadRequest, "education can not be blank", "education can not be blank")
	} else {
		ps, err := strconv.ParseUint(education, 10, 64)
		if err != nil || ps == 0 {
			log.Error("Wrong input for parameter: education")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: education", "Wrong input for parameter: education")
		}
		if len(education) > 11 {
			log.Error("Wrong input for parameter: education too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: education too long, max 11 character", "Missing required parameter: education too long, max 11 character")
		}
	}

	educationOther := c.FormValue("education_other")
	if educationOther != "" {
		var rowz []string
		rowz = append(rowz, "3")
		rowz = append(rowz, "0")
		rowz = append(rowz, educationOther)
		bindVarz = append(bindVarz, rowz)
	}

	//UPLOAD DOKUMEN FOTO E-KTP & SELFIE DENGAN KTP
	var file *multipart.FileHeader
	file, err = c.FormFile("pic_ktp")
	if file == nil {
		log.Error("Missing required parameter: pic_ktp")
		return lib.CustomError(http.StatusBadRequest, "pic_ktp can not be blank", "pic_ktp can not be blank")
	}

	var fileselfie *multipart.FileHeader
	fileselfie, err = c.FormFile("pic_selfie_ktp")
	if fileselfie == nil {
		log.Error("Missing required parameter: pic_selfie_ktp")
		return lib.CustomError(http.StatusBadRequest, "pic_selfie_ktp can not be blank", "pic_selfie_ktp can not be blank")
	}

	//URAIAN BIDANG USAHA DAN PEKERJAAN
	job := c.FormValue("job")
	if job == "" {
		log.Error("Missing required parameter: job")
		return lib.CustomError(http.StatusBadRequest, "job can not be blank", "job can not be blank")
	} else {
		if len(job) > 11 {
			log.Error("Wrong input for parameter: job too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: job too long, max 11 character", "Missing required parameter: job too long, max 11 character")
		}
		n, err := strconv.ParseUint(job, 10, 64)
		if err == nil && n > 0 {
			paramsOaPersonalData["occup_job"] = job
		} else {
			log.Error("Wrong input for parameter: job")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: job", "Wrong input for parameter: job")
		}
	}

	jobOther := c.FormValue("job_other")
	if jobOther != "" {
		var rowz []string
		rowz = append(rowz, "2")
		rowz = append(rowz, "0")
		rowz = append(rowz, jobOther)
		bindVarz = append(bindVarz, rowz)
	}

	company := c.FormValue("company")
	if company != "" {
		if len(company) > 50 {
			log.Error("Wrong input for parameter: company too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: company too long, max 50 character", "Missing required parameter: company too long, max 50 character")
		}
		paramsOaPersonalData["occup_company"] = company
	}

	position := c.FormValue("position")
	if position != "" {
		if len(position) > 11 {
			log.Error("Wrong input for parameter: position too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: position too long, max 11 character", "Missing required parameter: position too long, max 11 character")
		}
		n, err := strconv.ParseUint(position, 10, 64)
		if err == nil && n > 0 {
			paramsOaPersonalData["occup_position"] = position
		} else {
			log.Error("Wrong input for parameter: position")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: position", "Wrong input for parameter: position")
		}
	}

	positionOther := c.FormValue("position_other")
	if positionOther != "" {
		var rowz []string
		rowz = append(rowz, "10")
		rowz = append(rowz, "0")
		rowz = append(rowz, positionOther)
		bindVarz = append(bindVarz, rowz)
	}

	businessField := c.FormValue("business_field")
	if businessField != "" {
		if len(businessField) > 11 {
			log.Error("Wrong input for parameter: business_field too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: business_field too long, max 11 character", "Missing required parameter: business_field too long, max 11 character")
		}
		n, err := strconv.ParseUint(businessField, 10, 64)
		if err == nil && n > 0 {
			paramsOaPersonalData["occup_business_fields"] = businessField
		} else {
			log.Error("Wrong input for parameter: business_field")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: business_field", "Wrong input for parameter: business_field")
		}
	}

	businessFieldOther := c.FormValue("business_field_other")
	if businessFieldOther != "" {
		var rowz []string
		rowz = append(rowz, "9")
		rowz = append(rowz, "0")
		rowz = append(rowz, businessFieldOther)
		bindVarz = append(bindVarz, rowz)
	}

	annualIncome := c.FormValue("annual_income")
	if annualIncome == "" {
		log.Error("Missing required parameter: annual_income")
		return lib.CustomError(http.StatusBadRequest, "annual_income can not be blank", "annual_income can not be blank")
	} else {
		if len(annualIncome) > 11 {
			log.Error("Wrong input for parameter: annual_income too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: annual_income too long, max 11 character", "Missing required parameter: annual_income too long, max 11 character")
		}
		n, err := strconv.ParseUint(annualIncome, 10, 64)
		if err == nil && n > 0 {
			paramsOaPersonalData["annual_income"] = annualIncome
		} else {
			log.Error("Wrong input for parameter: annual_income")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: annual_income", "Wrong input for parameter: annual_income")
		}
	}

	fundSource := c.FormValue("fund_source")
	if fundSource == "" {
		log.Error("Missing required parameter: fund_source")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: fund_source", "Missing required parameter: fund_source")
	} else {
		if len(fundSource) > 11 {
			log.Error("Wrong input for parameter: fund_source too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: fund_source too long, max 11 character", "Missing required parameter: fund_source too long, max 11 character")
		}
		n, err := strconv.ParseUint(fundSource, 10, 64)
		if err == nil && n > 0 {
			paramsOaPersonalData["sourceof_fund"] = fundSource
		} else {
			log.Error("Wrong input for parameter: fund_source")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: fund_source", "Wrong input for parameter: fund_source")
		}
	}

	fundSourceOther := c.FormValue("fund_source_other")
	if fundSourceOther != "" {
		var row []string
		row = append(row, "5")
		row = append(row, "0")
		row = append(row, fundSourceOther)
		bindVarz = append(bindVarz, row)
	}

	pepStatus := c.FormValue("pep_status")
	if pepStatus == "" {
		log.Error("Missing required parameter: pep_status")
		return lib.CustomError(http.StatusBadRequest, "pep_status can not be blank", "pep_status can not be blank")
	} else {
		if len(pepStatus) > 11 {
			log.Error("Wrong input for parameter: pep_status too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: pep_status too long, max 11 character", "Missing required parameter: pep_status too long, max 11 character")
		}
		n, err := strconv.ParseUint(pepStatus, 10, 64)
		if err == nil && n > 0 {
			paramsOaPersonalData["pep_status"] = pepStatus
		} else {
			log.Error("Wrong input for parameter: pep_status")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: pep_status", "Wrong input for parameter: pep_status")
		}
	}

	objectives := c.FormValue("objectives")
	if objectives == "" {
		log.Error("Missing required parameter: objectives")
		return lib.CustomError(http.StatusBadRequest, "objectives can not be blank", "objectives can not be blank")
	} else {
		if len(objectives) > 11 {
			log.Error("Wrong input for parameter: objectives too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: objectives too long, max 11 character", "Missing required parameter: objectives too long, max 11 character")
		}
		n, err := strconv.ParseUint(objectives, 10, 64)
		if err == nil && n > 0 {
			paramsOaPersonalData["invesment_objectives"] = objectives
		} else {
			log.Error("Wrong input for parameter: objectives")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: objectives", "Wrong input for parameter: objectives")
		}
	}

	objectivesOther := c.FormValue("objectives_other")
	if objectivesOther != "" {
		var row []string
		row = append(row, "6")
		row = append(row, "0")
		row = append(row, objectivesOther)
		bindVarz = append(bindVarz, row)
	}

	corespondence := c.FormValue("corespondence")
	if corespondence != "" {
		if len(corespondence) > 11 {
			log.Error("Wrong input for parameter: corespondence too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: corespondence too long, max 11 character", "Missing required parameter: corespondence too long, max 11 character")
		}
		n, err := strconv.ParseUint(corespondence, 10, 64)
		if err == nil && n > 0 {
			paramsOaPersonalData["correspondence"] = corespondence
		} else {
			log.Error("Wrong input for parameter: corespondence")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: corespondence", "Wrong input for parameter: corespondence")
		}
	}

	//TAB 4
	motherMaidenName := c.FormValue("mother_maiden_name")
	if motherMaidenName == "" {
		log.Error("Missing required parameter: mother_maiden_name")
		return lib.CustomError(http.StatusBadRequest, "mother_maiden_name can not be blank", "mother_maiden_name can not be blank")
	} else {
		if len(motherMaidenName) > 50 {
			log.Error("Wrong input for parameter: mother_maiden_name too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: mother_maiden_name too long, max 50 character", "Missing required parameter: mother_maiden_name too long, max 50 character")
		}
	}

	relationOccupation := c.FormValue("relation_occupation")
	if relationOccupation == "" {
		log.Error("Missing required parameter: relation_occupation")
		return lib.CustomError(http.StatusBadRequest, "relation_occupation can not be blank", "relation_occupation can not be blank")
	} else {
		if len(relationOccupation) > 11 {
			log.Error("Wrong input for parameter: relation_occupation too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: relation_occupation too long, max 11 character", "Missing required parameter: relation_occupation too long, max 11 character")
		}
		n, err := strconv.ParseUint(relationOccupation, 10, 64)
		if err == nil && n > 0 {
			paramsOaPersonalData["relation_occupation"] = relationOccupation
		} else {
			log.Error("Wrong input for parameter: relation_occupation")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: relation_occupation", "Wrong input for parameter: relation_occupation")
		}
	}

	relationOccupationOther := c.FormValue("relation_occupation_other")
	if relationOccupationOther != "" {
		var row []string
		row = append(row, "8")
		row = append(row, "0")
		row = append(row, relationOccupationOther)
		bindVarz = append(bindVarz, row)
	}

	relationType := c.FormValue("relation_type")
	if relationType == "" {
		log.Error("Missing required parameter: relation_type")
		return lib.CustomError(http.StatusBadRequest, "relation_type can not be blank", "relation_type can not be blank")
	} else {
		if len(relationType) > 11 {
			log.Error("Wrong input for parameter: relation_type too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: relation_type too long, max 11 character", "Missing required parameter: relation_type too long, max 11 character")
		}
		n, err := strconv.ParseUint(relationType, 10, 64)
		if err == nil && n > 0 {
			paramsOaPersonalData["relation_type"] = relationType
		} else {
			log.Error("Wrong input for parameter: relation_type")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: relation_type", "Wrong input for parameter: relation_type")
		}
	}

	relationName := c.FormValue("relation_name")
	if relationName == "" {
		log.Error("Missing required parameter: relation_name")
		return lib.CustomError(http.StatusBadRequest, "relation_name can not be blank", "relation_name can not be blank")
	} else {
		if len(relationName) > 50 {
			log.Error("Wrong input for parameter: relation_name too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: relation_name too long, max 50 character", "Missing required parameter: relation_name too long, max 50 character")
		}
	}

	relationBusinessField := c.FormValue("relation_business_field")
	if relationBusinessField != "" {
		if len(relationBusinessField) > 11 {
			log.Error("Wrong input for parameter: relation_business_field too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: relation_business_field too long, max 11 character", "Missing required parameter: relation_business_field too long, max 11 character")
		}
		n, err := strconv.ParseUint(relationBusinessField, 10, 64)
		if err == nil && n > 0 {
			paramsOaPersonalData["relation_business_fields"] = relationBusinessField
		} else {
			log.Error("Wrong input for parameter: relation_business_field")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: relation_business_field", "Wrong input for parameter: relation_business_field")
		}
	}

	relationBusinessFieldOther := c.FormValue("relation_business_field_other")
	if relationBusinessFieldOther != "" {
		var row []string
		row = append(row, "9")
		row = append(row, "0")
		row = append(row, relationBusinessFieldOther)
		bindVarz = append(bindVarz, row)
	}

	emergencyName := c.FormValue("emergency_name")
	if emergencyName == "" {
		log.Error("Missing required parameter: emergency_name")
		return lib.CustomError(http.StatusBadRequest, "emergency_name can not be blank", "emergency_name can not be blank")
	} else {
		if len(emergencyName) > 50 {
			log.Error("Wrong input for parameter: emergency_name too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: emergency_name too long, max 50 character", "Missing required parameter: emergency_name too long, max 50 character")
		}
	}

	emergencyRelation := c.FormValue("emergency_relation")
	if emergencyRelation == "" {
		log.Error("Missing required parameter: emergency_relation")
		return lib.CustomError(http.StatusBadRequest, "emergency_relation can not be blank", "emergency_relation can not be blank")
	} else {
		if len(emergencyRelation) > 11 {
			log.Error("Wrong input for parameter: emergency_relation too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: emergency_relation too long, max 11 character", "Missing required parameter: emergency_relation too long, max 11 character")
		}
		n, err := strconv.ParseUint(emergencyRelation, 10, 64)
		if err == nil && n > 0 {
			paramsOaPersonalData["emergency_relation"] = emergencyRelation
		} else {
			log.Error("Wrong input for parameter: emergency_relation")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: emergency_relation", "Wrong input for parameter: emergency_relation")
		}
	}

	emergencyPhone := c.FormValue("emergency_phone")
	if emergencyPhone == "" {
		log.Error("Missing required parameter: emergency_phone")
		return lib.CustomError(http.StatusBadRequest, "emergency_phone can not be blank", "emergency_phone can not be blank")
	} else {
		if len(emergencyPhone) > 20 {
			log.Error("Wrong input for parameter: emergency_phone too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: emergency_phone too long, max 20 character", "Missing required parameter: emergency_phone too long, max 20 character")
		}
	}

	//TAB 5 REKENING DLL
	beneficialRelation := c.FormValue("beneficial_relation")
	if beneficialRelation == "" {
		log.Error("Missing required parameter: beneficial_relation")
		return lib.CustomError(http.StatusBadRequest, "beneficial_relation can not be blank", "beneficial_relation can not be blank")
	} else {
		if len(beneficialRelation) > 11 {
			log.Error("Wrong input for parameter: beneficial_relation too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: beneficial_relation too long, max 11 character", "Missing required parameter: beneficial_relation too long, max 11 character")
		}
		n, err := strconv.ParseUint(beneficialRelation, 10, 64)
		if err == nil && n > 0 {
			paramsOaPersonalData["beneficial_relation"] = beneficialRelation
		} else {
			log.Error("Wrong input for parameter: beneficial_relation")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: beneficial_relation", "Wrong input for parameter: beneficial_relation")
		}
	}

	beneficialName := c.FormValue("beneficial_name")
	if beneficialName == "" {
		log.Error("Missing required parameter: beneficial_name")
		return lib.CustomError(http.StatusBadRequest, "beneficial_name can not be blank", "beneficial_name can not be blank")
	} else {
		if len(beneficialName) > 50 {
			log.Error("Wrong input for parameter: beneficial_name too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: beneficial_name too long, max 50 character", "Missing required parameter: beneficial_name too long, max 50 character")
		}
	}

	bankKey := c.FormValue("bank_key")
	if bankKey == "" {
		log.Error("Missing required parameter: bank_key")
		return lib.CustomError(http.StatusBadRequest, "bank_key can not be blank", "bank_key can not be blank")
	} else {
		bank, err := strconv.ParseUint(bankKey, 10, 64)
		if err != nil || bank == 0 {
			log.Error("Wrong input for parameter: bank_key")
			return lib.CustomError(http.StatusBadRequest)
		}
		if len(bankKey) > 11 {
			log.Error("Wrong input for parameter: bank_key too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: bank_key too long, max 11 character", "Missing required parameter: bank_key too long, max 11 character")
		}
	}

	accountNo := c.FormValue("account_no")
	if accountNo == "" {
		log.Error("Missing required parameter: account_no")
		return lib.CustomError(http.StatusBadRequest, "account_no can not be blank", "account_no can not be blank")
	} else {
		if len(accountNo) > 50 {
			log.Error("Wrong input for parameter: account_no too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: account_no too long, max 50 character", "Missing required parameter: account_no too long, max 50 character")
		}
	}

	accountName := c.FormValue("account_name")
	if accountName == "" {
		log.Error("Missing required parameter: account_name")
		return lib.CustomError(http.StatusBadRequest, "account_name can not be blank", "account_name can not be blank")
	} else {
		if len(accountName) > 50 {
			log.Error("Wrong input for parameter: account_name too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: account_name too long, max 50 character", "Missing required parameter: account_name too long, max 50 character")
		}
	}

	branchName := c.FormValue("branch_name")
	if branchName == "" {
		log.Error("Missing required parameter: branch_name")
		return lib.CustomError(http.StatusBadRequest, "branch_name can not be blank", "branch_name can not be blank")
	} else {
		if len(branchName) > 50 {
			log.Error("Wrong input for parameter: branch_name too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: branch_name too long, max 50 character", "Missing required parameter: branch_name too long, max 50 character")
		}
	}

	quizOption := c.FormValue("quiz_option")
	if quizOption == "" {
		log.Error("Missing required parameter: quiz_option")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: quiz_option", "Missing required parameter: quiz_option")
	}

	s := strings.Split(quizOption, ",")
	var quizoptionkey []string

	for _, value := range s {
		is := strings.TrimSpace(value)
		if is != "" {
			if _, ok := lib.Find(quizoptionkey, is); !ok {
				quizoptionkey = append(quizoptionkey, is)
			}
		}
	}
	if len(quizoptionkey) <= 0 {
		log.Error("Missing required parameter: quiz_option")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: quiz_option", "Missing required parameter: quiz_option")
	}

	// Encrypt password
	encryptedPasswordByte := sha256.Sum256([]byte(password))
	encryptedPassword := hex.EncodeToString(encryptedPasswordByte[:])
	// Set expired for token
	date := time.Now().AddDate(0, 0, 1)
	dateLayout := "2006-01-02 15:04:05"
	expired := date.Format(dateLayout)
	// Generate verify key
	verifyKeyByte := sha256.Sum256([]byte(email + "_" + expired))
	verifyKey := hex.EncodeToString(verifyKeyByte[:])

	paramsUserLogin := make(map[string]string)

	//SC_USER_LOGIN
	paramsUserLogin["ulogin_name"] = email
	paramsUserLogin["ulogin_email"] = email
	paramsUserLogin["ulogin_full_name"] = fullname
	paramsUserLogin["ulogin_password"] = encryptedPassword
	paramsUserLogin["ulogin_must_changepwd"] = "1"
	paramsUserLogin["user_dept_key"] = "1"
	paramsUserLogin["last_password_changed"] = time.Now().Format(dateLayout)
	paramsUserLogin["verified_email"] = "1"
	paramsUserLogin["verified_mobileno"] = "1"
	paramsUserLogin["ulogin_mobileno"] = phone
	paramsUserLogin["ulogin_enabled"] = uloginenabled
	paramsUserLogin["ulogin_locked"] = "0"
	paramsUserLogin["ulogin_failed_count"] = "0"
	paramsUserLogin["user_category_key"] = "1"
	paramsUserLogin["last_access"] = time.Now().Format(dateLayout)
	paramsUserLogin["rec_created_date"] = time.Now().Format(dateLayout)
	paramsUserLogin["rec_created_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
	paramsUserLogin["accept_login_tnc"] = "1"
	paramsUserLogin["allowed_sharing_login"] = "1"
	paramsUserLogin["string_token"] = verifyKey
	paramsUserLogin["token_expired"] = expired
	paramsUserLogin["rec_status"] = "1"

	//OA_REQUEST
	layout := "2006-01-02 15:04:05"
	dateNow := time.Now().Format(layout)
	paramsOaRequest := make(map[string]string)
	paramsOaRequest["oa_status"] = "258"
	paramsOaRequest["oa_entry_start"] = dateNow
	if salesCode != "" {
		paramsOaRequest["sales_code"] = salesCode
	}
	paramsOaRequest["oa_entry_end"] = dateNow
	paramsOaRequest["branch_key"] = branchkey
	paramsOaRequest["agent_key"] = agentkey
	paramsOaRequest["oa_request_type"] = "127"
	paramsOaRequest["rec_status"] = "1"
	paramsOaRequest["rec_created_date"] = time.Now().Format(dateLayout)
	paramsOaRequest["rec_created_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

	//OA_POSTAL_ADDRESS CARD ID
	addressIDParams := make(map[string]string)
	addressIDParams["address_type"] = "17"
	addressIDParams["address_line1"] = addressid
	addressIDParams["kabupaten_key"] = kabupatenid
	addressIDParams["kecamatan_key"] = kecamatanid
	addressIDParams["postal_code"] = postalid
	addressIDParams["rec_status"] = "1"

	//OA_POSTAL_ADDRESS DOMICILE
	addressDomicileParams := make(map[string]string)
	addressDomicileParams["address_type"] = "18"
	addressDomicileParams["address_line1"] = addressdomicile
	addressDomicileParams["kabupaten_key"] = kabupatendomicile
	addressDomicileParams["kecamatan_key"] = kecamatandomicile
	addressDomicileParams["postal_code"] = postaldomicile
	addressDomicileParams["rec_status"] = "1"
	addressDomicileParams["rec_created_date"] = time.Now().Format(dateLayout)
	addressDomicileParams["rec_created_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

	//MS_BANK_ACCOUNT
	paramsBank := make(map[string]string)
	paramsBank["bank_key"] = bankKey
	paramsBank["account_no"] = accountNo
	paramsBank["account_holder_name"] = accountName
	paramsBank["branch_name"] = branchName
	paramsBank["currency_key"] = "1"
	paramsBank["bank_account_type"] = "1"
	paramsBank["rec_domain"] = "1"
	paramsBank["rec_status"] = "1"
	paramsBank["rec_created_date"] = time.Now().Format(dateLayout)
	paramsBank["rec_created_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

	//OA_PERSONAL_DATA
	log.Info("dateBirth: " + dateBirth)
	dateBirth += " 00:00:00"
	date, err = time.Parse(layout, dateBirth)
	dateStr := date.Format(layout)
	log.Info("dateBirth: " + dateStr)

	paramsOaPersonalData["full_name"] = fullname

	if nationality == "97" { //indonesia
		paramsOaPersonalData["idcard_type"] = "12"
	} else {
		paramsOaPersonalData["idcard_type"] = "13"
	}
	paramsOaPersonalData["place_birth"] = placeBirth
	paramsOaPersonalData["date_birth"] = dateStr
	paramsOaPersonalData["nationality"] = nationality
	paramsOaPersonalData["idcard_no"] = idcardNumber
	paramsOaPersonalData["gender"] = gender
	paramsOaPersonalData["marital_status"] = maritalStatus
	paramsOaPersonalData["phone_home"] = phoneHome
	paramsOaPersonalData["phone_mobile"] = phone
	paramsOaPersonalData["email_address"] = email
	paramsOaPersonalData["religion"] = religion
	paramsOaPersonalData["education"] = education
	paramsOaPersonalData["occup_job"] = job
	paramsOaPersonalData["occup_company"] = company
	paramsOaPersonalData["occup_position"] = position
	paramsOaPersonalData["beneficial_full_name"] = beneficialName
	paramsOaPersonalData["emergency_phone_no"] = emergencyPhone
	paramsOaPersonalData["relation_full_name"] = relationName
	paramsOaPersonalData["mother_maiden_name"] = motherMaidenName
	paramsOaPersonalData["emergency_full_name"] = emergencyName
	paramsOaPersonalData["rec_status"] = "1"
	paramsOaPersonalData["rec_created_date"] = time.Now().Format(dateLayout)
	paramsOaPersonalData["rec_created_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

	tx, _ := db.Db.Begin()

	//SAVE SC_USER_LOGIN
	status, err, idUserLogin := models.CreateScUserLoginReturnKey(paramsUserLogin)
	if err != nil {
		tx.Rollback()
		log.Error(err.Error())
		return lib.CustomError(http.StatusBadRequest, err.Error(), "Failed create data")
	}
	paramsOaRequest["user_login_key"] = idUserLogin
	paramsOaRequest["rec_attribute_id3"] = c.Request().UserAgent()

	//SAVE AO_PORTAL_ADDRESS IDCARD
	status, err, addressidID := models.CreateOaPostalAddress(addressIDParams)
	if err != nil {
		tx.Rollback()
		log.Error("Failed create adrress data idcard: " + err.Error())
		return lib.CustomError(status, err.Error(), "failed input data")
	}
	addressID, err := strconv.ParseUint(addressidID, 10, 64)
	if addressID == 0 {
		tx.Rollback()
		log.Error("Failed create adrress data idcard")
		return lib.CustomError(http.StatusBadGateway, "failed input data", "failed input data")
	}
	paramsOaPersonalData["idcard_address_key"] = addressidID

	//SAVE AO_PORTAL_ADDRESS DOMICILE
	status, err, addressDomicileID := models.CreateOaPostalAddress(addressDomicileParams)
	if err != nil {
		tx.Rollback()
		log.Error("Failed create adrress data domicile: " + err.Error())
		return lib.CustomError(status, err.Error(), "failed input data")
	}
	addressID, err = strconv.ParseUint(addressDomicileID, 10, 64)
	if addressID == 0 {
		tx.Rollback()
		log.Error("Failed create adrress data domicile")
		return lib.CustomError(http.StatusBadGateway, "failed input data", "failed input data")
	}
	paramsOaPersonalData["domicile_address_key"] = addressDomicileID

	//SAVE AO_PORTAL_ADDRESS COMPANY
	addressCompanyParams := make(map[string]string)
	companyAddress := c.FormValue("company_address")
	if companyAddress != "" {
		addressCompanyParams["address_type"] = "19"
		addressCompanyParams["address_line1"] = companyAddress
		addressCompanyParams["rec_status"] = "1"

		status, err, addressCompanyID := models.CreateOaPostalAddress(addressCompanyParams)
		if err != nil {
			tx.Rollback()
			log.Error("Failed create adrress data company: " + err.Error())
			return lib.CustomError(status, err.Error(), "failed input data")
		}
		addressID, err = strconv.ParseUint(addressCompanyID, 10, 64)
		if addressID == 0 {
			tx.Rollback()
			log.Error("Failed create adrress data company")
			return lib.CustomError(http.StatusBadGateway, "failed input data", "failed input data")
		}
		paramsOaPersonalData["occup_address_key"] = addressCompanyID
	}

	//SAVE MS_BANK_ACCOUNT
	status, err, bankAccountID := models.CreateMsBankAccount(paramsBank)
	if err != nil {
		tx.Rollback()
		log.Error("Failed create bank account data: " + err.Error())
		return lib.CustomError(status, err.Error(), "failed input data")
	}
	accountID, err := strconv.ParseUint(bankAccountID, 10, 64)
	if accountID == 0 {
		tx.Rollback()
		log.Error("Failed create bank account data")
		return lib.CustomError(http.StatusBadGateway, "failed input data", "failed input data")
	}
	paramsOaPersonalData["bank_account_key"] = bankAccountID

	//SAVE OA_REQUEST
	status, err, requestID := models.CreateOaRequest(paramsOaRequest)
	if err != nil {
		tx.Rollback()
		log.Error("Failed create oa request data: " + err.Error())
		return lib.CustomError(status, err.Error(), "failed input data")
	}
	request, err := strconv.ParseUint(requestID, 10, 64)
	if request == 0 {
		tx.Rollback()
		log.Error("Failed create oa request data")
		return lib.CustomError(http.StatusBadGateway, "failed input data", "failed input data")
	}

	//OA_REQUEST_BANK_ACCOUNT
	bankReq := make(map[string]string)
	bankReq["oa_request_key"] = requestID
	bankReq["bank_account_key"] = bankAccountID
	bankReq["flag_priority"] = "1"
	bankReq["bank_account_name"] = accountName
	bankReq["rec_status"] = "1"
	bankReq["rec_created_date"] = time.Now().Format(dateLayout)
	bankReq["rec_created_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
	log.Println("===== REQUEST BANK ACCOUNT ===== >>>", bankReq)
	status, err, _ = models.CreateOaRequestBankAccount(bankReq)
	if err != nil {
		tx.Rollback()
		log.Error("Failed create OA_REQUEST_BANK_ACCOUNT: " + err.Error())
		return lib.CustomError(status, err.Error(), "failed input data")
	}

	paramsOaPersonalData["oa_request_key"] = requestID

	//SAVE OA_PERSONAL_DATA
	err = os.MkdirAll(config.BasePathImage+"/images/user/"+idUserLogin, 0755)
	if err != nil {
		log.Error(err.Error())
	} else {
		var file *multipart.FileHeader
		file, err = c.FormFile("pic_ktp")
		if file != nil {
			if err != nil {
				return lib.CustomError(http.StatusBadRequest, err.Error(), "Missing required parameter: pic_ktp")
			}
			// Get file extension
			extension := filepath.Ext(file.Filename)
			// Generate filename
			var filename string
			for {
				filename = lib.RandStringBytesMaskImprSrc(20)
				log.Println("Generate filename:", filename)
				var personalData []models.OaPersonalData
				getParams := make(map[string]string)
				getParams["pic_ktp"] = filename + extension
				_, err := models.GetAllOaPersonalData(&personalData, 1, 0, getParams, false)
				if (err == nil && len(personalData) < 1) || err != nil {
					break
				}
			}
			// Upload image and move to proper directory
			err = lib.UploadImage(file, config.BasePathImage+"/images/user/"+idUserLogin+"/"+filename+extension)
			if err != nil {
				log.Println(err)
				return lib.CustomError(http.StatusInternalServerError)
			}
			paramsOaPersonalData["pic_ktp"] = filename + extension
		} else {
			return lib.CustomError(http.StatusBadRequest)
		}

		file, err = c.FormFile("pic_selfie_ktp")
		if file != nil {
			if err != nil {
				return lib.CustomError(http.StatusBadRequest, err.Error(), "Missing required parameter: pic_selfie_ktp")
			}
			// Get file extension
			extension := filepath.Ext(file.Filename)
			// Generate filename
			var filename string
			for {
				filename = lib.RandStringBytesMaskImprSrc(20)
				log.Println("Generate filename:", filename)
				var personalData []models.OaPersonalData
				getParams := make(map[string]string)
				getParams["pic_selfie_ktp"] = filename + extension
				_, err := models.GetAllOaPersonalData(&personalData, 1, 0, getParams, false)
				if (err == nil && len(personalData) < 1) || err != nil {
					break
				}
			}
			// Upload image and move to proper directory
			log.Println("===== LEWAT SINI =====")
			err = lib.UploadImage(file, config.BasePathImage+"/images/user/"+idUserLogin+"/"+filename+extension)
			if err != nil {
				log.Println(err)
				return lib.CustomError(http.StatusInternalServerError)
			}
			paramsOaPersonalData["pic_selfie_ktp"] = filename + extension
		} else {
			return lib.CustomError(http.StatusBadRequest)
		}
	}

	err = os.MkdirAll(config.BasePathImage+"/images/user/"+idUserLogin+"/signature", 0755)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(http.StatusBadGateway, err.Error(), err.Error())
	}
	file, err = c.FormFile("signature")
	if file != nil {
		if err != nil {
			return lib.CustomError(http.StatusBadRequest, err.Error(), "Missing required parameter: signature")
		}
		// Get file extension
		extension := filepath.Ext(file.Filename)
		// Generate filename
		var filename string
		for {
			filename = lib.RandStringBytesMaskImprSrc(20)
			log.Println("Generate filename:", filename)
			var personalData []models.OaPersonalData
			getParams := make(map[string]string)
			getParams["rec_image1"] = filename + extension
			_, err := models.GetAllOaPersonalData(&personalData, 1, 0, getParams, false)
			if (err == nil && len(personalData) < 1) || err != nil {
				break
			}
		}
		// Upload image and move to proper directory
		err = lib.UploadImage(file, config.BasePathImage+"/images/user/"+idUserLogin+"/signature/"+filename+extension)
		if err != nil {
			log.Println(err)
			return lib.CustomError(http.StatusInternalServerError)
		}
		paramsOaPersonalData["rec_image1"] = filename + extension
	} //else {
	// 	return lib.CustomError(http.StatusBadRequest)
	// }

	status, err, requestKey := models.CreateOaPersonalData(paramsOaPersonalData)
	if err != nil {
		tx.Rollback()
		log.Error("Failed create personal data: " + err.Error())
		return lib.CustomError(status, err.Error(), "failed input data")
	}

	var bindInterface []interface{}
	for i := 0; i < len(bindVarz); i++ {
		bindVarz[i][1] = requestKey
		bindInterface = append(bindInterface, bindVarz[i])
	}

	if len(bindInterface) > 0 {
		status, err = models.CreateMultipleUdfValue(bindInterface)
		if err != nil {
			log.Error(err.Error())
		}
	}

	//SAVE CMS_QUIZ_OPTIONS
	var questionOptions []models.QuestionOptionQuiz
	status, err = models.AdminGetQuestionOptionQuiz(&questionOptions, quizoptionkey)

	if err != nil {
		tx.Rollback()
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if len(questionOptions) < 1 {
		tx.Rollback()
		log.Error("Missing required parameter: quiz_option")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: quiz_option", "Missing required parameter: quiz_option")
	}

	var bindVar []interface{}
	var score uint64 = 0
	for _, val := range questionOptions {

		var row []string
		row = append(row, requestID)
		row = append(row, val.QuizQuestionKey)
		row = append(row, val.QuizOptionKey)
		row = append(row, strconv.FormatUint(val.QuizOptionScore, 10))
		row = append(row, "1")
		score += val.QuizOptionScore
		bindVar = append(bindVar, row)
	}

	var riskProfile models.MsRiskProfile
	scoreStr := strconv.FormatUint(score, 10)
	status, err = models.GetMsRiskProfileScore(&riskProfile, scoreStr)
	if err != nil {
		tx.Rollback()
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data risk profile")
	}

	paramsOaRiskProfile := make(map[string]string)
	paramsOaRiskProfile["oa_request_key"] = requestID
	paramsOaRiskProfile["risk_profile_key"] = strconv.FormatUint(riskProfile.RiskProfileKey, 10)
	paramsOaRiskProfile["score_result"] = scoreStr
	paramsOaRiskProfile["rec_status"] = "1"
	paramsOaRiskProfile["rec_created_date"] = time.Now().Format(dateLayout)
	paramsOaRiskProfile["rec_created_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

	status, err = models.CreateOaRiskProfile(paramsOaRiskProfile)
	if err != nil {
		tx.Rollback()
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed input data")
	}

	status, err = models.CreateMultipleOaRiskProfileQuiz(bindVar)
	if err != nil {
		tx.Rollback()
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed input data")
	}

	//create oa_request_bank_account
	paramsBankRquest := make(map[string]string)
	paramsBankRquest["oa_request_key"] = requestID
	paramsBankRquest["bank_account_key"] = bankAccountID
	paramsBankRquest["flag_priority"] = "1"
	paramsBankRquest["bank_account_name"] = accountName
	paramsBankRquest["rec_status"] = "1"
	paramsBankRquest["rec_created_date"] = time.Now().Format(dateLayout)
	paramsBankRquest["rec_created_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
	log.Println("===== REQUEST BANK ACCOUNT 2===== >>>", paramsBankRquest)

	status, err, _ = models.CreateOaRequestBankAccount(paramsBankRquest)
	if err != nil {
		tx.Rollback()
		log.Error("Failed create oa request bank acc: " + err.Error())
		return lib.CustomError(status, err.Error(), "failed input data")
	}

	tx.Commit()

	// Send email
	if config.Envi != "DEV" {

		//sent email to customer
		t := template.New("index-registration.html")

		t, err = t.ParseFiles(config.BasePath + "/mail/index-registration.html")
		if err != nil {
			log.Println(err)
		}

		var tpl bytes.Buffer
		if err := t.Execute(&tpl, struct {
			Name    string
			FileUrl string
		}{Name: fullname, FileUrl: config.FileUrl + "/images/mail"}); err != nil {
			log.Println(err)
		}

		result := tpl.String()

		mailer := gomail.NewMessage()
		mailer.SetHeader("From", config.EmailFrom)
		mailer.SetHeader("To", email)
		mailer.SetHeader("Subject", "[MotionFunds] Pembukaan Rekening Kamu sedang Diproses")
		mailer.SetBody("text/html", result)

		err = lib.SendEmail(mailer)
		if err != nil {
			log.Error("Error send email")
			log.Error(err)
		}
		// dialer := gomail.NewDialer(
		// 	config.EmailSMTPHost,
		// 	int(config.EmailSMTPPort),
		// 	config.EmailFrom,
		// 	config.EmailFromPassword,
		// )
		// dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

		// err = dialer.DialAndSend(mailer)
		// if err != nil {
		// 	log.Error("Error send email")
		// 	log.Error(err)
		// }

		//sent email to all CS & Agent
		var oaRequest models.OaRequest
		status, err = models.GetOaRequest(&oaRequest, requestID)
		if err != nil {
			log.Error(err.Error())
		} else {
			var personalData models.OaPersonalData
			status, err = models.GetOaPersonalData(&personalData, requestID, "oa_request_key")
			if err != nil {
				log.Error(err.Error())
			} else {
				SentEmailOaPengkinianToBackOffice(oaRequest, personalData, "11")
				SentEmailOaPengkinianToSales(oaRequest, personalData)
			}
		}
	}

	//insert message notif in app
	strIDUserLogin := strconv.FormatUint(lib.Profile.UserID, 10)
	paramsUserMessage := make(map[string]string)
	paramsUserMessage["umessage_type"] = "245"
	paramsUserMessage["umessage_recipient_key"] = idUserLogin
	paramsUserMessage["umessage_receipt_date"] = time.Now().Format(dateLayout)
	paramsUserMessage["flag_read"] = "0"
	paramsUserMessage["umessage_sent_date"] = time.Now().Format(dateLayout)
	paramsUserMessage["flag_sent"] = "1"
	subject := "Pembukaan Rekening sedang Diproses"
	body := "Terima kasih telah mendaftar. Kami sedang melakukan proses verifikasi data kamu max. 1X24 jam. Mohon ditunggu ya."
	paramsUserMessage["umessage_subject"] = subject
	paramsUserMessage["umessage_body"] = body

	paramsUserMessage["umessage_category"] = "248"
	paramsUserMessage["flag_archieved"] = "0"
	paramsUserMessage["archieved_date"] = time.Now().Format(dateLayout)
	paramsUserMessage["rec_status"] = "1"
	paramsUserMessage["rec_created_date"] = time.Now().Format(dateLayout)
	paramsUserMessage["rec_created_by"] = strIDUserLogin

	status, err = models.CreateScUserMessage(paramsUserMessage)
	if err != nil {
		log.Error("Error create user message")
	} else {
		log.Error("Sukses insert user message")
	}

	var responseData models.MsRiskProfileInfo

	responseData.RiskProfileKey = riskProfile.RiskProfileKey
	responseData.RiskCode = riskProfile.RiskCode
	responseData.RiskName = riskProfile.RiskName
	responseData.RiskDesc = riskProfile.RiskDesc
	responseData.Score = score

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = responseData

	return c.JSON(http.StatusOK, response)
}

func GetAdminListCustomerRedemption(c echo.Context) error {

	var err error
	var status int

	params := make(map[string]string)

	branchKey := c.QueryParam("branch_key")
	if branchKey != "" {
		params["c.openacc_branch_key"] = branchKey
	} else {
		//if user category  = 3 -> user branch, 2 = user HO
		var userCategory uint64
		userCategory = 3
		if lib.Profile.UserCategoryKey == userCategory {
			log.Println(lib.Profile)
			if lib.Profile.BranchKey != nil {
				strBranchKey := strconv.FormatUint(*lib.Profile.BranchKey, 10)
				params["c.openacc_branch_key"] = strBranchKey
			} else {
				log.Error("User Branch haven't Branch")
				return lib.CustomError(http.StatusBadRequest, "Wrong User Branch haven't Branch", "Wrong User Branch haven't Branch")
			}
		}
	}

	paramsLike := make(map[string]string)

	var customer []models.CustomerDropdown

	status, err = models.GetCustomerRedemptionDropdown(&customer, params, paramsLike)

	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if len(customer) < 1 {
		log.Error("Customer not found")
		return lib.CustomError(http.StatusNotFound, "Customer not found", "Customer not found")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = customer

	return c.JSON(http.StatusOK, response)
}

func GetAdminOaRequestPersonalDataRiskProfile(c echo.Context) error {
	customerKey := c.Param("key")
	key, _ := strconv.ParseUint(customerKey, 10, 64)
	if key == 0 {
		return lib.CustomError(http.StatusNotFound)
	}

	var oaRequestDB []models.OaRequest
	params := make(map[string]string)
	params["customer_key"] = customerKey
	params["orderBy"] = "oa_request_key"
	params["orderType"] = "DESC"
	status, err := models.GetAllOaRequest(&oaRequestDB, 0, 0, true, params)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Oa Request not found")
	}
	var requestKey string
	var oaData models.OaRequest
	if len(oaRequestDB) > 0 {
		oaData = oaRequestDB[0]
		if *oaData.Oastatus == uint64(258) || *oaData.Oastatus == uint64(259) {
			log.Error("oa in progress approval")
			return lib.CustomError(http.StatusNotFound, "Terdapat Data Request yang dalam approval. Mohon Tunggu proses approval untuk dapat melakukan pengkinian lagi.", "Terdapat Data Request yang dalam approval. Mohon Tunggu proses approval untuk dapat melakukan pengkinian lagi.")
		} else {
			requestKey = strconv.FormatUint(oaData.OaRequestKey, 10)
		}
	} else {
		log.Error("oa not found")
		return lib.CustomError(http.StatusNotFound, "Oa Request not found", "Oa Request not found")
	}

	var personalDataDB models.OaPersonalData
	if requestKey != "" {
		status, err = models.GetOaPersonalData(&personalDataDB, requestKey, "oa_request_key")
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Oa Request not found")
		}
	}

	responseData := make(map[string]interface{})
	responseData["full_name"] = personalDataDB.FullName
	if oaData.CustomerKey != nil {
		responseData["customer_key"] = *oaData.CustomerKey
	} else {
		responseData["customer_key"] = nil
	}
	if oaData.BranchKey != nil {
		responseData["branch_key"] = *oaData.BranchKey
	} else {
		responseData["branch_key"] = 1
	}
	if oaData.AgentKey != nil {
		responseData["agent_key"] = *oaData.AgentKey
	} else {
		responseData["agent_key"] = 1
	}

	responseData["user_login_key"] = oaData.UserLoginKey
	responseData["place_birth"] = personalDataDB.PlaceBirth
	responseData["date_birth"] = personalDataDB.DateBirth
	responseData["nationality"] = personalDataDB.Nationality
	responseData["idcard_type"] = personalDataDB.IDcardType
	responseData["idcard_no"] = personalDataDB.IDcardNo
	responseData["idcard_expired_date"] = personalDataDB.IDcardExpiredDate
	responseData["idcard_never_expired"] = personalDataDB.IDcardNeverExpired
	responseData["gender"] = personalDataDB.Gender
	responseData["marital_status"] = personalDataDB.MaritalStatus
	responseData["pep_status"] = personalDataDB.PepStatus
	var address models.OaPostalAddress
	_, err = models.GetOaPostalAddress(&address, strconv.FormatUint(*personalDataDB.IDcardAddressKey, 10))
	if err == nil {
		addressID := make(map[string]interface{})
		var city models.MsCity
		_, err = models.GetMsCityByParent(&city, strconv.FormatUint(*address.KabupatenKey, 10))
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
		addressID["provinsi_key"] = city.CityKey
		addressID["postal_address_key"] = address.PostalAddressKey
		addressID["kabupaten_key"] = address.KabupatenKey
		addressID["kecamatan_key"] = address.KecamatanKey
		addressID["address_line1"] = address.AddressLine1
		addressID["address_line2"] = address.AddressLine2
		addressID["address_line3"] = address.AddressLine3
		addressID["postal_code"] = address.PostalCode
		responseData["idcard_address"] = addressID
	}
	_, err = models.GetOaPostalAddress(&address, strconv.FormatUint(*personalDataDB.DomicileAddressKey, 10))
	if err == nil {
		var city models.MsCity
		_, err = models.GetMsCityByParent(&city, strconv.FormatUint(*address.KabupatenKey, 10))
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
		addressID := make(map[string]interface{})
		addressID["provinsi_key"] = city.CityKey
		addressID["postal_address_key"] = address.PostalAddressKey
		addressID["kabupaten_key"] = address.KabupatenKey
		addressID["kecamatan_key"] = address.KecamatanKey
		addressID["address_line1"] = address.AddressLine1
		addressID["address_line2"] = address.AddressLine2
		addressID["address_line3"] = address.AddressLine3
		addressID["postal_code"] = address.PostalCode
		responseData["domicile_address"] = addressID
	}

	dir := config.ImageUrl + "/images/user/" + strconv.FormatUint(*oaData.UserLoginKey, 10) + "/"

	responseData["phone_home"] = personalDataDB.PhoneHome
	responseData["phone_mobile"] = personalDataDB.PhoneMobile
	responseData["email"] = personalDataDB.EmailAddress
	responseData["religion"] = personalDataDB.Religion
	if *personalDataDB.Religion == uint64(26) {
		var ud models.UdfOtherValueStruct
		_, err = models.UdfOtherValueQuery(&ud, personalDataDB.PersonalDataKey, "1")
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
		responseData["religion_other"] = ud
	}
	responseData["pic_selfie"] = personalDataDB.PicSelfie
	responseData["pic_ktp"] = personalDataDB.PicKtp
	responseData["pic_selfie_ktp"] = personalDataDB.PicSelfieKtp
	responseData["signature"] = personalDataDB.RecImage1

	if personalDataDB.PicSelfie != nil && *personalDataDB.PicSelfie != "" {
		responseData["pic_selfie_path"] = dir + *personalDataDB.PicSelfie
	} else {
		responseData["pic_selfie_path"] = ""
	}

	if personalDataDB.PicKtp != nil && *personalDataDB.PicKtp != "" {
		responseData["pic_ktp_path"] = dir + *personalDataDB.PicKtp
	} else {
		responseData["pic_ktp_path"] = ""
	}

	if personalDataDB.PicSelfieKtp != nil && *personalDataDB.PicSelfieKtp != "" {
		responseData["pic_selfie_ktp_path"] = dir + *personalDataDB.PicSelfieKtp
	} else {
		responseData["pic_selfie_ktp_path"] = ""
	}

	if personalDataDB.RecImage1 != nil && *personalDataDB.RecImage1 != "" {
		responseData["signature_path"] = dir + "signature/" + *personalDataDB.RecImage1
	} else {
		responseData["signature_path"] = ""
	}

	responseData["education"] = personalDataDB.Education
	if *personalDataDB.Education == uint64(43) {
		var ud models.UdfOtherValueStruct
		_, err = models.UdfOtherValueQuery(&ud, personalDataDB.PersonalDataKey, "3")
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
		responseData["education_other"] = ud
	}
	responseData["occup_job"] = personalDataDB.OccupJob
	if *personalDataDB.OccupJob == uint64(35) {
		var ud models.UdfOtherValueStruct
		_, err = models.UdfOtherValueQuery(&ud, personalDataDB.PersonalDataKey, "2")
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
		responseData["occup_other"] = ud
	}
	responseData["occup_company"] = personalDataDB.OccupCompany
	responseData["occup_position"] = personalDataDB.OccupPosition
	if personalDataDB.OccupAddressKey != nil {
		_, err = models.GetOaPostalAddress(&address, strconv.FormatUint(*personalDataDB.OccupAddressKey, 10))
		if err == nil {
			log.Println("===== LEWAT SINI ===== >>")
			addressID := make(map[string]interface{})
			addressID["postal_address_key"] = address.PostalAddressKey
			addressID["kabupaten_key"] = address.KabupatenKey
			addressID["kecamatan_key"] = address.KecamatanKey
			addressID["address_line1"] = address.AddressLine1
			addressID["address_line2"] = address.AddressLine2
			addressID["address_line3"] = address.AddressLine3
			addressID["postal_code"] = address.PostalCode
			responseData["occup_address"] = addressID
		}
	}
	responseData["occup_business_field"] = personalDataDB.OccupBusinessFields
	if *personalDataDB.OccupBusinessFields == uint64(60) {
		var ud models.UdfOtherValueStruct
		_, err = models.UdfOtherValueQuery(&ud, personalDataDB.PersonalDataKey, "4")
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
		responseData["occup_business_field_other"] = ud
	}
	responseData["occup_phone"] = personalDataDB.OccupPhone
	responseData["occup_web_url"] = personalDataDB.OccupWebUrl
	responseData["correspondence"] = personalDataDB.Correspondence
	responseData["annual_income"] = personalDataDB.AnnualIncome
	responseData["sourceof_fund"] = personalDataDB.SourceofFund
	if *personalDataDB.SourceofFund == uint64(76) {
		var ud models.UdfOtherValueStruct
		_, err = models.UdfOtherValueQuery(&ud, personalDataDB.PersonalDataKey, "5")
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
		responseData["sourceof_fund_other"] = ud
	}
	responseData["invesment_objectives"] = personalDataDB.InvesmentObjectives
	if *personalDataDB.InvesmentObjectives == uint64(81) {
		var ud models.UdfOtherValueStruct
		_, err = models.UdfOtherValueQuery(&ud, personalDataDB.PersonalDataKey, "6")
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
		responseData["invesment_objectives_other"] = ud
	}
	responseData["relation_type"] = personalDataDB.RelationType
	responseData["relation_full_name"] = personalDataDB.RelationFullName
	responseData["relation_occupation"] = personalDataDB.RelationOccupation
	if *personalDataDB.RelationOccupation == uint64(157) {
		var ud models.UdfOtherValueStruct
		_, err = models.UdfOtherValueQuery(&ud, personalDataDB.PersonalDataKey, "8")
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
		responseData["relation_occupation_other"] = ud
	}
	responseData["relation_business_fields"] = personalDataDB.RelationBusinessFields
	if *personalDataDB.RelationOccupation == uint64(157) {
		var ud models.UdfOtherValueStruct
		_, err = models.UdfOtherValueQuery(&ud, personalDataDB.PersonalDataKey, "8")
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
		responseData["relation_occupation_other"] = ud
	}
	responseData["mother_maiden_name"] = personalDataDB.MotherMaidenName
	responseData["emergency_full_name"] = personalDataDB.EmergencyFullName
	responseData["emergency_relation"] = personalDataDB.EmergencyRelation
	responseData["emergency_phone_no"] = personalDataDB.EmergencyPhoneNo
	responseData["beneficial_full_name"] = personalDataDB.BeneficialFullName
	responseData["beneficial_relation"] = personalDataDB.BeneficialRelation
	if *personalDataDB.BeneficialRelation == uint64(106) {
		var ud models.UdfOtherValueStruct
		_, err = models.UdfOtherValueQuery(&ud, personalDataDB.PersonalDataKey, "7")
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
		responseData["beneficial_relation_other"] = ud
	}

	//get all bank request
	var accBank []models.OaRequestByField
	status, err = models.GetOaRequestBankByField(&accBank, "oa_request_key", strconv.FormatUint(personalDataDB.OaRequestKey, 10))
	if err != nil {
		responseData["bank_request"] = nil
	} else {
		if len(accBank) > 0 {
			var bankAccount []interface{}
			for _, value := range accBank {
				bank := make(map[string]interface{})
				bank["bank_key"] = value.BankKey
				bank["account_no"] = value.AccountNo
				bank["account_holder_name"] = value.AccountHolderName
				bank["branch_name"] = value.BranchName
				bank["flag_priority"] = value.FlagPriority
				bankAccount = append(bankAccount, bank)
			}
			responseData["bank_request"] = bankAccount
		} else {
			responseData["bank_request"] = nil
		}
	}

	var requestDB []models.OaRequest
	paramRequest := make(map[string]string)
	paramRequest["customer_key"] = customerKey
	paramRequest["orderBy"] = "oa_request_key"
	paramRequest["orderType"] = "DESC"
	_, err = models.GetAllOaRequest(&requestDB, 1, 0, false, paramRequest)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	request := requestDB[0]
	var quizDB []models.OaRiskProfileQuiz
	paramQuiz := make(map[string]string)
	paramQuiz["oa_request_key"] = strconv.FormatUint(request.OaRequestKey, 10)
	paramQuiz["orderBy"] = "oa_request_key"
	paramQuiz["orderType"] = "DESC"
	_, err = models.GetAllOaRiskProfileQuiz(&quizDB, 100, 0, paramQuiz, true)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if len(quizDB) < 1 {
		_, err = models.RiskProfileQuizIfNull(&quizDB, request.OaRequestKey)
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
	}

	var risk models.OaRiskProfile
	_, err = models.GetOaRiskProfile(&risk, strconv.FormatUint(request.OaRequestKey, 10), "oa_request_key")
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}

	var riskProfile models.MsRiskProfile
	_, err = models.GetMsRiskProfile(&riskProfile, strconv.FormatUint(risk.RiskProfileKey, 10))
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}

	riskProfileData := make(map[string]interface{})
	riskProfileData["score_result"] = risk.ScoreResult
	riskProfileData["risk_code"] = riskProfile.RiskCode
	riskProfileData["risk_name"] = riskProfile.RiskName
	riskProfileData["risk_desc"] = riskProfile.RiskDesc
	var quizData []interface{}
	for _, q := range quizDB {
		quiz := make(map[string]interface{})
		quiz["question_key"] = q.QuizQuestionKey
		quiz["option_key"] = q.QuizOptionKey
		quiz["option_score"] = q.QuizOptionScore
		quizData = append(quizData, quiz)
	}
	riskProfileData["quiz"] = quizData

	responseData["risk_profile"] = riskProfileData
	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = responseData
	return c.JSON(http.StatusOK, response)
}

func AdminSavePengkinianCustomerIndividu(c echo.Context) error {
	var err error
	var status int
	var bindVarz [][]string

	paramsOaPersonalData := make(map[string]string)

	//DEFAULT PARAM
	customerKey := c.FormValue("customer_key")
	if customerKey == "" {
		log.Error("Missing required parameter: customer_key")
		return lib.CustomError(http.StatusBadRequest, "customer_key can not be blank", "customer_key can not be blank")
	}

	var oaRequestDB []models.OaRequest
	params := make(map[string]string)
	params["customer_key"] = customerKey
	params["orderBy"] = "oa_request_key"
	params["orderType"] = "DESC"
	status, err = models.GetAllOaRequest(&oaRequestDB, 0, 0, true, params)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Customer not found")
	}

	var oaData models.OaRequest
	if len(oaRequestDB) > 0 {
		oaData = oaRequestDB[0]
		if *oaData.Oastatus == uint64(258) || *oaData.Oastatus == uint64(259) {
			log.Error("oa in progress approval")
			return lib.CustomError(http.StatusNotFound, "Terdapat Data Request yang dalam approval. Mohon Tunggu proses approval untuk dapat melakukan pengkinian lagi.", "Terdapat Data Request yang dalam approval. Mohon Tunggu proses approval untuk dapat melakukan pengkinian lagi.")
		}
	} else {
		log.Error("oa not found")
		return lib.CustomError(http.StatusNotFound, "Customer not found", "Customer not found")
	}

	oaRequestType := c.FormValue("oa_request_type")
	if oaRequestType == "" {
		log.Error("Missing required parameter: oa_request_type")
		return lib.CustomError(http.StatusBadRequest, "oa_request_type can not be blank", "oa_request_type can not be blank")
	}

	branchkey := c.FormValue("branch_key")
	if branchkey == "" {
		log.Error("Missing required parameter: branch_key")
		return lib.CustomError(http.StatusBadRequest, "branch_key can not be blank", "branch_key can not be blank")
	} else {
		n, err := strconv.ParseUint(branchkey, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: branch_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: branch_key", "Wrong input for parameter: branch_key")
		}
	}

	agentkey := c.FormValue("agent_key")
	if agentkey == "" {
		log.Error("Missing required parameter: agent_key")
		return lib.CustomError(http.StatusBadRequest, "agent_key can not be blank", "agent_key can not be blank")
	} else {
		n, err := strconv.ParseUint(agentkey, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: agent_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: agent_key", "Wrong input for parameter: agent_key")
		}
	}
	salesCode := c.FormValue("sales_code")

	//INFORMASI NASABAH
	email := c.FormValue("email")
	if email == "" {
		log.Error("Missing required parameter: email")
		return lib.CustomError(http.StatusBadRequest, "email can not be blank", "email can not be blank")
	}
	phone := c.FormValue("phone")
	if phone == "" {
		log.Error("Missing required parameter: phone")
		return lib.CustomError(http.StatusBadRequest, "phone can not be blank", "phone can not be blank")
	}
	fullname := c.FormValue("full_name")
	if fullname == "" {
		log.Error("Missing required parameter: full_name")
		return lib.CustomError(http.StatusBadRequest, "full_name can not be blank", "full_name can not be blank")
	}

	nationality := c.FormValue("nationality")
	if nationality == "" {
		log.Error("Missing required parameter: nationality")
		return lib.CustomError(http.StatusBadRequest, "nationality can not be blank", "nationality can not be blank")
	} else {
		n, err := strconv.ParseUint(nationality, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: nationality")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: nationality", "Wrong input for parameter: nationality")
		}
	}

	idcardNumber := c.FormValue("idcard_no")
	if idcardNumber == "" {
		log.Error("Missing required parameter: idcard_no")
		return lib.CustomError(http.StatusBadRequest, "idcard_no can not be blank", "idcard_no can not be blank")
	}

	gender := c.FormValue("gender")
	if gender == "" {
		log.Error("Missing required parameter: gender")
		return lib.CustomError(http.StatusBadRequest, "gender can not be blank", "gender can not be blank")
	} else {
		n, err := strconv.ParseUint(gender, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: gender")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: gender", "Wrong input for parameter: gender")
		}
	}

	placeBirth := c.FormValue("place_birth")
	if placeBirth == "" {
		log.Error("Missing required parameter: place_birth")
		return lib.CustomError(http.StatusBadRequest, "place_birth can not be blank", "place_birth can not be blank")
	}

	dateBirth := c.FormValue("date_birth")
	if dateBirth == "" {
		log.Error("Missing required parameter: date_birth")
		return lib.CustomError(http.StatusBadRequest, "date_birth can not be blank", "date_birth can not be blank")
	}

	maritalStatus := c.FormValue("marital_status")
	if maritalStatus == "" {
		log.Error("Missing required parameter: marital_status")
		return lib.CustomError(http.StatusBadRequest, "marital_status can not be blank", "marital_status can not be blank")
	} else {
		n, err := strconv.ParseUint(maritalStatus, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: marital_status")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: marital_status", "Wrong input for parameter: marital_status")
		}
	}

	addressid := c.FormValue("address_id")
	if addressid == "" {
		log.Error("Missing required parameter: address_id")
		return lib.CustomError(http.StatusBadRequest, "address_id can not be blank", "address_id can not be blank")
	}

	kabupatenid := c.FormValue("kabupaten_id")
	if kabupatenid != "" {
		n, err := strconv.ParseUint(kabupatenid, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: kabupaten_id")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: kabupaten_id", "Wrong input for parameter: kabupaten_id")
		}
	} else {
		log.Error("Missing required parameter: kabupaten_id")
		return lib.CustomError(http.StatusBadRequest, "kabupaten_id can not be blank", "kabupaten_id can not be blank")
	}

	kecamatanid := c.FormValue("kecamatan_id")
	if kecamatanid != "" {
		n, err := strconv.ParseUint(kecamatanid, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: kecamatan_id")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: kecamatan_id", "Wrong input for parameter: kecamatan_id")
		}
	} else {
		log.Error("Missing required parameter: kecamatan_id")
		return lib.CustomError(http.StatusBadRequest, "kecamatan_id can not be blank", "kecamatan_id can not be blank")
	}

	postalid := c.FormValue("postal_id")
	if postalid != "" {
		ps, err := strconv.ParseUint(postalid, 10, 64)
		if err != nil || ps == 0 {
			log.Error("Wrong input for parameter: postal_id")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: postal_id", "Wrong input for parameter: postal_id")
		}
	} else {
		log.Error("Missing required parameter: postal_id")
		return lib.CustomError(http.StatusBadRequest, "postal_id can not be blank", "postal_id can not be blank")
	}

	addressdomicile := c.FormValue("address_domicile")
	if addressdomicile == "" {
		log.Error("Missing required parameter: address_domicile")
		return lib.CustomError(http.StatusBadRequest, "address_domicile can not be blank", "address_domicile can not be blank")
	}

	kabupatendomicile := c.FormValue("kabupaten_domicile")
	if kabupatendomicile != "" {
		n, err := strconv.ParseUint(kabupatendomicile, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: kabupaten_domicile")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: kabupaten_domicile", "Wrong input for parameter: kabupaten_domicile")
		}
	} else {
		log.Error("Missing required parameter: kabupaten_domicile")
		return lib.CustomError(http.StatusBadRequest, "kabupaten_domicile can not be blank", "kabupaten_domicile can not be blank")
	}

	kecamatandomicile := c.FormValue("kecamatan_domicile")
	if kecamatandomicile != "" {
		n, err := strconv.ParseUint(kecamatandomicile, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: kecamatan_domicile")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: kecamatan_domicile", "Wrong input for parameter: kecamatan_domicile")
		}
	} else {
		log.Error("Missing required parameter: kecamatan_domicile")
		return lib.CustomError(http.StatusBadRequest, "kecamatan_domicile can not be blank", "kecamatan_domicile can not be blank")
	}

	postaldomicile := c.FormValue("postal_domicile")
	if postaldomicile != "" {
		ps, err := strconv.ParseUint(postaldomicile, 10, 64)
		if err != nil || ps == 0 {
			log.Error("Wrong input for parameter: postal_domicile")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: postal_domicile", "Wrong input for parameter: postal_domicile")
		}
	} else {
		log.Error("Missing required parameter: postal_domicile")
		return lib.CustomError(http.StatusBadRequest, "postal_domicile can not be blank", "postal_domicile can not be blank")
	}

	phoneHome := c.FormValue("phone_home")
	if phoneHome == "" {
		log.Error("Missing required parameter: phone_home")
		return lib.CustomError(http.StatusBadRequest, "phone_home can not be blank", "phone_home can not be blank")
	}

	religion := c.FormValue("religion")
	if religion == "" {
		log.Error("Missing required parameter: phone_home")
		return lib.CustomError(http.StatusBadRequest, "phone_home can not be blank", "phone_home can not be blank")
	} else {
		ps, err := strconv.ParseUint(religion, 10, 64)
		if err != nil || ps == 0 {
			log.Error("Wrong input for parameter: religion")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: religion", "Wrong input for parameter: religion")
		}
	}

	religionOther := c.FormValue("religion_other")
	if religionOther != "" {
		var rowz []string
		rowz = append(rowz, "1")
		rowz = append(rowz, "0")
		rowz = append(rowz, religionOther)
		bindVarz = append(bindVarz, rowz)
	}

	education := c.FormValue("education")
	if education == "" {
		log.Error("Missing required parameter: education")
		return lib.CustomError(http.StatusBadRequest, "education can not be blank", "education can not be blank")
	} else {
		ps, err := strconv.ParseUint(education, 10, 64)
		if err != nil || ps == 0 {
			log.Error("Wrong input for parameter: education")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: education", "Wrong input for parameter: education")
		}
	}

	educationOther := c.FormValue("education_other")
	if educationOther != "" {
		var rowz []string
		rowz = append(rowz, "3")
		rowz = append(rowz, "0")
		rowz = append(rowz, educationOther)
		bindVarz = append(bindVarz, rowz)
	}

	//UPLOAD DOKUMEN FOTO E-KTP & SELFIE DENGAN KTP

	picKtpDefault := c.FormValue("pic_ktp_default")
	picSelfieKtpDefault := c.FormValue("pic_selfie_ktp_default")
	signatureDefault := c.FormValue("signature_default")

	var file *multipart.FileHeader
	file, err = c.FormFile("pic_ktp")
	if file == nil && picKtpDefault == "" {
		log.Error("Missing required parameter: pic_ktp")
		return lib.CustomError(http.StatusBadRequest, "pic_ktp can not be blank", "pic_ktp can not be blank")
	}

	var fileselfie *multipart.FileHeader
	fileselfie, err = c.FormFile("pic_selfie_ktp")
	if fileselfie == nil && picSelfieKtpDefault == "" {
		log.Error("Missing required parameter: pic_selfie_ktp")
		return lib.CustomError(http.StatusBadRequest, "pic_selfie_ktp can not be blank", "pic_selfie_ktp can not be blank")
	}

	//URAIAN BIDANG USAHA DAN PEKERJAAN
	job := c.FormValue("job")
	if job == "" {
		log.Error("Missing required parameter: job")
		return lib.CustomError(http.StatusBadRequest, "job can not be blank", "job can not be blank")
	} else {
		n, err := strconv.ParseUint(job, 10, 64)
		if err == nil && n > 0 {
			paramsOaPersonalData["occup_job"] = job
		} else {
			log.Error("Wrong input for parameter: job")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: job", "Wrong input for parameter: job")
		}
	}

	jobOther := c.FormValue("job_other")
	if jobOther != "" {
		var rowz []string
		rowz = append(rowz, "2")
		rowz = append(rowz, "0")
		rowz = append(rowz, jobOther)
		bindVarz = append(bindVarz, rowz)
	}

	company := c.FormValue("company")
	if company != "" {
		paramsOaPersonalData["occup_company"] = company
	}

	position := c.FormValue("position")
	if position != "" {
		n, err := strconv.ParseUint(job, 10, 64)
		if err == nil && n > 0 {
			paramsOaPersonalData["occup_position"] = position
		} else {
			log.Error("Wrong input for parameter: position")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: position", "Wrong input for parameter: position")
		}
	}

	positionOther := c.FormValue("position_other")
	if positionOther != "" {
		var rowz []string
		rowz = append(rowz, "10")
		rowz = append(rowz, "0")
		rowz = append(rowz, positionOther)
		bindVarz = append(bindVarz, rowz)
	}

	businessField := c.FormValue("business_field")
	if businessField != "" {
		n, err := strconv.ParseUint(businessField, 10, 64)
		if err == nil && n > 0 {
			paramsOaPersonalData["occup_business_fields"] = businessField
		} else {
			log.Error("Wrong input for parameter: business_field")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: business_field", "Wrong input for parameter: business_field")
		}
	}

	businessFieldOther := c.FormValue("business_field_other")
	if businessFieldOther != "" {
		var rowz []string
		rowz = append(rowz, "9")
		rowz = append(rowz, "0")
		rowz = append(rowz, businessFieldOther)
		bindVarz = append(bindVarz, rowz)
	}

	annualIncome := c.FormValue("annual_income")
	if annualIncome == "" {
		log.Error("Missing required parameter: annual_income")
		return lib.CustomError(http.StatusBadRequest, "annual_income can not be blank", "annual_income can not be blank")
	} else {
		n, err := strconv.ParseUint(annualIncome, 10, 64)
		if err == nil && n > 0 {
			paramsOaPersonalData["annual_income"] = annualIncome
		} else {
			log.Error("Wrong input for parameter: annual_income")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: annual_income", "Wrong input for parameter: annual_income")
		}
	}

	fundSource := c.FormValue("fund_source")
	if fundSource == "" {
		log.Error("Missing required parameter: fund_source")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: fund_source", "Missing required parameter: fund_source")
	} else {
		n, err := strconv.ParseUint(fundSource, 10, 64)
		if err == nil && n > 0 {
			paramsOaPersonalData["sourceof_fund"] = fundSource
		} else {
			log.Error("Wrong input for parameter: fund_source")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: fund_source", "Wrong input for parameter: fund_source")
		}
	}

	fundSourceOther := c.FormValue("fund_source_other")
	if fundSourceOther != "" {
		var row []string
		row = append(row, "5")
		row = append(row, "0")
		row = append(row, fundSourceOther)
		bindVarz = append(bindVarz, row)
	}

	objectives := c.FormValue("objectives")
	if objectives == "" {
		log.Error("Missing required parameter: objectives")
		return lib.CustomError(http.StatusBadRequest, "objectives can not be blank", "objectives can not be blank")
	} else {
		n, err := strconv.ParseUint(objectives, 10, 64)
		if err == nil && n > 0 {
			paramsOaPersonalData["invesment_objectives"] = objectives
		} else {
			log.Error("Wrong input for parameter: objectives")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: objectives", "Wrong input for parameter: objectives")
		}
	}

	objectivesOther := c.FormValue("objectives_other")
	if objectivesOther != "" {
		var row []string
		row = append(row, "6")
		row = append(row, "0")
		row = append(row, objectivesOther)
		bindVarz = append(bindVarz, row)
	}

	corespondence := c.FormValue("corespondence")
	if corespondence != "" {
		n, err := strconv.ParseUint(corespondence, 10, 64)
		if err == nil && n > 0 {
			paramsOaPersonalData["correspondence"] = corespondence
		} else {
			log.Error("Wrong input for parameter: corespondence")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: corespondence", "Wrong input for parameter: corespondence")
		}
	}

	//TAB 4
	motherMaidenName := c.FormValue("mother_maiden_name")
	if motherMaidenName == "" {
		log.Error("Missing required parameter: mother_maiden_name")
		return lib.CustomError(http.StatusBadRequest, "mother_maiden_name can not be blank", "mother_maiden_name can not be blank")
	}

	relationName := c.FormValue("relation_name")
	if relationName == "" {
		log.Error("Missing required parameter: relation_name")
		return lib.CustomError(http.StatusBadRequest, "relation_name can not be blank", "relation_name can not be blank")
	}

	pepStatus := c.FormValue("pep_status")
	if pepStatus == "" {
		log.Error("Missing required parameter: pep_status")
		return lib.CustomError(http.StatusBadRequest, "pep_status can not be blank", "pep_status can not be blank")
	} else {
		n, err := strconv.ParseUint(pepStatus, 10, 64)
		if err == nil && n > 0 {
			paramsOaPersonalData["pep_status"] = pepStatus
		} else {
			log.Error("Wrong input for parameter: pep_status")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: pep_status", "Wrong input for parameter: pep_status")
		}
	}

	relationOccupation := c.FormValue("relation_occupation")
	if relationOccupation == "" {
		log.Error("Missing required parameter: relation_occupation")
		return lib.CustomError(http.StatusBadRequest, "relation_occupation can not be blank", "relation_occupation can not be blank")
	} else {
		n, err := strconv.ParseUint(relationOccupation, 10, 64)
		if err == nil && n > 0 {
			paramsOaPersonalData["relation_occupation"] = relationOccupation
		} else {
			log.Error("Wrong input for parameter: relation_occupation")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: relation_occupation", "Wrong input for parameter: relation_occupation")
		}
	}

	relationBusinessFieldOther := c.FormValue("relation_business_field_other")
	if relationBusinessFieldOther != "" {
		var row []string
		row = append(row, "10")
		row = append(row, "0")
		row = append(row, relationBusinessFieldOther)
		bindVarz = append(bindVarz, row)
	}

	relationType := c.FormValue("relation_type")
	if relationType == "" {
		log.Error("Missing required parameter: relation_type")
		return lib.CustomError(http.StatusBadRequest, "relation_type can not be blank", "relation_type can not be blank")
	} else {
		n, err := strconv.ParseUint(relationType, 10, 64)
		if err == nil && n > 0 {
			paramsOaPersonalData["relation_type"] = relationType
		} else {
			log.Error("Wrong input for parameter: relation_type")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: relation_type", "Wrong input for parameter: relation_type")
		}
	}

	relationBusinessField := c.FormValue("relation_business_field")
	if relationBusinessField != "" {
		n, err := strconv.ParseUint(relationBusinessField, 10, 64)
		if err == nil && n > 0 {
			paramsOaPersonalData["relation_business_fields"] = relationBusinessField
		} else {
			log.Error("Wrong input for parameter: relation_business_field")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: relation_business_field", "Wrong input for parameter: relation_business_field")
		}
	}

	emergencyName := c.FormValue("emergency_name")
	if emergencyName == "" {
		log.Error("Missing required parameter: emergency_name")
		return lib.CustomError(http.StatusBadRequest, "emergency_name can not be blank", "emergency_name can not be blank")
	}

	emergencyRelation := c.FormValue("emergency_relation")
	if emergencyRelation == "" {
		log.Error("Missing required parameter: emergency_relation")
		return lib.CustomError(http.StatusBadRequest, "emergency_relation can not be blank", "emergency_relation can not be blank")
	} else {
		n, err := strconv.ParseUint(emergencyRelation, 10, 64)
		if err == nil && n > 0 {
			paramsOaPersonalData["emergency_relation"] = emergencyRelation
		} else {
			log.Error("Wrong input for parameter: emergency_relation")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: emergency_relation", "Wrong input for parameter: emergency_relation")
		}
	}

	emergencyPhone := c.FormValue("emergency_phone")
	if emergencyPhone == "" {
		log.Error("Missing required parameter: emergency_phone")
		return lib.CustomError(http.StatusBadRequest, "emergency_phone can not be blank", "emergency_phone can not be blank")
	}

	//TAB 5 REKENING DLL
	beneficialRelation := c.FormValue("beneficial_relation")
	if beneficialRelation == "" {
		log.Error("Missing required parameter: beneficial_relation")
		return lib.CustomError(http.StatusBadRequest, "beneficial_relation can not be blank", "beneficial_relation can not be blank")
	} else {
		n, err := strconv.ParseUint(beneficialRelation, 10, 64)
		if err == nil && n > 0 {
			paramsOaPersonalData["beneficial_relation"] = beneficialRelation
		} else {
			log.Error("Wrong input for parameter: beneficial_relation")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: beneficial_relation", "Wrong input for parameter: beneficial_relation")
		}
	}

	beneficialName := c.FormValue("beneficial_name")
	if beneficialName == "" {
		log.Error("Missing required parameter: beneficial_name")
		return lib.CustomError(http.StatusBadRequest, "beneficial_name can not be blank", "beneficial_name can not be blank")
	}

	//BANK DETAIL
	bankKey := c.FormValue("bank_key")
	if bankKey == "" {
		log.Error("Missing required parameter: bank_key")
		return lib.CustomError(http.StatusBadRequest, "bank_key can not be blank", "bank_key can not be blank")
	} else {
		bank, err := strconv.ParseUint(bankKey, 10, 64)
		if err != nil || bank == 0 {
			log.Error("Wrong input for parameter: bank_key")
			return lib.CustomError(http.StatusBadRequest)
		}
	}

	accountNo := c.FormValue("account_no")
	if accountNo == "" {
		log.Error("Missing required parameter: account_no")
		return lib.CustomError(http.StatusBadRequest, "account_no can not be blank", "account_no can not be blank")
	}

	accountName := c.FormValue("account_name")
	if accountName == "" {
		log.Error("Missing required parameter: account_name")
		return lib.CustomError(http.StatusBadRequest, "account_name can not be blank", "account_name can not be blank")
	}

	branchName := c.FormValue("branch_name")
	if branchName == "" {
		log.Error("Missing required parameter: branch_name")
		return lib.CustomError(http.StatusBadRequest, "branch_name can not be blank", "branch_name can not be blank")
	}

	flagPriority := c.FormValue("flag_priority")
	if flagPriority == "" {
		log.Error("Missing required parameter: flag_priority")
		return lib.CustomError(http.StatusBadRequest, "flag_priority can not be blank", "flag_priority can not be blank")
	}

	//bank account 2
	bankKey2 := c.FormValue("bank_key_2")
	accountNo2 := c.FormValue("account_no_2")
	accountName2 := c.FormValue("account_name_2")
	branchName2 := c.FormValue("branch_name_2")
	flagPriority2 := c.FormValue("flag_priority_2")

	// if bankKey2 != "" || accountNo2 != "" || accountName2 != "" || branchName2 != "" || flagPriority2 != "" {
	// 	if bankKey2 == "" {
	// 		log.Error("Missing required parameter: bank_key_2")
	// 		return lib.CustomError(http.StatusBadRequest, "bank_key_2 can not be blank", "bank_key_2 can not be blank")
	// 	}

	// 	if accountNo2 == "" {
	// 		log.Error("Missing required parameter: account_no_2")
	// 		return lib.CustomError(http.StatusBadRequest, "account_no_2 can not be blank", "account_no_2 can not be blank")
	// 	}

	// 	if accountName2 == "" {
	// 		log.Error("Missing required parameter: account_name_2")
	// 		return lib.CustomError(http.StatusBadRequest, "account_name_2 can not be blank", "account_name_2 can not be blank")
	// 	}

	// 	if branchName2 == "" {
	// 		log.Error("Missing required parameter: branch_name_2")
	// 		return lib.CustomError(http.StatusBadRequest, "branch_name_2 can not be blank", "branch_name_2 can not be blank")
	// 	}

	// 	if flagPriority2 == "" {
	// 		log.Error("Missing required parameter: flag_priority_2")
	// 		return lib.CustomError(http.StatusBadRequest, "flag_priority_2 can not be blank", "flag_priority_2 can not be blank")
	// 	}
	// }

	//bank account 3
	bankKey3 := c.FormValue("bank_key_3")
	accountNo3 := c.FormValue("account_no_3")
	accountName3 := c.FormValue("account_name_3")
	branchName3 := c.FormValue("branch_name_3")
	flagPriority3 := c.FormValue("flag_priority_3")

	// if bankKey3 != "" || accountNo3 != "" || accountName3 != "" || branchName3 != "" || flagPriority3 != "" {
	// 	if bankKey3 == "" {
	// 		log.Error("Missing required parameter: bank_key_3")
	// 		return lib.CustomError(http.StatusBadRequest, "bank_key_3 can not be blank", "bank_key_3 can not be blank")
	// 	}

	// 	if accountNo3 == "" {
	// 		log.Error("Missing required parameter: account_no_3")
	// 		return lib.CustomError(http.StatusBadRequest, "account_no_3 can not be blank", "account_no_3 can not be blank")
	// 	}

	// 	if accountName3 == "" {
	// 		log.Error("Missing required parameter: account_name_3")
	// 		return lib.CustomError(http.StatusBadRequest, "account_name_3 can not be blank", "account_name_3 can not be blank")
	// 	}

	// 	if branchName3 == "" {
	// 		log.Error("Missing required parameter: branch_name_3")
	// 		return lib.CustomError(http.StatusBadRequest, "branch_name_3 can not be blank", "branch_name_3 can not be blank")
	// 	}

	// 	if flagPriority3 == "" {
	// 		log.Error("Missing required parameter: flag_priority_3")
	// 		return lib.CustomError(http.StatusBadRequest, "flag_priority_3 can not be blank", "flag_priority_3 can not be blank")
	// 	}
	// }

	countPriority := 0
	if flagPriority == "1" {
		countPriority++
	}
	if flagPriority2 == "1" {
		countPriority++
	}
	if flagPriority3 == "1" {
		countPriority++
	}

	if countPriority == 0 {
		log.Error("Missing required parameter: flag_prioriti must select 1")
		return lib.CustomError(http.StatusBadRequest, "flag_prioriti must select 1", "flag_prioriti must select 1")
	}

	if countPriority > 1 {
		log.Error("Missing required parameter: flag_prioriti must select 1")
		return lib.CustomError(http.StatusBadRequest, "flag_prioriti must select 1", "flag_prioriti must select 1")
	}

	//QUIZ
	quizOption := c.FormValue("quiz_option")
	if quizOption == "" {
		log.Error("Missing required parameter: quiz_option")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: quiz_option", "Missing required parameter: quiz_option")
	}

	s := strings.Split(quizOption, ",")
	var quizoptionkey []string

	for _, value := range s {
		is := strings.TrimSpace(value)
		if is != "" {
			if _, ok := lib.Find(quizoptionkey, is); !ok {
				quizoptionkey = append(quizoptionkey, is)
			}
		}
	}
	if len(quizoptionkey) <= 0 {
		log.Error("Missing required parameter: quiz_option")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: quiz_option", "Missing required parameter: quiz_option")
	}

	date := time.Now().AddDate(0, 0, 1)
	dateLayout := "2006-01-02 15:04:05"

	//OA_REQUEST
	layout := "2006-01-02 15:04:05"
	dateNow := time.Now().Format(layout)
	paramsOaRequest := make(map[string]string)
	paramsOaRequest["oa_status"] = "258"
	paramsOaRequest["oa_entry_start"] = dateNow
	paramsOaRequest["oa_entry_end"] = dateNow
	paramsOaRequest["oa_request_type"] = oaRequestType
	if salesCode != "" {
		paramsOaRequest["sales_code"] = salesCode
	}
	paramsOaRequest["branch_key"] = branchkey
	paramsOaRequest["agent_key"] = agentkey
	paramsOaRequest["customer_key"] = customerKey
	paramsOaRequest["rec_status"] = "1"

	//OA_POSTAL_ADDRESS CARD ID
	addressIDParams := make(map[string]string)
	addressIDParams["address_type"] = "17"
	addressIDParams["address_line1"] = addressid
	addressIDParams["kabupaten_key"] = kabupatenid
	addressIDParams["kecamatan_key"] = kecamatanid
	addressIDParams["postal_code"] = postalid
	addressIDParams["rec_status"] = "1"
	addressIDParams["rec_created_date"] = time.Now().Format(dateLayout)
	addressIDParams["rec_created_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

	//OA_POSTAL_ADDRESS DOMICILE
	addressDomicileParams := make(map[string]string)
	addressDomicileParams["address_type"] = "18"
	addressDomicileParams["address_line1"] = addressdomicile
	addressDomicileParams["kabupaten_key"] = kabupatendomicile
	addressDomicileParams["kecamatan_key"] = kecamatandomicile
	addressDomicileParams["postal_code"] = postaldomicile
	addressDomicileParams["rec_created_date"] = time.Now().Format(dateLayout)
	addressDomicileParams["rec_created_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
	addressDomicileParams["rec_status"] = "1"

	//OA_PERSONAL_DATA
	log.Info("dateBirth: " + dateBirth)
	dateBirth += " 00:00:00"
	date, err = time.Parse(layout, dateBirth)
	dateStr := date.Format(layout)
	log.Info("dateBirth: " + dateStr)

	paramsOaPersonalData["full_name"] = fullname
	if nationality == "97" { //indonesia
		paramsOaPersonalData["idcard_type"] = "12"
	} else {
		paramsOaPersonalData["idcard_type"] = "13"
	}
	paramsOaPersonalData["place_birth"] = placeBirth
	paramsOaPersonalData["date_birth"] = dateBirth
	paramsOaPersonalData["nationality"] = nationality
	paramsOaPersonalData["idcard_no"] = idcardNumber
	paramsOaPersonalData["gender"] = gender
	paramsOaPersonalData["marital_status"] = maritalStatus
	paramsOaPersonalData["phone_home"] = phoneHome
	paramsOaPersonalData["phone_mobile"] = phone
	paramsOaPersonalData["email_address"] = email
	paramsOaPersonalData["religion"] = religion
	paramsOaPersonalData["education"] = education
	paramsOaPersonalData["occup_job"] = job
	paramsOaPersonalData["occup_company"] = company
	paramsOaPersonalData["occup_position"] = position
	paramsOaPersonalData["beneficial_full_name"] = beneficialName
	paramsOaPersonalData["emergency_phone_no"] = emergencyPhone
	paramsOaPersonalData["relation_full_name"] = relationName
	paramsOaPersonalData["mother_maiden_name"] = motherMaidenName
	paramsOaPersonalData["emergency_full_name"] = emergencyName
	paramsOaPersonalData["rec_status"] = "1"
	paramsOaPersonalData["rec_created_date"] = time.Now().Format(dateLayout)
	paramsOaPersonalData["rec_created_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

	tx, _ := db.Db.Begin()

	//CEK SC_USER_LOGIN CUSTOMER
	var idUserLogin string
	var scUserLogin models.ScUserLogin
	status, err = models.GetScUserLoginByCustomerKey(&scUserLogin, customerKey)
	if err != nil {
		tx.Rollback()
		log.Error(err.Error())
		return lib.CustomError(http.StatusBadRequest, err.Error(), "Failed create data")
	} else {
		idUserLogin = strconv.FormatUint(scUserLogin.UserLoginKey, 10)
	}

	paramsOaRequest["user_login_key"] = idUserLogin
	paramsOaRequest["rec_created_date"] = time.Now().Format(dateLayout)
	paramsOaRequest["rec_created_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
	paramsOaRequest["rec_attribute_id3"] = c.Request().UserAgent()

	//SAVE AO_PORTAL_ADDRESS IDCARD
	status, err, addressidID := models.CreateOaPostalAddress(addressIDParams)
	if err != nil {
		tx.Rollback()
		log.Error("Failed create adrress data idcard: " + err.Error())
		return lib.CustomError(status, err.Error(), "failed input data")
	}
	addressID, err := strconv.ParseUint(addressidID, 10, 64)
	if addressID == 0 {
		tx.Rollback()
		log.Error("Failed create adrress data idcard")
		return lib.CustomError(http.StatusBadGateway, "failed input data", "failed input data")
	}
	paramsOaPersonalData["idcard_address_key"] = addressidID

	//SAVE AO_PORTAL_ADDRESS DOMICILE
	status, err, addressDomicileID := models.CreateOaPostalAddress(addressDomicileParams)
	if err != nil {
		tx.Rollback()
		log.Error("Failed create adrress data domicile: " + err.Error())
		return lib.CustomError(status, err.Error(), "failed input data")
	}
	addressID, err = strconv.ParseUint(addressDomicileID, 10, 64)
	if addressID == 0 {
		tx.Rollback()
		log.Error("Failed create adrress data domicile")
		return lib.CustomError(http.StatusBadGateway, "failed input data", "failed input data")
	}
	paramsOaPersonalData["domicile_address_key"] = addressDomicileID

	//SAVE AO_PORTAL_ADDRESS COMPANY
	addressCompanyParams := make(map[string]string)
	companyAddress := c.FormValue("company_address")
	if companyAddress != "" {
		addressCompanyParams["address_type"] = "19"
		addressCompanyParams["address_line1"] = companyAddress
		addressCompanyParams["rec_status"] = "1"
		addressCompanyParams["rec_created_date"] = time.Now().Format(dateLayout)
		addressCompanyParams["rec_created_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

		status, err, addressCompanyID := models.CreateOaPostalAddress(addressCompanyParams)
		if err != nil {
			tx.Rollback()
			log.Error("Failed create adrress data company: " + err.Error())
			return lib.CustomError(status, err.Error(), "failed input data")
		}
		addressID, err = strconv.ParseUint(addressCompanyID, 10, 64)
		if addressID == 0 {
			tx.Rollback()
			log.Error("Failed create adrress data company")
			return lib.CustomError(http.StatusBadGateway, "failed input data", "failed input data")
		}
		paramsOaPersonalData["occup_address_key"] = addressCompanyID
	}

	//SAVE OA_REQUEST
	status, err, requestID := models.CreateOaRequest(paramsOaRequest)
	if err != nil {
		tx.Rollback()
		log.Error("Failed create oa request data: " + err.Error())
		return lib.CustomError(status, err.Error(), "failed input data")
	}
	request, err := strconv.ParseUint(requestID, 10, 64)
	if request == 0 {
		tx.Rollback()
		log.Error("Failed create oa request data")
		return lib.CustomError(http.StatusBadGateway, "failed input data", "failed input data")
	}
	paramsOaPersonalData["oa_request_key"] = requestID

	//SAVE MS_BANK_ACCOUNT
	var bankAccountIDPriority string

	//MS_BANK_ACCOUNT 1
	paramsBank := make(map[string]string)
	paramsBank["bank_key"] = bankKey
	paramsBank["account_no"] = accountNo
	paramsBank["account_holder_name"] = accountName
	paramsBank["branch_name"] = branchName
	paramsBank["currency_key"] = "1"
	paramsBank["bank_account_type"] = "1"
	paramsBank["rec_domain"] = "1"
	paramsBank["rec_status"] = "1"
	paramsBank["rec_created_date"] = time.Now().Format(dateLayout)
	paramsBank["rec_created_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
	status, err, bankAccountID := models.CreateMsBankAccount(paramsBank)
	if err != nil {
		tx.Rollback()
		log.Error("Failed create bank account data: " + err.Error())
		return lib.CustomError(status, err.Error(), "failed input data")
	}
	accountID, err := strconv.ParseUint(bankAccountID, 10, 64)
	if accountID == 0 {
		tx.Rollback()
		log.Error("Failed create bank account data")
		return lib.CustomError(http.StatusBadGateway, "failed input data", "failed input data")
	}
	if flagPriority == "1" {
		bankAccountIDPriority = bankAccountID
	}

	//OA_REQUEST_BANK_ACCOUNT 1
	bankReq := make(map[string]string)
	bankReq["oa_request_key"] = requestID
	bankReq["bank_account_key"] = bankAccountID
	bankReq["flag_priority"] = flagPriority
	bankReq["bank_account_name"] = accountName
	bankReq["rec_status"] = "1"
	bankReq["rec_created_date"] = time.Now().Format(dateLayout)
	bankReq["rec_created_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
	status, err, _ = models.CreateOaRequestBankAccount(bankReq)
	if err != nil {
		tx.Rollback()
		log.Error("Failed create OA_REQUEST_BANK_ACCOUNT 1: " + err.Error())
		return lib.CustomError(status, err.Error(), "failed input data")
	}

	//MS_BANK_ACCOUNT 2
	if bankKey2 != "" || accountNo2 != "" || accountName2 != "" || branchName2 != "" || flagPriority2 != "" {
		//MS_BANK_ACCOUNT 2
		paramsBank2 := make(map[string]string)
		paramsBank2["bank_key"] = bankKey2
		paramsBank2["account_no"] = accountNo2
		paramsBank2["account_holder_name"] = accountName2
		paramsBank2["branch_name"] = branchName2
		paramsBank2["currency_key"] = "1"
		paramsBank2["bank_account_type"] = "1"
		paramsBank2["rec_domain"] = "1"
		paramsBank2["rec_status"] = "1"
		paramsBank2["rec_created_date"] = time.Now().Format(dateLayout)
		paramsBank2["rec_created_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
		status, err, bankAccountID2 := models.CreateMsBankAccount(paramsBank2)
		if err != nil {
			tx.Rollback()
			log.Error("Failed create bank account data: " + err.Error())
			return lib.CustomError(status, err.Error(), "failed input data")
		}
		accountID, err := strconv.ParseUint(bankAccountID, 10, 64)
		if accountID == 0 {
			tx.Rollback()
			log.Error("Failed create bank account data")
			return lib.CustomError(http.StatusBadGateway, "failed input data", "failed input data")
		}
		if flagPriority2 == "1" {
			bankAccountIDPriority = bankAccountID2
		}

		//OA_REQUEST_BANK_ACCOUNT 2
		bankReq := make(map[string]string)
		bankReq["oa_request_key"] = requestID
		bankReq["bank_account_key"] = bankAccountID2
		bankReq["flag_priority"] = flagPriority2
		bankReq["bank_account_name"] = accountName2
		bankReq["rec_status"] = "1"
		bankReq["rec_created_date"] = time.Now().Format(dateLayout)
		bankReq["rec_created_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
		status, err, _ = models.CreateOaRequestBankAccount(bankReq)
		if err != nil {
			tx.Rollback()
			log.Error("Failed create OA_REQUEST_BANK_ACCOUNT 2: " + err.Error())
			return lib.CustomError(status, err.Error(), "failed input data")
		}
	}

	//MS_BANK_ACCOUNT 3
	if bankKey3 != "" || accountNo3 != "" || accountName3 != "" || branchName3 != "" || flagPriority3 != "" {
		//MS_BANK_ACCOUNT 3
		paramsBank3 := make(map[string]string)
		paramsBank3["bank_key"] = bankKey3
		paramsBank3["account_no"] = accountNo3
		paramsBank3["account_holder_name"] = accountName3
		paramsBank3["branch_name"] = branchName3
		paramsBank3["currency_key"] = "1"
		paramsBank3["bank_account_type"] = "1"
		paramsBank3["rec_domain"] = "1"
		paramsBank3["rec_status"] = "1"
		paramsBank3["rec_created_date"] = time.Now().Format(dateLayout)
		paramsBank3["rec_created_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
		status, err, bankAccountID3 := models.CreateMsBankAccount(paramsBank3)
		if err != nil {
			tx.Rollback()
			log.Error("Failed create bank account data: " + err.Error())
			return lib.CustomError(status, err.Error(), "failed input data")
		}
		accountID, err := strconv.ParseUint(bankAccountID, 10, 64)
		if accountID == 0 {
			tx.Rollback()
			log.Error("Failed create bank account data")
			return lib.CustomError(http.StatusBadGateway, "failed input data", "failed input data")
		}
		if flagPriority3 == "1" {
			bankAccountIDPriority = bankAccountID3
		}

		//OA_REQUEST_BANK_ACCOUNT 3
		bankReq := make(map[string]string)
		bankReq["oa_request_key"] = requestID
		bankReq["bank_account_key"] = bankAccountID3
		bankReq["flag_priority"] = flagPriority3
		bankReq["bank_account_name"] = accountName3
		bankReq["rec_status"] = "1"
		bankReq["rec_created_date"] = time.Now().Format(dateLayout)
		bankReq["rec_created_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
		status, err, _ = models.CreateOaRequestBankAccount(bankReq)
		if err != nil {
			tx.Rollback()
			log.Error("Failed create OA_REQUEST_BANK_ACCOUNT 3: " + err.Error())
			return lib.CustomError(status, err.Error(), "failed input data")
		}
	}

	paramsOaPersonalData["bank_account_key"] = bankAccountIDPriority

	//SAVE OA_PERSONAL_DATA
	err = os.MkdirAll(config.BasePathImage+"/images/user/"+idUserLogin, 0755)
	if err != nil {
		log.Error(err.Error())
	} else {
		var file *multipart.FileHeader
		file, err = c.FormFile("pic_ktp")
		if file != nil {
			if err != nil {
				return lib.CustomError(http.StatusBadRequest, err.Error(), "Missing required parameter: pic_ktp")
			}
			// Get file extension
			extension := filepath.Ext(file.Filename)
			// Generate filename
			var filename string
			for {
				filename = lib.RandStringBytesMaskImprSrc(20)
				log.Println("Generate filename:", filename)
				var personalData []models.OaPersonalData
				getParams := make(map[string]string)
				getParams["pic_ktp"] = filename + extension
				_, err := models.GetAllOaPersonalData(&personalData, 1, 0, getParams, false)
				if (err == nil && len(personalData) < 1) || err != nil {
					break
				}
			}
			// Upload image and move to proper directory
			err = lib.UploadImage(file, config.BasePathImage+"/images/user/"+idUserLogin+"/"+filename+extension)
			if err != nil {
				log.Println(err)
				return lib.CustomError(http.StatusInternalServerError)
			}
			paramsOaPersonalData["pic_ktp"] = filename + extension
		} else {
			paramsOaPersonalData["pic_ktp"] = picKtpDefault
		}

		file, err = c.FormFile("pic_selfie_ktp")
		if file != nil {
			if err != nil {
				return lib.CustomError(http.StatusBadRequest, err.Error(), "Missing required parameter: pic_selfie_ktp")
			}
			// Get file extension
			extension := filepath.Ext(file.Filename)
			// Generate filename
			var filename string
			for {
				filename = lib.RandStringBytesMaskImprSrc(20)
				log.Println("Generate filename:", filename)
				var personalData []models.OaPersonalData
				getParams := make(map[string]string)
				getParams["pic_selfie_ktp"] = filename + extension
				_, err := models.GetAllOaPersonalData(&personalData, 1, 0, getParams, false)
				if (err == nil && len(personalData) < 1) || err != nil {
					break
				}
			}
			// Upload image and move to proper directory
			err = lib.UploadImage(file, config.BasePathImage+"/images/user/"+idUserLogin+"/"+filename+extension)
			if err != nil {
				log.Println(err)
				return lib.CustomError(http.StatusInternalServerError)
			}
			paramsOaPersonalData["pic_selfie_ktp"] = filename + extension
		} else {
			paramsOaPersonalData["pic_selfie_ktp"] = picSelfieKtpDefault
		}
	}

	err = os.MkdirAll(config.BasePathImage+"/images/user/"+idUserLogin+"/signature", 0755)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(http.StatusBadGateway, err.Error(), err.Error())
	}
	file, err = c.FormFile("signature")
	if file != nil {
		if err != nil {
			return lib.CustomError(http.StatusBadRequest, err.Error(), "Missing required parameter: signature")
		}
		// Get file extension
		extension := filepath.Ext(file.Filename)
		// Generate filename
		var filename string
		for {
			filename = lib.RandStringBytesMaskImprSrc(20)
			log.Println("Generate filename:", filename)
			var personalData []models.OaPersonalData
			getParams := make(map[string]string)
			getParams["rec_image1"] = filename + extension
			_, err := models.GetAllOaPersonalData(&personalData, 1, 0, getParams, false)
			if (err == nil && len(personalData) < 1) || err != nil {
				break
			}
		}
		// Upload image and move to proper directory
		err = lib.UploadImage(file, config.BasePathImage+"/images/user/"+idUserLogin+"/signature/"+filename+extension)
		if err != nil {
			log.Println(err)
			return lib.CustomError(http.StatusInternalServerError)
		}
		paramsOaPersonalData["rec_image1"] = filename + extension
	} else {
		paramsOaPersonalData["rec_image1"] = signatureDefault
	}

	status, err, _ = models.CreateOaPersonalData(paramsOaPersonalData)
	if err != nil {
		tx.Rollback()
		log.Error("Failed create personal data: " + err.Error())
		return lib.CustomError(status, err.Error(), "failed input data")
	}

	//SAVE CMS_QUIZ_OPTIONS
	var questionOptions []models.QuestionOptionQuiz
	status, err = models.AdminGetQuestionOptionQuiz(&questionOptions, quizoptionkey)

	if err != nil {
		tx.Rollback()
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if len(questionOptions) < 1 {
		tx.Rollback()
		log.Error("Missing required parameter: quiz_option")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: quiz_option", "Missing required parameter: quiz_option")
	}

	var bindVar []interface{}
	var score uint64 = 0
	for _, val := range questionOptions {

		var row []string
		row = append(row, requestID)
		row = append(row, val.QuizQuestionKey)
		row = append(row, val.QuizOptionKey)
		row = append(row, strconv.FormatUint(val.QuizOptionScore, 10))
		row = append(row, "1")
		score += val.QuizOptionScore
		bindVar = append(bindVar, row)
	}

	var riskProfile models.MsRiskProfile
	scoreStr := strconv.FormatUint(score, 10)
	status, err = models.GetMsRiskProfileScore(&riskProfile, scoreStr)
	if err != nil {
		tx.Rollback()
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data risk profile")
	}

	paramsOaRiskProfile := make(map[string]string)
	paramsOaRiskProfile["oa_request_key"] = requestID
	paramsOaRiskProfile["risk_profile_key"] = strconv.FormatUint(riskProfile.RiskProfileKey, 10)
	paramsOaRiskProfile["score_result"] = scoreStr
	paramsOaRiskProfile["rec_status"] = "1"
	paramsOaRiskProfile["rec_created_date"] = time.Now().Format(dateLayout)
	paramsOaRiskProfile["rec_created_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

	var bindInterface []interface{}
	for i := 0; i < len(bindVarz); i++ {
		bindVarz[i][1] = requestID
		bindInterface = append(bindInterface, bindVarz[i])
	}

	if len(bindInterface) > 0 {
		status, err = models.CreateMultipleUdfValue(bindInterface)
		if err != nil {
			log.Error(err.Error())
		}
	}

	status, err = models.CreateOaRiskProfile(paramsOaRiskProfile)
	if err != nil {
		tx.Rollback()
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed input data")
	}

	status, err = models.CreateMultipleOaRiskProfileQuiz(bindVar)
	if err != nil {
		tx.Rollback()
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed input data")
	}

	tx.Commit()

	// Send email

	if config.Envi != "DEV" {
		t := template.New("index-pengkinian.html")

		t, err = t.ParseFiles(config.BasePath + "/mail/index-pengkinian.html")
		if err != nil {
			log.Println(err)
		}

		var tpl bytes.Buffer
		if err := t.Execute(&tpl, struct {
			Name    string
			FileUrl string
		}{Name: fullname, FileUrl: config.FileUrl + "/images/mail"}); err != nil {
			log.Println(err)
		}

		result := tpl.String()

		mailer := gomail.NewMessage()
		mailer.SetHeader("From", config.EmailFrom)
		mailer.SetHeader("To", scUserLogin.UloginEmail)
		mailer.SetHeader("Subject", "[MotionFunds] Pengkinian Data Kamu sedang Diproses")
		mailer.SetBody("text/html", result)

		err = lib.SendEmail(mailer)
		if err != nil {
			log.Error("Error send email")
			log.Error(err)
			log.Error("Error send email")
		} else {
			log.Info("Email sent")
		}
		// dialer := gomail.NewDialer(
		// 	config.EmailSMTPHost,
		// 	int(config.EmailSMTPPort),
		// 	config.EmailFrom,
		// 	config.EmailFromPassword,
		// )
		// dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

		// err = dialer.DialAndSend(mailer)
		// if err != nil {
		// 	log.Error("Error send email")
		// 	log.Error(err)
		// 	log.Error("Error send email")
		// }
		// log.Info("Email sent")

		//sent email to all CS & Agent
		var oaRequest models.OaRequest
		status, err = models.GetOaRequest(&oaRequest, requestID)
		if err != nil {
			log.Error(err.Error())
		} else {
			var personalData models.OaPersonalData
			status, err = models.GetOaPersonalData(&personalData, requestID, "oa_request_key")
			if err != nil {
				log.Error(err.Error())
			} else {
				SentEmailOaPengkinianToBackOffice(oaRequest, personalData, "11")
				SentEmailOaPengkinianToSales(oaRequest, personalData)
			}
		}
	}

	//insert message notif in app
	strIDUserLogin := strconv.FormatUint(lib.Profile.UserID, 10)
	paramsUserMessage := make(map[string]string)
	paramsUserMessage["umessage_type"] = "245"
	paramsUserMessage["umessage_recipient_key"] = idUserLogin
	paramsUserMessage["umessage_receipt_date"] = time.Now().Format(dateLayout)
	paramsUserMessage["flag_read"] = "0"
	paramsUserMessage["umessage_sent_date"] = time.Now().Format(dateLayout)
	paramsUserMessage["flag_sent"] = "1"
	subject := "Pengkinian Data sedang Diproses"
	body := "Terima kasih telah melakukan Pengkinian Data. Kami sedang melakukan proses verifikasi data kamu max. 1X24 jam. Mohon ditunggu ya."
	paramsUserMessage["umessage_subject"] = subject
	paramsUserMessage["umessage_body"] = body

	paramsUserMessage["umessage_category"] = "248"
	paramsUserMessage["flag_archieved"] = "0"
	paramsUserMessage["archieved_date"] = time.Now().Format(dateLayout)
	paramsUserMessage["rec_status"] = "1"
	paramsUserMessage["rec_created_date"] = time.Now().Format(dateLayout)
	paramsUserMessage["rec_created_by"] = strIDUserLogin

	status, err = models.CreateScUserMessage(paramsUserMessage)
	if err != nil {
		log.Error("Error create user message")
	} else {
		log.Error("Sukses insert user message")
	}
	lib.CreateNotifCustomerFromAdminByUserLoginKey(idUserLogin, subject, body, "HOME")

	var responseData models.MsRiskProfileInfo

	responseData.RiskProfileKey = riskProfile.RiskProfileKey
	responseData.RiskCode = riskProfile.RiskCode
	responseData.RiskName = riskProfile.RiskName
	responseData.RiskDesc = riskProfile.RiskDesc
	responseData.Score = score

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = responseData

	return c.JSON(http.StatusOK, response)
}

func CheckUniqueEmailNoHp(c echo.Context) error {
	field := c.FormValue("field")
	if field == "" {
		log.Error("Wrong input for parameter: field")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: field", "Missing required parameter: field")
	}
	value := c.FormValue("value")
	if value == "" {
		log.Error("Wrong input for parameter: value")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: value", "Missing required parameter: value")
	}

	var userKey *string
	userLoginKey := c.FormValue("user_login_key")
	if userLoginKey != "" {
		userKey = &userLoginKey
	}

	var countData models.CountData
	status, err := models.ValidateUniqueData(&countData, field, value, userKey)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}

	var ff string

	if field == "email" {
		ff = "Email"
	} else {
		ff = "Mobile Number"
	}

	var valid bool
	var message string
	if int(countData.CountData) > int(0) {
		valid = false
		message = ff + " sudah digunakan"
	} else {
		valid = true
		message = ff + " valid"
	}
	responseData := make(map[string]interface{})
	responseData["valid"] = valid
	responseData["message"] = message

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = responseData

	return c.JSON(http.StatusOK, response)
}

func CheckUniqueNoId(c echo.Context) error {
	noId := c.FormValue("no_id")
	if noId == "" {
		log.Error("Wrong input for parameter: no_id")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: no_id", "Missing required parameter: no_id")
	}

	var cusKey *string
	customerKey := c.FormValue("customer_key")
	if customerKey != "" {
		cusKey = &customerKey
	}

	var countData models.CountData
	status, err := models.ValidateUniquePersonalData(&countData, "b.idcard_no", noId, cusKey, nil)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	var valid bool
	var message string
	if int(countData.CountData) > int(0) {
		valid = false
		message = "ID Card sudah digunakan"
	} else {
		valid = true
		message = "ID Card valid"
	}
	responseData := make(map[string]interface{})
	responseData["valid"] = valid
	responseData["message"] = message

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = responseData

	return c.JSON(http.StatusOK, response)
}

func IndividuSendAccountStatement(c echo.Context) error {
	var err error
	decimal.MarshalJSONWithoutQuotes = true

	customerKey := c.QueryParam("customer_key")
	if customerKey == "" {
		log.Error("Missing required parameter: customer_key")
		return lib.CustomError(http.StatusBadRequest, "customer_key can not be blank", "customer_key can not be blank")
	} else {
		n, err := strconv.ParseUint(customerKey, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: customer_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: customer_key", "Wrong input for parameter: customer_key")
		}
	}

	var customer models.HeaderCustomerDetailAccountStatement
	_, err = models.GetHeaderCustomerDetailAccountStatement(&customer, customerKey)
	if err != nil {
		log.Error("Customer not found")
		return lib.CustomError(http.StatusBadRequest, "Customer not found", "Customer not found")
	}

	var datefrom string
	var dateto string

	datefrom = c.QueryParam("date_from")
	if datefrom == "" {
		log.Error("Missing required parameter: date_from")
		return lib.CustomError(http.StatusBadRequest, "date_from can not be blank", "date_from can not be blank")
	}

	dateto = c.QueryParam("date_to")
	if dateto == "" {
		log.Error("Missing required parameter: date_to")
		return lib.CustomError(http.StatusBadRequest, "date_to can not be blank", "date_to can not be blank")
	}
	layoutISO := "2006-01-02"

	from, _ := time.Parse(layoutISO, datefrom)
	from = time.Date(from.Year(), from.Month(), from.Day(), 0, 0, 0, 0, time.UTC)

	to, _ := time.Parse(layoutISO, dateto)
	to = time.Date(to.Year(), to.Month(), to.Day(), 0, 0, 0, 0, time.UTC)

	var dateawal string
	var dateakhir string

	if from.Before(to) {
		dateawal = datefrom
		dateakhir = dateto
	}

	if from.After(to) {
		dateawal = dateto
		dateakhir = datefrom
	}

	responseData := make(map[string]interface{})

	//get header
	layout := "2006-01-02"
	newLayout := "02 Jan 2006"
	header := make(map[string]interface{})
	dateParem, _ := time.Parse(layout, dateawal)
	header["date_from"] = dateParem.Format(newLayout)
	dateParem, _ = time.Parse(layout, dateakhir)
	header["date_to"] = dateParem.Format(newLayout)
	header["customer_key"] = customer.CustomerKey
	header["cif"] = customer.Cif
	header["sid"] = customer.Sid
	header["full_name"] = customer.FullName
	header["email"] = customer.Email
	header["address"] = customer.Address

	responseData["header"] = header

	var transactions []models.AccountStatementCustomerProduct

	status, err := models.AdminGetAllAccountStatementCustomerProduct(&transactions, customerKey, dateawal, dateakhir)

	if err != nil {
		if err != sql.ErrNoRows {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data transaction")
		}
	}

	if len(transactions) > 0 {
		var datatrans []interface{}
		var productKey uint64
		var transGroupProduct []interface{}
		product := make(map[string]interface{})
		count := make(map[string]interface{})
		var totalSubs decimal.Decimal
		var totalRedm decimal.Decimal
		nol := decimal.NewFromInt(0)

		var accKeyLast uint64
		var productKeyLast uint64

		for idx, tr := range transactions {
			if idx != 0 {
				if productKey != tr.ProductKey {
					var balanceEnding models.BeginningEndingBalance
					_, err = models.GetBeginningEndingBalanceAcc(&balanceEnding, "ENDING BALANCE", dateakhir, strconv.FormatUint(accKeyLast, 10), strconv.FormatUint(productKeyLast, 10))

					endingbalance := make(map[string]interface{})
					if err != nil {
						var lastavgnav models.NavValue
						_, err = models.AdminLastAvgNav(&lastavgnav, strconv.FormatUint(productKeyLast, 10), customerKey, dateakhir)
						if err != nil {
							endingbalance["avg_nav"] = nol
						} else {
							endingbalance["avg_nav"] = lastavgnav.NavValue.Truncate(2)
						}

						var lastnav models.NavValue
						_, err = models.AdminLastNavValue(&lastnav, strconv.FormatUint(productKeyLast, 10), dateakhir)
						if err != nil {
							endingbalance["nav_value"] = nol
						} else {
							endingbalance["nav_value"] = lastnav.NavValue.Truncate(2)
						}

						dateParem, _ = time.Parse(layout, dateakhir)
						endingbalance["date"] = dateParem.Format(newLayout)
						endingbalance["description"] = "ENDING BALANCE"
						endingbalance["amount"] = nol
						endingbalance["unit"] = nol
						endingbalance["fee"] = nol
						transGroupProduct = append(transGroupProduct, endingbalance)
					} else {
						endingbalance["date"] = balanceEnding.Tanggal
						endingbalance["description"] = balanceEnding.Description
						endingbalance["amount"] = balanceEnding.Amount.Truncate(*balanceEnding.DecAmount)
						endingbalance["nav_value"] = balanceEnding.NavValue.Truncate(*balanceEnding.DecNav)
						endingbalance["unit"] = balanceEnding.Unit.Truncate(*balanceEnding.DecUnit)
						endingbalance["avg_nav"] = balanceEnding.AvgNav.Truncate(*balanceEnding.DecNav)
						endingbalance["fee"] = balanceEnding.Fee.Truncate(0)
						transGroupProduct = append(transGroupProduct, endingbalance)
					}

					row := make(map[string]interface{})
					row["count"] = count
					row["product"] = product
					row["list"] = transGroupProduct
					datatrans = append(datatrans, row)
					productKey = tr.ProductKey

					transGroupProduct = nil
					product = make(map[string]interface{})
					count = make(map[string]interface{})
					totalSubs = nol
					totalRedm = nol

					var balance models.BeginningEndingBalance
					_, err = models.GetBeginningEndingBalanceAcc(&balance, "BEGINNING BALANCE", dateawal, strconv.FormatUint(tr.AccKey, 10), strconv.FormatUint(tr.ProductKey, 10))

					beginning := make(map[string]interface{})
					if err != nil {
						var lastavgnav models.NavValue
						_, err = models.AdminLastAvgNav(&lastavgnav, strconv.FormatUint(tr.ProductKey, 10), customerKey, dateawal)
						if err != nil {
							beginning["avg_nav"] = nol
						} else {
							beginning["avg_nav"] = lastavgnav.NavValue.Truncate(2)
						}

						var lastnav models.NavValue
						_, err = models.AdminLastNavValue(&lastnav, strconv.FormatUint(tr.ProductKey, 10), dateawal)
						if err != nil {
							beginning["nav_value"] = nol
						} else {
							beginning["nav_value"] = lastnav.NavValue.Truncate(2)
						}
						dateParem, _ = time.Parse(layout, dateawal)
						beginning["date"] = dateParem.Format(newLayout)
						beginning["description"] = "BEGINNING BALANCE"
						beginning["amount"] = nol
						beginning["unit"] = nol
						beginning["fee"] = nol
						transGroupProduct = append(transGroupProduct, beginning)
					} else {
						beginning["date"] = balance.Tanggal
						beginning["description"] = balance.Description
						beginning["amount"] = balance.Amount.Truncate(*balance.DecAmount)
						beginning["nav_value"] = balance.NavValue.Truncate(*balance.DecNav)
						beginning["unit"] = balance.Unit.Truncate(*balance.DecUnit)
						beginning["avg_nav"] = balance.AvgNav.Truncate(*balance.DecNav)
						beginning["fee"] = balance.Fee.Truncate(0)
						transGroupProduct = append(transGroupProduct, beginning)
					}

					trans := make(map[string]interface{})
					trans["date"] = tr.NavDate
					trans["description"] = tr.Trans
					trans["amount"] = tr.Amount.Round(0).Truncate(*tr.DecAmount)
					trans["nav_value"] = tr.NavValue.Truncate(*tr.DecNav)
					trans["unit"] = tr.Unit.Truncate(*tr.DecUnit)
					trans["avg_nav"] = tr.AvgNavSub.Truncate(*tr.DecNav)
					trans["fee"] = tr.Fee.Truncate(0)
					transGroupProduct = append(transGroupProduct, trans)

					if (tr.TransTypeKey == uint64(1)) || (tr.TransTypeKey == uint64(4)) {
						totalSubs = totalSubs.Add(tr.Amount.Round(0)).Truncate(0)
					}
					if (tr.TransTypeKey == uint64(2)) || (tr.TransTypeKey == uint64(3)) {
						totalRedm = totalRedm.Add(tr.Amount.Round(0)).Truncate(0)
					}
					count["subs"] = totalSubs
					count["redm"] = totalRedm

					product["product_id"] = tr.ProductKey
					product["product_name"] = tr.ProductName
					product["product_bank_name"] = tr.BankName
					product["product_bank_account_name"] = tr.AccountName
					product["product_bank_account_no"] = tr.AccountNo
					product["product_code"] = tr.ProductCode
					product["currency"] = tr.Currency

					if idx == (len(transactions) - 1) {
						var balanceEndingLast models.BeginningEndingBalance
						_, err = models.GetBeginningEndingBalanceAcc(&balanceEndingLast, "ENDING BALANCE", dateakhir, strconv.FormatUint(accKeyLast, 10), strconv.FormatUint(productKeyLast, 10))

						endingbalancelast := make(map[string]interface{})
						if err != nil {

							var lastavgnav models.NavValue
							_, err = models.AdminLastAvgNav(&lastavgnav, strconv.FormatUint(productKeyLast, 10), customerKey, dateakhir)
							if err != nil {
								endingbalancelast["avg_nav"] = nol
							} else {
								endingbalancelast["avg_nav"] = lastavgnav.NavValue.Truncate(2)
							}

							var lastnav models.NavValue
							_, err = models.AdminLastNavValue(&lastnav, strconv.FormatUint(productKeyLast, 10), dateakhir)
							if err != nil {
								endingbalancelast["nav_value"] = nol
							} else {
								endingbalancelast["nav_value"] = lastnav.NavValue.Truncate(2)
							}
							dateParem, _ = time.Parse(layout, dateakhir)
							endingbalancelast["date"] = dateParem.Format(newLayout)
							endingbalancelast["description"] = "ENDING BALANCE"
							endingbalancelast["amount"] = nol
							endingbalancelast["unit"] = nol
							endingbalancelast["fee"] = nol
							transGroupProduct = append(transGroupProduct, endingbalancelast)
						} else {
							endingbalancelast["date"] = balanceEndingLast.Tanggal
							endingbalancelast["description"] = balanceEndingLast.Description
							endingbalancelast["amount"] = balanceEndingLast.Amount.Truncate(*balanceEndingLast.DecAmount)
							endingbalancelast["nav_value"] = balanceEndingLast.NavValue.Truncate(*balanceEndingLast.DecNav)
							endingbalancelast["unit"] = balanceEndingLast.Unit.Truncate(*balanceEndingLast.DecUnit)
							endingbalancelast["avg_nav"] = balanceEndingLast.AvgNav.Truncate(*balanceEndingLast.DecNav)
							endingbalancelast["fee"] = balanceEndingLast.Fee.Truncate(0)
							transGroupProduct = append(transGroupProduct, endingbalancelast)
						}

						row := make(map[string]interface{})
						row["count"] = count
						row["product"] = product
						row["list"] = transGroupProduct
						datatrans = append(datatrans, row)
					}
				} else {
					trans := make(map[string]interface{})
					trans["date"] = tr.NavDate
					trans["description"] = tr.Trans
					trans["amount"] = tr.Amount.Round(0).Truncate(*tr.DecAmount)
					trans["nav_value"] = tr.NavValue.Truncate(*tr.DecNav)
					trans["unit"] = tr.Unit.Truncate(*tr.DecUnit)
					trans["avg_nav"] = tr.AvgNavSub.Truncate(*tr.DecNav)
					trans["fee"] = tr.Fee.Truncate(0)
					transGroupProduct = append(transGroupProduct, trans)

					if (tr.TransTypeKey == 1) || (tr.TransTypeKey == 4) {
						totalSubs = totalSubs.Add(tr.Amount.Round(0)).Truncate(0)
					}
					if (tr.TransTypeKey == 2) || (tr.TransTypeKey == 3) {
						totalRedm = totalRedm.Add(tr.Amount.Round(0)).Truncate(0)
					}
					count["subs"] = totalSubs
					count["redm"] = totalRedm

					product["product_id"] = tr.ProductKey
					product["product_name"] = tr.ProductName
					product["product_bank_name"] = tr.BankName
					product["product_bank_account_name"] = tr.AccountName
					product["product_bank_account_no"] = tr.AccountNo
					product["product_code"] = tr.ProductCode
					product["currency"] = tr.Currency

					if idx == (len(transactions) - 1) {
						var balanceEndingLast models.BeginningEndingBalance
						_, err = models.GetBeginningEndingBalanceAcc(&balanceEndingLast, "ENDING BALANCE", dateakhir, strconv.FormatUint(accKeyLast, 10), strconv.FormatUint(productKeyLast, 10))

						endingbalancelast := make(map[string]interface{})
						if err != nil {

							endingbalancelast["avg_nav"] = trans["avg_nav"]

							var lastnav models.NavValue
							_, err = models.AdminLastNavValue(&lastnav, strconv.FormatUint(productKeyLast, 10), dateakhir)
							if err != nil {
								endingbalancelast["nav_value"] = nol
							} else {
								endingbalancelast["nav_value"] = lastnav.NavValue.Truncate(2)
							}
							dateParem, _ = time.Parse(layout, dateakhir)
							endingbalancelast["date"] = dateParem.Format(newLayout)
							endingbalancelast["description"] = "ENDING BALANCE"
							endingbalancelast["amount"] = nol
							endingbalancelast["unit"] = nol
							endingbalancelast["fee"] = nol
							transGroupProduct = append(transGroupProduct, endingbalancelast)
						} else {
							endingbalancelast["date"] = balanceEndingLast.Tanggal
							endingbalancelast["description"] = balanceEndingLast.Description
							endingbalancelast["amount"] = balanceEndingLast.Amount.Truncate(0)
							endingbalancelast["nav_value"] = balanceEndingLast.NavValue.Truncate(*balanceEndingLast.DecNav)
							endingbalancelast["unit"] = balanceEndingLast.Unit.Truncate(*balanceEndingLast.DecUnit)
							endingbalancelast["avg_nav"] = balanceEndingLast.AvgNav.Truncate(*balanceEndingLast.DecNav)
							endingbalancelast["fee"] = balanceEndingLast.Fee.Truncate(0)
							transGroupProduct = append(transGroupProduct, endingbalancelast)
						}

						row := make(map[string]interface{})
						row["count"] = count
						row["product"] = product
						row["list"] = transGroupProduct
						datatrans = append(datatrans, row)
					}

				}
			} else {

				var balance models.BeginningEndingBalance
				_, err = models.GetBeginningEndingBalanceAcc(&balance, "BEGINNING BALANCE", dateawal, strconv.FormatUint(tr.AccKey, 10), strconv.FormatUint(tr.ProductKey, 10))

				beginning := make(map[string]interface{})
				if err != nil {
					var lastavgnav models.NavValue
					_, err = models.AdminLastAvgNav(&lastavgnav, strconv.FormatUint(tr.ProductKey, 10), customerKey, dateawal)
					if err != nil {
						beginning["avg_nav"] = nol
					} else {
						beginning["avg_nav"] = lastavgnav.NavValue.Truncate(2)
					}

					var lastnav models.NavValue
					_, err = models.AdminLastNavValue(&lastnav, strconv.FormatUint(tr.ProductKey, 10), dateawal)
					if err != nil {
						beginning["nav_value"] = nol
					} else {
						beginning["nav_value"] = lastnav.NavValue.Truncate(2)
					}
					dateParem, _ = time.Parse(layout, dateawal)
					beginning["date"] = dateParem.Format(newLayout)
					beginning["description"] = "BEGINNING BALANCE"
					beginning["amount"] = nol
					beginning["unit"] = nol
					beginning["fee"] = nol
					transGroupProduct = append(transGroupProduct, beginning)
				} else {
					beginning["date"] = balance.Tanggal
					beginning["description"] = balance.Description
					beginning["amount"] = balance.Amount.Truncate(*balance.DecAmount)
					beginning["nav_value"] = balance.NavValue.Truncate(*balance.DecNav)
					beginning["unit"] = balance.Unit.Truncate(*balance.DecUnit)
					beginning["avg_nav"] = balance.AvgNav.Truncate(*balance.DecNav)
					beginning["fee"] = balance.Fee.Truncate(0)
					transGroupProduct = append(transGroupProduct, beginning)
				}

				trans := make(map[string]interface{})
				trans["date"] = tr.NavDate
				trans["description"] = tr.Trans
				trans["amount"] = tr.Amount.Round(0).Truncate(*balance.DecAmount)
				trans["nav_value"] = tr.NavValue.Truncate(*balance.DecNav)
				trans["unit"] = tr.Unit.Truncate(*balance.DecUnit)
				trans["avg_nav"] = tr.AvgNavSub.Truncate(*balance.DecNav)
				trans["fee"] = tr.Fee.Truncate(0)
				transGroupProduct = append(transGroupProduct, trans)

				if (tr.TransTypeKey == 1) || (tr.TransTypeKey == 4) {
					totalSubs = totalSubs.Add(tr.Amount.Round(0)).Truncate(0)
				}
				if (tr.TransTypeKey == 2) || (tr.TransTypeKey == 3) {
					totalRedm = totalRedm.Add(tr.Amount.Round(0)).Truncate(0)
				}
				count["subs"] = totalSubs
				count["redm"] = totalRedm

				product["product_id"] = tr.ProductKey
				product["product_name"] = tr.ProductName
				product["product_bank_name"] = tr.BankName
				product["product_bank_account_name"] = tr.AccountName
				product["product_bank_account_no"] = tr.AccountNo
				product["product_code"] = tr.ProductCode
				product["currency"] = tr.Currency

				if idx == (len(transactions) - 1) {
					var balanceEndingLast models.BeginningEndingBalance
					_, err = models.GetBeginningEndingBalanceAcc(&balanceEndingLast, "ENDING BALANCE", dateakhir, strconv.FormatUint(tr.AccKey, 10), strconv.FormatUint(tr.ProductKey, 10))

					endingbalancelast := make(map[string]interface{})
					if err != nil {
						var lastavgnav models.NavValue
						_, err = models.AdminLastAvgNav(&lastavgnav, strconv.FormatUint(tr.ProductKey, 10), customerKey, dateakhir)
						if err != nil {
							endingbalancelast["avg_nav"] = nol
						} else {
							endingbalancelast["avg_nav"] = lastavgnav.NavValue.Truncate(2)
						}

						var lastnav models.NavValue
						_, err = models.AdminLastNavValue(&lastnav, strconv.FormatUint(tr.ProductKey, 10), dateakhir)
						if err != nil {
							endingbalancelast["nav_value"] = nol
						} else {
							endingbalancelast["nav_value"] = lastnav.NavValue.Truncate(2)
						}
						dateParem, _ = time.Parse(layout, dateakhir)
						endingbalancelast["date"] = dateParem.Format(newLayout)
						endingbalancelast["description"] = "ENDING BALANCE"
						endingbalancelast["amount"] = nol
						endingbalancelast["unit"] = nol
						endingbalancelast["fee"] = nol
						transGroupProduct = append(transGroupProduct, endingbalancelast)
					} else {
						endingbalancelast["date"] = balanceEndingLast.Tanggal
						endingbalancelast["description"] = balanceEndingLast.Description
						endingbalancelast["amount"] = balanceEndingLast.Amount.Truncate(*balanceEndingLast.DecAmount)
						endingbalancelast["nav_value"] = balanceEndingLast.NavValue.Truncate(*balanceEndingLast.DecNav)
						endingbalancelast["unit"] = balanceEndingLast.Unit.Truncate(*balanceEndingLast.DecUnit)
						endingbalancelast["avg_nav"] = balanceEndingLast.AvgNav.Truncate(*balanceEndingLast.DecNav)
						endingbalancelast["fee"] = balanceEndingLast.Fee.Truncate(0)
						transGroupProduct = append(transGroupProduct, endingbalancelast)
					}

					row := make(map[string]interface{})
					row["count"] = count
					row["product"] = product
					row["list"] = transGroupProduct
					datatrans = append(datatrans, row)
				}
			}
			accKeyLast = tr.AccKey
			productKeyLast = tr.ProductKey
			productKey = tr.ProductKey
		}

		responseData["transaction"] = datatrans

	}

	//========== GENERATE HALAMAN HTML DAHULU ==========
	tm := template.New("account-statement-template.html")

	tm, err = tm.ParseFiles(config.BasePath + "/mail/account-statement-template.html")
	if err != nil {
		log.Println(err)
	}
	f, err := os.Create(config.BasePath + "/mail/account-statement-" + customerKey + ".html")
	if err != nil {
		log.Println("create file: ", err)
	}
	if err := tm.Execute(f, responseData); err != nil {
		log.Println(err)
	}
	f.Close()

	//========== GENERATE HALAMAN HTML DAHULU ==========

	pdfg, err := wkhtml.NewPDFGenerator()
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed send email")
	}

	// Set global options
	pdfg.Dpi.Set(300)
	pdfg.Orientation.Set(wkhtml.OrientationLandscape)
	pdfg.Grayscale.Set(false)

	// Create a new input page from an URL
	page := wkhtml.NewPage(config.BasePath + "/mail/account-statement-" + customerKey + ".html")
	// log.Println("========= LEWAT SINI ==========")

	// Set options for this page
	page.FooterRight.Set("[page]")
	page.FooterFontSize.Set(10)
	page.Zoom.Set(0.95)
	page.Allow.Set(config.BasePath + "/mail/images")

	// Add to document
	pdfg.AddPage(page)

	// Create PDF document in internal buffer
	err = pdfg.Create()
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed send email")
	}
	err = os.MkdirAll(config.BasePath+"/files/"+customerKey, 0755)
	// Write buffer contents to file on disk
	err = pdfg.WriteFile(config.BasePath + "/files/" + customerKey + "/account-statement.pdf")
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed send email")
	}
	log.Info("Success create file")

	t := template.New("index-portofolio.html")

	t, err = t.ParseFiles(config.BasePath + "/mail/index-portofolio.html")
	if err != nil {
		log.Println(err)
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, struct{ FileUrl string }{FileUrl: config.FileUrl + "/images/mail"}); err != nil {
		log.Println(err)
	}

	result := tpl.String()

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", config.EmailFrom)
	mailer.SetHeader("To", customer.Email)
	mailer.SetHeader("Subject", "[MotionFunds] Laporan Akun")
	mailer.SetBody("text/html", result)
	mailer.Attach(config.BasePath + "/files/" + customerKey + "/account-statement.pdf")

	err = lib.SendEmail(mailer)
	if err != nil {
		log.Error(err)
		return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed send email")
	} else {
		log.Info("Email sent")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = ""

	return c.JSON(http.StatusOK, response)
}
