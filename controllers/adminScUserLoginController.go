package controllers

import (
	"bytes"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"html/template"
	"math"
	"mf-bo-api/config"
	"mf-bo-api/lib"
	"mf-bo-api/models"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/badoux/checkmail"
	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

func LogoutAdmin(c echo.Context) error {
	var err error

	strIDUserLogin := strconv.FormatUint(lib.Profile.UserID, 10)

	dateLayout := "2006-01-02 15:04:05"
	paramsSession := make(map[string]string)
	paramsSession["user_login_key"] = strIDUserLogin
	paramsSession["logout_date"] = time.Now().Format(dateLayout)
	paramsSession["login_session_key"] = ""

	_, err = models.UpdateScLoginSession(paramsSession)
	if err != nil {
		log.Error("Error update session in logout")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = ""
	return c.JSON(http.StatusOK, response)
}

func GetListScUserLoginAdmin(c echo.Context) error {

	// errorAuth := initAuthHoIt()
	// if errorAuth != nil {
	// 	log.Error("User Autorizer")
	// 	return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	// }

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

	items := []string{"user_login_key", "ucategory_name", "user_dept_name", "ulogin_name", "ulogin_full_name", "ulogin_email", "role_name", "rec_created_date"}

	params := make(map[string]string)
	orderBy := c.QueryParam("order_by")
	if orderBy != "" {
		_, found := lib.Find(items, orderBy)
		if found {
			var orderByJoin string
			orderByJoin = "u.user_login_key"
			if orderBy == "user_login_key" {
				orderByJoin = "u.user_login_key"
			}
			if orderBy == "ucategory_name" {
				orderByJoin = "cat.ucategory_name"
			}
			if orderBy == "user_dept_name" {
				orderByJoin = "dept.user_dept_name"
			}
			if orderBy == "ulogin_name" {
				orderByJoin = "u.ulogin_name"
			}
			if orderBy == "ulogin_full_name" {
				orderByJoin = "u.ulogin_full_name"
			}
			if orderBy == "ulogin_email" {
				orderByJoin = "u.ulogin_email"
			}
			if orderBy == "role_name" {
				orderByJoin = "role.role_name"
			}
			if orderBy == "rec_created_date" {
				orderByJoin = "u.rec_created_date"
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
		params["orderBy"] = "u.user_login_key"
		params["orderType"] = "ASC"
	}

	var searchData *string

	search := c.QueryParam("search")
	if search != "" {
		searchData = &search
	}

	userCategoryKey := c.QueryParam("user_category_key")
	if userCategoryKey != "" {
		params["u.user_category_key"] = userCategoryKey
	}

	rolekey := c.QueryParam("role_key")
	if rolekey != "" {
		rolekeyCek, err := strconv.ParseUint(rolekey, 10, 64)
		if err == nil && rolekeyCek > 0 {
			params["role.role_key"] = rolekey
		} else {
			log.Error("Wrong input for parameter: role_key")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: role_key", "Missing required parameter: role_key")
		}
	}
	//mapping scUserLogin
	var scUserLogin []models.AdminListScUserLogin
	status, err = models.AdminGetAllScUserLogin(&scUserLogin, limit, offset, params, noLimit, searchData)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}

	var countData models.CountData
	var pagination int
	if limit > 0 {
		status, err = models.AdminCountDataGetAllScUserlogin(&countData, params, searchData)
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
	response.Data = scUserLogin

	return c.JSON(http.StatusOK, response)
}

func GetDetailScUserLoginAdmin(c echo.Context) error {
	var err error
	var status int

	keyStr := c.Param("key")
	key, _ := strconv.ParseUint(keyStr, 10, 64)
	if key == 0 {
		return lib.CustomError(http.StatusNotFound)
	}

	var scUserLogin models.ScUserLogin
	status, err = models.GetScUserLoginByKey(&scUserLogin, keyStr)
	if err != nil {
		return lib.CustomError(http.StatusNotFound)
	}

	var responseData models.AdminDetailScUserLogin

	if scUserLogin.CustomerKey != nil {
		var customer models.MsCustomer
		status, err = models.GetMsCustomer(&customer, strconv.FormatUint(*scUserLogin.CustomerKey, 10))
		if err == nil {
			responseData.CustomerKey = &customer.CustomerKey
			responseData.CustomerName = &customer.FullName
		}

		var oa models.OaRequestKeyLastHistory
		status, err = models.AdminGetLastOaRequest(&oa, strconv.FormatUint(*scUserLogin.CustomerKey, 10))
		if err == nil {
			responseData.OaRequestKey = &oa.OaRequestKey
		}

		var ifua models.TrAccountIfua
		status, err = models.GetCifTrAccountByField(&ifua, strconv.FormatUint(*scUserLogin.CustomerKey, 10), "customer_key")
		if err == nil {
			responseData.Cif = ifua.IfuaNo
		}
	}

	responseData.UserLoginKey = scUserLogin.UserLoginKey
	responseData.NoHp = scUserLogin.UloginMobileno

	var scUserCategory models.ScUserCategory
	strUCKey := strconv.FormatUint(scUserLogin.UserCategoryKey, 10)
	status, err = models.GetScUserCategory(&scUserCategory, strUCKey)
	if err != nil {
		if err != sql.ErrNoRows {
			return lib.CustomError(status)
		}
	} else {
		var ucat models.ScUserCategoryInfo
		ucat.UserCategoryKey = scUserCategory.UserCategoryKey
		ucat.UcategoryCode = scUserCategory.UcategoryCode
		ucat.UcategoryName = scUserCategory.UcategoryName
		ucat.UcategoryDesc = scUserCategory.UcategoryDesc
		responseData.UserCategory = ucat
	}

	if scUserLogin.UserDeptKey != nil {
		var scUserDept models.ScUserDept
		strUDept := strconv.FormatUint(*scUserLogin.UserDeptKey, 10)
		status, err = models.GetScUserDept(&scUserDept, strUDept)
		if err != nil {
			if err != sql.ErrNoRows {
				return lib.CustomError(status)
			}
		} else {
			var udept models.ScUserDeptInfo
			udept.UserDeptKey = scUserDept.UserDeptKey
			udept.UserDeptCode = scUserDept.UserDeptCode
			udept.UserDeptName = scUserDept.UserDeptName
			udept.UserDeptDesc = scUserDept.UserDeptDesc
			responseData.UserDept = &udept
		}
	}

	responseData.UloginName = scUserLogin.UloginName
	responseData.UloginFullName = scUserLogin.UloginFullName
	responseData.UloginEmail = scUserLogin.UloginEmail

	if scUserLogin.RoleKey != nil {
		var scRole models.ScRole
		strRoleKey := strconv.FormatUint(*scUserLogin.RoleKey, 10)
		status, err = models.GetScRole(&scRole, strRoleKey)
		if err != nil {
			if err != sql.ErrNoRows {
				return lib.CustomError(status)
			}
		} else {
			var role models.ScRoleInfoLogin
			role.RoleKey = scRole.RoleKey
			role.RoleCode = scRole.RoleCode
			role.RoleName = scRole.RoleName
			role.RoleDesc = scRole.RoleDesc
			responseData.Role = &role
		}
	}

	if scUserLogin.UloginEnabled == uint8(1) {
		responseData.Enabled = true
	} else {
		responseData.Enabled = false
	}

	if scUserLogin.UloginLocked == uint8(1) {
		responseData.Locked = true
	} else {
		responseData.Locked = false
	}

	if scUserLogin.VerifiedEmail != nil {
		if *scUserLogin.VerifiedEmail == uint8(1) {
			responseData.VerifiedEmail = true
		} else {
			responseData.VerifiedEmail = false
		}
	} else {
		responseData.VerifiedEmail = false
	}

	if scUserLogin.VerifiedMobileno == uint8(1) {
		responseData.VerifiedMobileno = true
	} else {
		responseData.VerifiedMobileno = false
	}

	layout := "2006-01-02 15:04:05"
	newLayout := "02 Jan 2006"
	if scUserLogin.RecCreatedDate != nil {
		date, err := time.Parse(layout, *scUserLogin.RecCreatedDate)
		if err == nil {
			oke := date.Format(newLayout)
			responseData.CreatedDate = &oke
		}
	}

	if scUserLogin.RecImage1 != nil && *scUserLogin.RecImage1 != "" {
		responseData.RecImage = config.ImageUrl + "/images/user/" + strconv.FormatUint(scUserLogin.UserLoginKey, 10) + "/profile/" + *scUserLogin.RecImage1
	} else {
		responseData.RecImage = config.ImageUrl + "/images/post/default.png"
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = responseData

	return c.JSON(http.StatusOK, response)
}

func DisableEnableUser(c echo.Context) error {
	var err error

	params := make(map[string]string)

	key := c.FormValue("key")
	if key == "" {
		log.Error("Missing required parameter: key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: key", "Missing required parameter: key")
	}

	keyCek, err := strconv.ParseUint(key, 10, 64)
	if err == nil && keyCek > 0 {
		params["user_login_key"] = key
	} else {
		log.Error("Wrong input for parameter: key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: key", "Missing required parameter: key")
	}

	var scUserLogin models.ScUserLogin
	status, err := models.GetScUserLoginByKey(&scUserLogin, key)
	if err != nil {
		log.Error("User login not found")
		return lib.CustomError(http.StatusNotFound)
	}

	dateLayout := "2006-01-02 15:04:05"

	if scUserLogin.UloginEnabled == 1 {
		params["ulogin_enabled"] = "0"
	} else {
		params["ulogin_enabled"] = "1"
	}
	params["rec_modified_date"] = time.Now().Format(dateLayout)
	params["rec_modified_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

	status, err = models.UpdateScUserLogin(params)
	if err != nil {
		log.Error("Failed create request data: " + err.Error())
		return lib.CustomError(status, err.Error(), "failed input data")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = nil
	return c.JSON(http.StatusOK, response)

}

func LockUnlockUser(c echo.Context) error {
	var err error

	params := make(map[string]string)

	key := c.FormValue("key")
	if key == "" {
		log.Error("Missing required parameter: key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: key", "Missing required parameter: key")
	}

	keyCek, err := strconv.ParseUint(key, 10, 64)
	if err == nil && keyCek > 0 {
		params["user_login_key"] = key
	} else {
		log.Error("Wrong input for parameter: key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: key", "Missing required parameter: key")
	}

	var scUserLogin models.ScUserLogin
	status, err := models.GetScUserLoginByKey(&scUserLogin, key)
	if err != nil {
		log.Error("User login not found")
		return lib.CustomError(http.StatusNotFound)
	}

	dateLayout := "2006-01-02 15:04:05"

	if scUserLogin.UloginLocked == 1 {
		params["ulogin_locked"] = "0"
	} else {
		params["ulogin_locked"] = "1"
	}
	params["rec_modified_date"] = time.Now().Format(dateLayout)
	params["rec_modified_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

	status, err = models.UpdateScUserLogin(params)
	if err != nil {
		log.Error("Failed create request data: " + err.Error())
		return lib.CustomError(status, err.Error(), "failed input data")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = nil
	return c.JSON(http.StatusOK, response)

}

func CreateAdminScUserLogin(c echo.Context) error {
	var err error
	var status int

	params := make(map[string]string)

	//user_category_key
	usercategorykey := c.FormValue("user_category_key")
	if usercategorykey == "" {
		log.Error("Missing required parameter: user_category_key cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: user_category_key cann't be blank", "Missing required parameter: user_category_key cann't be blank")
	}
	sub, err := strconv.ParseUint(usercategorykey, 10, 64)
	if err == nil && sub > 0 {
		params["user_category_key"] = usercategorykey
	} else {
		log.Error("Wrong input for parameter: user_category_key number")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: user_category_key must number", "Missing required parameter: user_category_key number")
	}

	//user_dept_key
	userdeptkey := c.FormValue("user_dept_key")
	if userdeptkey == "" {
		log.Error("Missing required parameter: user_dept_key cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: user_dept_key cann't be blank", "Missing required parameter: user_dept_key cann't be blank")
	}
	sub, err = strconv.ParseUint(userdeptkey, 10, 64)
	if err == nil && sub > 0 {
		params["user_dept_key"] = userdeptkey
	} else {
		log.Error("Wrong input for parameter: user_dept_key number")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: user_dept_key must number", "Missing required parameter: user_dept_key number")
	}

	//role_key
	rolekey := c.FormValue("role_key")
	if rolekey == "" {
		log.Error("Missing required parameter: role_key cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: role_key cann't be blank", "Missing required parameter: role_key cann't be blank")
	}
	sub, err = strconv.ParseUint(rolekey, 10, 64)
	if err == nil && sub > 0 {
		params["role_key"] = rolekey
	} else {
		log.Error("Wrong input for parameter: role_key number")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: role_key must number", "Missing required parameter: role_key number")
	}

	//ulogin_name / ulogin_email
	email := c.FormValue("email")
	if email == "" {
		log.Error("Missing required parameter: email cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: email cann't be blank", "Missing required parameter: email cann't be blank")
	}
	//ulogin_name / ulogin_email
	uloginname := c.FormValue("username")
	if uloginname == "" {
		log.Error("Missing required parameter: username cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: username cann't be blank", "Missing required parameter: username cann't be blank")
	} else {
		if len(uloginname) < 5 {
			log.Error("Wrong input for parameter: username kurang dari 5 digit")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: username kurang dari 5 digit", "Missing required parameter: username kurang dari 5 digit")
		}
		if len(uloginname) > 50 {
			log.Error("Wrong input for parameter: username lebih dari 50 digit")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: username lebih dari 50 digit", "Missing required parameter: username lebih dari 50 digit")
		}
	}
	params["ulogin_name"] = uloginname
	params["ulogin_email"] = email

	//check unique ulogin_name / ulogin_email
	paramsScUserLogin := make(map[string]string)
	paramsScUserLogin["ulogin_name"] = uloginname
	paramsScUserLogin["ulogin_email"] = email

	var countDataExisting models.CountData
	status, err = models.AdminGetValidateUniqueInsertUpdateScUserLogin(&countDataExisting, paramsScUserLogin, nil)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if int(countDataExisting.CountData) > 0 {
		log.Error("Missing required parameter: username or email already existing, use other username or email")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: username or email already existing, use other username or email", "Missing required parameter: username or email already existing, use other username or email")
	}

	// Validate email
	err = checkmail.ValidateFormat(email)
	if err != nil {
		log.Error("email format is not valid, must email")
		return lib.CustomError(http.StatusBadRequest, "email format is not valid, email must email format", "Email format is not valid, Email must email format")
	}

	if strings.Contains(uloginname, " ") {
		log.Error("Username format is not valid")
		return lib.CustomError(http.StatusBadRequest, "Username format is not valid", "Username format is not valid")
	}

	err = checkmail.ValidateFormat(email)
	if err != nil {
		log.Error("Username format is not valid, must email")
		return lib.CustomError(http.StatusBadRequest, "Username format is not valid, Username must email format", "Username format is not valid, Username must email format")
	}

	//ulogin_full_name
	uloginfullname := c.FormValue("ulogin_full_name")
	if uloginfullname == "" {
		log.Error("Missing required parameter: ulogin_full_name cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: ulogin_full_name cann't be blank", "Missing required parameter: ulogin_full_name cann't be blank")
	}
	params["ulogin_full_name"] = uloginfullname

	password := c.FormValue("password")
	if password == "" {
		log.Error("Missing required parameter: password cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: password cann't be blank", "Missing required parameter: password cann't be blank")
	} else {
		if len(password) < 8 {
			log.Error("Wrong input for parameter: password kurang dari 8 digit")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: password kurang dari 8 digit", "Missing required parameter: password kurang dari 8 digit")
		}
	}

	repassword := c.FormValue("re_password")
	if repassword == "" {
		log.Error("Missing required parameter: re_password cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: re_password cann't be blank", "Missing required parameter: re_password cann't be blank")
	}

	if password != repassword {
		log.Error("Missing required parameter: password & re_password harus sama")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: password & re_password harus sama", "Missing required parameter: password & re_password harus sama")
	}

	// Validate password
	length, number, upper, special := verifyPassword(password)
	if length == false || number == false || upper == false || special == false {
		log.Error("Password does meet the criteria")
		return lib.CustomError(http.StatusBadRequest, "Password does meet the criteria", "Your password need at least 8 character length, has lower and upper case letter, has numeric letter, and has special character")
	}
	// Encrypt password
	encryptedPasswordByte := sha256.Sum256([]byte(password))
	encryptedPassword := hex.EncodeToString(encryptedPasswordByte[:])
	params["ulogin_password"] = encryptedPassword

	//phone_number
	nohp := c.FormValue("no_hp")
	if nohp != "" {
		params["ulogin_mobileno"] = nohp
	}

	//ulogin_enabled
	enabled := c.FormValue("enabled")
	var enabledBool bool
	if enabled != "" {
		enabledBool, err = strconv.ParseBool(enabled)
		if err != nil {
			log.Error("enabled parameter should be true/false")
			return lib.CustomError(http.StatusBadRequest, "enabled parameter should be true/false", "enabled parameter should be true/false")
		}
		if enabledBool == true {
			params["ulogin_enabled"] = "1"
		} else {
			params["ulogin_enabled"] = "0"
		}
	} else {
		log.Error("enabled parameter should be true/false")
		return lib.CustomError(http.StatusBadRequest, "enabled parameter should be true/false", "enabled parameter should be true/false")
	}

	//locked
	locked := c.FormValue("locked")
	var lockedBool bool
	if locked != "" {
		lockedBool, err = strconv.ParseBool(locked)
		if err != nil {
			log.Error("locked parameter should be true/false")
			return lib.CustomError(http.StatusBadRequest, "locked parameter should be true/false", "locked parameter should be true/false")
		}
		if lockedBool == true {
			params["ulogin_locked"] = "1"
		} else {
			params["ulogin_locked"] = "0"
		}
	} else {
		log.Error("locked parameter should be true/false")
		return lib.CustomError(http.StatusBadRequest, "locked parameter should be true/false", "locked parameter should be true/false")
	}

	// Set expired for token
	date := time.Now().AddDate(0, 0, 1)
	dateLayout := "2006-01-02 15:04:05"
	expired := date.Format(dateLayout)

	// Generate verify key
	verifyKeyByte := sha256.Sum256([]byte(uloginname + "_" + expired))
	verifyKey := hex.EncodeToString(verifyKeyByte[:])

	// Input to database
	params["ulogin_must_changepwd"] = "0"
	params["last_password_changed"] = time.Now().Format(dateLayout)
	params["verified_email"] = "1"
	params["verified_mobileno"] = "1"
	params["ulogin_enabled"] = "1"
	params["ulogin_failed_count"] = "0"
	params["last_access"] = time.Now().Format(dateLayout)
	params["accept_login_tnc"] = "1"
	params["allowed_sharing_login"] = "1"
	params["string_token"] = verifyKey
	params["token_expired"] = expired

	params["rec_order"] = "0"
	params["rec_status"] = "1"
	params["rec_created_date"] = time.Now().Format(dateLayout)
	params["rec_created_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

	status, err = models.CreateScUserLogin(params)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(http.StatusBadRequest, err.Error(), "Failed create user")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = nil
	return c.JSON(http.StatusOK, response)

}

func UpdateAdminScUserLogin(c echo.Context) error {
	var err error
	params := make(map[string]string)

	userloginkey := c.FormValue("key")
	if userloginkey == "" {
		log.Error("Missing required parameter: key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: key", "Missing required parameter: key")
	}
	struserloginkey, err := strconv.ParseUint(userloginkey, 10, 64)
	if err == nil && struserloginkey > 0 {
		params["user_login_key"] = userloginkey
	} else {
		log.Error("Wrong input for parameter: key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: key", "Missing required parameter: key")
	}

	var scUserLogin models.ScUserLogin
	_, err = models.GetScUserLoginByKey(&scUserLogin, userloginkey)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(http.StatusBadRequest)
	}

	//user_category_key
	usercategorykey := c.FormValue("user_category_key")
	if usercategorykey == "" {
		log.Error("Missing required parameter: user_category_key cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: user_category_key cann't be blank", "Missing required parameter: user_category_key cann't be blank")
	}
	sub, err := strconv.ParseUint(usercategorykey, 10, 64)
	if err == nil && sub > 0 {
		params["user_category_key"] = usercategorykey
	} else {
		log.Error("Wrong input for parameter: user_category_key number")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: user_category_key must number", "Missing required parameter: user_category_key number")
	}

	//user_dept_key
	userdeptkey := c.FormValue("user_dept_key")
	if userdeptkey == "" {
		log.Error("Missing required parameter: user_dept_key cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: user_dept_key cann't be blank", "Missing required parameter: user_dept_key cann't be blank")
	}
	sub, err = strconv.ParseUint(userdeptkey, 10, 64)
	if err == nil && sub > 0 {
		params["user_dept_key"] = userdeptkey
	} else {
		log.Error("Wrong input for parameter: user_dept_key number")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: user_dept_key must number", "Missing required parameter: user_dept_key number")
	}

	//role_key
	rolekey := c.FormValue("role_key")
	if rolekey == "" {
		log.Error("Missing required parameter: role_key cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: role_key cann't be blank", "Missing required parameter: role_key cann't be blank")
	}
	sub, err = strconv.ParseUint(rolekey, 10, 64)
	if err == nil && sub > 0 {
		params["role_key"] = rolekey
	} else {
		log.Error("Wrong input for parameter: role_key number")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: role_key must number", "Missing required parameter: role_key number")
	}

	//ulogin_full_name
	uloginfullname := c.FormValue("ulogin_full_name")
	if uloginfullname == "" {
		log.Error("Missing required parameter: ulogin_full_name cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: ulogin_full_name cann't be blank", "Missing required parameter: ulogin_full_name cann't be blank")
	}
	params["ulogin_full_name"] = uloginfullname

	//phone_number
	nohp := c.FormValue("no_hp")
	if nohp != "" {
		params["ulogin_mobileno"] = nohp
	}

	//ulogin_enabled
	recstatus := c.FormValue("enabled")
	var recstatusBool bool
	if recstatus != "" {
		recstatusBool, err = strconv.ParseBool(recstatus)
		if err != nil {
			log.Error("enabled parameter should be true/false")
			return lib.CustomError(http.StatusBadRequest, "enabled parameter should be true/false", "enabled parameter should be true/false")
		}
		if recstatusBool == true {
			params["ulogin_enabled"] = "1"
		} else {
			params["ulogin_enabled"] = "0"
		}
	} else {
		log.Error("enabled parameter should be true/false")
		return lib.CustomError(http.StatusBadRequest, "enabled parameter should be true/false", "enabled parameter should be true/false")
	}

	//ulogin_locked
	uloginlocked := c.FormValue("locked")
	var uloginlockedBool bool
	if uloginlocked != "" {
		uloginlockedBool, err = strconv.ParseBool(uloginlocked)
		if err != nil {
			log.Error("locked parameter should be true/false")
			return lib.CustomError(http.StatusBadRequest, "locked parameter should be true/false", "locked parameter should be true/false")
		}
		if uloginlockedBool == true {
			params["ulogin_locked"] = "1"
		} else {
			params["ulogin_locked"] = "0"
		}
	} else {
		log.Error("locked parameter should be true/false")
		return lib.CustomError(http.StatusBadRequest, "locked parameter should be true/false", "locked parameter should be true/false")
	}

	//verified_email
	verifiedemail := c.FormValue("verified_email")
	var verifiedemailBool bool
	if verifiedemail != "" {
		verifiedemailBool, err = strconv.ParseBool(verifiedemail)
		if err != nil {
			log.Error("verified_email parameter should be true/false")
			return lib.CustomError(http.StatusBadRequest, "verified_email parameter should be true/false", "verified_email parameter should be true/false")
		}
		if verifiedemailBool == true {
			params["verified_email"] = "1"
		} else {
			params["verified_email"] = "0"
		}
	} else {
		log.Error("verified_email parameter should be true/false")
		return lib.CustomError(http.StatusBadRequest, "verified_email parameter should be true/false", "verified_email parameter should be true/false")
	}

	//verified_mobileno
	verifiedmobileno := c.FormValue("verified_mobileno")
	var verifiedmobilenoBool bool
	if verifiedmobileno != "" {
		verifiedmobilenoBool, err = strconv.ParseBool(verifiedmobileno)
		if err != nil {
			log.Error("verified_mobileno parameter should be true/false")
			return lib.CustomError(http.StatusBadRequest, "verified_mobileno parameter should be true/false", "verified_mobileno parameter should be true/false")
		}
		if verifiedmobilenoBool == true {
			params["verified_mobileno"] = "1"
		} else {
			params["verified_mobileno"] = "0"
		}
	} else {
		log.Error("verified_mobileno parameter should be true/false")
		return lib.CustomError(http.StatusBadRequest, "verified_mobileno parameter should be true/false", "verified_mobileno parameter should be true/false")
	}

	dateLayout := "2006-01-02 15:04:05"
	params["rec_order"] = "0"
	params["rec_modified_date"] = time.Now().Format(dateLayout)
	params["rec_modified_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

	_, err = models.UpdateScUserLogin(params)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(http.StatusBadRequest, err.Error(), "Failed create user")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = nil
	return c.JSON(http.StatusOK, response)

}

func ChangePasswordUser(c echo.Context) error {
	var err error
	params := make(map[string]string)

	userloginkey := c.FormValue("key")
	if userloginkey == "" {
		log.Error("Missing required parameter: key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: key", "Missing required parameter: key")
	}
	struserloginkey, err := strconv.ParseUint(userloginkey, 10, 64)
	if err == nil && struserloginkey > 0 {
		params["user_login_key"] = userloginkey
	} else {
		log.Error("Wrong input for parameter: key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: key", "Missing required parameter: key")
	}

	var scUserLogin models.ScUserLogin
	_, err = models.GetScUserLoginByKey(&scUserLogin, userloginkey)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(http.StatusBadRequest)
	}

	password := lib.RandStringBytesMaskImprSrc(8)

	// Encrypt password
	encryptedPasswordByte := sha256.Sum256([]byte(password))
	encryptedPassword := hex.EncodeToString(encryptedPasswordByte[:])
	params["ulogin_password"] = encryptedPassword
	params["rec_attribute_id1"] = password
	params["ulogin_must_changepwd"] = "1"

	dateLayout := "2006-01-02 15:04:05"
	params["rec_order"] = "0"
	params["rec_modified_date"] = time.Now().Format(dateLayout)
	params["rec_modified_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

	_, err = models.UpdateScUserLogin(params)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(http.StatusBadRequest, err.Error(), "Failed create user")
	}

	//send email to user
	bodyEmail := "<p>Password Baru anda : <b>" + password + "<b><p/> <br/> <p>Login dengan password diatas dan ganti password anda dengan password baru.<p/>"
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", config.EmailFrom)
	mailer.SetHeader("To", scUserLogin.UloginEmail)
	mailer.SetHeader("Subject", "[MotionFunds] Change Password")
	mailer.SetBody("text/html", bodyEmail)

	err = lib.SendEmail(mailer)
	if err != nil {
		log.Error("Error send email")
		log.Error(err)
	}

	// e := email.NewEmail()
	// e.From = config.EmailFrom
	// e.To = []string{scUserLogin.UloginEmail}
	// e.Subject = "[MotionFunds] Change Password"
	// e.HTML = []byte(bodyEmail)

	// byteSlice, _ := e.Bytes()

	// privateKey, err := ioutil.ReadFile(config.BasePath + "/pkey-dkim_motionfunds.txt")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// options := dkim.NewSigOptions()
	// options.PrivateKey = privateKey
	// options.Domain = "motionfunds.id"
	// options.Selector = "dkim"
	// options.SignatureExpireIn = 3600
	// options.BodyLength = uint(len([]byte(bodyEmail)))
	// options.Headers = []string{"from", "date", "mime-version", "received", "received"}
	// options.AddSignatureTimestamp = true
	// options.Canonicalization = "relaxed/relaxed"

	// err = dkim.Sign(&byteSlice, options)
	// if err != nil {
	// 	log.Error("Error send email no password")
	// 	log.Error(err)
	// 	return err
	// }
	// tlsconfig := &tls.Config{
	// 	InsecureSkipVerify: true,
	// }
	// add := config.EmailSMTPHost + ":" + strconv.FormatUint(config.EmailSMTPPort, 10)
	// // err = e.Send(add, smtp.PlainAuth("", config.EmailFrom, config.EmailFromPassword, config.EmailSMTPHost))
	// err = e.SendWithStartTLS(add, smtp.PlainAuth("", config.EmailFrom, config.EmailFromPassword, config.EmailSMTPHost), tlsconfig)
	// if err != nil {
	// 	log.Error("Error send email dkim")
	// 	log.Error(err)
	// }

	// d := &gomail.Dialer{Host: config.EmailSMTPHost, Port: int(config.EmailSMTPPort)}
	// d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	// if err := d.DialAndSend(mailer); err != nil {
	// 	log.Error("Error send email no password")
	// 	log.Error(err)
	// }

	//end send email to user

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = nil
	return c.JSON(http.StatusOK, response)

}

func ChangeRoleUser(c echo.Context) error {
	var err error

	errorAuth := initAuthHoIt()
	if errorAuth != nil {
		log.Error("User Autorizer")
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}

	params := make(map[string]string)

	userloginkey := c.FormValue("key")
	if userloginkey == "" {
		log.Error("Missing required parameter: key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: key", "Missing required parameter: key")
	}
	struserloginkey, err := strconv.ParseUint(userloginkey, 10, 64)
	if err == nil && struserloginkey > 0 {
		params["user_login_key"] = userloginkey
	} else {
		log.Error("Wrong input for parameter: key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: key", "Missing required parameter: key")
	}

	var scUserLogin models.ScUserLogin
	_, err = models.GetScUserLoginByKey(&scUserLogin, userloginkey)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(http.StatusBadRequest)
	}

	//role_key
	rolekey := c.FormValue("role_key")
	if rolekey == "" {
		log.Error("Missing required parameter: role_key cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: role_key cann't be blank", "Missing required parameter: role_key cann't be blank")
	}
	sub, err := strconv.ParseUint(rolekey, 10, 64)
	if err == nil && sub > 0 {
		params["role_key"] = rolekey
	} else {
		log.Error("Wrong input for parameter: role_key number")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: role_key must number", "Missing required parameter: role_key number")
	}

	dateLayout := "2006-01-02 15:04:05"
	params["rec_order"] = "0"
	params["rec_modified_date"] = time.Now().Format(dateLayout)
	params["rec_modified_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

	_, err = models.UpdateScUserLogin(params)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(http.StatusBadRequest, err.Error(), "Failed create user")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = nil
	return c.JSON(http.StatusOK, response)

}

func DeleteUser(c echo.Context) error {
	var err error

	errorAuth := initAuthHoIt()
	if errorAuth != nil {
		log.Error("User Autorizer")
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}

	params := make(map[string]string)

	key := c.FormValue("key")
	if key == "" {
		log.Error("Missing required parameter: key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: key", "Missing required parameter: key")
	}

	keyCek, err := strconv.ParseUint(key, 10, 64)
	if err == nil && keyCek > 0 {
		params["user_login_key"] = key
	} else {
		log.Error("Wrong input for parameter: key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: key", "Missing required parameter: key")
	}

	var scUserLogin models.ScUserLogin
	status, err := models.GetScUserLoginByKey(&scUserLogin, key)
	if err != nil {
		log.Error("User login not found")
		return lib.CustomError(http.StatusNotFound)
	}

	dateLayout := "2006-01-02 15:04:05"

	params["rec_status"] = "0"
	params["rec_deleted_date"] = time.Now().Format(dateLayout)
	params["rec_deleted_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

	status, err = models.UpdateScUserLogin(params)
	if err != nil {
		log.Error("Failed create request data: " + err.Error())
		return lib.CustomError(status, err.Error(), "failed input data")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = nil
	return c.JSON(http.StatusOK, response)

}

func GetDetailScUserLogin(c echo.Context) error {
	var err error
	var status int

	strKey := strconv.FormatUint(lib.Profile.UserID, 10)

	var scUserLogin models.ScUserLogin
	status, err = models.GetScUserLoginByKey(&scUserLogin, strKey)
	if err != nil {
		return lib.CustomError(http.StatusNotFound)
	}

	var responseData models.AdminDetailScUserLogin

	responseData.UserLoginKey = scUserLogin.UserLoginKey
	responseData.NoHp = scUserLogin.UloginMobileno

	var scUserCategory models.ScUserCategory
	strUCKey := strconv.FormatUint(scUserLogin.UserCategoryKey, 10)
	status, err = models.GetScUserCategory(&scUserCategory, strUCKey)
	if err != nil {
		if err != sql.ErrNoRows {
			return lib.CustomError(status)
		}
	} else {
		var ucat models.ScUserCategoryInfo
		ucat.UserCategoryKey = scUserCategory.UserCategoryKey
		ucat.UcategoryCode = scUserCategory.UcategoryCode
		ucat.UcategoryName = scUserCategory.UcategoryName
		ucat.UcategoryDesc = scUserCategory.UcategoryDesc
		responseData.UserCategory = ucat
	}

	if scUserLogin.UserDeptKey != nil {
		var scUserDept models.ScUserDept
		strUDept := strconv.FormatUint(*scUserLogin.UserDeptKey, 10)
		status, err = models.GetScUserDept(&scUserDept, strUDept)
		if err != nil {
			if err != sql.ErrNoRows {
				return lib.CustomError(status)
			}
		} else {
			var udept models.ScUserDeptInfo
			udept.UserDeptKey = scUserDept.UserDeptKey
			udept.UserDeptCode = scUserDept.UserDeptCode
			udept.UserDeptName = scUserDept.UserDeptName
			udept.UserDeptDesc = scUserDept.UserDeptDesc
			responseData.UserDept = &udept
		}
	}

	responseData.UloginName = scUserLogin.UloginName
	responseData.UloginFullName = scUserLogin.UloginFullName
	responseData.UloginEmail = scUserLogin.UloginEmail

	if scUserLogin.RoleKey != nil {
		var scRole models.ScRole
		strRoleKey := strconv.FormatUint(*scUserLogin.RoleKey, 10)
		status, err = models.GetScRole(&scRole, strRoleKey)
		if err != nil {
			if err != sql.ErrNoRows {
				return lib.CustomError(status)
			}
		} else {
			var role models.ScRoleInfoLogin
			role.RoleKey = scRole.RoleKey
			role.RoleCode = scRole.RoleCode
			role.RoleName = scRole.RoleName
			role.RoleDesc = scRole.RoleDesc
			responseData.Role = &role
		}
	}

	if scUserLogin.UloginEnabled == uint8(1) {
		responseData.Enabled = true
	} else {
		responseData.Enabled = false
	}

	if scUserLogin.UloginLocked == uint8(1) {
		responseData.Locked = true
	} else {
		responseData.Locked = false
	}

	if scUserLogin.VerifiedEmail != nil {
		if *scUserLogin.VerifiedEmail == uint8(1) {
			responseData.VerifiedEmail = true
		} else {
			responseData.VerifiedEmail = false
		}
	} else {
		responseData.VerifiedEmail = false
	}

	if scUserLogin.VerifiedMobileno == uint8(1) {
		responseData.VerifiedMobileno = true
	} else {
		responseData.VerifiedMobileno = false
	}

	layout := "2006-01-02 15:04:05"
	newLayout := "02 Jan 2006"
	if scUserLogin.RecCreatedDate != nil {
		date, err := time.Parse(layout, *scUserLogin.RecCreatedDate)
		if err == nil {
			oke := date.Format(newLayout)
			responseData.CreatedDate = &oke
		}
	}

	if scUserLogin.RecImage1 != nil && *scUserLogin.RecImage1 != "" {
		responseData.RecImage = config.ImageUrl + "/images/user/" + strconv.FormatUint(scUserLogin.UserLoginKey, 10) + "/profile/" + *scUserLogin.RecImage1
	} else {
		responseData.RecImage = config.ImageUrl + "/images/post/default.png"
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = responseData

	return c.JSON(http.StatusOK, response)
}

func AdminChangePasswordUserLogin(c echo.Context) error {
	var err error
	params := make(map[string]string)

	currentPassword := c.FormValue("current_password")
	if currentPassword == "" {
		log.Error("Missing required parameter: current_password")
		return lib.CustomError(http.StatusBadRequest, "current password can not be blank", "current password can not be blank")
	}

	newPassword := c.FormValue("new_password")
	if newPassword == "" {
		log.Error("Missing required parameter: new_password")
		return lib.CustomError(http.StatusBadRequest, "new password can not be blank", "new password can not be blank")
	}

	confirmPassword := c.FormValue("confirm_password")
	if confirmPassword == "" {
		log.Error("Missing required parameter: confirm_password")
		return lib.CustomError(http.StatusBadRequest, "confirm new password can not be blank", "confirm new password can not be blank")
	}

	strKey := strconv.FormatUint(lib.Profile.UserID, 10)

	var scUserLogin models.ScUserLogin
	_, err = models.GetScUserLoginByKey(&scUserLogin, strKey)
	if err != nil {
		return lib.CustomError(http.StatusNotFound)
	}
	params["user_login_key"] = strKey

	length, number, upper, special := verifyPassword(newPassword)
	if length == false || number == false || upper == false || special == false {
		log.Error("New Password does meet the criteria")
		return lib.CustomError(http.StatusBadRequest, "New Password does meet the criteria", "Your New password need at least 8 character length, has lower and upper case letter, has numeric letter, and has special character")
	}

	if strings.Contains(newPassword, confirmPassword) == false {
		log.Error("Missing required parameter: conf_password must equal with password")
		return lib.CustomError(http.StatusBadRequest, "conf_password must equal with password", "conf_password must equal with password")
	}
	// Check valid password
	encryptedPasswordByte := sha256.Sum256([]byte(currentPassword))
	encryptedPassword := hex.EncodeToString(encryptedPasswordByte[:])
	if encryptedPassword != scUserLogin.UloginPassword {
		log.Error("Missing required parameter: wrong current password")
		return lib.CustomError(http.StatusBadRequest, "wrong current password", "wrong current password")
	}

	encryptedNewPasswordByte := sha256.Sum256([]byte(newPassword))
	encryptedNewPassword := hex.EncodeToString(encryptedNewPasswordByte[:])

	params["ulogin_password"] = encryptedNewPassword
	params["ulogin_must_changepwd"] = "0"

	dateLayout := "2006-01-02 15:04:05"
	params["rec_modified_date"] = time.Now().Format(dateLayout)
	params["rec_modified_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

	_, err = models.UpdateScUserLogin(params)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(http.StatusBadRequest, err.Error(), "Failed update password")
	}

	if config.Envi == "PROD" {
		//send email to user
		mailer := gomail.NewMessage()
		mailer.SetHeader("From", config.EmailFrom)
		mailer.SetHeader("To", scUserLogin.UloginEmail)
		mailer.SetHeader("Subject", "[MotionFunds] Change Password")
		mailer.SetBody("text/html", "<p>Password berhasil diubah.<p/><p>Apabila kamu tidak merasa mengganti password dengan password baru, segera hubungi admin MotionFunds.<p/>")

		err = lib.SendEmail(mailer)
		if err != nil {
			log.Error("Error send email")
			log.Error(err)
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
		// }
		//end send email to user
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = nil
	return c.JSON(http.StatusOK, response)
}

func AdminChangeDataUserLogin(c echo.Context) error {
	var err error
	params := make(map[string]string)

	username := c.FormValue("username")
	if username == "" {
		log.Error("Missing required parameter: username")
		return lib.CustomError(http.StatusBadRequest, "username can not be blank", "username can not be blank")
	}
	params["ulogin_name"] = username

	name := c.FormValue("name")
	if name == "" {
		log.Error("Missing required parameter: name")
		return lib.CustomError(http.StatusBadRequest, "name can not be blank", "name can not be blank")
	}
	params["ulogin_full_name"] = name

	email := c.FormValue("email")
	if email == "" {
		log.Error("Missing required parameter: email")
		return lib.CustomError(http.StatusBadRequest, "email can not be blank", "email can not be blank")
	}

	err = checkmail.ValidateFormat(email)
	if err != nil {
		log.Error("Email format is not valid")
		return lib.CustomError(http.StatusBadRequest, "Email format is not valid", "Email format is not valid")
	}
	params["ulogin_email"] = email

	no_hp := c.FormValue("no_hp")
	if no_hp == "" {
		log.Error("Missing required parameter: no_hp")
		return lib.CustomError(http.StatusBadRequest, "no_hp can not be blank", "no_hp can not be blank")
	}
	params["ulogin_mobileno"] = no_hp

	strKey := strconv.FormatUint(lib.Profile.UserID, 10)

	//check unique ulogin_name
	paramsScUserLogin := make(map[string]string)
	paramsScUserLogin["ulogin_name"] = username
	paramsScUserLogin["rec_status"] = "1"

	var countDataExisting models.CountData
	status, err := models.AdminGetValidateUniqueInsertUpdateScUserLogin(&countDataExisting, paramsScUserLogin, &strKey)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if int(countDataExisting.CountData) > 0 {
		log.Error("Missing required parameter: username or email already existing, use other username or email")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: username or email already existing, use other username or email", "Missing required parameter: username or email already existing, use other username or email")
	}

	//check unique ulogin_email
	paramsEmail := make(map[string]string)
	paramsEmail["ulogin_email"] = email

	var countDataEmail models.CountData
	status, err = models.AdminGetValidateUniqueInsertUpdateScUserLogin(&countDataEmail, paramsEmail, &strKey)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if int(countDataEmail.CountData) > 0 {
		log.Error("Missing required parameter: email already existing, use other email")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: email already existing, use other email", "Missing required parameter: email already existing, use other email")
	}

	var scUserLogin models.ScUserLogin
	_, err = models.GetScUserLoginByKey(&scUserLogin, strKey)
	if err != nil {
		return lib.CustomError(http.StatusNotFound)
	}
	params["user_login_key"] = strKey

	dateLayout := "2006-01-02 15:04:05"
	params["rec_modified_date"] = time.Now().Format(dateLayout)
	params["rec_modified_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

	_, err = models.UpdateScUserLogin(params)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(http.StatusBadRequest, err.Error(), "Failed create user")
	}

	if config.Envi == "PROD" {
		//send email to user
		mailer := gomail.NewMessage()
		mailer.SetHeader("From", config.EmailFrom)
		mailer.SetHeader("To", email)
		mailer.SetHeader("Subject", "[MotionFunds] Change Password")
		mailer.SetBody("text/html", "<p>Data berhasil diubah.<p/><p>Apabila kamu tidak merasa mengganti data dengan data baru, segera hubungi admin MotionFunds.<p/>")

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
		//end send email to user
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = nil
	return c.JSON(http.StatusOK, response)
}

func UpdateAuthScUserLogin(c echo.Context) error {
	var err error
	params := make(map[string]string)

	userloginkey := c.FormValue("key")
	if userloginkey == "" {
		log.Error("Missing required parameter: key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: key", "Missing required parameter: key")
	}

	if len(userloginkey) > 11 {
		log.Error("key must maximal 11 character")
		return lib.CustomError(http.StatusBadRequest, "key must maximal 11 character", "key must maximal 11 character")
	}
	struserloginkey, err := strconv.ParseUint(userloginkey, 10, 64)
	if err == nil && struserloginkey > 0 {
		params["user_login_key"] = userloginkey
	} else {
		log.Error("Wrong input for parameter: key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: key", "Missing required parameter: key")
	}

	var scUserLogin models.ScUserLogin
	_, err = models.GetScUserLoginByKey(&scUserLogin, userloginkey)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(http.StatusBadRequest)
	}
	ulogin_enable_prev_state := scUserLogin.UloginEnabled

	//username
	username := c.FormValue("username")
	if username == "" {
		log.Error("Missing required parameter: username cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: username cann't be blank", "Missing required parameter: username cann't be blank")
	}
	if len(username) > 50 {
		log.Error("username must maximal 50 character")
		return lib.CustomError(http.StatusBadRequest, "username must maximal 50 character", "username must maximal 50 character")
	}
	if strings.Contains(username, " ") {
		log.Error("username tidak boleh terdapat spasi")
		return lib.CustomError(http.StatusBadRequest, "username tidak boleh terdapat spasi", "username tidak boleh terdapat spasi")
	}

	//check unique ulogin_name
	var countDataUsername models.CountData
	status, err := models.ValidateUniqueData(&countDataUsername, "ulogin_name", username, &userloginkey)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if int(countDataUsername.CountData) > 0 {
		log.Error("Missing required parameter: username already existing, use other username")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: username already existing, use other username", "Missing required parameter: username already existing, use other username")
	}
	params["ulogin_name"] = username

	//email
	email := c.FormValue("email")
	if email == "" {
		log.Error("Missing required parameter: email cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: email cann't be blank", "Missing required parameter: email cann't be blank")
	}
	if len(email) > 50 {
		log.Error("email must maximal 50 character")
		return lib.CustomError(http.StatusBadRequest, "email must maximal 50 character", "email must maximal 50 character")
	}
	err = checkmail.ValidateFormat(email)
	if err != nil {
		log.Error("Email format is not valid")
		return lib.CustomError(http.StatusBadRequest, "Email format is not valid", "Email format is not valid")
	}

	var countDataEmail models.CountData
	status, err = models.ValidateUniqueData(&countDataEmail, "ulogin_email", email, &userloginkey)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if int(countDataEmail.CountData) > 0 {
		log.Error("Missing required parameter: email already existing, use other email")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: email already existing, use other email", "Missing required parameter: email already existing, use other email")
	}
	params["ulogin_email"] = email

	//no_hp
	noHp := c.FormValue("no_hp")
	if noHp == "" {
		log.Error("Missing required parameter: no_hp cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: no_hp cann't be blank", "Missing required parameter: no_hp cann't be blank")
	}
	if len(noHp) > 18 {
		log.Error("no_hp must maximal 18 character")
		return lib.CustomError(http.StatusBadRequest, "no_hp must maximal 18 character", "no_hp must maximal 18 character")
	}

	var countDataNoHp models.CountData
	status, err = models.ValidateUniqueData(&countDataNoHp, "ulogin_mobileno", noHp, &userloginkey)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if int(countDataNoHp.CountData) > 0 {
		log.Error("Missing required parameter: no_hp already existing, use other no_hp")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: no_hp already existing, use other no_hp", "Missing required parameter: no_hp already existing, use other no_hp")
	}
	params["ulogin_mobileno"] = noHp

	//ulogin_full_name
	uloginFullName := c.FormValue("ulogin_full_name")
	if uloginFullName == "" {
		log.Error("Missing required parameter: ulogin_full_name cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: ulogin_full_name cann't be blank", "Missing required parameter: ulogin_full_name cann't be blank")
	}
	if len(uloginFullName) > 50 {
		log.Error("ulogin_full_name must maximal 50 character")
		return lib.CustomError(http.StatusBadRequest, "ulogin_full_name must maximal 50 character", "ulogin_full_name must maximal 50 character")
	}
	params["ulogin_full_name"] = uloginFullName

	//ulogin_enabled
	statusEnable := c.FormValue("status_enabled")
	if statusEnable == "" {
		log.Error("Missing required parameter: status_enabled cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: status_enabled cann't be blank", "Missing required parameter: status_enabled cann't be blank")
	}
	if len(statusEnable) > 2 {
		log.Error("status_enabled must maximal 2 character")
		return lib.CustomError(http.StatusBadRequest, "status_enabled must maximal 2 character", "status_enabled must maximal 2 character")
	}

	items := []string{"0", "1", "9"}
	//0 = Disabled
	//1 = Enabled/Online
	//9 = Offline

	_, found := lib.Find(items, statusEnable)
	if !found {
		log.Error("status_enabled tidak sesuai format")
		return lib.CustomError(http.StatusBadRequest, "status_enabled tidak sesuai format", "status_enabled tidak sesuai format")
	}
	params["verified_email"] = "0"
	params["verified_mobileno"] = "0"

	// Set expired for token
	date := time.Now().AddDate(0, 0, 30)
	dateLayout := "2006-01-02 15:04:05"
	expired := date.Format(dateLayout)

	// Generate verify key
	verifyKeyByte := sha256.Sum256([]byte(email + "_" + expired))
	verifyKey := hex.EncodeToString(verifyKeyByte[:])

	params["ulogin_locked"] = "0"
	params["ulogin_must_changepwd"] = "0"
	params["ulogin_failed_count"] = "0"
	params["string_token"] = verifyKey
	params["token_expired"] = expired
	params["ulogin_enabled"] = statusEnable
	params["rec_order"] = "0"
	params["rec_status"] = "1"
	params["rec_modified_date"] = time.Now().Format(dateLayout)
	params["rec_modified_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
	params["last_access"] = time.Now().Format(dateLayout)

	_, err = models.UpdateScUserLogin(params)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(http.StatusBadRequest, err.Error(), "Failed update user")
	}

	// Send email
	t := template.New("index-email-activation.html")

	t, err = t.ParseFiles(config.BasePath + "/mail/index-email-activation.html")
	if err != nil {
		log.Println(err)
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, struct {
		Url     string
		FileUrl string
	}{Url: config.BaseUrl + "/verifyuser?token=" + verifyKey, FileUrl: config.ImageUrl + "/images/mail"}); err != nil {
		log.Println(err)
	}

	//kirim email: JIKA nasabah Offline berubah menjadi ONLINE
	if statusEnable == "1" && ulogin_enable_prev_state == 9 {
		result := tpl.String()

		mailer := gomail.NewMessage()
		mailer.SetHeader("From", config.EmailFrom)
		mailer.SetHeader("To", email)
		mailer.SetHeader("Subject", "[MotionFunds] Verifikasi Email Kamu")
		mailer.SetBody("text/html", result)

		err = lib.SendEmail(mailer)
		if err != nil {
			log.Error(err)
			return lib.CustomError(http.StatusInternalServerError, err.Error(), "Error send email")
		} else {
			log.Info("Email sent")
		}
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = nil
	return c.JSON(http.StatusOK, response)

}
