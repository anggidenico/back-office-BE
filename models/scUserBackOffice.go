package models

import (
	"log"
	"mf-bo-api/db"
)

type UserCategoryBackOfficeList struct {
	UserLoginKey   *uint64 `db:"user_login_key" json:"user_login_key"`
	Email          *string `db:"email" json:"email"`
	CategoryName   *string `db:"category_name" json:"category_name"`
	DepartmentName *string `db:"department_name" json:"department_name"`
	RoleName       *string `db:"role_name" json:"role_name"`
	EnableStatus   *bool   `db:"enable_status" json:"enable_status"`
	LockedStatus   *bool   `db:"lock_status" json:"lock_status"`
}

func GetUserBackOfficeCategoryList() []UserCategoryBackOfficeList {
	query := `SELECT t1.user_login_key, t1.ulogin_email AS email, t4.role_name,
	t2.ucategory_name AS category_name, t3.user_dept_name AS department_name,
	CASE
    	WHEN t1.ulogin_enabled = 1 THEN 'true'
    	ELSE 'false'
	END AS enable_status, 
	CASE
    	WHEN t1.ulogin_locked = 1 THEN 'true'
    	ELSE 'false'
	END AS lock_status
	FROM sc_user_login t1
	INNER JOIN sc_user_category t2 ON t2.user_category_key = t1.user_category_key 
	INNER JOIN sc_user_dept t3 ON t3.user_dept_key = t1.user_dept_key
	INNER JOIN sc_role t4 ON t4.role_key = t1.role_key
	WHERE t1.rec_status = 1 AND t2.user_category_key IN (2,3)`

	var result []UserCategoryBackOfficeList
	err := db.Db.Select(&result, query)
	if err != nil {
		log.Println(err.Error())
	}

	return result
}
