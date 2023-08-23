package controllers

import (
	"fmt"
	"math"
	"mf-bo-api/config"
	"mf-bo-api/lib"
	"mf-bo-api/models"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/labstack/echo"
	"github.com/shopspring/decimal"
)

func GetListCustomerIndividuStatusSuspend(c echo.Context) error {

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

	items := []string{"cif", "full_name", "sid", "date_birth", "customer_key", "phone_mobile", "email_address", "cif_suspend_flag", "ktp"}

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
		params["orderBy"] = "customer_key"
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
		paramsLike["c.full_name"] = fullname
	}
	ktp := c.QueryParam("ktp")
	if ktp != "" {
		paramsLike["pd.idcard_no"] = ktp
	}
	branchKey := c.QueryParam("branch_key")
	if branchKey != "" {
		params["c.openacc_branch_key"] = branchKey
	}
	suspendFlag := c.QueryParam("cif_suspend_flag")
	if suspendFlag != "" {
		params["c.cif_suspend_flag"] = suspendFlag
	}

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

	var customers []models.CustomerIndividuStatusSuspend

	status, err = models.AdminGetAllCustomerStatusSuspend(&customers, limit, offset, params, paramsLike, noLimit)

	if err != nil {
		// log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if len(customers) < 1 {
		// log.Error("Customer not found")
		return lib.CustomError(http.StatusNotFound, "Customer not found", "Customer not found")
	}

	var countData models.CountData
	var pagination int
	if limit > 0 {
		status, err = models.CountAdminGetAllCustomerStatusSuspend(&countData, params, paramsLike)
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
	response.Data = customers

	return c.JSON(http.StatusOK, response)
}

