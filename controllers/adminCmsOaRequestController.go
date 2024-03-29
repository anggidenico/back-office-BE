package controllers

import (
	"bytes"
	"database/sql"
	"html/template"
	"log"
	"math"
	"mf-bo-api/config"
	"mf-bo-api/db"
	"mf-bo-api/lib"
	"mf-bo-api/models"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo"
	"github.com/shopspring/decimal"
	"gopkg.in/gomail.v2"
)

var File_Directory_OaRequest = "/images/oa_request/"

func initAuthCs() error {
	// log.Println("initAuthCs")
	var roleKeyCs uint64
	roleKeyCs = 11

	if lib.Profile.RoleKey != roleKeyCs {
		// return lib.CustomError(http.StatusBadRequest, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}
	return nil
}

func initAuthKyc() error {
	// log.Println("initAuthKyc")

	var roleKeyKyc uint64
	roleKeyKyc = 12

	if lib.Profile.RoleKey != roleKeyKyc {
		// return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}
	return nil
}

func initAuthFundAdmin() error {
	// log.Println("initAuthFundAdmin")

	var roleKeyFundAdmin uint64
	roleKeyFundAdmin = 13

	if lib.Profile.RoleKey != roleKeyFundAdmin {
		// log.Error("User Autorizer")
		// return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}
	return nil
}

func initAuthCsKyc() error {
	// log.Println("initAuthCsKyc")

	var roleKeyCs uint64
	roleKeyCs = 11
	var roleKeyKyc uint64
	roleKeyKyc = 12

	if (lib.Profile.RoleKey != roleKeyCs) && (lib.Profile.RoleKey != roleKeyKyc) {
		// return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}
	return nil
}

func initAuthCsKycFundAdmin() error {
	// log.Println("initAuthCsKycFundAdmin")

	var roleKeyCs uint64
	roleKeyCs = 11
	var roleKeyKyc uint64
	roleKeyKyc = 12
	var roleKeyFundAdmin uint64
	roleKeyFundAdmin = 13

	if (lib.Profile.RoleKey != roleKeyCs) && (lib.Profile.RoleKey != roleKeyKyc) && (lib.Profile.RoleKey != roleKeyFundAdmin) {
		// return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}
	return nil
}

func GetOaRequestList(c echo.Context) error {
	oaRequestType := "127"
	return GetOaRequestListAdmin(c, oaRequestType)
}

func GetListPengkinianRiskProfile(c echo.Context) error {
	oaRequestType := "128"
	return GetOaRequestListAdmin(c, oaRequestType)
}

func GetListPengkinianPersonalData(c echo.Context) error {
	oaRequestType := "296"
	return GetOaRequestListAdmin(c, oaRequestType)
}

func GetOaRequestListAdmin(c echo.Context, oaRequestType string) error {

	errorAuth := initAuthCsKyc()
	if errorAuth != nil {
		// log.Error("User Autorizer")
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}

	var roleKeyCs uint64
	roleKeyCs = 11
	var roleKeyKyc uint64
	roleKeyKyc = 12

	var err error
	var status int

	oaStatusCs := "258"
	oaStatusKyc := "259"

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
			// log.Error("Limit should be number")
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
			// log.Error("Page should be number")
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
			// log.Error("Nolimit parameter should be true/false")
			return lib.CustomError(http.StatusBadRequest, "Nolimit parameter should be true/false", "Nolimit parameter should be true/false")
		}
	} else {
		noLimit = false
	}

	items := []string{"oa_request_key", "oa_request_type", "oa_entry_start", "oa_entry_end", "oa_status", "rec_order", "rec_status", "oa_risk_level"}

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
			// log.Error("Wrong input for parameter order_by")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter order_by", "Wrong input for parameter order_by")
		}
	} else {
		params["orderBy"] = "oa_request_key"
		params["orderType"] = "DESC"
	}

	//if user approval CS
	if lib.Profile.RoleKey == roleKeyCs {
		params["oa_status"] = oaStatusCs
	}
	//if user approval KYC / Complainer
	if lib.Profile.RoleKey == roleKeyKyc {
		params["oa_status"] = oaStatusKyc
	}
	params["oa_request_type"] = oaRequestType
	params["rec_status"] = "1"

	var oaRequestDB []models.OaRequest
	status, err = models.GetAllOaRequestIndividu(&oaRequestDB, limit, offset, noLimit, params, strconv.FormatUint(lib.Profile.UserID, 10))
	if err != nil {
		// log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if len(oaRequestDB) < 1 {
		// log.Error("oa not found")
		return lib.CustomError(http.StatusNotFound, "Oa Request not found", "Oa Request not found")
	}

	var lookupIds []string
	var branchIds []string
	branchIds = append(branchIds, strconv.FormatUint(uint64(1), 10))
	var agentIds []string
	agentIds = append(agentIds, strconv.FormatUint(uint64(1), 10))
	var oaRequestIds []string
	var userApprovalIds []string
	for _, oareq := range oaRequestDB {

		if oareq.Oastatus != nil {
			if _, ok := lib.Find(lookupIds, strconv.FormatUint(*oareq.Oastatus, 10)); !ok {
				lookupIds = append(lookupIds, strconv.FormatUint(*oareq.Oastatus, 10))
			}
		}

		if oareq.OaSource != nil {
			if _, ok := lib.Find(lookupIds, strconv.FormatUint(*oareq.OaSource, 10)); !ok {
				lookupIds = append(lookupIds, strconv.FormatUint(*oareq.OaSource, 10))
			}
		}

		if _, ok := lib.Find(oaRequestIds, strconv.FormatUint(oareq.OaRequestKey, 10)); !ok {
			oaRequestIds = append(oaRequestIds, strconv.FormatUint(oareq.OaRequestKey, 10))
		}

		if oareq.BranchKey != nil {
			if _, ok := lib.Find(branchIds, strconv.FormatUint(*oareq.BranchKey, 10)); !ok {
				branchIds = append(branchIds, strconv.FormatUint(*oareq.BranchKey, 10))
			}
		}

		if oareq.AgentKey != nil {
			if _, ok := lib.Find(agentIds, strconv.FormatUint(*oareq.AgentKey, 10)); !ok {
				agentIds = append(agentIds, strconv.FormatUint(*oareq.AgentKey, 10))
			}
		}

		if oareq.RecCreatedBy != nil {
			userkyc, _ := strconv.ParseUint(*oareq.RecCreatedBy, 10, 64)
			if userkyc > 0 {
				if _, ok := lib.Find(userApprovalIds, strconv.FormatUint(userkyc, 10)); !ok {
					userApprovalIds = append(userApprovalIds, strconv.FormatUint(userkyc, 10))
				}
			}
		}
	}

	var userappr []models.ScUserLogin
	if len(userApprovalIds) > 0 {
		status, err = models.GetScUserLoginIn(&userappr, userApprovalIds, "user_login_key")
		if err != nil {
			if err != sql.ErrNoRows {
				// log.Error(err.Error())
				return lib.CustomError(status, err.Error(), "Failed get data")
			}
		}
	}

	usrData := make(map[uint64]models.ScUserLogin)
	for _, usr := range userappr {
		usrData[usr.UserLoginKey] = usr
	}

	//mapping lookup
	var genLookup []models.GenLookup
	if len(lookupIds) > 0 {
		status, err = models.GetGenLookupIn(&genLookup, lookupIds, "lookup_key")
		if err != nil {
			// log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
	}
	gData := make(map[uint64]models.GenLookup)
	for _, gen := range genLookup {
		gData[gen.LookupKey] = gen
	}

	//mapping branch
	var branchs []models.MsBranch
	if len(branchIds) > 0 {
		status, err = models.GetMsBranchIn(&branchs, branchIds, "branch_key")
		if err != nil {
			// log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
	}
	bData := make(map[uint64]models.MsBranch)
	for _, br := range branchs {
		bData[br.BranchKey] = br
	}

	//mapping agent
	var agents []models.MsAgent
	if len(agentIds) > 0 {
		status, err = models.GetMsAgentIn(&agents, agentIds, "agent_key")
		if err != nil {
			// log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
	}
	aData := make(map[uint64]models.MsAgent)
	for _, ag := range agents {
		aData[ag.AgentKey] = ag
	}

	//mapping personal data
	var oaPersonalData []models.OaPersonalData
	if len(oaRequestIds) > 0 {
		status, err = models.GetOaPersonalDataIn(&oaPersonalData, oaRequestIds, "oa_request_key")
		if err != nil {
			// log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
	}
	pdData := make(map[uint64]models.OaPersonalData)
	for _, oaPD := range oaPersonalData {
		pdData[oaPD.OaRequestKey] = oaPD
	}

	var responseData []models.OaRequestListResponse
	for _, oareq := range oaRequestDB {
		var data models.OaRequestListResponse

		if oareq.Oastatus != nil {
			if n, ok := gData[*oareq.Oastatus]; ok {
				data.Oastatus = *n.LkpName
			}
		}

		if oareq.OaSource != nil {
			if n, ok := gData[*oareq.OaSource]; ok {
				data.OaSource = *n.LkpName
			}
		}

		data.OaRequestKey = oareq.OaRequestKey

		layout := "2006-01-02 15:04:05"
		newLayout := "02 Jan 2006 15:04"
		newLayoutDateBirth := "02 Jan 2006"
		date, _ := time.Parse(layout, oareq.OaEntryEnd)
		data.OaDate = date.Format(newLayout)
		data.CreatedBy = ""
		if oareq.RecCreatedBy != nil {
			usercreate, _ := strconv.ParseUint(*oareq.RecCreatedBy, 10, 64)
			if usercreate > 0 {
				if n, ok := usrData[usercreate]; ok {
					data.CreatedBy = n.UloginName
				}
			}

		}

		if n, ok := pdData[oareq.OaRequestKey]; ok {
			data.EmailAddress = n.EmailAddress
			data.PhoneNumber = n.PhoneMobile
			date, _ = time.Parse(layout, n.DateBirth)
			data.DateBirth = date.Format(newLayoutDateBirth)
			data.FullName = n.FullName
			data.IDCardNo = n.IDcardNo
		}

		var branchKey uint64
		if oareq.BranchKey != nil {
			branchKey = *oareq.BranchKey
		} else {
			branchKey = uint64(1)
		}

		var agentKey uint64
		if oareq.AgentKey != nil {
			agentKey = *oareq.AgentKey
		} else {
			agentKey = uint64(1)
		}

		if b, ok := bData[branchKey]; ok {
			data.Branch = b.BranchName
		}

		if a, ok := aData[agentKey]; ok {
			data.Agent = a.AgentName
		}

		responseData = append(responseData, data)
	}

	var countData models.OaRequestCountData
	var pagination int
	if limit > 0 {
		status, err = models.GetCountOaRequestIndividu(&countData, params, strconv.FormatUint(lib.Profile.UserID, 10))
		if err != nil {
			// log.Error(err.Error())
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
	response.Data = responseData

	return c.JSON(http.StatusOK, response)
}

func GetOaRequestData(c echo.Context) error {
	keyStr := c.Param("key")
	return ResultOaRequestData(keyStr, c, false)
}

func GetLastHistoryOaRequestData(c echo.Context) error {
	keyStr := c.Param("key")
	key, _ := strconv.ParseUint(keyStr, 10, 64)
	if key == 0 {
		return lib.CustomError(http.StatusNotFound)
	}

	var oareq models.OaRequest
	status, err := models.GetOaRequest(&oareq, keyStr)
	if err != nil {
		return lib.CustomError(status)
	}

	var lastKeyStr string

	if oareq.OaRequestType == nil {
		// log.Error("OA Request Type Null")
		return lib.CustomError(http.StatusBadRequest)
	} else if *oareq.OaRequestType == 127 { //NEW error tidak ada history
		// log.Error("OA Request Type NEW harusnya UPDATE")
		return lib.CustomError(http.StatusBadRequest)
	} else if *oareq.OaRequestType == 128 {
		if oareq.CustomerKey == nil { //Error jika belum jadi customer
			return lib.CustomError(http.StatusBadRequest)
		}
		var lastHistoryOareq models.OaRequestKeyLastHistory
		customerKey := strconv.FormatUint(*oareq.CustomerKey, 10)
		status, err := models.AdminGetLastHistoryOaRequest(&lastHistoryOareq, customerKey, keyStr)
		if err != nil {
			return lib.CustomError(status)
		}
		lastKeyStr = strconv.FormatUint(lastHistoryOareq.OaRequestKey, 10)
	}

	if lastKeyStr == "" {
		return lib.CustomError(http.StatusBadRequest)
	}

	return ResultOaRequestData(lastKeyStr, c, true)
}

func ResultOaRequestData(keyStr string, c echo.Context, isHistory bool) error {
	errorAuth := initAuthCsKycFundAdmin()
	if errorAuth != nil {
		// log.Error("User Autorizer")
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}
	var err error
	var status int
	decimal.MarshalJSONWithoutQuotes = true
	//Get parameter limit
	// keyStr := c.Param("key")
	key, _ := strconv.ParseUint(keyStr, 10, 64)
	if key == 0 {
		return lib.CustomError(http.StatusNotFound)
	}

	var oareq models.OaRequest
	status, err = models.GetOaRequest(&oareq, keyStr)
	if err != nil {
		return lib.CustomError(status)
	}

	var roleKeyCs uint64
	roleKeyCs = 11
	var roleKeyKyc uint64
	roleKeyKyc = 12
	var roleKeyFundAdmin uint64
	roleKeyFundAdmin = 13

	// log.Println(lib.Profile.RoleKey)

	strOaKey := strconv.FormatUint(*oareq.Oastatus, 10)

	if lib.Profile.RoleKey == roleKeyCs {
		if isHistory == false {
			oaStatusCs := strconv.FormatUint(uint64(258), 10)
			if strOaKey != oaStatusCs {
				// log.Error("User Autorizer")
				return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
			}
		}
	}

	if lib.Profile.RoleKey == roleKeyKyc {
		if isHistory == false {
			oaStatusKyc := strconv.FormatUint(uint64(259), 10)
			if strOaKey != oaStatusKyc {
				// log.Error("User Autorizer")
				return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
			}
		}
	}

	if lib.Profile.RoleKey == roleKeyFundAdmin {
		if isHistory == false {
			oaStatusFundAdmin1 := strconv.FormatUint(uint64(260), 10)
			oaStatusFundAdmin2 := strconv.FormatUint(uint64(261), 10)
			if (strOaKey != oaStatusFundAdmin1) && (strOaKey != oaStatusFundAdmin2) {
				// log.Error("User Autorizer")
				return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
			}
		}
	}

	var responseData models.OaRequestDetailResponse

	layout := "2006-01-02 15:04:05"
	newLayout := "02 Jan 2006"
	newLayoutOACreate := "02 Jan 2006 15:04"

	responseData.OaRequestKey = oareq.OaRequestKey
	date, _ := time.Parse(layout, oareq.OaEntryStart)
	responseData.OaEntryStart = date.Format(newLayoutOACreate)
	date, _ = time.Parse(layout, oareq.OaEntryEnd)
	responseData.OaEntryEnd = date.Format(newLayout)
	responseData.SalesCode = oareq.SalesCode

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
	if oareq.SiteReferer != nil {
		if _, ok := lib.Find(oaRequestLookupIds, strconv.FormatUint(*oareq.SiteReferer, 10)); !ok {
			oaRequestLookupIds = append(oaRequestLookupIds, strconv.FormatUint(*oareq.SiteReferer, 10))
		}
	}
	if oareq.OaSource != nil {
		if _, ok := lib.Find(oaRequestLookupIds, strconv.FormatUint(*oareq.OaSource, 10)); !ok {
			oaRequestLookupIds = append(oaRequestLookupIds, strconv.FormatUint(*oareq.OaSource, 10))
		}
	}

	//gen lookup oa request
	var lookupOaReq []models.GenLookup
	if len(oaRequestLookupIds) > 0 {
		status, err = models.GetGenLookupIn(&lookupOaReq, oaRequestLookupIds, "lookup_key")
		if err != nil {
			if err != sql.ErrNoRows {
				// log.Error(err.Error())
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

	if oareq.SiteReferer != nil {
		if n, ok := gData[*oareq.SiteReferer]; ok {
			responseData.SiteReferer = n.LkpName
		}
	}

	if oareq.OaSource != nil {
		if n, ok := gData[*oareq.OaSource]; ok {
			responseData.OaSource = n.LkpName
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

		dir := config.ImageUrl + File_Directory_OaRequest + strconv.FormatUint(oareq.OaRequestKey, 10) + "/"

		if oapersonal.PicKtp != nil && *oapersonal.PicKtp != "" {
			path := dir + *oapersonal.PicKtp
			responseData.PicKtp = &path
		}

		if oapersonal.PicSelfie != nil && *oapersonal.PicSelfie != "" {
			path := dir + *oapersonal.PicSelfie
			responseData.PicSelfie = &path
		}

		if oapersonal.RecImage1 != nil && *oapersonal.RecImage1 != "" {
			// path := dir + "signature/" + *oapersonal.RecImage1
			path := dir + *oapersonal.RecImage1
			responseData.Signature = &path
		}

		if oapersonal.PicSelfieKtp != nil && *oapersonal.PicSelfieKtp != "" {
			path := dir + *oapersonal.PicSelfieKtp
			responseData.PicSelfieKtp = &path
		}

		responseData.OccupCompany = oapersonal.OccupCompany
		responseData.OccupPhone = oapersonal.OccupPhone
		responseData.OccupWebURL = oapersonal.OccupWebUrl
		responseData.MotherMaidenName = oapersonal.MotherMaidenName
		responseData.BeneficialFullName = oapersonal.BeneficialFullName
		responseData.RelationFullName = oapersonal.RelationFullName
		responseData.PepName = oapersonal.PepName
		responseData.PepPosition = oapersonal.PepPosition

		var getFiles []models.MsFileModels
		prmGetFile := make(map[string]string)
		prmGetFile["ref_fk_domain"] = "oa_request"
		prmGetFile["ref_fk_key"] = strconv.FormatUint(oareq.OaRequestKey, 10)
		_, err = models.GetMsFileDataWithCondition(&getFiles, prmGetFile)
		if err != nil {
			return lib.CustomError(http.StatusNotFound, err.Error(), err.Error())
		}

		if len(getFiles) > 0 {
			for _, data := range getFiles {
				var files models.OaFiles
				files.FileName = *data.FileName
				files.FileRemarks = data.FileNotes
				files.FileType = data.RecAttributeId2
				files.FileUrl = *data.FileUrl
				responseData.OaFiles = append(responseData.OaFiles, files)
			}
		}

		//mapping gen lookup
		var personalDataLookupIds []string

		if _, ok := lib.Find(personalDataLookupIds, strconv.FormatUint(oapersonal.IDcardType, 10)); !ok {
			personalDataLookupIds = append(personalDataLookupIds, strconv.FormatUint(oapersonal.IDcardType, 10))
		}
		if oapersonal.Gender != nil {
			if _, ok := lib.Find(personalDataLookupIds, strconv.FormatUint(*oapersonal.Gender, 10)); !ok {
				personalDataLookupIds = append(personalDataLookupIds, strconv.FormatUint(*oapersonal.Gender, 10))
			}
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
					// log.Error(err.Error())
					return lib.CustomError(status, err.Error(), "Failed get data")
				}
			}
		}

		pData := make(map[uint64]models.GenLookup)
		for _, genLook := range lookupPersonData {
			pData[genLook.LookupKey] = genLook
		}

		if n, ok := pData[oapersonal.IDcardType]; ok {
			responseData.IDCardType = n.LkpName
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
				// log.Error("Error Personal Data not Found")
				return lib.CustomError(status, err.Error(), "Personal data not found")
			}
		} else {
			responseData.Nationality = &country.CouName
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
					// log.Error(err.Error())
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
								// log.Error(err.Error())
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
							responseData.IDcardAddress.Kabupaten = &c.CityName
						}
					}

					if p.KecamatanKey != nil {
						if c, ok := cityData[*p.KecamatanKey]; ok {
							responseData.IDcardAddress.Kecamatan = &c.CityName
						}
					}

					if p.KabupatenKey != nil {
						var city models.MsCity
						_, err = models.GetMsCityByParent(&city, strconv.FormatUint(*p.KabupatenKey, 10))
						if err != nil {
							// log.Error(err.Error())
						} else {
							responseData.IDcardAddress.Provinsi = &city.CityName
						}
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
								// log.Error(err.Error())
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

					if p.KabupatenKey != nil {
						var city models.MsCity
						_, err = models.GetMsCityByParent(&city, strconv.FormatUint(*p.KabupatenKey, 10))
						if err != nil {
							// log.Error(err.Error())
						} else {
							responseData.DomicileAddress.Provinsi = &city.CityName
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
								// log.Error(err.Error())
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

		//set bank_request
		var accBank []models.OaRequestByField
		status, err = models.GetOaRequestBankByField(&accBank, "oa_request_key", strconv.FormatUint(oareq.OaRequestKey, 10))
		if err != nil {
			responseData.BankRequest = nil
		} else {
			responseData.BankRequest = &accBank
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
				// log.Error(err.Error())
				return lib.CustomError(status, err.Error(), "Failed get data")
			}
		}
		responseData.RiskProfile = oaRiskProfile

		//mapping oa risk profile quiz
		var oaRiskProfileQuiz []models.AdminOaRiskProfileQuiz
		status, err = models.AdminGetOaRiskProfileQuizByOaRequestKey(&oaRiskProfileQuiz, strKey)
		if err != nil {
			if err != sql.ErrNoRows {
				// log.Error(err.Error())
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
					// log.Error(err.Error())
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
					// log.Error(err.Error())
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
						// log.Error(err.Error())
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
						// log.Error(err.Error())
						return lib.CustomError(status, err.Error(), "Failed get data")
					}
				}
				responseData.TinIssuanceCountry = &country.CouName
			}
		} else {
			sliceName := strings.Fields(oapersonal.FullName)
			if len(sliceName) > 0 {
				responseData.FirstName = &sliceName[0]
				if len(sliceName) > 1 {
					if len(sliceName) == 2 {
						responseData.LastName = &sliceName[1]
					} else {
						responseData.MiddleName = &sliceName[1]
						if len(sliceName) > 2 {
							lastName := strings.Join(sliceName[2:], " ")
							responseData.LastName = &lastName
						}
					}
				}
			}
		}

		//set branch
		var branch_key string
		if oareq.BranchKey != nil {
			branch_key = strconv.FormatUint(*oareq.BranchKey, 10)
		} else {
			branch_key = "1"
		}
		var branch models.MsBranch
		status, err = models.GetMsBranch(&branch, branch_key)
		if err == nil {
			// log.Println(branch.BranchKey)
			var b models.MsBranchDropdown
			b.BranchKey = branch.BranchKey
			b.BranchName = branch.BranchName
			responseData.Branch = &b
		}

		//set agent
		var agent_key string
		if oareq.AgentKey != nil {
			agent_key = strconv.FormatUint(*oareq.AgentKey, 10)
		} else {
			agent_key = "1"
		}
		var agent models.MsAgent
		status, err = models.GetMsAgent(&agent, agent_key)
		if err == nil {
			var a models.MsAgentDropdown
			a.AgentKey = agent.AgentKey
			a.AgentName = agent.AgentName
			responseData.Agent = &a
		}

		responseData.ReligionOther = nil
		responseData.JobOther = nil
		responseData.EducationOther = nil
		responseData.BusinessFieldOther = nil
		// responseData.RelationBusinessFieldOther = nil
		// responseData.RelationOccupationOther = nil
		// responseData.PositionOther = nil
		responseData.BeneficialRelationOther = nil
		responseData.ObjectivesOther = nil
		responseData.FundSourceOther = nil

		udfVal := make(map[uint64]models.UdfValue)
		paramsUdf := make(map[string]string)
		paramsUdf["u.row_data_key"] = strconv.FormatUint(oapersonal.PersonalDataKey, 10)
		paramsUdf["ui.udf_category_key"] = "1"
		var udf []models.UdfValue
		_, err := models.GetAllUdfValue(&udf, paramsUdf)
		if err == nil && len(udf) > 0 {
			for _, usr := range udf {
				udfVal[usr.UdfInfoKey] = usr
			}
		}
		if len(udfVal) > 0 {
			if ed, ok := udfVal[1]; ok { //1 = RELIGION
				if ed.UdfValues != nil {
					responseData.ReligionOther = ed.UdfValues
				}
			}
			if ed, ok := udfVal[2]; ok { //2 = OCCUPATION
				if ed.UdfValues != nil {
					responseData.JobOther = ed.UdfValues
				}
			}
			if ed, ok := udfVal[3]; ok { //3 = EDUCATION
				if ed.UdfValues != nil {
					responseData.EducationOther = ed.UdfValues
				}
			}
			if ed, ok := udfVal[4]; ok { //4 = BUSINESS_FIELDS
				if ed.UdfValues != nil {
					responseData.BusinessFieldOther = ed.UdfValues
				}
			}
			if ed, ok := udfVal[5]; ok { //5 = SOURCEOF_INCOME
				if ed.UdfValues != nil {
					responseData.FundSourceOther = ed.UdfValues
				}
			}
			if ed, ok := udfVal[6]; ok { //6 = INVESTMENT_OBJECTIVES
				if ed.UdfValues != nil {
					responseData.ObjectivesOther = ed.UdfValues
				}
			}
			if ed, ok := udfVal[7]; ok { //7 = BENEFICIAL_RELATION
				if ed.UdfValues != nil {
					responseData.BeneficialRelationOther = ed.UdfValues
				}
			}
			if ed, ok := udfVal[10]; ok { //10 = OCCUP_POSITION
				if ed.UdfValues != nil {
					responseData.JobPositionOther = ed.UdfValues
				}
			}
			if ed, ok := udfVal[8]; ok { //8 = RELATION_OCCUPATION
				if ed.UdfValues != nil {
					responseData.Relation.RelationOccupationOther = ed.UdfValues
				}
			}
			if ed, ok := udfVal[9]; ok { //9 = RELATION_BUSINESS_FIELDS
				if ed.UdfValues != nil {
					responseData.Relation.RelationBusinessFieldOther = ed.UdfValues
				}
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

func UpdateStatusApprovalCS(c echo.Context) error {
	errorAuth := initAuthCs()
	if errorAuth != nil {
		// log.Error("User Autorizer")
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}
	var err error
	var status int

	params := make(map[string]string)

	oastatus := c.FormValue("oa_status") //259 = approve --------- 258 = reject
	if oastatus == "" {
		// log.Error("Missing required parameter: oa_status")
		return lib.CustomError(http.StatusBadRequest)
	}
	n, err := strconv.ParseUint(oastatus, 10, 64)
	if err == nil && n > 0 {
		if (oastatus != "259") && (oastatus != "258") {
			// log.Error("Wrong input for parameter: oa_status must 259/258")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: oa_status", "Wrong input for parameter: oa_status")
		}
		params["oa_status"] = oastatus
	} else {
		// log.Error("Wrong input for parameter: oa_status")
		return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: oa_status", "Wrong input for parameter: oa_status")
	}

	dateLayout := "2006-01-02 15:04:05"
	strKey := strconv.FormatUint(lib.Profile.UserID, 10)

	check1notes := c.FormValue("notes")
	params["check1_notes"] = check1notes

	if oastatus != "259" { //jika reject
		if check1notes == "" {
			// log.Error("Missing required parameter notes: Notes tidak boleh kosong")
			return lib.CustomError(http.StatusBadRequest, "Notes tidak boleh kosong", "Notes tidak boleh kosong")
		}
		params["rec_status"] = "0"
		params["rec_deleted_date"] = time.Now().Format(dateLayout)
		params["rec_deleted_by"] = strKey
		params["check1_flag"] = "0"
	} else {
		params["check1_flag"] = "1"
	}

	oarequestkey := c.FormValue("oa_request_key")
	if oarequestkey == "" {
		// log.Error("Missing required parameter: oa_request_key")
		return lib.CustomError(http.StatusBadRequest)
	}
	n, err = strconv.ParseUint(oarequestkey, 10, 64)
	if err == nil && n > 0 {
		params["oa_request_key"] = oarequestkey
	} else {
		// log.Error("Wrong input for parameter: oa_request_key")
		return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: oa_request_key", "Wrong input for parameter: oa_request_key")
	}

	params["check1_date"] = time.Now().Format(dateLayout)
	params["rec_modified_date"] = time.Now().Format(dateLayout)
	params["check1_references"] = strKey
	params["rec_modified_by"] = strKey

	var oareq models.OaRequest
	status, err = models.GetOaRequest(&oareq, oarequestkey)
	if err != nil {
		return lib.CustomError(status)
	}

	strOaKey := strconv.FormatUint(*oareq.Oastatus, 10)

	oaStatusCs := strconv.FormatUint(uint64(258), 10)
	if strOaKey != oaStatusCs {
		// log.Error("User Autorizer")
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}

	_, err = models.UpdateOaRequest(params)
	if err != nil {
		// log.Error("Error update oa request")
		return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed update data")
	}

	// log.Info("Success update approved CS")

	if *oareq.OaRequestType != uint64(lib.OA_REQ_TYPE_PENGKINIAN_RISIKO_INT) {
		var oapersonal models.OaPersonalData
		strKeyOa := strconv.FormatUint(oareq.OaRequestKey, 10)
		status, err = models.GetOaPersonalDataByOaRequestKey(&oapersonal, strKeyOa)
		if err != nil {
			// log.Error("Error Personal Data not Found")
			return lib.CustomError(status, err.Error(), "Personal data not found")
		}

		if oastatus != "259" {
			strUserLoginKey := strconv.FormatUint(*oareq.UserLoginKey, 10)

			paramsUserMessage := make(map[string]string)
			paramsUserMessage["umessage_type"] = "245"
			if oareq.UserLoginKey != nil {
				paramsUserMessage["umessage_recipient_key"] = strUserLoginKey
			} else {
				paramsUserMessage["umessage_recipient_key"] = "0"
			}
			paramsUserMessage["umessage_receipt_date"] = time.Now().Format(dateLayout)
			paramsUserMessage["flag_read"] = "0"
			paramsUserMessage["umessage_sender_key"] = strKey
			paramsUserMessage["umessage_sent_date"] = time.Now().Format(dateLayout)
			paramsUserMessage["flag_sent"] = "1"
			var subject, body string
			if *oareq.OaRequestType == uint64(127) {
				subject = "Pembukaan Rekening kamu ditolak"
				body = check1notes + " Silakan menghubungi Customer Service untuk informasi lebih lanjut."
			} else {
				if *oareq.OaRequestType == uint64(128) {
					subject = "Pengkinian Profile Risiko kamu ditolak"
					body = check1notes + " Silakan menghubungi Customer Service untuk informasi lebih lanjut."
				} else {
					subject = "Pengkinian Data kamu ditolak"
					body = check1notes + " Silakan menghubungi Customer Service untuk informasi lebih lanjut."
				}
			}
			paramsUserMessage["umessage_body"] = body
			paramsUserMessage["umessage_subject"] = subject
			paramsUserMessage["umessage_category"] = "248"
			paramsUserMessage["flag_archieved"] = "0"
			paramsUserMessage["archieved_date"] = time.Now().Format(dateLayout)
			paramsUserMessage["rec_status"] = "1"
			paramsUserMessage["rec_created_date"] = time.Now().Format(dateLayout)
			paramsUserMessage["rec_created_by"] = strKey

			status, err = models.CreateScUserMessage(paramsUserMessage)
			if err != nil {
				log.Println(err.Error())
			}

			if *oareq.OaRequestType != uint64(lib.OA_REQ_TYPE_PENGKINIAN_RISIKO_INT) {
				//update personal data -> delete
				paramsPersonalDataDelete := make(map[string]string)
				paramsPersonalDataDelete["personal_data_key"] = strconv.FormatUint(oapersonal.PersonalDataKey, 10)
				paramsPersonalDataDelete["rec_status"] = "0"
				paramsPersonalDataDelete["rec_deleted_date"] = time.Now().Format(dateLayout)
				paramsPersonalDataDelete["rec_deleted_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

				_, err = models.UpdateOaPersonalData(paramsPersonalDataDelete)
				if err != nil {
					// log.Error("Error update personal data delete")
					return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed update data")
				}
				SentEmailRejectOaPengkinianToCustomer(oareq, oapersonal, check1notes)

			} else {

				getCust := make(map[string]string)
				getCust["user_login_key"] = strconv.FormatUint(*oareq.UserLoginKey, 10)
				customerData := models.GetCustomerDetailWithParams(getCust)

				sendOneSignal := make(map[string]string)
				sendOneSignal["token_notif"] = *customerData.TokenNotif
				sendOneSignal["phone_number"] = *customerData.Phone
				sendOneSignal["description"] = subject
				err = lib.CreateNotifOneSignal(sendOneSignal)
				if err != nil {
					log.Println(err.Error())
				}

				err = SendEmailRejectRiskProfilePengkinianToCustomer(*customerData.Email, customerData.FullName)
				if err != nil {
					log.Println(err.Error())
				}

			}
		}
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = nil
	return c.JSON(http.StatusOK, response)
}

func UpdateStatusApprovalCompliance(c echo.Context) error {
	errorAuth := initAuthKyc()
	if errorAuth != nil {
		// log.Error("User Autorizer")
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}
	var err error
	var status int

	params := make(map[string]string)

	oastatus := c.FormValue("oa_status") //260 = app --- 258 = reject
	if oastatus == "" {
		// log.Error("Missing required parameter: oa_status")
		return lib.CustomError(http.StatusBadRequest)
	}
	n, err := strconv.ParseUint(oastatus, 10, 64)
	if err == nil && n > 0 {
		if (oastatus != "260") && (oastatus != "258") {
			// log.Error("Wrong input for parameter: oa_status must 260/258")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: oa_status", "Wrong input for parameter: oa_status")
		}
		if oastatus == "260" {
			params["oa_status"] = oastatus
		}
	} else {
		// log.Error("Wrong input for parameter: oa_status")
		return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: oa_status", "Wrong input for parameter: oa_status")
	}

	check2notes := c.FormValue("notes")
	params["check2_notes"] = check2notes

	dateLayout := "2006-01-02 15:04:05"
	strKey := strconv.FormatUint(lib.Profile.UserID, 10)

	if oastatus != "260" { //jika reject
		if check2notes == "" {
			// log.Error("Missing required parameter notes: Notes tidak boleh kosong")
			return lib.CustomError(http.StatusBadRequest, "Notes tidak boleh kosong", "Notes tidak boleh kosong")
		}
		params["rec_status"] = "0"
		params["rec_deleted_date"] = time.Now().Format(dateLayout)
		params["rec_deleted_by"] = strKey
	}

	oarequestkey := c.FormValue("oa_request_key")
	if oarequestkey == "" {
		// log.Error("Missing required parameter: oa_request_key")
		return lib.CustomError(http.StatusBadRequest)
	}
	n, err = strconv.ParseUint(oarequestkey, 10, 64)
	if err == nil && n > 0 {
		params["oa_request_key"] = oarequestkey
	} else {
		// log.Error("Wrong input for parameter: oa_request_key")
		return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: oa_request_key", "Wrong input for parameter: oa_request_key")
	}

	params["check2_date"] = time.Now().Format(dateLayout)
	params["rec_modified_date"] = time.Now().Format(dateLayout)
	params["rec_modified_by"] = strKey
	params["check2_references"] = strKey
	if oastatus == "260" { //approve
		params["check2_flag"] = "1"
	} else { //reject
		params["check2_flag"] = "0"
		params["rec_status"] = "0"
		params["rec_deleted_date"] = time.Now().Format(dateLayout)
		params["rec_deleted_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
	}

	var oareq models.OaRequest
	status, err = models.GetOaRequest(&oareq, oarequestkey)
	if err != nil {
		return lib.CustomError(status)
	}

	oarisklevel := c.FormValue("oa_risk_level")

	if *oareq.OaRequestType != uint64(lib.OA_REQ_TYPE_PENGKINIAN_RISIKO_INT) {
		if oarisklevel == "" {
			// log.Error("Missing required parameter: oa_risk_level")
			return lib.CustomError(http.StatusBadRequest)
		}
		n, err = strconv.ParseUint(oarisklevel, 10, 64)
		if err == nil && n > 0 {
			params["oa_risk_level"] = oarisklevel
		} else {
			// log.Error("Wrong input for parameter: oa_risk_level")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: oa_risk_level", "Wrong input for parameter: oa_risk_level")
		}
	}

	strOaKey := strconv.FormatUint(*oareq.Oastatus, 10)

	oaStatusKyc := strconv.FormatUint(uint64(259), 10)
	if strOaKey != oaStatusKyc {
		// log.Error("User Autorizer")
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}

	tx, err := db.Db.Begin()

	//cek rec order
	if oareq.CustomerKey == nil {
		params["rec_order"] = "0"
	} else {
		params["rec_order"] = "0"

		var lastHistoryOareq models.OaRequestKeyLastHistory
		customerKey := strconv.FormatUint(*oareq.CustomerKey, 10)
		_, err := models.AdminGetLastHistoryOaRequest(&lastHistoryOareq, customerKey, oarequestkey)
		if err == nil {
			if lastHistoryOareq.RecOrder != nil {
				lastOrder := *lastHistoryOareq.RecOrder + 1
				params["rec_order"] = strconv.FormatUint(lastOrder, 10)
			}
		}
	}

	//update oa request
	_, err = models.UpdateOaRequest(params)
	if err != nil {
		tx.Rollback()
		// log.Error("Error update oa request")
		return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed update data")
	}
	// log.Info("Success update approved Compliance Transaction")

	if *oareq.OaRequestType != uint64(lib.OA_REQ_TYPE_PENGKINIAN_RISIKO_INT) {

		var oapersonal models.OaPersonalData
		strKeyOa := strconv.FormatUint(oareq.OaRequestKey, 10)
		status, err = models.GetOaPersonalDataByOaRequestKey(&oapersonal, strKeyOa)
		if err != nil {
			tx.Rollback()
			// log.Error("Error Personal Data not Found")
			return lib.CustomError(status, err.Error(), "Personal data not found")
		}

		if oastatus == "260" { // JIKA APPROVE

			var customerKey string

			if oareq.CustomerKey == nil { //NEW OA
				//create customer

				paramsCustomer := make(map[string]string)
				paramsCustomer["id_customer"] = "0"

				year, month, _ := time.Now().Date()

				var customer models.MsCustomer
				tahun := strconv.FormatUint(uint64(year), 10)
				bulan := strconv.FormatUint(uint64(month), 10)
				if len(bulan) == 1 {
					bulan = "0" + bulan
				}
				unitHolderLike := tahun + bulan
				status, err = models.GetLastUnitHolder(&customer, unitHolderLike)
				if err != nil {
					paramsCustomer["unit_holder_idno"] = unitHolderLike + "000001"
				} else {
					dgt := customer.UnitHolderIDno[len(customer.UnitHolderIDno)-6:]
					productKeyCek, _ := strconv.ParseUint(dgt, 10, 64)
					hasil := strconv.FormatUint(productKeyCek+1, 10)
					if len(hasil) == 1 {
						hasil = "00000" + hasil
					} else if len(hasil) == 2 {
						hasil = "0000" + hasil
					} else if len(hasil) == 3 {
						hasil = "000" + hasil
					} else if len(hasil) == 4 {
						hasil = "00" + hasil
					} else if len(hasil) == 5 {
						hasil = "0" + hasil
					}

					resultData := unitHolderLike + hasil
					paramsCustomer["unit_holder_idno"] = resultData
				}

				paramsCustomer["full_name"] = oapersonal.FullName
				paramsCustomer["investor_type"] = "263"
				paramsCustomer["customer_category"] = "265"
				paramsCustomer["cif_suspend_flag"] = "0"

				if oareq.BranchKey == nil {
					paramsCustomer["openacc_branch_key"] = "1"
				} else {
					paramsCustomer["openacc_branch_key"] = strconv.FormatUint(*oareq.BranchKey, 10)
				}

				if oareq.AgentKey == nil {
					paramsCustomer["openacc_agent_key"] = "1"
				} else {
					paramsCustomer["openacc_agent_key"] = strconv.FormatUint(*oareq.AgentKey, 10)
				}

				paramsCustomer["openacc_date"] = time.Now().Format(dateLayout)
				paramsCustomer["flag_employee"] = "0"
				paramsCustomer["flag_group"] = "0"
				paramsCustomer["merging_flag"] = "0"
				paramsCustomer["rec_status"] = "1"
				paramsCustomer["rec_created_date"] = time.Now().Format(dateLayout)
				paramsCustomer["rec_created_by"] = strKey

				sliceName := strings.Fields(oapersonal.FullName)

				if len(sliceName) > 0 {
					if len(sliceName) == 1 {
						paramsCustomer["first_name"] = sliceName[0]
						paramsCustomer["last_name"] = sliceName[0]
					}
					if len(sliceName) == 2 {
						paramsCustomer["first_name"] = sliceName[0]
						paramsCustomer["last_name"] = sliceName[1]
					}
					if len(sliceName) > 2 {
						ln := len(sliceName)
						paramsCustomer["first_name"] = sliceName[0]
						paramsCustomer["middle_name"] = sliceName[1]
						lastName := strings.Join(sliceName[2:ln], " ")
						paramsCustomer["last_name"] = lastName
					}
				}

				strNationality := strconv.FormatUint(oapersonal.Nationality, 10)
				if strNationality == "97" {
					paramsCustomer["fatca_status"] = "278"
				} else if strNationality == "225" {
					paramsCustomer["fatca_status"] = "279"
				} else {
					paramsCustomer["fatca_status"] = "280"
				}

				// CREATE CUSTOMER CODE
				NewClientCode := models.NewClientCode()
				paramsCustomer["client_code"] = NewClientCode

				status, err, requestID := models.CreateMsCustomer(paramsCustomer)
				// log.Println("========== PARAMETER  INSERT CUSTOMER ==========>>>", paramsCustomer)
				if err != nil {
					tx.Rollback()
					// log.Error("Error create customer", err.Error())
					return lib.CustomError(status, err.Error(), "failed input data")
				} else {
					// log.Println("========== BERHASIL CREATE KE MS CUSTOMER ========== ")
				}
				request, err := strconv.ParseUint(requestID, 10, 64)
				if request == 0 {
					tx.Rollback()
					// log.Error("Failed create customer")
					return lib.CustomError(http.StatusBadGateway, "failed input data", "failed input data")
				}

				customerKey = requestID

				// paramOaUpdate := make(map[string]string)
				// paramOaUpdate["customer_key"] = requestID
				// paramOaUpdate["oa_request_key"] = oarequestkey
				// paramOaUpdate["oa_status"] = "261" // Customer_Build

				// _, err = models.UpdateOaRequest(paramOaUpdate)
				// if err != nil {

				// 	tx.Rollback()
				// 	return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed update data")

				// } else {

				strUserLoginKey := strconv.FormatUint(*oareq.UserLoginKey, 10)

				paramsUserLogin := make(map[string]string)
				paramsUserLogin["customer_key"] = requestID
				paramsUserLogin["rec_modified_date"] = time.Now().Format(dateLayout)
				paramsUserLogin["rec_modified_by"] = strKey
				paramsUserLogin["ulogin_full_name"] = oapersonal.FullName
				paramsUserLogin["role_key"] = "1"
				strUserLoginKeyOa := strconv.FormatUint(*oareq.UserLoginKey, 10)
				paramsUserLogin["user_login_key"] = strUserLoginKeyOa
				_, err = models.UpdateScUserLogin(paramsUserLogin)
				if err != nil {
					tx.Rollback()
					return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed update data")
				}

				paramsAgentCustomer := make(map[string]string)
				paramsAgentCustomer["customer_key"] = requestID
				paramsAgentCustomer["agent_key"] = "1"
				if oareq.SalesCode != nil {
					var agent models.MsAgent
					var salcode string
					salcode = *oareq.SalesCode
					status, err = models.GetMsAgentByField(&agent, salcode, "agent_code")
					if err == nil {
						paramsAgentCustomer["agent_key"] = strconv.FormatUint(agent.AgentKey, 10)
					}
				}

				paramsAgentCustomer["rec_status"] = "1"
				paramsAgentCustomer["eff_date"] = oareq.OaEntryStart
				paramsAgentCustomer["rec_created_date"] = time.Now().Format(dateLayout)
				paramsAgentCustomer["rec_created_by"] = strKey
				status, err = models.CreateMsAgentCustomer(paramsAgentCustomer)
				if err != nil {
					tx.Rollback()
					// log.Error("Error create agent customer")
					return lib.CustomError(status, err.Error(), "failed input data")
				}

				tx.Commit()

				var userData models.ScUserLogin
				status, err = models.GetScUserLoginByCustomerKey(&userData, requestID)
				if err == nil {
					sendEmailApproveOa(oapersonal.FullName, userData.UloginEmail)
				}

				paramsUserMessage := make(map[string]string)
				paramsUserMessage["umessage_type"] = "245"
				if oareq.UserLoginKey != nil {
					paramsUserMessage["umessage_recipient_key"] = strUserLoginKey
				} else {
					paramsUserMessage["umessage_recipient_key"] = "0"
				}
				paramsUserMessage["umessage_receipt_date"] = time.Now().Format(dateLayout)
				paramsUserMessage["flag_read"] = "0"
				paramsUserMessage["umessage_sender_key"] = strKey
				paramsUserMessage["umessage_sent_date"] = time.Now().Format(dateLayout)
				paramsUserMessage["flag_sent"] = "1"
				subject := "Selamat! Pembukaan Rekening telah Disetujui"
				body := "Saat ini akun kamu sudah aktif dan bisa melakukan transaksi. Yuk mulai investasi sekarang juga."
				paramsUserMessage["umessage_subject"] = subject
				paramsUserMessage["umessage_body"] = body
				paramsUserMessage["umessage_category"] = "248"
				paramsUserMessage["flag_archieved"] = "0"
				paramsUserMessage["archieved_date"] = time.Now().Format(dateLayout)
				paramsUserMessage["rec_status"] = "1"
				paramsUserMessage["rec_created_date"] = time.Now().Format(dateLayout)
				paramsUserMessage["rec_created_by"] = strKey
				status, err = models.CreateScUserMessage(paramsUserMessage)
				if err != nil {
					tx.Rollback()
					return lib.CustomError(status, err.Error(), "failed input data")
				}

				getCust := make(map[string]string)
				getCust["user_login_key"] = strconv.FormatUint(*oareq.UserLoginKey, 10)
				customerData := models.GetCustomerDetailWithParams(getCust)
				sendOneSignal := make(map[string]string)
				if customerData.TokenNotif != nil {
					sendOneSignal["token_notif"] = *customerData.TokenNotif
				}
				sendOneSignal["phone_number"] = *customerData.Phone
				sendOneSignal["description"] = subject
				err = lib.CreateNotifOneSignal(sendOneSignal)
				if err != nil {
					log.Println(err.Error())
				}
				// lib.CreateNotifCustomerFromAdminByUserLoginKey(strUserLoginKey, subject, body, "TRANSACTION")

				// }

			} else { // PENGKINIAN PERSONAL DATA

				customerKey = strconv.FormatUint(*oareq.CustomerKey, 10)
				strUserLoginKey := strconv.FormatUint(*oareq.UserLoginKey, 10)

				//UPDATE SC_USER_LOGIN
				paramsScUserLogin := make(map[string]string)
				paramsScUserLogin["user_login_key"] = strUserLoginKey
				paramsScUserLogin["ulogin_name"] = oapersonal.EmailAddress
				paramsScUserLogin["ulogin_full_name"] = oapersonal.EmailAddress
				paramsScUserLogin["ulogin_email"] = oapersonal.EmailAddress
				paramsScUserLogin["ulogin_mobileno"] = oapersonal.PhoneMobile
				paramsScUserLogin["rec_modified_date"] = time.Now().Format(dateLayout)
				paramsScUserLogin["rec_modified_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
				_, err = models.UpdateScUserLogin(paramsScUserLogin)
				if err != nil {
					// log.Error("Error update user data")
					return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed update data")
				}

				// paramOaUpdate := make(map[string]string)
				// paramOaUpdate["oa_request_key"] = oarequestkey
				// paramOaUpdate["oa_status"] = "261" // Customer_Build
				// _, err = models.UpdateOaRequest(paramOaUpdate)
				// if err != nil {
				// 	tx.Rollback()
				// 	log.Println(err.Error())
				// }

				paramsUserMessage := make(map[string]string)
				paramsUserMessage["umessage_type"] = "245"
				if oareq.UserLoginKey != nil {
					paramsUserMessage["umessage_recipient_key"] = strUserLoginKey
				} else {
					paramsUserMessage["umessage_recipient_key"] = "0"
				}
				paramsUserMessage["umessage_receipt_date"] = time.Now().Format(dateLayout)
				paramsUserMessage["flag_read"] = "0"
				paramsUserMessage["umessage_sender_key"] = strKey
				paramsUserMessage["umessage_sent_date"] = time.Now().Format(dateLayout)
				paramsUserMessage["flag_sent"] = "1"
				subject := "Selamat! Pengkinian Data telah Disetujui"
				body := "Saat ini pengkinian data kamu sudah disetujui. Yuk investasi sekarang juga."
				paramsUserMessage["umessage_subject"] = subject
				paramsUserMessage["umessage_body"] = body
				paramsUserMessage["umessage_category"] = "248"
				paramsUserMessage["flag_archieved"] = "0"
				paramsUserMessage["archieved_date"] = time.Now().Format(dateLayout)
				paramsUserMessage["rec_status"] = "1"
				paramsUserMessage["rec_created_date"] = time.Now().Format(dateLayout)
				paramsUserMessage["rec_created_by"] = strKey
				status, err = models.CreateScUserMessage(paramsUserMessage)
				if err != nil {
					tx.Rollback()
					return lib.CustomError(status, err.Error(), "failed input data")
				}
				getCust := make(map[string]string)
				getCust["user_login_key"] = strconv.FormatUint(*oareq.UserLoginKey, 10)
				customerData := models.GetCustomerDetailWithParams(getCust)
				sendOneSignal := make(map[string]string)
				if customerData.TokenNotif != nil {
					sendOneSignal["token_notif"] = *customerData.TokenNotif
				}
				sendOneSignal["phone_number"] = *customerData.Phone
				sendOneSignal["description"] = subject
				err = lib.CreateNotifOneSignal(sendOneSignal)
				if err != nil {
					log.Println(err.Error())
				}
				// lib.CreateNotifCustomerFromAdminByUserLoginKey(strUserLoginKey, subject, body, "TRANSACTION")
			}

			//delete all ms_customer_bank_account by customer
			deleteParam := make(map[string]string)
			deleteParam["rec_status"] = "0"
			deleteParam["rec_modified_date"] = time.Now().Format(dateLayout)
			deleteParam["rec_modified_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
			deleteParam["rec_deleted_date"] = time.Now().Format(dateLayout)
			deleteParam["rec_deleted_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
			_, err = models.UpdateDataByField(deleteParam, "customer_key", customerKey)
			if err != nil {
				// log.Error("Error delete all ms_customer_bank_account")
			}
			//create all ms_customer_bank_account by oa_req_key
			var accBank []models.OaRequestByField
			status, err = models.GetOaRequestBankByField(&accBank, "oa_request_key", strconv.FormatUint(oareq.OaRequestKey, 10))
			if err != nil {
				// log.Error(err.Error())
			}
			if len(accBank) > 0 {
				var bindVarMsBank []interface{}
				for _, value := range accBank {
					var row []string
					row = append(row, customerKey)                                  //customer_key
					row = append(row, strconv.FormatUint(value.BankAccountKey, 10)) //bank_account_key
					row = append(row, strconv.FormatUint(value.FlagPriority, 10))   //flag_priority
					row = append(row, *value.AccountHolderName)                     //bank_account_name
					row = append(row, "1")                                          //rec_status
					row = append(row, time.Now().Format(dateLayout))                //rec_created_date
					row = append(row, strconv.FormatUint(lib.Profile.UserID, 10))   //rec_created_by
					bindVarMsBank = append(bindVarMsBank, row)
				}
				_, err = models.CreateMultipleMsCustomerBankkAccount(bindVarMsBank)
				if err != nil {
					// log.Error("========== FAILED CREATE MS CUSTOMER BANK ACCOUNT ==========" + err.Error())
					return lib.CustomError(status, err.Error(), "failed input data")
				}
			}

		} else { // JIKA REJECT

			paramsPersonalDataDelete := make(map[string]string)
			paramsPersonalDataDelete["personal_data_key"] = strconv.FormatUint(oapersonal.PersonalDataKey, 10)
			paramsPersonalDataDelete["rec_status"] = "0"
			paramsPersonalDataDelete["rec_deleted_date"] = time.Now().Format(dateLayout)
			paramsPersonalDataDelete["rec_deleted_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
			_, err = models.UpdateOaPersonalData(paramsPersonalDataDelete)
			if err != nil {
				tx.Rollback()
				// log.Error("Error update personal data delete")
				return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed update data")
			}

			paramsUserMessage := make(map[string]string)
			paramsUserMessage["umessage_type"] = "245"
			strUserLoginKey := strconv.FormatUint(*oareq.UserLoginKey, 10)
			if oareq.UserLoginKey != nil {
				paramsUserMessage["umessage_recipient_key"] = strUserLoginKey
			} else {
				paramsUserMessage["umessage_recipient_key"] = "0"
			}
			paramsUserMessage["umessage_receipt_date"] = time.Now().Format(dateLayout)
			paramsUserMessage["flag_read"] = "0"
			paramsUserMessage["umessage_sender_key"] = strKey
			paramsUserMessage["umessage_sent_date"] = time.Now().Format(dateLayout)
			paramsUserMessage["flag_sent"] = "1"
			var subject string
			body := check2notes + " Silakan menghubungi Customer Service untuk informasi lebih lanjut."
			paramsUserMessage["umessage_body"] = body
			if oareq.CustomerKey == nil { //NEW OA
				paramsUserMessage["umessage_subject"] = "Pembukaan Rekening kamu ditolak"
				subject = "Pembukaan Rekening kamu ditolak"
			} else {
				paramsUserMessage["umessage_subject"] = "Pengkinian Data kamu ditolak"
				subject = "Pengkinian Data kamu ditolak"
			}
			// log.Println(subject)
			paramsUserMessage["umessage_body"] = body
			paramsUserMessage["umessage_category"] = "248"
			paramsUserMessage["flag_archieved"] = "0"
			paramsUserMessage["archieved_date"] = time.Now().Format(dateLayout)
			paramsUserMessage["rec_status"] = "1"
			paramsUserMessage["rec_created_date"] = time.Now().Format(dateLayout)
			paramsUserMessage["rec_created_by"] = strKey
			status, err = models.CreateScUserMessage(paramsUserMessage)
			if err != nil {
				tx.Rollback()
				// log.Error("Error create user message")
			}

			SentEmailRejectOaPengkinianToCustomer(oareq, oapersonal, check2notes)
			getCust := make(map[string]string)
			getCust["user_login_key"] = strconv.FormatUint(*oareq.UserLoginKey, 10)
			customerData := models.GetCustomerDetailWithParams(getCust)
			sendOneSignal := make(map[string]string)
			if customerData.TokenNotif != nil {
				sendOneSignal["token_notif"] = *customerData.TokenNotif
			}
			sendOneSignal["phone_number"] = *customerData.Phone
			sendOneSignal["description"] = subject
			err = lib.CreateNotifOneSignal(sendOneSignal)
			if err != nil {
				log.Println(err.Error())
			}
			// lib.CreateNotifCustomerFromAdminByUserLoginKey(strUserLoginKey, subject, body, "TRANSACTION")
		}

		tx.Commit()

	} else { // JIKA REQUEST TIPE PENGKINIAN RISK PROFILE

		if oastatus == "260" { //JIKA REQUEST DI APPROVE
			log.Println("pengkinian risk profile aprrove")
			// paramOaUpdate := make(map[string]string)
			// paramOaUpdate["oa_request_key"] = oarequestkey
			// paramOaUpdate["oa_status"] = "261" // Customer_Build
			// _, err = models.UpdateOaRequest(paramOaUpdate)
			// if err != nil {
			// 	tx.Rollback()
			// 	log.Println(err.Error())
			// }

			subject := "Selamat! pengkinian profil risiko telah disetujui"
			body := "Saat ini pengkinian profil risiko kamu sudah disetujui. Yuk investasi sekarang juga."
			paramsUserMessage := make(map[string]string)
			paramsUserMessage["umessage_type"] = "245"
			strUserLoginKey := strconv.FormatUint(*oareq.UserLoginKey, 10)
			if oareq.UserLoginKey != nil {
				paramsUserMessage["umessage_recipient_key"] = strUserLoginKey
			} else {
				paramsUserMessage["umessage_recipient_key"] = "0"
			}
			paramsUserMessage["umessage_receipt_date"] = time.Now().Format(dateLayout)
			paramsUserMessage["flag_read"] = "0"
			paramsUserMessage["umessage_sender_key"] = strKey
			paramsUserMessage["umessage_sent_date"] = time.Now().Format(dateLayout)
			paramsUserMessage["flag_sent"] = "1"
			paramsUserMessage["umessage_subject"] = subject
			paramsUserMessage["umessage_body"] = body
			paramsUserMessage["umessage_category"] = "248"
			paramsUserMessage["flag_archieved"] = "0"
			paramsUserMessage["archieved_date"] = time.Now().Format(dateLayout)
			paramsUserMessage["rec_status"] = "1"
			paramsUserMessage["rec_created_date"] = time.Now().Format(dateLayout)
			paramsUserMessage["rec_created_by"] = strKey
			status, err = models.CreateScUserMessage(paramsUserMessage)
			if err != nil {
				// tx.Rollback()
				log.Println(err.Error())
			}

			getCust := make(map[string]string)
			getCust["user_login_key"] = strconv.FormatUint(*oareq.UserLoginKey, 10)
			customerData := models.GetCustomerDetailWithParams(getCust)
			sendOneSignal := make(map[string]string)
			if customerData.TokenNotif != nil {
				sendOneSignal["token_notif"] = *customerData.TokenNotif
			}
			sendOneSignal["phone_number"] = *customerData.Phone
			sendOneSignal["description"] = subject
			err = lib.CreateNotifOneSignal(sendOneSignal)
			if err != nil {
				log.Println(err.Error())
			}

		} else {
			log.Println("pengkinian risk profile reject")

			paramsUserMessage := make(map[string]string)
			paramsUserMessage["umessage_type"] = "245"
			strUserLoginKey := strconv.FormatUint(*oareq.UserLoginKey, 10)
			if oareq.UserLoginKey != nil {
				paramsUserMessage["umessage_recipient_key"] = strUserLoginKey
			} else {
				paramsUserMessage["umessage_recipient_key"] = "0"
			}
			paramsUserMessage["umessage_receipt_date"] = time.Now().Format(dateLayout)
			paramsUserMessage["flag_read"] = "0"
			paramsUserMessage["umessage_sender_key"] = strKey
			paramsUserMessage["umessage_sent_date"] = time.Now().Format(dateLayout)
			paramsUserMessage["flag_sent"] = "1"
			var subject string
			subject = "Pengkinian Data kamu ditolak"
			body := check2notes + " Silakan menghubungi Customer Service untuk informasi lebih lanjut."
			paramsUserMessage["umessage_body"] = body
			paramsUserMessage["umessage_subject"] = "Pengkinian Data kamu ditolak"
			paramsUserMessage["umessage_body"] = body
			paramsUserMessage["umessage_category"] = "248"
			paramsUserMessage["flag_archieved"] = "0"
			paramsUserMessage["archieved_date"] = time.Now().Format(dateLayout)
			paramsUserMessage["rec_status"] = "1"
			paramsUserMessage["rec_created_date"] = time.Now().Format(dateLayout)
			paramsUserMessage["rec_created_by"] = strKey
			status, err = models.CreateScUserMessage(paramsUserMessage)
			if err != nil {
				// tx.Rllback()
				log.Println(err.Error())
			}

			getCust := make(map[string]string)
			getCust["user_login_key"] = strconv.FormatUint(*oareq.UserLoginKey, 10)
			customerData := models.GetCustomerDetailWithParams(getCust)
			sendOneSignal := make(map[string]string)
			if customerData.TokenNotif != nil {
				sendOneSignal["token_notif"] = *customerData.TokenNotif
			}
			sendOneSignal["phone_number"] = *customerData.Phone
			sendOneSignal["description"] = subject
			err = lib.CreateNotifOneSignal(sendOneSignal)
			if err != nil {
				log.Println(err.Error())
			}

		}

	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = ""
	return c.JSON(http.StatusOK, response)
}

func GetOaRequestListDoTransaction(c echo.Context) error {

	errorAuth := initAuthFundAdmin()
	if errorAuth != nil {
		// log.Error("User Autorizer")
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}

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
			// log.Error("Limit should be number")
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
			// log.Error("Page should be number")
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
			// log.Error("Nolimit parameter should be true/false")
			return lib.CustomError(http.StatusBadRequest, "Nolimit parameter should be true/false", "Nolimit parameter should be true/false")
		}
	} else {
		noLimit = false
	}

	items := []string{"oa_request_key", "oa_request_type", "oa_entry_start", "oa_entry_end", "oa_status", "rec_order", "rec_status", "oa_risk_level"}

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
			// log.Error("Wrong input for parameter order_by")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter order_by", "Wrong input for parameter order_by")
		}
	} else {
		params["orderBy"] = "oa_request_key"
		params["orderType"] = "DESC"
	}

	params["rec_status"] = "1"

	var oaStatusIn []string
	oaStatusIn = append(oaStatusIn, "260") //KYC Approve
	// oaStatusIn = append(oaStatusIn, "261") //CUST BUILD
	// oaStatusIn = append(oaStatusIn, "262") //SINVEST DONE

	var oaRequestDB []models.OaRequest
	status, err = models.GetAllOaRequestDoTransaction(&oaRequestDB, limit, offset, noLimit, params, oaStatusIn, "oa_status")
	if err != nil {
		// log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if len(oaRequestDB) < 1 {
		// log.Error("oa not found")
		return lib.CustomError(http.StatusNotFound, "Oa Request not found", "Oa Request not found")
	}

	var lookupIds []string
	var oaRequestIds []string
	var branchIds []string
	branchIds = append(branchIds, strconv.FormatUint(uint64(1), 10))
	var agentIds []string
	agentIds = append(agentIds, strconv.FormatUint(uint64(1), 10))
	var userApprovalIds []string
	for _, oareq := range oaRequestDB {

		if oareq.Oastatus != nil {
			if _, ok := lib.Find(lookupIds, strconv.FormatUint(*oareq.Oastatus, 10)); !ok {
				lookupIds = append(lookupIds, strconv.FormatUint(*oareq.Oastatus, 10))
			}
		}

		if _, ok := lib.Find(oaRequestIds, strconv.FormatUint(oareq.OaRequestKey, 10)); !ok {
			oaRequestIds = append(oaRequestIds, strconv.FormatUint(oareq.OaRequestKey, 10))
		}

		if oareq.BranchKey != nil {
			if _, ok := lib.Find(branchIds, strconv.FormatUint(*oareq.BranchKey, 10)); !ok {
				branchIds = append(branchIds, strconv.FormatUint(*oareq.BranchKey, 10))
			}
		}

		if oareq.AgentKey != nil {
			if _, ok := lib.Find(agentIds, strconv.FormatUint(*oareq.AgentKey, 10)); !ok {
				agentIds = append(agentIds, strconv.FormatUint(*oareq.AgentKey, 10))
			}
		}

		if oareq.RecCreatedBy != nil {
			userkyc, _ := strconv.ParseUint(*oareq.RecCreatedBy, 10, 64)
			if userkyc > 0 {
				if _, ok := lib.Find(userApprovalIds, strconv.FormatUint(userkyc, 10)); !ok {
					userApprovalIds = append(userApprovalIds, strconv.FormatUint(userkyc, 10))
				}
			}
		}
	}

	var userappr []models.ScUserLogin
	if len(userApprovalIds) > 0 {
		status, err = models.GetScUserLoginIn(&userappr, userApprovalIds, "user_login_key")
		if err != nil {
			if err != sql.ErrNoRows {
				// log.Error(err.Error())
				return lib.CustomError(status, err.Error(), "Failed get data")
			}
		}
	}

	usrData := make(map[uint64]models.ScUserLogin)
	for _, usr := range userappr {
		usrData[usr.UserLoginKey] = usr
	}

	//mapping lookup
	var genLookup []models.GenLookup
	if len(lookupIds) > 0 {
		status, err = models.GetGenLookupIn(&genLookup, lookupIds, "lookup_key")
		if err != nil {
			// log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
	}
	gData := make(map[uint64]models.GenLookup)
	for _, gen := range genLookup {
		gData[gen.LookupKey] = gen
	}

	//mapping personal data
	var oaPersonalData []models.OaPersonalData
	if len(oaRequestIds) > 0 {
		status, err = models.GetOaPersonalDataIn(&oaPersonalData, oaRequestIds, "oa_request_key")
		if err != nil {
			// log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
	}
	pdData := make(map[uint64]models.OaPersonalData)
	for _, oaPD := range oaPersonalData {
		pdData[oaPD.OaRequestKey] = oaPD
	}

	//mapping branch
	var branchs []models.MsBranch
	if len(branchIds) > 0 {
		status, err = models.GetMsBranchIn(&branchs, branchIds, "branch_key")
		if err != nil {
			// log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
	}
	bData := make(map[uint64]models.MsBranch)
	for _, br := range branchs {
		bData[br.BranchKey] = br
	}

	//mapping agent
	var agents []models.MsAgent
	if len(agentIds) > 0 {
		status, err = models.GetMsAgentIn(&agents, agentIds, "agent_key")
		if err != nil {
			// log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
	}
	aData := make(map[uint64]models.MsAgent)
	for _, ag := range agents {
		aData[ag.AgentKey] = ag
	}

	var responseData []models.OaRequestListResponse
	for _, oareq := range oaRequestDB {
		var data models.OaRequestListResponse

		if oareq.Oastatus != nil {
			if n, ok := gData[*oareq.Oastatus]; ok {
				data.Oastatus = *n.LkpName
			}
		}

		data.OaRequestKey = oareq.OaRequestKey

		layout := "2006-01-02 15:04:05"
		newLayout := "02 Jan 2006 15:04"
		newLayoutDateBirth := "02 Jan 2006"
		date, _ := time.Parse(layout, oareq.OaEntryEnd)
		data.OaDate = date.Format(newLayout)
		data.CreatedBy = ""
		if oareq.RecCreatedBy != nil {
			usercreate, _ := strconv.ParseUint(*oareq.RecCreatedBy, 10, 64)
			if usercreate > 0 {
				if n, ok := usrData[usercreate]; ok {
					data.CreatedBy = n.UloginName
				}
			}

		}

		if n, ok := pdData[oareq.OaRequestKey]; ok {
			data.EmailAddress = n.EmailAddress
			data.PhoneNumber = n.PhoneMobile
			date, _ = time.Parse(layout, n.DateBirth)
			data.DateBirth = date.Format(newLayoutDateBirth)
			data.FullName = n.FullName
			data.IDCardNo = n.IDcardNo
		}

		var branchKey uint64
		if oareq.BranchKey != nil {
			branchKey = *oareq.BranchKey
		} else {
			branchKey = uint64(1)
		}

		var agentKey uint64
		if oareq.AgentKey != nil {
			agentKey = *oareq.AgentKey
		} else {
			agentKey = uint64(1)
		}

		if b, ok := bData[branchKey]; ok {
			data.Branch = b.BranchName
		}

		if a, ok := aData[agentKey]; ok {
			data.Agent = a.AgentName
		}

		responseData = append(responseData, data)
	}

	var countData models.OaRequestCountData
	var pagination int
	if limit > 0 {
		status, err = models.GetCountOaRequestDoTransaction(&countData, params, oaStatusIn, "oa_status")
		if err != nil {
			// log.Error(err.Error())
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
	response.Data = responseData

	return c.JSON(http.StatusOK, response)
}

func sendEmailApproveOa(fullName string, email string) {
	// Send email
	t := template.New("email-sukses-verifikasi.html")

	t, err := t.ParseFiles(config.BasePath + "/mail/email-sukses-verifikasi.html")
	if err != nil {
		log.Println(err.Error())
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl,
		struct {
			Name    string
			FileUrl string
		}{
			Name:    fullName,
			FileUrl: config.ImageUrl + "/images/mail"}); err != nil {
		log.Println(err.Error())
	}

	result := tpl.String()

	mailer := gomail.NewMessage()
	// mailer.SetHeader("From", config.EmailFrom)
	mailer.SetHeader("To", email)
	mailer.SetHeader("Subject", "[MotionFunds] Pembukaan Rekening Kamu telah Disetujui")
	mailer.SetBody("text/html", result)

	err = lib.SendEmail(mailer)
	if err != nil {
		log.Println(err.Error())
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
	// 	// log.Error(err)
	// }
}

func GetDetailPengkinianProfileRisiko(c echo.Context) error {
	keyStr := c.Param("key")
	return ResultOaProfileRisiko(keyStr, c, false)
}

func GetDetailPengkinianProfileRisikoLastHistory(c echo.Context) error {
	keyStr := c.Param("key")
	key, _ := strconv.ParseUint(keyStr, 10, 64)
	if key == 0 {
		return lib.CustomError(http.StatusNotFound)
	}

	var oareq models.OaRequest
	status, err := models.GetOaRequest(&oareq, keyStr)
	if err != nil {
		return lib.CustomError(status)
	}

	var lastKeyStr string

	if oareq.OaRequestType == nil {
		// log.Error("OA Request Type Null")
		return lib.CustomError(http.StatusBadRequest)
	} else {
		if oareq.CustomerKey == nil { //Error jika belum jadi customer
			return lib.CustomError(http.StatusBadRequest)
		}
		var lastHistoryOareq models.OaRequestKeyLastHistory
		customerKey := strconv.FormatUint(*oareq.CustomerKey, 10)
		status, err := models.AdminGetLastHistoryOaRequest(&lastHistoryOareq, customerKey, keyStr)
		if err != nil {
			return lib.CustomError(status)
		}
		lastKeyStr = strconv.FormatUint(lastHistoryOareq.OaRequestKey, 10)
	}

	if lastKeyStr == "" {
		return lib.CustomError(http.StatusBadRequest)
	}

	return ResultOaProfileRisiko(lastKeyStr, c, true)
}

func ResultOaProfileRisiko(keyStr string, c echo.Context, isHistory bool) error {
	var err error
	var status int
	decimal.MarshalJSONWithoutQuotes = true
	//Get parameter limit
	key, _ := strconv.ParseUint(keyStr, 10, 64)
	if key == 0 {
		return lib.CustomError(http.StatusNotFound)
	}

	var oareq models.OaRequest
	status, err = models.GetOaRequest(&oareq, keyStr)
	if err != nil {
		return lib.CustomError(status)
	}

	if isHistory == false {
		strRequestType := strconv.FormatUint(*oareq.OaRequestType, 10)
		if strRequestType != "128" { //Pengkinian Risk Profile
			// log.Error("Data tidak ditemukan. Data bukan pengkinian Risk Profile")
			return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
		}
	}

	var roleKeyCs uint64
	roleKeyCs = 11
	var roleKeyKyc uint64
	roleKeyKyc = 12
	var roleKeyFundAdmin uint64
	roleKeyFundAdmin = 13

	// log.Println(lib.Profile.RoleKey)

	strOaKey := strconv.FormatUint(*oareq.Oastatus, 10)

	if lib.Profile.RoleKey == roleKeyCs {
		if isHistory == false {
			oaStatusCs := strconv.FormatUint(uint64(258), 10)
			if strOaKey != oaStatusCs {
				// log.Error("User Autorizer")
				return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
			}
		}
	}

	if lib.Profile.RoleKey == roleKeyKyc {
		if isHistory == false {
			oaStatusKyc := strconv.FormatUint(uint64(259), 10)
			if strOaKey != oaStatusKyc {
				// log.Error("User Autorizer")
				return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
			}
		}
	}

	if lib.Profile.RoleKey == roleKeyFundAdmin {
		if isHistory == false {
			oaStatusFundAdmin1 := strconv.FormatUint(uint64(260), 10)
			oaStatusFundAdmin2 := strconv.FormatUint(uint64(261), 10)
			if (strOaKey != oaStatusFundAdmin1) && (strOaKey != oaStatusFundAdmin2) {
				// log.Error("User Autorizer")
				return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
			}
		}
	}

	var responseData models.OaRequestDetailRiskProfil

	layout := "2006-01-02 15:04:05"
	newLayout := "02 Jan 2006"
	newLayoutOACreate := "02 Jan 2006 15:04"

	responseData.OaRequestKey = oareq.OaRequestKey
	date, _ := time.Parse(layout, oareq.OaEntryStart)
	responseData.OaEntryStart = date.Format(newLayoutOACreate)
	date, _ = time.Parse(layout, oareq.OaEntryEnd)
	responseData.OaEntryEnd = date.Format(newLayout)

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
				// log.Error(err.Error())
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

		//mapping gen lookup
		var personalDataLookupIds []string
		if _, ok := lib.Find(personalDataLookupIds, strconv.FormatUint(oapersonal.IDcardType, 10)); !ok {
			personalDataLookupIds = append(personalDataLookupIds, strconv.FormatUint(oapersonal.IDcardType, 10))
		}
		if oapersonal.Gender != nil {
			if _, ok := lib.Find(personalDataLookupIds, strconv.FormatUint(*oapersonal.Gender, 10)); !ok {
				personalDataLookupIds = append(personalDataLookupIds, strconv.FormatUint(*oapersonal.Gender, 10))
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
		//gen lookup personal data
		var lookupPersonData []models.GenLookup
		if len(personalDataLookupIds) > 0 {
			status, err = models.GetGenLookupIn(&lookupPersonData, personalDataLookupIds, "lookup_key")
			if err != nil {
				if err != sql.ErrNoRows {
					// log.Error(err.Error())
					return lib.CustomError(status, err.Error(), "Failed get data")
				}
			}
		}

		pData := make(map[uint64]models.GenLookup)
		for _, genLook := range lookupPersonData {
			pData[genLook.LookupKey] = genLook
		}

		if n, ok := pData[oapersonal.IDcardType]; ok {
			responseData.IDCardType = *n.LkpName
		}

		if oapersonal.Gender != nil {
			if n, ok := pData[*oapersonal.Gender]; ok {
				responseData.Gender = n.LkpName
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
				// log.Error("Error Personal Data not Found")
				return lib.CustomError(status, err.Error(), "Personal data not found")
			}
		} else {
			responseData.Nationality = &country.CouName
		}

		if oapersonal.Education != nil {
			if n, ok := pData[*oapersonal.Education]; ok {
				responseData.Education = n.LkpName
			}
		}

		var oaRiskProfile []models.AdminOaRiskProfile
		status, err = models.AdminGetOaRiskProfile(&oaRiskProfile, strKey)
		if err != nil {
			if err != sql.ErrNoRows {
				// log.Error(err.Error())
				return lib.CustomError(status, err.Error(), "Failed get data")
			}
		}
		responseData.RiskProfile = oaRiskProfile[0]

		//mapping oa risk profile quiz
		var oaRiskProfileQuiz []models.AdminOaRiskProfileQuiz
		status, err = models.AdminGetOaRiskProfileQuizByOaRequestKey(&oaRiskProfileQuiz, strKey)
		if err != nil {
			if err != sql.ErrNoRows {
				// log.Error(err.Error())
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
					// log.Error(err.Error())
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

		//set branch
		var branch_key string
		if oareq.BranchKey != nil {
			branch_key = strconv.FormatUint(*oareq.BranchKey, 10)
		} else {
			branch_key = "1"
		}
		var branch models.MsBranch
		status, err = models.GetMsBranch(&branch, branch_key)
		if err == nil {
			// log.Println(branch.BranchKey)
			var b models.MsBranchDropdown
			b.BranchKey = branch.BranchKey
			b.BranchName = branch.BranchName
			responseData.Branch = &b
		}

		//set agent
		var agent_key string
		if oareq.AgentKey != nil {
			agent_key = strconv.FormatUint(*oareq.AgentKey, 10)
		} else {
			agent_key = "1"
		}
		var agent models.MsAgent
		status, err = models.GetMsAgent(&agent, agent_key)
		if err == nil {
			var a models.MsAgentDropdown
			a.AgentKey = agent.AgentKey
			a.AgentName = agent.AgentName
			responseData.Agent = &a
		}
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = responseData

	return c.JSON(http.StatusOK, response)
}

func GetDetailPengkinianPersonalData(c echo.Context) error {
	keyStr := c.Param("key")
	return ResultOaPersonalData(keyStr, c, false)
}

func GetDetailPengkinianPersonalDataLastHistory(c echo.Context) error {
	keyStr := c.Param("key")
	key, _ := strconv.ParseUint(keyStr, 10, 64)
	if key == 0 {
		return lib.CustomError(http.StatusNotFound)
	}

	var oareq models.OaRequest
	status, err := models.GetOaRequest(&oareq, keyStr)
	if err != nil {
		return lib.CustomError(status)
	}

	var lastKeyStr string

	if oareq.OaRequestType == nil {
		// log.Error("OA Request Type Null")
		return lib.CustomError(http.StatusBadRequest)
	} else {
		if oareq.CustomerKey == nil { //Error jika belum jadi customer
			// log.Println("Customer Null")
			return lib.CustomError(http.StatusBadRequest)
		}
		var lastHistoryOareq models.OaRequestKeyLastHistory
		customerKey := strconv.FormatUint(*oareq.CustomerKey, 10)
		status, err := models.AdminGetLastHistoryOaRequest(&lastHistoryOareq, customerKey, keyStr)
		if err != nil {
			// log.Println(err)
			return lib.CustomError(status)
		}
		lastKeyStr = strconv.FormatUint(lastHistoryOareq.OaRequestKey, 10)
	}

	if lastKeyStr == "" {
		return lib.CustomError(http.StatusBadRequest)
	}

	return ResultOaPersonalData(lastKeyStr, c, true)
}

func ResultOaPersonalData(keyStr string, c echo.Context, isHistory bool) error {
	var err error
	var status int
	decimal.MarshalJSONWithoutQuotes = true
	key, _ := strconv.ParseUint(keyStr, 10, 64)
	if key == 0 {
		return lib.CustomError(http.StatusNotFound)
	}

	var oareq models.OaRequest
	status, err = models.GetOaRequest(&oareq, keyStr)
	if err != nil {
		return lib.CustomError(status)
	}

	var roleKeyCs uint64
	roleKeyCs = 11
	var roleKeyKyc uint64
	roleKeyKyc = 12
	var roleKeyFundAdmin uint64
	roleKeyFundAdmin = 13

	// log.Println(lib.Profile.RoleKey)

	strOaKey := strconv.FormatUint(*oareq.Oastatus, 10)

	if lib.Profile.RoleKey == roleKeyCs {
		if isHistory == false {
			oaStatusCs := strconv.FormatUint(uint64(258), 10)
			if strOaKey != oaStatusCs {
				// log.Error("User Autorizer")
				return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
			}
		}
	}

	if lib.Profile.RoleKey == roleKeyKyc {
		if isHistory == false {
			oaStatusKyc := strconv.FormatUint(uint64(259), 10)
			if strOaKey != oaStatusKyc {
				// log.Error("User Autorizer")
				return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
			}
		}
	}

	if lib.Profile.RoleKey == roleKeyFundAdmin {
		if isHistory == false {
			oaStatusFundAdmin1 := strconv.FormatUint(uint64(260), 10)
			oaStatusFundAdmin2 := strconv.FormatUint(uint64(261), 10)
			if (strOaKey != oaStatusFundAdmin1) && (strOaKey != oaStatusFundAdmin2) {
				// log.Error("User Autorizer")
				return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
			}
		}
	}

	var responseData models.OaRequestDetailResponse

	layout := "2006-01-02 15:04:05"
	newLayout := "02 Jan 2006"
	newLayoutOACreate := "02 Jan 2006 15:04"

	responseData.OaRequestKey = oareq.OaRequestKey
	date, _ := time.Parse(layout, oareq.OaEntryStart)
	responseData.OaEntryStart = date.Format(newLayoutOACreate)
	date, _ = time.Parse(layout, oareq.OaEntryEnd)
	responseData.OaEntryEnd = date.Format(newLayout)

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
	if oareq.SiteReferer != nil {
		if _, ok := lib.Find(oaRequestLookupIds, strconv.FormatUint(*oareq.SiteReferer, 10)); !ok {
			oaRequestLookupIds = append(oaRequestLookupIds, strconv.FormatUint(*oareq.SiteReferer, 10))
		}
	}

	//gen lookup oa request
	var lookupOaReq []models.GenLookup
	if len(oaRequestLookupIds) > 0 {
		status, err = models.GetGenLookupIn(&lookupOaReq, oaRequestLookupIds, "lookup_key")
		if err != nil {
			if err != sql.ErrNoRows {
				// log.Error(err.Error())
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

	if oareq.SiteReferer != nil {
		if n, ok := gData[*oareq.SiteReferer]; ok {
			responseData.SiteReferer = n.LkpName
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
			path := dir + *oapersonal.RecImage1
			responseData.Signature = &path
		}

		responseData.OccupCompany = oapersonal.OccupCompany
		responseData.OccupPhone = oapersonal.OccupPhone
		responseData.OccupWebURL = oapersonal.OccupWebUrl
		responseData.MotherMaidenName = oapersonal.MotherMaidenName
		responseData.BeneficialFullName = oapersonal.BeneficialFullName
		responseData.RelationFullName = oapersonal.RelationFullName
		responseData.PepName = oapersonal.PepName
		responseData.PepPosition = oapersonal.PepPosition

		//mapping gen lookup
		var personalDataLookupIds []string

		if _, ok := lib.Find(personalDataLookupIds, strconv.FormatUint(oapersonal.IDcardType, 10)); !ok {
			personalDataLookupIds = append(personalDataLookupIds, strconv.FormatUint(oapersonal.IDcardType, 10))
		}
		if oapersonal.Gender != nil {
			if _, ok := lib.Find(personalDataLookupIds, strconv.FormatUint(*oapersonal.Gender, 10)); !ok {
				personalDataLookupIds = append(personalDataLookupIds, strconv.FormatUint(*oapersonal.Gender, 10))
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
		if oapersonal.PepStatus != nil {
			if _, ok := lib.Find(personalDataLookupIds, strconv.FormatUint(*oapersonal.PepStatus, 10)); !ok {
				personalDataLookupIds = append(personalDataLookupIds, strconv.FormatUint(*oapersonal.PepStatus, 10))
			}
		}
		//gen lookup personal data
		var lookupPersonData []models.GenLookup
		if len(personalDataLookupIds) > 0 {
			status, err = models.GetGenLookupIn(&lookupPersonData, personalDataLookupIds, "lookup_key")
			if err != nil {
				if err != sql.ErrNoRows {
					// log.Error(err.Error())
					return lib.CustomError(status, err.Error(), "Failed get data")
				}
			}
		}

		pData := make(map[uint64]models.GenLookup)
		for _, genLook := range lookupPersonData {
			pData[genLook.LookupKey] = genLook
		}

		if n, ok := pData[oapersonal.IDcardType]; ok {
			responseData.IDCardType = n.LkpName
		}

		if oapersonal.Gender != nil {
			if n, ok := pData[*oapersonal.Gender]; ok {
				responseData.Gender = n.LkpName
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
		if oapersonal.PepStatus != nil {
			if n, ok := pData[*oapersonal.PepStatus]; ok {
				responseData.PepStatus = n.LkpName
			}
		}
		responseData.PepName = oapersonal.PepName

		var country models.MsCountry

		strCountry := strconv.FormatUint(oapersonal.Nationality, 10)
		status, err = models.GetMsCountry(&country, strCountry)
		if err != nil {
			if err != sql.ErrNoRows {
				// log.Error("Error Personal Data not Found")
				return lib.CustomError(status, err.Error(), "Personal data not found")
			}
		} else {
			responseData.Nationality = &country.CouName
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
					// log.Error(err.Error())
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
								// log.Error(err.Error())
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
							responseData.IDcardAddress.Kabupaten = &c.CityName
						}
					}

					if p.KecamatanKey != nil {
						if c, ok := cityData[*p.KecamatanKey]; ok {
							responseData.IDcardAddress.Kecamatan = &c.CityName
						}
					}

					var city models.MsCity
					_, err = models.GetMsCityByParent(&city, strconv.FormatUint(*p.KabupatenKey, 10))
					if err != nil {
						// log.Error(err.Error())
					} else {
						responseData.IDcardAddress.Provinsi = &city.CityName
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
								// log.Error(err.Error())
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

					var city models.MsCity
					if p.KabupatenKey != nil {
						_, err = models.GetMsCityByParent(&city, strconv.FormatUint(*p.KabupatenKey, 10))
					}
					if err != nil {
						// log.Error(err.Error())
					} else {
						responseData.DomicileAddress.Provinsi = &city.CityName
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
								// log.Error(err.Error())
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

		//set bank_request
		var accBank []models.OaRequestByField
		status, err = models.GetOaRequestBankByField(&accBank, "oa_request_key", strconv.FormatUint(oareq.OaRequestKey, 10))
		log.Println(accBank)
		log.Println(err)
		if err != nil {
			responseData.BankRequest = nil
		} else {
			responseData.BankRequest = &accBank
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

		//add response field Sinvest
		if oareq.CustomerKey != nil {
			var customer models.MsCustomer
			strCustomerKey := strconv.FormatUint(*oareq.CustomerKey, 10)
			status, err = models.GetMsCustomer(&customer, strCustomerKey)
			if err != nil {
				if err != sql.ErrNoRows {
					// log.Error(err.Error())
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
						// log.Error(err.Error())
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
						// log.Error(err.Error())
						return lib.CustomError(status, err.Error(), "Failed get data")
					}
				}
				responseData.TinIssuanceCountry = &country.CouName
			}
		} else {
			sliceName := strings.Fields(oapersonal.FullName)
			if len(sliceName) > 0 {
				responseData.FirstName = &sliceName[0]
				if len(sliceName) > 1 {
					responseData.MiddleName = &sliceName[1]
					if len(sliceName) > 1 {
						if len(sliceName) == 2 {
							responseData.LastName = &sliceName[1]
						} else {
							responseData.MiddleName = &sliceName[1]
							if len(sliceName) > 2 {
								lastName := strings.Join(sliceName[2:], " ")
								responseData.LastName = &lastName
							}
						}
					}
				}
			}
		}

		//set branch
		var branch_key string
		if oareq.BranchKey != nil {
			branch_key = strconv.FormatUint(*oareq.BranchKey, 10)
		} else {
			branch_key = "1"
		}
		var branch models.MsBranch
		status, err = models.GetMsBranch(&branch, branch_key)
		if err == nil {
			// log.Println(branch.BranchKey)
			var b models.MsBranchDropdown
			b.BranchKey = branch.BranchKey
			b.BranchName = branch.BranchName
			responseData.Branch = &b
		}

		//set agent
		var agent_key string
		if oareq.AgentKey != nil {
			agent_key = strconv.FormatUint(*oareq.AgentKey, 10)
		} else {
			agent_key = "1"
		}
		var agent models.MsAgent
		status, err = models.GetMsAgent(&agent, agent_key)
		if err == nil {
			var a models.MsAgentDropdown
			a.AgentKey = agent.AgentKey
			a.AgentName = agent.AgentName
			responseData.Agent = &a
		}

		responseData.ReligionOther = nil
		responseData.JobOther = nil
		responseData.EducationOther = nil
		responseData.BusinessFieldOther = nil
		// responseData.RelationBusinessFieldOther = nil
		// responseData.RelationOccupationOther = nil
		// responseData.PositionOther = nil
		responseData.BeneficialRelationOther = nil
		responseData.ObjectivesOther = nil
		responseData.FundSourceOther = nil

		udfVal := make(map[uint64]models.UdfValue)
		paramsUdf := make(map[string]string)
		paramsUdf["u.row_data_key"] = strconv.FormatUint(oapersonal.PersonalDataKey, 10)
		paramsUdf["ui.udf_category_key"] = "1"
		var udf []models.UdfValue
		_, err := models.GetAllUdfValue(&udf, paramsUdf)
		if err == nil && len(udf) > 0 {
			for _, usr := range udf {
				udfVal[usr.UdfInfoKey] = usr
			}
		}
		if len(udfVal) > 0 {
			if ed, ok := udfVal[1]; ok { //1 = RELIGION
				if ed.UdfValues != nil {
					responseData.ReligionOther = ed.UdfValues
				}
			}
			if ed, ok := udfVal[2]; ok { //2 = OCCUPATION
				if ed.UdfValues != nil {
					responseData.JobOther = ed.UdfValues
				}
			}
			if ed, ok := udfVal[3]; ok { //3 = EDUCATION
				if ed.UdfValues != nil {
					responseData.EducationOther = ed.UdfValues
				}
			}
			if ed, ok := udfVal[4]; ok { //4 = BUSINESS_FIELDS
				if ed.UdfValues != nil {
					responseData.BusinessFieldOther = ed.UdfValues
				}
			}
			if ed, ok := udfVal[5]; ok { //5 = SOURCEOF_INCOME
				if ed.UdfValues != nil {
					responseData.FundSourceOther = ed.UdfValues
				}
			}
			if ed, ok := udfVal[6]; ok { //6 = INVESTMENT_OBJECTIVES
				if ed.UdfValues != nil {
					responseData.ObjectivesOther = ed.UdfValues
				}
			}
			if ed, ok := udfVal[7]; ok { //7 = BENEFICIAL_RELATION
				if ed.UdfValues != nil {
					responseData.BeneficialRelationOther = ed.UdfValues
				}
			}
			if ed, ok := udfVal[10]; ok { //10 = OCCUP_POSITION
				if ed.UdfValues != nil {
					responseData.JobPositionOther = ed.UdfValues
				}
			}
			if ed, ok := udfVal[8]; ok { //8 = RELATION_OCCUPATION
				if ed.UdfValues != nil {
					responseData.Relation.RelationOccupationOther = ed.UdfValues
				}
			}
			if ed, ok := udfVal[9]; ok { //9 = RELATION_BUSINESS_FIELDS
				if ed.UdfValues != nil {
					responseData.Relation.RelationBusinessFieldOther = ed.UdfValues
				}
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

func SentEmailOaPengkinianToBackOffice(
	oaRequest models.OaRequest,
	oaPersonalData models.OaPersonalData,
	roleKey string) {

	var err error
	var mailTemp, subject string
	mailParam := make(map[string]string)
	if roleKey == "11" {
		mailParam["BackOfficeGroup"] = "Customer Service"
	} else if roleKey == "12" {
		mailParam["BackOfficeGroup"] = "Compliance"
	} else if roleKey == "13" {
		mailParam["BackOfficeGroup"] = "FundAdmin"
	}

	layout := "2006-01-02 15:04:05"
	dateLayout := "02 Jan 2006"
	date, _ := time.Parse(layout, oaPersonalData.DateBirth)
	mailParam["Name"] = oaPersonalData.FullName
	mailParam["Identitas"] = oaPersonalData.IDcardNo
	mailParam["Npwp"] = "-"
	mailParam["TanggalLahir"] = date.Format(dateLayout)
	mailParam["Email"] = oaPersonalData.EmailAddress
	mailParam["NoHp"] = oaPersonalData.PhoneMobile
	mailParam["FileUrl"] = config.ImageUrl + "/images/mail"

	if *oaRequest.OaRequestType == uint64(127) { // oa new
		subject = "[MotionFunds] Mohon Verifikasi Pembukaan Rekening Reksa Dana"
		mailTemp = "email-oa-to-cs-kyc-fundadmin.html"
	} else { // pengkinian
		if *oaRequest.OaRequestType == uint64(128) {
			mailParam["JenisPengkinian"] = "Pengkinian Profile Risiko"
			subject = "[MotionFunds] Mohon Verifikasi Pengkinian Profile Risiko"
		} else {
			mailParam["JenisPengkinian"] = "Pengkinian Personal Data"
			subject = "[MotionFunds] Mohon Verifikasi Pengkinian Personal Data"
		}
		mailTemp = "email-pengkinian-to-cs-kyc-fundadmin.html"
	}

	paramsScLogin := make(map[string]string)
	paramsScLogin["role_key"] = roleKey
	paramsScLogin["rec_status"] = "1"
	var userLogin []models.ScUserLogin
	_, err = models.GetAllScUserLogin(&userLogin, 0, 0, paramsScLogin, true)
	if err != nil {
		// log.Error("User BO tidak ditemukan")
		// log.Error(err)
	} else {
		t := template.New(mailTemp)

		t, err = t.ParseFiles(config.BasePath + "/mail/" + mailTemp)
		if err != nil {
			// log.Error("Failed send mail: " + err.Error())
		} else {
			for _, scLogin := range userLogin {
				strUserCat := strconv.FormatUint(scLogin.UserCategoryKey, 10)
				if (strUserCat == "2") || (strUserCat == "3") {
					var tpl bytes.Buffer
					if err := t.Execute(&tpl, mailParam); err != nil {
						// log.Error("Failed send mail: " + err.Error())
					} else {
						result := tpl.String()

						mailer := gomail.NewMessage()
						// mailer.SetHeader("From", config.EmailFrom)
						mailer.SetHeader("To", scLogin.UloginEmail)
						mailer.SetHeader("Subject", subject)
						mailer.SetBody("text/html", result)

						err = lib.SendEmail(mailer)
						if err != nil {
							// log.Error("Failed send mail to: " + scLogin.UloginEmail)
							// log.Error("Failed send mail: " + err.Error())
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
						// 	// log.Error("Failed send mail to: " + scLogin.UloginEmail)
						// 	// log.Error("Failed send mail: " + err.Error())
						// }
					}
				}
			}
		}
	}
}

func SentEmailOaPengkinianToSales(
	oaRequest models.OaRequest,
	oaPersonalData models.OaPersonalData) {

	var err error
	var mailTemp, subject string
	mailParam := make(map[string]string)

	layout := "2006-01-02 15:04:05"
	dateLayout := "02 Jan 2006"
	date, _ := time.Parse(layout, oaPersonalData.DateBirth)
	mailParam["Name"] = oaPersonalData.FullName
	mailParam["Identitas"] = oaPersonalData.IDcardNo
	mailParam["Npwp"] = "-"
	mailParam["TanggalLahir"] = date.Format(dateLayout)
	mailParam["Email"] = oaPersonalData.EmailAddress
	mailParam["NoHp"] = oaPersonalData.PhoneMobile
	mailParam["FileUrl"] = config.ImageUrl + "/images/mail"

	if *oaRequest.OaRequestType == uint64(127) { // oa new
		subject = "[MotionFunds] Mohon Verifikasi Pembukaan Rekening Reksa Dana"
		mailTemp = "email-oa-to-sales.html"
	} else { // pengkinian
		if *oaRequest.OaRequestType == uint64(128) {
			mailParam["JenisPengkinian"] = "Pengkinian Profile Risiko"
			subject = "[MotionFunds] Mohon Verifikasi Pengkinian Profile Risiko"
		} else {
			mailParam["JenisPengkinian"] = "Pengkinian Personal Data"
			subject = "[MotionFunds] Mohon Verifikasi Pengkinian Personal Data"
		}
		mailTemp = "email-pengkinian-to-sales.html"
	}

	var agentKey string
	if oaRequest.AgentKey == nil {
		agentKey = "1"
	} else {
		agentKey = strconv.FormatUint(*oaRequest.AgentKey, 10)
	}

	var agent models.MsAgent
	_, err = models.GetMsAgent(&agent, agentKey)
	if err != nil {
		// log.Error("Agent not found")
	} else {
		if agent.AgentEmail != nil {
			t := template.New(mailTemp)

			t, err = t.ParseFiles(config.BasePath + "/mail/" + mailTemp)
			if err != nil {
				// log.Error("Failed send mail to sales: " + err.Error())
			} else {
				var tpl bytes.Buffer
				if err := t.Execute(&tpl, mailParam); err != nil {
					// log.Error("Failed send mail to sales: " + err.Error())
				} else {
					result := tpl.String()

					mailer := gomail.NewMessage()
					// mailer.SetHeader("From", config.EmailFrom)
					mailer.SetHeader("To", *agent.AgentEmail)
					mailer.SetHeader("Subject", subject)
					mailer.SetBody("text/html", result)

					err = lib.SendEmail(mailer)
					if err != nil {
						// log.Error("Failed send mail to sales : " + *agent.AgentEmail)
						// log.Error("Failed send mail: " + err.Error())
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
					// 	// log.Error("Failed send mail to sales : " + *agent.AgentEmail)
					// 	// log.Error("Failed send mail: " + err.Error())
					// }
				}
			}
		} else {
			// log.Error("Agent tidak punya email")
		}
	}
}

func SentEmailRejectOaPengkinianToCustomer(
	oaRequest models.OaRequest,
	oaPersonalData models.OaPersonalData,
	status string) {

	var err error
	var mailTemp, subject string
	mailParam := make(map[string]string)

	layout := "2006-01-02 15:04:05"
	dateLayout := "02 Jan 2006"
	date, _ := time.Parse(layout, oaRequest.OaEntryStart)
	mailParam["FileUrl"] = config.ImageUrl + "/images/mail"
	mailParam["Nama"] = oaPersonalData.FullName
	mailParam["KTP"] = oaPersonalData.IDcardNo
	mailParam["Tanggal"] = date.Format(dateLayout)
	mailParam["Status"] = status

	mailTemp = "email-oa-pengkinian-rejected.html"
	if *oaRequest.OaRequestType == uint64(127) { // oa new
		subject = "[MotionFunds] Pembukaan Rekening Kamu Belum Dapat Disetujui"
		mailParam["Judul"] = "Pendaftaran Rekening Belum Dapat Disetujui"
		mailParam["Keterangan"] = "pendaftaran rekening yang kamu ajukan belum dapat dapat disetujui:"
	} else { // pengkinian
		subject = "[MotionFunds] Mohon Verifikasi Pembukaan Rekening Reksa Dana"
		mailParam["Judul"] = "Pengkinian Data Kamu Tidak Dapat Diproses"
		mailParam["Keterangan"] = "pengkinian dara yan kamu ajukan belum dapat kami proses lebih lanjut:"
	}

	var userLogin models.ScUserLogin
	_, err = models.GetScUserKey(&userLogin, strconv.FormatUint(*oaRequest.UserLoginKey, 10))
	if err != nil {
		// log.Error("User tidak ditemukan")
		// log.Error(err)
		return
	} else {
		// log.Println("User Ada")
		t := template.New(mailTemp)

		t, err = t.ParseFiles(config.BasePath + "/mail/" + mailTemp)
		if err != nil {
			// log.Error("Failed send mail: " + err.Error())
		} else {
			// log.Println("File Template Ada")
			var tpl bytes.Buffer
			if err := t.Execute(&tpl, mailParam); err != nil {
				// log.Error("Failed send mail: " + err.Error())
			} else {
				result := tpl.String()

				mailer := gomail.NewMessage()
				// mailer.SetHeader("From", config.EmailFrom)
				mailer.SetHeader("To", userLogin.UloginEmail)
				mailer.SetHeader("Subject", subject)
				mailer.SetBody("text/html", result)

				err = lib.SendEmail(mailer)
				if err != nil {
					// log.Error("Failed send mail to: " + userLogin.UloginEmail)
					// log.Error("Failed send mail: " + err.Error())
				} else {
					// log.Println("Sukses kirim email: " + userLogin.UloginEmail)
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
				// 	// log.Error("Failed send mail to: " + userLogin.UloginEmail)
				// 	// log.Error("Failed send mail: " + err.Error())
				// } else {
				// 	// log.Println("Sukses kirim email: " + userLogin.UloginEmail)
				// }
			}
		}
	}
}
