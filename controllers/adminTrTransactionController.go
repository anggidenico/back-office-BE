package controllers

import (
	"bytes"
	"crypto/tls"
	"database/sql"
	"fmt"
	"html/template"
	"log"
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

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/labstack/echo"
	"github.com/leekchan/accounting"
	"github.com/shopspring/decimal"
	"gopkg.in/gomail.v2"
)

func initAuthBranchEntryHoEntry() error {
	var roleKeyBranchEntry uint64
	roleKeyBranchEntry = 7
	var roleKeyHoEntry uint64
	roleKeyHoEntry = 10

	if (lib.Profile.RoleKey != roleKeyBranchEntry) && (lib.Profile.RoleKey != roleKeyHoEntry) {
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}
	return nil
}

func initAuthTransactionAdmin() error {
	roles := []string{"7", "10", "11", "12", "13"}
	strRoleLogin := strconv.FormatUint(lib.Profile.RoleKey, 10)
	_, found := lib.Find(roles, strRoleLogin)
	if !found {
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}

	return nil
}

func GetTransactionApprovalList(c echo.Context) error {
	errorAuth := initAuthCsKyc()
	if errorAuth != nil {
		// log.Error("User Autorizer")
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}

	var roleKeyCs uint64
	roleKeyCs = 11
	var roleKeyKyc uint64
	roleKeyKyc = 12

	var transStatusKey []string

	//if user approval CS
	if lib.Profile.RoleKey == roleKeyCs {
		transStatusKey = append(transStatusKey, "2")
	}
	//if user approval KYC / Complainer
	if lib.Profile.RoleKey == roleKeyKyc {
		transStatusKey = append(transStatusKey, "4")
	}

	return getListAdmin(transStatusKey, c, nil)
}

func GetTransactionCutOffList(c echo.Context) error {
	errorAuth := initAuthFundAdmin()
	if errorAuth != nil {
		// log.Error("User Autorizer")
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}

	var transStatusKey []string
	transStatusKey = append(transStatusKey, "5")

	return getListAdmin(transStatusKey, c, nil)
}

func GetTransactionBatchList(c echo.Context) error {
	errorAuth := initAuthFundAdmin()
	if errorAuth != nil {
		// log.Error("User Autorizer")
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}

	//date
	postnavdate := c.QueryParam("nav_date")
	if postnavdate == "" {
		// log.Error("Missing required parameter: nav_date")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: nav_date", "Missing required parameter: nav_date")
	}

	var transStatusKey []string
	transStatusKey = append(transStatusKey, "6")

	return getListAdmin(transStatusKey, c, &postnavdate)
}

func GetTransactionConfirmationList(c echo.Context) error {
	errorAuth := initAuthFundAdmin()
	if errorAuth != nil {
		// log.Error("User Autorizer")
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}

	var transStatusKey []string
	transStatusKey = append(transStatusKey, "7")

	return getListAdmin(transStatusKey, c, nil)
}

func GetTransactionCorrectionAdminList(c echo.Context) error {
	var transStatusKey []string
	transStatusKey = append(transStatusKey, "1")

	return getListAdmin(transStatusKey, c, nil)
}

func GetTransactionPostingList(c echo.Context) error {
	errorAuth := initAuthFundAdmin()
	if errorAuth != nil {
		// log.Error("User Autorizer")
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}

	var transStatusKey []string
	transStatusKey = append(transStatusKey, "8")

	//date
	postnavdate := c.QueryParam("nav_date")
	if postnavdate != "" {
		return getListAdmin(transStatusKey, c, &postnavdate)
	} else {
		return getListAdmin(transStatusKey, c, nil)
	}

}

func GetTransactionUnpostingList(c echo.Context) error {
	errorAuth := initAuthFundAdmin()
	if errorAuth != nil {
		// log.Error("User Autorizer")
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}

	var transStatusKey []string
	transStatusKey = append(transStatusKey, "9")

	return getListAdmin(transStatusKey, c, nil)
}