func AdminSuspendUnsuspendCustomer(c echo.Context) error {

	params := make(map[string]string)

	customerKeyStr := c.FormValue("customer_key")
	if customerKeyStr != "" {
		customerKey, err := strconv.ParseUint(customerKeyStr, 10, 64)
		if err != nil || customerKey == 0 {
			// log.Error("Wrong input for parameter: customer_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: customer_key", "Wrong input for parameter: customer_key")
		}
	} else {
		// log.Error("Missing required parameter: customer_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: customer_key", "Missing required parameter: customer_key")
	}

	suspendFlag := c.FormValue("suspend_flag")
	if suspendFlag == "" {
		// log.Error("Missing required parameter: suspend_flag")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: suspend_flag", "Missing required parameter: suspend_flag")
	}

	reason := c.FormValue("reason")
	if suspendFlag == "1" && reason == "" {
		// log.Error("Missing required parameter: reason")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: reason", "Missing required parameter: reason")
	}

	var cus models.MsCustomer
	_, err := models.GetMsCustomer(&cus, customerKeyStr)
	if err != nil {
		// log.Error(err.Error())
		return lib.CustomError(http.StatusBadRequest, err.Error(), "Customer tidak ditemukan")
	}

	dateLayout := "2006-01-02 15:04:05"
	params["customer_key"] = customerKeyStr
	params["cif_suspend_modified_date"] = time.Now().Format(dateLayout)
	params["cif_suspend_reason"] = reason
	params["cif_suspend_flag"] = suspendFlag
	params["rec_modified_date"] = time.Now().Format(dateLayout)
	params["rec_modified_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
	_, err = models.UpdateMsCustomer(params)
	if err != nil {
		// log.Error("Error update suspend ms_customer")
		return lib.CustomError(http.StatusBadRequest, err.Error(), "Failed update data")
	}

	//update tr_account
	paramsTrAccount := make(map[string]string)
	paramsTrAccount["sub_suspend_flag"] = suspendFlag
	paramsTrAccount["sub_suspend_modified_date"] = time.Now().Format(dateLayout)
	paramsTrAccount["sub_suspend_reason"] = reason
	paramsTrAccount["sub_suspend_reference"] = strconv.FormatUint(lib.Profile.UserID, 10)
	paramsTrAccount["red_suspend_flag"] = suspendFlag
	paramsTrAccount["red_suspend_modified_date"] = time.Now().Format(dateLayout)
	paramsTrAccount["red_suspend_reason"] = reason
	paramsTrAccount["red_suspend_reference"] = strconv.FormatUint(lib.Profile.UserID, 10)
	paramsTrAccount["rec_modified_date"] = time.Now().Format(dateLayout)
	paramsTrAccount["rec_modified_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
	_, err = models.UpdateTrAccountUploadSinvest(paramsTrAccount, "customer_key", customerKeyStr)
	if err != nil {
		// log.Println(err.Error())
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = nil
	return c.JSON(http.StatusOK, response)

}

func AdminGetDetailCustomer(c echo.Context) error {
	customerKeyStr := c.Param("customer_key")
	if customerKeyStr != "" {
		customerKey, err := strconv.ParseUint(customerKeyStr, 10, 64)
		if err != nil || customerKey == 0 {
			// log.Error("Wrong input for parameter: customer_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: customer_key", "Wrong input for parameter: customer_key")
		}
	} else {
		// log.Error("Missing required parameter: customer_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: customer_key", "Missing required parameter: customer_key")
	}

	params := make(map[string]string)
	params["c.customer_key"] = customerKeyStr
	params["c.investor_type"] = "263"

	var customer models.CustomerIndividuStatusSuspend
	_, err := models.AdminGetDetailCustomerStatusSuspend(&customer, params)
	if err != nil {
		// log.Error("Error get data ms_customer")
		return lib.CustomError(http.StatusBadRequest, err.Error(), "Failed get data")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = customer
	return c.JSON(http.StatusOK, response)

}

func UploadExcelCIFSuspendUnsuspend(c echo.Context) error {
	var err error
	err = os.MkdirAll(config.BasePath+"/cif_suspend/", 0755)
	if err != nil {
		// log.Error(err.Error())
		return lib.CustomError(http.StatusUnauthorized, "Terjadi kesalahan.", "Terjadi kesalahan.")
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
			var filename string
			filename = lib.RandStringBytesMaskImprSrc(20)
			// log.Println("Generate filename:", filename)
			// Upload image and move to proper directory
			err = lib.UploadImage(file, config.BasePath+"/cif_suspend/"+filename+extension)
			if err != nil {
				// log.Println(err)
				return lib.CustomError(http.StatusInternalServerError)
			}

			xlsx, err := excelize.OpenFile(config.BasePath + "/cif_suspend/" + filename + extension)
			if err != nil {
				// log.Println(config.BasePath + "/cif_suspend/" + filename + extension)
				// log.Fatal("ERROR ", err.Error())
				return lib.CustomError(http.StatusInternalServerError)
			}

			sheet1Name := xlsx.GetSheetName(1)

			isRowNoValue := 0

			for i := 2; i < 100000; i++ {
				no := xlsx.GetCellValue(sheet1Name, fmt.Sprintf("A%d", i))
				classification := xlsx.GetCellValue(sheet1Name, fmt.Sprintf("B%d", i))
				saCode := xlsx.GetCellValue(sheet1Name, fmt.Sprintf("C%d", i))
				saName := xlsx.GetCellValue(sheet1Name, fmt.Sprintf("D%d", i))
				sid := xlsx.GetCellValue(sheet1Name, fmt.Sprintf("E%d", i))
				ifua := xlsx.GetCellValue(sheet1Name, fmt.Sprintf("F%d", i))
				ifuaName := xlsx.GetCellValue(sheet1Name, fmt.Sprintf("G%d", i))
				suspendType := xlsx.GetCellValue(sheet1Name, fmt.Sprintf("H%d", i))
				suspendReason := xlsx.GetCellValue(sheet1Name, fmt.Sprintf("I%d", i))
				description := xlsx.GetCellValue(sheet1Name, fmt.Sprintf("J%d", i))
				effDate := xlsx.GetCellValue(sheet1Name, fmt.Sprintf("K%d", i))
				effDateInputTime := xlsx.GetCellValue(sheet1Name, fmt.Sprintf("L%d", i))
				closingDate := xlsx.GetCellValue(sheet1Name, fmt.Sprintf("M%d", i))
				closingDateInputTime := xlsx.GetCellValue(sheet1Name, fmt.Sprintf("N%d", i))

				// log.Println(classification)
				// log.Println(ifua)
				// log.Println(suspendType)
				// log.Println(suspendReason)

				if no == "" && classification == "" && saCode == "" && saName == "" &&
					sid == "" && ifua == "" && ifuaName == "" && suspendType == "" &&
					suspendReason == "" && description == "" && effDate == "" && effDateInputTime == "" &&
					closingDate == "" && closingDateInputTime == "" {
					isRowNoValue++
				}

				if isRowNoValue > 5 {
					break
				}

				if classification != "" && ifua != "" && suspendType != "" {
					var acc models.TrAccount
					_, err := models.GetTrAccountByField(&acc, ifua, "ifua_no")
					if err == nil {
						dateLayout := "2006-01-02 15:04:05"
						paramsTrAccount := make(map[string]string)
						paramsCustomer := make(map[string]string)
						if classification == "Suspension" {
							paramsCustomer["cif_suspend_flag"] = "1"
							if suspendType == "Subscription" {
								paramsTrAccount["sub_suspend_flag"] = "1"
								paramsTrAccount["sub_suspend_modified_date"] = time.Now().Format(dateLayout)
								paramsTrAccount["sub_suspend_reason"] = suspendReason
							} else {
								paramsTrAccount["red_suspend_flag"] = "1"
								paramsTrAccount["red_suspend_modified_date"] = time.Now().Format(dateLayout)
								paramsTrAccount["red_suspend_reason"] = suspendReason
							}
						} else {
							paramsCustomer["cif_suspend_flag"] = "0"
							if suspendType == "Subscription" {
								paramsTrAccount["sub_suspend_flag"] = "0"
								paramsTrAccount["sub_suspend_modified_date"] = time.Now().Format(dateLayout)
								paramsTrAccount["sub_suspend_reason"] = suspendReason
							} else {
								paramsTrAccount["red_suspend_flag"] = "0"
								paramsTrAccount["red_suspend_modified_date"] = time.Now().Format(dateLayout)
								paramsTrAccount["red_suspend_reason"] = suspendReason
							}
						}

						paramsTrAccount["ifua_no"] = ifua
						paramsTrAccount["ifua_name"] = ifuaName
						paramsTrAccount["rec_modified_date"] = time.Now().Format(dateLayout)
						paramsTrAccount["rec_modified_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
						_, err = models.UpdateTrAccountUploadSinvest(paramsTrAccount, "customer_key", strconv.FormatUint(acc.CustomerKey, 10))
						if err != nil {
							// log.Error("Error update suspend flag")
							// log.Error(err.Error())
							continue
						}

						//update customer
						paramsCustomer["cif_suspend_modified_date"] = time.Now().Format(dateLayout)
						paramsCustomer["cif_suspend_reason"] = suspendReason
						paramsCustomer["rec_modified_date"] = time.Now().Format(dateLayout)
						paramsCustomer["rec_modified_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
						_, err = models.UpdateMsCustomerByField(paramsCustomer, "customer_key", strconv.FormatUint(acc.CustomerKey, 10))
						if err != nil {
							// log.Error("Error update cms_customer")
							// log.Error(err.Error())
							continue
						}
					} else {
						// log.Error("Error get tr_account")
						// log.Error(err.Error())
						continue
					}
				}
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
	response.Data = ""
	return c.JSON(http.StatusOK, response)

}
