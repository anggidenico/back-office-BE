package models

import (
	"database/sql"
	"fmt"
	"log"
	"mf-bo-api/db"
	"net/http"
)

/* Menu model */

type ScMenuModel struct {
	MenuKey    uint64        `db:"menu_key" json:"menu_key"`
	MenuParent *uint64       `db:"menu_parent" json:"menu_parent"`
	MenuCode   string        `db:"menu_code" json:"menu_code"`
	MenuName   string        `db:"menu_name" json:"menu_name"`
	MenuPage   *string       `db:"menu_page" json:"menu_page"`
	MenuUrl    string        `db:"menu_url" json:"menu_url"`
	MenuDesc   *string       `db:"menu_desc" json:"menu_desc"`
	MenuIcon   *string       `db:"rec_attribute_id1" json:"menu_icon"`
	RecOrder   *int32        `db:"rec_order" json:"rec_order"`
	ChildList  []ScMenuModel `json:"child_list"`
}

type SideBarMenuModel struct {
	RootMenuKey      uint64        `db:"root_key" json:"root_key"`
	RootMenuCode     string        `db:"root_code" json:"root_code"`
	RootMenuName     string        `db:"root_name" json:"root_name"`
	RootMenuDesc     *string       `db:"root_desc" json:"root_desc"`
	RootHomeURL      *string       `db:"root_url" json:"root_url"`
	RootLoginURL     *string       `db:"url_login" json:"url_login"`
	RootLogoutURL    *string       `db:"url_logout" json:"url_logout"`
	RootDashboardURL *string       `db:"url_dashboard" json:"url_dashboard"`
	RootFolderPage   *string       `db:"root_folder" json:"root_folder"`
	RootMenuIcon     *string       `db:"root_icon" json:"root_icon"`
	RootMenuLogo     *string       `db:"root_logo" json:"root_logo"`
	RootMainColor    *string       `db:"color_primary" json:"color_primary"`
	RootAltColor     *string       `db:"color_secondary" json:"color_secondary"`
	RoleKey          uint64        `db:"role_key" json:"role_key"`
	RoleCode         string        `db:"role_code" json:"role_code"`
	RoleName         string        `db:"role_name" json:"role_name"`
	RoleDesc         string        `db:"role_desc" json:"role_desc"`
	MenuList         []ScMenuModel `db:"menu_list" json:"menu_list"`
}

/* Menu queries */

func GetRootMenu(ctx *SideBarMenuModel, pModule string, pRole_key uint16) (int, error) {
	module_code := "MFBO"
	if pModule != "" {
		module_code = pModule
	}
	role_key := pRole_key
	query := `select 
	COALESCE(am.rec_attribute_id3,'0') AS root_key,
	am.app_module_code AS root_code,
	am.app_module_name AS root_name,
	am.app_module_desc AS root_desc, 	
	am.url_module_home AS root_url,
	am.url_module_login AS url_login,
	am.url_module_logout AS url_logout,
	am.url_module_dasboard AS url_dashboard,
	am.app_module_prefix AS root_folder,
	am.rec_image1 AS root_icon,
	am.rec_image2 AS root_logo,
	COALESCE(am.rec_attribute_id1,'FFF') as color_primary,
	COALESCE(am.rec_attribute_id2,'CCC') AS color_secondary,	
	r.role_key,
	r.role_code, 
	r.role_name,
	r.role_desc
from sc_app_module am 
CROSS JOIN sc_role r
WHERE am.rec_status = 1 
AND r.rec_status = 1
AND am.app_module_code = '%v'
AND r.role_key = %v;`

	s_sql := fmt.Sprintf(query, module_code, role_key)
	//log.Println("[----GetRootMenu: ----]", s_sql)

	err := db.Db.Get(ctx, s_sql)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Println(err)
			return http.StatusBadGateway, err
		}
	}
	return http.StatusOK, nil
}

func GetMenuTree(ctx *[]ScMenuModel, p_role_key uint16, p_parent_key uint64) (int, error) {
	query :=
		`SELECT
		sm.menu_key,
		sm.menu_parent,
		sm.menu_code,
		sm.menu_name,
		sm.menu_page,
		COALESCE(sm.menu_url,'') as menu_url,
		sm.menu_desc,
		sm.rec_order,
		sm.rec_attribute_id1 
		FROM sc_menu sm 
INNER JOIN sc_role_menu rm ON (rm.menu_key=sm.menu_key AND rm.rec_status=1)
WHERE sm.rec_status = 1 
AND rm.role_key = %v
AND COALESCE(sm.menu_parent,0) = %v
ORDER BY sm.rec_order;`

	s_sql := fmt.Sprintf(query, p_role_key, p_parent_key)
	//log.Println("[----GetMenuTree: ----]", s_sql)
	err := db.Db.Select(ctx, s_sql)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Println(err)
			return http.StatusBadGateway, err
		}
	}

	return http.StatusOK, nil
}
