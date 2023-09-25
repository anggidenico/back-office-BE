package models

import (
	"log"
	"math"
	"mf-bo-api/db"
	"strconv"
)

type UserCategoryCustomerList struct {
	UserLoginKey     *uint64 `db:"user_login_key" json:"user_login_key"`
	CustomerKey      *uint64 `db:"customer_key" json:"customer_key"`
	FullName         *string `db:"full_name" json:"full_name"`
	Email            *string `db:"email" json:"email"`
	EmailVerifStatus *bool   `db:"email_verif_status" json:"email_verif_status"`
	PhoneMobile      *string `db:"phone_mobile" json:"phone_mobile"`
	PhoneVerifStatus *bool   `db:"phone_verif_status" json:"phone_verif_status"`
	CategoryName     *string `db:"category_name" json:"category_name"`
	DepartmentName   *string `db:"department_name" json:"department_name"`
	RoleName         *string `db:"role_name" json:"role_name"`
	EnableStatus     *bool   `db:"enable_status" json:"enable_status"`
	LockedStatus     *bool   `db:"lock_status" json:"lock_status"`
	// Source           *string `db:"source" json:"source"`
}

func GetUserCategoryCustomerList(params map[string]string, limit uint64, offset uint64) ([]UserCategoryCustomerList, int) {
	query := `SELECT t1.user_login_key, t1.customer_key, t1.ulogin_full_name AS full_name, 
	t1.ulogin_email AS email, t1.ulogin_mobileno AS phone_mobile,
	CASE WHEN t1.verified_email = 1 THEN 'true' ELSE 'false' END AS email_verif_status, 
	CASE WHEN t1.verified_mobileno = 1 THEN 'true' ELSE 'false' END AS phone_verif_status, 
	t4.role_name, t2.ucategory_name AS category_name, t3.user_dept_name AS department_name,
	CASE WHEN t1.ulogin_enabled = 1 THEN 'true' ELSE 'false' END AS enable_status, 
	CASE WHEN t1.ulogin_locked = 1 THEN 'true' ELSE 'false' END AS lock_status
	FROM sc_user_login t1
	INNER JOIN sc_user_category t2 ON t2.user_category_key = t1.user_category_key 
	INNER JOIN sc_user_dept t3 ON t3.user_dept_key = t1.user_dept_key
	INNER JOIN sc_role t4 ON t4.role_key = t1.role_key
	INNER JOIN ms_customer t5 ON t5.customer_key = t1.customer_key
	WHERE t1.rec_status = 1 `

	if valueMap, ok := params["email"]; ok {
		query += `AND t1.ulogin_email LIKE '%` + valueMap + `%'`
	}
	if valueMap, ok := params["full_name"]; ok {
		query += `AND t1.ulogin_full_name LIKE '%` + valueMap + `%'`
	}
	if valueMap, ok := params["phone_mobile"]; ok {
		query += `AND t1.ulogin_mobileno LIKE '%` + valueMap + `%'`
	}
	if valueMap, ok := params["user_category"]; ok {
		query += `AND t1.user_category_key = ` + valueMap
	}

	queryCountPage := `SELECT count(*) FROM ( ` + query + `) t1`

	if limit > 0 {
		query += " LIMIT " + strconv.FormatUint(limit, 10)
		if offset > 0 {
			query += " OFFSET " + strconv.FormatUint(offset, 10)
		}
	}

	var result []UserCategoryCustomerList
	err := db.Db.Select(&result, query)
	if err != nil {
		log.Println(err.Error())
	}

	var pagination int
	var count uint64
	err = db.Db.Get(&count, queryCountPage)
	if limit > 0 {
		if count < limit {
			pagination = 1
		} else {
			calc := math.Ceil(float64(count) / float64(limit))
			pagination = int(calc)
		}
	}

	return result, pagination
}

type UserCategoryBackOfficeList struct {
	UserLoginKey   *uint64 `db:"user_login_key" json:"user_login_key"`
	Email          *string `db:"email" json:"email"`
	CategoryName   *string `db:"category_name" json:"category_name"`
	DepartmentName *string `db:"department_name" json:"department_name"`
	RoleName       *string `db:"role_name" json:"role_name"`
	EnableStatus   *bool   `db:"enable_status" json:"enable_status"`
	LockedStatus   *bool   `db:"lock_status" json:"lock_status"`
}

func GetUserCategoryBackOfficeList(params map[string]string, limit uint64, offset uint64) ([]UserCategoryBackOfficeList, int) {
	query := `SELECT t1.user_login_key, t1.ulogin_email AS email, t4.role_name,
	t2.ucategory_name AS category_name, t3.user_dept_name AS department_name,
	CASE WHEN t1.ulogin_enabled = 1 THEN 'true' ELSE 'false' END AS enable_status, 
	CASE WHEN t1.ulogin_locked = 1 THEN 'true' ELSE 'false' END AS lock_status
	FROM sc_user_login t1
	INNER JOIN sc_user_category t2 ON t2.user_category_key = t1.user_category_key 
	INNER JOIN sc_user_dept t3 ON t3.user_dept_key = t1.user_dept_key
	INNER JOIN sc_role t4 ON t4.role_key = t1.role_key
	WHERE t1.rec_status = 1 AND t2.user_category_key IN (2,3)`

	if valueMap, ok := params["email"]; ok {
		query += `AND t1.ulogin_email LIKE '%` + valueMap + `%'`
		// queryCountPage += `AND t3.idcard_no = ` + valueMap
	}

	queryCountPage := `SELECT count(*) FROM
	( ` + query + `) t1`

	if limit > 0 {
		query += " LIMIT " + strconv.FormatUint(limit, 10)
		if offset > 0 {
			query += " OFFSET " + strconv.FormatUint(offset, 10)
		}
	}

	var result []UserCategoryBackOfficeList
	err := db.Db.Select(&result, query)
	if err != nil {
		log.Println(err.Error())
	}

	var pagination int

	var count uint64
	err = db.Db.Get(&count, queryCountPage)
	if err != nil {
		log.Println(err.Error())
	}

	if limit > 0 {
		if count < limit {
			pagination = 1
		} else {
			calc := math.Ceil(float64(count) / float64(limit))
			pagination = int(calc)
		}
	}

	return result, pagination
}
