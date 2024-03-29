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
	"github.com/leekchan/accounting"
	"github.com/shopspring/decimal"
)

func GetListTrNavAdmin(c echo.Context) error {

	var err error
	var status int
	decimal.MarshalJSONWithoutQuotes = true

	errorAuth := IsMenuAccessAllowed(39)
	if errorAuth != nil {
		// log.Error("User Autorizer")
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}

	params := make(map[string]string)

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

	navdate := c.QueryParam("nav_date")
	if navdate != "" {
		params["nav_date"] = navdate
	}

	if (productKey == "") && (navdate == "") {
		// log.Error("Wrong input for parameter: product_key atau nav_date harus salah satu diisi")
		return lib.CustomError(http.StatusBadRequest, "Mohon input Produk atau Nav Date", "Mohon input Produk atau Nav Date")
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

	items := []string{"nav_key", "product_key", "nav_date", "nav_value", "original_value", "nav_status", "prod_aum_total", "prod_unit_total", "publish_mode"}

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
		params["orderBy"] = "nav_date"
		params["orderType"] = "DESC"
	}

	params["rec_status"] = "1"

	var trNav []models.TrNav

	status, err = models.GetAllTrNav(&trNav, limit, offset, params, noLimit)

	if err != nil {
		if err != sql.ErrNoRows {
			// log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
	}
	if len(trNav) < 1 {
		// log.Error("nav not found")
		return lib.CustomError(http.StatusNotFound, "Nav not found", "Nav not found")
	}

	var genLookupIds []string
	var productIds []string
	for _, nav := range trNav {
		if _, ok := lib.Find(genLookupIds, strconv.FormatUint(nav.NavStatus, 10)); !ok {
			genLookupIds = append(genLookupIds, strconv.FormatUint(nav.NavStatus, 10))
		}
		if _, ok := lib.Find(genLookupIds, strconv.FormatUint(nav.PublishMode, 10)); !ok {
			genLookupIds = append(genLookupIds, strconv.FormatUint(nav.PublishMode, 10))
		}
		if _, ok := lib.Find(productIds, strconv.FormatUint(nav.ProductKey, 10)); !ok {
			productIds = append(productIds, strconv.FormatUint(nav.ProductKey, 10))
		}
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

	//gen msProduct
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

	proData := make(map[uint64]models.MsProduct)
	for _, pro := range msProduct {
		proData[pro.ProductKey] = pro
	}

	var responseData []models.TrNavList
	for _, nv := range trNav {
		var data models.TrNavList

		data.NavKey = nv.NavKey
		if n, ok := proData[nv.ProductKey]; ok {
			data.ProductName = n.ProductNameAlt
		}

		layout := "2006-01-02 15:04:05"
		newLayout := "02 Jan 2006"
		date, err := time.Parse(layout, nv.NavDate)
		if err == nil {
			data.NavDate = date.Format(newLayout)
		}

		ac := accounting.Accounting{Symbol: "", Precision: 4, Thousand: ".", Decimal: ","}

		data.NavValue = ac.FormatMoney(nv.NavValue)
		data.OriginalValue = ac.FormatMoney(nv.OriginalValue)
		if n, ok := gData[nv.NavStatus]; ok {
			data.NavStatus = n.LkpName
		}
		data.ProdAumTotal = ac.FormatMoney(nv.ProdAumTotal)
		data.ProdUnitTotal = ac.FormatMoney(nv.ProdUnitTotal)
		if n, ok := gData[nv.PublishMode]; ok {
			data.PublishMode = n.LkpName
		}

		responseData = append(responseData, data)
	}

	var countData models.CountData
	var pagination int
	if limit > 0 {
		status, err = models.GetAllTrNavCount(&countData, params)
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

func GetNavDetailAdmin(c echo.Context) error {
	var err error
	var status int
	decimal.MarshalJSONWithoutQuotes = true

	errorAuth := IsMenuAccessAllowed(39)
	if errorAuth != nil {
		// log.Error("User Autorizer")
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}

	keyStr := c.Param("key")
	key, err := strconv.ParseUint(keyStr, 10, 64)
	if err == nil && key == 0 {
		// log.Error("Wrong input for parameter: nav_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: key", "Missing required parameter: key")
	}

	var trNav models.TrNav
	status, err = models.GetTrNavByKey(&trNav, keyStr)
	if err != nil {
		return lib.CustomError(http.StatusNotFound)
	}

	var responseData models.TrNavDetail

	var lookupIds []string

	if _, ok := lib.Find(lookupIds, strconv.FormatUint(trNav.NavStatus, 10)); !ok {
		lookupIds = append(lookupIds, strconv.FormatUint(trNav.NavStatus, 10))
	}
	if _, ok := lib.Find(lookupIds, strconv.FormatUint(trNav.PublishMode, 10)); !ok {
		lookupIds = append(lookupIds, strconv.FormatUint(trNav.PublishMode, 10))
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

	keyProductStr := strconv.FormatUint(trNav.ProductKey, 10)
	var product models.MsProduct
	status, err = models.GetMsProduct(&product, keyProductStr)
	if err != nil {
		return lib.CustomError(http.StatusNotFound)
	}

	responseData.NavKey = trNav.NavKey

	//set product
	var pro models.MsProductInfo
	pro.ProductKey = product.ProductKey
	pro.ProductCode = product.ProductCode
	pro.ProductName = product.ProductName
	pro.ProductNameAlt = product.ProductNameAlt
	responseData.Product = pro

	layout := "2006-01-02 15:04:05"
	newLayout := "02 Jan 2006"
	date, err := time.Parse(layout, trNav.NavDate)
	if err == nil {
		responseData.NavDate = date.Format(newLayout)
	}

	responseData.NavValue = trNav.NavValue
	responseData.OriginalValue = trNav.OriginalValue
	responseData.NavValue = trNav.NavValue
	responseData.ProdAumTotal = trNav.ProdAumTotal
	responseData.ProdUnitTotal = trNav.ProdUnitTotal

	if n, ok := gData[trNav.NavStatus]; ok {
		var trc models.LookupTrans

		trc.LookupKey = n.LookupKey
		trc.LkpGroupKey = n.LkpGroupKey
		trc.LkpCode = n.LkpCode
		trc.LkpName = n.LkpName
		responseData.NavStatus = trc
	}

	if n, ok := gData[trNav.PublishMode]; ok {
		var trc models.LookupTrans

		trc.LookupKey = n.LookupKey
		trc.LkpGroupKey = n.LkpGroupKey
		trc.LkpCode = n.LkpCode
		trc.LkpName = n.LkpName
		responseData.PublishMode = trc
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = responseData

	return c.JSON(http.StatusOK, response)
}

func DeleteNavAdmin(c echo.Context) error {
	var err error
	decimal.MarshalJSONWithoutQuotes = true

	errorAuth := IsMenuAccessAllowed(39)
	if errorAuth != nil {
		// log.Error("User Autorizer")
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}

	params := make(map[string]string)

	navKey := c.FormValue("nav_key")
	if navKey == "" {
		// log.Error("Missing required parameter: nav_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: nav_key", "Missing required parameter: nav_key")
	}

	navKeyCek, err := strconv.ParseUint(navKey, 10, 64)
	if err == nil && navKeyCek > 0 {
		params["nav_key"] = navKey
	} else {
		// log.Error("Wrong input for parameter: nav_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: nav_key", "Missing required parameter: nav_key")
	}

	var trNav models.TrNav
	_, err = models.GetTrNavByKey(&trNav, navKey)
	if err != nil {
		return lib.CustomError(http.StatusNotFound)
	}

	dateLayout := "2006-01-02 15:04:05"
	params["rec_status"] = "0"
	params["rec_deleted_date"] = time.Now().Format(dateLayout)
	params["rec_deleted_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

	status, err := models.UpdateTrNav(params)
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

func CreateAdminTrNav(c echo.Context) error {
	var err error
	var status int

	errorAuth := IsMenuAccessAllowed(39)
	if errorAuth != nil {
		// log.Error("User Autorizer")
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}

	params := make(map[string]string)

	//product_key
	productkey := c.FormValue("product_key")
	if productkey == "" {
		// log.Error("Missing required parameter: product_key cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: product_key cann't be blank", "Missing required parameter: product_key cann't be blank")
	}

	sub, err := strconv.ParseUint(productkey, 10, 64)
	if err == nil && sub > 0 {
		params["product_key"] = productkey
	} else {
		// log.Error("Wrong input for parameter: product_key number")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: product_key must number", "Missing required parameter: product_key number")
	}

	//nav_date
	navdate := c.FormValue("nav_date")
	if navdate == "" {
		// log.Error("Missing required parameter: nav_date cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: nav_date cann't be blank", "Missing required parameter: nav_date cann't be blank")
	}
	params["nav_date"] = navdate + " 00:00:00"

	//cek exist by nav_date & product_key

	var trNav models.TrNav
	_, err = models.GetNavByProductKeyAndNavDate(&trNav, productkey, navdate)
	if err == nil {
		// log.Error("Missing required parameter: Data nav dengan product dan nav date tersebut sudah ada.")
		return lib.CustomError(http.StatusBadRequest, "Data nav dengan product dan nav date tersebut sudah ada.", "Data nav dengan product dan nav date tersebut sudah ada.")
	}

	//nav_value
	navvalue := c.FormValue("nav_value")
	if navvalue == "" {
		// log.Error("Missing required parameter: nav_value cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: nav_value cann't be blank", "Missing required parameter: nav_value cann't be blank")
	}
	navvalueFloat, err := strconv.ParseFloat(navvalue, 64)
	if err == nil {
		if navvalueFloat < 0 {
			// log.Error("Wrong input for parameter: nav_value cann't negatif")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: nav_value must cann't negatif", "Missing required parameter: nav_value cann't negatif")
		}
		params["nav_value"] = navvalue
	} else {
		// log.Error("Wrong input for parameter: nav_value number")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: nav_value must number", "Missing required parameter: nav_value number")
	}

	//original_value
	originalvalue := c.FormValue("original_value")
	if originalvalue == "" {
		// log.Error("Missing required parameter: original_value cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: original_value cann't be blank", "Missing required parameter: original_value cann't be blank")
	}
	originalvalueFloat, err := strconv.ParseFloat(originalvalue, 64)
	if err == nil {
		if originalvalueFloat < 0 {
			// log.Error("Wrong input for parameter: original_value cann't negatif")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: original_value must cann't negatif", "Missing required parameter: original_value cann't negatif")
		}
		params["original_value"] = originalvalue
	} else {
		// log.Error("Wrong input for parameter: original_value number")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: original_value must number", "Missing required parameter: original_value number")
	}

	//nav_status
	navstatus := c.FormValue("nav_status")
	if navstatus == "" {
		// log.Error("Missing required parameter: nav_status cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: nav_status cann't be blank", "Missing required parameter: nav_status cann't be blank")
	}

	sub, err = strconv.ParseUint(navstatus, 10, 64)
	if err == nil && sub > 0 {
		params["nav_status"] = navstatus
	} else {
		// log.Error("Wrong input for parameter: nav_status number")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: nav_status must number", "Missing required parameter: nav_status number")
	}

	//prod_aum_total
	prodaumtotal := c.FormValue("prod_aum_total")
	if prodaumtotal == "" {
		// log.Error("Missing required parameter: prod_aum_total cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: prod_aum_total cann't be blank", "Missing required parameter: prod_aum_total cann't be blank")
	}
	prodaumtotalFloat, err := strconv.ParseFloat(prodaumtotal, 64)
	if err == nil {
		if prodaumtotalFloat < 0 {
			// log.Error("Wrong input for parameter: prod_aum_total cann't negatif")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: prod_aum_total must cann't negatif", "Missing required parameter: prod_aum_total cann't negatif")
		}
		params["prod_aum_total"] = prodaumtotal
	} else {
		// log.Error("Wrong input for parameter: prod_aum_total number")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: prod_aum_total must number", "Missing required parameter: prod_aum_total number")
	}

	//prod_unit_total
	produnittotal := c.FormValue("prod_unit_total")
	if produnittotal == "" {
		// log.Error("Missing required parameter: prod_unit_total cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: prod_unit_total cann't be blank", "Missing required parameter: prod_unit_total cann't be blank")
	}
	produnittotalFloat, err := strconv.ParseFloat(produnittotal, 64)
	if err == nil {
		if produnittotalFloat < 0 {
			// log.Error("Wrong input for parameter: prod_unit_total cann't negatif")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: prod_unit_total must cann't negatif", "Missing required parameter: prod_unit_total cann't negatif")
		}
		params["prod_unit_total"] = produnittotal
	} else {
		// log.Error("Wrong input for parameter: prod_unit_total number")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: prod_unit_total must number", "Missing required parameter: prod_unit_total number")
	}

	//publish_mode
	publishmode := c.FormValue("publish_mode")
	if publishmode == "" {
		// log.Error("Missing required parameter: publish_mode cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: publish_mode cann't be blank", "Missing required parameter: publish_mode cann't be blank")
	}

	sub, err = strconv.ParseUint(publishmode, 10, 64)
	if err == nil && sub > 0 {
		params["publish_mode"] = publishmode
	} else {
		// log.Error("Wrong input for parameter: publish_mode number")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: publish_mode must number", "Missing required parameter: publish_mode number")
	}

	dateLayout := "2006-01-02 15:04:05"
	params["rec_order"] = "0"
	params["rec_status"] = "1"
	params["rec_created_date"] = time.Now().Format(dateLayout)
	params["rec_created_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

	status, err = models.CreateTrNav(params)
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

func UpdateAdminTrNav(c echo.Context) error {
	var err error
	var status int

	errorAuth := IsMenuAccessAllowed(39)
	if errorAuth != nil {
		// log.Error("User Autorizer")
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}

	params := make(map[string]string)

	navKey := c.FormValue("nav_key")
	if navKey == "" {
		// log.Error("Missing required parameter: nav_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: nav_key", "Missing required parameter: nav_key")
	}

	navKeyCek, err := strconv.ParseUint(navKey, 10, 64)
	if err == nil && navKeyCek > 0 {
		params["nav_key"] = navKey
	} else {
		// log.Error("Wrong input for parameter: nav_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: nav_key", "Missing required parameter: nav_key")
	}

	var trNav models.TrNav
	_, err = models.GetTrNavByKey(&trNav, navKey)
	if err != nil {
		return lib.CustomError(http.StatusNotFound)
	}

	//product_key
	productkey := c.FormValue("product_key")
	if productkey == "" {
		// log.Error("Missing required parameter: product_key cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: product_key cann't be blank", "Missing required parameter: product_key cann't be blank")
	}

	sub, err := strconv.ParseUint(productkey, 10, 64)
	if err == nil && sub > 0 {
		params["product_key"] = productkey
	} else {
		// log.Error("Wrong input for parameter: product_key number")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: product_key must number", "Missing required parameter: product_key number")
	}

	//nav_date
	navdate := c.FormValue("nav_date")
	if navdate == "" {
		// log.Error("Missing required parameter: nav_date cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: nav_date cann't be blank", "Missing required parameter: nav_date cann't be blank")
	}
	params["nav_date"] = navdate + " 00:00:00"

	//cek exist by nav_date & product_key except nav_key
	var trNavCheck models.TrNav
	_, err = models.GetNavByProductKeyAndNavDateExcept(&trNavCheck, productkey, navdate, navKey)
	if err == nil {
		// log.Error("Missing required parameter: Data nav dengan product dan nav date tersebut sudah ada.")
		return lib.CustomError(http.StatusBadRequest, "Data nav dengan product dan nav date tersebut sudah ada.", "Data nav dengan product dan nav date tersebut sudah ada.")
	}

	//nav_value
	navvalue := c.FormValue("nav_value")
	if navvalue == "" {
		// log.Error("Missing required parameter: nav_value cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: nav_value cann't be blank", "Missing required parameter: nav_value cann't be blank")
	}
	navvalueFloat, err := strconv.ParseFloat(navvalue, 64)
	if err == nil {
		if navvalueFloat < 0 {
			// log.Error("Wrong input for parameter: nav_value cann't negatif")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: nav_value must cann't negatif", "Missing required parameter: nav_value cann't negatif")
		}
		params["nav_value"] = navvalue
	} else {
		// log.Error("Wrong input for parameter: nav_value number")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: nav_value must number", "Missing required parameter: nav_value number")
	}

	//original_value
	originalvalue := c.FormValue("original_value")
	if originalvalue == "" {
		// log.Error("Missing required parameter: original_value cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: original_value cann't be blank", "Missing required parameter: original_value cann't be blank")
	}
	originalvalueFloat, err := strconv.ParseFloat(originalvalue, 64)
	if err == nil {
		if originalvalueFloat < 0 {
			// log.Error("Wrong input for parameter: original_value cann't negatif")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: original_value must cann't negatif", "Missing required parameter: original_value cann't negatif")
		}
		params["original_value"] = originalvalue
	} else {
		// log.Error("Wrong input for parameter: original_value number")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: original_value must number", "Missing required parameter: original_value number")
	}

	//nav_status
	navstatus := c.FormValue("nav_status")
	if navstatus == "" {
		// log.Error("Missing required parameter: nav_status cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: nav_status cann't be blank", "Missing required parameter: nav_status cann't be blank")
	}

	sub, err = strconv.ParseUint(navstatus, 10, 64)
	if err == nil && sub > 0 {
		params["nav_status"] = navstatus
	} else {
		// log.Error("Wrong input for parameter: nav_status number")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: nav_status must number", "Missing required parameter: nav_status number")
	}

	//prod_aum_total
	prodaumtotal := c.FormValue("prod_aum_total")
	if prodaumtotal == "" {
		// log.Error("Missing required parameter: prod_aum_total cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: prod_aum_total cann't be blank", "Missing required parameter: prod_aum_total cann't be blank")
	}
	prodaumtotalFloat, err := strconv.ParseFloat(prodaumtotal, 64)
	if err == nil {
		if prodaumtotalFloat < 0 {
			// log.Error("Wrong input for parameter: prod_aum_total cann't negatif")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: prod_aum_total must cann't negatif", "Missing required parameter: prod_aum_total cann't negatif")
		}
		params["prod_aum_total"] = prodaumtotal
	} else {
		// log.Error("Wrong input for parameter: prod_aum_total number")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: prod_aum_total must number", "Missing required parameter: prod_aum_total number")
	}

	//prod_unit_total
	produnittotal := c.FormValue("prod_unit_total")
	if produnittotal == "" {
		// log.Error("Missing required parameter: prod_unit_total cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: prod_unit_total cann't be blank", "Missing required parameter: prod_unit_total cann't be blank")
	}
	produnittotalFloat, err := strconv.ParseFloat(produnittotal, 64)
	if err == nil {
		if produnittotalFloat < 0 {
			// log.Error("Wrong input for parameter: prod_unit_total cann't negatif")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: prod_unit_total must cann't negatif", "Missing required parameter: prod_unit_total cann't negatif")
		}
		params["prod_unit_total"] = produnittotal
	} else {
		// log.Error("Wrong input for parameter: prod_unit_total number")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: prod_unit_total must number", "Missing required parameter: prod_unit_total number")
	}

	//publish_mode
	publishmode := c.FormValue("publish_mode")
	if publishmode == "" {
		// log.Error("Missing required parameter: publish_mode cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: publish_mode cann't be blank", "Missing required parameter: publish_mode cann't be blank")
	}

	sub, err = strconv.ParseUint(publishmode, 10, 64)
	if err == nil && sub > 0 {
		params["publish_mode"] = publishmode
	} else {
		// log.Error("Wrong input for parameter: publish_mode number")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: publish_mode must number", "Missing required parameter: publish_mode number")
	}

	dateLayout := "2006-01-02 15:04:05"
	params["rec_modified_date"] = time.Now().Format(dateLayout)
	params["rec_modified_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

	status, err = models.UpdateTrNav(params)
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

func GetNavPrice(c echo.Context) error {
	// log.Println("SINI")
	decimal.MarshalJSONWithoutQuotes = true
	productKey := c.QueryParam("product_key")
	if productKey == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing product_key", "Missing product_key")
	}
	navDate := c.QueryParam("nav_date")
	if navDate == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing nav_date", "Missing nav_date")
	}

	var trNav models.TrNav
	_, err := models.GetNavByProductKeyAndNavDate(&trNav, productKey, navDate)
	if err != nil {
		return lib.CustomError(http.StatusBadGateway, err.Error(), err.Error())
	}
	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = trNav
	return c.JSON(http.StatusOK, response)
}
