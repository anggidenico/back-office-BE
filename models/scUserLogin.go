package models

import (
	"database/sql"
	"mf-bo-api/db"
	"net/http"
	"strconv"
	"strings"
)

type ScUserLoginRegister struct {
	UloginEmail    string `json:"ulogin_email"`
	UloginMobileno string `json:"ulogin_mobileno"`
}

type ScUserLogin struct {
	UserLoginKey         uint64  `db:"user_login_key"            json:"user_login_key"`
	UserCategoryKey      uint64  `db:"user_category_key"         json:"user_category_key"`
	CustomerKey          *uint64 `db:"customer_key"              json:"customer_key"`
	RoleKey              *uint64 `db:"role_key"                  json:"role_key"`
	UserDeptKey          *uint64 `db:"user_dept_key"             json:"user_dept_key"`
	UserDeptKey1         *uint64 `db:"user_dept_key1"            json:"user_dept_key1"`
	UloginName           string  `db:"ulogin_name"               json:"ulogin_name"`
	UloginFullName       string  `db:"ulogin_full_name"          json:"ulogin_full_name"`
	UloginPassword       string  `db:"ulogin_password"           json:"ulogin_password"`
	UloginEmail          string  `db:"ulogin_email"              json:"ulogin_email"`
	VerifiedEmail        *uint8  `db:"verified_email"            json:"verified_email"`
	LastVerifiedEmail    *string `db:"last_verified_email"       json:"last_verified_email"`
	StringToken          *string `db:"string_token"              json:"string_token"`
	TokenExpired         *string `db:"token_expired"             json:"token_expired"`
	UloginPin            *string `db:"ulogin_pin"                json:"ulogin_pin"`
	LastChangedPin       *string `db:"last_changed_pin"          json:"last_changed_pin"`
	MustChangePin        *uint8  `db:"must_change_pin"           json:"must_change_pin"`
	UloginMobileno       *string `db:"ulogin_mobileno"           json:"ulogin_mobileno"`
	LastVerifiedMobileno *string `db:"last_verified_mobileno"    json:"last_verified_mobileno"`
	VerifiedMobileno     uint8   `db:"verified_mobileno"         json:"verified_mobileno"`
	UloginMustChangepwd  uint8   `db:"ulogin_must_changepwd"     json:"ulogin_must_changepwd"`
	LastPasswordChanged  *string `db:"last_password_changed"     json:"last_password_changed"`
	OtpNumber            *string `db:"otp_number"                json:"otp_number"`
	OtpNumberExpired     *string `db:"otp_number_expired"        json:"otp_number_expired"`
	UloginLocked         uint8   `db:"ulogin_locked"             json:"ulogin_locked"`
	UloginEnabled        uint8   `db:"ulogin_enabled"            json:"ulogin_enabled"`
	UloginFailedCount    uint64  `db:"ulogin_failed_count"       json:"ulogin_failed_count"`
	LastAccess           *string `db:"last_access"               json:"last_access"`
	AcceptLoginTnc       uint8   `db:"accept_login_tnc"          json:"accept_login_tnc"`
	AllowedSharingLogin  *uint8  `db:"allowed_sharing_login"     json:"allowed_sharing_login"`
	TokenNotif           *string `db:"token_notif"               json:"token_notif"`
	LastUpdateTokenNotif *string `db:"last_update_token_notif"   json:"last_update_token_notif"`
	LockedDate           *string `db:"locked_date"               json:"locked_date"`
	RecOrder             *uint64 `db:"rec_order"                 json:"rec_order"`
	RecStatus            uint8   `db:"rec_status"                json:"rec_status"`
	RecCreatedDate       *string `db:"rec_created_date"          json:"rec_created_date"`
	RecCreatedBy         *string `db:"rec_created_by"            json:"rec_created_by"`
	RecModifiedDate      *string `db:"rec_modified_date"         json:"rec_modified_date"`
	RecModifiedBy        *string `db:"rec_modified_by"           json:"rec_modified_by"`
	RecImage1            *string `db:"rec_image1"                json:"rec_image1"`
	RecImage2            *string `db:"rec_image2"                json:"rec_image2"`
	RecApprovalStatus    *uint8  `db:"rec_approval_status"       json:"rec_approval_status"`
	RecApprovalStage     *uint64 `db:"rec_approval_stage"        json:"rec_approval_stage"`
	RecApprovedDate      *string `db:"rec_approved_date"         json:"rec_approved_date"`
	RecApprovedBy        *string `db:"rec_approved_by"           json:"rec_approved_by"`
	RecDeletedDate       *string `db:"rec_deleted_date"          json:"rec_deleted_date"`
	RecDeletedBy         *string `db:"rec_deleted_by"            json:"rec_deleted_by"`
	RecAttributeID1      *string `db:"rec_attribute_id1"         json:"rec_attribute_id1"`
	RecAttributeID2      *string `db:"rec_attribute_id2"         json:"rec_attribute_id2"`
	RecAttributeID3      *string `db:"rec_attribute_id3"         json:"rec_attribute_id3"`
}

