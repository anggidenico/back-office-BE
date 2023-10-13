package models

import (
	"log"
	"mf-bo-api/db"
)

type scEndpoint struct {
	EndpointCode       string `db:"endpoint_code" json:"endpoint_code"`
	EndpointName       string `db:"endpoint_name" json:"endpoint_name"`
	Method             string `db:"endpoint_verb" json:"method"`
	EndpointUrl        string `db:"endpoint_uri" json:"url"`
	EndpointCtgCode    string `db:"endpoint_ctg_code" json:"endpoint_ctg_code"`
	EndpointCtgDesc    string `db:"endpoint_ctg_desc" json:"endpoint_ctg_desc"`
	EndpointCtgPurpose string `db:"endpoint_ctg_purpose" json:"endpoint_ctg_purpose"`
	MenuCode           string `db:"menu_code" json:"menu_code"`
	MenuName           string `db:"menu_name" json:"menu_name"`
	MenuPage           string `db:"menu_page" json:"menu_page"`
	MenuUrl            string `db:"menu_url" json:"menu_url"`
	MenuDesc           string `db:"menu_desc" json:"menu_desc"`
	EndpointVersion    string `db:"endpoint_version" json:"endpoint_version"`
	PrivilegesKey      string `db:"privileges_key" json:"privileges_key"`
}

type scEndpointDetail struct {
	EndpointKey         string  `db:"endpoint_key" json:"endpoint_key"`
	EndpointCategoryKey string  `db:"endpoint_category_key" json:"endpoint_category_key"`
	EndPointCode        string  `db:"endpoint_code" json:"endpoint_code"`
	EndpointName        *string `db:"endpoint_name" json:"endpoint_name"`
	MenuKey             string  `db:"menu_key" json:"menu_key"`
	EndpointVerb        string  `db:"endpoint_verb" json:"endpoint_verb"`
	EndPointUrl         string  `db:"endpoint_uri" json:"endpoint_url"`
	EndpointVersion     string  `db:"endpoint_version" json:"endpoint_version"`
	PrivilegesKey       string  `db:"privileges_key" json:"privileges_key"`
}

func GetEndpointscModels() (result []scEndpoint) {
	query := `SELECT a.endpoint_code, a.endpoint_name, a.endpoint_verb, a.endpoint_uri,a.endpoint_version,a.endpoint_version,a.privileges_key, b.endpoint_ctg_code, b.endpoint_ctg_desc, b.endpoint_ctg_purpose, c.menu_code, c.menu_name, c.menu_page, c.menu_url, c.menu_desc FROM sc_endpoint AS a 
	JOIN sc_endpoint_category AS b ON a.endpoint_category_key = b.endpoint_category_key 
	JOIN sc_menu AS c ON a.menu_key = c.menu_key;`
	log.Println("====================>>>", query)
	err := db.Db.Select(&result, query)
	if err != nil {
		log.Println(err.Error())
		// return http.StatusBadGateway, err
	}
	return

}

func GetDetailEndpointModels(EndPointKey string) (result scEndpointDetail) {
	query := `SELECT endpoint_key,endpoint_category_key,endpoint_code,endpoint_name,menu_key,endpoint_verb,endpoint_uri,endpoint_version,privileges_key FROM sc_endpoint WHERE endpoint_key = ` + EndPointKey
	log.Println("====================>>>", query)
	err := db.Db.Get(&result, query)
	if err != nil {
		log.Println(err.Error())
		// return http.StatusBadGateway, err
	}
	return
}