func getListAdmin(transStatusKey []string, c echo.Context, postnavdate *string) error {

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

	items := []string{"transaction_key", "branch_key", "agent_key", "customer_key", "product_key", "trans_date", "trans_amount", "trans_bank_key"}

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
		params["orderBy"] = "transaction_key"
		params["orderType"] = "ASC"
	}

	params["rec_status"] = "1"
	if postnavdate != nil {
		params["nav_date"] = *postnavdate
	}

	//if user admin role 7 branch
	//if user category  = 3 -> user branch, 2 = user HO
	var userCategory uint64
	userCategory = 3
	if lib.Profile.UserCategoryKey == userCategory {
		// log.Println(lib.Profile)
		if lib.Profile.BranchKey != nil {
			strBranchKey := strconv.FormatUint(*lib.Profile.BranchKey, 10)
			params["c.openacc_branch_key"] = strBranchKey
		} else {
			// log.Error("User Branch haven't Branch")
			return lib.CustomError(http.StatusBadRequest, "Wrong User Branch haven't Branch", "Wrong User Branch haven't Branch")
		}
	}

	transTypeKey := c.QueryParam("trans_type_key")
	if transTypeKey != "" {
		params["trans_type_key"] = transTypeKey
	}

	var trTransaction []models.TrTransaction

	_, cekConfirm := lib.Find(transStatusKey, "7")
	_, cekPosting := lib.Find(transStatusKey, "8")
	_, cekCutoff := lib.Find(transStatusKey, "5")
	_, cekBatch := lib.Find(transStatusKey, "6")
	_, cekPosted := lib.Find(transStatusKey, "9")
	if !cekConfirm && !cekPosting && !cekCutoff && !cekPosted && !cekBatch {
		status, err = models.AdminGetAllTrTransaction(&trTransaction, limit, offset, noLimit, params, transStatusKey, "trans_status_key", false, strconv.FormatUint(lib.Profile.UserID, 10))
	} else {
		status, err = models.AdminGetAllTrTransaction(&trTransaction, limit, offset, noLimit, params, transStatusKey, "trans_status_key", true, strconv.FormatUint(lib.Profile.UserID, 10))
	}

	if err != nil {
		// log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if len(trTransaction) < 1 {
		// log.Error("transaction not found")
		return lib.CustomError(http.StatusNotFound, "Transaction not found", "Transaction not found")
	}

	var branchIds []string
	var agentIds []string
	var customerIds []string
	var productIds []string
	var transTypeIds []string
	var bankIds []string
	var lookupIds []string
	var transIds []string
	var parentTransIds []string
	for _, tr := range trTransaction {
		if tr.TransSource != nil {
			if _, ok := lib.Find(lookupIds, strconv.FormatUint(*tr.TransSource, 10)); !ok {
				lookupIds = append(lookupIds, strconv.FormatUint(*tr.TransSource, 10))
			}
		}

		if tr.BranchKey != nil {
			if _, ok := lib.Find(branchIds, strconv.FormatUint(*tr.BranchKey, 10)); !ok {
				branchIds = append(branchIds, strconv.FormatUint(*tr.BranchKey, 10))
			}
		}
		if tr.AgentKey != nil {
			if _, ok := lib.Find(agentIds, strconv.FormatUint(*tr.AgentKey, 10)); !ok {
				agentIds = append(agentIds, strconv.FormatUint(*tr.AgentKey, 10))
			}
		}
		if _, ok := lib.Find(customerIds, strconv.FormatUint(tr.CustomerKey, 10)); !ok {
			customerIds = append(customerIds, strconv.FormatUint(tr.CustomerKey, 10))
		}
		if _, ok := lib.Find(productIds, strconv.FormatUint(tr.ProductKey, 10)); !ok {
			productIds = append(productIds, strconv.FormatUint(tr.ProductKey, 10))
		}
		if _, ok := lib.Find(transTypeIds, strconv.FormatUint(tr.TransTypeKey, 10)); !ok {
			transTypeIds = append(transTypeIds, strconv.FormatUint(tr.TransTypeKey, 10))
		}
		if tr.TransBankKey != nil {
			if _, ok := lib.Find(bankIds, strconv.FormatUint(*tr.TransBankKey, 10)); !ok {
				bankIds = append(bankIds, strconv.FormatUint(*tr.TransBankKey, 10))
			}
		}

		strTransTypeKey := strconv.FormatUint(tr.TransTypeKey, 10)
		if strTransTypeKey == "4" {
			if tr.ParentKey != nil {
				if _, ok := lib.Find(parentTransIds, strconv.FormatUint(*tr.ParentKey, 10)); !ok {
					parentTransIds = append(parentTransIds, strconv.FormatUint(*tr.ParentKey, 10))
				}
			}
		}
		if strTransTypeKey == "1" { //
			if tr.PaymentMethod != nil {
				if _, ok := lib.Find(lookupIds, strconv.FormatUint(*tr.PaymentMethod, 10)); !ok {
					lookupIds = append(lookupIds, strconv.FormatUint(*tr.PaymentMethod, 10))
				}
			}
			if _, ok := lib.Find(transIds, strconv.FormatUint(tr.TransactionKey, 10)); !ok {
				transIds = append(transIds, strconv.FormatUint(tr.TransactionKey, 10))
			}
		}
	}

	//mapping branch
	var msBranch []models.MsBranch
	if len(branchIds) > 0 {
		status, err = models.GetMsBranchIn(&msBranch, branchIds, "branch_key")
		if err != nil {
			// log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
	}
	branchData := make(map[uint64]models.MsBranch)
	for _, b := range msBranch {
		branchData[b.BranchKey] = b
	}

	//gen lookup
	var lookupOaReq []models.GenLookup
	if len(lookupIds) > 0 {
		status, err = models.GetGenLookupIn(&lookupOaReq, lookupIds, "lookup_key")
		if err != nil {
			if err != sql.ErrNoRows {
				// log.Error(err.Error())
			}
		}
	}
	gData := make(map[uint64]models.GenLookup)
	for _, gen := range lookupOaReq {
		gData[gen.LookupKey] = gen
	}

	//tr settlement branch
	var trSettle []models.TrTransactionSettlement
	if len(transIds) > 0 {
		status, err = models.GetTrTransactionSettlementIn(&trSettle, transIds, "transaction_key")
		if err != nil {
			// log.Error(err.Error())
		}
	}
	var paymentChannelIds []string
	settleData := make(map[uint64]models.TrTransactionSettlement)
	for _, sd := range trSettle {
		if sd.TransactionKey != nil {
			settleData[*sd.TransactionKey] = sd
		}
		if _, ok := lib.Find(paymentChannelIds, strconv.FormatUint(sd.SettlePaymentMethod, 10)); !ok {
			paymentChannelIds = append(paymentChannelIds, strconv.FormatUint(sd.SettlePaymentMethod, 10))
		}
	}

	//ms payment channel branch
	var pChannel []models.MsPaymentChannel
	if len(paymentChannelIds) > 0 {
		status, err = models.GetMsPaymentChannelIn(&pChannel, paymentChannelIds, "pchannel_key")
		if err != nil {
			// log.Error(err.Error())
		}
	}
	channelData := make(map[uint64]models.MsPaymentChannel)
	for _, pc := range pChannel {
		channelData[pc.PchannelKey] = pc
	}

	//mapping agent
	var msAgent []models.MsAgent
	if len(agentIds) > 0 {
		status, err = models.GetMsAgentIn(&msAgent, agentIds, "agent_key")
		if err != nil {
			// log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
	}
	agentData := make(map[uint64]models.MsAgent)
	for _, a := range msAgent {
		agentData[a.AgentKey] = a
	}

	//mapping customer
	var msCustomer []models.MsCustomer
	if len(customerIds) > 0 {
		status, err = models.GetMsCustomerIn(&msCustomer, customerIds, "customer_key")
		if err != nil {
			// log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
	}
	customerData := make(map[uint64]models.MsCustomer)
	for _, c := range msCustomer {
		customerData[c.CustomerKey] = c
	}

	//user customer
	var userLogin []models.ScUserLogin
	if len(customerIds) > 0 {
		status, err = models.GetScUserLoginIn(&userLogin, customerIds, "customer_key")
		if err != nil {
			// log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
	}
	userLoginData := make(map[uint64]models.ScUserLogin)
	for _, c := range userLogin {
		userLoginData[*c.CustomerKey] = c
	}

	//mapping parent transaction
	var parentTrans []models.TrTransaction
	if len(parentTransIds) > 0 {
		status, err = models.GetTrTransactionIn(&parentTrans, parentTransIds, "transaction_key")
		if err != nil {
			// log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
	}
	parentTransData := make(map[uint64]models.TrTransaction)
	for _, pt := range parentTrans {
		parentTransData[pt.TransactionKey] = pt

		if _, ok := lib.Find(productIds, strconv.FormatUint(pt.ProductKey, 10)); !ok {
			productIds = append(productIds, strconv.FormatUint(pt.ProductKey, 10))
		}
	}

	//mapping product
	var msProduct []models.MsProduct
	if len(productIds) > 0 {
		status, err = models.GetMsProductIn(&msProduct, productIds, "product_key")
		if err != nil {
			// log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
	}
	productData := make(map[uint64]models.MsProduct)
	for _, p := range msProduct {
		productData[p.ProductKey] = p
	}

	//mapping Trans type
	var transactionType []models.TrTransactionType
	if len(transTypeIds) > 0 {
		status, err = models.GetMsTransactionTypeIn(&transactionType, transTypeIds, "trans_type_key")
		if err != nil {
			// log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
	}
	transactionTypeData := make(map[uint64]models.TrTransactionType)
	for _, t := range transactionType {
		transactionTypeData[t.TransTypeKey] = t
	}

	//mapping ms bank
	var msBank []models.MsBank
	if len(bankIds) > 0 {
		status, err = models.GetMsBankIn(&msBank, bankIds, "bank_key")
		if err != nil {
			// log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
	}
	bankData := make(map[uint64]models.MsBank)
	for _, b := range msBank {
		bankData[b.BankKey] = b
	}

	//mapping trans status
	var trTransactionStatus []models.TrTransactionStatus
	if len(transStatusKey) > 0 {
		status, err = models.GetMsTransactionStatusIn(&trTransactionStatus, transStatusKey, "trans_status_key")
		if err != nil {
			// log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
	}
	transStatusData := make(map[uint64]models.TrTransactionStatus)
	for _, ts := range trTransactionStatus {
		transStatusData[ts.TransStatusKey] = ts
	}

	var responseData []models.AdminTrTransactionList
	for _, tr := range trTransaction {
		var data models.AdminTrTransactionList

		data.TransactionKey = tr.TransactionKey

		if tr.BranchKey != nil {
			if n, ok := branchData[*tr.BranchKey]; ok {
				data.BranchName = n.BranchName
			}
		}

		if tr.AgentKey != nil {
			if n, ok := agentData[*tr.AgentKey]; ok {
				data.AgentName = n.AgentName
			}
		}

		if n, ok := customerData[tr.CustomerKey]; ok {
			data.CustomerName = n.FullName
		}

		if n, ok := productData[tr.ProductKey]; ok {
			data.ProductName = n.ProductNameAlt
		}

		if n, ok := transStatusData[tr.TransStatusKey]; ok {
			data.TransStatus = *n.StatusCode
		}

		layout := "2006-01-02 15:04:05"
		newLayout := "02 Jan 2006"
		date, _ := time.Parse(layout, tr.TransDate)
		data.TransDate = date.Format(newLayout)
		date, _ = time.Parse(layout, tr.NavDate)
		data.NavDate = date.Format(newLayout)

		if n, ok := transactionTypeData[tr.TransTypeKey]; ok {
			data.TransType = *n.TypeDescription
		}

		data.TransAmount = tr.TransAmount
		data.TransUnit = tr.TransUnit
		data.TotalAmount = tr.TotalAmount

		if tr.TransBankKey != nil {
			if n, ok := bankData[*tr.TransBankKey]; ok {
				data.TransBankName = n.BankName
			}
		}

		data.TransBankAccNo = tr.TransBankAccNo
		data.TransBankaccName = tr.TransBankaccName

		strTransTypeKey := strconv.FormatUint(tr.TransTypeKey, 10)

		if strTransTypeKey == "4" {
			data.ProductIn = &data.ProductName
			if tr.ParentKey != nil {
				if n, ok := parentTransData[*tr.ParentKey]; ok {
					if pd, ok := productData[n.ProductKey]; ok {
						data.ProductOut = &pd.ProductName
					}
				}
			}
		}

		if strTransTypeKey == "1" { // subs/topup
			kosong := ""
			data.PaymentMethod = &kosong
			data.PaymentChannel = &kosong
			if tr.PaymentMethod != nil {
				if n, ok := gData[*tr.PaymentMethod]; ok {
					data.PaymentMethod = n.LkpName
				}
			}
			if sd, ok := settleData[tr.TransactionKey]; ok {
				if pc, ok := channelData[sd.SettlePaymentMethod]; ok {
					data.PaymentChannel = pc.PchannelName
				} else {
					// log.Println("no channel")
				}
			} else {
				// log.Println("no settle")
			}
		}

		if tr.TransSource != nil {
			if n, ok := gData[*tr.TransSource]; ok {
				data.TransSource = n.LkpName
			}
		}

		data.RecImage1 = nil

		if tr.TransTypeKey == uint64(1) && tr.PaymentMethod != nil {
			if *tr.PaymentMethod == uint64(284) { //transfer manual
				if tr.RecImage1 != nil {
					if n, ok := userLoginData[tr.CustomerKey]; ok {
						path := config.ImageUrl + "/images/user/" + strconv.FormatUint(n.UserLoginKey, 10) + "/transfer/" + *tr.RecImage1
						data.RecImage1 = &path
					}
				}
			}
		}

		responseData = append(responseData, data)
	}

	var countData models.CountData
	var pagination int
	if limit > 0 {
		status, err = models.AdminGetCountTrTransaction(&countData, params, transStatusKey, "trans_status_key", strconv.FormatUint(lib.Profile.UserID, 10))
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

func GetTransactionDetail(c echo.Context) error {
	var err error
	var status int
	decimal.MarshalJSONWithoutQuotes = true

	keyStr := c.Param("key")
	key, _ := strconv.ParseUint(keyStr, 10, 64)
	if key == 0 {
		return lib.CustomError(http.StatusNotFound)
	}

	var transaction models.TrTransaction
	status, err = models.GetTrTransaction(&transaction, keyStr)
	if err != nil {
		return lib.CustomError(status)
	}

	strTransStatusKey := strconv.FormatUint(transaction.TransStatusKey, 10)

	var roleKeyCs uint64
	roleKeyCs = 11
	var roleKeyKyc uint64
	roleKeyKyc = 12
	var roleKeyFundAdmin uint64
	roleKeyFundAdmin = 13

	if lib.Profile.RoleKey == roleKeyCs {
		statusCs := strconv.FormatUint(uint64(2), 10)
		if statusCs != strTransStatusKey {
			// log.Error("User Autorizer")
			return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
		}
	}
	if lib.Profile.RoleKey == roleKeyKyc {
		statusKyc := strconv.FormatUint(uint64(4), 10)
		if statusKyc != strTransStatusKey {
			// log.Error("User Autorizer")
			return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
		}
	}
	if lib.Profile.RoleKey == roleKeyFundAdmin {
		status := []string{"5", "6", "7", "8", "9"}
		_, found := lib.Find(status, strTransStatusKey)
		if !found {
			// log.Error("User Autorizer")
			return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
		}
	}

	//if user category  = 3 -> user branch, 2 = user HO
	var userCategory uint64
	userCategory = 3
	if lib.Profile.UserCategoryKey == userCategory {
		var cus models.MsCustomer
		strCusKey := strconv.FormatUint(transaction.CustomerKey, 10)
		status, err = models.GetMsCustomer(&cus, strCusKey)
		if err != nil {
			// log.Error(err.Error())
			return lib.CustomError(http.StatusBadRequest, err.Error(), "Customer tidak ditemukan")
		} else {
			if cus.OpenaccBranchKey == nil {
				// log.Error("Customer Branch null, not match with user branch")
				return lib.CustomError(http.StatusNotFound)
			} else {
				strCusBranch := strconv.FormatUint(*cus.OpenaccBranchKey, 10)
				strUserBranch := strconv.FormatUint(*lib.Profile.BranchKey, 10)
				// log.Error("Customer branch " + strCusBranch)
				// log.Error("User branch " + strUserBranch)
				if strCusBranch != strUserBranch {
					// log.Error("User Branch not match with customer branch")
					return lib.CustomError(http.StatusNotFound)
				}
			}
		}

	}

	var responseData models.AdminTransactionDetail

	var lookupIds []string
	var pChannelIds []string

	var trSettlement []models.TrTransactionSettlement
	paramSettlement := make(map[string]string)
	paramSettlement["rec_status"] = "1"
	paramSettlement["transaction_key"] = strconv.FormatUint(transaction.TransactionKey, 10)
	status, err = models.GetAllTrTransactionSettlement(&trSettlement, paramSettlement)
	if err == nil {
		if len(trSettlement) > 0 {
			for _, settlement := range trSettlement {
				if _, ok := lib.Find(lookupIds, strconv.FormatUint(settlement.SettlePurposed, 10)); !ok {
					lookupIds = append(lookupIds, strconv.FormatUint(settlement.SettlePurposed, 10))
				}

				if _, ok := lib.Find(lookupIds, strconv.FormatUint(settlement.SettleStatus, 10)); !ok {
					lookupIds = append(lookupIds, strconv.FormatUint(settlement.SettleStatus, 10))
				}

				if _, ok := lib.Find(lookupIds, strconv.FormatUint(settlement.SettleChannel, 10)); !ok {
					lookupIds = append(lookupIds, strconv.FormatUint(settlement.SettleChannel, 10))
				}

				if _, ok := lib.Find(pChannelIds, strconv.FormatUint(settlement.SettlePaymentMethod, 10)); !ok {
					pChannelIds = append(pChannelIds, strconv.FormatUint(settlement.SettlePaymentMethod, 10))
				}
			}
		}
	}

	if transaction.TrxCode != nil {
		if _, ok := lib.Find(lookupIds, strconv.FormatUint(*transaction.TrxCode, 10)); !ok {
			lookupIds = append(lookupIds, strconv.FormatUint(*transaction.TrxCode, 10))
		}
	}
	if transaction.EntryMode != nil {
		if _, ok := lib.Find(lookupIds, strconv.FormatUint(*transaction.EntryMode, 10)); !ok {
			lookupIds = append(lookupIds, strconv.FormatUint(*transaction.EntryMode, 10))
		}
	}
	if transaction.PaymentMethod != nil {
		if _, ok := lib.Find(lookupIds, strconv.FormatUint(*transaction.PaymentMethod, 10)); !ok {
			lookupIds = append(lookupIds, strconv.FormatUint(*transaction.PaymentMethod, 10))
		}
	}
	if transaction.TrxRiskLevel != nil {
		if _, ok := lib.Find(lookupIds, strconv.FormatUint(*transaction.TrxRiskLevel, 10)); !ok {
			lookupIds = append(lookupIds, strconv.FormatUint(*transaction.TrxRiskLevel, 10))
		}
	}
	if transaction.TransSource != nil {
		if _, ok := lib.Find(lookupIds, strconv.FormatUint(*transaction.TransSource, 10)); !ok {
			lookupIds = append(lookupIds, strconv.FormatUint(*transaction.TransSource, 10))
		}
	}

	//ms payment channel branch
	var pChannel []models.MsPaymentChannel
	if len(pChannelIds) > 0 {
		status, err = models.GetMsPaymentChannelIn(&pChannel, pChannelIds, "pchannel_key")
		if err != nil {
			// log.Error(err.Error())
		}
	}
	channelData := make(map[uint64]models.MsPaymentChannel)
	for _, pc := range pChannel {
		channelData[pc.PchannelKey] = pc
	}

	//gen lookup oa request
	var lookupOaReq []models.GenLookup
	if len(lookupIds) > 0 {
		status, err = models.GetGenLookupIn(&lookupOaReq, lookupIds, "lookup_key")
		if err != nil {
			if err != sql.ErrNoRows {
				// log.Error(err.Error())
				return lib.CustomError(status, err.Error(), "Failed get data")
			}
		}
	}

	layout := "2006-01-02 15:04:05"
	newLayout := "02 Jan 2006"

	gData := make(map[uint64]models.GenLookup)
	for _, gen := range lookupOaReq {
		gData[gen.LookupKey] = gen
	}

	//set settlement
	var settlementTransactionList []models.TransactionSettlement
	if len(trSettlement) > 0 {
		for _, settlement := range trSettlement {
			var data models.TransactionSettlement
			data.SettlementKey = settlement.SettlementKey
			date, _ := time.Parse(layout, settlement.SettleDate)
			data.SettleDate = date.Format(newLayout)
			data.SettleNominal = settlement.SettleNominal
			if settlement.SettleRealizedDate != nil {
				date, _ = time.Parse(layout, *settlement.SettleRealizedDate)
				data.SettleRealizedDate = date.Format(newLayout)
			} else {
				data.SettleRealizedDate = ""
			}
			data.SettleRemarks = settlement.SettleRemarks
			data.SettleReference = settlement.SettleReference
			if n, ok := gData[settlement.SettlePurposed]; ok {
				data.SettlePurposed = *n.LkpName
			}
			if n, ok := gData[settlement.SettleStatus]; ok {
				data.SettleStatus = *n.LkpName
			}
			if n, ok := gData[settlement.SettleChannel]; ok {
				data.SettleChannel = *n.LkpName
			}

			data.SettlePaymentMethod = ""
			if n, ok := channelData[settlement.SettlePaymentMethod]; ok {
				data.SettlePaymentMethod = *n.PchannelName
			}

			settlementTransactionList = append(settlementTransactionList, data)
		}
	}
	responseData.TransactionSettlement = &settlementTransactionList

	if transaction.TrxCode != nil {
		if n, ok := gData[*transaction.TrxCode]; ok {
			var trc models.LookupTrans

			trc.LookupKey = n.LookupKey
			trc.LkpGroupKey = n.LkpGroupKey
			trc.LkpCode = n.LkpCode
			trc.LkpName = n.LkpName
			responseData.TrxCode = &trc
		}
	}

	if transaction.EntryMode != nil {
		if n, ok := gData[*transaction.EntryMode]; ok {
			var entm models.LookupTrans

			entm.LookupKey = n.LookupKey
			entm.LkpGroupKey = n.LkpGroupKey
			entm.LkpCode = n.LkpCode
			entm.LkpName = n.LkpName
			responseData.EntryMode = &entm
		}
	}

	if transaction.TransSource != nil {
		if n, ok := gData[*transaction.TransSource]; ok {
			responseData.TransSource = n.LkpName
		}
	}

	if transaction.PaymentMethod != nil {
		if n, ok := gData[*transaction.PaymentMethod]; ok {
			var pm models.LookupTrans
			pm.LookupKey = n.LookupKey
			pm.LkpGroupKey = n.LkpGroupKey
			pm.LkpCode = n.LkpCode
			pm.LkpName = n.LkpName
			responseData.PaymentMethod = &pm
		}
	}

	if transaction.TrxRiskLevel != nil {
		if n, ok := gData[*transaction.TrxRiskLevel]; ok {
			var risk models.LookupTrans

			risk.LookupKey = n.LookupKey
			risk.LkpGroupKey = n.LkpGroupKey
			risk.LkpCode = n.LkpCode
			risk.LkpName = n.LkpName
			responseData.TrxRiskLevel = &risk
		}
	}

	responseData.TransactionKey = transaction.TransactionKey
	date, _ := time.Parse(layout, transaction.TransDate)
	responseData.TransDate = date.Format(newLayout)
	date, _ = time.Parse(layout, transaction.NavDate)
	responseData.NavDate = date.Format(newLayout)
	if transaction.RecCreatedDate != nil {
		date, err = time.Parse(layout, *transaction.RecCreatedDate)
		if err == nil {
			oke := date.Format(newLayout)
			responseData.RecCreatedDate = &oke
		}
	}
	responseData.RecCreatedBy = transaction.RecCreatedBy
	responseData.TransAmount = transaction.TransAmount
	responseData.TransUnit = transaction.TransUnit
	responseData.TransUnitPercent = transaction.TransUnitPercent
	if transaction.FlagRedemtAll != nil {
		if int(*transaction.FlagRedemtAll) > 0 {
			responseData.FlagRedemtAll = true
		}
	}
	if transaction.FlagNewSub != nil {
		if int(*transaction.FlagNewSub) > 0 {
			responseData.FlagNewSub = true
		}
	}
	responseData.TransFeePercent = transaction.TransFeePercent
	responseData.TransFeeAmount = transaction.TransFeeAmount
	responseData.ChargesFeeAmount = transaction.ChargesFeeAmount
	responseData.ServicesFeeAmount = transaction.ServicesFeeAmount
	responseData.TotalAmount = transaction.TotalAmount
	responseData.SettlementDate = transaction.SettlementDate
	responseData.TransBankAccNo = transaction.TransBankAccNo
	responseData.TransBankaccName = transaction.TransBankaccName
	responseData.TransRemarks = transaction.TransRemarks
	responseData.TransReferences = transaction.TransReferences
	responseData.PromoCode = transaction.PromoCode
	responseData.SalesCode = transaction.SalesCode
	if transaction.RiskWaiver > 0 {
		responseData.RiskWaiver = true
	}
	// responseData.FileUploadDate = transaction.FileUploadDate
	responseData.ProceedDate = transaction.ProceedDate
	responseData.ProceedAmount = transaction.ProceedAmount
	responseData.SentDate = transaction.SentDate
	responseData.SentReferences = transaction.SentReferences
	responseData.ConfirmedDate = transaction.ConfirmedDate
	responseData.PostedDate = transaction.PostedDate
	responseData.PostedUnits = transaction.PostedUnits
	responseData.SettledDate = transaction.SettledDate
	responseData.StampFeeAmount = transaction.StampFeeAmount

	strCustomer := strconv.FormatUint(transaction.CustomerKey, 10)

	dir := ""
	if transaction.RecApprovalStage == nil {
		var userData models.ScUserLogin
		status, err = models.GetScUserLoginByCustomerKey(&userData, strCustomer)
		if err != nil {
			if err != sql.ErrNoRows {
				return lib.CustomError(status)
			}
		}
		dir = config.ImageUrl + "/images/user/" + strconv.FormatUint(userData.UserLoginKey, 10) + "/transfer/"
	} else {
		dir = config.ImageUrl + "/images/user/institusi/" + strCustomer + "/transfer/"
	}

	if transaction.RecImage1 != nil {
		path := dir + *transaction.RecImage1
		responseData.RecImage1 = &path
	}

	if transaction.BranchKey != nil {
		var branch models.MsBranch
		strBranch := strconv.FormatUint(*transaction.BranchKey, 10)
		status, err = models.GetMsBranch(&branch, strBranch)
		if err != nil {
			if err != sql.ErrNoRows {
				return lib.CustomError(status)
			}
		} else {
			var br models.BranchTrans
			br.BranchKey = branch.BranchKey
			br.BranchCode = branch.BranchCode
			br.BranchName = branch.BranchName
			responseData.Branch = &br
		}
	}

	//check agent
	if transaction.AgentKey != nil {
		var agent models.MsAgent
		strAgent := strconv.FormatUint(*transaction.AgentKey, 10)
		status, err = models.GetMsAgent(&agent, strAgent)
		if err != nil {
			if err != sql.ErrNoRows {
				return lib.CustomError(status)
			}
		} else {
			var a models.AgentTrans
			a.AgentKey = agent.AgentKey
			a.AgentCode = agent.AgentCode
			a.AgentName = agent.AgentName
			responseData.Agent = &a
		}
	}

	//check customer
	var customer models.MsCustomer
	strCus := strconv.FormatUint(transaction.CustomerKey, 10)
	status, err = models.GetMsCustomer(&customer, strCus)
	if err != nil {
		if err != sql.ErrNoRows {
			return lib.CustomError(status)
		}
	} else {
		responseData.Customer.CustomerKey = customer.CustomerKey
		responseData.Customer.FullName = customer.FullName
		responseData.Customer.SidNo = customer.SidNo
		responseData.Customer.UnitHolderIDno = customer.UnitHolderIDno
	}

	//check product
	var product models.MsProduct
	strPro := strconv.FormatUint(transaction.ProductKey, 10)
	status, err = models.GetMsProduct(&product, strPro)
	if err != nil {
		if err != sql.ErrNoRows {
			return lib.CustomError(status)
		}
	} else {
		if product.FundTypeKey != nil {
			responseData.FundTypeKey = product.FundTypeKey
		}
		responseData.Product.ProductKey = product.ProductKey
		responseData.Product.ProductCode = product.ProductCode
		responseData.Product.ProductName = product.ProductName
	}

	//check trans status
	var transStatus models.TrTransactionStatus
	strTrSt := strconv.FormatUint(transaction.TransStatusKey, 10)
	status, err = models.GetTrTransactionStatus(&transStatus, strTrSt)
	if err != nil {
		if err != sql.ErrNoRows {
			return lib.CustomError(status)
		}
	} else {
		responseData.TransStatus.TransStatusKey = transStatus.TransStatusKey
		responseData.TransStatus.StatusCode = transStatus.StatusCode
		responseData.TransStatus.StatusDescription = transStatus.StatusDescription
	}
	responseData.TransCalcMethod = transaction.TransCalcMethod
	responseData.NavDateReal = transaction.NavDate

	//check trans type
	var transType models.TrTransactionType
	strTrTy := strconv.FormatUint(transaction.TransTypeKey, 10)
	status, err = models.GetMsTransactionType(&transType, strTrTy)
	if err != nil {
		if err != sql.ErrNoRows {
			return lib.CustomError(status)
		}
	} else {
		responseData.TransType.TransTypeKey = transType.TransTypeKey
		responseData.TransType.TypeCode = transType.TypeCode
		responseData.TransType.TypeDescription = transType.TypeDescription
	}

	//check bank
	var bank models.MsBank
	if transaction.TransBankKey != nil {
		strBank := strconv.FormatUint(*transaction.TransBankKey, 10)
		status, err = models.GetMsBank(&bank, strBank)
		if err != nil {
			if err != sql.ErrNoRows {
				return lib.CustomError(status)
			}
		} else {
			var tb models.TransBank
			tb.BankKey = bank.BankKey
			tb.BankCode = bank.BankCode
			tb.BankName = bank.BankName
			responseData.TransBank = &tb
		}
	}

	//check aca
	if transaction.AcaKey != nil {
		var aca models.TrAccountAgent
		strAca := strconv.FormatUint(*transaction.AcaKey, 10)
		status, err = models.GetTrAccountAgent(&aca, strAca)
		if err != nil {
			if err != sql.ErrNoRows {
				return lib.CustomError(status)
			}
		} else {
			var ac models.AcaTrans
			ac.AcaKey = aca.AcaKey
			ac.AccKey = aca.AccKey
			var agent models.MsAgent
			strAgent := strconv.FormatUint(aca.AgentKey, 10)
			status, err = models.GetMsAgent(&agent, strAgent)
			if err != nil {
				if err != sql.ErrNoRows {
					return lib.CustomError(status)
				}
			} else {
				ac.AgentKey = agent.AgentKey
				ac.AgentCode = agent.AgentCode
				ac.AgentName = agent.AgentName
			}

			responseData.Aca = &ac
		}
	}

	//check transaction confirmation
	strTrKey := strconv.FormatUint(transaction.TransactionKey, 10)
	var tc models.TrTransactionConfirmation
	status, err = models.GetTrTransactionConfirmationByTransactionKey(&tc, strTrKey)
	if err != nil {
		if err != sql.ErrNoRows {
			return lib.CustomError(status)
		}
	} else {
		var transTc models.TransactionConfirmation
		transTc.TcKey = tc.TcKey
		transTc.ConfirmDate = tc.ConfirmDate
		transTc.ConfirmedAmount = tc.ConfirmedAmount
		transTc.ConfirmedUnit = tc.ConfirmedUnit

		responseData.TransactionConfirmation = &transTc

		var transTcInfo models.TrTransactionConfirmationInfo
		transTcInfo.TcKey = tc.TcKey
		date, _ := time.Parse(layout, tc.ConfirmDate)
		transTcInfo.ConfirmDate = date.Format(newLayout)
		transTcInfo.ConfirmedAmount = tc.ConfirmedAmount
		transTcInfo.ConfirmedUnit = tc.ConfirmedUnit
		transTcInfo.ConfirmedAmountDiff = tc.ConfirmedAmountDiff
		transTcInfo.ConfirmedUnitDiff = tc.ConfirmedUnitDiff

		responseData.TransactionConfirmationInfo = &transTcInfo
	}

	//cek promo
	if transaction.PromoCode != nil {
		var promo models.TrPromoData
		status, err = models.AdminGetDetailTransactionPromo(&promo, strTrKey, *transaction.PromoCode)
		if err == nil {
			responseData.Promo = &promo
		}
	}

	//bank transaction customer
	var trBankAccount models.TrTransactionBankAccount
	status, err = models.GetTrTransactionBankAccountByField(&trBankAccount, strTrKey, "transaction_key")
	if err == nil {
		strCustBankAcc := strconv.FormatUint(trBankAccount.CustBankaccKey, 10)
		var trBankCust models.MsCustomerBankAccountInfo
		status, err = models.GetMsCustomerBankAccountTransactionByKey(&trBankCust, strCustBankAcc)
		if err == nil {
			responseData.CustomerBankAccount = &trBankCust
		}

		strProdBankAcc := strconv.FormatUint(trBankAccount.ProdBankaccKey, 10)
		var prodBankAccount models.MsProductBankAccountTransactionInfo
		status, err = models.GetMsProductBankAccountTransactionByKey(&prodBankAccount, strProdBankAcc)
		if err == nil {
			responseData.ProductBankAccount = &prodBankAccount
		}
	}

	responseData.IsEnableUnposting = false
	responseData.MessageEnableUnposting = ""

	if strTrSt == "9" {
		responseData.MessageEnableUnposting = "Transaksi tidak dapat di Un-posting karena bukan data terakhir dari customer dan produk yang sama."
		var transAfter models.TrTransaction
		status, err = models.CheckTrTransactionLastProductCustomer(&transAfter, strCustomer, strPro, keyStr)
		if err != nil {
			if err == sql.ErrNoRows {
				responseData.IsEnableUnposting = true
				responseData.MessageEnableUnposting = ""
			}
		}
	}

	prmGetFile := make(map[string]string)
	prmGetFile["ref_fk_domain"] = "tr_transaction"
	prmGetFile["ref_fk_key"] = strTrKey
	var fls []models.MsFileModels
	status, err = models.GetMsFileDataWithCondition(&fls, prmGetFile)
	if len(fls) > 0 {
		for _, fl := range fls {
			aa := config.ImageUrl + *fl.FilePath
			responseData.UrlUpload = append(responseData.UrlUpload, &aa)
		}
		responseData.FileUploadDate = &*fls[0].RecCreatedDate
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = responseData

	return c.JSON(http.StatusOK, response)
}

func TransactionApprovalCs(c echo.Context) error {
	errorAuth := initAuthCs()
	if errorAuth != nil {
		// log.Error("User Autorizer")
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}

	transStatusKeyDefault := "2"
	transStatusIds := []string{"1", "3", "4"}

	return ProsesApproval(transStatusKeyDefault, transStatusIds, c)
}

func TransactionApprovalCompliance(c echo.Context) error {
	errorAuth := initAuthKyc()
	if errorAuth != nil {
		// log.Error("User Autorizer")
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}

	transStatusKeyDefault := "4"
	transStatusIds := []string{"1", "3", "5"}

	return ProsesApproval(transStatusKeyDefault, transStatusIds, c)
}

func ProsesApproval(transStatusKeyDefault string, transStatusIds []string, c echo.Context) error {
	var err error
	var status int
	decimal.MarshalJSONWithoutQuotes = true

	var roleKeyCs uint64
	roleKeyCs = 11
	var roleKeyKyc uint64
	roleKeyKyc = 12

	params := make(map[string]string)

	transStatus := c.FormValue("trans_status_key")
	if transStatus == "" {
		// log.Error("Missing required parameter: trans_status_key")
		return lib.CustomError(http.StatusBadRequest)
	} else {
		_, found := lib.Find(transStatusIds, transStatus)
		if !found {
			// log.Error("Missing required parameter: trans_status_key")
			return lib.CustomError(http.StatusBadRequest)
		}
	}

	n, err := strconv.ParseUint(transStatus, 10, 64)
	if err == nil && n > 0 {
		params["trans_status_key"] = transStatus
	} else {
		// log.Error("Wrong input for parameter: trans_status_key")
		return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: trans_status_key", "Wrong input for parameter: trans_status_key")
	}

	notes := c.FormValue("notes")

	if (transStatus == "1") || (transStatus == "3") { //CORRECTED / DELETED
		if notes == "" {
			// log.Error("Missing required parameter notes: Notes tidak boleh kosong")
			return lib.CustomError(http.StatusBadRequest, "Notes tidak boleh kosong", "Notes tidak boleh kosong")
		}
	}

	if lib.Profile.RoleKey == roleKeyKyc {
		trxrisklevel := c.FormValue("trx_risk_level")
		if trxrisklevel == "" {
			// log.Error("Missing required parameter: trx_risk_level")
			return lib.CustomError(http.StatusBadRequest, "trx_risk_level can not be blank", "trx_risk_level can not be blank")
		} else {

			listLevelOption := []string{"114", "115"} //lookup group key 24
			_, found := lib.Find(listLevelOption, trxrisklevel)
			if !found {
				// log.Error("Missing required parameter: trx_risk_level")
				return lib.CustomError(http.StatusBadRequest)
			}
		}

		n, err := strconv.ParseUint(trxrisklevel, 10, 64)
		if err == nil && n > 0 {
			params["trx_risk_level"] = trxrisklevel
		} else {
			// log.Error("Wrong input for parameter: trx_risk_level")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: trx_risk_level", "Wrong input for parameter: trx_risk_level")
		}
	}

	transactionkey := c.FormValue("transaction_key")
	if transactionkey == "" {
		// log.Error("Missing required parameter: transaction_key")
		return lib.CustomError(http.StatusBadRequest)
	}

	n, err = strconv.ParseUint(transactionkey, 10, 64)
	if err == nil && n > 0 {
		params["transaction_key"] = transactionkey
	} else {
		// log.Error("Wrong input for parameter: transaction_key")
		return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: transaction_key", "Wrong input for parameter: transaction_key")
	}

	var transaction models.TrTransaction
	status, err = models.GetTrTransaction(&transaction, transactionkey)
	if err != nil {
		return lib.CustomError(status)
	}

	strTransStatusKey := strconv.FormatUint(transaction.TransStatusKey, 10)

	strTransTypeKey := strconv.FormatUint(transaction.TransTypeKey, 10)

	// if strTransTypeKey == "3" {
	// 	// log.Error("Transaction not found")
	// 	return lib.CustomError(http.StatusBadRequest)
	// }

	if transStatusKeyDefault != strTransStatusKey {
		// log.Error("User Autorizer")
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}

	dateLayout := "2006-01-02 15:04:05"
	strIDUserLogin := strconv.FormatUint(lib.Profile.UserID, 10)

	if lib.Profile.RoleKey == roleKeyCs {
		params["check1_notes"] = notes
		params["check1_date"] = time.Now().Format(dateLayout)
		params["check1_flag"] = "1"
		params["check1_references"] = strIDUserLogin
	}

	if lib.Profile.RoleKey == roleKeyKyc {
		params["check2_notes"] = notes
		params["check2_date"] = time.Now().Format(dateLayout)
		params["check2_flag"] = "1"
		params["check2_references"] = strIDUserLogin
	}

	params["rec_modified_by"] = strIDUserLogin
	params["rec_modified_date"] = time.Now().Format(dateLayout)

	_, err = models.UpdateTrTransaction(params)
	if err != nil {
		// log.Error("Error update tr transaction")
		return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed update data")
	}

	// if strTransTypeKey == "4" {
	// 	if transaction.ParentKey != nil {
	// 		strParentKey := strconv.FormatUint(*transaction.ParentKey, 10)
	// 		params["transaction_key"] = strParentKey
	// 		// log.Println(params)
	// 		_, err = models.UpdateTrTransaction(params)
	// 		if err != nil {
	// 			// log.Error("Error update tr transaction parent")
	// 			return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed update data")
	// 		}
	// 	}
	// }

	//send email to KYC / fund admin
	if (lib.Profile.RoleKey == roleKeyCs) || (lib.Profile.RoleKey == roleKeyKyc) { //cek user
		if (transStatus == "4") || (transStatus == "5") { //cek jika approve
			//check customer
			var customer models.MsCustomer
			strCus := strconv.FormatUint(transaction.CustomerKey, 10)
			status, err = models.GetMsCustomer(&customer, strCus)
			if err != nil {
				if err != sql.ErrNoRows {
					return lib.CustomError(status)
				}
			}

			//check product
			var product models.MsProduct
			strPro := strconv.FormatUint(transaction.ProductKey, 10)
			status, err = models.GetMsProduct(&product, strPro)
			if err != nil {
				if err != sql.ErrNoRows {
					return lib.CustomError(status)
				}
			}

			if lib.Profile.RoleKey == roleKeyCs {
				// SentEmailTransactionToBackOffice(strconv.FormatUint(transaction.TransactionKey, 10), "12")
			}
			if lib.Profile.RoleKey == roleKeyKyc {
				// SentEmailTransactionToBackOffice(strconv.FormatUint(transaction.TransactionKey, 10), "12")
			}
		}

	}
	// log.Info("Success update transaksi")

	//notif if reject
	if transStatus == "3" {

		var customer models.MsCustomer
		strCustomerKey := strconv.FormatUint(transaction.CustomerKey, 10)
		status, err = models.GetMsCustomer(&customer, strCustomerKey)
		if err == nil {
			if customer.InvestorType == "263" { //individu
				var userLogin models.ScUserLogin
				_, err := models.GetScUserLoginByCustomerKey(&userLogin, strCustomerKey)
				if err != nil {
					// log.Error(err.Error())
				}
				//create user message
				CreateNotifRejected(strCustomerKey, strIDUserLogin, notes, strTransTypeKey, transaction, userLogin)

				//notip email ke customer
				SendEmailRejected(strCustomerKey, strIDUserLogin, notes, strTransTypeKey, transaction, userLogin)

				//notip email ke sales
				SentEmailTransactionRejectToSales(transactionkey, notes)
			} else if customer.InvestorType == "264" { //institusi
				SentEmailTransactionInstitutionRejectBackOfficeToUserCcSales(transactionkey, strCustomerKey, notes)
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

func SendEmailRejected(strCustomerKey string, strIDUserLogin string,
	notes string, strTransTypeKey string, transaction models.TrTransaction,
	userLogin models.ScUserLogin) {
	decimal.MarshalJSONWithoutQuotes = true

	var subject string
	var tpl bytes.Buffer
	layout := "2006-01-02 15:04:05"
	newLayout := "02 Jan 2006"
	timeLayout := "15:04"

	productFrom := "-"

	var customer models.MsCustomer
	_, err := models.GetMsCustomer(&customer, strCustomerKey)
	if err != nil {
		if err != sql.ErrNoRows {
			// log.Error(err.Error())
		}
	}

	strProductKey := strconv.FormatUint(transaction.ProductKey, 10)

	var product models.MsProduct
	_, err = models.GetMsProduct(&product, strProductKey)
	if err != nil {
		// log.Error(err.Error())
	}

	var currencyDB models.MsCurrency
	_, err = models.GetMsCurrency(&currencyDB, strconv.FormatUint(*product.CurrencyKey, 10))
	if err != nil {
		// log.Error(err.Error())
	}

	date, _ := time.Parse(layout, transaction.TransDate)

	var transactionParent models.TrTransaction
	if strTransTypeKey == "4" { // SWITCH IN
		if transaction.ParentKey != nil {
			strTrParentKey := strconv.FormatUint(*transaction.ParentKey, 10)
			_, err := models.GetTrTransaction(&transactionParent, strTrParentKey)
			if err == nil {
				strProductParentKey := strconv.FormatUint(transactionParent.ProductKey, 10)

				var productparent models.MsProduct
				_, err = models.GetMsProduct(&productparent, strProductParentKey)
				if err != nil {
					productFrom = productparent.ProductNameAlt
				}
			}
		}
	}

	bankNameCus := "-"
	noRekCus := "-"
	namaRekCus := "-"
	cabangRek := "-"

	//bank account customer
	strTrKey := strconv.FormatUint(transaction.TransactionKey, 10)
	if strTransTypeKey == "2" { //redm
		var trBankAccount models.TrTransactionBankAccount
		_, err := models.GetTrTransactionBankAccountByField(&trBankAccount, strTrKey, "transaction_key")
		if err == nil {
			strCustBankAcc := strconv.FormatUint(trBankAccount.CustBankaccKey, 10)
			var trBankCust models.MsCustomerBankAccountInfo
			_, err = models.GetMsCustomerBankAccountTransactionByKey(&trBankCust, strCustBankAcc)
			if err == nil {
				bankNameCus = trBankCust.BankName
				noRekCus = trBankCust.AccountNo
				namaRekCus = trBankCust.AccountName
				if trBankCust.BranchName != nil {
					cabangRek = *trBankCust.BranchName
				}
			}
		}
	}

	transType := "Subscription"
	transTypeKecil := "subscription"
	if strTransTypeKey == "1" { // SUBS
		if transaction.FlagNewSub != nil {
			if *transaction.FlagNewSub == 0 {
				transType = "Top Up"
				transTypeKecil = "top up"
			}
		}
	}

	ac := accounting.Accounting{Symbol: "", Precision: 2, Thousand: ".", Decimal: ","}

	dataReplace := struct {
		FileUrl        string
		Name           string
		Cif            string
		Date           string
		Time           string
		ProductName    string
		Symbol         *string
		Amount         string
		Fee            string
		RedmUnit       string
		BankName       string
		NoRek          string
		NamaRek        string
		Cabang         string
		ProductFrom    string
		ProductTo      string
		TransType      string
		Notes          string
		TransTypeKecil string
	}{
		FileUrl:        config.ImageUrl + "/images/mail",
		Name:           customer.FullName,
		Cif:            customer.UnitHolderIDno,
		Date:           date.Format(newLayout),
		Time:           date.Format(timeLayout) + " WIB",
		ProductName:    product.ProductNameAlt,
		Symbol:         currencyDB.Symbol,
		Amount:         ac.FormatMoney(transaction.TransAmount),
		Fee:            ac.FormatMoney(transaction.TransFeeAmount),
		RedmUnit:       ac.FormatMoney(transaction.TransUnit),
		BankName:       bankNameCus,
		NoRek:          noRekCus,
		NamaRek:        namaRekCus,
		Cabang:         cabangRek,
		ProductFrom:    productFrom,
		ProductTo:      product.ProductNameAlt,
		TransType:      transType,
		Notes:          notes,
		TransTypeKecil: transTypeKecil}

	if strTransTypeKey == "1" { // SUBS
		if transaction.FlagNewSub != nil {
			if *transaction.FlagNewSub == 1 {
				subject = "[MotionFunds] Subscription Kamu Gagal"
			} else {
				subject = "[MotionFunds] Top Up Kamu telah Berhasil"
			}
		} else {
			subject = "[MotionFunds] Top Up Kamu Gagal"
		}

		t := template.New("email-subscription-rejected.html")

		t, err := t.ParseFiles(config.BasePath + "/mail/email-subscription-rejected.html")
		if err != nil {
			// log.Println(err)
		}

		if err := t.Execute(&tpl, dataReplace); err != nil {
			// log.Println(err)
		}
	}

	if strTransTypeKey == "2" { // REDM
		subject = "[MotionFunds] Redemption Kamu Gagal"
		t := template.New("email-redemption-rejected.html")

		t, err := t.ParseFiles(config.BasePath + "/mail/email-redemption-rejected.html")
		if err != nil {
			// log.Println(err)
		}

		if err := t.Execute(&tpl, dataReplace); err != nil {
			// log.Println(err)
		}

	}

	if strTransTypeKey == "4" { // SWITCH
		subject = "[MotionFunds] Switching Kamu Gagal"
		t := template.New("email-switching-rejected.html")

		t, err := t.ParseFiles(config.BasePath + "/mail/email-switching-rejected.html")
		if err != nil {
			// log.Println(err)
		}

		if err := t.Execute(&tpl, dataReplace); err != nil {
			// log.Println(err)
		}

	}

	// Send email
	result := tpl.String()

	mailer := gomail.NewMessage()
	// mailer.SetHeader("From", config.EmailFrom)
	mailer.SetHeader("To", userLogin.UloginEmail)
	mailer.SetHeader("Subject", subject)
	mailer.SetBody("text/html", result)

	err = lib.SendEmail(mailer)
	if err != nil {
		// log.Info("Email sent error")
		// log.Error(err)
	} else {
		// log.Info("Email sent")
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
	// 	// log.Info("Email sent error")
	// 	// log.Error(err)
	// } else {
	// 	// log.Info("Email sent")
	// }
}

func CreateNotifRejected(strCustomerKey string, strIDUserLogin string,
	notes string, strTransTypeKey string, transaction models.TrTransaction,
	userLogin models.ScUserLogin) {

	dateLayout := "2006-01-02 15:04:05"

	paramsUserMessage := make(map[string]string)
	paramsUserMessage["umessage_type"] = "245"

	strUserLoginKey := strconv.FormatUint(userLogin.UserLoginKey, 10)
	paramsUserMessage["umessage_recipient_key"] = strUserLoginKey
	paramsUserMessage["umessage_receipt_date"] = time.Now().Format(dateLayout)
	paramsUserMessage["flag_read"] = "0"
	paramsUserMessage["umessage_sender_key"] = strIDUserLogin
	paramsUserMessage["umessage_sent_date"] = time.Now().Format(dateLayout)
	paramsUserMessage["flag_sent"] = "1"
	var subject string
	if strTransTypeKey == "1" { // SUBS
		if transaction.FlagNewSub != nil {
			if *transaction.FlagNewSub == 1 {
				subject = "Subscription Tidak Dapat Diproses"
				paramsUserMessage["umessage_subject"] = "Subscription Tidak Dapat Diproses"
			} else {
				subject = "Top Up Tidak Dapat Diproses"
				paramsUserMessage["umessage_subject"] = "Top Up Tidak Dapat Diproses"
			}
		} else {
			subject = "Top Up Tidak Dapat Diproses"
			paramsUserMessage["umessage_subject"] = "Top Up Tidak Dapat Diproses"
		}
	}

	if strTransTypeKey == "2" { // REDM
		subject = "Redemption Tidak Dapat Diproses"
		paramsUserMessage["umessage_subject"] = "Redemption Tidak Dapat Diproses"
	}
	if strTransTypeKey == "4" { // SWITCH
		subject = "Switching Tidak Dapat Diproses"
		paramsUserMessage["umessage_subject"] = "Switching Tidak Dapat Diproses"
	}

	body := notes + " Silakan menghubungi Customer Service untuk informasi lebih lanjut."
	paramsUserMessage["umessage_body"] = body

	paramsUserMessage["umessage_category"] = "248"
	paramsUserMessage["flag_archieved"] = "0"
	paramsUserMessage["archieved_date"] = time.Now().Format(dateLayout)
	paramsUserMessage["rec_status"] = "1"
	paramsUserMessage["rec_created_date"] = time.Now().Format(dateLayout)
	paramsUserMessage["rec_created_by"] = strIDUserLogin

	_, err := models.CreateScUserMessage(paramsUserMessage)
	if err != nil {
		// log.Error(err.Error())
		// log.Error("Error create user message")
	} else {
		// log.Println("Success create user message")
	}
	lib.CreateNotifCustomerFromAdminByUserLoginKey(strUserLoginKey, subject, body, "TRANSACTION")
}

func UpdateNavDate(c echo.Context) error {
	errorAuth := initAuthFundAdmin()
	if errorAuth != nil {
		// log.Error("User Autorizer")
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}
	decimal.MarshalJSONWithoutQuotes = true

	var err error
	var status int

	params := make(map[string]string)

	//date
	postnavdate := c.FormValue("nav_date")
	if postnavdate == "" {
		// log.Error("Missing required parameter: nav_date")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: nav_date", "Missing required parameter: nav_date")
	}

	layoutISO := "2006-01-02"

	t, _ := time.Parse(layoutISO, postnavdate)
	t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
	w := lib.IsWeekend(t)
	if w {
		// log.Error("Missing required parameter: nav_date cann't Weekend")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: nav_date cann't Weekend", "Missing required parameter: nav_date cann't Weekend")
	}

	paramHoliday := make(map[string]string)
	paramHoliday["holiday_date"] = postnavdate

	var holiday []models.MsHoliday
	status, err = models.GetAllMsHoliday(&holiday, paramHoliday)
	if err != nil {
		if err != sql.ErrNoRows {
			// log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
	}
	if len(holiday) > 0 {
		// log.Error("nav_date is Bursa Holiday")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: nav_date is Bursa Holiday", "Missing required parameter: nav_date is Bursa Holiday")
	}

	params["nav_date"] = postnavdate

	//trans_fee_percent
	transfeepercent := c.FormValue("trans_fee_percent")
	if transfeepercent != "" {
		transfeepercentFloat, err := strconv.ParseFloat(transfeepercent, 64)
		if err == nil {
			if transfeepercentFloat < 0 {
				// log.Error("Wrong input for parameter: trans_fee_percent cann't negatif")
				return lib.CustomError(http.StatusBadRequest, "Missing required parameter: trans_fee_percent must cann't negatif", "Missing required parameter: trans_fee_percent cann't negatif")
			}
			params["trans_fee_percent"] = transfeepercent
		} else {
			// log.Error("Wrong input for parameter: trans_fee_percent number")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: trans_fee_percent must number", "Missing required parameter: trans_fee_percent number")
		}
	}

	//trans_fee_amount
	transfeeamount := c.FormValue("trans_fee_amount")
	if transfeeamount != "" {
		transfeeamountFloat, err := strconv.ParseFloat(transfeeamount, 64)
		if err == nil {
			if transfeeamountFloat < 0 {
				// log.Error("Wrong input for parameter: trans_fee_amount cann't negatif")
				return lib.CustomError(http.StatusBadRequest, "Missing required parameter: trans_fee_amount must cann't negatif", "Missing required parameter: trans_fee_amount cann't negatif")
			}
			params["trans_fee_amount"] = transfeeamount
		} else {
			// log.Error("Wrong input for parameter: trans_fee_amount number")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: trans_fee_amount must number", "Missing required parameter: trans_fee_amount number")
		}
	}

	//charges_fee_amount
	chargesfeeamount := c.FormValue("charges_fee_amount")
	if chargesfeeamount != "" {
		chargesfeeamountFloat, err := strconv.ParseFloat(chargesfeeamount, 64)
		if err == nil {
			if chargesfeeamountFloat < 0 {
				// log.Error("Wrong input for parameter: charges_fee_amount cann't negatif")
				return lib.CustomError(http.StatusBadRequest, "Missing required parameter: charges_fee_amount must cann't negatif", "Missing required parameter: charges_fee_amount cann't negatif")
			}
			params["charges_fee_amount"] = chargesfeeamount
		} else {
			// log.Error("Wrong input for parameter: charges_fee_amount number")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: charges_fee_amount must number", "Missing required parameter: charges_fee_amount number")
		}
	}

	//services_fee_amount
	servicesfeeamount := c.FormValue("services_fee_amount")
	if servicesfeeamount != "" {
		servicesfeeamountFloat, err := strconv.ParseFloat(servicesfeeamount, 64)
		if err == nil {
			if servicesfeeamountFloat < 0 {
				// log.Error("Wrong input for parameter: services_fee_amount cann't negatif")
				return lib.CustomError(http.StatusBadRequest, "Missing required parameter: services_fee_amount must cann't negatif", "Missing required parameter: services_fee_amount cann't negatif")
			}
			params["services_fee_amount"] = servicesfeeamount
		} else {
			// log.Error("Wrong input for parameter: services_fee_amount number")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: services_fee_amount must number", "Missing required parameter: services_fee_amount number")
		}
	}

	var productbankacckey string
	var customerbankacckey string

	//prod_bankacc_key
	prodbankacckey := c.FormValue("prod_bankacc_key")
	if prodbankacckey == "" {
		// log.Error("Missing required parameter: prod_bankacc_key cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: prod_bankacc_key cann't be blank", "Missing required parameter: prod_bankacc_key cann't be blank")
	}
	strprodbankacckey, err := strconv.ParseUint(prodbankacckey, 10, 64)
	if err == nil && strprodbankacckey > 0 {
		productbankacckey = prodbankacckey
	} else {
		// log.Error("Wrong input for parameter: prod_bankacc_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: prod_bankacc_key", "Missing required parameter: prod_bankacc_key")
	}

	//cust_bankacc_key
	custbankacckey := c.FormValue("cust_bankacc_key")
	if custbankacckey == "" {
		// log.Error("Missing required parameter: cust_bankacc_key cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: cust_bankacc_key cann't be blank", "Missing required parameter: cust_bankacc_key cann't be blank")
	}
	strcustbankacckey, err := strconv.ParseUint(custbankacckey, 10, 64)
	if err == nil && strcustbankacckey > 0 {
		customerbankacckey = custbankacckey
	} else {
		// log.Error("Wrong input for parameter: cust_bankacc_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: cust_bankacc_key", "Missing required parameter: cust_bankacc_key")
	}

	transactionkey := c.FormValue("transaction_key")
	if transactionkey == "" {
		// log.Error("Missing required parameter: transaction_key")
		return lib.CustomError(http.StatusBadRequest)
	}

	n, err := strconv.ParseUint(transactionkey, 10, 64)
	if err == nil && n > 0 {
		params["transaction_key"] = transactionkey
	} else {
		// log.Error("Wrong input for parameter: transaction_key")
		return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: transaction_key", "Wrong input for parameter: transaction_key")
	}

	var transaction models.TrTransaction
	status, err = models.GetTrTransaction(&transaction, transactionkey)
	if err != nil {
		return lib.CustomError(status)
	}

	strTransType := strconv.FormatUint(transaction.TransTypeKey, 10)

	strTransStatusKey := strconv.FormatUint(transaction.TransStatusKey, 10)

	strStatusCutOff := "5"

	if strTransStatusKey != strStatusCutOff {
		// log.Error("Data not found")
		return lib.CustomError(http.StatusUnauthorized, "Data not found", "Data not found")
	}

	dateLayout := "2006-01-02 15:04:05"
	strIDUserLogin := strconv.FormatUint(lib.Profile.UserID, 10)

	params["rec_modified_by"] = strIDUserLogin
	params["rec_modified_date"] = time.Now().Format(dateLayout)

	params["settlement_date"] = postnavdate

	//set settlement_date by settlement period product
	if strTransType == "2" { //REDM
		//check product
		var product models.MsProduct
		strPro := strconv.FormatUint(transaction.ProductKey, 10)
		status, err = models.GetMsProduct(&product, strPro)
		if err != nil {
			if err != sql.ErrNoRows {
				return lib.CustomError(status)
			}
		} else {
			layoutISO := "2006-01-02"
			t, _ := time.Parse(layoutISO, postnavdate)

			if (product.SettlementPeriod != nil) && (*product.SettlementPeriod > 0) {
				params["settlement_date"] = SettDate(t, int(*product.SettlementPeriod))
			}

		}
	}

	_, err = models.UpdateTrTransaction(params)
	if err != nil {
		// log.Error("Error update tr transaction")
		return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed update data")
	}

	if strTransType == "4" {
		if transaction.ParentKey != nil {
			strParentKey := strconv.FormatUint(*transaction.ParentKey, 10)
			params["transaction_key"] = strParentKey
			// log.Println(params)
			_, err = models.UpdateTrTransaction(params)
			if err != nil {
				// log.Error("Error update tr transaction parent")
				return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed update data")
			}
		}
	}

	//cek tr_transaction_bank_account by transaction_key
	var trBankAccount models.TrTransactionBankAccount
	status, err = models.GetTrTransactionBankAccountByField(&trBankAccount, transactionkey, "transaction_key")
	paramsTrBankAcc := make(map[string]string)
	paramsTrBankAcc["prod_bankacc_key"] = productbankacckey
	paramsTrBankAcc["cust_bankacc_key"] = customerbankacckey
	if err == nil { //exis -> update
		paramsTrBankAcc["rec_modified_by"] = strIDUserLogin
		paramsTrBankAcc["rec_modified_date"] = time.Now().Format(dateLayout)
		strTrBankAccKey := strconv.FormatUint(trBankAccount.TransBankaccKey, 10)
		_, err = models.UpdateTrTransactionBankAccount(paramsTrBankAcc, strTrBankAccKey, "trans_bankacc_key")
		if err != nil {
			// log.Error(err.Error())
			// log.Error("Error update tr transaction bank account parent")
		}

	} else { //null -> insert
		paramsTrBankAcc["transaction_key"] = transactionkey
		paramsTrBankAcc["rec_modified_by"] = strIDUserLogin
		paramsTrBankAcc["rec_modified_date"] = time.Now().Format(dateLayout)
		paramsTrBankAcc["rec_status"] = "1"
		_, err = models.CreateTrTransactionBankAccount(paramsTrBankAcc)
		if err != nil {
			// log.Error(err.Error())
			// log.Error("Error insert tr transaction bank account parent")
		}

	}

	// log.Info("Success update transaksi")

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = ""
	return c.JSON(http.StatusOK, response)
}

func TransactionApprovalCutOff(c echo.Context) error {
	errorAuth := initAuthFundAdmin()
	if errorAuth != nil {
		// log.Error("User Autorizer")
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}

	//list id
	transIds := c.FormValue("trans_ids")
	if transIds == "" {
		// log.Error("Missing required parameter: trans_ids")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: trans_ids", "Missing required parameter: trans_ids")
	}

	s := strings.Split(transIds, ",")

	var transParamIds []string

	for _, value := range s {
		is := strings.TrimSpace(value)
		if is != "" {
			if _, ok := lib.Find(transParamIds, is); !ok {
				transParamIds = append(transParamIds, is)
			}
		}
	}

	var transactionList []models.TrTransaction
	if len(transParamIds) > 0 {
		status, err := models.GetTrTransactionIn(&transactionList, transParamIds, "transaction_key")
		if err != nil {
			// log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
		if len(transParamIds) != len(transactionList) {
			// log.Error("Missing required parameter: trans_ids")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: Jumlah Data & Parameter berbeda", "Missing required parameter: Jumlah Data & Parameter berbeda")
		}
	} else {
		// log.Error("Missing required parameter: trans_ids")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: trans_ids", "Missing required parameter: trans_ids")
	}

	strStatusCutOff := "5"

	trBankAccDone := 0
	trBankAccNotDone := 0
	var transParamIdsValid []string

	for _, tr := range transactionList {
		strTransStatusKey := strconv.FormatUint(tr.TransStatusKey, 10)
		if strTransStatusKey != strStatusCutOff {
			// log.Error("Missing required parameter: trans_ids")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: trans_ids ", "Missing required parameter: trans_ids")
		}

		//bank account customer
		strTrKey := strconv.FormatUint(tr.TransactionKey, 10)
		var trBankAccount models.TrTransactionBankAccount
		_, err := models.GetTrTransactionBankAccountByField(&trBankAccount, strTrKey, "transaction_key")
		if err == nil {
			if _, ok := lib.Find(transParamIdsValid, strTrKey); !ok {
				transParamIdsValid = append(transParamIdsValid, strTrKey)
			}
			trBankAccDone++
		} else {
			trBankAccNotDone++
		}
	}

	paramsUpdate := make(map[string]string)

	paramsUpdate["trans_status_key"] = "6"
	dateLayout := "2006-01-02 15:04:05"
	paramsUpdate["rec_modified_date"] = time.Now().Format(dateLayout)
	strKey := strconv.FormatUint(lib.Profile.UserID, 10)
	paramsUpdate["rec_modified_by"] = strKey

	if len(transParamIdsValid) > 0 {
		_, err := models.UpdateTrTransactionByKeyIn(paramsUpdate, transParamIdsValid, "transaction_key")
		if err != nil {
			// log.Error("Error update oa request")
			return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed update data")
		}
	}

	if trBankAccNotDone > 0 {
		// log.Error("Ada Transaksi yang belum memiliki transaksi bank akun. Silakan update terlebih dahulu.")
		return lib.CustomError(http.StatusInternalServerError, "Ada Transaksi yang belum memiliki transaksi bank akun. Silakan update terlebih dahulu.", "Ada Transaksi yang belum memiliki transaksi bank akun. Silakan update terlebih dahulu.")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = ""
	return c.JSON(http.StatusOK, response)
}

func SettDate(t time.Time, period int) string {
	var dateNext string

	layoutISO := "2006-01-02"
	dateNext = t.Format(layoutISO)

	for i := 0; i < period; i++ {

		t, _ := time.Parse(layoutISO, dateNext)
		dateAfter := t.AddDate(0, 0, 1)
		t = time.Date(dateAfter.Year(), dateAfter.Month(), dateAfter.Day(), 0, 0, 0, 0, time.UTC)
		dateAfter = SkipWeekend(dateAfter)

		strDate := dateAfter.Format(layoutISO)
		dateNext = CheckHolidayBursa(strDate)
	}
	return dateNext
}

func SkipWeekend(t time.Time) time.Time {
	dateAfter := t
	t = t.UTC()

	switch t.Weekday() {
	case time.Saturday:
		dateAfter = t.AddDate(0, 0, 2)
	case time.Sunday:
		dateAfter = t.AddDate(0, 0, 1)
	}

	return dateAfter
}

func CheckHolidayBursa(date string) string {
	dateStr := date
	layoutISO := "2006-01-02"
	paramHoliday := make(map[string]string)
	paramHoliday["holiday_date"] = date

	var holiday []models.MsHoliday
	_, err := models.GetAllMsHoliday(&holiday, paramHoliday)
	if err != nil {
		if err == sql.ErrNoRows {
			t, _ := time.Parse(layoutISO, date)
			t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
			dateStr = date
		}
	}

	if len(holiday) > 0 {
		t, _ := time.Parse(layoutISO, dateStr)
		dateAfter := t.AddDate(0, 0, 1)
		strDate := dateAfter.Format(layoutISO)
		dateStr = CheckHolidayBursa(strDate)
	}

	w, _ := time.Parse(layoutISO, dateStr)
	w = time.Date(w.Year(), w.Month(), w.Day(), 0, 0, 0, 0, time.UTC)
	cek := lib.IsWeekend(w)
	if cek {
		dateSkip := SkipWeekend(w)

		dateStr = dateSkip.Format(layoutISO)
		dateStr = CheckHolidayBursa(dateStr)
	}
	return dateStr
}

func GetFormatExcelDownloadList(c echo.Context) error {
	errorAuth := initAuthFundAdmin()
	if errorAuth != nil {
		// log.Error("User Autorizer")
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}
	decimal.MarshalJSONWithoutQuotes = true

	var transStatusKey []string
	transStatusKey = append(transStatusKey, "7")

	var err error
	var status int

	params := make(map[string]string)

	//date
	postnavdate := c.FormValue("nav_date")
	if postnavdate == "" {
		// log.Error("Missing required parameter: nav_date")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: nav_date", "Missing required parameter: nav_date")
	}

	transactiontype := c.FormValue("transaction_type")
	if transactiontype == "" {
		// log.Error("Missing required parameter: transaction_type")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: transaction_type", "Missing required parameter: transaction_type")
	}

	rolestransactiontype := []string{"1", "2", "3", "4", "13"}
	_, found := lib.Find(rolestransactiontype, transactiontype)
	if !found {
		return lib.CustomError(http.StatusUnauthorized, "Missing parameter: transaction_type", "Missing parameter: transaction_type")
	}

	params["rec_status"] = "1"
	params["nav_date"] = postnavdate
	params["trans_type_key"] = transactiontype

	var trTransaction []models.TrTransaction
	status, err = models.AdminGetAllTrTransaction(&trTransaction, 0, 0, true, params, transStatusKey, "trans_status_key", true, strconv.FormatUint(lib.Profile.UserID, 10))
	if err != nil {
		// log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if len(trTransaction) < 1 {
		// log.Error("transaction not found")
		return lib.CustomError(http.StatusNotFound, "Transaction not found", "Transaction not found")
	}

	var transTypeIds []string
	var customerIds []string
	var productIds []string
	for _, tr := range trTransaction {
		if _, ok := lib.Find(customerIds, strconv.FormatUint(tr.CustomerKey, 10)); !ok {
			customerIds = append(customerIds, strconv.FormatUint(tr.CustomerKey, 10))
		}
		if _, ok := lib.Find(productIds, strconv.FormatUint(tr.ProductKey, 10)); !ok {
			productIds = append(productIds, strconv.FormatUint(tr.ProductKey, 10))
		}
		if _, ok := lib.Find(transTypeIds, strconv.FormatUint(tr.TransTypeKey, 10)); !ok {
			transTypeIds = append(transTypeIds, strconv.FormatUint(tr.TransTypeKey, 10))
		}
	}

	//mapping customer
	var msCustomer []models.MsCustomer
	if len(customerIds) > 0 {
		status, err = models.GetMsCustomerIn(&msCustomer, customerIds, "customer_key")
		if err != nil {
			// log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
	}
	customerData := make(map[uint64]models.MsCustomer)
	for _, c := range msCustomer {
		customerData[c.CustomerKey] = c
	}

	//mapping product
	var msProduct []models.MsProduct
	if len(productIds) > 0 {
		status, err = models.GetMsProductIn(&msProduct, productIds, "product_key")
		if err != nil {
			// log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
	}
	var custodianIds []string
	productData := make(map[uint64]models.MsProduct)
	for _, p := range msProduct {
		productData[p.ProductKey] = p

		if p.CustodianKey != nil {
			if _, ok := lib.Find(custodianIds, strconv.FormatUint(*p.CustodianKey, 10)); !ok {
				custodianIds = append(custodianIds, strconv.FormatUint(*p.CustodianKey, 10))
			}
		}
	}

	//mapping Trans type
	var transactionType []models.TrTransactionType
	if len(transTypeIds) > 0 {
		status, err = models.GetMsTransactionTypeIn(&transactionType, transTypeIds, "trans_type_key")
		if err != nil {
			// log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
	}
	transactionTypeData := make(map[uint64]models.TrTransactionType)
	for _, t := range transactionType {
		transactionTypeData[t.TransTypeKey] = t
	}

	//mapping ms custodian bank
	var custodianBank []models.MsCustodianBank
	if len(custodianIds) > 0 {
		status, err = models.GetMsCustodianBankIn(&custodianBank, custodianIds, "custodian_key")
		if err != nil {
			// log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
	}
	custodianBankData := make(map[uint64]models.MsCustodianBank)
	for _, cb := range custodianBank {
		custodianBankData[cb.CustodianKey] = cb
	}

	//mapping tr nav
	var trNav []models.TrNav
	if len(productIds) > 0 {
		status, err = models.GetAllTrNavBetween(&trNav, postnavdate, postnavdate, productIds)
		if err != nil {
			// log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
	}

	trNavData := make(map[uint64]models.TrNav)
	for _, tr := range trNav {
		trNavData[tr.ProductKey] = tr
	}

	var responseData []models.DownloadFormatExcelList
	for _, tr := range trTransaction {
		var data models.DownloadFormatExcelList

		data.IDTransaction = tr.TransactionKey
		if n, ok := transactionTypeData[tr.TransTypeKey]; ok {
			data.IDCategory = *n.TypeCode
		}
		if n, ok := productData[tr.ProductKey]; ok {
			data.ProductName = n.ProductName

			if n.CustodianKey != nil {
				if cb, ok := custodianBankData[*n.CustodianKey]; ok {
					data.Keterangan = cb.CustodianCode
				}
			}
		}

		if n, ok := customerData[tr.CustomerKey]; ok {
			data.FullName = n.FullName
		}

		layout := "2006-01-02 15:04:05"
		newLayout := "01/02/2006"
		date, _ := time.Parse(layout, tr.NavDate)
		data.NavDate = date.Format(newLayout)
		date, _ = time.Parse(layout, tr.TransDate)
		data.TransactionDate = date.Format(newLayout)

		data.Units = tr.TransUnit
		data.NetAmount = tr.TransAmount

		data.NavValue = nil
		if nv, ok := trNavData[tr.ProductKey]; ok {
			data.NavValue = &nv.NavValue
		} else {
			data.Keterangan = "NAV VALUE NOT EXIST"
		}
		data.ApproveUnits = tr.TransUnit
		data.ApproveAmount = tr.TransAmount
		data.Result = ""

		responseData = append(responseData, data)
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = responseData

	return c.JSON(http.StatusOK, response)
}

func UploadExcelConfirmation(c echo.Context) error {
	var err error
	decimal.MarshalJSONWithoutQuotes = true

	var responseData []models.DownloadFormatExcelList

	zero := decimal.NewFromInt(0)

	err = os.MkdirAll(config.BasePath+"/transaksi_upload/confirmation/", 0755)
	if err != nil {
		// log.Error(err.Error())
	} else {
		var file *multipart.FileHeader
		file, err = c.FormFile("excel")

		if file != nil {
			if err != nil {
				return lib.CustomError(http.StatusBadRequest)
			}
			// Get file extension
			extension := filepath.Ext(file.Filename)
			// log.Println(extension)
			roles := []string{".xlsx", ".XLSX"}
			_, found := lib.Find(roles, extension)
			if !found {
				return lib.CustomError(http.StatusUnauthorized, "Format file must .xlsx", "Format file must .xlsx")
			}

			// Generate filename
			//var filename string
			filename := lib.RandStringBytesMaskImprSrc(20)
			// log.Println("Generate filename:", filename)
			// Upload image and move to proper directory
			err = lib.UploadImage(file, config.BasePath+"/transaksi_upload/confirmation/"+filename+extension)
			if err != nil {
				// log.Println(err)
				return lib.CustomError(http.StatusInternalServerError)
			}

			xlsx, err := excelize.OpenFile(config.BasePath + "/transaksi_upload/confirmation/" + filename + extension)
			if err != nil {
				// log.Println(config.BasePath + "/transaksi_upload/confirmation/" + filename + extension)
				// // log.Fatal("ERROR", err.Error())
				return lib.CustomError(http.StatusInternalServerError)
			}

			sheet1Name := xlsx.GetSheetName(1)

			for i := 2; i < 1000; i++ {
				var data models.DownloadFormatExcelList

				iDTransaction := xlsx.GetCellValue(sheet1Name, fmt.Sprintf("A%d", i))
				iDCategory := xlsx.GetCellValue(sheet1Name, fmt.Sprintf("B%d", i))
				productName := xlsx.GetCellValue(sheet1Name, fmt.Sprintf("C%d", i))
				fullName := xlsx.GetCellValue(sheet1Name, fmt.Sprintf("D%d", i))
				navDate := xlsx.GetCellValue(sheet1Name, fmt.Sprintf("E%d", i))
				transactionDate := xlsx.GetCellValue(sheet1Name, fmt.Sprintf("F%d", i))
				units := xlsx.GetCellValue(sheet1Name, fmt.Sprintf("G%d", i))
				netAmount := xlsx.GetCellValue(sheet1Name, fmt.Sprintf("H%d", i))
				navValue := xlsx.GetCellValue(sheet1Name, fmt.Sprintf("I%d", i))
				approveUnits := xlsx.GetCellValue(sheet1Name, fmt.Sprintf("J%d", i))
				approveAmount := xlsx.GetCellValue(sheet1Name, fmt.Sprintf("K%d", i))
				keterangan := xlsx.GetCellValue(sheet1Name, fmt.Sprintf("L%d", i))
				result := xlsx.GetCellValue(sheet1Name, fmt.Sprintf("M%d", i))

				// log.Println(navDate)
				// log.Println(transactionDate)
				if iDTransaction == "" {
					break
				}

				key, _ := strconv.ParseUint(iDTransaction, 10, 64)
				if key == 0 {
					return lib.CustomError(http.StatusNotFound)
				}

				data.IDTransaction = key
				data.IDCategory = iDCategory
				data.ProductName = productName
				data.FullName = fullName
				data.NavDate = navDate
				data.TransactionDate = transactionDate

				if units != "" {
					if unitsFloat, err := decimal.NewFromString(units); err == nil {
						data.Units = unitsFloat
					}
				}

				if netAmount != "" {
					if netAmountFloat, err := decimal.NewFromString(netAmount); err == nil {
						data.NetAmount = netAmountFloat
					}
				}

				if navValue != "" {
					if navValueFloat, err := decimal.NewFromString(navValue); err == nil {
						nav := navValueFloat
						data.NavValue = &nav
					}
				}

				var transUnitFifo decimal.Decimal

				if approveUnits != "" {
					if approveUnitsFloat, err := decimal.NewFromString(approveUnits); err == nil {
						data.ApproveUnits = approveUnitsFloat
						transUnitFifo = data.ApproveUnits
					}
				}

				if approveAmount != "" {
					if approveAmountFloat, err := decimal.NewFromString(approveAmount); err == nil {
						data.ApproveAmount = approveAmountFloat
					}
				}

				data.Keterangan = keterangan
				data.Result = result

				//cek transaction
				var transaction models.TrTransaction
				_, err := models.GetTrTransaction(&transaction, iDTransaction)
				if err != nil {
					if err == sql.ErrNoRows {
						data.Result = "Data Transaction Not Found"
					} else {
						data.Result = err.Error()
					}
					fmt.Printf("%v \n", data)
					responseData = append(responseData, data)
					continue
				}

				//cek trans status
				strTransStatusKey := strconv.FormatUint(transaction.TransStatusKey, 10)
				if strTransStatusKey != "7" {
					data.Result = "Status Transaction Not in CONFIRMED"
					fmt.Printf("%v \n", data)
					responseData = append(responseData, data)
					continue
				}

				//cek transaction confirmation
				var transactionConf models.TrTransactionConfirmation
				_, err = models.GetTrTransactionConfirmationByTransactionKey(&transactionConf, iDTransaction)
				if err != nil {
					if err != sql.ErrNoRows {
						data.Result = err.Error()
						fmt.Printf("%v \n", data)
						responseData = append(responseData, data)
						continue
					}
				} else {
					data.Result = "TC already exists"
					fmt.Printf("%v \n", data)
					responseData = append(responseData, data)
					continue
				}

				strProductKey := strconv.FormatUint(transaction.ProductKey, 10)
				strCustomerKey := strconv.FormatUint(transaction.CustomerKey, 10)

				var trNav []models.TrNav
				_, err = models.GetTrNavByProductKeyAndNavDate(&trNav, strProductKey, transaction.NavDate)
				if err != nil {
					if err != sql.ErrNoRows {
						data.Result = "NAV VALUE NOT EXIST"
						fmt.Printf("%v \n", data)
						responseData = append(responseData, data)
						continue
					} else {
						data.Result = err.Error()
						fmt.Printf("%v \n", data)
						responseData = append(responseData, data)
						continue
					}
				}
				strTransTypeKey := strconv.FormatUint(transaction.TransTypeKey, 10)

				var trBalanceCustomer []models.TrBalanceCustomerProduk

				//redm cek balance / saldo aktif
				if (strTransTypeKey == "2") || (strTransTypeKey == "3") { // REDM
					_, err = models.GetLastBalanceCustomerByProductKey(&trBalanceCustomer, strCustomerKey, strProductKey)
					if err != nil {
						if err != sql.ErrNoRows {
							data.Result = "Balance is empty"
							fmt.Printf("%v \n", data)
							responseData = append(responseData, data)
							continue
						} else {
							data.Result = err.Error()
							fmt.Printf("%v \n", data)
							responseData = append(responseData, data)
							continue
						}
					}
				}

				//redm cek balance / saldo aktif di parent jika switch
				var transactionParent models.TrTransaction
				if strTransTypeKey == "4" { // SWITCH IN
					if transaction.ParentKey == nil {
						data.Result = "Parent Transaction is empty"
						fmt.Printf("%v \n", data)
						responseData = append(responseData, data)
						continue
					}

					strTrParentKey := strconv.FormatUint(*transaction.ParentKey, 10)
					_, err := models.GetTrTransaction(&transactionParent, strTrParentKey)
					if err != nil {
						if err == sql.ErrNoRows {
							data.Result = "Data Parent Transaction Not Found"
						} else {
							data.Result = err.Error()
						}
						fmt.Printf("%v \n", data)
						responseData = append(responseData, data)
						continue
					}
				}

				//data valid 1. create tr_transaction_confirmation, 2. update trans status, 3. create tr_transaction_fifo
				//1. create tr_transaction_confirmation
				dateLayout := "2006-01-02 15:04:05"
				params := make(map[string]string)
				params["transaction_key"] = iDTransaction
				params["confirmed_amount"] = approveAmount
				params["confirmed_unit"] = approveUnits
				params["confirm_result"] = "208"

				approveUnitsFloat := decimal.NewFromInt(0)
				if approveUnits != "" {
					if appUnits, err := decimal.NewFromString(approveUnits); err == nil {
						approveUnitsFloat = appUnits
					}
				}
				if transaction.TransUnit.Cmp(approveUnitsFloat) == 1 {
					// strTransUnit := fmt.Sprintf("%g", transaction.TransUnit.Sub(approveUnitsFloat))
					params["confirmed_unit_diff"] = transaction.TransUnit.Sub(approveUnitsFloat).String()
				} else if transaction.TransUnit.Cmp(approveUnitsFloat) == -1 {
					// strTransUnit := fmt.Sprintf("%g", approveUnitsFloat.Sub(transaction.TransUnit))
					params["confirmed_unit_diff"] = approveUnitsFloat.Sub(transaction.TransUnit).String()
				} else {
					params["confirmed_unit_diff"] = "0"
				}

				approveAmountFloat := decimal.NewFromInt(0)
				if approveUnits != "" {
					if appAmount, err := decimal.NewFromString(approveAmount); err == nil {
						approveAmountFloat = appAmount
					}
				}
				if transaction.TransAmount.Cmp(approveAmountFloat) == 1 {
					// strTransAmount := fmt.Sprintf("%g", transaction.TransAmount.Sub(approveAmountFloat))
					params["confirmed_amount_diff"] = transaction.TransAmount.Sub(approveAmountFloat).String()
				} else if transaction.TransAmount.Cmp(approveAmountFloat) == -1 {
					// strTransAmount := fmt.Sprintf("%g", approveAmountFloat.Sub(transaction.TransAmount))
					params["confirmed_amount_diff"] = approveAmountFloat.Sub(transaction.TransAmount).String()
				} else {
					params["confirmed_amount_diff"] = "0"
				}

				params["rec_status"] = "1"
				params["rec_created_date"] = time.Now().Format(dateLayout)
				params["rec_created_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

				if (strTransTypeKey == "2") || (strTransTypeKey == "3") { // REDM --> AVG_NAV ambil dari trBalance LAST
					var avgNav models.AvgNav
					_, err = models.GetLastAvgNavTrBalanceCustomerByProductKey(&avgNav, strCustomerKey, strProductKey)
					if err != nil {
						params["avg_nav"] = "0"
					} else {
						if avgNav.AvgNav != nil {
							//var avgAvg *decimal.Decimal
							avgAvg := avgNav.AvgNav
							params["avg_nav"] = avgAvg.String()
						} else {
							params["avg_nav"] = "0"
						}
						// strAvgNav := fmt.Sprintf("%g", *avgNav.AvgNav)
						// params["avg_nav"] = strAvgNav
					}
				} else { //SUBS
					if approveUnits != "" {
						if approveUnitsFloat, err := decimal.NewFromString(approveUnits); err == nil {
							avgNav := transaction.TotalAmount.Div(approveUnitsFloat)
							// strAvgNav := fmt.Sprintf("%g", avgNav)
							params["avg_nav"] = avgNav.String()
						}
					}
				}

				layout := "2006-01-02 15:04:05"
				newLayout := "2006-01-02"
				date, _ := time.Parse(layout, transaction.NavDate)

				layoutISO := "2006-01-02"
				t, _ := time.Parse(layoutISO, date.Format(newLayout))

				params["confirm_date"] = SettDate(t, int(1)) + " 00:00:00"

				status, err, trConfirmationID := models.CreateTrTransactionConfirmation(params)
				if err != nil {
					// log.Error(err.Error())
					return lib.CustomError(status, err.Error(), "Failed input data")
				}

				// 2. update trans status
				paramsTrans := make(map[string]string)
				paramsTrans["trans_status_key"] = "8"
				paramsTrans["confirmed_date"] = time.Now().Format(dateLayout)
				paramsTrans["rec_modified_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
				paramsTrans["rec_modified_date"] = time.Now().Format(dateLayout)
				paramsTrans["transaction_key"] = iDTransaction
				_, err = models.UpdateTrTransaction(paramsTrans)
				if err != nil {
					// log.Error("Error update tr transaction")
					return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed update data")
				}

				//3. create tr_transaction_fifo
				if strTransTypeKey == "1" ||
					strTransTypeKey == "13" ||
					strTransTypeKey == "4" { // SUB | TOPUP |  switchin
					paramsFifo := make(map[string]string)
					paramsFifo["trans_conf_sub_key"] = trConfirmationID
					if transaction.AcaKey != nil {
						strAcaKey := strconv.FormatUint(*transaction.AcaKey, 10)
						paramsFifo["sub_aca_key"] = strAcaKey
					}
					paramsFifo["holding_days"] = "0"
					paramsFifo["trans_unit"] = approveUnits
					paramsFifo["fee_nav_mode"] = "207"

					var transAmountFifo decimal.Decimal
					transAmountFifo = transUnitFifo.Mul(trNav[0].NavValue)
					// strTransAmountFifo := fmt.Sprintf("%g", transAmountFifo)
					paramsFifo["trans_amount"] = transAmountFifo.String()

					var feeTypeStr string

					if strTransTypeKey == "1" {
						feeTypeStr = "183"
					}
					if strTransTypeKey == "4" {
						feeTypeStr = "185"
					}

					var feeItem models.MsProductFeeItem

					// _, err = models.GetMsProductFeeItemCalculateFifoWithLimit(&feeItem, strProductKey, strTransAmountFifo, feeTypeStr)
					_, err = models.GetMsProductFeeItemCalculateFifoWithLimit(&feeItem, strProductKey, transAmountFifo.String(), feeTypeStr)
					if err != nil {
						if err == sql.ErrNoRows {
							_, err = models.GetMsProductFeeItemLastCalculateFifo(&feeItem, strProductKey, feeTypeStr)
							if err != nil {
								// log.Error(err.Error())
								paramsFifo["trans_fee_amount"] = "0"
								paramsFifo["trans_nett_amount"] = "0"
							} else {
								transfeeamount := feeItem.FeeValue.Div(decimal.NewFromInt(100)).Mul(transAmountFifo)
								// strTransfeeamount := fmt.Sprintf("%g", transfeeamount)
								paramsFifo["trans_fee_amount"] = transfeeamount.String()

								transnett := transAmountFifo.Add(transfeeamount)
								// strTransnett := fmt.Sprintf("%g", transnett)
								paramsFifo["trans_nett_amount"] = transnett.String()
							}
						} else {
							// log.Error(err.Error())
							paramsFifo["trans_fee_amount"] = "0"
							paramsFifo["trans_nett_amount"] = "0"
						}
					} else {
						transfeeamount := feeItem.FeeValue.Div(decimal.NewFromInt(100)).Mul(transAmountFifo)
						// strTransfeeamount := fmt.Sprintf("%g", transfeeamount)
						paramsFifo["trans_fee_amount"] = transfeeamount.String()

						transnett := transAmountFifo.Add(transfeeamount)
						// strTransnett := fmt.Sprintf("%g", transnett)
						paramsFifo["trans_nett_amount"] = transnett.String()
					}

					paramsFifo["trans_fee_tax"] = "0"
					paramsFifo["rec_status"] = "1"
					paramsFifo["rec_created_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
					paramsFifo["rec_created_date"] = time.Now().Format(dateLayout)
					_, err = models.CreateTrTransactionFifo(paramsFifo)
					if err != nil {
						// log.Error("Error update tr transaction")
						return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed insert data")
					}
				}

				if (strTransTypeKey == "2") || (strTransTypeKey == "3") { // REDM / switchout
					sisaFifo := transUnitFifo
					for _, trBalance := range trBalanceCustomer {
						if sisaFifo.Cmp(decimal.NewFromInt(0)) == 1 {
							var balanceUsed decimal.Decimal

							if trBalance.BalanceUnit.Cmp(sisaFifo) == 1 {
								balanceUsed = sisaFifo
								sisaFifo = zero
							}

							if trBalance.BalanceUnit.Cmp(sisaFifo) == -1 {
								balanceUsed = trBalance.BalanceUnit
								sisaFifo = sisaFifo.Sub(trBalance.BalanceUnit)
							}

							if trBalance.BalanceUnit == sisaFifo {
								balanceUsed = trBalance.BalanceUnit
								sisaFifo = zero
							}

							paramsFifo := make(map[string]string)
							paramsFifo["trans_conf_red_key"] = trConfirmationID
							strTcKeySub := strconv.FormatUint(trBalance.TcKey, 10)
							paramsFifo["trans_conf_sub_key"] = strTcKeySub
							if transaction.AcaKey != nil {
								strAcaKey := strconv.FormatUint(trBalance.AcaKey, 10)
								paramsFifo["sub_aca_key"] = strAcaKey
							}

							day1, _ := time.Parse(dateLayout, transaction.NavDate)
							day2, _ := time.Parse(dateLayout, trBalance.NavDate)

							days := day1.Sub(day2).Hours() / 24
							strDays := fmt.Sprintf("%g", days)

							paramsFifo["holding_days"] = strDays

							// strUnitUsed := fmt.Sprintf("%g", balanceUsed)
							paramsFifo["trans_unit"] = balanceUsed.String()
							paramsFifo["fee_nav_mode"] = "207"

							var transAmountFifo decimal.Decimal
							transAmountFifo = transUnitFifo.Mul(trNav[0].NavValue)
							// strTransAmountFifo := fmt.Sprintf("%g", transAmountFifo)
							paramsFifo["trans_amount"] = transAmountFifo.String()

							var feeTypeStr string

							if strTransTypeKey == "2" {
								feeTypeStr = "184"
							}
							if strTransTypeKey == "3" {
								feeTypeStr = "185"
							}

							var feeItem models.MsProductFeeItem

							// _, err = models.GetMsProductFeeItemCalculateFifoWithLimit(&feeItem, strProductKey, strTransAmountFifo, feeTypeStr)
							_, err = models.GetMsProductFeeItemCalculateFifoWithLimit(&feeItem, strProductKey, transAmountFifo.String(), feeTypeStr)
							if err != nil {
								if err == sql.ErrNoRows {
									_, err = models.GetMsProductFeeItemLastCalculateFifo(&feeItem, strProductKey, feeTypeStr)
									if err != nil {
										// log.Error(err.Error())
										paramsFifo["trans_fee_amount"] = "0"
										paramsFifo["trans_nett_amount"] = "0"
									} else {
										transfeeamount := feeItem.FeeValue.Div(decimal.NewFromInt(100)).Mul(transAmountFifo)
										// strTransfeeamount := fmt.Sprintf("%g", transfeeamount)
										paramsFifo["trans_fee_amount"] = transfeeamount.String()

										transnett := transAmountFifo.Add(transfeeamount)
										// strTransnett := fmt.Sprintf("%g", transnett)
										paramsFifo["trans_nett_amount"] = transnett.String()
									}
								} else {
									// log.Error(err.Error())
									paramsFifo["trans_fee_amount"] = "0"
									paramsFifo["trans_nett_amount"] = "0"
								}
							} else {
								transfeeamount := feeItem.FeeValue.Div(decimal.NewFromInt(100)).Mul(transAmountFifo)
								// strTransfeeamount := fmt.Sprintf("%g", transfeeamount)
								paramsFifo["trans_fee_amount"] = transfeeamount.String()

								transnett := transAmountFifo.Sub(transfeeamount)
								// strTransnett := fmt.Sprintf("%g", transnett)
								paramsFifo["trans_nett_amount"] = transnett.String()
							}

							paramsFifo["trans_fee_tax"] = "0"
							paramsFifo["rec_status"] = "1"
							paramsFifo["rec_created_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
							paramsFifo["rec_created_date"] = time.Now().Format(dateLayout)
							_, err = models.CreateTrTransactionFifo(paramsFifo)
							if err != nil {
								// log.Error("Error update tr transaction")
								return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed insert data")
							}
						} else {
							break
						}
					}
				}

				data.Keterangan = ""
				data.Result = "SUCCESS"

				responseData = append(responseData, data)
			}
		} else {
			// log.Error("File cann't be blank")
			return lib.CustomError(http.StatusNotFound, "File can not be blank", "File can not be blank")
		}
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = responseData
	return c.JSON(http.StatusOK, response)

}

func ProsesPosting(c echo.Context) error {
	errorAuth := initAuthFundAdmin()

	zero := decimal.NewFromInt(0)
	if errorAuth != nil {
		// log.Error("User Autorizer")
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}

	var err error
	var status int

	params := make(map[string]string)

	transactionkey := c.FormValue("transaction_key")
	if transactionkey == "" {
		// log.Error("Missing required parameter: transaction_key")
		return lib.CustomError(http.StatusBadRequest)
	}

	n, err := strconv.ParseUint(transactionkey, 10, 64)
	if err == nil && n > 0 {
		params["transaction_key"] = transactionkey
	} else {
		// log.Error("Wrong input for parameter: transaction_key")
		return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: transaction_key", "Wrong input for parameter: transaction_key")
	}

	var transaction models.TrTransaction
	status, err = models.GetTrTransaction(&transaction, transactionkey)
	if err != nil {
		return lib.CustomError(status)
	}

	strTransStatusKey := strconv.FormatUint(transaction.TransStatusKey, 10)

	if strTransStatusKey != "8" {
		// log.Error("User Autorizer")
		return lib.CustomError(http.StatusBadRequest, "Invalid transaction status. transaction_status="+strTransStatusKey, "Invalid transaction status")
	}

	dateLayout := "2006-01-02 15:04:05"
	strIDUserLogin := strconv.FormatUint(lib.Profile.UserID, 10)
	strTransTypeKey := strconv.FormatUint(transaction.TransTypeKey, 10)

	var transactionConf models.TrTransactionConfirmation
	strTransactionKey := strconv.FormatUint(transaction.TransactionKey, 10)
	_, err = models.GetTrTransactionConfirmationByTransactionKey(&transactionConf, strTransactionKey)
	if err != nil {
		// log.Error(err.Error())
		return lib.CustomError(http.StatusBadRequest, "GetTrTransactionConfirmationByTransactionKey:"+err.Error(), "Err GetTrTransactionConfirmationByTransactionKey")
	}

	var trBalanceCustomer []models.TrBalanceCustomerProduk
	strProductKey := strconv.FormatUint(transaction.ProductKey, 10)
	strCustomerKey := strconv.FormatUint(transaction.CustomerKey, 10)

	_, err = models.GetLastBalanceCustomerByProductKey(&trBalanceCustomer, strCustomerKey, strProductKey)

	if err != nil {
		if strTransTypeKey == "2" || strTransTypeKey == "3" { // REDM
			if err != sql.ErrNoRows {
				// log.Error("Transaction have not balance")
				return lib.CustomError(http.StatusBadRequest)
			} else {
				// log.Error(err.Error())
				return lib.CustomError(http.StatusBadRequest)
			}
		}
	}

	// strTransUnit := fmt.Sprintf("%g", transactionConf.ConfirmedUnit)

	//create tr_balance
	if strTransTypeKey == "1" || strTransTypeKey == "13" || strTransTypeKey == "4" { // SUB & SWIN
		paramsBalance := make(map[string]string)
		strAcaKey := strconv.FormatUint(*transaction.AcaKey, 10)
		paramsBalance["aca_key"] = strAcaKey
		strTransactionConf := strconv.FormatUint(transactionConf.TcKey, 10)
		paramsBalance["tc_key"] = strTransactionConf

		newlayout := "2006-01-02"
		t, _ := time.Parse(dateLayout, transactionConf.ConfirmDate)
		balanceDate := t.Format(newlayout)

		paramsBalance["balance_date"] = balanceDate + " 00:00:00"
		paramsBalance["balance_unit"] = transactionConf.ConfirmedUnit.String()
		paramsBalance["rec_order"] = "0"
		paramsBalance["rec_status"] = "1"
		paramsBalance["rec_created_date"] = time.Now().Format(dateLayout)
		paramsBalance["rec_created_by"] = strIDUserLogin

		//calculate avg_nag tr_balance
		//sum balance unit

		balanceUnitSum := zero
		for _, trBalance := range trBalanceCustomer {
			balanceUnitSum = balanceUnitSum.Add(trBalance.BalanceUnit)
		}

		//avg nav balance last
		avgNavLast := zero
		var avgNavMod models.AvgNav
		_, err = models.GetLastAvgNavTrBalanceCustomerByProductKey(&avgNavMod, strCustomerKey, strProductKey)
		if err == nil {
			if avgNavMod.AvgNav != nil {
				avgNavLast = *avgNavMod.AvgNav
			}
		}

		variable1 := balanceUnitSum.Mul(avgNavLast)
		variable2 := transactionConf.ConfirmedUnit.Mul(*transactionConf.AvgNav)

		balanceUnitSumAll := balanceUnitSum.Add(transactionConf.ConfirmedUnit)

		newAvgNavValue := variable1.Add(variable2).Div(balanceUnitSumAll)
		// strAvgNav := fmt.Sprintf("%g", countAvgNavBalance)
		paramsBalance["avg_nav"] = newAvgNavValue.String()
		//end calculate avg_nag tr_balance
		log.Println("CreateTrBalance")
		status, err := models.CreateTrBalance(paramsBalance)
		if err != nil {
			// log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed input data")
		}
	}

	if (strTransTypeKey == "2") || (strTransTypeKey == "3") { // REDM & SWOUT
		sisaFifo := transactionConf.ConfirmedUnit
		//avg nav balance last
		avgNavLast := zero
		var avgNav models.AvgNav
		_, err = models.GetLastAvgNavTrBalanceCustomerByProductKey(&avgNav, strCustomerKey, strProductKey)
		if err == nil {
			if avgNav.AvgNav != nil {
				avgNavLast = *avgNav.AvgNav
			}
		}
		// strAvgNav := fmt.Sprintf("%g", avgNavLast)

		for _, trBalance := range trBalanceCustomer {
			if sisaFifo.Cmp(zero) == 1 {
				var sisaBalance decimal.Decimal

				if trBalance.BalanceUnit.Cmp(sisaFifo) == 1 {
					sisaBalance = trBalance.BalanceUnit.Sub(sisaFifo)
					sisaFifo = zero
				}

				if trBalance.BalanceUnit.Cmp(sisaFifo) == -1 {
					sisaBalance = zero
					sisaFifo = sisaFifo.Sub(trBalance.BalanceUnit)
				}

				if trBalance.BalanceUnit == sisaFifo {
					sisaBalance = zero
					sisaFifo = zero
				}

				paramsBalance := make(map[string]string)
				strAcaKey := strconv.FormatUint(*&trBalance.AcaKey, 10)
				paramsBalance["aca_key"] = strAcaKey
				strTransactionSubs := strconv.FormatUint(trBalance.TcKey, 10)
				paramsBalance["tc_key"] = strTransactionSubs
				strTransactionRed := strconv.FormatUint(transactionConf.TcKey, 10)
				paramsBalance["tc_key_red"] = strTransactionRed

				newlayout := "2006-01-02"
				t, _ := time.Parse(dateLayout, transactionConf.ConfirmDate)
				balanceDate := t.Format(newlayout)

				// strTransUnitSisa := fmt.Sprintf("%g", sisaBalance)

				paramsBalance["balance_date"] = balanceDate + " 00:00:00"
				paramsBalance["balance_unit"] = sisaBalance.String()
				paramsBalance["avg_nav"] = avgNavLast.String()

				var balance models.TrBalance
				_, err = models.GetLastTrBalanceByTcRed(&balance, strTransactionRed)
				if err != nil {
					paramsBalance["rec_order"] = "0"
				} else {
					if balance.RecOrder == nil {
						paramsBalance["rec_order"] = "0"
					} else {
						orderNext := int(*balance.RecOrder) + 1
						strOrderNext := strconv.FormatInt(int64(orderNext), 10)
						paramsBalance["rec_order"] = strOrderNext
					}
				}

				paramsBalance["rec_status"] = "1"
				paramsBalance["rec_created_date"] = time.Now().Format(dateLayout)
				paramsBalance["rec_created_by"] = strIDUserLogin
				log.Println("CreateTrBalance")
				status, err := models.CreateTrBalance(paramsBalance)
				if err != nil {
					// log.Error(err.Error())
					return lib.CustomError(status, err.Error(), "Failed input data")
				}
			} else {
				break
			}
		}
	}

	//update tr_transaction
	params["posted_units"] = transactionConf.ConfirmedUnit.String()
	params["trans_status_key"] = "9"
	params["posted_date"] = time.Now().Format(dateLayout)
	params["rec_modified_by"] = strIDUserLogin
	params["rec_modified_date"] = time.Now().Format(dateLayout)
	log.Println("UpdateTrTransaction")
	_, err = models.UpdateTrTransaction(params)
	if err != nil {
		// log.Error("Error update tr transaction")
		return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed update data")
	}

	//create user message
	if strTransTypeKey != "3" { // SELAIN SWITCH-OUT
		var customer models.MsCustomer
		strCustomerKey := strconv.FormatUint(transaction.CustomerKey, 10)
		_, err = models.GetMsCustomer(&customer, strCustomerKey)
		if err == nil {
			if customer.InvestorType == "263" { //individu
				paramsUserMessage := make(map[string]string)
				paramsUserMessage["umessage_type"] = "245"

				var userLogin models.ScUserLogin
				_, err = models.GetScUserLoginByCustomerKey(&userLogin, strCustomerKey)
				if err != nil {
					// log.Error(err.Error())
					return lib.CustomError(http.StatusBadRequest)
				}

				strUserLoginKey := strconv.FormatUint(userLogin.UserLoginKey, 10)
				paramsUserMessage["umessage_recipient_key"] = strUserLoginKey
				paramsUserMessage["umessage_receipt_date"] = time.Now().Format(dateLayout)
				paramsUserMessage["flag_read"] = "0"
				paramsUserMessage["umessage_sender_key"] = strIDUserLogin
				paramsUserMessage["umessage_sent_date"] = time.Now().Format(dateLayout)
				paramsUserMessage["flag_sent"] = "1"
				var subject string
				var body string
				if strTransTypeKey == "1" { // SUBS
					if transaction.FlagNewSub != nil {
						if *transaction.FlagNewSub == 1 {
							subject = "Subscription Berhasil"
							body = "Subscription kamu telah efektif dibukukan. Silakan cek portofolio di akun kamu untuk melihat transaksi."
							paramsUserMessage["umessage_subject"] = subject
							paramsUserMessage["umessage_body"] = body
						} else {
							subject = "Top Up Berhasil"
							body = "Top Up kamu telah efektif dibukukan. Silakan cek portofolio di akun kamu untuk melihat transaksi."
							paramsUserMessage["umessage_subject"] = subject
							paramsUserMessage["umessage_body"] = body
						}
					} else {
						subject = "Top Up Berhasil"
						body = "Top Up kamu telah efektif dibukukan. Silakan cek portofolio di akun kamu untuk melihat transaksi."
						paramsUserMessage["umessage_subject"] = subject
						paramsUserMessage["umessage_body"] = body
					}
				}

				if strTransTypeKey == "2" { // REDM
					subject = "Redemption Berhasil"
					body = "Redemption kamu telah berhasil dijalankan. Dana akan ditransfer ke rekening bank kamu maks. 7 hari bursa. Silakan cek portofolio di akun kamu untuk melihat transaksi."
					paramsUserMessage["umessage_subject"] = subject
					paramsUserMessage["umessage_body"] = body
				}
				if strTransTypeKey == "4" { // SWITCH
					subject = "Switching Berhasil"
					body = "Switching kamu telah berhasil dijalankan. Silakan cek portofolio di akun kamu untuk melihat transaksi."
					paramsUserMessage["umessage_subject"] = subject
					paramsUserMessage["umessage_body"] = body
				}

				paramsUserMessage["umessage_category"] = "248"
				paramsUserMessage["flag_archieved"] = "0"
				paramsUserMessage["archieved_date"] = time.Now().Format(dateLayout)
				paramsUserMessage["rec_status"] = "1"
				paramsUserMessage["rec_created_date"] = time.Now().Format(dateLayout)
				paramsUserMessage["rec_created_by"] = strIDUserLogin

				status, err = models.CreateScUserMessage(paramsUserMessage)
				if err != nil {
					// log.Error("Error create user message")
					return lib.CustomError(status, err.Error(), "failed input data user message")
				}
				sendEmailTransactionPosted(transaction, transactionConf, userLogin, strCustomerKey, strTransTypeKey)
				getCust := make(map[string]string)
				getCust["customer_key"] = strCustomerKey
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
			} else if customer.InvestorType == "264" { //institusi
				SentEmailTransactionInstitutionPostingBackOfficeToUserCcSales(transactionkey, strCustomerKey)
			}
		}

	}

	// log.Info("Success update transaksi")

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = ""
	return c.JSON(http.StatusOK, response)
}

func sendEmailTransactionPosted(
	transaction models.TrTransaction,
	transactionConf models.TrTransactionConfirmation,
	userLogin models.ScUserLogin,
	strCustomerKey string,
	strTransTypeKey string) {
	if config.Envi != "DEV" {
		decimal.MarshalJSONWithoutQuotes = true

		var subject string
		var tpl bytes.Buffer
		layout := "2006-01-02 15:04:05"
		newLayout := "02 Jan 2006"
		timeLayout := "15:04"

		productFrom := "-"

		var customer models.MsCustomer
		_, err := models.GetMsCustomer(&customer, strCustomerKey)
		if err != nil {
			if err != sql.ErrNoRows {
				// log.Error(err.Error())
			}
		}

		strProductKey := strconv.FormatUint(transaction.ProductKey, 10)

		var product models.MsProduct
		_, err = models.GetMsProduct(&product, strProductKey)
		if err != nil {
			// log.Error(err.Error())
		}

		var trNavData models.TrNav

		var trNav []models.TrNav
		_, err = models.GetTrNavByProductKeyAndNavDate(&trNav, strProductKey, transaction.NavDate)

		if err == nil {
			trNavData = trNav[0]
		}

		var currencyDB models.MsCurrency
		_, err = models.GetMsCurrency(&currencyDB, strconv.FormatUint(*product.CurrencyKey, 10))
		if err != nil {
			// log.Error(err.Error())
		}

		date, _ := time.Parse(layout, transaction.TransDate)
		nabDate, _ := time.Parse(layout, transaction.NavDate)

		var transactionParent models.TrTransaction
		if strTransTypeKey == "4" { // SWITCH IN
			if transaction.ParentKey != nil {
				strTrParentKey := strconv.FormatUint(*transaction.ParentKey, 10)
				_, err := models.GetTrTransaction(&transactionParent, strTrParentKey)
				if err == nil {
					strProductParentKey := strconv.FormatUint(transactionParent.ProductKey, 10)

					var productparent models.MsProduct
					_, err = models.GetMsProduct(&productparent, strProductParentKey)
					if err == nil {
						productFrom = productparent.ProductNameAlt
					}
				}
			}
		}

		bankNameCus := "-"
		noRekCus := "-"
		namaRekCus := "-"
		cabangRek := "-"

		//bank account customer
		strTrKey := strconv.FormatUint(transaction.TransactionKey, 10)
		if strTransTypeKey == "2" { //redm
			var trBankAccount models.TrTransactionBankAccount
			_, err := models.GetTrTransactionBankAccountByField(&trBankAccount, strTrKey, "transaction_key")
			if err == nil {
				strCustBankAcc := strconv.FormatUint(trBankAccount.CustBankaccKey, 10)
				var trBankCust models.MsCustomerBankAccountInfo
				_, err = models.GetMsCustomerBankAccountTransactionByKey(&trBankCust, strCustBankAcc)
				if err == nil {
					bankNameCus = trBankCust.BankName
					noRekCus = trBankCust.AccountNo
					namaRekCus = trBankCust.AccountName
					if trBankCust.BranchName != nil {
						cabangRek = *trBankCust.BranchName
					}
				}
			}
		}

		transType := "Subscription"
		if strTransTypeKey == "1" { // SUBS
			if transaction.FlagNewSub != nil {
				if *transaction.FlagNewSub == 0 {
					transType = "Top Up"
				}
			}
		}

		ac := accounting.Accounting{Symbol: "", Precision: 2, Thousand: ".", Decimal: ","}

		dataReplace := struct {
			FileUrl        string
			Name           string
			Cif            string
			Date           string
			Time           string
			ProductName    string
			Symbol         *string
			Amount         string
			Fee            string
			RedmUnit       string
			NabUnit        string
			BankName       string
			NoRek          string
			NamaRek        string
			Cabang         string
			ProductFrom    string
			ProductTo      string
			NABDate        string
			UnitPenyertaan string
			TransType      string
		}{
			FileUrl:        config.ImageUrl + "/images/mail",
			Name:           customer.FullName,
			Cif:            customer.UnitHolderIDno,
			Date:           date.Format(newLayout),
			Time:           date.Format(timeLayout) + " WIB",
			ProductName:    product.ProductNameAlt,
			Symbol:         currencyDB.Symbol,
			Amount:         ac.FormatMoney(transactionConf.ConfirmedAmount),
			Fee:            ac.FormatMoney(transaction.TransFeeAmount),
			RedmUnit:       ac.FormatMoney(transactionConf.ConfirmedUnit),
			NabUnit:        ac.FormatMoney(trNavData.NavValue),
			BankName:       bankNameCus,
			NoRek:          noRekCus,
			NamaRek:        namaRekCus,
			Cabang:         cabangRek,
			ProductFrom:    productFrom,
			ProductTo:      product.ProductNameAlt,
			NABDate:        nabDate.Format(newLayout),
			UnitPenyertaan: ac.FormatMoney(transactionConf.ConfirmedUnit),
			TransType:      transType}

		if strTransTypeKey == "1" { // SUBS
			if transaction.FlagNewSub != nil {
				if *transaction.FlagNewSub == 1 {
					subject = "[MotionFunds] Subscription Kamu telah Berhasil"
				} else {
					subject = "[MotionFunds] Top Up Kamu telah Berhasil"
				}
			} else {
				subject = "[MotionFunds] Top Up Kamu telah Berhasil"
			}

			t := template.New("email-subscription-posted.html")

			t, err := t.ParseFiles(config.BasePath + "/mail/email-subscription-posted.html")
			if err != nil {
				// log.Println(err)
			}

			if err := t.Execute(&tpl, dataReplace); err != nil {
				// log.Println(err)
			}
		}

		if strTransTypeKey == "2" { // REDM
			subject = "[MotionFunds] Redemption Kamu teleh Berhasil"
			t := template.New("email-redemption-posted.html")

			t, err := t.ParseFiles(config.BasePath + "/mail/email-redemption-posted.html")
			if err != nil {
				// log.Println(err)
			}

			if err := t.Execute(&tpl, dataReplace); err != nil {
				// log.Println(err)
			}

		}

		if strTransTypeKey == "4" { // SWITCH
			subject = "[MotionFunds] Switching Kamu telah Berhasil"
			t := template.New("email-switching-posted.html")

			t, err := t.ParseFiles(config.BasePath + "/mail/email-switching-posted.html")
			if err != nil {
				// log.Println(err)
			}

			if err := t.Execute(&tpl, dataReplace); err != nil {
				// log.Println(err)
			}

		}

		// Send email
		result := tpl.String()

		mailer := gomail.NewMessage()
		// mailer.SetHeader("From", config.EmailFrom)
		mailer.SetHeader("To", userLogin.UloginEmail)
		mailer.SetHeader("Subject", subject)
		mailer.SetBody("text/html", result)

		err = lib.SendEmail(mailer)
		if err != nil {
			// log.Info("Email sent error")
			// log.Error(err)
		} else {
			// log.Info("Email sent")
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
		// 	// log.Info("Email sent error")
		// 	// log.Error(err)
		// } else {
		// 	// log.Info("Email sent")
		// }
	}
}

func GetCustomerBankAccount(c echo.Context) error {

	var err error
	var status int

	keyStr := c.Param("key")
	key, _ := strconv.ParseUint(keyStr, 10, 64)
	if key == 0 {
		return lib.CustomError(http.StatusNotFound)
	}

	var transaction models.TrTransaction
	status, err = models.GetTrTransaction(&transaction, keyStr)
	if err != nil {
		return lib.CustomError(status)
	}

	strCusKey := strconv.FormatUint(transaction.CustomerKey, 10)

	var customerBankAccountInfo []models.MsCustomerBankAccountInfo
	status, err = models.GetAllMsCustomerBankAccountTransaction(&customerBankAccountInfo, strCusKey)
	if err != nil {
		if err != sql.ErrNoRows {
			// log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = customerBankAccountInfo

	return c.JSON(http.StatusOK, response)
}

func GetProductBankAccount(c echo.Context) error {

	var err error
	var status int

	keyStr := c.Param("key")
	key, _ := strconv.ParseUint(keyStr, 10, 64)
	if key == 0 {
		return lib.CustomError(http.StatusNotFound)
	}

	var transaction models.TrTransaction
	status, err = models.GetTrTransaction(&transaction, keyStr)
	if err != nil {
		return lib.CustomError(status)
	}

	strProductKey := strconv.FormatUint(transaction.ProductKey, 10)

	var lookupTransType string

	if transaction.TransTypeKey == 1 || transaction.TransTypeKey == 4 { //sub + sw.in
		lookupTransType = "269"
	} else { //red + sw.out
		lookupTransType = "270"
	}

	var bankAccountTransactionInfo []models.MsProductBankAccountTransactionInfo

	status, err = models.GetAllMsProductBankAccountTransaction(&bankAccountTransactionInfo, strProductKey, lookupTransType)
	if err != nil {
		if err != sql.ErrNoRows {
			// log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = bankAccountTransactionInfo

	return c.JSON(http.StatusOK, response)
}

func ProsesUnposting(c echo.Context) error {
	errorAuth := initAuthFundAdmin()
	if errorAuth != nil {
		// log.Error("User Autorizer")
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}

	var err error
	var status int

	params := make(map[string]string)

	transactionkey := c.FormValue("transaction_key")
	if transactionkey == "" {
		// log.Error("Missing required parameter: transaction_key")
		return lib.CustomError(http.StatusBadRequest)
	}

	n, err := strconv.ParseUint(transactionkey, 10, 64)
	if err == nil && n > 0 {
		params["transaction_key"] = transactionkey
	} else {
		// log.Error("Wrong input for parameter: transaction_key")
		return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: transaction_key", "Wrong input for parameter: transaction_key")
	}

	var transaction models.TrTransaction
	status, err = models.GetTrTransaction(&transaction, transactionkey)
	if err != nil {
		// log.Error("Transaction not exist")
		return lib.CustomError(status)
	}

	strTransStatusKey := strconv.FormatUint(transaction.TransStatusKey, 10)

	if strTransStatusKey != "9" {
		// log.Error("User Autorizer, status transaksi not posted")
		return lib.CustomError(http.StatusBadRequest)
	}

	strCustomer := strconv.FormatUint(transaction.CustomerKey, 10)
	strPro := strconv.FormatUint(transaction.ProductKey, 10)

	var transAfter models.TrTransaction
	status, err = models.CheckTrTransactionLastProductCustomer(&transAfter, strCustomer, strPro, transactionkey)
	if err == nil {
		// log.Error("Transaksi tidak dapat di Un-posting karena bukan data terakhir dari customer dan produk yang sama.")
		return lib.CustomError(http.StatusBadRequest, "Transaksi tidak dapat di Un-posting karena bukan data terakhir dari customer dan produk yang sama.", "Transaksi tidak dapat di Un-posting karena bukan data terakhir dari customer dan produk yang sama.")
	}

	dateLayout := "2006-01-02 15:04:05"
	strIDUserLogin := strconv.FormatUint(lib.Profile.UserID, 10)

	strTransTypeKey := strconv.FormatUint(transaction.TransTypeKey, 10)

	var transactionConf models.TrTransactionConfirmation
	strTransactionKey := strconv.FormatUint(transaction.TransactionKey, 10)
	_, err = models.GetTrTransactionConfirmationByTransactionKey(&transactionConf, strTransactionKey)
	if err != nil {
		// log.Error(err.Error())
		return lib.CustomError(http.StatusBadRequest)
	}

	//update tr_transaction
	params["trans_status_key"] = "7"
	params["rec_modified_by"] = strIDUserLogin
	params["rec_modified_date"] = time.Now().Format(dateLayout)

	_, err = models.UpdateTrTransaction(params)
	if err != nil {
		// log.Error("Error update tr transaction")
		return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed update data")
	}

	//update delete tr_transaction_confirmation
	strTcKey := strconv.FormatUint(transactionConf.TcKey, 10)
	paramsConf := make(map[string]string)
	paramsConf["tc_key"] = strTcKey
	paramsConf["rec_status"] = "0"
	paramsConf["rec_deleted_by"] = strIDUserLogin
	paramsConf["rec_deleted_date"] = time.Now().Format(dateLayout)

	_, err = models.UpdateTrTransactionConfirmation(paramsConf)
	if err != nil {
		// log.Error("Error delete tr transaction condirmation")
		return lib.CustomError(http.StatusInternalServerError, err.Error(), "Error delete tr transaction condirmation")
	}

	//update delete tr_transaction_fifo
	paramsFifo := make(map[string]string)
	paramsFifo["rec_status"] = "0"
	paramsFifo["rec_deleted_by"] = strIDUserLogin
	paramsFifo["rec_deleted_date"] = time.Now().Format(dateLayout)

	var field string
	if strTransTypeKey == "1" || strTransTypeKey == "4" {
		field = "trans_conf_sub_key"
	} else {
		field = "trans_conf_red_key"
	}

	_, err = models.UpdateTrTransactionFifo(paramsFifo, strTcKey, field)
	if err != nil {
		// log.Error("Error delete tr transaction fifo")
		return lib.CustomError(http.StatusInternalServerError, err.Error(), "Error delete tr transaction fifo")
	}

	//update delete tr_balance
	paramsBalance := make(map[string]string)
	paramsBalance["rec_status"] = "0"
	paramsBalance["rec_deleted_by"] = strIDUserLogin
	paramsBalance["rec_deleted_date"] = time.Now().Format(dateLayout)

	var fieldBalance string
	if strTransTypeKey == "1" || strTransTypeKey == "4" {
		fieldBalance = "tc_key"
	} else {
		fieldBalance = "tc_key_red"
	}

	_, err = models.UpdateTrBalance(paramsFifo, strTcKey, fieldBalance)
	if err != nil {
		// log.Error("Error delete tr balance")
		return lib.CustomError(http.StatusInternalServerError, err.Error(), "Error delete tr balance")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = ""
	return c.JSON(http.StatusOK, response)
}

func DataTransaksiInquiry(c echo.Context) error {
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

	var isAll bool
	if noLimit == true {
		isAll = true
	} else {
		isAll = false
	}

	items := []string{"transaction_key", "branch_key", "agent_key", "customer_key", "product_key", "trans_date", "trans_amount", "trans_bank_key"}

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
		params["orderBy"] = "transaction_key"
		params["orderType"] = "ASC"
	}

	productKey := c.QueryParam("product_key")
	if productKey != "" {
		productKeyCek, err := strconv.ParseUint(productKey, 10, 64)
		if err == nil && productKeyCek > 0 {
			params["product_key"] = productKey
		} else {
			// log.Error("Wrong input for parameter: product_key")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: product_key", "Missing required parameter: product_key")
		}
	}

	transstatuskey := c.QueryParam("trans_status_key")
	if transstatuskey != "" {
		transstatuskeyCek, err := strconv.ParseUint(transstatuskey, 10, 64)
		if err == nil && transstatuskeyCek > 0 {
			params["trans_status_key"] = transstatuskey
		} else {
			// log.Error("Wrong input for parameter: trans_status_key")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: trans_status_key", "Missing required parameter: trans_status_key")
		}
	}

	transtypekey := c.QueryParam("trans_type_key")
	if transtypekey != "" {
		transtypekeyCek, err := strconv.ParseUint(transtypekey, 10, 64)
		if err == nil && transtypekeyCek > 0 {
			params["trans_type_key"] = transtypekey
		} else {
			// log.Error("Wrong input for parameter: trans_type_key")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: trans_type_key", "Missing required parameter: trans_type_key")
		}
	}

	navdate := c.QueryParam("nav_date")
	if navdate != "" {
		params["nav_date"] = navdate
	}

	customerkey := c.QueryParam("customer_key")
	if customerkey != "" {
		params["customer_key"] = customerkey
	}

	if lib.Profile.UserCategoryKey == 3 { //user branch
		// log.Println(lib.Profile)
		if lib.Profile.BranchKey != nil {
			params["branch_key"] = strconv.FormatUint(*lib.Profile.BranchKey, 10)
		}
	}

	var transStatusKey []string

	var trTransaction []models.TrTransaction

	status, err = models.AdminGetAllTrTransaction(&trTransaction, limit, offset, noLimit, params, transStatusKey, "trans_status_key", isAll, strconv.FormatUint(lib.Profile.UserID, 10))

	if err != nil {
		// log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if len(trTransaction) < 1 {
		// log.Error("transaction not found")
		return lib.CustomError(http.StatusNotFound, "Transaction not found", "Transaction not found")
	}

	var branchIds []string
	var lookupIds []string
	var agentIds []string
	var customerIds []string
	var productIds []string
	var transTypeIds []string
	var transStatusIds []string
	var transactionIds []string
	for _, tr := range trTransaction {

		if tr.BranchKey != nil {
			if _, ok := lib.Find(branchIds, strconv.FormatUint(*tr.BranchKey, 10)); !ok {
				branchIds = append(branchIds, strconv.FormatUint(*tr.BranchKey, 10))
			}
		}
		if tr.AgentKey != nil {
			if _, ok := lib.Find(agentIds, strconv.FormatUint(*tr.AgentKey, 10)); !ok {
				agentIds = append(agentIds, strconv.FormatUint(*tr.AgentKey, 10))
			}
		}
		if _, ok := lib.Find(transactionIds, strconv.FormatUint(tr.TransactionKey, 10)); !ok {
			transactionIds = append(transactionIds, strconv.FormatUint(tr.TransactionKey, 10))
		}
		if _, ok := lib.Find(customerIds, strconv.FormatUint(tr.CustomerKey, 10)); !ok {
			customerIds = append(customerIds, strconv.FormatUint(tr.CustomerKey, 10))
		}
		if _, ok := lib.Find(productIds, strconv.FormatUint(tr.ProductKey, 10)); !ok {
			productIds = append(productIds, strconv.FormatUint(tr.ProductKey, 10))
		}
		if _, ok := lib.Find(transTypeIds, strconv.FormatUint(tr.TransTypeKey, 10)); !ok {
			transTypeIds = append(transTypeIds, strconv.FormatUint(tr.TransTypeKey, 10))
		}
		if _, ok := lib.Find(transStatusIds, strconv.FormatUint(tr.TransStatusKey, 10)); !ok {
			transStatusIds = append(transStatusIds, strconv.FormatUint(tr.TransStatusKey, 10))
		}
		if tr.TransSource != nil {
			if _, ok := lib.Find(lookupIds, strconv.FormatUint(*tr.TransSource, 10)); !ok {
				lookupIds = append(lookupIds, strconv.FormatUint(*tr.TransSource, 10))
			}
		}
	}

	//gen lookup oa request
	var lookupOaReq []models.GenLookup
	if len(lookupIds) > 0 {
		status, err = models.GetGenLookupIn(&lookupOaReq, lookupIds, "lookup_key")
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

	//mapping branch
	var msBranch []models.MsBranch
	if len(branchIds) > 0 {
		status, err = models.GetMsBranchIn(&msBranch, branchIds, "branch_key")
		if err != nil {
			// log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
	}
	branchData := make(map[uint64]models.MsBranch)
	for _, b := range msBranch {
		branchData[b.BranchKey] = b
	}

	//mapping agent
	var msAgent []models.MsAgent
	if len(agentIds) > 0 {
		status, err = models.GetMsAgentIn(&msAgent, agentIds, "agent_key")
		if err != nil {
			// log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
	}
	agentData := make(map[uint64]models.MsAgent)
	for _, a := range msAgent {
		agentData[a.AgentKey] = a
	}

	//mapping customer
	var msCustomer []models.MsCustomer
	if len(customerIds) > 0 {
		status, err = models.GetMsCustomerIn(&msCustomer, customerIds, "customer_key")
		if err != nil {
			// log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
	}
	customerData := make(map[uint64]models.MsCustomer)
	for _, c := range msCustomer {
		customerData[c.CustomerKey] = c
	}

	//user customer
	var userLogin []models.ScUserLogin
	if len(customerIds) > 0 {
		status, err = models.GetScUserLoginIn(&userLogin, customerIds, "customer_key")
		if err != nil {
			// log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
	}
	userLoginData := make(map[uint64]models.ScUserLogin)
	for _, c := range userLogin {
		userLoginData[*c.CustomerKey] = c
	}

	//mapping product
	var msProduct []models.MsProduct
	if len(productIds) > 0 {
		status, err = models.GetMsProductIn(&msProduct, productIds, "product_key")
		if err != nil {
			if err != sql.ErrNoRows {
				// log.Error(err.Error())
				return lib.CustomError(status, err.Error(), "Failed get data")
			}
		}
	}
	productData := make(map[uint64]models.MsProduct)
	for _, p := range msProduct {
		productData[p.ProductKey] = p
	}

	//mapping Trans type
	var transactionType []models.TrTransactionType
	if len(transTypeIds) > 0 {
		status, err = models.GetMsTransactionTypeIn(&transactionType, transTypeIds, "trans_type_key")
		if err != nil {
			if err != sql.ErrNoRows {
				// log.Error(err.Error())
				return lib.CustomError(status, err.Error(), "Failed get data")
			}
		}
	}
	transactionTypeData := make(map[uint64]models.TrTransactionType)
	for _, t := range transactionType {
		transactionTypeData[t.TransTypeKey] = t
	}

	//mapping trans status
	var trTransactionStatus []models.TrTransactionStatus
	if len(transStatusIds) > 0 {
		status, err = models.GetMsTransactionStatusIn(&trTransactionStatus, transStatusIds, "trans_status_key")
		if err != nil {
			if err != sql.ErrNoRows {
				// log.Error(err.Error())
				return lib.CustomError(status, err.Error(), "Failed get data")
			}
		}
	}
	transStatusData := make(map[uint64]models.TrTransactionStatus)
	for _, ts := range trTransactionStatus {
		transStatusData[ts.TransStatusKey] = ts
	}

	//mapping tc confirmation
	var transConf []models.TrTransactionConfirmation
	if len(transactionIds) > 0 {
		status, err = models.GetTrTransactionConfirmationIn(&transConf, transactionIds, "transaction_key")
		if err != nil {
			if err != sql.ErrNoRows {
				// log.Error(err.Error())
				return lib.CustomError(status, err.Error(), "Failed get TC data")
			}
		}
	}
	tcData := make(map[uint64]models.TrTransactionConfirmation)
	for _, tc := range transConf {
		tcData[tc.TransactionKey] = tc
	}

	var responseData []models.AdminTrTransactionInquiryList
	for _, tr := range trTransaction {
		var data models.AdminTrTransactionInquiryList

		data.TransactionKey = tr.TransactionKey

		if tr.BranchKey != nil {
			if n, ok := branchData[*tr.BranchKey]; ok {
				data.BranchName = n.BranchName
			}
		}

		if tr.AgentKey != nil {
			if n, ok := agentData[*tr.AgentKey]; ok {
				data.AgentName = n.AgentName
			}
		}

		if n, ok := customerData[tr.CustomerKey]; ok {
			data.CustomerName = n.FullName
		}

		if n, ok := productData[tr.ProductKey]; ok {
			data.ProductName = n.ProductNameAlt
		}

		if n, ok := transStatusData[tr.TransStatusKey]; ok {
			data.TransStatus = *n.StatusCode
		}

		if tr.TransSource != nil {
			if n, ok := gData[*tr.TransSource]; ok {
				data.TransSource = n.LkpName
			}
		}

		layout := "2006-01-02 15:04:05"
		newLayout := "02 Jan 2006"
		date, _ := time.Parse(layout, tr.TransDate)
		data.TransDate = date.Format(newLayout)
		date, _ = time.Parse(layout, tr.NavDate)
		data.NavDate = date.Format(newLayout)

		if n, ok := transactionTypeData[tr.TransTypeKey]; ok {
			data.TransType = *n.TypeDescription
		}

		if tc, ok := tcData[tr.TransactionKey]; ok {
			data.TransAmount = tc.ConfirmedAmount
			data.TransUnit = tc.ConfirmedUnit
		} else {
			data.TransAmount = tr.TransAmount
			data.TransUnit = tr.TransUnit
		}

		data.TotalAmount = tr.TotalAmount

		data.RecImage1 = nil

		if tr.TransTypeKey == uint64(1) && tr.PaymentMethod != nil {
			if *tr.PaymentMethod == uint64(284) { //transfer manual
				if tr.RecImage1 != nil {
					if n, ok := userLoginData[tr.CustomerKey]; ok {
						path := config.ImageUrl + "/images/user/" + strconv.FormatUint(n.UserLoginKey, 10) + "/transfer/" + *tr.RecImage1
						data.RecImage1 = &path
					}
				}
			}
		}

		responseData = append(responseData, data)
	}

	var countData models.CountData
	var pagination int
	if limit > 0 {
		status, err = models.AdminGetCountTrTransaction(&countData, params, transStatusKey, "trans_status_key", strconv.FormatUint(lib.Profile.UserID, 10))
		if err != nil {
			// log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
		if int(countData.CountData) < int(limit) {
			pagination = 1
		} else {
			// log.Println("JUMLAH DATANYA => ", countData.CountData)
			// log.Println("LIMIT => ", limit)
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

func DetailTransaksiInquiry(c echo.Context) error {

	var err error
	var status int
	decimal.MarshalJSONWithoutQuotes = true

	keyStr := c.Param("key")
	key, _ := strconv.ParseUint(keyStr, 10, 64)
	if key == 0 {
		return lib.CustomError(http.StatusNotFound)
	}

	var transaction models.TrTransaction
	status, err = models.GetTrTransaction(&transaction, keyStr)
	if err != nil {
		return lib.CustomError(status)
	}

	var responseData models.AdminTransactionDetail

	var lookupIds []string
	var pChannelIds []string

	var trSettlement []models.TrTransactionSettlement
	paramSettlement := make(map[string]string)
	paramSettlement["rec_status"] = "1"
	paramSettlement["transaction_key"] = strconv.FormatUint(transaction.TransactionKey, 10)
	status, err = models.GetAllTrTransactionSettlement(&trSettlement, paramSettlement)
	if err == nil {
		if len(trSettlement) > 0 {
			for _, settlement := range trSettlement {
				if _, ok := lib.Find(lookupIds, strconv.FormatUint(settlement.SettlePurposed, 10)); !ok {
					lookupIds = append(lookupIds, strconv.FormatUint(settlement.SettlePurposed, 10))
				}

				if _, ok := lib.Find(lookupIds, strconv.FormatUint(settlement.SettleStatus, 10)); !ok {
					lookupIds = append(lookupIds, strconv.FormatUint(settlement.SettleStatus, 10))
				}

				if _, ok := lib.Find(lookupIds, strconv.FormatUint(settlement.SettleChannel, 10)); !ok {
					lookupIds = append(lookupIds, strconv.FormatUint(settlement.SettleChannel, 10))
				}

				if _, ok := lib.Find(pChannelIds, strconv.FormatUint(settlement.SettlePaymentMethod, 10)); !ok {
					pChannelIds = append(pChannelIds, strconv.FormatUint(settlement.SettlePaymentMethod, 10))
				}
			}
		}
	}

	if transaction.TrxCode != nil {
		if _, ok := lib.Find(lookupIds, strconv.FormatUint(*transaction.TrxCode, 10)); !ok {
			lookupIds = append(lookupIds, strconv.FormatUint(*transaction.TrxCode, 10))
		}
	}
	if transaction.EntryMode != nil {
		if _, ok := lib.Find(lookupIds, strconv.FormatUint(*transaction.EntryMode, 10)); !ok {
			lookupIds = append(lookupIds, strconv.FormatUint(*transaction.EntryMode, 10))
		}
	}
	if transaction.PaymentMethod != nil {
		if _, ok := lib.Find(lookupIds, strconv.FormatUint(*transaction.PaymentMethod, 10)); !ok {
			lookupIds = append(lookupIds, strconv.FormatUint(*transaction.PaymentMethod, 10))
		}
	}
	if transaction.TrxRiskLevel != nil {
		if _, ok := lib.Find(lookupIds, strconv.FormatUint(*transaction.TrxRiskLevel, 10)); !ok {
			lookupIds = append(lookupIds, strconv.FormatUint(*transaction.TrxRiskLevel, 10))
		}
	}
	if transaction.TransSource != nil {
		if _, ok := lib.Find(lookupIds, strconv.FormatUint(*transaction.TransSource, 10)); !ok {
			lookupIds = append(lookupIds, strconv.FormatUint(*transaction.TransSource, 10))
		}
	}

	//ms payment channel branch
	var pChannel []models.MsPaymentChannel
	if len(pChannelIds) > 0 {
		status, err = models.GetMsPaymentChannelIn(&pChannel, pChannelIds, "pchannel_key")
		if err != nil {
			// log.Error(err.Error())
		}
	}
	channelData := make(map[uint64]models.MsPaymentChannel)
	for _, pc := range pChannel {
		channelData[pc.PchannelKey] = pc
	}

	//gen lookup oa request
	var lookupOaReq []models.GenLookup
	if len(lookupIds) > 0 {
		status, err = models.GetGenLookupIn(&lookupOaReq, lookupIds, "lookup_key")
		if err != nil {
			if err != sql.ErrNoRows {
				// log.Error(err.Error())
				return lib.CustomError(status, err.Error(), "Failed get data")
			}
		}
	}

	layout := "2006-01-02 15:04:05"
	newLayout := "02 Jan 2006"

	gData := make(map[uint64]models.GenLookup)
	for _, gen := range lookupOaReq {
		gData[gen.LookupKey] = gen
	}

	//set settlement
	var settlementTransactionList []models.TransactionSettlement
	if len(trSettlement) > 0 {
		for _, settlement := range trSettlement {
			var data models.TransactionSettlement
			data.SettlementKey = settlement.SettlementKey
			date, _ := time.Parse(layout, settlement.SettleDate)
			data.SettleDate = date.Format(newLayout)
			data.SettleNominal = settlement.SettleNominal
			if settlement.SettleRealizedDate != nil {
				date, _ = time.Parse(layout, *settlement.SettleRealizedDate)
				data.SettleRealizedDate = date.Format(newLayout)
			} else {
				data.SettleRealizedDate = ""
			}
			data.SettleRealizedDate = date.Format(newLayout)
			data.SettleRemarks = settlement.SettleRemarks
			data.SettleReference = settlement.SettleReference
			if n, ok := gData[settlement.SettlePurposed]; ok {
				data.SettlePurposed = *n.LkpName
			}
			if n, ok := gData[settlement.SettleStatus]; ok {
				data.SettleStatus = *n.LkpName
			}
			if n, ok := gData[settlement.SettleChannel]; ok {
				data.SettleChannel = *n.LkpName
			}

			data.SettlePaymentMethod = ""
			if n, ok := channelData[settlement.SettlePaymentMethod]; ok {
				data.SettlePaymentMethod = *n.PchannelName
			}

			settlementTransactionList = append(settlementTransactionList, data)
		}
	}
	responseData.TransactionSettlement = &settlementTransactionList

	if transaction.TrxCode != nil {
		if n, ok := gData[*transaction.TrxCode]; ok {
			var trc models.LookupTrans

			trc.LookupKey = n.LookupKey
			trc.LkpGroupKey = n.LkpGroupKey
			trc.LkpCode = n.LkpCode
			trc.LkpName = n.LkpName
			responseData.TrxCode = &trc
		}
	}

	if transaction.EntryMode != nil {
		if n, ok := gData[*transaction.EntryMode]; ok {
			var entm models.LookupTrans

			entm.LookupKey = n.LookupKey
			entm.LkpGroupKey = n.LkpGroupKey
			entm.LkpCode = n.LkpCode
			entm.LkpName = n.LkpName
			responseData.EntryMode = &entm
		}
	}

	if transaction.TransSource != nil {
		if n, ok := gData[*transaction.TransSource]; ok {
			responseData.TransSource = n.LkpName
		}
	}

	if transaction.PaymentMethod != nil {
		if n, ok := gData[*transaction.PaymentMethod]; ok {
			var pm models.LookupTrans
			pm.LookupKey = n.LookupKey
			pm.LkpGroupKey = n.LkpGroupKey
			pm.LkpCode = n.LkpCode
			pm.LkpName = n.LkpName
			responseData.PaymentMethod = &pm
		}
	}

	if transaction.TrxRiskLevel != nil {
		if n, ok := gData[*transaction.TrxRiskLevel]; ok {
			var risk models.LookupTrans

			risk.LookupKey = n.LookupKey
			risk.LkpGroupKey = n.LkpGroupKey
			risk.LkpCode = n.LkpCode
			risk.LkpName = n.LkpName
			responseData.TrxRiskLevel = &risk
		}
	}

	responseData.TransactionKey = transaction.TransactionKey
	date, _ := time.Parse(layout, transaction.TransDate)
	responseData.TransDate = date.Format(newLayout)
	date, _ = time.Parse(layout, transaction.NavDate)
	responseData.NavDate = date.Format(newLayout)
	if transaction.RecCreatedDate != nil {
		date, err = time.Parse(layout, *transaction.RecCreatedDate)
		if err == nil {
			oke := date.Format(newLayout)
			responseData.RecCreatedDate = &oke
		}
	}
	responseData.RecCreatedBy = transaction.RecCreatedBy
	responseData.TransAmount = transaction.TransAmount
	responseData.TransUnit = transaction.TransUnit
	responseData.TransUnitPercent = transaction.TransUnitPercent
	if transaction.FlagRedemtAll != nil {
		if int(*transaction.FlagRedemtAll) > 0 {
			responseData.FlagRedemtAll = true
		}
	}
	if transaction.FlagNewSub != nil {
		if int(*transaction.FlagNewSub) > 0 {
			responseData.FlagNewSub = true
		}
	}
	responseData.TransEntry = &transaction.TransEntry
	responseData.TransFeePercent = transaction.TransFeePercent
	responseData.TransFeeAmount = transaction.TransFeeAmount
	responseData.ChargesFeeAmount = transaction.ChargesFeeAmount
	responseData.ServicesFeeAmount = transaction.ServicesFeeAmount
	responseData.TotalAmount = transaction.TotalAmount
	responseData.SettlementDate = transaction.SettlementDate
	responseData.TransBankAccNo = transaction.TransBankAccNo
	responseData.TransBankaccName = transaction.TransBankaccName
	responseData.TransRemarks = transaction.TransRemarks
	responseData.TransReferences = transaction.TransReferences
	responseData.PromoCode = transaction.PromoCode
	responseData.SalesCode = transaction.SalesCode
	if transaction.RiskWaiver > 0 {
		responseData.RiskWaiver = true
	}
	responseData.FileUploadDate = transaction.FileUploadDate
	responseData.ProceedDate = transaction.ProceedDate
	responseData.ProceedAmount = transaction.ProceedAmount
	responseData.SentDate = transaction.SentDate
	responseData.SentReferences = transaction.SentReferences
	responseData.ConfirmedDate = transaction.ConfirmedDate
	responseData.PostedDate = transaction.PostedDate
	responseData.PostedUnits = transaction.PostedUnits
	responseData.SettledDate = transaction.SettledDate
	responseData.StampFeeAmount = transaction.StampFeeAmount

	strCustomer := strconv.FormatUint(transaction.CustomerKey, 10)

	dir := ""

	if transaction.RecApprovalStage == nil {
		var userData models.ScUserLogin
		status, err = models.GetScUserLoginByCustomerKey(&userData, strCustomer)
		if err != nil {
			if err != sql.ErrNoRows {
				return lib.CustomError(status)
			}
		}
		dir = config.ImageUrl + "/images/user/" + strconv.FormatUint(userData.UserLoginKey, 10) + "/transfer/"
	} else {
		dir = config.ImageUrl + "/images/user/institusi/" + strCustomer + "/transfer/"
	}

	if transaction.RecImage1 != nil {
		path := dir + *transaction.RecImage1
		responseData.RecImage1 = &path
	}

	if transaction.BranchKey != nil {
		var branch models.MsBranch
		strBranch := strconv.FormatUint(*transaction.BranchKey, 10)
		status, err = models.GetMsBranch(&branch, strBranch)
		if err != nil {
			if err != sql.ErrNoRows {
				return lib.CustomError(status)
			}
		} else {
			var br models.BranchTrans
			br.BranchKey = branch.BranchKey
			br.BranchCode = branch.BranchCode
			br.BranchName = branch.BranchName
			responseData.Branch = &br
		}
	}

	//check agent
	if transaction.AgentKey != nil {
		var agent models.MsAgent
		strAgent := strconv.FormatUint(*transaction.AgentKey, 10)
		status, err = models.GetMsAgent(&agent, strAgent)
		if err != nil {
			if err != sql.ErrNoRows {
				return lib.CustomError(status)
			}
		} else {
			var a models.AgentTrans
			a.AgentKey = agent.AgentKey
			a.AgentCode = agent.AgentCode
			a.AgentName = agent.AgentName
			responseData.Agent = &a
		}
	}

	//check customer
	var customer models.MsCustomer
	strCus := strconv.FormatUint(transaction.CustomerKey, 10)
	status, err = models.GetMsCustomer(&customer, strCus)
	if err != nil {
		if err != sql.ErrNoRows {
			return lib.CustomError(status)
		}
	} else {
		responseData.Customer.CustomerKey = customer.CustomerKey
		responseData.Customer.FullName = customer.FullName
		responseData.Customer.SidNo = customer.SidNo
		responseData.Customer.UnitHolderIDno = customer.UnitHolderIDno
	}

	//check product
	var product models.MsProduct
	strPro := strconv.FormatUint(transaction.ProductKey, 10)
	status, err = models.GetMsProduct(&product, strPro)
	if err != nil {
		if err != sql.ErrNoRows {
			return lib.CustomError(status)
		}
	} else {
		if product.FundTypeKey != nil {
			responseData.FundTypeKey = product.FundTypeKey
		}
		responseData.Product.ProductKey = product.ProductKey
		responseData.Product.ProductCode = product.ProductCode
		responseData.Product.ProductName = product.ProductName
	}

	//check trans status
	var transStatus models.TrTransactionStatus
	strTrSt := strconv.FormatUint(transaction.TransStatusKey, 10)
	status, err = models.GetTrTransactionStatus(&transStatus, strTrSt)
	if err != nil {
		if err != sql.ErrNoRows {
			return lib.CustomError(status)
		}
	} else {
		responseData.TransStatus.TransStatusKey = transStatus.TransStatusKey
		responseData.TransStatus.StatusCode = transStatus.StatusCode
		responseData.TransStatus.StatusDescription = transStatus.StatusDescription
	}

	//check trans type
	var transType models.TrTransactionType
	strTrTy := strconv.FormatUint(transaction.TransTypeKey, 10)
	status, err = models.GetMsTransactionType(&transType, strTrTy)
	if err != nil {
		if err != sql.ErrNoRows {
			return lib.CustomError(status)
		}
	} else {
		responseData.TransType.TransTypeKey = transType.TransTypeKey
		responseData.TransType.TypeCode = transType.TypeCode
		responseData.TransType.TypeDescription = transType.TypeDescription
	}

	//check bank
	var bank models.MsBank
	if transaction.TransBankKey != nil {
		strBank := strconv.FormatUint(*transaction.TransBankKey, 10)
		status, err = models.GetMsBank(&bank, strBank)
		if err != nil {
			if err != sql.ErrNoRows {
				return lib.CustomError(status)
			}
		} else {
			var tb models.TransBank
			tb.BankKey = bank.BankKey
			tb.BankCode = bank.BankCode
			tb.BankName = bank.BankName
			responseData.TransBank = &tb
		}
	}

	//check aca
	if transaction.AcaKey != nil {
		var aca models.TrAccountAgent
		strAca := strconv.FormatUint(*transaction.AcaKey, 10)
		status, err = models.GetTrAccountAgent(&aca, strAca)
		if err != nil {
			if err != sql.ErrNoRows {
				return lib.CustomError(status)
			}
		} else {
			var ac models.AcaTrans
			ac.AcaKey = aca.AcaKey
			ac.AccKey = aca.AccKey
			var agent models.MsAgent
			strAgent := strconv.FormatUint(aca.AgentKey, 10)
			status, err = models.GetMsAgent(&agent, strAgent)
			if err != nil {
				if err != sql.ErrNoRows {
					return lib.CustomError(status)
				}
			} else {
				ac.AgentKey = agent.AgentKey
				ac.AgentCode = agent.AgentCode
				ac.AgentName = agent.AgentName
			}

			responseData.Aca = &ac
		}
	}

	//check transaction confirmation
	strTrKey := strconv.FormatUint(transaction.TransactionKey, 10)
	// if transaction.AcaKey != nil {
	var tc models.TrTransactionConfirmation
	status, err = models.GetTrTransactionConfirmationByTransactionKey(&tc, strTrKey)
	if err != nil {
		if err != sql.ErrNoRows {
			return lib.CustomError(status)
		}
	} else {
		var transTc models.TrTransactionConfirmationInfo
		transTc.TcKey = tc.TcKey
		date, _ := time.Parse(layout, tc.ConfirmDate)
		transTc.ConfirmDate = date.Format(newLayout)
		transTc.ConfirmedAmount = tc.ConfirmedAmount
		transTc.ConfirmedUnit = tc.ConfirmedUnit
		transTc.ConfirmedAmountDiff = tc.ConfirmedAmountDiff
		transTc.ConfirmedUnitDiff = tc.ConfirmedUnitDiff

		responseData.TransactionConfirmationInfo = &transTc
	}
	responseData.TransCalcMethod = transaction.TransCalcMethod
	responseData.NavDateReal = transaction.NavDate
	// }

	//cek promo
	if transaction.PromoCode != nil {
		var promo models.TrPromoData
		status, err = models.AdminGetDetailTransactionPromo(&promo, strTrKey, *transaction.PromoCode)
		if err == nil {
			responseData.Promo = &promo
		}
	}

	//bank transaction customer
	var trBankAccount models.TrTransactionBankAccount
	status, err = models.GetTrTransactionBankAccountByField(&trBankAccount, strTrKey, "transaction_key")
	if err == nil {
		strCustBankAcc := strconv.FormatUint(trBankAccount.CustBankaccKey, 10)
		var trBankCust models.MsCustomerBankAccountInfo
		status, err = models.GetMsCustomerBankAccountTransactionByKey(&trBankCust, strCustBankAcc)
		if err == nil {
			responseData.CustomerBankAccount = &trBankCust
		}

		strProdBankAcc := strconv.FormatUint(trBankAccount.ProdBankaccKey, 10)
		var prodBankAccount models.MsProductBankAccountTransactionInfo
		status, err = models.GetMsProductBankAccountTransactionByKey(&prodBankAccount, strProdBankAcc)
		if err == nil {
			responseData.ProductBankAccount = &prodBankAccount
		}
	}

	responseData.IsEnableUnposting = false
	responseData.MessageEnableUnposting = ""

	if strTrSt == "9" {
		responseData.MessageEnableUnposting = "Transaksi tidak dapat di Un-posting karena bukan data terakhir dari customer dan produk yang sama."
		var transAfter models.TrTransaction
		status, err = models.CheckTrTransactionLastProductCustomer(&transAfter, strCustomer, strPro, keyStr)
		if err != nil {
			if err == sql.ErrNoRows {
				responseData.IsEnableUnposting = true
				responseData.MessageEnableUnposting = ""
			}
		}
	}

	paramsFile := make(map[string]string)
	paramsFile["ref_fk_key"] = strconv.FormatUint(transaction.TransactionKey, 10)
	paramsFile["ref_fk_domain"] = "tr_transaction"
	paramsFile["rec_status"] = "1"
	var filez []models.MsFile
	_, err = models.GetAllMsFile(&filez, 100, 100, paramsFile, true)
	// // log.Println("========== jumlah record ==========", len(filez))
	if len(filez) > 0 {
		// log.Println("========== sudah upload bukti trf ==========")
		for _, fl := range filez {
			aa := config.ImageUrl + *fl.FilePath
			responseData.UrlUpload = append(responseData.UrlUpload, &aa)
		}
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = responseData

	return c.JSON(http.StatusOK, response)
}

func ProsesPostingAll(c echo.Context) error {
	zero := decimal.NewFromInt(0)

	var err error
	var status int

	paramsSearch := make(map[string]string)

	transTypeKey := c.FormValue("trans_type_key")
	if transTypeKey == "" {
		// log.Error("Missing required parameter: trans_type_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: trans_type_key", "Missing required parameter: trans_type_key")
	}

	navDate := c.FormValue("nav_date")
	if navDate == "" {
		// log.Error("Missing required parameter: nav_date")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: nav_date", "Missing required parameter: nav_date")
	}
	paramsSearch["nav_date"] = navDate
	paramsSearch["rec_status"] = "1"
	paramsSearch["trans_status_key"] = "8"

	var trTransaction []models.TrTransaction
	var transTypeIds []string

	if transTypeKey == "3" {
		paramsSearch["orderBy"] = "trans_type_key"
		paramsSearch["orderType"] = "ASC"
		transTypeIds = append(transTypeIds, "3")
		transTypeIds = append(transTypeIds, "4")
		status, err = models.AdminGetAllTrTransactionPosting(&trTransaction, paramsSearch, transTypeIds, "trans_type_key", true)
		if err != nil {
			// log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
		if len(trTransaction) < 1 {
			// log.Error("transaction switching not found")
			return lib.CustomError(http.StatusNotFound, "Transaction not found", "Transaction not found")
		}
	} else {
		paramsSearch["orderBy"] = "transaction_key"
		paramsSearch["orderType"] = "ASC"
		transTypeIds = append(transTypeIds, transTypeKey)
		status, err = models.AdminGetAllTrTransactionPosting(&trTransaction, paramsSearch, transTypeIds, "trans_type_key", false)
		if err != nil {
			// log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
		if len(trTransaction) < 1 {
			// log.Error("transaction sub/redm not found")
			return lib.CustomError(http.StatusNotFound, "Transaction not found", "Transaction not found")
		}
	}

	dateLayout := "2006-01-02 15:04:05"
	strIDUserLogin := strconv.FormatUint(lib.Profile.UserID, 10)

	for _, transaction := range trTransaction {
		strTransTypeKey := strconv.FormatUint(transaction.TransTypeKey, 10)

		var transactionConf models.TrTransactionConfirmation
		strTransactionKey := strconv.FormatUint(transaction.TransactionKey, 10)
		_, err = models.GetTrTransactionConfirmationByTransactionKey(&transactionConf, strTransactionKey)
		if err != nil {
			// log.Error(err.Error())
			return lib.CustomError(http.StatusBadRequest)
		}

		var trBalanceCustomer []models.TrBalanceCustomerProduk
		strProductKey := strconv.FormatUint(transaction.ProductKey, 10)
		strCustomerKey := strconv.FormatUint(transaction.CustomerKey, 10)

		if strTransTypeKey == "2" || strTransTypeKey == "3" { // REDM
			_, err = models.GetLastBalanceCustomerByProductKey(&trBalanceCustomer, strCustomerKey, strProductKey)
			if err != nil {
				if err != sql.ErrNoRows {
					// log.Error("Transaction have not balance")
					return lib.CustomError(http.StatusBadRequest)
				} else {
					// log.Error(err.Error())
					return lib.CustomError(http.StatusBadRequest)
				}
			}
		}

		// strTransUnit := fmt.Sprintf("%g", transactionConf.ConfirmedUnit)

		//create tr_balance
		if (strTransTypeKey == "1") || (strTransTypeKey == "4") { // SUB & SWIN
			paramsBalance := make(map[string]string)
			strAcaKey := strconv.FormatUint(*transaction.AcaKey, 10)
			paramsBalance["aca_key"] = strAcaKey
			strTransactionConf := strconv.FormatUint(transactionConf.TcKey, 10)
			paramsBalance["tc_key"] = strTransactionConf

			newlayout := "2006-01-02"
			t, _ := time.Parse(dateLayout, transactionConf.ConfirmDate)
			balanceDate := t.Format(newlayout)

			paramsBalance["balance_date"] = balanceDate + " 00:00:00"
			paramsBalance["balance_unit"] = transactionConf.ConfirmedUnit.String()
			paramsBalance["rec_order"] = "0"
			paramsBalance["rec_status"] = "1"
			paramsBalance["rec_created_date"] = time.Now().Format(dateLayout)
			paramsBalance["rec_created_by"] = strIDUserLogin

			//calculate avg_nag tr_balance
			//sum balance unit

			balanceUnitSum := zero
			for _, trBalance := range trBalanceCustomer {
				balanceUnitSum = balanceUnitSum.Add(trBalance.BalanceUnit)
			}

			//avg nav balance last
			avgNavLast := zero
			var avgNav models.AvgNav
			_, err = models.GetLastAvgNavTrBalanceCustomerByProductKey(&avgNav, strCustomerKey, strProductKey)
			if err == nil {
				if avgNav.AvgNav != nil {
					avgNavLast = *avgNav.AvgNav
				}
			}

			variable1 := balanceUnitSum.Mul(avgNavLast)
			variable2 := transactionConf.ConfirmedUnit.Mul(*transactionConf.AvgNav)

			balanceUnitSumAll := balanceUnitSum.Add(transactionConf.ConfirmedUnit)

			countAvgNavBalance := variable1.Add(variable2).Div(balanceUnitSumAll)
			// strAvgNav := fmt.Sprintf("%g", countAvgNavBalance)
			paramsBalance["avg_nav"] = countAvgNavBalance.String()
			//end calculate avg_nag tr_balance

			status, err := models.CreateTrBalance(paramsBalance)
			if err != nil {
				// log.Error(err.Error())
				return lib.CustomError(status, err.Error(), "Failed input data")
			}
		}

		if (strTransTypeKey == "2") || (strTransTypeKey == "3") { // REDM & SWOUT
			sisaFifo := transactionConf.ConfirmedUnit
			//avg nav balance last
			avgNavLast := zero
			var avgNav models.AvgNav
			_, err = models.GetLastAvgNavTrBalanceCustomerByProductKey(&avgNav, strCustomerKey, strProductKey)
			if err == nil {
				if avgNav.AvgNav != nil {
					avgNavLast = *avgNav.AvgNav
				}
			}
			// strAvgNav := fmt.Sprintf("%g", avgNavLast)

			for _, trBalance := range trBalanceCustomer {
				if sisaFifo.Cmp(zero) == 1 {
					var sisaBalance decimal.Decimal

					if trBalance.BalanceUnit.Cmp(sisaFifo) == 1 {
						sisaBalance = trBalance.BalanceUnit.Sub(sisaFifo)
						sisaFifo = zero
					}

					if trBalance.BalanceUnit.Cmp(sisaFifo) == -1 {
						sisaBalance = zero
						sisaFifo = sisaFifo.Sub(trBalance.BalanceUnit)
					}

					if trBalance.BalanceUnit == sisaFifo {
						sisaBalance = zero
						sisaFifo = zero
					}

					paramsBalance := make(map[string]string)
					strAcaKey := strconv.FormatUint(*transaction.AcaKey, 10)
					paramsBalance["aca_key"] = strAcaKey
					strTransactionSubs := strconv.FormatUint(trBalance.TcKey, 10)
					paramsBalance["tc_key"] = strTransactionSubs
					strTransactionRed := strconv.FormatUint(transactionConf.TcKey, 10)
					paramsBalance["tc_key_red"] = strTransactionRed

					newlayout := "2006-01-02"
					t, _ := time.Parse(dateLayout, transactionConf.ConfirmDate)
					balanceDate := t.Format(newlayout)

					// strTransUnitSisa := fmt.Sprintf("%g", sisaBalance)

					paramsBalance["balance_date"] = balanceDate + " 00:00:00"
					paramsBalance["balance_unit"] = sisaBalance.String()
					paramsBalance["avg_nav"] = avgNavLast.String()

					var balance models.TrBalance
					status, err = models.GetLastTrBalanceByTcRed(&balance, strTransactionRed)
					if err != nil {
						paramsBalance["rec_order"] = "0"
					} else {
						if balance.RecOrder == nil {
							paramsBalance["rec_order"] = "0"
						} else {
							orderNext := int(*balance.RecOrder) + 1
							strOrderNext := strconv.FormatInt(int64(orderNext), 10)
							paramsBalance["rec_order"] = strOrderNext
						}
					}

					paramsBalance["rec_status"] = "1"
					paramsBalance["rec_created_date"] = time.Now().Format(dateLayout)
					paramsBalance["rec_created_by"] = strIDUserLogin
					status, err := models.CreateTrBalance(paramsBalance)
					if err != nil {
						// log.Error(err.Error())
						return lib.CustomError(status, err.Error(), "Failed input data")
					}
				} else {
					break
				}
			}
		}

		//update tr_transaction
		params := make(map[string]string)
		params["posted_units"] = transactionConf.ConfirmedUnit.String()
		params["trans_status_key"] = "9"
		params["transaction_key"] = strTransactionKey
		params["posted_date"] = time.Now().Format(dateLayout)
		params["rec_modified_by"] = strIDUserLogin
		params["rec_modified_date"] = time.Now().Format(dateLayout)

		_, err = models.UpdateTrTransaction(params)
		if err != nil {
			// log.Error("Error update tr transaction")
			return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed update data")
		}

		//create user message
		if strTransTypeKey != "3" { // SELAIN SWITCH-OUT
			var customer models.MsCustomer
			strCustomerKey := strconv.FormatUint(transaction.CustomerKey, 10)
			status, err = models.GetMsCustomer(&customer, strCustomerKey)
			if err == nil {
				if customer.InvestorType == "263" { //individu
					paramsUserMessage := make(map[string]string)
					paramsUserMessage["umessage_type"] = "245"

					var userLogin models.ScUserLogin
					_, err = models.GetScUserLoginByCustomerKey(&userLogin, strCustomerKey)
					if err != nil {
						// log.Error(err.Error())
						return lib.CustomError(http.StatusBadRequest)
					}

					strUserLoginKey := strconv.FormatUint(userLogin.UserLoginKey, 10)
					paramsUserMessage["umessage_recipient_key"] = strUserLoginKey
					paramsUserMessage["umessage_receipt_date"] = time.Now().Format(dateLayout)
					paramsUserMessage["flag_read"] = "0"
					paramsUserMessage["umessage_sender_key"] = strIDUserLogin
					paramsUserMessage["umessage_sent_date"] = time.Now().Format(dateLayout)
					paramsUserMessage["flag_sent"] = "1"
					var subject string
					var body string
					if strTransTypeKey == "1" { // SUBS
						if transaction.FlagNewSub != nil {
							if *transaction.FlagNewSub == 1 {
								subject = "Subscription Berhasil"
								body = "Subscription kamu telah efektif dibukukan. Silakan cek portofolio di akun kamu untuk melihat transaksi."
								paramsUserMessage["umessage_subject"] = subject
								paramsUserMessage["umessage_body"] = body
							} else {
								subject = "Top Up Berhasil"
								body = "Top Up kamu telah efektif dibukukan. Silakan cek portofolio di akun kamu untuk melihat transaksi."
								paramsUserMessage["umessage_subject"] = subject
								paramsUserMessage["umessage_body"] = body
							}
						} else {
							subject = "Top Up Berhasil"
							body = "Top Up kamu telah efektif dibukukan. Silakan cek portofolio di akun kamu untuk melihat transaksi."
							paramsUserMessage["umessage_subject"] = subject
							paramsUserMessage["umessage_body"] = body
						}
					}

					if strTransTypeKey == "2" { // REDM
						subject = "Redemption Berhasil"
						body = "Redemption kamu telah berhasil dijalankan. Dana akan ditransfer ke rekening bank kamu maks. 7 hari bursa. Silakan cek portofolio di akun kamu untuk melihat transaksi."
						paramsUserMessage["umessage_subject"] = subject
						paramsUserMessage["umessage_body"] = body
					}
					if strTransTypeKey == "4" { // SWITCH
						subject = "Switching Berhasil"
						body = "Switching kamu telah berhasil dijalankan. Silakan cek portofolio di akun kamu untuk melihat transaksi."
						paramsUserMessage["umessage_subject"] = subject
						paramsUserMessage["umessage_body"] = body
					}

					paramsUserMessage["umessage_category"] = "248"
					paramsUserMessage["flag_archieved"] = "0"
					paramsUserMessage["archieved_date"] = time.Now().Format(dateLayout)
					paramsUserMessage["rec_status"] = "1"
					paramsUserMessage["rec_created_date"] = time.Now().Format(dateLayout)
					paramsUserMessage["rec_created_by"] = strIDUserLogin

					status, err = models.CreateScUserMessage(paramsUserMessage)
					if err != nil {
						// log.Error("Error create user message")
						return lib.CustomError(status, err.Error(), "failed input data user message")
					}
					lib.CreateNotifCustomerFromAdminByUserLoginKey(strUserLoginKey, subject, body, "TRANSACTION")

					sendEmailTransactionPosted(transaction, transactionConf, userLogin, strCustomerKey, strTransTypeKey)
				} else if customer.InvestorType == "264" { //institusi
					SentEmailTransactionInstitutionPostingBackOfficeToUserCcSales(strTransactionKey, strCustomerKey)
				}
			}
		}

		// log.Info("Success update transaksi")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = ""
	return c.JSON(http.StatusOK, response)
}

func TestSentEmail(c echo.Context) error {
	trKey := c.FormValue("tr_key")
	if trKey == "" {
		// log.Error("Missing required parameter: tr_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: tr_key", "Missing required parameter: tr_key")
	}
	role := c.FormValue("role")
	if role == "" {
		// log.Error("Missing required parameter: role")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: role", "Missing required parameter: role")
	}

	// SentEmailTransactionToBackOffice(trKey, role)
	// SentEmailTransactionToBackOfficeAndSales(trKey, role)
	SentEmailTransactionToBackOfficeWithDb(trKey, role)

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = nil
	return c.JSON(http.StatusOK, response)
}

func SentEmailTransactionToBackOffice(transactionKey string, roleKey string) {
	var err error
	var transaction models.DetailTransactionDataSentEmail
	_, err = models.AdminDetailTransactionDataSentEmail(&transaction, transactionKey)
	if err != nil {
		// log.Error("Failed get transaction: " + err.Error())
		return
	}
	var mailTemp, subject string
	mailParam := make(map[string]string)
	if roleKey == "11" {
		mailParam["BackOfficeGroup"] = "Customer Service"
	} else if roleKey == "12" {
		mailParam["BackOfficeGroup"] = "Compliance"
	} else if roleKey == "13" {
		mailParam["BackOfficeGroup"] = "FundAdmin"
	}

	mailParam["FileUrl"] = config.ImageUrl + "/images/mail"
	mailParam["NamaLengkap"] = transaction.FullName
	mailParam["CIF"] = *transaction.Cif
	mailParam["TanggalTransaksi"] = transaction.TransDate
	mailParam["WaktuTransaksi"] = transaction.TransTime
	mailParam["Sales"] = *transaction.Sales
	ac0 := accounting.Accounting{Symbol: "", Precision: 0, Thousand: ".", Decimal: ","}
	if *transaction.EntryMode == uint64(140) { //Amount
		mailParam["JumlahTransaksi"] = transaction.CurrencySymbol + ". " + ac0.FormatMoneyDecimal(transaction.TransAmount.Truncate(0))
	} else { //Unit
		mailParam["JumlahTransaksi"] = ac0.FormatMoneyDecimal(transaction.TransUnit.Truncate(2)) + " Unit"
	}

	mailParam["BiayaTransaksi"] = transaction.CurrencySymbol + ". " + ac0.FormatMoneyDecimal(transaction.Fee.Truncate(0))

	if transaction.TransTypeKey == uint64(1) { // subs
		subject = "[MotionFunds] Mohon Verifikasi Transaksi Subscription"

		mailParam["TipeTransaksi"] = "Subscription"
		mailParam["NamaProduk"] = transaction.ProductName
		mailParam["MetodePembayaran"] = *transaction.PaymentMethodName
		mailParam["RekeningBankKustodian"] = *transaction.RekBankCustodian
		if *transaction.PaymentMethod == uint64(284) { //manual
			// log.Println("MANUAL TRANSFER")
			mailTemp = "email-new-subs-to-cs-kyc-fundadmin.html"
			var trDef models.TrTransaction
			_, err := models.GetTrTransaction(&trDef, transactionKey)
			if err == nil {
				if trDef.RecApprovalStage == nil {
					if transaction.BuktiTransafer != nil {
						mailParam["BuktiTransfer"] = config.ImageUrl + "/images/user/" + transaction.UserLoginKey + "/transfer/" + *transaction.BuktiTransafer
					} else {
						mailParam["BuktiTransfer"] = ""
					}
				} else {
					if transaction.BuktiTransafer != nil {
						mailParam["BuktiTransfer"] = config.ImageUrl + "/images/user/institusi/" + strconv.FormatUint(trDef.CustomerKey, 10) + "/transfer/" + *transaction.BuktiTransafer
					} else {
						mailParam["BuktiTransfer"] = ""
					}
				}
			}
		} else {
			// log.Println("NON MANUAL TRANSFER")
			mailTemp = "email-new-subs-to-cs-kyc-fundadmin-non-manual-transfer.html"
			mailParam["BuktiTransfer"] = "-"
		}
	} else if transaction.TransTypeKey == uint64(13) { // top up
		subject = "[MotionFunds] Mohon Verifikasi Transaksi Subscription"

		mailParam["TipeTransaksi"] = "Top Up"
		mailParam["NamaProduk"] = transaction.ProductName
		mailParam["MetodePembayaran"] = *transaction.PaymentMethodName
		mailParam["RekeningBankKustodian"] = *transaction.RekBankCustodian
		if *transaction.PaymentMethod == uint64(284) { //manual
			// log.Println("MANUAL TRANSFER")
			mailTemp = "email-new-subs-to-cs-kyc-fundadmin.html"
			var trDef models.TrTransaction
			_, err := models.GetTrTransaction(&trDef, transactionKey)
			if err == nil {
				if trDef.RecApprovalStage == nil {
					if transaction.BuktiTransafer != nil {
						mailParam["BuktiTransfer"] = config.ImageUrl + "/images/user/" + transaction.UserLoginKey + "/transfer/" + *transaction.BuktiTransafer
					} else {
						mailParam["BuktiTransfer"] = ""
					}
				} else {
					if transaction.BuktiTransafer != nil {
						mailParam["BuktiTransfer"] = config.ImageUrl + "/images/user/institusi/" + strconv.FormatUint(trDef.CustomerKey, 10) + "/transfer/" + *transaction.BuktiTransafer
					} else {
						mailParam["BuktiTransfer"] = ""
					}
				}
			}
		} else {
			// log.Println("NON MANUAL TRANSFER")
			mailTemp = "email-new-subs-to-cs-kyc-fundadmin-non-manual-transfer.html"
			mailParam["BuktiTransfer"] = "-"
		}
	} else if transaction.TransTypeKey == uint64(2) { // redm
		subject = "[MotionFunds] Mohon Verifikasi Transaksi Redemption"
		mailTemp = "email-new-redm-to-cs-kyc-fundadmin.html"

		mailParam["TipeTransaksi"] = "Redemption"
		mailParam["NamaProduk"] = transaction.ProductName
		mailParam["NamaBank"] = *transaction.BankRekBankCustomer
		mailParam["NoRekeningBank"] = *transaction.NoRekBankCustomer
		mailParam["NamaPadaRekeningBank"] = *transaction.NameRekBankCustomer
		mailParam["Cabang"] = *transaction.CabangRekBankCustomer
	} else { //switching
		subject = "[MotionFunds] Mohon Verifikasi Transaksi Switching"
		mailTemp = "email-new-switching-to-cs-kyc-fundadmin.html"

		mailParam["TipeTransaksi"] = "Switching"
		mailParam["ProdukAsal"] = transaction.ProductName
		mailParam["ProdukTujuan"] = *transaction.ProductTujuan
	}

	paramsScLogin := make(map[string]string)
	paramsScLogin["role_key"] = roleKey
	paramsScLogin["rec_status"] = "1"
	var userLogin []models.ScUserLogin
	_, err = models.GetAllScUserLogin(&userLogin, 0, 0, paramsScLogin, true)
	if err != nil {
		// log.Println("Email BO kosong")
		// log.Error("User BO tidak ditemukan")
		// log.Error(err)
	} else {
		// log.Println("Data User BO tersedia")
		// log.Println(len(userLogin))
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
						} else {
							// log.Println("Sukses email BO : " + scLogin.UloginEmail)
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
						// } else {
						// 	// log.Println("Sukses email BO : " + scLogin.UloginEmail)
						// }
					}
				}
			}
		}
	}
}

func SentEmailTransactionToBackOfficeAndSales(transactionKey string, roleKey string) {
	var err error
	var transaction models.DetailTransactionDataSentEmail
	_, err = models.AdminDetailTransactionDataSentEmail(&transaction, transactionKey)
	if err != nil {
		// log.Error("Failed get transaction: " + err.Error())
		return
	}
	var mailTemp, subject string
	mailParam := make(map[string]string)
	if roleKey == "11" {
		mailParam["BackOfficeGroup"] = "Customer Service"
	} else if roleKey == "12" {
		mailParam["BackOfficeGroup"] = "Compliance"
	} else if roleKey == "13" {
		mailParam["BackOfficeGroup"] = "FundAdmin"
	}

	mailParam["FileUrl"] = config.ImageUrl + "/images/mail"
	mailParam["NamaLengkap"] = transaction.FullName
	mailParam["CIF"] = *transaction.Cif
	mailParam["TanggalTransaksi"] = transaction.TransDate
	mailParam["WaktuTransaksi"] = transaction.TransTime
	mailParam["Sales"] = *transaction.Sales
	ac0 := accounting.Accounting{Symbol: "", Precision: 0, Thousand: ".", Decimal: ","}
	if *transaction.EntryMode == uint64(140) { //Amount
		mailParam["JumlahTransaksi"] = transaction.CurrencySymbol + ". " + ac0.FormatMoneyDecimal(transaction.TransAmount.Truncate(0))
	} else { //Unit
		mailParam["JumlahTransaksi"] = ac0.FormatMoneyDecimal(transaction.TransUnit.Truncate(2)) + " Unit"
	}

	mailParam["BiayaTransaksi"] = transaction.CurrencySymbol + ". " + ac0.FormatMoneyDecimal(transaction.Fee.Truncate(0))

	if transaction.TransTypeKey == uint64(1) { // subs
		subject = "[MotionFunds] Mohon Verifikasi Transaksi Subscription"

		mailParam["TipeTransaksi"] = "Subscription"
		mailParam["NamaProduk"] = transaction.ProductName
		mailParam["MetodePembayaran"] = *transaction.PaymentMethodName
		mailParam["RekeningBankKustodian"] = *transaction.RekBankCustodian
		if *transaction.PaymentMethod == uint64(284) { //manual
			// log.Println("MANUAL TRANSFER")
			mailTemp = "email-new-subs-to-cs-kyc-fundadmin.html"
			var trDef models.TrTransaction
			_, err := models.GetTrTransaction(&trDef, transactionKey)
			if err == nil {
				if trDef.RecApprovalStage == nil {
					if transaction.BuktiTransafer != nil {
						mailParam["BuktiTransfer"] = config.ImageUrl + "/images/user/" + transaction.UserLoginKey + "/transfer/" + *transaction.BuktiTransafer
					} else {
						mailParam["BuktiTransfer"] = ""
					}
				} else {
					if transaction.BuktiTransafer != nil {
						mailParam["BuktiTransfer"] = config.ImageUrl + "/images/user/institusi/" + strconv.FormatUint(trDef.CustomerKey, 10) + "/transfer/" + *transaction.BuktiTransafer
					} else {
						mailParam["BuktiTransfer"] = ""
					}
				}
			}
		} else {
			// log.Println("NON MANUAL TRANSFER")
			mailTemp = "email-new-subs-to-cs-kyc-fundadmin-non-manual-transfer.html"
			mailParam["BuktiTransfer"] = "-"
		}
	} else if transaction.TransTypeKey == uint64(2) { // redm
		subject = "[MotionFunds] Mohon Verifikasi Transaksi Redemption"
		mailTemp = "email-new-redm-to-cs-kyc-fundadmin.html"

		mailParam["TipeTransaksi"] = "Redemption"
		mailParam["NamaProduk"] = transaction.ProductName
		mailParam["NamaBank"] = *transaction.BankRekBankCustomer
		mailParam["NoRekeningBank"] = *transaction.NoRekBankCustomer
		mailParam["NamaPadaRekeningBank"] = *transaction.NameRekBankCustomer
		mailParam["Cabang"] = *transaction.CabangRekBankCustomer
	} else { //switching
		subject = "[MotionFunds] Mohon Verifikasi Transaksi Switching"
		mailTemp = "email-new-switching-to-cs-kyc-fundadmin.html"

		mailParam["TipeTransaksi"] = "Switching"
		mailParam["ProdukAsal"] = transaction.ProductName
		mailParam["ProdukTujuan"] = *transaction.ProductTujuan
	}

	paramsScLogin := make(map[string]string)
	paramsScLogin["role_key"] = roleKey
	paramsScLogin["rec_status"] = "1"
	var userLogin []models.ScUserLogin
	_, err = models.GetAllScUserLogin(&userLogin, 0, 0, paramsScLogin, true)
	if err != nil {
		// log.Println("Email BO kosong")
		// log.Error("User BO tidak ditemukan")
		// log.Error(err)
	} else {
		// log.Println("Data User BO tersedia")
		// log.Println(len(userLogin))
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
						} else {
							// log.Println("Sukses email BO : " + scLogin.UloginEmail)
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
						// } else {
						// 	// log.Println("Sukses email BO : " + scLogin.UloginEmail)
						// }
					}
				}
			}
		}
	}

	// to sales
	if transaction.SalesEmail != nil {
		var mailTempSales, subject string
		if transaction.TransTypeKey == uint64(1) { // subs
			subject = "[MotionFunds] Mohon Verifikasi Transaksi Subscription"
			if *transaction.PaymentMethod == uint64(284) { //manual
				// log.Println("MANUAL TRANSFER")
				mailTempSales = "email-new-subs-to-sales.html"
			} else {
				// log.Println("NON MANUAL TRANSFER")
				mailTempSales = "email-new-subs-to-sales-non-manual-transfer.html"
			}
		} else if transaction.TransTypeKey == uint64(2) { // redm
			subject = "[MotionFunds] Mohon Verifikasi Transaksi Redemption"
			mailTempSales = "email-new-redm-to-sales.html"
		} else { //switching
			subject = "[MotionFunds] Mohon Verifikasi Transaksi Switching"
			mailTempSales = "email-new-switching-to-sales.html"
		}

		t := template.New(mailTempSales)

		t, err = t.ParseFiles(config.BasePath + "/mail/" + mailTempSales)
		if err != nil {
			// log.Error("Failed send mail: " + err.Error())
			return
		}
		var tpl bytes.Buffer
		if err := t.Execute(&tpl, mailParam); err != nil {
			// log.Error("Failed send mail: " + err.Error())
			return
		}
		result := tpl.String()

		mailer := gomail.NewMessage()
		// mailer.SetHeader("From", config.EmailFrom)
		mailer.SetHeader("To", *transaction.SalesEmail)
		mailer.SetHeader("Subject", subject)
		mailer.SetBody("text/html", result)

		err = lib.SendEmail(mailer)
		if err != nil {
			// log.Error("Failed send mail to sales: " + *transaction.SalesEmail)
			// log.Error("Failed send mail to sales: " + err.Error())
		} else {
			// log.Println("Sukses email Sales : " + *transaction.SalesEmail)
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
		// 	// log.Error("Failed send mail to sales: " + *transaction.SalesEmail)
		// 	// log.Error("Failed send mail to sales: " + err.Error())
		// } else {
		// 	// log.Println("Sukses email Sales : " + *transaction.SalesEmail)
		// }
	} else {
		// log.Println("Data Sales tidak ada email")
	}
}

func SentEmailTransactionRejectToSales(transactionKey string, notes string) {
	var err error
	var transaction models.DetailTransactionDataSentEmail
	_, err = models.AdminDetailTransactionDataSentEmail(&transaction, transactionKey)
	if err != nil {
		// log.Error("Failed get transaction: " + err.Error())
		return
	}
	if transaction.SalesEmail != nil {
		var mailTemp, subject string
		mailParam := make(map[string]string)

		mailParam["FileUrl"] = config.ImageUrl + "/images/mail"
		mailParam["NamaLengkap"] = transaction.FullName
		mailParam["CIF"] = *transaction.Cif
		mailParam["TanggalTransaksi"] = transaction.TransDate
		mailParam["WaktuTransaksi"] = transaction.TransTime
		mailParam["Sales"] = *transaction.Sales
		ac0 := accounting.Accounting{Symbol: "", Precision: 0, Thousand: ".", Decimal: ","}
		if *transaction.EntryMode == uint64(140) { //Amount
			mailParam["JumlahTransaksi"] = transaction.CurrencySymbol + ". " + ac0.FormatMoneyDecimal(transaction.TransAmount.Truncate(0))
		} else { //Unit
			mailParam["JumlahTransaksi"] = ac0.FormatMoneyDecimal(transaction.TransUnit.Truncate(2)) + " Unit"
		}

		mailParam["BiayaTransaksi"] = transaction.CurrencySymbol + ". " + ac0.FormatMoneyDecimal(transaction.Fee.Truncate(0))

		if transaction.TransTypeKey == uint64(1) { // subs
			subject = "[MotionFunds] Transaksi Subscription Ditolak"

			mailParam["TipeTransaksi"] = "Subscription"
			mailParam["NamaProduk"] = transaction.ProductName
			mailParam["MetodePembayaran"] = *transaction.PaymentMethodName
			mailParam["RekeningBankKustodian"] = *transaction.RekBankCustodian
			if *transaction.PaymentMethod == uint64(284) { //manual
				// log.Println("MANUAL TRANSFER")
				mailTemp = "email-subs-rejected-to-sales.html"

				var trDef models.TrTransaction
				_, err := models.GetTrTransaction(&trDef, transactionKey)
				if err == nil {
					if trDef.RecApprovalStage == nil {
						if transaction.BuktiTransafer != nil {
							mailParam["BuktiTransfer"] = config.ImageUrl + "/images/user/" + transaction.UserLoginKey + "/transfer/" + *transaction.BuktiTransafer
						} else {
							mailParam["BuktiTransfer"] = ""
						}
					} else {
						if transaction.BuktiTransafer != nil {
							mailParam["BuktiTransfer"] = config.ImageUrl + "/images/user/institusi/" + strconv.FormatUint(trDef.CustomerKey, 10) + "/transfer/" + *transaction.BuktiTransafer
						} else {
							mailParam["BuktiTransfer"] = ""
						}
					}
				}
			} else {
				// log.Println("NON MANUAL TRANSFER")
				mailTemp = "email-subs-rejected-to-sales-non-manual-transfer.html"
				mailParam["BuktiTransfer"] = "-"
			}
		} else if transaction.TransTypeKey == uint64(2) { // redm
			subject = "[MotionFunds] Transaksi Redemption Ditolak"
			mailTemp = "email-redm-rejected-to-sales.html"

			mailParam["TipeTransaksi"] = "Redemption"
			mailParam["NamaProduk"] = transaction.ProductName
			mailParam["NamaBank"] = *transaction.BankRekBankCustomer
			mailParam["NoRekeningBank"] = *transaction.NoRekBankCustomer
			mailParam["NamaPadaRekeningBank"] = *transaction.NameRekBankCustomer
			mailParam["Cabang"] = *transaction.CabangRekBankCustomer
		} else { //switching
			subject = "[MotionFunds] Transaksi Switching Ditolak"
			mailTemp = "email-switching-rejected-to-sales.html"

			mailParam["TipeTransaksi"] = "Switching"
			mailParam["ProdukAsal"] = transaction.ProductName
			mailParam["ProdukTujuan"] = *transaction.ProductTujuan
		}
		mailParam["Status"] = "Ditolak"
		mailParam["Alasan"] = notes

		t := template.New(mailTemp)

		t, err = t.ParseFiles(config.BasePath + "/mail/" + mailTemp)
		if err != nil {
			// log.Error("Failed send mail: " + err.Error())
			return
		}
		var tpl bytes.Buffer
		if err := t.Execute(&tpl, mailParam); err != nil {
			// log.Error("Failed send mail: " + err.Error())
			return
		}
		result := tpl.String()

		mailer := gomail.NewMessage()
		// mailer.SetHeader("From", config.EmailFrom)
		mailer.SetHeader("To", *transaction.SalesEmail)
		mailer.SetHeader("Subject", subject)
		mailer.SetBody("text/html", result)

		err = lib.SendEmail(mailer)
		if err != nil {
			// log.Error("Failed send mail to sales: " + *transaction.SalesEmail)
			// log.Error("Failed send mail to sales: " + err.Error())
		} else {
			// log.Println("Sukses email Sales : " + *transaction.SalesEmail)
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
		// 	// log.Error("Failed send mail to sales: " + *transaction.SalesEmail)
		// 	// log.Error("Failed send mail to sales: " + err.Error())
		// } else {
		// 	// log.Println("Sukses email Sales : " + *transaction.SalesEmail)
		// }
	} else {
		// log.Println("Data Sales tidak ada email")
	}
}

func SentEmailTransactionToBackOfficeWithDb(transactionKey string, roleKey string) {
	var err error
	var transaction models.DetailTransactionDataSentEmail
	_, err = models.AdminDetailTransactionDataSentEmail(&transaction, transactionKey)
	if err != nil {
		// log.Error("Failed get transaction: " + err.Error())
		return
	}
	var mailTemp, subject string
	mailParam := make(map[string]string)
	if roleKey == "11" {
		mailParam["BackOfficeGroup"] = "Customer Service"
	} else if roleKey == "12" {
		mailParam["BackOfficeGroup"] = "Compliance"
	} else if roleKey == "13" {
		mailParam["BackOfficeGroup"] = "FundAdmin"
	}

	mailParam["FileUrl"] = config.ImageUrl + "/images/mail"
	mailParam["NamaLengkap"] = transaction.FullName
	mailParam["CIF"] = *transaction.Cif
	mailParam["TanggalTransaksi"] = transaction.TransDate
	mailParam["WaktuTransaksi"] = transaction.TransTime
	mailParam["Sales"] = *transaction.Sales
	ac0 := accounting.Accounting{Symbol: "", Precision: 0, Thousand: ".", Decimal: ","}
	if *transaction.EntryMode == uint64(140) { //Amount
		mailParam["JumlahTransaksi"] = transaction.CurrencySymbol + ". " + ac0.FormatMoneyDecimal(transaction.TransAmount.Truncate(0))
	} else { //Unit
		mailParam["JumlahTransaksi"] = ac0.FormatMoneyDecimal(transaction.TransUnit.Truncate(2)) + " Unit"
	}

	mailParam["BiayaTransaksi"] = transaction.CurrencySymbol + ". " + ac0.FormatMoneyDecimal(transaction.Fee.Truncate(0))

	if transaction.TransTypeKey == uint64(1) { // subs
		// subject = "[MotionFunds] Mohon Verifikasi Transaksi Subscription"

		mailParam["TipeTransaksi"] = "Subscription"
		mailParam["NamaProduk"] = transaction.ProductName
		mailParam["MetodePembayaran"] = *transaction.PaymentMethodName
		mailParam["RekeningBankKustodian"] = *transaction.RekBankCustodian
		if *transaction.PaymentMethod == uint64(284) { //manual
			// log.Println("MANUAL TRANSFER")
			mailTemp = "NEW-SUBS-TO-BACKOFFICE-MANUAL-TRANSFER"

			var trDef models.TrTransaction
			_, err := models.GetTrTransaction(&trDef, transactionKey)
			if err == nil {
				if trDef.RecApprovalStage == nil {
					if transaction.BuktiTransafer != nil {
						mailParam["BuktiTransfer"] = config.ImageUrl + "/images/user/" + transaction.UserLoginKey + "/transfer/" + *transaction.BuktiTransafer
					} else {
						mailParam["BuktiTransfer"] = ""
					}
				} else {
					if transaction.BuktiTransafer != nil {
						mailParam["BuktiTransfer"] = config.ImageUrl + "/images/user/institusi/" + strconv.FormatUint(trDef.CustomerKey, 10) + "/transfer/" + *transaction.BuktiTransafer
					} else {
						mailParam["BuktiTransfer"] = ""
					}
				}
			}
		} else {
			// log.Println("NON MANUAL TRANSFER")
			mailTemp = "NEW-SUBS-TO-BACKOFFICE-NON-MANUAL-TRANSFER"
			mailParam["BuktiTransfer"] = "-"
		}
	} else if transaction.TransTypeKey == uint64(2) { // redm
		mailTemp = "NEW-REDM-TO-BACKOFFICE"

		mailParam["TipeTransaksi"] = "Redemption"
		mailParam["NamaProduk"] = transaction.ProductName
		mailParam["NamaBank"] = *transaction.BankRekBankCustomer
		mailParam["NoRekeningBank"] = *transaction.NoRekBankCustomer
		mailParam["NamaPadaRekeningBank"] = *transaction.NameRekBankCustomer
		mailParam["Cabang"] = *transaction.CabangRekBankCustomer
	} else { //switching
		mailTemp = "NEW-SWITCHING-TO-BACKOFFICE"

		mailParam["TipeTransaksi"] = "Switching"
		mailParam["ProdukAsal"] = transaction.ProductName
		mailParam["ProdukTujuan"] = *transaction.ProductTujuan
	}

	var mail models.MmMailMaster
	_, err = models.GetMmMailMaster(&mail, "mail_template_name", mailTemp)
	if err != nil {
		// log.Error("Mail Template Name : " + mailTemp)
		// log.Error("Mail Template Name tidak di temukan: " + err.Error())
		return
	} else {
		subject = *mail.MailSubject

		paramsScLogin := make(map[string]string)
		paramsScLogin["role_key"] = roleKey
		paramsScLogin["rec_status"] = "1"
		var userLogin []models.ScUserLogin
		_, err = models.GetAllScUserLogin(&userLogin, 0, 0, paramsScLogin, true)
		if err != nil {
			// log.Println("Email BO kosong")
			// log.Error("User BO tidak ditemukan")
			// log.Error(err)
		} else {
			// log.Println("Data User BO tersedia")
			// log.Println(len(userLogin))
			t := template.New(mailTemp)

			t, err = t.Parse(*mail.MailBody)
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
							result = strings.ReplaceAll(result, "\r\n", "")
							result = strings.ReplaceAll(result, "\"", "")
							// log.Println(result)

							mailer := gomail.NewMessage()
							// mailer.SetHeader("From", config.EmailFrom)
							mailer.SetHeader("To", scLogin.UloginEmail)
							mailer.SetHeader("Subject", subject)
							mailer.SetBody("text/html", result)

							paramsLog := make(map[string]string)
							err = lib.SendEmail(mailer)
							if err != nil {
								// log.Error("Failed send mail to: " + scLogin.UloginEmail)
								// log.Error("Failed send mail: " + err.Error())
								paramsLog["job_error_log"] = err.Error()
							} else {
								// log.Println("Sukses email BO : " + scLogin.UloginEmail)
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
							// 	paramsLog["job_error_log"] = err.Error()
							// } else {
							// 	// log.Println("Sukses email BO : " + scLogin.UloginEmail)
							// }

							//save to mail log
							dateLayout := "2006-01-02 15:04:05"
							paramsLog["mail_master_key"] = strconv.FormatUint(mail.MailMasterKey, 10)
							paramsLog["mail_master_type"] = strconv.FormatUint(mail.MailMasterType, 10)
							if mail.MailMasterCategory != nil {
								paramsLog["mail_master_category"] = strconv.FormatUint(*mail.MailMasterCategory, 10)
							}
							paramsLog["mail_template_name"] = mail.MailTemplateName
							if mail.MailAccountKey != nil {
								paramsLog["mail_account_key"] = strconv.FormatUint(*mail.MailAccountKey, 10)
							}
							if mail.MailToGroupKey != nil {
								paramsLog["mail_to_group_key"] = strconv.FormatUint(*mail.MailToGroupKey, 10)
							}
							if mail.MailToMailKey != nil {
								paramsLog["mail_to_mail_key"] = strconv.FormatUint(*mail.MailToMailKey, 10)
							}
							if mail.MailCcGroupKey != nil {
								paramsLog["mail_cc_group_key"] = strconv.FormatUint(*mail.MailCcGroupKey, 10)
							}
							if mail.MailCcMailKey != nil {
								paramsLog["mail_cc_mail_key"] = strconv.FormatUint(*mail.MailCcMailKey, 10)
							}
							if mail.MailBccGroupKey != nil {
								paramsLog["mail_bcc_group_key"] = strconv.FormatUint(*mail.MailBccGroupKey, 10)
							}
							if mail.MailBccMailKey != nil {
								paramsLog["mail_bcc_mail_key"] = strconv.FormatUint(*mail.MailBccMailKey, 10)
							}
							paramsLog["mail_subject"] = subject
							paramsLog["mail_body"] = result
							paramsLog["mail_has_attachment"] = strconv.FormatUint(uint64(mail.MailHasAttachment), 10)
							paramsLog["mail_flag_html"] = strconv.FormatUint(uint64(mail.MailFlagHtml), 10)
							paramsLog["job_is_execute"] = "1"
							paramsLog["job_sent_date"] = time.Now().Format(dateLayout)
							paramsLog["job_sent_count"] = "1"
							paramsLog["job_execute_date"] = time.Now().Format(dateLayout)
							paramsLog["rec_created_date"] = time.Now().Format(dateLayout)
							paramsLog["rec_created_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
							paramsLog["rec_status"] = "1"
							_, err, _ := models.CreateMmMailSentLog(paramsLog)
							if err != nil {
								// log.Error(err.Error())
								// log.Error("Error create log email")
							} else {
								// log.Println("Success create log email")
							}
						}
					}
				}
			}
		}
	}
}

func SentEmailTransactionInstitutionRejectBackOfficeToUserCcSales(
	transactionKey string,
	customerKey string,
	notes string) {

	var err error

	var transaction models.DetailTransactionDataSentEmail
	_, err = models.AdminDetailTransactionDataSentEmail(&transaction, transactionKey)
	if err != nil {
		// log.Error("Failed get transaction: " + err.Error())
		return
	}
	var mailTemp, subject string
	mailParam := make(map[string]string)

	mailParam["FileUrl"] = config.ImageUrl + "/images/mail"

	mailParam["NamaLengkap"] = transaction.FullName
	mailParam["CIF"] = *transaction.Cif
	mailParam["TanggalTransaksi"] = transaction.TransDate
	mailParam["WaktuTransaksi"] = transaction.TransTime
	ac0 := accounting.Accounting{Symbol: "", Precision: 0, Thousand: ".", Decimal: ","}
	if *transaction.EntryMode == uint64(140) { //Amount
		mailParam["JumlahTransaksi"] = transaction.CurrencySymbol + ". " + ac0.FormatMoneyDecimal(transaction.TransAmount.Truncate(0))
	} else { //Unit
		mailParam["JumlahTransaksi"] = ac0.FormatMoneyDecimal(transaction.TransUnit.Truncate(2)) + " Unit"
	}

	mailParam["BiayaTransaksi"] = transaction.CurrencySymbol + ". " + ac0.FormatMoneyDecimal(transaction.Fee.Truncate(0))

	mailParam["Status"] = "Gagal"

	mailParam["Keterangan"] = notes

	if transaction.TransTypeKey == uint64(1) { // subs
		mailTemp = "email-transaksi-reject-subs-redm-institusi-cs-kyc.html"
		subject = "[MotionFunds] Transaksi Subscription Dibatalkan"

		mailParam["TypeTransaksi"] = "Subscription"
		if transaction.FlagNewSub != nil && *transaction.FlagNewSub == uint8(0) {
			mailParam["TypeTransaksi"] = "Topup"
		}
		mailParam["NamaProduk"] = transaction.ProductName
	} else if transaction.TransTypeKey == uint64(2) { // redm
		mailTemp = "email-transaksi-reject-subs-redm-institusi-cs-kyc.html"
		subject = "[MotionFunds] Transaksi Redemption Dibatalkan"

		mailParam["TypeTransaksi"] = "Redemption"
		mailParam["NamaProduk"] = transaction.ProductName
	} else { //switching
		mailTemp = "email-transaksi-reject-switching-institusi-cs-kyc.html"
		subject = "[MotionFunds] Transaksi Switching Dibatalkan"

		mailParam["TypeTransaksi"] = "Switching"
		mailParam["ProdukAsal"] = transaction.ProductName
		mailParam["ProdukTujuan"] = *transaction.ProductTujuan
	}

	paramsScLoginNext := make(map[string]string)
	paramsScLoginNext["customer_key"] = customerKey
	paramsScLoginNext["rec_status"] = "1"
	paramsScLoginNext["verified_email"] = "1"
	paramsScLoginNext["verified_mobileno"] = "1"
	paramsScLoginNext["ulogin_enabled"] = "1"
	var userTujuan []models.ScUserLogin
	_, err = models.GetAllScUserLogin(&userTujuan, 0, 0, paramsScLoginNext, true)
	if err == nil {
		// log.Println("Data User BO tersedia")
		t := template.New(mailTemp)

		t, err = t.ParseFiles(config.BasePath + "/mail/" + mailTemp)
		if err != nil {
			// log.Error("Failed send mail: " + err.Error())
		} else {
			var tpl bytes.Buffer
			if err := t.Execute(&tpl, mailParam); err != nil {
				// log.Error("Failed send mail: " + err.Error())
			} else {
				result := tpl.String()

				mailer := gomail.NewMessage()
				// mailer.SetHeader("From", config.EmailFrom)

				addresses := make([]string, len(userTujuan))
				for i, scLogin := range userTujuan {
					addresses[i] = mailer.FormatAddress(scLogin.UloginEmail, "")
				}

				mailer.SetHeader("To", addresses...)

				mailer.SetHeader("Subject", subject)
				mailer.SetBody("text/html", result)

				if transaction.SalesEmail != nil {
					mailer.SetAddressHeader("Cc", *transaction.SalesEmail, *transaction.Sales)
				}

				xPort, _ := strconv.ParseInt(config.MF_EmailSMTPPort, 10, 32)
				dialer := gomail.NewDialer(
					config.MF_EmailSMTPHost,
					int(xPort),
					config.MF_EmailFrom,
					config.MF_EmailFromPassword,
				)
				dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

				err = dialer.DialAndSend(mailer)
				if err != nil {
					// log.Error("Failed send mail to user institution")
					// log.Error("Failed send mail: " + err.Error())
				} else {
					// log.Println("Sukses email internal intitution : customer_key = " + customerKey)
				}
			}
		}
	}
}

func SentEmailTransactionInstitutionPostingBackOfficeToUserCcSales(
	transactionKey string,
	customerKey string) {

	var err error

	var transaction models.DetailTransactionDataSentEmail
	_, err = models.AdminDetailTransactionDataSentEmail(&transaction, transactionKey)
	if err != nil {
		// log.Error("Failed get transaction: " + err.Error())
		return
	}
	var mailTemp, subject string
	mailParam := make(map[string]string)

	mailParam["FileUrl"] = config.ImageUrl + "/images/mail"

	mailParam["NamaLengkap"] = transaction.FullName
	mailParam["CIF"] = *transaction.Cif
	mailParam["TanggalTransaksi"] = transaction.TransDate
	mailParam["WaktuTransaksi"] = transaction.TransTime
	ac0 := accounting.Accounting{Symbol: "", Precision: 0, Thousand: ".", Decimal: ","}
	if *transaction.EntryMode == uint64(140) { //Amount
		mailParam["JumlahTransaksi"] = transaction.CurrencySymbol + ". " + ac0.FormatMoneyDecimal(transaction.TransAmount.Truncate(0))
	} else { //Unit
		mailParam["JumlahTransaksi"] = ac0.FormatMoneyDecimal(transaction.TransUnit.Truncate(2)) + " Unit"
	}

	mailParam["BiayaTransaksi"] = transaction.CurrencySymbol + ". " + ac0.FormatMoneyDecimal(transaction.Fee.Truncate(0))

	mailParam["NABUnit"] = ac0.FormatMoneyDecimal(transaction.NavValue.Truncate(2))
	mailParam["TanggalNAB"] = transaction.NavDate
	var tc models.TrTransactionConfirmation
	_, err = models.GetTrTransactionConfirmationByTransactionKey(&tc, transactionKey)
	if err == nil {
		mailParam["UnitPenyertaan"] = ac0.FormatMoneyDecimal(tc.ConfirmedUnit.Truncate(2))
	} else {
		mailParam["UnitPenyertaan"] = "-"
	}

	if transaction.TransTypeKey == uint64(1) { // subs
		mailTemp = "email-transaksi-posting-subs-redm-institusi-cs-kyc.html"
		subject = "[MotionFunds] Transaksi Subscription Anda Berhasil"

		mailParam["TypeTransaksi"] = "Subscription"
		if transaction.FlagNewSub != nil && *transaction.FlagNewSub == uint8(0) {
			mailParam["TypeTransaksi"] = "Topup"
		}
		mailParam["NamaProduk"] = transaction.ProductName
	} else if transaction.TransTypeKey == uint64(2) { // redm
		mailTemp = "email-transaksi-posting-subs-redm-institusi-cs-kyc.html"
		subject = "[MotionFunds] Transaksi Redemption Anda Berhasil"

		mailParam["TypeTransaksi"] = "Redemption"
		mailParam["NamaProduk"] = transaction.ProductName
	} else { //switching
		mailTemp = "email-transaksi-posting-switching-institusi-cs-kyc.html"
		subject = "[MotionFunds] Transaksi Switching Anda Berhasil"

		mailParam["TypeTransaksi"] = "Switching"
		mailParam["ProdukAsal"] = transaction.ProductName
		mailParam["ProdukTujuan"] = *transaction.ProductTujuan
	}

	paramsScLoginNext := make(map[string]string)
	paramsScLoginNext["customer_key"] = customerKey
	paramsScLoginNext["rec_status"] = "1"
	paramsScLoginNext["verified_email"] = "1"
	paramsScLoginNext["verified_mobileno"] = "1"
	paramsScLoginNext["ulogin_enabled"] = "1"
	var userTujuan []models.ScUserLogin
	_, err = models.GetAllScUserLogin(&userTujuan, 0, 0, paramsScLoginNext, true)
	if err == nil {
		// log.Println("Data User BO tersedia")
		t := template.New(mailTemp)

		t, err = t.ParseFiles(config.BasePath + "/mail/" + mailTemp)
		if err != nil {
			// log.Error("Failed send mail: " + err.Error())
		} else {
			var tpl bytes.Buffer
			if err := t.Execute(&tpl, mailParam); err != nil {
				// log.Error("Failed send mail: " + err.Error())
			} else {
				result := tpl.String()

				mailer := gomail.NewMessage()
				// mailer.SetHeader("From", config.EmailFrom)

				addresses := make([]string, len(userTujuan))
				for i, scLogin := range userTujuan {
					addresses[i] = mailer.FormatAddress(scLogin.UloginEmail, "")
				}

				mailer.SetHeader("To", addresses...)

				mailer.SetHeader("Subject", subject)
				mailer.SetBody("text/html", result)

				if transaction.SalesEmail != nil {
					mailer.SetAddressHeader("Cc", *transaction.SalesEmail, *transaction.Sales)
				}

				err = lib.SendEmail(mailer)
				if err != nil {
					// log.Error("Failed send mail to user institution")
					// log.Error("Failed send mail: " + err.Error())
				} else {
					// log.Println("Sukses email internal intitution : customer_key = " + customerKey)
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
				// 	// log.Error("Failed send mail to user institution")
				// 	// log.Error("Failed send mail: " + err.Error())
				// } else {
				// 	// log.Println("Sukses email internal intitution : customer_key = " + customerKey)
				// }
			}
		}
	}
}

func ProsesCorrection(c echo.Context) error {
	var err error
	params := make(map[string]string)

	transactionkey := c.FormValue("transaction_key")
	if transactionkey == "" {
		// log.Error("Missing required parameter: transaction_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: transaction_key", "Missing required parameter: transaction_key")
	}

	n, err := strconv.ParseUint(transactionkey, 10, 64)
	if err == nil && n > 0 {
		params["transaction_key"] = transactionkey
	} else {
		// log.Error("Wrong input for parameter: transaction_key")
		return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: transaction_key", "Wrong input for parameter: transaction_key")
	}

	var transaction models.TrTransaction
	_, err = models.GetTrTransaction(&transaction, transactionkey)
	if err != nil {
		// log.Error(err.Error())
		return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: transaction_key", "Wrong input for parameter: transaction_key")
	}

	transStatusIds := []string{"2", "4", "5"}

	_, found := lib.Find(transStatusIds, strconv.FormatUint(transaction.TransStatusKey, 10))
	if !found {
		// log.Error("trans_status_key buka 2,4,5")
		return lib.CustomError(http.StatusBadRequest, "Transaction Not Allowed", "Transaction Not Allowed")
	}

	notes := c.FormValue("notes")
	if notes == "" {
		// log.Error("Missing required parameter notes: Notes tidak boleh kosong")
		return lib.CustomError(http.StatusBadRequest, "Notes tidak boleh kosong", "Notes tidak boleh kosong")
	} else {
		if len(notes) > 250 {
			// log.Error("Missing required parameter: notes too long, max 250 character")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: notes too long, max 250 character", "Missing required parameter: notes too long, max 250 character")
		}
	}

	dateLayout := "2006-01-02 15:04:05"
	strIDUserLogin := strconv.FormatUint(lib.Profile.UserID, 10)

	params["trans_status_key"] = "1"
	params["check2_notes"] = notes
	params["check2_date"] = time.Now().Format(dateLayout)
	params["check2_flag"] = "1"
	params["check2_references"] = strIDUserLogin

	params["rec_modified_by"] = strIDUserLogin
	params["rec_modified_date"] = time.Now().Format(dateLayout)

	_, err = models.UpdateTrTransaction(params)
	if err != nil {
		// log.Error("Error update tr transaction")
		return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed update data")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = ""
	return c.JSON(http.StatusOK, response)
}

func GetTransactionStampsAdmin(c echo.Context) error {

	// var cKey string
	var err error
	var status int
	var str_message string
	params := make(map[string]string)
	decimal.MarshalJSONWithoutQuotes = true

	customerKeystr := c.QueryParam("customer_key")
	if customerKeystr != "" {
		params["customer_key"] = customerKeystr
	} else {
		return lib.CustomError(http.StatusNotFound)
	}
	params["rec_status"] = "1"

	var appConfig models.ScAppConfig
	status, err = models.GetScAppConfigByCode(&appConfig, "TRX_STAMP_MIN_VALUE_IDR")
	if err != nil {
		str_message = err.Error()
		// log.Error(str_message)
		return lib.CustomError(http.StatusBadRequest, str_message, "Fail to get Config TRX_STAMP_MIN_VALUE_IDR")
	}
	min_trx_amount, _ := strconv.ParseInt(*appConfig.AppConfigValue, 10, 64)

	// log.Info("min_trx_amount %T, %d", min_trx_amount, min_trx_amount)
	paramsStamp := make(map[string]string)

	//paramsStamp["currency_key"] = "1" //mata uang materai
	var stampValues []models.StampDutyValue
	_, err = models.GetStampNominal(&stampValues, paramsStamp)
	if err != nil {
		str_message = err.Error()
		// log.Error(str_message)
		return lib.CustomError(http.StatusBadRequest, str_message, "Fail to get Stamp Nominal")
	}
	if len(stampValues) < 1 {
		str_message = err.Error()
		// log.Error(str_message)
		return lib.CustomError(http.StatusBadRequest, str_message, "Fail to get Stamp Nominal")
	}

	var stampData models.TransStampData
	status, err = models.GetTransactionStamps(&stampData, params)
	if err != nil {
		str_message = err.Error()
		// log.Error(str_message)
		return lib.CustomError(status, str_message, "Failed get StampDataInfo")
	}
	//nilai materai
	stampData.StampDutyValues = stampValues

	//stampData.CustomerKey = *lib.Profile.CustomerKey
	//stampData.NavDate = time.Now().Format(dateLayout)

	stampData.HasStamp = false
	if stampData.StampFeeAmount.Cmp(decimal.Zero) == 1 {
		stampData.HasStamp = true
		stampData.StampMessageInfo = fmt.Sprintf("Biaya Materai dikenakan jika akumulasi transaksi per hari melebihi Rp.%9s", lib.FormatNumber(min_trx_amount))
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = stampData
	return c.JSON(http.StatusOK, response)

}