type AdminListScUserLogin struct {
	UserLoginKey     uint64  `db:"user_login_key"      json:"user_login_key"`
	UloginEnabled    uint64  `db:"ulogin_enabled"      json:"ulogin_enabled"`
	UCategoryName    string  `db:"ucategory_name"      json:"ucategory_name"`
	UserDeptName     *string `db:"user_dept_name"      json:"user_dept_name"`
	UloginName       string  `db:"ulogin_name"         json:"ulogin_name"`
	UloginFullName   string  `db:"ulogin_full_name"    json:"ulogin_full_name"`
	UloginEmail      string  `db:"ulogin_email"        json:"ulogin_email"`
	IsLogin          string  `db:"is_login"            json:"is_login"`
	RoleKey          *uint64 `db:"role_key"            json:"role_key"`
	RoleName         *string `db:"role_name"           json:"role_name"`
	Enabled          string  `db:"enabled"             json:"enabled"`
	Locked           string  `db:"locked"              json:"locked"`
	CreatedDate      *string `db:"created_date"        json:"created_date"`
	UloginMobileno   *string `db:"ulogin_mobileno"     json:"ulogin_mobileno"`
	VerifiedMobileno *uint8  `db:"verified_mobileno"   json:"verified_mobileno"`
	VerifiedEmail    *uint8  `db:"verified_email"      json:"verified_email"`
	CustomerKey      *uint64 `db:"customer_key"       json:"customer_key"`
}

type AdminDetailScUserLogin struct {
	UserLoginKey     uint64             `json:"user_login_key"`
	UserCategory     ScUserCategoryInfo `json:"user_category"`
	UserDept         *ScUserDeptInfo    `json:"user_dept"`
	UloginName       string             `json:"ulogin_name"`
	UloginFullName   string             `json:"ulogin_full_name"`
	UloginEmail      string             `json:"ulogin_email"`
	Role             *ScRoleInfoLogin   `json:"role"`
	Enabled          bool               `json:"enabled"`
	Locked           bool               `json:"locked"`
	VerifiedEmail    bool               `json:"verified_email"`
	VerifiedMobileno bool               `json:"verified_mobileno"`
	CreatedDate      *string            `json:"created_date"`
	RecImage         string             `json:"rec_image"`
	NoHp             *string            `json:"no_hp"`
	Cif              *string            `json:"cif"`
	CustomerName     *string            `json:"customer_name"`
	CustomerKey      *uint64            `json:"customer_key"`
	OaRequestKey     *uint64            `json:"oa_request_key"`
}

type AdminListScUserLoginRole struct {
	UloginName     string `json:"ulogin_name"`
	UloginFullName string `json:"ulogin_full_name"`
	UloginEmail    string `json:"ulogin_email"`
}

type UserLoginKeyLocked struct {
	UserLoginKey uint64 `db:"user_login_key"            json:"user_login_key"`
}
type UserBlastPromo struct {
	UserLoginKey uint64  `db:"user_login_key"       json:"user_login_key"`
	TokenNotif   *string `db:"token_notif"          json:"token_notif"`
	FirstName    *string `db:"first_name"           json:"first_name"`
}

