package controllers

import (
	"encoding/json"
	"math"
	"mf-bo-api/config"
	"mf-bo-api/lib"
	"mf-bo-api/models"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo"
	"github.com/shopspring/decimal"
)

func GetListAgentDropdown(c echo.Context) error {

	var err error
	var status int

	var msAgent []models.MsAgentDropdown

	status, err = models.GetMsAgentDropdown(&msAgent)

	if err != nil {
		// log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if len(msAgent) < 1 {
		// log.Error("Agent not found")
		return lib.CustomError(http.StatusNotFound, "Agent not found", "Agent not found")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = msAgent

	return c.JSON(http.StatusOK, response)
}

func GetListAgentLastBranch(c echo.Context) error {

	var err error
	var status int

	keyStr := c.Param("branch_key")
	key, _ := strconv.ParseUint(keyStr, 10, 64)
	if key == 0 {
		return lib.CustomError(http.StatusNotFound)
	}

	var agent []models.MsAgentLastBranch

	status, err = models.GetMsAgentLastBranch(&agent, keyStr)

	if err != nil {
		// log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if len(agent) < 1 {
		// log.Error("Agent not found")
		return lib.CustomError(http.StatusNotFound, "Agent not found", "Agent not found")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = agent

	return c.JSON(http.StatusOK, response)
}

func AdminGetListMsAgent(c echo.Context) error {

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

	items := []string{"branch_name", "agent_id", "agent_code", "agent_name", "agent_category", "agent_channel"}

	params := make(map[string]string)
	orderBy := c.QueryParam("order_by")
	if orderBy != "" {
		_, found := lib.Find(items, orderBy)
		if found {

			var ord string
			if orderBy == "branch_name" {
				ord = "b.branch_name"
			} else if orderBy == "agent_id" {
				ord = "a.agent_id"
			} else if orderBy == "agent_code" {
				ord = "a.agent_code"
			} else if orderBy == "agent_name" {
				ord = "a.agent_name"
			} else if orderBy == "agent_category" {
				ord = "cat.lkp_name"
			} else if orderBy == "agent_channel" {
				ord = "cha.lkp_name"
			}
			params["orderBy"] = ord
			orderType := c.QueryParam("order_type")
			if (orderType == "asc") || (orderType == "ASC") || (orderType == "desc") || (orderType == "DESC") {
				params["orderType"] = orderType
			}
		} else {
			// log.Error("Wrong input for parameter order_by")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter order_by", "Wrong input for parameter order_by")
		}
	} else {
		params["orderBy"] = "a.agent_key"
		params["orderType"] = "DESC"
	}

	searchLike := c.QueryParam("search_like")
	branchKey := c.QueryParam("branch_key")
	if branchKey != "" {
		params["b.branch_key"] = branchKey
	}

	var agent []models.ListAgentAdmin

	status, err = models.AdminGetListAgent(&agent, limit, offset, params, searchLike, noLimit)

	if err != nil {
		// log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if len(agent) < 1 {
		// log.Error("Agent not found")
		return lib.CustomError(http.StatusNotFound, "Agent not found", "Agent not found")
	}

	var countData models.CountData
	var pagination int
	if limit > 0 {
		status, err = models.CountAdminGetListAgent(&countData, params, searchLike)
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
	response.Data = agent

	return c.JSON(http.StatusOK, response)
}

func AdminDeleteMsAgent(c echo.Context) error {
	var err error

	params := make(map[string]string)

	keyStr := c.FormValue("agent_key")
	key, _ := strconv.ParseUint(keyStr, 10, 64)
	if key == 0 {
		// log.Error("Missing required parameter: agent_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: agent_key", "Missing required parameter: agent_key")
	}

	dateLayout := "2006-01-02 15:04:05"
	params["agent_key"] = keyStr
	params["rec_status"] = "0"
	params["rec_deleted_date"] = time.Now().Format(dateLayout)
	params["rec_deleted_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

	_, err = models.UpdateMsAgent(params)
	if err != nil {
		// log.Error("Error delete ms_branch")
		return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed delete data")
	}

	//delete ms_agent_branch
	paramsAgentBranch := make(map[string]string)
	paramsAgentBranch["rec_status"] = "0"
	paramsAgentBranch["rec_deleted_date"] = time.Now().Format(dateLayout)
	paramsAgentBranch["rec_deleted_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

	_, err = models.UpdateDeleteBranchAgent(paramsAgentBranch, "agent_key", keyStr)
	if err != nil {
		// log.Error("Error delete ms_agent_branch")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = nil
	return c.JSON(http.StatusOK, response)
}

func AdminCreateMsAgent(c echo.Context) error {
	var err error
	var status int

	params := make(map[string]string)

	branchKey := c.FormValue("branch_key")
	if branchKey != "" {
		n, err := strconv.ParseUint(branchKey, 10, 64)
		if err != nil || n == 0 {
			// log.Error("Wrong input for parameter: branch_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: branch_key", "Wrong input for parameter: branch_key")
		}
	} else {
		// log.Error("Missing required parameter: branch_key")
		return lib.CustomError(http.StatusBadRequest, "branch_key can not be blank", "branch_key can not be blank")
	}

	agentId := c.FormValue("agent_id")
	if agentId == "" {
		// log.Error("Missing required parameter: agent_id")
		return lib.CustomError(http.StatusBadRequest, "agent_id can not be blank", "agent_id can not be blank")
	} else {
		n, err := strconv.ParseUint(agentId, 10, 64)
		if err != nil || n == 0 {
			// log.Error("Wrong input for parameter: agent_id")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: agent_id", "Wrong input for parameter: agent_id")
		}
		//validate unique agent_id
		var countData models.CountData
		status, err = models.CountMsAgentValidateUnique(&countData, "agent_id", agentId, "")
		if err != nil {
			// log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
		if int(countData.CountData) > int(0) {
			// log.Error("Missing required parameter: agent_id")
			return lib.CustomError(http.StatusBadRequest, "agent_id already used", "agent_id already used")
		}
		params["agent_id"] = agentId
	}

	agentCode := c.FormValue("agent_code")
	if agentCode == "" {
		// log.Error("Missing required parameter: agent_code")
		return lib.CustomError(http.StatusBadRequest, "agent_code can not be blank", "agent_code can not be blank")
	} else {
		//validate unique agent_code
		var countData models.CountData
		status, err = models.CountMsAgentValidateUnique(&countData, "agent_code", agentCode, "")
		if err != nil {
			// log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
		if int(countData.CountData) > int(0) {
			// log.Error("Missing required parameter: agent_code")
			return lib.CustomError(http.StatusBadRequest, "agent_code already used", "agent_code already used")
		}
		params["agent_code"] = agentCode
	}

	agentName := c.FormValue("agent_name")
	if agentName == "" {
		// log.Error("Missing required parameter: agent_name")
		return lib.CustomError(http.StatusBadRequest, "agent_name can not be blank", "agent_name can not be blank")
	} else {
		//validate unique agent_name
		var countData models.CountData
		status, err = models.CountMsAgentValidateUnique(&countData, "agent_name", agentName, "")
		if err != nil {
			// log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
		if int(countData.CountData) > int(0) {
			// log.Error("Missing required parameter: agent_name")
			return lib.CustomError(http.StatusBadRequest, "agent_name already used", "agent_name already used")
		}
		params["agent_name"] = agentName
	}

	agentEmail := c.FormValue("agent_email")
	if agentEmail == "" {
		// log.Error("Missing required parameter: agent_email")
		return lib.CustomError(http.StatusBadRequest, "agent_email can not be blank", "agent_email can not be blank")
	} else {
		params["agent_email"] = agentEmail
	}

	agentShortName := c.FormValue("agent_short_name")
	if agentShortName != "" {
		params["agent_short_name"] = agentShortName
	}

	agentCategory := c.FormValue("agent_category")
	if agentCategory != "" {
		n, err := strconv.ParseUint(agentCategory, 10, 64)
		if err != nil || n == 0 {
			// log.Error("Wrong input for parameter: agent_category")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: agent_category", "Wrong input for parameter: agent_category")
		}
		params["agent_category"] = agentCategory
	}

	agentChannel := c.FormValue("agent_channel")
	if agentChannel != "" {
		n, err := strconv.ParseUint(agentChannel, 10, 64)
		if err != nil || n == 0 {
			// log.Error("Wrong input for parameter: agent_channel")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: agent_channel", "Wrong input for parameter: agent_channel")
		}
		params["agent_channel"] = agentChannel
	}

	referenceCode := c.FormValue("reference_code")
	if referenceCode != "" {
		params["reference_code"] = referenceCode
	}

	remarks := c.FormValue("remarks")
	if remarks != "" {
		params["remarks"] = remarks
	}

	recOrder := c.FormValue("rec_order")
	if recOrder != "" {
		n, err := strconv.ParseUint(recOrder, 10, 64)
		if err != nil || n == 0 {
			// log.Error("Wrong input for parameter: rec_order")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: rec_order", "Wrong input for parameter: rec_order")
		}
		params["rec_order"] = recOrder
	}

	dateLayout := "2006-01-02 15:04:05"
	params["rec_created_date"] = time.Now().Format(dateLayout)
	params["rec_created_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
	params["rec_status"] = "1"

	status, err, lastID := models.CreateMsAgent(params)
	if err != nil {
		// log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed input data")
	}

	dateLayoutNow := "2006-01-02"
	//insert ms_agent_branch
	paramsAgentBranch := make(map[string]string)
	paramsAgentBranch["agent_key"] = lastID
	paramsAgentBranch["branch_key"] = branchKey
	paramsAgentBranch["eff_date"] = time.Now().Format(dateLayoutNow) + " 00:00:00"
	paramsAgentBranch["rec_created_date"] = time.Now().Format(dateLayout)
	paramsAgentBranch["rec_created_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
	paramsAgentBranch["rec_status"] = "1"
	paramsAgentBranch["remarks"] = remarks
	status, err = models.CreateMsAgentBranch(paramsAgentBranch)
	if err != nil {
		// log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed input data")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = ""

	return c.JSON(http.StatusOK, response)
}

func AdminUpdateMsAgent(c echo.Context) error {
	var err error
	var status int

	params := make(map[string]string)

	agentKey := c.FormValue("agent_key")
	if agentKey != "" {
		n, err := strconv.ParseUint(agentKey, 10, 64)
		if err != nil || n == 0 {
			// log.Error("Wrong input for parameter: agent_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: agent_key", "Wrong input for parameter: agent_key")
		}
		params["agent_key"] = agentKey
	} else {
		// log.Error("Missing required parameter: agent_key")
		return lib.CustomError(http.StatusBadRequest, "agent_key can not be blank", "agent_key can not be blank")
	}

	branchKey := c.FormValue("branch_key")
	if branchKey != "" {
		n, err := strconv.ParseUint(branchKey, 10, 64)
		if err != nil || n == 0 {
			// log.Error("Wrong input for parameter: branch_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: branch_key", "Wrong input for parameter: branch_key")
		}
	} else {
		// log.Error("Missing required parameter: branch_key")
		return lib.CustomError(http.StatusBadRequest, "branch_key can not be blank", "branch_key can not be blank")
	}

	agentId := c.FormValue("agent_id")
	if agentId == "" {
		// log.Error("Missing required parameter: agent_id")
		return lib.CustomError(http.StatusBadRequest, "agent_id can not be blank", "agent_id can not be blank")
	} else {
		n, err := strconv.ParseUint(agentId, 10, 64)
		if err != nil || n == 0 {
			// log.Error("Wrong input for parameter: agent_id")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: agent_id", "Wrong input for parameter: agent_id")
		}
		//validate unique agent_id
		var countData models.CountData
		status, err = models.CountMsAgentValidateUnique(&countData, "agent_id", agentId, agentKey)
		if err != nil {
			// log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
		if int(countData.CountData) > int(0) {
			// log.Error("Missing required parameter: agent_id")
			return lib.CustomError(http.StatusBadRequest, "agent_id already used", "agent_id already used")
		}
		params["agent_id"] = agentId
	}

	agentCode := c.FormValue("agent_code")
	if agentCode == "" {
		// log.Error("Missing required parameter: agent_code")
		return lib.CustomError(http.StatusBadRequest, "agent_code can not be blank", "agent_code can not be blank")
	} else {
		//validate unique agent_code
		var countData models.CountData
		status, err = models.CountMsAgentValidateUnique(&countData, "agent_code", agentCode, agentKey)
		if err != nil {
			// log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
		if int(countData.CountData) > int(0) {
			// log.Error("Missing required parameter: agent_code")
			return lib.CustomError(http.StatusBadRequest, "agent_code already used", "agent_code already used")
		}
		params["agent_code"] = agentCode
	}

	agentName := c.FormValue("agent_name")
	if agentName == "" {
		// log.Error("Missing required parameter: agent_name")
		return lib.CustomError(http.StatusBadRequest, "agent_name can not be blank", "agent_name can not be blank")
	} else {
		//validate unique agent_name
		var countData models.CountData
		status, err = models.CountMsAgentValidateUnique(&countData, "agent_name", agentName, agentKey)
		if err != nil {
			// log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
		if int(countData.CountData) > int(0) {
			// log.Error("Missing required parameter: agent_name")
			return lib.CustomError(http.StatusBadRequest, "agent_name already used", "agent_name already used")
		}
		params["agent_name"] = agentName
	}

	agentEmail := c.FormValue("agent_email")
	if agentEmail == "" {
		// log.Error("Missing required parameter: agent_email")
		return lib.CustomError(http.StatusBadRequest, "agent_email can not be blank", "agent_email can not be blank")
	} else {
		params["agent_email"] = agentEmail
	}

	agentShortName := c.FormValue("agent_short_name")
	if agentShortName != "" {
		params["agent_short_name"] = agentShortName
	}

	agentCategory := c.FormValue("agent_category")
	if agentCategory != "" {
		n, err := strconv.ParseUint(agentCategory, 10, 64)
		if err != nil || n == 0 {
			// log.Error("Wrong input for parameter: agent_category")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: agent_category", "Wrong input for parameter: agent_category")
		}
		params["agent_category"] = agentCategory
	}

	agentChannel := c.FormValue("agent_channel")
	if agentChannel != "" {
		n, err := strconv.ParseUint(agentChannel, 10, 64)
		if err != nil || n == 0 {
			// log.Error("Wrong input for parameter: agent_channel")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: agent_channel", "Wrong input for parameter: agent_channel")
		}
		params["agent_channel"] = agentChannel
	}

	referenceCode := c.FormValue("reference_code")
	if referenceCode != "" {
		params["reference_code"] = referenceCode
	}

	remarks := c.FormValue("remarks")
	if remarks != "" {
		params["remarks"] = remarks
	}

	recOrder := c.FormValue("rec_order")
	if recOrder != "" {
		n, err := strconv.ParseUint(recOrder, 10, 64)
		if err != nil || n == 0 {
			// log.Error("Wrong input for parameter: rec_order")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: rec_order", "Wrong input for parameter: rec_order")
		}
		params["rec_order"] = recOrder
	}

	dateLayout := "2006-01-02 15:04:05"
	params["rec_modified_date"] = time.Now().Format(dateLayout)
	params["rec_modified_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
	params["rec_status"] = "1"

	status, err = models.UpdateMsAgent(params)
	if err != nil {
		// log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed input data")
	}

	dateLayoutNow := "2006-01-02"
	//insert/update ms_agent_branch
	paramsGet := make(map[string]string)
	paramsGet["agent_key"] = agentKey
	paramsGet["branch_key"] = branchKey
	paramsGet["rec_status"] = "1"

	var msAgentBranch []models.MsAgentBranch

	status, err = models.GetAllMsAgentBranch(&msAgentBranch, 10, 10, paramsGet, true)

	if err != nil {
		// log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if len(msAgentBranch) < 1 {
		//delete
		paramsDeleteAgentBranch := make(map[string]string)
		paramsDeleteAgentBranch["rec_status"] = "0"
		paramsDeleteAgentBranch["rec_deleted_date"] = time.Now().Format(dateLayout)
		paramsDeleteAgentBranch["rec_deleted_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

		_, err = models.UpdateDeleteBranchAgent(paramsDeleteAgentBranch, "agent_key", agentKey)
		if err != nil {
			// log.Error("Error delete ms_agent_branch")
		}

		paramsAgentBranch := make(map[string]string)
		paramsAgentBranch["agent_key"] = agentKey
		paramsAgentBranch["branch_key"] = branchKey
		paramsAgentBranch["eff_date"] = time.Now().Format(dateLayoutNow) + " 00:00:00"
		paramsAgentBranch["rec_created_date"] = time.Now().Format(dateLayout)
		paramsAgentBranch["rec_created_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
		paramsAgentBranch["rec_status"] = "1"
		paramsAgentBranch["remarks"] = remarks
		status, err = models.CreateMsAgentBranch(paramsAgentBranch)
		if err != nil {
			// log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed input data")
		}
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = ""

	return c.JSON(http.StatusOK, response)
}

func AdminDetailMsAgent(c echo.Context) error {
	var err error

	agentKey := c.Param("agent_key")
	if agentKey == "" {
		// log.Error("Missing required parameter: agent_key")
		return lib.CustomError(http.StatusBadRequest, "agent_key can not be blank", "agent_key can not be blank")
	} else {
		n, err := strconv.ParseUint(agentKey, 10, 64)
		if err != nil || n == 0 {
			// log.Error("Wrong input for parameter: agent_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: agent_key", "Wrong input for parameter: agent_key")
		}
	}

	var branch models.MsAgentBranchDetail
	_, err = models.AdminGetDetailAgent(&branch, agentKey)
	if err != nil {
		// log.Error("Agent not found")
		return lib.CustomError(http.StatusBadRequest, "Agent not found", "Agent not found")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = branch

	return c.JSON(http.StatusOK, response)
}

func AdminGetListMsAgentCustomer(c echo.Context) error {
	decimal.MarshalJSONWithoutQuotes = true

	var err error
	var status int
	var paramsToResponse []interface{}
	timeNow := lib.TIMENOW_DATEONLYFORMAT
	var total_aum decimal.Decimal

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

	// ISI SEARCH
	searchLike := c.QueryParam("search_like")

	// ISI TYPE SEARCH
	searchType := c.QueryParam("search_type")

	// ISI Effective Date
	eff_date := c.QueryParam("effective_date")

	var agent_cust []models.AgentCustomerDetailList

	status, err = models.AdminGetListAgentCustomer(&agent_cust, searchLike, searchType, eff_date, limit, offset, noLimit)

	if err != nil {
		// log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if len(agent_cust) < 1 {
		// log.Error("Agent and Customer data not found")
		return lib.CustomError(http.StatusNotFound, "Agent and Customer data not found", "Agent and Customer data not found")
	}

	for _, dt_agent_cust := range agent_cust {
		paramsToAppend := make(map[string]interface{})
		paramsToAppend["effective_date"] = dt_agent_cust.EffectiveDate
		paramsToAppend["agent_key"] = dt_agent_cust.AgentKey
		paramsToAppend["agent_name"] = dt_agent_cust.AgentName
		paramsToAppend["agent_code"] = dt_agent_cust.AgentCode
		// paramsToAppend["agent_short_name"] = dt_agent_cust.AgentShortName
		paramsToAppend["customer_key"] = dt_agent_cust.CustomerKey
		paramsToAppend["customer_name"] = dt_agent_cust.CustomerName
		paramsToAppend["unit_holder_id"] = dt_agent_cust.UnitHolderID
		// ========== AUM HARDCODE AUM SEMENTARA ==========
		// ac := accounting.Accounting{Precision: 0, Thousand: "."}
		// paramsToAppend["aum"] = ac.FormatMoney(rand.Intn(999999999-100000) + 100000)

		//========== PERHITUNGAN AUM REAL ==========

		// GET TR_ACCOUNT KEY
		var accKey string
		paramsAcc := make(map[string]string)
		paramsAcc["rec_status"] = "1"
		paramsAcc["customer_key"] = dt_agent_cust.CustomerKey
		var trAccountDB []models.TrAccount
		status, err = models.GetAllTrAccount(&trAccountDB, paramsAcc)
		if len(trAccountDB) > 0 {
			// GET TR_ACCOUNT_AGENT KEY
			accKey = strconv.FormatUint(trAccountDB[0].AccKey, 10)
			if trAccountDB[0].SubSuspendFlag != nil && *trAccountDB[0].SubSuspendFlag == 1 {
				// log.Error("Account suspended")
				return lib.CustomError(status, "Account suspended ", "Account suspended ")
			}

			// GET ACCOUNT AGENT ACA KEY
			paramsAccAgent := make(map[string]string)
			paramsAccAgent["acc_key"] = accKey
			paramsAccAgent["rec_status"] = "1"
			// var acaKey string
			var accountAgentDB []models.TrAccountAgent
			status, err = models.GetAllTrAccountAgent(&accountAgentDB, paramsAccAgent)
			if len(accountAgentDB) > 0 {
				// acaKey = strconv.FormatUint(accountAgentDB[0].AcaKey, 10)
			}
			// Println(acaKey)
			// GET BALANCE AMOUNT
			var bal models.AumBalanceUnitStruct
			_, err = models.AumBalanceQuery(&bal, accKey, timeNow)
			if err != nil {
				// log.Println(err)
			}
			balance_unit := bal.BalanceUnit

			for _, dt := range trAccountDB {
				// GET NAV VALUE
				var navVal models.AumNavValueStruct
				prodKey := strconv.FormatUint(dt.ProductKey, 10)
				_, err = models.AumNavValueQuery(&navVal, prodKey, timeNow)
				if err != nil {
					// log.Println(err)
				}
				navValues := navVal.NavValue
				// SUM THE AUM
				aums := navValues.Mul(balance_unit)
				total_aum = total_aum.Add(aums)
			}

			paramsToAppend["aum"] = total_aum.Truncate(2)
		}
		// log.Println(accKey)

		paramsToResponse = append(paramsToResponse, paramsToAppend)
	}

	var countData models.CountData
	var pagination int
	if limit > 0 {
		status, err = models.CountAdminGetListAgentCustomer(&countData, searchLike, searchType)
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
	response.Data = paramsToResponse

	return c.JSON(http.StatusOK, response)
}

func AdminAgentCustomerMovement(c echo.Context) error {
	// err := json.NewDecoder(c.Request().Body).Decode(&jsonBody)
	var paramsToResponse []interface{}

	rec_created_date := time.Now().Format(lib.TIMESTAMPFORMAT)
	rec_created_by := lib.Profile.Name

	new_agent_key := c.FormValue("agent_key")

	new_eff_date := c.FormValue("effective_date")

	dataAgentCust := c.FormValue("data")

	if dataAgentCust != "" {
		var dataAgentCustSlice []interface{}
		err := json.Unmarshal([]byte(dataAgentCust), &dataAgentCustSlice)
		if err != nil {
			// log.Error(err.Error())
			// log.Error("Missing required parameter: data")
			return lib.CustomError(http.StatusBadRequest, err.Error(), "Wrong input for parameter: bank_account")
		}

		if len(dataAgentCustSlice) > 0 {
			for _, val := range dataAgentCustSlice {
				dataToAppend := make(map[string]string)
				valueMap := val.(map[string]interface{})

				if val, ok := valueMap["customer_key"]; ok {
					dataToAppend["customer_key"] = val.(string)
				}

				dataToAppend["agent_key"] = new_agent_key
				dataToAppend["effective_date"] = new_eff_date

				status, err := models.AdminAddAgentCustomer(dataToAppend["agent_key"], dataToAppend["customer_key"], dataToAppend["effective_date"], rec_created_date, rec_created_by)
				if err != nil {
					// log.Error("Failed insert to ms_agent_customer: " + err.Error())
					return lib.CustomError(status, err.Error(), "failed input data agent customer movement")
				}
				paramsToResponse = append(paramsToResponse, dataToAppend)

			}
		}
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = paramsToResponse
	return c.JSON(http.StatusOK, response)
}
