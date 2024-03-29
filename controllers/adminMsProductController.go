package controllers

import (
	"database/sql"
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

// func initAuthHoIt() error {
// 	var roleKeyHoIt uint64
// 	roleKeyHoIt = 15

// 	if lib.Profile.RoleKey != roleKeyHoIt {
// 		// log.Error("User Autorizer")
// 		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
// 	}
// 	return nil
// }

func GetListProductAdmin(c echo.Context) error {
	PAGE_MENU_KEY := 71
	var err error
	var status int
	decimal.MarshalJSONWithoutQuotes = true

	errorAuth := IsMenuAccessAllowed(PAGE_MENU_KEY)
	if errorAuth != nil {
		// log.Error("User Autorizer")
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}

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

	items := []string{"product_key", "product_code", "product_name", "launch_date", "inception_date", "isin_code", "flag_syariah", "sinvest_fund_code"}

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
		params["orderBy"] = "product_key"
		params["orderType"] = "ASC"
	}

	params["rec_status"] = "1"

	paramsLike := make(map[string]string)

	productName := c.QueryParam("product_name")
	if productName != "" {
		paramsLike["product_name"] = productName
	}
	productCode := c.QueryParam("product_code")
	if productCode != "" {
		paramsLike["product_code"] = productCode
	}
	isinCode := c.QueryParam("isin_code")
	if isinCode != "" {
		paramsLike["isin_code"] = isinCode
	}

	// var RiskProfileName string

	var msProduct []models.MsProduct

	status, err = models.AdminGetAllMsProductWithLike(&msProduct, limit, offset, params, paramsLike, noLimit)

	if err != nil {
		// log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if len(msProduct) < 1 {
		// log.Error("product not found")
		return lib.CustomError(http.StatusNotFound, "Product not found", "Product not found")
	}

	var currencyIds []string
	var productCategoryIds []string
	var productTypeIds []string
	var genLookupIds []string
	var custodianIds []string
	var riskProfileIds []string

	for _, pro := range msProduct {
		if pro.RiskProfileKey != nil {
			if _, ok := lib.Find(riskProfileIds, strconv.FormatUint(*pro.RiskProfileKey, 10)); !ok {
				riskProfileIds = append(riskProfileIds, strconv.FormatUint(*pro.RiskProfileKey, 10))
			}
			// var msRiskProfileData []models.MsRiskProfile
			// status, err = models.GetMsRiskProfileIn(&msRiskProfileData, riskProfileIds)
			// if len(msRiskProfileData) > 0 {
			// 	RiskProfileName = *msRiskProfileData[0].RiskName
			// }
		}

		if pro.CurrencyKey != nil {
			if _, ok := lib.Find(currencyIds, strconv.FormatUint(*pro.CurrencyKey, 10)); !ok {
				currencyIds = append(currencyIds, strconv.FormatUint(*pro.CurrencyKey, 10))
			}
		}

		if pro.ProductCategoryKey != nil {
			if _, ok := lib.Find(productCategoryIds, strconv.FormatUint(*pro.ProductCategoryKey, 10)); !ok {
				productCategoryIds = append(productCategoryIds, strconv.FormatUint(*pro.ProductCategoryKey, 10))
			}
		}

		if pro.ProductTypeKey != nil {
			if _, ok := lib.Find(productTypeIds, strconv.FormatUint(*pro.ProductTypeKey, 10)); !ok {
				productTypeIds = append(productTypeIds, strconv.FormatUint(*pro.ProductTypeKey, 10))
			}
		}

		if pro.CustodianKey != nil {
			if _, ok := lib.Find(custodianIds, strconv.FormatUint(*pro.CustodianKey, 10)); !ok {
				custodianIds = append(custodianIds, strconv.FormatUint(*pro.CustodianKey, 10))
			}
		}

	}

	//mapping currency
	var msCurrency []models.MsCurrency
	if len(currencyIds) > 0 {
		status, err = models.GetMsCurrencyIn(&msCurrency, currencyIds, "currency_key")
		if err != nil {
			// log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
	}
	currencyData := make(map[uint64]models.MsCurrency)
	for _, b := range msCurrency {
		currencyData[b.CurrencyKey] = b
	}

	var msRiskProfileData []models.MsRiskProfile
	if len(riskProfileIds) > 0 {
		status, err = models.GetMsRiskProfileIn(&msRiskProfileData, riskProfileIds)
		if err != nil {
			// log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
	}
	RiskProfileData := make(map[uint64]models.MsRiskProfile)
	for _, rp := range msRiskProfileData {
		RiskProfileData[rp.RiskProfileKey] = rp
	}

	//mapping MsProductCategory
	var msProductCategory []models.MsProductCategory
	if len(productCategoryIds) > 0 {
		status, err = models.GetMsProductCategoryIn(&msProductCategory, productCategoryIds, "product_category_key")
		if err != nil {
			// log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
	}
	proCatData := make(map[uint64]models.MsProductCategory)
	for _, a := range msProductCategory {
		proCatData[a.ProductCategoryKey] = a
	}

	//mapping product_type
	var msProductType []models.MsProductType
	if len(productTypeIds) > 0 {
		status, err = models.GetMsProductTypeIn(&msProductType, productTypeIds, "product_type_key")
		if err != nil {
			// log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
	}
	productTypeData := make(map[uint64]models.MsProductType)
	for _, p := range msProductType {
		productTypeData[p.ProductTypeKey] = p
	}

	//gen lookup
	var lookup []models.GenLookup
	if len(genLookupIds) > 0 {
		status, err = models.GetGenLookupIn(&lookup, genLookupIds, "lookup_key")
		if err != nil {
			if err != sql.ErrNoRows {
				// log.Error(err.Error())
				return lib.CustomError(status, err.Error(), "Failed get data")
			}
		}
	}

	gData := make(map[uint64]models.GenLookup)
	for _, gen := range lookup {
		gData[gen.LookupKey] = gen
	}

	//mapping parent custodian
	var msCustoBank []models.MsCustodianBank
	if len(custodianIds) > 0 {
		status, err = models.GetMsCustodianBankIn(&msCustoBank, custodianIds, "custodian_key")
		if err != nil {
			// log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
	}
	custoData := make(map[uint64]models.MsCustodianBank)
	for _, cus := range msCustoBank {
		custoData[cus.CustodianKey] = cus
	}

	var responseData []models.AdminMsProductList
	for _, pro := range msProduct {
		var data models.AdminMsProductList

		data.ProductKey = pro.ProductKey
		data.ProductCode = pro.ProductCode
		data.ProductName = pro.ProductName
		data.ProductNameAlt = pro.ProductNameAlt
		if pro.CurrencyKey != nil {
			if n, ok := currencyData[*pro.CurrencyKey]; ok {
				data.CurrencyName = n.Name
			}
		}
		if pro.ProductCategoryKey != nil {
			if n, ok := proCatData[*pro.ProductCategoryKey]; ok {
				data.ProductCategoryName = n.CategoryName
			}
		}
		if pro.ProductTypeKey != nil {
			if n, ok := productTypeData[*pro.ProductTypeKey]; ok {
				data.ProductTypeName = n.ProductTypeName
			}
		}
		if pro.RiskProfileKey != nil {
			if n, ok := RiskProfileData[*pro.RiskProfileKey]; ok {
				// log.Println(n.RiskCode, n.RiskName)
				data.RiskProfileName = n.RiskName
			}
		}
		layout := "2006-01-02 15:04:05"
		newLayout := "02 Jan 2006"
		if pro.LaunchDate != nil {
			date, err := time.Parse(layout, *pro.LaunchDate)
			if err == nil {
				oke := date.Format(newLayout)
				data.LaunchDate = &oke
			}
		}
		if pro.InceptionDate != nil {
			date, _ := time.Parse(layout, *pro.InceptionDate)
			if err == nil {
				oke := date.Format(newLayout)
				data.InceptionDate = &oke
			}
		}

		data.IsinCode = pro.IsinCode

		if pro.FlagSyariah == 1 {
			data.Syariah = true
		} else {
			data.Syariah = false
		}

		if pro.CustodianKey != nil {
			if n, ok := custoData[*pro.CustodianKey]; ok {
				data.CustodianFullName = n.CustodianFullName
			}
		}

		data.SinvestFundCode = pro.SinvestFundCode

		if pro.FlagEnabled == 1 {
			data.Enabled = true
		} else {
			data.Enabled = false
		}

		if pro.FlagSubscription == 1 {
			data.Subscription = true
		} else {
			data.Subscription = false
		}

		if pro.FlagRedemption == 1 {
			data.Redemption = true
		} else {
			data.Redemption = false
		}

		if pro.FlagSwitchOut == 1 {
			data.SwitchOut = true
		} else {
			data.SwitchOut = false
		}

		if pro.FlagSwitchIn == 1 {
			data.SwitchIn = true
		} else {
			data.SwitchIn = false
		}

		data.StatusUpdate = models.ProductStatusUpdate(strconv.FormatUint(data.ProductKey, 10))

		responseData = append(responseData, data)
	}

	var countData models.CountData
	var pagination int
	if limit > 0 {
		status, err = models.AdminGetCountMsProductWithLike(&countData, params, paramsLike)
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

func GetProductDetailAdmin(c echo.Context) error {
	PAGE_MENU_KEY := 71
	var err error
	var status int
	decimal.MarshalJSONWithoutQuotes = true

	errorAuth := IsMenuAccessAllowed(PAGE_MENU_KEY)
	if errorAuth != nil {
		// log.Error("User Autorizer")
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}

	keyStr := c.Param("key")
	key, _ := strconv.ParseUint(keyStr, 10, 64)
	if key == 0 {
		return lib.CustomError(http.StatusNotFound)
	}

	var product models.MsProduct
	status, err = models.GetMsProduct(&product, keyStr)
	if err != nil {
		return lib.CustomError(http.StatusNotFound)
	}

	var responseData models.AdminMsProductDetail

	var lookupIds []string

	if product.ProductPhase != nil {
		if _, ok := lib.Find(lookupIds, strconv.FormatUint(*product.ProductPhase, 10)); !ok {
			lookupIds = append(lookupIds, strconv.FormatUint(*product.ProductPhase, 10))
		}
	}
	if product.NavValuationType != nil {
		if _, ok := lib.Find(lookupIds, strconv.FormatUint(*product.NavValuationType, 10)); !ok {
			lookupIds = append(lookupIds, strconv.FormatUint(*product.NavValuationType, 10))
		}
	}

	//gen lookup oa request
	var lookupProduct []models.GenLookup
	if len(lookupIds) > 0 {
		status, err = models.GetGenLookupIn(&lookupProduct, lookupIds, "lookup_key")
		if err != nil {
			if err != sql.ErrNoRows {
				// log.Error(err.Error())
				return lib.CustomError(status, err.Error(), "Failed get data")
			}
		}
	}

	gData := make(map[uint64]models.GenLookup)
	for _, gen := range lookupProduct {
		gData[gen.LookupKey] = gen
	}

	responseData.ProductKey = product.ProductKey
	responseData.ProductCode = product.ProductCode
	responseData.ProductName = product.ProductName
	responseData.ProductNameAlt = product.ProductNameAlt
	if product.CurrencyKey != nil {
		var currency models.MsCurrency
		strCurrency := strconv.FormatUint(*product.CurrencyKey, 10)
		status, err = models.GetMsCurrency(&currency, strCurrency)
		if err != nil {
			if err != sql.ErrNoRows {
				return lib.CustomError(status)
			}
		} else {
			var cr models.MsCurrencyInfo
			cr.CurrencyKey = currency.CurrencyKey
			cr.Code = currency.Code
			cr.Symbol = currency.Symbol
			cr.Name = currency.Name
			cr.FlagBase = currency.FlagBase
			responseData.Currency = &cr
		}
	}

	if product.ProductCategoryKey != nil {
		var msProductCategory models.MsProductCategory
		strProCatKey := strconv.FormatUint(*product.ProductCategoryKey, 10)
		status, err = models.GetMsProductCategory(&msProductCategory, strProCatKey)
		if err != nil {
			if err != sql.ErrNoRows {
				return lib.CustomError(status)
			}
		} else {
			var cr models.MsProductCategoryInfo
			cr.ProductCategoryKey = msProductCategory.ProductCategoryKey
			cr.CategoryCode = msProductCategory.CategoryCode
			cr.CategoryName = msProductCategory.CategoryName
			cr.CategoryDesc = msProductCategory.CategoryDesc
			responseData.ProductCategory = &cr
		}
	}

	if product.ProductTypeKey != nil {
		var msProductType models.MsProductType
		strProTypeKey := strconv.FormatUint(*product.ProductTypeKey, 10)
		status, err = models.GetMsProductType(&msProductType, strProTypeKey)
		if err != nil {
			if err != sql.ErrNoRows {
				return lib.CustomError(status)
			}
		} else {
			var cr models.MsProductTypeInfo
			cr.ProductTypeKey = msProductType.ProductTypeKey
			cr.ProductTypeCode = msProductType.ProductTypeCode
			cr.ProductTypeName = msProductType.ProductTypeName
			cr.ProductTypeDesc = msProductType.ProductTypeDesc
			responseData.ProductType = &cr
		}
	}

	if product.FundTypeKey != nil {
		var fundType models.MsFundType
		strFundTypeKey := strconv.FormatUint(*product.FundTypeKey, 10)
		status, err = models.GetMsFundType(&fundType, strFundTypeKey)
		if err != nil {
			if err != sql.ErrNoRows {
				return lib.CustomError(status)
			}
		} else {
			var cr models.MsFundTypeInfo
			cr.FundTypeKey = fundType.FundTypeKey
			cr.FundTypeCode = fundType.FundTypeCode
			cr.FundTypeName = fundType.FundTypeName
			responseData.FundType = &cr
		}
	}

	if product.FundStructureKey != nil {
		var msFundStructure models.MsFundStructure
		strKeyFk := strconv.FormatUint(*product.FundStructureKey, 10)
		status, err = models.GetMsFundStructure(&msFundStructure, strKeyFk)
		if err != nil {
			if err != sql.ErrNoRows {
				return lib.CustomError(status)
			}
		} else {
			var cr models.MsFundStructureInfo
			cr.FundStructureKey = msFundStructure.FundStructureKey
			cr.FundStructureCode = msFundStructure.FundStructureCode
			cr.FundStructureName = msFundStructure.FundStructureName
			cr.FundStructureDesc = msFundStructure.FundStructureDesc
			responseData.FundStructure = &cr
		}
	}

	if product.RiskProfileKey != nil {
		var riskProfile models.MsRiskProfile
		strRiskProfileKey := strconv.FormatUint(*product.RiskProfileKey, 10)
		status, err = models.GetMsRiskProfile(&riskProfile, strRiskProfileKey)
		if err != nil {
			if err != sql.ErrNoRows {
				return lib.CustomError(status)
			}
		} else {
			var rs models.MsRiskProfileInfoAdmin
			rs.RiskProfileKey = riskProfile.RiskProfileKey
			rs.RiskCode = riskProfile.RiskCode
			rs.RiskName = riskProfile.RiskName
			rs.RiskDesc = riskProfile.RiskDesc
			responseData.RiskProfile = &rs
		}
	}

	responseData.ProductProfile = product.ProductProfile
	responseData.InvestmentObjectives = product.InvestmentObjectives

	if product.ProductPhase != nil {
		if n, ok := gData[*product.ProductPhase]; ok {
			var trc models.LookupTrans

			trc.LookupKey = n.LookupKey
			trc.LkpGroupKey = n.LkpGroupKey
			trc.LkpCode = n.LkpCode
			trc.LkpName = n.LkpName
			responseData.ProductPhase = &trc
		}
	}

	if product.NavValuationType != nil {
		if n, ok := gData[*product.NavValuationType]; ok {
			var trc models.LookupTrans

			trc.LookupKey = n.LookupKey
			trc.LkpGroupKey = n.LkpGroupKey
			trc.LkpCode = n.LkpCode
			trc.LkpName = n.LkpName
			responseData.NavValuationType = &trc
		}
	}

	responseData.ProspectusLink = product.ProspectusLink

	layout := "2006-01-02 15:04:05"
	newLayout := "02 Jan 2006"
	if product.LaunchDate != nil {
		date, err := time.Parse(layout, *product.LaunchDate)
		if err == nil {
			oke := date.Format(newLayout)
			responseData.LaunchDate = &oke
		}
	}
	if product.InceptionDate != nil {
		date, _ := time.Parse(layout, *product.InceptionDate)
		if err == nil {
			oke := date.Format(newLayout)
			responseData.InceptionDate = &oke
		}
	}

	responseData.IsinCode = product.IsinCode

	if product.FlagSyariah == 1 {
		responseData.FlagSyariah = true
	} else {
		responseData.FlagSyariah = false
	}

	responseData.MaxSubFee = product.MaxSubFee
	responseData.MaxRedFee = product.MaxRedFee
	responseData.MaxSwiFee = product.MaxSwiFee
	responseData.MinSubAmount = product.MinSubAmount
	responseData.MinRedAmount = product.MinRedAmount
	responseData.MinRedUnit = product.MinRedUnit
	responseData.MinUnitAfterRed = product.MinUnitAfterRed
	responseData.ManagementFee = product.ManagementFee
	responseData.CustodianFee = product.CustodianFee

	if product.CustodianKey != nil {
		var msCustodianBank models.MsCustodianBank
		strKeyFk := strconv.FormatUint(*product.CustodianKey, 10)
		status, err = models.GetMsCustodianBank(&msCustodianBank, strKeyFk)
		if err != nil {
			if err != sql.ErrNoRows {
				return lib.CustomError(status)
			}
		} else {
			var cr models.MsCustodianBankInfoList
			cr.CustodianKey = msCustodianBank.CustodianKey
			cr.CustodianCode = msCustodianBank.CustodianCode
			cr.CustodianShortName = msCustodianBank.CustodianShortName
			cr.CustodianFullName = msCustodianBank.CustodianFullName
			responseData.Custodian = &cr
		}
	}

	responseData.OjkFee = product.OjkFee
	responseData.ProductFeeAmount = product.ProductFeeAmount

	if product.OverwriteTransactFlag == 1 {
		responseData.OverwriteTransactFlag = true
	} else {
		responseData.OverwriteTransactFlag = false
	}

	if product.OverwriteFeeFlag == 1 {
		responseData.OverwriteFeeFlag = true
	} else {
		responseData.OverwriteFeeFlag = false
	}
	responseData.OtherFeeAmount = product.OtherFeeAmount
	responseData.SettlementPeriod = product.SettlementPeriod
	responseData.SinvestFundCode = product.SinvestFundCode

	if product.FlagEnabled == 1 {
		responseData.FlagEnabled = true
	} else {
		responseData.FlagEnabled = false
	}

	if product.FlagSubscription == 1 {
		responseData.FlagSubscription = true
	} else {
		responseData.FlagSubscription = false
	}

	if product.FlagRedemption == 1 {
		responseData.FlagRedemption = true
	} else {
		responseData.FlagRedemption = false
	}

	if product.FlagSwitchOut == 1 {
		responseData.FlagSwitchOut = true
	} else {
		responseData.FlagSwitchOut = false
	}

	if product.FlagSwitchIn == 1 {
		responseData.FlagSwitchIn = true
	} else {
		responseData.FlagSwitchIn = false
	}

	responseData.MinTopupAmount = product.MinTopUpAmount
	responseData.MinAmountAfterRed = product.MinAmountAfterRed
	responseData.DecAmount = uint64(product.DecAmount)
	responseData.DecNav = uint64(product.DecNav)
	responseData.DecPerformance = uint64(product.DecPerformance)
	responseData.DecUnit = uint64(product.DecUnit)
	responseData.NpwpDateReg = product.NpwpDateReg
	responseData.NpwpName = product.NpwpName
	responseData.NpwpNumber = product.NpwpNumber
	responseData.PortfolioCode = product.PortfolioCode

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = responseData

	return c.JSON(http.StatusOK, response)
}

func DeleteProductAdmin(c echo.Context) error {
	PAGE_MENU_KEY := 71
	var err error
	decimal.MarshalJSONWithoutQuotes = true

	errorAuth := IsMenuAccessAllowed(PAGE_MENU_KEY)
	if errorAuth != nil {
		// log.Error("User Autorizer")
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}

	params := make(map[string]string)

	productKey := c.FormValue("product_key")
	if productKey == "" {
		// log.Error("Missing required parameter: product_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: product_key", "Missing required parameter: product_key")
	}

	productKeyCek, err := strconv.ParseUint(productKey, 10, 64)
	if err == nil && productKeyCek > 0 {
		params["product_key"] = productKey
	} else {
		// log.Error("Wrong input for parameter: product_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: product_key", "Missing required parameter: product_key")
	}

	var product models.MsProduct
	status, err := models.GetMsProduct(&product, productKey)
	if err != nil {
		// log.Error("Product not found")
		return lib.CustomError(status)
	}

	dateLayout := "2006-01-02 15:04:05"
	params["rec_status"] = "0"
	params["rec_deleted_date"] = time.Now().Format(dateLayout)
	params["rec_deleted_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

	status, err = models.UpdateMsProduct(params)
	if err != nil {
		// log.Error("Failed create request data: " + err.Error())
		return lib.CustomError(status, err.Error(), "failed input data")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = nil
	return c.JSON(http.StatusOK, response)

}

func CreateAdminMsProduct(c echo.Context) error {
	PAGE_MENU_KEY := 71
	var err error
	var status int

	errorAuth := IsMenuAccessAllowed(PAGE_MENU_KEY)
	if errorAuth != nil {
		// log.Error("User Autorizer")
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}

	params := make(map[string]string)
	paramsCheckValidateAnd := make(map[string]string)
	paramsCheckValidateAnd["rec_status"] = "1"

	//product_code
	productcode := c.FormValue("product_code")
	if productcode == "" {
		// log.Error("Missing required parameter: product_code cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: product_code cann't be blank", "Missing required parameter: product_code cann't be blank")
	}
	params["product_code"] = productcode

	//check unique product_code
	paramsProdukCode := make(map[string]string)
	paramsProdukCode["product_code"] = productcode

	var countDataExisting models.CountData
	status, err = models.AdminGetValidateUniqueDataInsertUpdate(&countDataExisting, paramsProdukCode, paramsCheckValidateAnd, nil)
	if err != nil {
		// log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if int(countDataExisting.CountData) > 0 {
		// log.Error("Missing required parameter: product_code already existing, use other product_code")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: product_code already existing, use other product_code", "Missing required parameter: product_code already existing, use other product_code")
	}

	//product_id
	params["product_id"] = "0"

	//product_name
	productname := c.FormValue("product_name")
	if productname == "" {
		// log.Error("Missing required parameter: product_name cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: product_name cann't be blank", "Missing required parameter: product_name cann't be blank")
	}
	params["product_name"] = productname

	//check unique product_name
	paramsProductName := make(map[string]string)
	paramsProductName["product_name"] = productname

	status, err = models.AdminGetValidateUniqueDataInsertUpdate(&countDataExisting, paramsProductName, paramsCheckValidateAnd, nil)
	if err != nil {
		// log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if int(countDataExisting.CountData) > 0 {
		// log.Error("Missing required parameter: product_name already existing, use other product_name")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: product_name already existing, use other product_name", "Missing required parameter: product_name already existing, use other product_name")
	}

	//product_name_alt
	productnamealt := c.FormValue("product_name_alt")
	if productnamealt == "" {
		// log.Error("Missing required parameter: product_name_alt cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: product_name_alt cann't be blank", "Missing required parameter: product_name_alt cann't be blank")
	}
	params["product_name_alt"] = productnamealt

	//currency_key
	currencykey := c.FormValue("currency_key")
	if currencykey != "" {
		sub, err := strconv.ParseUint(currencykey, 10, 64)
		if err == nil && sub > 0 {
			params["currency_key"] = currencykey
		} else {
			// log.Error("Wrong input for parameter: currency_key number")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: currency_key must number", "Missing required parameter: currency_key number")
		}
	}

	//product_category_key
	productcategorykey := c.FormValue("product_category_key")
	if productcategorykey != "" {
		sub, err := strconv.ParseUint(productcategorykey, 10, 64)
		if err == nil && sub > 0 {
			params["product_category_key"] = productcategorykey
		} else {
			// log.Error("Wrong input for parameter: product_category_key number")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: product_category_key must number", "Missing required parameter: product_category_key number")
		}
	}

	//product_type_key
	// producttypekey := c.FormValue("product_type_key")
	// if producttypekey != "" {
	// 	sub, err := strconv.ParseUint(producttypekey, 10, 64)
	// 	if err == nil && sub > 0 {
	// 		params["product_type_key"] = producttypekey
	// 	} else {
	// 		// log.Error("Wrong input for parameter: product_type_key number")
	// 		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: product_type_key must number", "Missing required parameter: product_type_key number")
	// 	}
	// }

	//fund_type_key
	fundtypekey := c.FormValue("fund_type_key")
	if fundtypekey != "" {
		sub, err := strconv.ParseUint(fundtypekey, 10, 64)
		if err == nil && sub > 0 {
			params["fund_type_key"] = fundtypekey
		} else {
			// log.Error("Wrong input for parameter: fund_type_key number")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: fund_type_key must number", "Missing required parameter: fund_type_key number")
		}
	}

	//fund_structure_key
	// fundstructurekey := c.FormValue("fund_structure_key")
	// if fundstructurekey != "" {
	// 	sub, err := strconv.ParseUint(fundstructurekey, 10, 64)
	// 	if err == nil && sub > 0 {
	// 		params["fund_structure_key"] = fundstructurekey
	// 	} else {
	// 		// log.Error("Wrong input for parameter: fund_structure_key number")
	// 		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: fund_structure_key must number", "Missing required parameter: fund_structure_key number")
	// 	}
	// }

	//risk_profile_key
	// riskprofilekey := c.FormValue("risk_profile_key")
	// if riskprofilekey != "" {
	// 	sub, err := strconv.ParseUint(riskprofilekey, 10, 64)
	// 	if err == nil && sub > 0 {
	// 		params["risk_profile_key"] = riskprofilekey
	// 	} else {
	// 		// log.Error("Wrong input for parameter: risk_profile_key number")
	// 		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: risk_profile_key must number", "Missing required parameter: risk_profile_key number")
	// 	}
	// }
	riskProfileKey := models.GetRiskProfileByFundType(fundtypekey)
	params["risk_profile_key"] = strconv.FormatUint(riskProfileKey, 10)

	custodiankey := c.FormValue("custodian_key")
	if custodiankey != "" {
		sub, err := strconv.ParseUint(custodiankey, 10, 64)
		if err == nil && sub > 0 {
			params["custodian_key"] = custodiankey
		} else {
			// log.Error("Wrong input for parameter: custodian_key number")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: custodian_key must number", "Missing required parameter: custodian_key number")
		}
	}

	//product_profile
	productprofile := c.FormValue("product_profile")
	if productprofile != "" {
		params["product_profile"] = productprofile
	}

	//investment_objectives
	investmentobjectives := c.FormValue("investment_objectives")
	if investmentobjectives != "" {
		params["investment_objectives"] = investmentobjectives
	}

	//product_phase
	productphase := c.FormValue("product_phase")
	if productphase != "" {
		sub, err := strconv.ParseUint(productphase, 10, 64)
		if err == nil && sub > 0 {
			params["product_phase"] = productphase
		} else {
			// log.Error("Wrong input for parameter: product_phase number")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: product_phase must number", "Missing required parameter: product_phase number")
		}
	}

	//nav_valuation_type
	navvaluationtype := c.FormValue("nav_valuation_type")
	if navvaluationtype != "" {
		sub, err := strconv.ParseUint(navvaluationtype, 10, 64)
		if err == nil && sub > 0 {
			params["nav_valuation_type"] = navvaluationtype
		} else {
			// log.Error("Wrong input for parameter: nav_valuation_type number")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: nav_valuation_type must number", "Missing required parameter: nav_valuation_type number")
		}
	}

	//prospectus_link
	prospectuslink := c.FormValue("prospectus_link")
	if prospectuslink != "" {
		length := len(prospectuslink)
		if length > 255 {
			// log.Error("Wrong input for parameter: prospectus_link number")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: prospectus_link too long, max 255 character", "Missing required parameter: prospectus_link too long, max 255 character")
		}
		params["prospectus_link"] = prospectuslink
	}

	//launch_date
	launchdate := c.FormValue("launch_date")
	if launchdate != "" {
		params["launch_date"] = launchdate + " 00:00:00"
	}

	//inception_date
	inceptiondate := c.FormValue("inception_date")
	if inceptiondate != "" {
		params["inception_date"] = inceptiondate + " 00:00:00"
	}

	//isin_code
	isincode := c.FormValue("isin_code")
	if isincode != "" {
		params["isin_code"] = isincode

		//check unique isin_code
		paramsIsinCode := make(map[string]string)
		paramsIsinCode["isin_code"] = isincode

		status, err = models.AdminGetValidateUniqueDataInsertUpdate(&countDataExisting, paramsIsinCode, paramsCheckValidateAnd, nil)
		if err != nil {
			// log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
		if int(countDataExisting.CountData) > 0 {
			// log.Error("Missing required parameter: isin_code already existing, use other isin_code")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: isin_code already existing, use other isin_code", "Missing required parameter: isin_code already existing, use other isin_code")
		}
	}

	//flag_syariah
	flagsyariah := c.FormValue("flag_syariah")
	var flagsyariahBool bool
	if flagsyariah != "" {
		flagsyariahBool, err = strconv.ParseBool(flagsyariah)
		if err != nil {
			// log.Error("flag_syariah parameter should be true/false")
			return lib.CustomError(http.StatusBadRequest, "flag_syariah parameter should be true/false", "flag_syariah parameter should be true/false")
		}
		if flagsyariahBool == true {
			params["flag_syariah"] = "1"
		} else {
			params["flag_syariah"] = "0"
		}
	} else {
		// log.Error("flag_syariah parameter should be true/false")
		return lib.CustomError(http.StatusBadRequest, "flag_syariah parameter should be true/false", "flag_syariah parameter should be true/false")
	}

	//max_sub_fee
	maxsubfee := c.FormValue("max_sub_fee")
	if maxsubfee == "" {
		// log.Error("Missing required parameter: max_sub_fee cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: max_sub_fee cann't be blank", "Missing required parameter: max_sub_fee cann't be blank")
	}
	maxsubfeeFloat, err := strconv.ParseFloat(maxsubfee, 64)
	if err == nil {
		if maxsubfeeFloat < 0 {
			// log.Error("Wrong input for parameter: max_sub_fee cann't negatif")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: max_sub_fee must cann't negatif", "Missing required parameter: max_sub_fee cann't negatif")
		}
		params["max_sub_fee"] = maxsubfee
	} else {
		// log.Error("Wrong input for parameter: max_sub_fee number")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: max_sub_fee must number", "Missing required parameter: max_sub_fee number")
	}

	//max_red_fee
	maxredfee := c.FormValue("max_red_fee")
	if maxredfee == "" {
		// log.Error("Missing required parameter: max_red_fee cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: max_red_fee cann't be blank", "Missing required parameter: max_red_fee cann't be blank")
	}
	maxredfeeFloat, err := strconv.ParseFloat(maxredfee, 64)
	if err == nil {
		if maxredfeeFloat < 0 {
			// log.Error("Wrong input for parameter: max_red_fee cann't negatif")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: max_red_fee must cann't negatif", "Missing required parameter: max_red_fee cann't negatif")
		}
		params["max_red_fee"] = maxredfee
	} else {
		// log.Error("Wrong input for parameter: max_red_fee number")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: max_red_fee must number", "Missing required parameter: max_red_fee number")
	}

	//max_swi_fee
	maxswifee := c.FormValue("max_swi_fee")
	if maxswifee == "" {
		// log.Error("Missing required parameter: max_swi_fee cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: max_swi_fee cann't be blank", "Missing required parameter: max_swi_fee cann't be blank")
	}
	maxswifeeFloat, err := strconv.ParseFloat(maxswifee, 64)
	if err == nil {
		if maxswifeeFloat < 0 {
			// log.Error("Wrong input for parameter: max_swi_fee cann't negatif")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: max_swi_fee must cann't negatif", "Missing required parameter: max_swi_fee cann't negatif")
		}
		params["max_swi_fee"] = maxswifee
	} else {
		// log.Error("Wrong input for parameter: max_swi_fee number")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: max_swi_fee must number", "Missing required parameter: max_swi_fee number")
	}

	//min_sub_amount
	minsubamount := c.FormValue("min_sub_amount")
	if minsubamount == "" {
		// log.Error("Missing required parameter: min_sub_amount cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: min_sub_amount cann't be blank", "Missing required parameter: min_sub_amount cann't be blank")
	}
	minsubamountFloat, err := strconv.ParseFloat(minsubamount, 64)
	if err == nil {
		if minsubamountFloat < 0 {
			// log.Error("Wrong input for parameter: min_sub_amount cann't negatif")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: min_sub_amount must cann't negatif", "Missing required parameter: min_sub_amount cann't negatif")
		}
		params["min_sub_amount"] = minsubamount
	} else {
		// log.Error("Wrong input for parameter: min_sub_amount number")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: min_sub_amount must number", "Missing required parameter: min_sub_amount number")
	}

	mintopupamount := c.FormValue("min_topup_amount")
	if mintopupamount == "" {
		// log.Error("Missing required parameter: min_topup_amount cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: min_topup_amount cann't be blank", "Missing required parameter: min_topup_amount cann't be blank")
	}
	mintopupamountFloat, err := strconv.ParseFloat(mintopupamount, 64)
	if err == nil {
		if mintopupamountFloat < 0 {
			// log.Error("Wrong input for parameter: min_topup_amount cann't negatif")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: min_topup_amount must cann't negatif", "Missing required parameter: min_topup_amount cann't negatif")
		}
		params["min_topup_amount"] = mintopupamount
	} else {
		// log.Error("Wrong input for parameter: min_topup_amount number")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: min_topup_amount must number", "Missing required parameter: min_topup_amount number")
	}

	//min_red_amount
	minredamount := c.FormValue("min_red_amount")
	if minredamount == "" {
		// log.Error("Missing required parameter: min_red_amount cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: min_red_amount cann't be blank", "Missing required parameter: min_red_amount cann't be blank")
	}
	minredamountFloat, err := strconv.ParseFloat(minredamount, 64)
	if err == nil {
		if minredamountFloat < 0 {
			// log.Error("Wrong input for parameter: min_red_amount cann't negatif")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: min_red_amount must cann't negatif", "Missing required parameter: min_red_amount cann't negatif")
		}
		params["min_red_amount"] = minredamount
	} else {
		// log.Error("Wrong input for parameter: min_red_amount number")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: min_red_amount must number", "Missing required parameter: min_red_amount number")
	}

	//min_red_unit
	minredunit := c.FormValue("min_red_unit")
	if minredunit == "" {
		// log.Error("Missing required parameter: min_red_unit cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: min_red_unit cann't be blank", "Missing required parameter: min_red_unit cann't be blank")
	}
	minredunitFloat, err := strconv.ParseFloat(minredunit, 64)
	if err == nil {
		if minredunitFloat < 0 {
			// log.Error("Wrong input for parameter: min_red_unit cann't negatif")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: min_red_unit must cann't negatif", "Missing required parameter: min_red_unit cann't negatif")
		}
		params["min_red_unit"] = minredunit
	} else {
		// log.Error("Wrong input for parameter: min_red_unit number")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: min_red_unit must number", "Missing required parameter: min_red_unit number")
	}

	//min_unit_after_red
	minunitafterred := c.FormValue("min_unit_after_red")
	if minunitafterred == "" {
		// log.Error("Missing required parameter: min_unit_after_red cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: min_unit_after_red cann't be blank", "Missing required parameter: min_unit_after_red cann't be blank")
	}
	minunitafterredFloat, err := strconv.ParseFloat(minunitafterred, 64)
	if err == nil {
		if minunitafterredFloat < 0 {
			// log.Error("Wrong input for parameter: min_unit_after_red cann't negatif")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: min_unit_after_red must cann't negatif", "Missing required parameter: min_unit_after_red cann't negatif")
		}
		params["min_unit_after_red"] = minunitafterred
	} else {
		// log.Error("Wrong input for parameter: min_unit_after_red number")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: min_unit_after_red must number", "Missing required parameter: min_unit_after_red number")
	}

	minamountafterred := c.FormValue("min_amount_after_red")
	if minamountafterred == "" {
		// log.Error("Missing required parameter: min_amount_after_red cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: min_amount_after_red cann't be blank", "Missing required parameter: min_amount_after_red cann't be blank")
	}
	minamountafterredFloat, err := strconv.ParseFloat(minamountafterred, 64)
	if err == nil {
		if minamountafterredFloat < 0 {
			// log.Error("Wrong input for parameter: min_amount_after_red cann't negatif")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: min_amount_after_red must cann't negatif", "Missing required parameter: min_amount_after_red cann't negatif")
		}
		params["min_amount_after_red"] = minamountafterred
	} else {
		// log.Error("Wrong input for parameter: min_amount_after_red number")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: min_amount_after_red must number", "Missing required parameter: min_amount_after_red number")
	}

	//management_fee
	managementfee := c.FormValue("management_fee")
	if managementfee == "" {
		// log.Error("Missing required parameter: management_fee cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: management_fee cann't be blank", "Missing required parameter: management_fee cann't be blank")
	}
	managementfeeFloat, err := strconv.ParseFloat(managementfee, 64)
	if err == nil {
		if managementfeeFloat < 0 {
			// log.Error("Wrong input for parameter: management_fee cann't negatif")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: management_fee must cann't negatif", "Missing required parameter: management_fee cann't negatif")
		}
		params["management_fee"] = managementfee
	} else {
		// log.Error("Wrong input for parameter: management_fee number")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: management_fee must number", "Missing required parameter: management_fee number")
	}

	//custodian_fee
	custodianfee := c.FormValue("custodian_fee")
	if custodianfee == "" {
		// log.Error("Missing required parameter: custodian_fee cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: custodian_fee cann't be blank", "Missing required parameter: custodian_fee cann't be blank")
	}
	custodianfeeFloat, err := strconv.ParseFloat(custodianfee, 64)
	if err == nil {
		if custodianfeeFloat < 0 {
			// log.Error("Wrong input for parameter: custodian_fee cann't negatif")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: custodian_fee must cann't negatif", "Missing required parameter: custodian_fee cann't negatif")
		}
		params["custodian_fee"] = custodianfee
	} else {
		// log.Error("Wrong input for parameter: custodian_fee number")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: custodian_fee must number", "Missing required parameter: custodian_fee number")
	}

	//ojk_fee
	// ojkfee := c.FormValue("ojk_fee")
	// if ojkfee == "" {
	// 	// log.Error("Missing required parameter: ojk_fee cann't be blank")
	// 	return lib.CustomError(http.StatusBadRequest, "Missing required parameter: ojk_fee cann't be blank", "Missing required parameter: ojk_fee cann't be blank")
	// }
	// ojkfeeFloat, err := strconv.ParseFloat(ojkfee, 64)
	// if err == nil {
	// 	if ojkfeeFloat < 0 {
	// 		// log.Error("Wrong input for parameter: ojk_fee cann't negatif")
	// 		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: ojk_fee must cann't negatif", "Missing required parameter: ojk_fee cann't negatif")
	// 	}
	// 	params["ojk_fee"] = ojkfee
	// } else {
	// 	// log.Error("Wrong input for parameter: ojk_fee number")
	// 	return lib.CustomError(http.StatusBadRequest, "Missing required parameter: ojk_fee must number", "Missing required parameter: ojk_fee number")
	// }

	//product_fee_amount
	// productfeeamount := c.FormValue("product_fee_amount")
	// if productfeeamount == "" {
	// 	// log.Error("Missing required parameter: product_fee_amount cann't be blank")
	// 	return lib.CustomError(http.StatusBadRequest, "Missing required parameter: product_fee_amount cann't be blank", "Missing required parameter: product_fee_amount cann't be blank")
	// }
	// productfeeamountFloat, err := strconv.ParseFloat(productfeeamount, 64)
	// if err == nil {
	// 	if productfeeamountFloat < 0 {
	// 		// log.Error("Wrong input for parameter: product_fee_amount cann't negatif")
	// 		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: product_fee_amount must cann't negatif", "Missing required parameter: product_fee_amount cann't negatif")
	// 	}
	// 	params["product_fee_amount"] = productfeeamount
	// } else {
	// 	// log.Error("Wrong input for parameter: product_fee_amount number")
	// 	return lib.CustomError(http.StatusBadRequest, "Missing required parameter: product_fee_amount must number", "Missing required parameter: product_fee_amount number")
	// }

	//overwrite_transact_flag
	// overwritetransactflag := c.FormValue("overwrite_transact_flag")
	// var overwritetransactflagBool bool
	// if overwritetransactflag != "" {
	// 	overwritetransactflagBool, err = strconv.ParseBool(overwritetransactflag)
	// 	if err != nil {
	// 		// log.Error("overwrite_transact_flag parameter should be true/false")
	// 		return lib.CustomError(http.StatusBadRequest, "overwrite_transact_flag parameter should be true/false", "overwrite_transact_flag parameter should be true/false")
	// 	}
	// 	if overwritetransactflagBool == true {
	// 		params["overwrite_transact_flag"] = "1"
	// 	} else {
	// 		params["overwrite_transact_flag"] = "0"
	// 	}
	// } else {
	// 	// log.Error("overwrite_transact_flag parameter should be true/false")
	// 	return lib.CustomError(http.StatusBadRequest, "overwrite_transact_flag parameter should be true/false", "overwrite_transact_flag parameter should be true/false")
	// }

	//overwrite_fee_flag
	// overwritefeeflag := c.FormValue("overwrite_fee_flag")
	// var overwritefeeflagBool bool
	// if overwritefeeflag != "" {
	// 	overwritefeeflagBool, err = strconv.ParseBool(overwritefeeflag)
	// 	if err != nil {
	// 		// log.Error("overwrite_fee_flag parameter should be true/false")
	// 		return lib.CustomError(http.StatusBadRequest, "overwrite_fee_flag parameter should be true/false", "overwrite_fee_flag parameter should be true/false")
	// 	}
	// 	if overwritefeeflagBool == true {
	// 		params["overwrite_fee_flag"] = "1"
	// 	} else {
	// 		params["overwrite_fee_flag"] = "0"
	// 	}
	// } else {
	// 	// log.Error("overwrite_fee_flag parameter should be true/false")
	// 	return lib.CustomError(http.StatusBadRequest, "overwrite_fee_flag parameter should be true/false", "overwrite_fee_flag parameter should be true/false")
	// }

	//other_fee_amount
	// otherfeeamount := c.FormValue("other_fee_amount")
	// if otherfeeamount == "" {
	// 	// log.Error("Missing required parameter: other_fee_amount cann't be blank")
	// 	return lib.CustomError(http.StatusBadRequest, "Missing required parameter: other_fee_amount cann't be blank", "Missing required parameter: other_fee_amount cann't be blank")
	// }
	// otherfeeamountFloat, err := strconv.ParseFloat(otherfeeamount, 64)
	// if err == nil {
	// 	if otherfeeamountFloat < 0 {
	// 		// log.Error("Wrong input for parameter: other_fee_amount cann't negatif")
	// 		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: other_fee_amount must cann't negatif", "Missing required parameter: other_fee_amount cann't negatif")
	// 	}
	// 	params["other_fee_amount"] = otherfeeamount
	// } else {
	// 	// log.Error("Wrong input for parameter: other_fee_amount number")
	// 	return lib.CustomError(http.StatusBadRequest, "Missing required parameter: other_fee_amount must number", "Missing required parameter: other_fee_amount number")
	// }

	//settlement_period
	settlementperiod := c.FormValue("settlement_period")
	if settlementperiod != "" {
		_, err = strconv.ParseUint(settlementperiod, 10, 64)
		if err == nil {
			params["settlement_period"] = settlementperiod
		} else {
			// log.Error("Wrong input for parameter: settlement_period number")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: settlement_period must number", "Missing required parameter: settlement_period number")
		}
	}

	//sinvest_fund_code
	sinvestfundcode := c.FormValue("sinvest_fund_code")
	if sinvestfundcode != "" {
		params["sinvest_fund_code"] = sinvestfundcode

		//check unique sinvest_fund_code
		paramsSinvest := make(map[string]string)
		paramsSinvest["sinvest_fund_code"] = sinvestfundcode

		status, err = models.AdminGetValidateUniqueDataInsertUpdate(&countDataExisting, paramsSinvest, paramsCheckValidateAnd, nil)
		if err != nil {
			// log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
		if int(countDataExisting.CountData) > 0 {
			// log.Error("Missing required parameter: sinvest_fund_code already existing, use other sinvest_fund_code")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: sinvest_fund_code already existing, use other sinvest_fund_code", "Missing required parameter: sinvest_fund_code already existing, use other sinvest_fund_code")
		}
	}

	//flag_enabled
	// flagenabled := c.FormValue("flag_enabled")
	// var flagenabledBool bool
	// if flagenabled != "" {
	// 	flagenabledBool, err = strconv.ParseBool(flagenabled)
	// 	if err != nil {
	// 		// log.Error("flag_enabled parameter should be true/false")
	// 		return lib.CustomError(http.StatusBadRequest, "flag_enabled parameter should be true/false", "flag_enabled parameter should be true/false")
	// 	}
	// 	if flagenabledBool == true {
	// 		params["flag_enabled"] = "1"
	// 	} else {
	// 		params["flag_enabled"] = "0"
	// 	}
	// } else {
	// 	// log.Error("flag_enabled parameter should be true/false")
	// 	return lib.CustomError(http.StatusBadRequest, "flag_enabled parameter should be true/false", "flag_enabled parameter should be true/false")
	// }
	params["flag_enabled"] = "1"

	//flag_subscription
	// flagsubscription := c.FormValue("flag_subscription")
	// var flagsubscriptionBool bool
	// if flagsubscription != "" {
	// 	flagsubscriptionBool, err = strconv.ParseBool(flagsubscription)
	// 	if err != nil {
	// 		// log.Error("flag_subscription parameter should be true/false")
	// 		return lib.CustomError(http.StatusBadRequest, "flag_subscription parameter should be true/false", "flag_subscription parameter should be true/false")
	// 	}
	// 	if flagsubscriptionBool == true {
	// 		params["flag_subscription"] = "1"
	// 	} else {
	// 		params["flag_subscription"] = "0"
	// 	}
	// } else {
	// 	// log.Error("flag_subscription parameter should be true/false")
	// 	return lib.CustomError(http.StatusBadRequest, "flag_subscription parameter should be true/false", "flag_subscription parameter should be true/false")
	// }
	params["flag_subscription"] = "1"

	//flag_redemption
	// flagredemption := c.FormValue("flag_redemption")
	// var flagredemptionBool bool
	// if flagredemption != "" {
	// 	flagredemptionBool, err = strconv.ParseBool(flagredemption)
	// 	if err != nil {
	// 		// log.Error("flag_redemption parameter should be true/false")
	// 		return lib.CustomError(http.StatusBadRequest, "flag_redemption parameter should be true/false", "flag_redemption parameter should be true/false")
	// 	}
	// 	if flagredemptionBool == true {
	// 		params["flag_redemption"] = "1"
	// 	} else {
	// 		params["flag_redemption"] = "0"
	// 	}
	// } else {
	// 	// log.Error("flag_redemption parameter should be true/false")
	// 	return lib.CustomError(http.StatusBadRequest, "flag_redemption parameter should be true/false", "flag_redemption parameter should be true/false")
	// }
	params["flag_redemption"] = "1"

	//flag_switch_out
	// flagswitchout := c.FormValue("flag_switch_out")
	// var flagswitchoutBool bool
	// if flagswitchout != "" {
	// 	flagswitchoutBool, err = strconv.ParseBool(flagswitchout)
	// 	if err != nil {
	// 		// log.Error("flag_switch_out parameter should be true/false")
	// 		return lib.CustomError(http.StatusBadRequest, "flag_switch_out parameter should be true/false", "flag_switch_out parameter should be true/false")
	// 	}
	// 	if flagswitchoutBool == true {
	// 		params["flag_switch_out"] = "1"
	// 	} else {
	// 		params["flag_switch_out"] = "0"
	// 	}
	// } else {
	// 	// log.Error("flag_switch_out parameter should be true/false")
	// 	return lib.CustomError(http.StatusBadRequest, "flag_switch_out parameter should be true/false", "flag_switch_out parameter should be true/false")
	// }
	params["flag_switch_out"] = "1"

	//flag_switch_in
	// flagswitchin := c.FormValue("flag_switch_in")
	// var flagswitchinBool bool
	// if flagswitchin != "" {
	// 	flagswitchinBool, err = strconv.ParseBool(flagswitchin)
	// 	if err != nil {
	// 		// log.Error("flag_switch_in parameter should be true/false")
	// 		return lib.CustomError(http.StatusBadRequest, "flag_switch_in parameter should be true/false", "flag_switch_in parameter should be true/false")
	// 	}
	// 	if flagswitchinBool == true {
	// 		params["flag_switch_in"] = "1"
	// 	} else {
	// 		params["flag_switch_in"] = "0"
	// 	}
	// } else {
	// 	// log.Error("flag_switch_in parameter should be true/false")
	// 	return lib.CustomError(http.StatusBadRequest, "flag_switch_in parameter should be true/false", "flag_switch_in parameter should be true/false")
	// }
	params["flag_switch_in"] = "1"

	// decNav := c.FormValue("dec_nav")
	// if decNav == "" {
	// 	return lib.CustomError(http.StatusBadRequest, "Missing dec_nav", "Missing dec_nav")
	// } else {
	// 	params["dec_nav"] = decNav
	// }
	params["dec_nav"] = "4"

	// decunit := c.FormValue("dec_unit")
	// if decunit == "" {
	// 	return lib.CustomError(http.StatusBadRequest, "Missing dec_unit", "Missing dec_unit")
	// } else {
	// 	params["dec_unit"] = decNav
	// }
	params["dec_unit"] = "4"

	// decperformance := c.FormValue("dec_performance")
	// if decperformance == "" {
	// 	return lib.CustomError(http.StatusBadRequest, "Missing dec_performance", "Missing dec_performance")
	// } else {
	// 	params["dec_performance"] = decperformance
	// }
	params["dec_performance"] = "4"

	// decamount := c.FormValue("dec_amount")
	// if decamount == "" {
	// 	return lib.CustomError(http.StatusBadRequest, "Missing dec_amount", "Missing dec_amount")
	// } else {
	// 	params["dec_amount"] = decamount
	// }
	params["dec_amount"] = "4"

	// npwp_number := c.FormValue("npwp_number")
	// params["npwp_number"] = npwp_number

	// npwp_date_reg := c.FormValue("npwp_date_reg")
	// params["npwp_date_reg"] = npwp_date_reg

	// npwp_name := c.FormValue("npwp_name")
	// params["npwp_name"] = npwp_name

	// portfolio_code := c.FormValue("portfolio_code")
	// params["portfolio_code"] = portfolio_code

	dateLayout := "2006-01-02 15:04:05"
	params["rec_status"] = "1"
	params["rec_order"] = "0"
	params["rec_created_date"] = time.Now().Format(dateLayout)
	params["rec_created_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

	status, err = models.CreateMsProduct(params)
	if err != nil {
		// log.Error("Failed create request data: " + err.Error())
		return lib.CustomError(status, err.Error(), "failed input data")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = nil
	return c.JSON(http.StatusOK, response)

}

func UpdateAdminMsProduct(c echo.Context) error {
	PAGE_MENU_KEY := 71
	var err error
	var status int
	decimal.MarshalJSONWithoutQuotes = true

	errorAuth := IsMenuAccessAllowed(PAGE_MENU_KEY)
	if errorAuth != nil {
		// log.Error("User Autorizer")
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}

	params := make(map[string]string)
	paramsCheckValidateAnd := make(map[string]string)
	paramsCheckValidateAnd["rec_status"] = "1"

	productkey := c.FormValue("product_key")
	if productkey == "" {
		// log.Error("Missing required parameter: product_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: product_key", "Missing required parameter: product_key")
	}
	strproductkey, err := strconv.ParseUint(productkey, 10, 64)
	if err == nil && strproductkey > 0 {
		params["product_key"] = productkey
	} else {
		// log.Error("Wrong input for parameter: product_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: product_key", "Missing required parameter: product_key")
	}

	var product models.MsProduct
	status, err = models.GetMsProduct(&product, productkey)
	if err != nil {
		// log.Error(err.Error())
		return lib.CustomError(http.StatusBadRequest)
	}

	if product.RecStatus == 0 {
		// log.Error("Product not found, rec_status = 0")
		return lib.CustomError(http.StatusBadRequest)
	}

	//product_code
	productcode := c.FormValue("product_code")
	if productcode == "" {
		// log.Error("Missing required parameter: product_code cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: product_code cann't be blank", "Missing required parameter: product_code cann't be blank")
	}
	params["product_code"] = productcode

	//check unique product_code
	paramsProdukCode := make(map[string]string)
	paramsProdukCode["product_code"] = productcode

	var countDataExisting models.CountData
	status, err = models.AdminGetValidateUniqueDataInsertUpdate(&countDataExisting, paramsProdukCode, paramsCheckValidateAnd, &productkey)
	if err != nil {
		// log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if int(countDataExisting.CountData) > 0 {
		// log.Error("Missing required parameter: product_code already existing, use other product_code")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: product_code already existing, use other product_code", "Missing required parameter: product_code already existing, use other product_code")
	}

	var sub uint64

	//product_name
	productname := c.FormValue("product_name")
	if productname == "" {
		// log.Error("Missing required parameter: product_name cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: product_name cann't be blank", "Missing required parameter: product_name cann't be blank")
	}
	params["product_name"] = productname

	//check unique product_name
	paramsProductName := make(map[string]string)
	paramsProductName["product_name"] = productname

	status, err = models.AdminGetValidateUniqueDataInsertUpdate(&countDataExisting, paramsProductName, paramsCheckValidateAnd, &productkey)
	if err != nil {
		// log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if int(countDataExisting.CountData) > 0 {
		// log.Error("Missing required parameter: product_name already existing, use other product_name")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: product_name already existing, use other product_name", "Missing required parameter: product_name already existing, use other product_name")
	}

	//product_name_alt
	productnamealt := c.FormValue("product_name_alt")
	if productnamealt == "" {
		// log.Error("Missing required parameter: product_name_alt cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: product_name_alt cann't be blank", "Missing required parameter: product_name_alt cann't be blank")
	}
	params["product_name_alt"] = productnamealt

	//currency_key
	currencykey := c.FormValue("currency_key")
	if currencykey != "" {
		sub, err = strconv.ParseUint(currencykey, 10, 64)
		if err == nil && sub > 0 {
			params["currency_key"] = currencykey
		} else {
			// log.Error("Wrong input for parameter: currency_key number")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: currency_key must number", "Missing required parameter: currency_key number")
		}
	}

	//product_category_key
	productcategorykey := c.FormValue("product_category_key")
	if productcategorykey != "" {
		sub, err = strconv.ParseUint(productcategorykey, 10, 64)
		if err == nil && sub > 0 {
			params["product_category_key"] = productcategorykey
		} else {
			// log.Error("Wrong input for parameter: product_category_key number")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: product_category_key must number", "Missing required parameter: product_category_key number")
		}
	}

	//product_type_key
	producttypekey := c.FormValue("product_type_key")
	if producttypekey != "" {
		sub, err = strconv.ParseUint(producttypekey, 10, 64)
		if err == nil && sub > 0 {
			params["product_type_key"] = producttypekey
		} else {
			// log.Error("Wrong input for parameter: product_type_key number")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: product_type_key must number", "Missing required parameter: product_type_key number")
		}
	}

	//fund_type_key
	fundtypekey := c.FormValue("fund_type_key")
	if fundtypekey != "" {
		sub, err = strconv.ParseUint(fundtypekey, 10, 64)
		if err == nil && sub > 0 {
			params["fund_type_key"] = fundtypekey
		} else {
			// log.Error("Wrong input for parameter: fund_type_key number")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: fund_type_key must number", "Missing required parameter: fund_type_key number")
		}
	}

	//fund_structure_key
	fundstructurekey := c.FormValue("fund_structure_key")
	if fundstructurekey != "" {
		sub, err = strconv.ParseUint(fundstructurekey, 10, 64)
		if err == nil && sub > 0 {
			params["fund_structure_key"] = fundstructurekey
		} else {
			// log.Error("Wrong input for parameter: fund_structure_key number")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: fund_structure_key must number", "Missing required parameter: fund_structure_key number")
		}
	}

	//risk_profile_key
	riskprofilekey := c.FormValue("risk_profile_key")
	if riskprofilekey != "" {
		sub, err = strconv.ParseUint(riskprofilekey, 10, 64)
		if err == nil && sub > 0 {
			params["risk_profile_key"] = riskprofilekey
		} else {
			// log.Error("Wrong input for parameter: risk_profile_key number")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: risk_profile_key must number", "Missing required parameter: risk_profile_key number")
		}
	}

	//product_profile
	productprofile := c.FormValue("product_profile")
	if productprofile != "" {
		params["product_profile"] = productprofile
	}

	//investment_objectives
	investmentobjectives := c.FormValue("investment_objectives")
	if investmentobjectives != "" {
		params["investment_objectives"] = investmentobjectives
	}

	//product_phase
	productphase := c.FormValue("product_phase")
	if productphase != "" {
		sub, err = strconv.ParseUint(productphase, 10, 64)
		if err == nil && sub > 0 {
			params["product_phase"] = productphase
		} else {
			// log.Error("Wrong input for parameter: product_phase number")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: product_phase must number", "Missing required parameter: product_phase number")
		}
	}

	//nav_valuation_type
	navvaluationtype := c.FormValue("nav_valuation_type")
	if navvaluationtype != "" {
		sub, err = strconv.ParseUint(navvaluationtype, 10, 64)
		if err == nil && sub > 0 {
			params["nav_valuation_type"] = navvaluationtype
		} else {
			// log.Error("Wrong input for parameter: nav_valuation_type number")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: nav_valuation_type must number", "Missing required parameter: nav_valuation_type number")
		}
	}

	//prospectus_link
	prospectuslink := c.FormValue("prospectus_link")
	if prospectuslink != "" {
		length := len(prospectuslink)
		if length > 255 {
			// log.Error("Wrong input for parameter: prospectus_link number")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: prospectus_link too long, max 255 character", "Missing required parameter: prospectus_link too long, max 255 character")
		}
		params["prospectus_link"] = prospectuslink
	}

	//launch_date
	launchdate := c.FormValue("launch_date")
	if launchdate != "" {
		params["launch_date"] = launchdate + " 00:00:00"
	}

	//inception_date
	inceptiondate := c.FormValue("inception_date")
	if inceptiondate != "" {
		params["inception_date"] = inceptiondate + " 00:00:00"
	}

	//isin_code
	isincode := c.FormValue("isin_code")
	if isincode != "" {
		params["isin_code"] = isincode

		//check unique isin_code
		paramsIsinCode := make(map[string]string)
		paramsIsinCode["isin_code"] = isincode

		status, err = models.AdminGetValidateUniqueDataInsertUpdate(&countDataExisting, paramsIsinCode, paramsCheckValidateAnd, &productkey)
		if err != nil {
			// log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
		if int(countDataExisting.CountData) > 0 {
			// log.Error("Missing required parameter: isin_code already existing, use other isin_code")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: isin_code already existing, use other isin_code", "Missing required parameter: isin_code already existing, use other isin_code")
		}
	}

	//flag_syariah
	flagsyariah := c.FormValue("flag_syariah")
	var flagsyariahBool bool
	if flagsyariah != "" {
		flagsyariahBool, err = strconv.ParseBool(flagsyariah)
		if err != nil {
			// log.Error("flag_syariah parameter should be true/false")
			return lib.CustomError(http.StatusBadRequest, "flag_syariah parameter should be true/false", "flag_syariah parameter should be true/false")
		}
		if flagsyariahBool == true {
			params["flag_syariah"] = "1"
		} else {
			params["flag_syariah"] = "0"
		}
	} else {
		// log.Error("flag_syariah parameter should be true/false")
		return lib.CustomError(http.StatusBadRequest, "flag_syariah parameter should be true/false", "flag_syariah parameter should be true/false")
	}

	//max_sub_fee
	maxsubfee := c.FormValue("max_sub_fee")
	if maxsubfee == "" {
		// log.Error("Missing required parameter: max_sub_fee cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: max_sub_fee cann't be blank", "Missing required parameter: max_sub_fee cann't be blank")
	}
	maxsubfeeFloat, err := strconv.ParseFloat(maxsubfee, 64)
	if err == nil {
		if maxsubfeeFloat < 0 {
			// log.Error("Wrong input for parameter: max_sub_fee cann't negatif")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: max_sub_fee must cann't negatif", "Missing required parameter: max_sub_fee cann't negatif")
		}
		params["max_sub_fee"] = maxsubfee
	} else {
		// log.Error("Wrong input for parameter: max_sub_fee number")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: max_sub_fee must number", "Missing required parameter: max_sub_fee number")
	}

	//max_red_fee
	maxredfee := c.FormValue("max_red_fee")
	if maxredfee == "" {
		// log.Error("Missing required parameter: max_red_fee cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: max_red_fee cann't be blank", "Missing required parameter: max_red_fee cann't be blank")
	}
	maxredfeeFloat, err := strconv.ParseFloat(maxredfee, 64)
	if err == nil {
		if maxredfeeFloat < 0 {
			// log.Error("Wrong input for parameter: max_red_fee cann't negatif")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: max_red_fee must cann't negatif", "Missing required parameter: max_red_fee cann't negatif")
		}
		params["max_red_fee"] = maxredfee
	} else {
		// log.Error("Wrong input for parameter: max_red_fee number")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: max_red_fee must number", "Missing required parameter: max_red_fee number")
	}

	//max_swi_fee
	maxswifee := c.FormValue("max_swi_fee")
	if maxswifee == "" {
		// log.Error("Missing required parameter: max_swi_fee cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: max_swi_fee cann't be blank", "Missing required parameter: max_swi_fee cann't be blank")
	}
	maxswifeeFloat, err := strconv.ParseFloat(maxswifee, 64)
	if err == nil {
		if maxswifeeFloat < 0 {
			// log.Error("Wrong input for parameter: max_swi_fee cann't negatif")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: max_swi_fee must cann't negatif", "Missing required parameter: max_swi_fee cann't negatif")
		}
		params["max_swi_fee"] = maxswifee
	} else {
		// log.Error("Wrong input for parameter: max_swi_fee number")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: max_swi_fee must number", "Missing required parameter: max_swi_fee number")
	}

	//min_sub_amount
	minsubamount := c.FormValue("min_sub_amount")
	if minsubamount == "" {
		// log.Error("Missing required parameter: min_sub_amount cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: min_sub_amount cann't be blank", "Missing required parameter: min_sub_amount cann't be blank")
	}
	minsubamountFloat, err := strconv.ParseFloat(minsubamount, 64)
	if err == nil {
		if minsubamountFloat < 0 {
			// log.Error("Wrong input for parameter: min_sub_amount cann't negatif")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: min_sub_amount must cann't negatif", "Missing required parameter: min_sub_amount cann't negatif")
		}
		params["min_sub_amount"] = minsubamount
	} else {
		// log.Error("Wrong input for parameter: min_sub_amount number")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: min_sub_amount must number", "Missing required parameter: min_sub_amount number")
	}

	//min_red_amount
	minredamount := c.FormValue("min_red_amount")
	if minredamount == "" {
		// log.Error("Missing required parameter: min_red_amount cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: min_red_amount cann't be blank", "Missing required parameter: min_red_amount cann't be blank")
	}
	minredamountFloat, err := strconv.ParseFloat(minredamount, 64)
	if err == nil {
		if minredamountFloat < 0 {
			// log.Error("Wrong input for parameter: min_red_amount cann't negatif")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: min_red_amount must cann't negatif", "Missing required parameter: min_red_amount cann't negatif")
		}
		params["min_red_amount"] = minredamount
	} else {
		// log.Error("Wrong input for parameter: min_red_amount number")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: min_red_amount must number", "Missing required parameter: min_red_amount number")
	}

	//min_red_unit
	minredunit := c.FormValue("min_red_unit")
	if minredunit == "" {
		// log.Error("Missing required parameter: min_red_unit cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: min_red_unit cann't be blank", "Missing required parameter: min_red_unit cann't be blank")
	}
	minredunitFloat, err := strconv.ParseFloat(minredunit, 64)
	if err == nil {
		if minredunitFloat < 0 {
			// log.Error("Wrong input for parameter: min_red_unit cann't negatif")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: min_red_unit must cann't negatif", "Missing required parameter: min_red_unit cann't negatif")
		}
		params["min_red_unit"] = minredunit
	} else {
		// log.Error("Wrong input for parameter: min_red_unit number")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: min_red_unit must number", "Missing required parameter: min_red_unit number")
	}

	//min_unit_after_red
	minunitafterred := c.FormValue("min_unit_after_red")
	if minunitafterred == "" {
		// log.Error("Missing required parameter: min_unit_after_red cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: min_unit_after_red cann't be blank", "Missing required parameter: min_unit_after_red cann't be blank")
	}
	minunitafterredFloat, err := strconv.ParseFloat(minunitafterred, 64)
	if err == nil {
		if minunitafterredFloat < 0 {
			// log.Error("Wrong input for parameter: min_unit_after_red cann't negatif")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: min_unit_after_red must cann't negatif", "Missing required parameter: min_unit_after_red cann't negatif")
		}
		params["min_unit_after_red"] = minunitafterred
	} else {
		// log.Error("Wrong input for parameter: min_unit_after_red number")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: min_unit_after_red must number", "Missing required parameter: min_unit_after_red number")
	}

	//management_fee
	managementfee := c.FormValue("management_fee")
	if managementfee == "" {
		// log.Error("Missing required parameter: management_fee cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: management_fee cann't be blank", "Missing required parameter: management_fee cann't be blank")
	}
	managementfeeFloat, err := strconv.ParseFloat(managementfee, 64)
	if err == nil {
		if managementfeeFloat < 0 {
			// log.Error("Wrong input for parameter: management_fee cann't negatif")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: management_fee must cann't negatif", "Missing required parameter: management_fee cann't negatif")
		}
		params["management_fee"] = managementfee
	} else {
		// log.Error("Wrong input for parameter: management_fee number")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: management_fee must number", "Missing required parameter: management_fee number")
	}

	//custodian_fee
	custodianfee := c.FormValue("custodian_fee")
	if custodianfee == "" {
		// log.Error("Missing required parameter: custodian_fee cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: custodian_fee cann't be blank", "Missing required parameter: custodian_fee cann't be blank")
	}
	custodianfeeFloat, err := strconv.ParseFloat(custodianfee, 64)
	if err == nil {
		if custodianfeeFloat < 0 {
			// log.Error("Wrong input for parameter: custodian_fee cann't negatif")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: custodian_fee must cann't negatif", "Missing required parameter: custodian_fee cann't negatif")
		}
		params["custodian_fee"] = custodianfee
	} else {
		// log.Error("Wrong input for parameter: custodian_fee number")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: custodian_fee must number", "Missing required parameter: custodian_fee number")
	}

	//custodian_key
	custodiankey := c.FormValue("custodian_key")
	if custodiankey != "" {
		sub, err = strconv.ParseUint(custodiankey, 10, 64)
		if err == nil && sub > 0 {
			params["custodian_key"] = custodiankey
		} else {
			// log.Error("Wrong input for parameter: custodian_key number")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: custodian_key must number", "Missing required parameter: custodian_key number")
		}
	}

	//ojk_fee
	ojkfee := c.FormValue("ojk_fee")
	if ojkfee == "" {
		// log.Error("Missing required parameter: ojk_fee cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: ojk_fee cann't be blank", "Missing required parameter: ojk_fee cann't be blank")
	}
	ojkfeeFloat, err := strconv.ParseFloat(ojkfee, 64)
	if err == nil {
		if ojkfeeFloat < 0 {
			// log.Error("Wrong input for parameter: ojk_fee cann't negatif")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: ojk_fee must cann't negatif", "Missing required parameter: ojk_fee cann't negatif")
		}
		params["ojk_fee"] = ojkfee
	} else {
		// log.Error("Wrong input for parameter: ojk_fee number")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: ojk_fee must number", "Missing required parameter: ojk_fee number")
	}

	//product_fee_amount
	productfeeamount := c.FormValue("product_fee_amount")
	if productfeeamount == "" {
		// log.Error("Missing required parameter: product_fee_amount cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: product_fee_amount cann't be blank", "Missing required parameter: product_fee_amount cann't be blank")
	}
	productfeeamountFloat, err := strconv.ParseFloat(productfeeamount, 64)
	if err == nil {
		if productfeeamountFloat < 0 {
			// log.Error("Wrong input for parameter: product_fee_amount cann't negatif")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: product_fee_amount must cann't negatif", "Missing required parameter: product_fee_amount cann't negatif")
		}
		params["product_fee_amount"] = productfeeamount
	} else {
		// log.Error("Wrong input for parameter: product_fee_amount number")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: product_fee_amount must number", "Missing required parameter: product_fee_amount number")
	}

	//overwrite_transact_flag
	overwritetransactflag := c.FormValue("overwrite_transact_flag")
	var overwritetransactflagBool bool
	if overwritetransactflag != "" {
		overwritetransactflagBool, err = strconv.ParseBool(overwritetransactflag)
		if err != nil {
			// log.Error("overwrite_transact_flag parameter should be true/false")
			return lib.CustomError(http.StatusBadRequest, "overwrite_transact_flag parameter should be true/false", "overwrite_transact_flag parameter should be true/false")
		}
		if overwritetransactflagBool == true {
			params["overwrite_transact_flag"] = "1"
		} else {
			params["overwrite_transact_flag"] = "0"
		}
	} else {
		// log.Error("overwrite_transact_flag parameter should be true/false")
		return lib.CustomError(http.StatusBadRequest, "overwrite_transact_flag parameter should be true/false", "overwrite_transact_flag parameter should be true/false")
	}

	//overwrite_fee_flag
	overwritefeeflag := c.FormValue("overwrite_fee_flag")
	var overwritefeeflagBool bool
	if overwritefeeflag != "" {
		overwritefeeflagBool, err = strconv.ParseBool(overwritefeeflag)
		if err != nil {
			// log.Error("overwrite_fee_flag parameter should be true/false")
			return lib.CustomError(http.StatusBadRequest, "overwrite_fee_flag parameter should be true/false", "overwrite_fee_flag parameter should be true/false")
		}
		if overwritefeeflagBool == true {
			params["overwrite_fee_flag"] = "1"
		} else {
			params["overwrite_fee_flag"] = "0"
		}
	} else {
		// log.Error("overwrite_fee_flag parameter should be true/false")
		return lib.CustomError(http.StatusBadRequest, "overwrite_fee_flag parameter should be true/false", "overwrite_fee_flag parameter should be true/false")
	}

	//other_fee_amount
	otherfeeamount := c.FormValue("other_fee_amount")
	if otherfeeamount == "" {
		// log.Error("Missing required parameter: other_fee_amount cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: other_fee_amount cann't be blank", "Missing required parameter: other_fee_amount cann't be blank")
	}
	otherfeeamountFloat, err := strconv.ParseFloat(otherfeeamount, 64)
	if err == nil {
		if otherfeeamountFloat < 0 {
			// log.Error("Wrong input for parameter: other_fee_amount cann't negatif")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: other_fee_amount must cann't negatif", "Missing required parameter: other_fee_amount cann't negatif")
		}
		params["other_fee_amount"] = otherfeeamount
	} else {
		// log.Error("Wrong input for parameter: other_fee_amount number")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: other_fee_amount must number", "Missing required parameter: other_fee_amount number")
	}

	//settlement_period
	settlementperiod := c.FormValue("settlement_period")
	if settlementperiod != "" {
		sub, err = strconv.ParseUint(settlementperiod, 10, 64)
		if err == nil {
			params["settlement_period"] = settlementperiod
		} else {
			// log.Error("Wrong input for parameter: settlement_period number")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: settlement_period must number", "Missing required parameter: settlement_period number")
		}
	}

	//sinvest_fund_code
	sinvestfundcode := c.FormValue("sinvest_fund_code")
	if sinvestfundcode != "" {
		params["sinvest_fund_code"] = sinvestfundcode

		//check unique sinvest_fund_code
		paramsSinvest := make(map[string]string)
		paramsSinvest["sinvest_fund_code"] = sinvestfundcode

		status, err = models.AdminGetValidateUniqueDataInsertUpdate(&countDataExisting, paramsSinvest, paramsCheckValidateAnd, &productkey)
		if err != nil {
			// log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
		if int(countDataExisting.CountData) > 0 {
			// log.Error("Missing required parameter: sinvest_fund_code already existing, use other sinvest_fund_code")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: sinvest_fund_code already existing, use other sinvest_fund_code", "Missing required parameter: sinvest_fund_code already existing, use other sinvest_fund_code")
		}
	}

	//flag_enabled
	flagenabled := c.FormValue("flag_enabled")
	var flagenabledBool bool
	if flagenabled != "" {
		flagenabledBool, err = strconv.ParseBool(flagenabled)
		if err != nil {
			// log.Error("flag_enabled parameter should be true/false")
			return lib.CustomError(http.StatusBadRequest, "flag_enabled parameter should be true/false", "flag_enabled parameter should be true/false")
		}
		if flagenabledBool == true {
			params["flag_enabled"] = "1"
		} else {
			params["flag_enabled"] = "0"
		}
	} else {
		// log.Error("flag_enabled parameter should be true/false")
		return lib.CustomError(http.StatusBadRequest, "flag_enabled parameter should be true/false", "flag_enabled parameter should be true/false")
	}

	//flag_subscription
	flagsubscription := c.FormValue("flag_subscription")
	var flagsubscriptionBool bool
	if flagsubscription != "" {
		flagsubscriptionBool, err = strconv.ParseBool(flagsubscription)
		if err != nil {
			// log.Error("flag_subscription parameter should be true/false")
			return lib.CustomError(http.StatusBadRequest, "flag_subscription parameter should be true/false", "flag_subscription parameter should be true/false")
		}
		if flagsubscriptionBool == true {
			params["flag_subscription"] = "1"
		} else {
			params["flag_subscription"] = "0"
		}
	} else {
		// log.Error("flag_subscription parameter should be true/false")
		return lib.CustomError(http.StatusBadRequest, "flag_subscription parameter should be true/false", "flag_subscription parameter should be true/false")
	}

	//flag_redemption
	flagredemption := c.FormValue("flag_redemption")
	var flagredemptionBool bool
	if flagredemption != "" {
		flagredemptionBool, err = strconv.ParseBool(flagredemption)
		if err != nil {
			// log.Error("flag_redemption parameter should be true/false")
			return lib.CustomError(http.StatusBadRequest, "flag_redemption parameter should be true/false", "flag_redemption parameter should be true/false")
		}
		if flagredemptionBool == true {
			params["flag_redemption"] = "1"
		} else {
			params["flag_redemption"] = "0"
		}
	} else {
		// log.Error("flag_redemption parameter should be true/false")
		return lib.CustomError(http.StatusBadRequest, "flag_redemption parameter should be true/false", "flag_redemption parameter should be true/false")
	}

	//flag_switch_out
	flagswitchout := c.FormValue("flag_switch_out")
	var flagswitchoutBool bool
	if flagswitchout != "" {
		flagswitchoutBool, err = strconv.ParseBool(flagswitchout)
		if err != nil {
			// log.Error("flag_switch_out parameter should be true/false")
			return lib.CustomError(http.StatusBadRequest, "flag_switch_out parameter should be true/false", "flag_switch_out parameter should be true/false")
		}
		if flagswitchoutBool == true {
			params["flag_switch_out"] = "1"
		} else {
			params["flag_switch_out"] = "0"
		}
	} else {
		// log.Error("flag_switch_out parameter should be true/false")
		return lib.CustomError(http.StatusBadRequest, "flag_switch_out parameter should be true/false", "flag_switch_out parameter should be true/false")
	}

	//flag_switch_in
	flagswitchin := c.FormValue("flag_switch_in")
	var flagswitchinBool bool
	if flagswitchin != "" {
		flagswitchinBool, err = strconv.ParseBool(flagswitchin)
		if err != nil {
			// log.Error("flag_switch_in parameter should be true/false")
			return lib.CustomError(http.StatusBadRequest, "flag_switch_in parameter should be true/false", "flag_switch_in parameter should be true/false")
		}
		if flagswitchinBool == true {
			params["flag_switch_in"] = "1"
		} else {
			params["flag_switch_in"] = "0"
		}
	} else {
		// log.Error("flag_switch_in parameter should be true/false")
		return lib.CustomError(http.StatusBadRequest, "flag_switch_in parameter should be true/false", "flag_switch_in parameter should be true/false")
	}

	dateLayout := "2006-01-02 15:04:05"
	params["rec_status"] = "1"
	params["rec_order"] = "0"
	params["rec_modified_date"] = time.Now().Format(dateLayout)
	params["rec_modified_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

	status, err = models.UpdateMsProduct(params)
	if err != nil {
		// log.Error("Failed create request data: " + err.Error())
		return lib.CustomError(status, err.Error(), "failed input data")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = nil
	return c.JSON(http.StatusOK, response)

}

func GetListProductAdminDropdown(c echo.Context) error {

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

	items := []string{"product_key", "product_code", "product_name", "launch_date", "inception_date", "isin_code", "flag_syariah", "sinvest_fund_code"}

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
		params["orderBy"] = "product_key"
		params["orderType"] = "ASC"
	}

	params["rec_status"] = "1"

	paramsLike := make(map[string]string)

	var msProduct []models.MsProduct

	status, err = models.AdminGetAllMsProductWithLike(&msProduct, limit, offset, params, paramsLike, noLimit)

	if err != nil {
		// log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if len(msProduct) < 1 {
		// log.Error("product not found")
		return lib.CustomError(http.StatusNotFound, "Product not found", "Product not found")
	}

	var responseData []models.MsProductListDropdown
	for _, pro := range msProduct {
		var data models.MsProductListDropdown
		data.ProductKey = pro.ProductKey
		data.ProductCode = pro.ProductCode
		data.ProductName = pro.ProductName
		responseData = append(responseData, data)
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = responseData

	return c.JSON(http.StatusOK, response)
}

func AdminGetProductSubscription(c echo.Context) error {
	var err error
	var status int

	//fund_type_key
	fundtype := c.Param("fund_type_key")
	if fundtype != "" {
		sub, err := strconv.ParseUint(fundtype, 10, 64)
		if err != nil || sub == 0 {
			// log.Error("Wrong input for parameter: fund_type_key number")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: fund_type_key must number", "Missing required parameter: fund_type_key number")
		}
	} else {
		// log.Error("Wrong input for parameter: fund_type_key number")
		return lib.CustomError(http.StatusBadRequest, "fund_type_key required", "fund_type_key required")
	}

	var products []models.ProductSubscriptionFundType
	status, err = models.AdminGetProductSubscription(&products, fundtype)

	if err != nil {
		// log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if len(products) < 1 {
		// log.Error("Product not found")
		return lib.CustomError(http.StatusNotFound, "Product not found", "Product not found")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = products

	return c.JSON(http.StatusOK, response)
}

func GetBankProductSubscription(c echo.Context) error {

	var err error
	var status int

	productStr := c.Param("product_key")
	key, _ := strconv.ParseUint(productStr, 10, 64)
	if key == 0 {
		return lib.CustomError(http.StatusNotFound)
	}

	lookupTransType := "269"

	var bankAccountTransactionInfo []models.MsProductBankAccountTransactionInfo

	status, err = models.GetAllMsProductBankAccountTransaction(&bankAccountTransactionInfo, productStr, lookupTransType)
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

func AdminGetProductRedemption(c echo.Context) error {
	var err error
	var status int
	decimal.MarshalJSONWithoutQuotes = true

	customerKeyStr := c.Param("customer_key")
	var cus models.MsCustomer
	if customerKeyStr != "" {
		customerKey, err := strconv.ParseUint(customerKeyStr, 10, 64)
		if err == nil && customerKey > 0 {
			status, err = models.GetMsCustomer(&cus, customerKeyStr)
			if err != nil {
				// log.Error(err.Error())
				return lib.CustomError(http.StatusBadRequest, err.Error(), "Customer tidak ditemukan")
			}
			if cus.CifSuspendFlag == uint8(1) {
				// log.Error("Customer Suspended")
				return lib.CustomError(http.StatusBadRequest, "Customer Suspended", "Customer Suspended")
			}
		} else {
			// log.Error("Wrong input for parameter: customer_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: customer_key", "Wrong input for parameter: customer_key")
		}
	} else {
		// log.Error("Missing required parameter: customer_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: customer_key", "Missing required parameter: customer_key")
	}

	var products []models.ProductRedemption
	status, err = models.AdminGetProductRedemption(&products, customerKeyStr)

	if err != nil {
		// log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if len(products) < 1 {
		// log.Error("Product not found")
		return lib.CustomError(http.StatusNotFound, "Product not found", "Product not found")
	}

	var acaIds []string
	for _, pr := range products {
		if _, ok := lib.Find(acaIds, strconv.FormatUint(pr.AcaKey, 10)); !ok {
			acaIds = append(acaIds, strconv.FormatUint(pr.AcaKey, 10))
		}
	}

	//mapping Balance Unit
	var balances []models.SumBalanceUnit
	if len(acaIds) > 0 {
		status, err = models.GetSumBalanceUnit(&balances, acaIds)
		if err != nil {
			// log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
	}
	balanceData := make(map[uint64]models.SumBalanceUnit)
	for _, b := range balances {
		balanceData[b.AcaKey] = b
	}

	agentName := ""

	var agent models.CustomerAgent
	status, err = models.GetCustomerLastAgent(&agent, customerKeyStr)
	if err == nil {
		if agent.AgentName != nil {
			agentName = *agent.AgentCode + " - " + *agent.AgentName
		}
	}

	var productList []models.ProductRedemption
	zero := decimal.NewFromInt(0)
	for _, pr := range products {
		if n, ok := balanceData[pr.AcaKey]; ok {
			if n.Unit.Cmp(zero) == 1 {
				var prod models.ProductRedemption
				prod.ProductKey = pr.ProductKey
				prod.RiskProfileKey = pr.RiskProfileKey
				prod.FlagRedemption = pr.FlagRedemption
				prod.FlagSwitchOut = pr.FlagSwitchOut
				prod.FundTypeName = pr.FundTypeName
				prod.ProductName = pr.ProductName
				prod.NavDate = pr.NavDate
				prod.NavValue = pr.NavValue.Truncate(2)
				prod.MinRedAmount = pr.MinRedAmount.Truncate(2)
				prod.MinRedUnit = pr.MinRedUnit.Truncate(2)
				prod.MinUnitAfterRed = pr.MinUnitAfterRed.Truncate(2)
				prod.RiskName = pr.RiskName
				prod.AcaKey = pr.AcaKey
				prod.Unit = n.Unit.Truncate(2)
				prod.NilaiInvestasi = n.Unit.Mul(prod.NavValue).Truncate(0)
				prod.SalesName = &agentName
				productList = append(productList, prod)
			}
		}
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = productList

	return c.JSON(http.StatusOK, response)
}

func GetProductDetailTransactionSubscription(c echo.Context) error {
	decimal.MarshalJSONWithoutQuotes = true

	productKeyStr := c.Param("product_key")
	var product models.ProductSubscription
	if productKeyStr != "" {
		productKey, err := strconv.ParseUint(productKeyStr, 10, 64)
		if err == nil && productKey > 0 {
			_, err = models.AdminGetProductSubscriptionByProductKey(&product, productKeyStr)
			if err != nil {
				// log.Error(err.Error())
				return lib.CustomError(http.StatusBadRequest, err.Error(), "Product tidak ditemukan")
			}
		} else {
			// log.Error("Wrong input for parameter: product_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: product_key", "Wrong input for parameter: product_key")
		}
	} else {
		// log.Error("Missing required parameter: product_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: product_key", "Missing required parameter: product_key")
	}

	var productResponse models.ProductSubscription
	productResponse.ProductKey = product.ProductKey
	productResponse.FundTypeName = product.FundTypeName
	productResponse.ProductName = product.ProductName
	productResponse.NavDate = product.NavDate
	productResponse.NavValue = product.NavValue.Truncate(2)
	productResponse.ProductImage = product.ProductImage
	productResponse.MinSubAmount = product.MinSubAmount.Truncate(2)
	productResponse.MinRedAmount = product.MinRedAmount.Truncate(2)
	productResponse.MinRedUnit = product.MinRedUnit.Truncate(2)
	productResponse.ProspectusLink = product.ProspectusLink
	productResponse.FfsLink = product.FfsLink
	productResponse.RiskName = product.RiskName
	productResponse.CurrencyKey = product.CurrencyKey
	productResponse.Symbol = product.Symbol
	productResponse.CurrencyCode = product.CurrencyCode
	productResponse.CurrencyName = product.CurrencyName
	productResponse.FlagShowOntnc = product.FlagShowOntnc
	productResponse.FeeAnnotation = product.FeeAnnotation
	productResponse.FeeValue = product.FeeValue.Truncate(2)

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = productResponse
	return c.JSON(http.StatusOK, response)
}

func AdminGetProductSwitchIn(c echo.Context) error {
	var err error
	var status int
	decimal.MarshalJSONWithoutQuotes = true

	customerKeyStr := c.Param("customer_key")
	var cus models.MsCustomer
	if customerKeyStr != "" {
		customerKey, err := strconv.ParseUint(customerKeyStr, 10, 64)
		if err == nil && customerKey > 0 {
			status, err = models.GetMsCustomer(&cus, customerKeyStr)
			if err != nil {
				// log.Error(err.Error())
				return lib.CustomError(http.StatusBadRequest, err.Error(), "Customer tidak ditemukan")
			}
			if cus.CifSuspendFlag == uint8(1) {
				// log.Error("Customer Suspended")
				return lib.CustomError(http.StatusBadRequest, "Customer Suspended", "Customer Suspended")
			}
		} else {
			// log.Error("Wrong input for parameter: customer_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: customer_key", "Wrong input for parameter: customer_key")
		}
	} else {
		// log.Error("Missing required parameter: customer_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: customer_key", "Missing required parameter: customer_key")
	}

	productSwOutKey := c.Param("product_switch_out_key")
	if productSwOutKey == "" {
		// log.Error("Missing required parameter: product_switch_out_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: product_switch_out_key", "Missing required parameter: product_switch_out_key")
	} else {
		prodKey, err := strconv.ParseUint(productSwOutKey, 10, 64)
		if err != nil || prodKey == 0 {
			// log.Error("Wrong input for parameter: product_switch_out_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: product_switch_out_key", "Wrong input for parameter: product_switch_out_key")
		}
	}

	var products []models.ProductHaveBalanceSwitchIn
	status, err = models.AdminGetProductHaveBalanceSwitchIn(&products, customerKeyStr, productSwOutKey)

	if err != nil {
		if err != sql.ErrNoRows {
			// log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
	}

	var acaIds []string
	if len(products) > 0 {
		for _, pr := range products {
			if pr.AcaKey != nil {
				if _, ok := lib.Find(acaIds, strconv.FormatUint(*pr.AcaKey, 10)); !ok {
					acaIds = append(acaIds, strconv.FormatUint(*pr.AcaKey, 10))
				}
			}
		}
	}

	//mapping Balance Unit
	var balances []models.SumBalanceUnit
	if len(acaIds) > 0 {
		status, err = models.GetSumBalanceUnit(&balances, acaIds)
		if err != nil {
			// log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
	}
	balanceData := make(map[uint64]models.SumBalanceUnit)
	if len(balances) > 0 {
		for _, b := range balances {
			balanceData[b.AcaKey] = b
		}
	}

	var productKeyExis []string
	productKeyExis = append(productKeyExis, productSwOutKey)

	var productList []models.ProductHaveBalanceSwitchIn
	zero := decimal.NewFromInt(0)
	if len(products) > 0 {
		for _, pr := range products {
			if pr.AcaKey != nil {
				if n, ok := balanceData[*pr.AcaKey]; ok {
					if n.Unit.Cmp(zero) == 1 {
						if _, ok := lib.Find(productKeyExis, strconv.FormatUint(pr.ProductKey, 10)); !ok {
							productKeyExis = append(productKeyExis, strconv.FormatUint(pr.ProductKey, 10))
						}
						var prod models.ProductHaveBalanceSwitchIn
						prod.ProductKey = pr.ProductKey
						prod.RiskProfileKey = pr.RiskProfileKey
						prod.FlagSwitchIn = pr.FlagSwitchIn
						prod.FundTypeName = pr.FundTypeName
						prod.ProductName = pr.ProductName
						prod.NavDate = pr.NavDate
						prod.NavValue = pr.NavValue.Truncate(2)
						prod.MinSubAmount = pr.MinSubAmount.Truncate(2)
						prod.RiskName = pr.RiskName
						prod.AcaKey = pr.AcaKey
						prod.Unit = n.Unit.Truncate(2)
						prod.NilaiInvestasi = n.Unit.Mul(prod.NavValue).Truncate(0)
						productList = append(productList, prod)
					}
				}
			}
		}
	}

	var productsElse []models.ProductHaveBalanceSwitchIn
	status, err = models.AdminGetProductNotInBalanceSwitchIn(&productsElse, productKeyExis)
	if err != nil {
		if err != sql.ErrNoRows {
			// log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
	}
	if len(productsElse) > 0 {
		for _, prd := range productsElse {
			productList = append(productList, prd)
		}
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = productList

	return c.JSON(http.StatusOK, response)
}