func GetAllScUserLogin(c *[]ScUserLogin, limit uint64, offset uint64, params map[string]string, nolimit bool) (int, error) {
	query := `SELECT
              sc_user_login.* FROM 
			  sc_user_login`
	var present bool
	var whereClause []string
	var condition string

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, "sc_user_login."+field+" = '"+value+"'")
		}
	}

	// Combile where clause
	if len(whereClause) > 0 {
		condition += " WHERE "
		for index, where := range whereClause {
			condition += where
			if (len(whereClause) - 1) > index {
				condition += " AND "
			}
		}
	}
	// Check order by
	var orderBy string
	var orderType string
	if orderBy, present = params["orderBy"]; present == true {
		condition += " ORDER BY " + orderBy
		if orderType, present = params["orderType"]; present == true {
			condition += " " + orderType
		}
	}
	query += condition

	// Query limit and offset
	if !nolimit {
		query += " LIMIT " + strconv.FormatUint(limit, 10)
		if offset > 0 {
			query += " OFFSET " + strconv.FormatUint(offset, 10)
		}
	}

	// Main query
	// log.Println("==========  ==========>>>", query)
	err := db.Db.Select(c, query)
	if err != nil {
		// log.Error(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func GetScUserLogin(c *ScUserLogin, email string) (int, error) {
	query := `SELECT sc_user_login.* WHERE sc_user_login.ulogin_email = ` + email
	// log.Println("==========  ==========>>>", query)
	err := db.Db.Get(c, query)
	if err != nil {
		// log.Error(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}

func CreateScUserLogin(params map[string]string) (int, error) {
	query := "INSERT INTO sc_user_login"
	// Get params
	var fields, values string
	var bindvars []interface{}
	for key, value := range params {
		fields += key + ", "
		values += "?, "
		bindvars = append(bindvars, value)
	}
	fields = fields[:(len(fields) - 2)]
	values = values[:(len(values) - 2)]

	// Combine params to build query
	query += "(" + fields + ") VALUES(" + values + ")"
	// log.Println("==========  ==========>>>", query)

	tx, err := db.Db.Begin()
	if err != nil {
		// log.Error(err)
		return http.StatusBadGateway, err
	}
	_, err = tx.Exec(query, bindvars...)
	tx.Commit()
	if err != nil {
		// log.Error(err)
		return http.StatusBadRequest, err
	}
	return http.StatusOK, nil
}

func UpdateScUserLogin(params map[string]string) (int, error) {
	query := "UPDATE sc_user_login SET "
	// Get params
	i := 0
	for key, value := range params {
		if key != "user_login_key" {

			query += key + " = '" + value + "'"

			if (len(params) - 2) > i {
				query += ", "
			}
			i++
		}
	}
	query += " WHERE user_login_key = " + params["user_login_key"]
	// log.Println("========== QUERY UPDATE SC USER LOGIN ==========>>>", query)

	tx, err := db.Db.Begin()
	if err != nil {
		// log.Error(err)
		return http.StatusBadGateway, err
	}
	// var ret sql.Result
	_, err = tx.Exec(query)

	// // log.Infoln("========== query dan parameter upload image ==========")
	// log.Println("========== QUERY UPDATE SC USER LOGIN ==========", query)
	// // log.Infoln("========================================")
	// // log.Infoln(params)
	// // log.Infoln("========================================")

	//banyak transaction di DB ke lock, sementara di disabled dlu
	// row, _ := ret.RowsAffected()
	// if row > 0 {
	// 	tx.Commit()
	// } else {
	// 	return http.StatusNotFound, err
	// }
	if err != nil {
		tx.Rollback()
		// log.Error(err)
		return http.StatusBadRequest, err
	}
	tx.Commit()
	return http.StatusOK, nil
}

func GetScUserLoginIn(c *[]ScUserLogin, value []string, field string) (int, error) {
	inQuery := strings.Join(value, ",")
	query2 := `SELECT 
				sc_user_login.* FROM 
				sc_user_login`
	query := query2 + " WHERE sc_user_login." + field + " IN(" + inQuery + ")"

	// Main query
	// log.Println("==========  ==========>>>", query)
	err := db.Db.Select(c, query)
	if err != nil {
		// log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func GetScUserLoginByKey(c *ScUserLogin, key string) (int, error) {
	query := `SELECT sc_user_login.* FROM sc_user_login WHERE sc_user_login.rec_status = 1 AND sc_user_login.user_login_key = ` + key
	// log.Println("==========  ==========>>>", query)
	err := db.Db.Get(c, query)
	if err != nil {
		// log.Println(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}

func GetScUserLoginByCustomerKey(c *ScUserLogin, key string) (int, error) {
	query := `SELECT sc_user_login.* FROM sc_user_login WHERE sc_user_login.rec_status = 1 AND sc_user_login.customer_key = ` + key
	// log.Println("==========  ==========>>>", query)
	err := db.Db.Get(c, query)
	if err != nil {
		// log.Println(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}

func AdminGetAllScUserLogin(c *[]AdminListScUserLogin, limit uint64, offset uint64, params map[string]string, nolimit bool, searchLike *string) (int, error) {
	query := `SELECT
				u.user_login_key AS user_login_key, 
				u.ulogin_enabled AS ulogin_enabled,
				cat.ucategory_name AS ucategory_name,
				dept.user_dept_name AS user_dept_name, 
				u.ulogin_name AS ulogin_name, 
				u.ulogin_full_name AS ulogin_full_name,
				u.ulogin_email AS ulogin_email,
				(CASE
					WHEN ses.login_session_key IS NULL THEN "No"
					ELSE 
					(CASE
					WHEN ses.logout_date IS NOT NULL THEN "No"
					WHEN DATE_ADD(ses.login_date, INTERVAL 2 HOUR) > NOW() THEN "Yes"
					ELSE "No"
					END)
				END) AS is_login,
				role.role_key AS role_key,
				role.role_name AS role_name,
				(CASE
					WHEN u.ulogin_enabled = '1' THEN 'Yes'
					ELSE 'No'
				END) AS enabled,
				(CASE
					WHEN u.ulogin_locked = '1' THEN 'Yes'
					ELSE 'No'
				END) AS locked,
				DATE_FORMAT(u.rec_created_date, '%d %M %Y') AS created_date,
				u.verified_email,
				u.ulogin_mobileno,
				u.verified_mobileno,
				u.customer_key 
			  FROM sc_user_login AS u 
			  LEFT JOIN sc_role AS role ON u.role_key = role.role_key 
			  LEFT JOIN sc_user_category AS cat ON cat.user_category_key = u.user_category_key 
			  LEFT JOIN sc_user_dept AS dept ON dept.user_dept_key = u.user_dept_key 
			  LEFT JOIN sc_login_session AS ses ON ses.user_login_key = u.user_login_key 
			  WHERE u.rec_status = 1`
	var present bool
	var whereClause []string
	var condition string

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, field+" = '"+value+"'")
		}
	}

	// Combile where clause
	if len(whereClause) > 0 {
		condition += " AND "
		for index, where := range whereClause {
			condition += where
			if (len(whereClause) - 1) > index {
				condition += " AND "
			}
		}
	}

	//search like all
	if searchLike != nil {
		condition += " AND ("
		condition += " u.user_login_key LIKE '%" + *searchLike + "%' OR"
		condition += " cat.ucategory_name LIKE '%" + *searchLike + "%' OR"
		condition += " dept.user_dept_name LIKE '%" + *searchLike + "%' OR"
		condition += " u.ulogin_name LIKE '%" + *searchLike + "%' OR"
		condition += " u.ulogin_full_name LIKE '%" + *searchLike + "%' OR"
		condition += " u.ulogin_email LIKE '%" + *searchLike + "%' OR"
		condition += " role.role_name LIKE '%" + *searchLike + "%' OR"
		condition += " DATE_FORMAT(u.rec_created_date, '%d %M %Y') LIKE '%" + *searchLike + "%' )"
	}

	// Check order by
	var orderBy string
	var orderType string
	if orderBy, present = params["orderBy"]; present == true {
		condition += " ORDER BY " + orderBy
		if orderType, present = params["orderType"]; present == true {
			condition += " " + orderType
		}
	}
	query += condition

	// Query limit and offset
	if !nolimit {
		query += " LIMIT " + strconv.FormatUint(limit, 10)
		if offset > 0 {
			query += " OFFSET " + strconv.FormatUint(offset, 10)
		}
	}

	// Main query
	// log.Println("==========  ==========>>>", query)
	err := db.Db.Select(c, query)
	if err != nil {
		// log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func AdminCountDataGetAllScUserlogin(c *CountData, params map[string]string, searchLike *string) (int, error) {
	query := `SELECT
	            count(u.user_login_key) AS count_data
			  FROM sc_user_login AS u 
			  LEFT JOIN sc_role AS role ON u.role_key = role.role_key 
			  LEFT JOIN sc_user_category AS cat ON cat.user_category_key = u.user_category_key 
			  LEFT JOIN sc_user_dept AS dept ON dept.user_dept_key = u.user_dept_key 
			  LEFT JOIN sc_login_session AS ses ON ses.user_login_key = u.user_login_key 
			  WHERE u.rec_status = 1`
	var whereClause []string
	var condition string

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, field+" = '"+value+"'")
		}
	}

	// Combile where clause
	if len(whereClause) > 0 {
		condition += " AND "
		for index, where := range whereClause {
			condition += where
			if (len(whereClause) - 1) > index {
				condition += " AND "
			}
		}
	}

	//search like all
	if searchLike != nil {
		condition += " AND ("
		condition += " u.user_login_key LIKE '%" + *searchLike + "%' OR"
		condition += " cat.ucategory_name LIKE '%" + *searchLike + "%' OR"
		condition += " dept.user_dept_name LIKE '%" + *searchLike + "%' OR"
		condition += " u.ulogin_name LIKE '%" + *searchLike + "%' OR"
		condition += " u.ulogin_full_name LIKE '%" + *searchLike + "%' OR"
		condition += " u.ulogin_email LIKE '%" + *searchLike + "%' OR"
		condition += " role.role_name LIKE '%" + *searchLike + "%' OR"
		condition += " DATE_FORMAT(u.rec_created_date, '%d %M %Y') LIKE '%" + *searchLike + "%' )"
	}

	query += condition

	// Main query
	// log.Println("==========  ==========>>>", query)
	err := db.Db.Get(c, query)
	if err != nil {
		// log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func AdminGetValidateUniqueInsertUpdateScUserLogin(c *CountData, paramsOr map[string]string, updateKey *string) (int, error) {
	query := `SELECT
			  count(sc_user_login.user_login_key) as count_data 
			  FROM sc_user_login `
	var orWhereClause []string
	var condition string

	for fieldOr, valueOr := range paramsOr {
		orWhereClause = append(orWhereClause, "sc_user_login."+fieldOr+" = '"+valueOr+"'")
	}

	// Combile where Or clause
	if len(orWhereClause) > 0 {
		condition += " WHERE ("
		for index, where := range orWhereClause {
			condition += where
			if (len(orWhereClause) - 1) > index {
				condition += " OR "
			} else {
				condition += ") "
			}
		}
	}

	if updateKey != nil {
		if len(orWhereClause) > 0 {
			condition += " AND "
		} else {
			condition += " WHERE "
		}

		condition += " sc_user_login.user_login_key != '" + *updateKey + "'"
	}

	query += condition

	// Main query
	// log.Println("==========  ==========>>>", query)
	err := db.Db.Get(c, query)
	if err != nil {
		// log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func GetCountScUserLogin(c *CountData, params map[string]string) (int, error) {
	query := `SELECT
			  count(sc_user_login.user_login_key) as count_data 
			  FROM 
			  sc_user_login`
	var whereClause []string
	var condition string

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, "sc_user_login."+field+" = '"+value+"'")
		}
	}

	// Combile where clause
	if len(whereClause) > 0 {
		condition += " WHERE "
		for index, where := range whereClause {
			condition += where
			if (len(whereClause) - 1) > index {
				condition += " AND "
			}
		}
	}

	query += condition

	// Main query
	// log.Println("==========  ==========>>>", query)
	err := db.Db.Get(c, query)
	if err != nil {
		// log.Error(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func GetAllScUserLoginByNameOrEmail(c *[]ScUserLogin, limit uint64, offset uint64, params map[string]string, nolimit bool) (int, error) {
	query := `SELECT
              sc_user_login.* FROM 
			  sc_user_login`
	var present bool
	var whereClause []string
	var orClause []string
	var condition string

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType" || field == "ulogin_name" || field == "ulogin_email") {
			whereClause = append(whereClause, "sc_user_login."+field+" = '"+value+"'")
		}
		if field == "ulogin_name" || field == "ulogin_email" {
			orClause = append(orClause, "sc_user_login."+field+" = '"+value+"'")
		}
	}

	// Combile where clause
	if len(whereClause) > 0 {
		condition += " WHERE "
		for index, where := range whereClause {
			condition += where
			if (len(whereClause) - 1) > index {
				condition += " AND "
			}
		}
	}

	// Combile where or clause
	if len(whereClause) > 0 {
		condition += " AND ("
		for index, where := range orClause {
			condition += where
			if (len(orClause) - 1) > index {
				condition += " OR "
			} else {
				condition += " ) "
			}
		}
	} else {
		condition += " WHERE ("
		for index, where := range orClause {
			condition += where
			if (len(orClause) - 1) > index {
				condition += " OR "
			} else {
				condition += " ) "
			}
		}
	}

	condition += " AND sc_user_login.user_category_key IN (2,3) "

	// Check order by
	var orderBy string
	var orderType string
	if orderBy, present = params["orderBy"]; present == true {
		condition += " ORDER BY " + orderBy
		if orderType, present = params["orderType"]; present == true {
			condition += " " + orderType
		}
	}
	query += condition

	// Query limit and offset
	if !nolimit {
		query += " LIMIT " + strconv.FormatUint(limit, 10)
		if offset > 0 {
			query += " OFFSET " + strconv.FormatUint(offset, 10)
		}
	}

	// Main query
	// log.Println("========= login ============", query)
	err := db.Db.Select(c, query)
	if err != nil {
		// log.Error(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func CreateScUserLoginReturnKey(params map[string]string) (int, error, string) {
	query := "INSERT INTO sc_user_login"
	// Get params
	var fields, values string
	var bindvars []interface{}
	for key, value := range params {
		fields += key + ", "
		values += "?, "
		bindvars = append(bindvars, value)
	}
	fields = fields[:(len(fields) - 2)]
	values = values[:(len(values) - 2)]

	// Combine params to build query
	query += "(" + fields + ") VALUES(" + values + ")"
	// log.Println("==========  ==========>>>", query)

	tx, err := db.Db.Begin()
	if err != nil {
		// log.Error(err)
		return http.StatusBadGateway, err, "0"
	}
	var ret sql.Result
	ret, err = tx.Exec(query, bindvars...)
	tx.Commit()
	if err != nil {
		// log.Error(err)
		return http.StatusBadRequest, err, "0"
	}
	lastID, _ := ret.LastInsertId()
	return http.StatusOK, nil, strconv.FormatInt(lastID, 10)
}

func GetUserLocked(c *[]UserLoginKeyLocked) (int, error) {
	query := `SELECT 
				user_login_key 
			FROM sc_user_login 
			WHERE ulogin_locked = 1 AND rec_status = 1
			AND DATE_ADD(locked_date, INTERVAL 1 HOUR) < NOW()`
	// log.Println("==========  ==========>>>", query)
	err := db.Db.Select(c, query)
	if err != nil {
		// log.Println(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}

func UpdateScUserLoginByKeyIn(params map[string]string, valueIn []string, fieldIn string) (int, error) {
	query := "UPDATE sc_user_login SET "
	// Get params
	i := 0
	for key, value := range params {
		query += key + " = '" + value + "'"

		if (len(params) - 1) > i {
			query += ", "
		}
		i++
	}

	inQuery := strings.Join(valueIn, ",")
	query += " WHERE sc_user_login." + fieldIn + " IN(" + inQuery + ")"

	// log.Println("==========  ==========>>>", query)

	tx, err := db.Db.Begin()
	if err != nil {
		// log.Error(err)
		return http.StatusBadGateway, err
	}
	var ret sql.Result
	ret, err = tx.Exec(query)
	row, _ := ret.RowsAffected()
	tx.Commit()
	if row > 0 {
	} else {
		return http.StatusNotFound, err
	}
	if err != nil {
		// log.Error(err)
		return http.StatusBadRequest, err
	}
	return http.StatusOK, nil
}

func AdminGetAllUserBlastPromo(c *[]UserBlastPromo) (int, error) {
	query := `SELECT
				u.user_login_key,
				u.token_notif,
				c.first_name 
			FROM sc_user_login AS u 
			LEFT JOIN ms_customer AS c ON u.customer_key = c.customer_key
			WHERE u.user_category_key = 1 AND u.rec_status = 1 AND u.token_notif IS NOT NULL`

	// Main query
	// log.Println("==========  ==========>>>", query)
	err := db.Db.Select(c, query)
	if err != nil {
		// log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func ValidateUniqueData(c *CountData, field string, value string, userLoginKey *string) (int, error) {
	var query string
	query = `SELECT
				count(user_login_key) AS count_data
              FROM sc_user_login where rec_status = '1' AND ` + field + ` = '` + value + `'`

	if userLoginKey != nil {
		query += ` AND user_login_key != '` + *userLoginKey + `'`
	}

	// Main query
	// log.Println("==========  ==========>>>", query)
	err := db.Db.Get(c, query)
	if err != nil {
		// log.Error(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func CheckCreatePin(c *CountData, userLoginKey string) (int, error) {
	var query string
	query = `SELECT 
				COUNT(u.user_login_key) AS count_data
			FROM sc_user_login AS u
			INNER JOIN oa_request AS o ON o.user_login_key = u.user_login_key
			INNER JOIN ms_customer AS c ON c.customer_key = o.customer_key
			WHERE u.rec_status = 1 AND o.rec_status = 1 AND c.rec_status = 1 
			AND (u.ulogin_pin IS NULL OR u.ulogin_pin = "" OR u.must_change_pin = "1") AND c.customer_key IS NOT NULL 
			AND u.user_login_key = '` + userLoginKey + `' GROUP BY u.user_login_key`

	// Main query
	// log.Println("==========  ==========>>>", query)
	err := db.Db.Get(c, query)
	if err != nil {
		// log.Error(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func SetNullTokenNotif(tokenNotif string) (int, error) {
	query := `UPDATE sc_user_login SET token_notif = NULL WHERE token_notif = "` + tokenNotif + `"`

	tx, err := db.Db.Begin()
	if err != nil {
		// log.Error(err)
		return http.StatusBadGateway, err
	}

	_, err = tx.Exec(query)

	if err != nil {
		tx.Rollback()
		// log.Error(err)
		return http.StatusBadRequest, err
	}
	tx.Commit()
	return http.StatusOK, nil
}

func GetScUserKey(c *ScUserLogin, key string) (int, error) {
	query := `SELECT sc_user_login.* FROM sc_user_login WHERE sc_user_login.user_login_key = ` + key
	// log.Println("==========  ==========>>>", query)
	err := db.Db.Get(c, query)
	if err != nil {
		// log.Println(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}

func CreateScUserLoginWithReturnPK(params map[string]string) (int, error, string) {
	query := "INSERT INTO sc_user_login"
	// Get params
	var fields, values string
	var bindvars []interface{}
	for key, value := range params {
		fields += key + ", "
		values += "?, "
		bindvars = append(bindvars, value)
	}
	fields = fields[:(len(fields) - 2)]
	values = values[:(len(values) - 2)]

	// Combine params to build query
	query += "(" + fields + ") VALUES(" + values + ")"
	// log.Println("==========  ==========>>>", query)

	tx, err := db.Db.Begin()
	if err != nil {
		// log.Error(err)
		return http.StatusBadGateway, err, "0"
	}
	ret, err := tx.Exec(query, bindvars...)
	tx.Commit()
	if err != nil {
		// log.Error(err)
		return http.StatusBadRequest, err, "0"
	}
	lastID, _ := ret.LastInsertId()
	return http.StatusOK, nil, strconv.FormatInt(lastID, 10)
}

func UpdateDeleteUserNotOA(days string) (int, error) {
	query := `UPDATE sc_user_login SET rec_deleted_by = "CRON DELETE USER", 
				rec_deleted_date = NOW(), rec_status = 0 
			WHERE user_login_key IN (
				SELECT 
					u.user_login_key 
				FROM sc_user_login AS u 
				LEFT JOIN oa_request AS oa ON oa.user_login_key = u.user_login_key 
				WHERE (u.role_key IS NULL OR u.role_key IN (1)) 
				AND u.user_category_key = 1 
				AND u.user_dept_key = 1 AND u.rec_status = 1
				AND oa.oa_request_key IS NULL 
				AND DATE_ADD(u.rec_created_date, INTERVAL ` + days + ` DAY) < NOW()
			)`

	tx, err := db.Db.Begin()
	if err != nil {
		// log.Error(err)
		return http.StatusBadGateway, err
	}
	var ret sql.Result
	ret, err = tx.Exec(query)
	row, _ := ret.RowsAffected()
	tx.Commit()
	if row > 0 {
	} else {
		return http.StatusNotFound, err
	}
	if err != nil {
		// log.Error(err)
		return http.StatusBadRequest, err
	}
	return http.StatusOK, nil
}
