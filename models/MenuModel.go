package models

import (
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
	RootMenuKey      uint64        `db:"menu_key" json:"root_key"`
	RootMenuCode     string        `db:"menu_code" json:"root_code"`
	RootMenuName     string        `db:"menu_name" json:"root_name"`
	RootMenuDesc     *string       `db:"menu_desc" json:"root_desc"`
	RootHomeURL      string        `db:"menu_url" json:"root_url"`
	RootMenuPage     *string       `db:"menu_page" json:"root_folder"`
	RootMenuIcon     *string       `db:"rec_image1" json:"root_icon"`
	RootMenuLogo     *string       `db:"rec_image2" json:"root_logo"`
	RootMainColor    *string       `db:"rec_attribute_id1" json:"color_primary"`
	RootAltColor     *string       `db:"rec_attribute_id2" json:"color_secondary"`
	RootModuleCode   *string       `db:"app_module_code" json:"module_code"`
	RootMenuTypeCode *string       `db:"menu_type_code" json:"menu_type_code"`
	MenuList         []ScMenuModel `json:"menu_list"`
}

/* Menu queries */

func GetRootMenu(ctx *SideBarMenuModel) (int, error) {
	//menu_code := ""
	//AND sm.menu_code = '%v'

	role_key := 15
	query := `select 
	sm.menu_key,
	sm.menu_code,
	sm.menu_name,
	sm.menu_desc, 
	COALESCE(sm.menu_url,'') as menu_url,
	COALESCE(sm.menu_page,'/') as menu_page,
	sm.rec_image1,
	sm.rec_image2, 
	sm.rec_attribute_id1, 
	sm.rec_attribute_id2,
	am.app_module_code,
	mt.menu_type_code  
from sc_menu sm 
INNER JOIN sc_role_menu rm ON (rm.menu_key=sm.menu_key AND rm.rec_status=0)
INNER JOIN sc_menu_type mt ON (sm.menu_type_key = mt.menu_type_key AND mt.rec_status=1 AND mt.menu_type_key=2)
INNER JOIN sc_app_module am ON (am.app_module_key=sm.app_module_key AND am.rec_status = 1)
where sm.rec_status = 1 
AND am.app_module_key NOT IN (1)
AND COALESCE(sm.menu_parent, 0) = 0
AND rm.role_key = %v
order BY am.app_module_key, sm.rec_order;
`

	s_sql := fmt.Sprintf(query, role_key)
	log.Println("[----GetRootMenu: ----]", s_sql)

	err := db.Db.Get(ctx, s_sql)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func GetMenuTree(ctx *[]ScMenuModel, p_parent_key uint64) (int, error) {
	query :=
		`select
		sm.menu_key,
		sm.menu_parent,
		sm.menu_code,
		sm.menu_name,
		sm.menu_page,
		COALESCE(sm.menu_url,'') as menu_url,
		sm.menu_desc,
		sm.rec_order,
		sm.rec_attribute_id1 
	from sc_menu sm 
	where sm.rec_status = 1 
	and COALESCE(menu_parent, 0) = %v 
	order by rec_order`

	s_sql := fmt.Sprintf(query, p_parent_key)
	log.Println("[----GetMenuTree: ----]", s_sql)

	err := db.Db.Select(ctx, s_sql)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}
