package models

import (
	"database/sql"
	"fmt"
	"mf-bo-api/db"
	"net/http"
)

type MenuPageAccess struct {
	AccessEnable    bool `json:"access_enable"`
	AllowedCreate   bool `json:"allowed_create"`
	AllowedRead     bool `json:"allowed_read"`
	AllowedUpdate   bool `json:"allowed_update"`
	AllowedDelete   bool `json:"allowed_delete"`
	AllowedApproval bool `json:"allowed_approval"`
}

type RoleMenuList struct {
	RoleMenuKey     uint64  `db:"role_menu_key" json:"role_menu_key"`
	RoleKey         uint64  `db:"role_key" json:"role_key"`
	RoleCode        string  `db:"role_code" json:"role_code"`
	RoleName        string  `db:"role_name" json:"role_name"`
	MenuKey         uint64  `db:"menu_key" json:"menu_key"`
	MenuParentKey   *uint64 `db:"menu_parent" json:"menu_parent"`
	MenuCode        string  `db:"menu_code" json:"menu_code"`
	MenuName        string  `db:"menu_name" json:"menu_name"`
	AllowedCreate   bool    `db:"flag_create" json:"allowed_create"`
	AllowedRead     bool    `db:"flag_read" json:"allowed_read"`
	AllowedUpdate   bool    `db:"flag_update" json:"allowed_update"`
	AllowedDelete   bool    `db:"flag_approval" json:"allowed_delete"`
	AllowedApproval bool    `db:"flag_approval" json:"allowed_approval"`
	HasEndpoint     bool    `db:"has_endpoint" json:"has_endpoint"`
}

type RoleMenuDetail struct {
	RoleMenuKey     uint64  `db:"role_menu_key" json:"role_menu_key"`
	RoleKey         uint64  `db:"role_key" json:"role_key"`
	MenuKey         uint64  `db:"menu_key" json:"menu_key"`
	MenuParentKey   *uint64 `db:"menu_parent" json:"menu_parent"`
	MenuCode        string  `db:"menu_code" json:"menu_code"`
	MenuName        string  `db:"menu_name" json:"menu_name"`
	AccessEnable    bool    `db:"rec_status" json:"access_enable"`
	AllowedCreate   bool    `db:"flag_create" json:"allowed_create"`
	AllowedRead     bool    `db:"flag_read" json:"allowed_read"`
	AllowedUpdate   bool    `db:"flag_update" json:"allowed_update"`
	AllowedDelete   bool    `db:"flag_delete" json:"allowed_delete"`
	AllowedApproval bool    `db:"flag_approval" json:"allowed_approval"`
	HasEndpoint     bool    `db:"has_endpoint" json:"has_endpoint"`
}

/*
get a single row of RoleMenu

dipakai utk valiasi useraccess terhadap page/menu tertentu
*/
func GetRoleMenuAccess(ctx *RoleMenuDetail, pRoleKey uint64, pMenuKey uint64) (int, error) {
	query := `SELECT 
	a.role_menu_key, 
	a.menu_key, 
	b.menu_parent, 
	b.menu_code, 
	b.menu_name, 
	b.has_endpoint, 
	a.flag_approval, 
	a.flag_create, 
	a.flag_read, 
	a.flag_update, 
	a.flag_delete, 
	a.rec_status 
	FROM sc_role_menu a
	INNER JOIN sc_menu b ON (a.menu_key = b.menu_key)
	WHERE b.rec_status = 1
	AND a.role_key = %v
	AND a.menu_key = %v;`

	// Main query
	s_sql := fmt.Sprintf(query, pRoleKey, pMenuKey)
	//log.Println(s_sql)
	err := db.Db.Get(ctx, s_sql)
	if err != nil {
		if err != sql.ErrNoRows {
			//log.Println(err.Error())
			return http.StatusBadGateway, err
		}
	}

	return http.StatusOK, nil

}
