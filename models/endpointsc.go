package models

import (
	"fmt"
	"log"
	"mf-bo-api/db"
	"net/http"
)

type ScEndpointt struct {
	EndpointCode       string  `db:"endpoint_code" json:"endpoint_code"`
	EndpointName       string  `db:"endpoint_name" json:"endpoint_name"`
	Method             string  `db:"endpoint_verb" json:"method"`
	EndpointUrl        string  `db:"endpoint_uri" json:"url"`
	EndpointCtgCode    string  `db:"endpoint_ctg_code" json:"endpoint_ctg_code"`
	EndpointCtgDesc    string  `db:"endpoint_ctg_desc" json:"endpoint_ctg_desc"`
	EndpointCtgPurpose string  `db:"endpoint_ctg_purpose" json:"endpoint_ctg_purpose"`
	MenuCode           string  `db:"menu_code" json:"menu_code"`
	MenuName           string  `db:"menu_name" json:"menu_name"`
	MenuPage           string  `db:"menu_page" json:"menu_page"`
	MenuUrl            *string `db:"menu_url" json:"menu_url"`
	MenuDesc           string  `db:"menu_desc" json:"menu_desc"`
	EndpointVersion    string  `db:"endpoint_version" json:"endpoint_version"`
	PrivilegesKey      string  `db:"privileges_key" json:"privileges_key"`
}

type ScEndpointDetail struct {
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

func GetEndpointscModels(c *[]ScEndpointt) (int, error) {
	query := `SELECT a.endpoint_code, a.endpoint_name, a.endpoint_verb, a.endpoint_uri,a.endpoint_version,a.endpoint_version,a.privileges_key, b.endpoint_ctg_code, b.endpoint_ctg_desc, b.endpoint_ctg_purpose, c.menu_code, c.menu_name, c.menu_page, c.menu_url, c.menu_desc FROM sc_endpoint AS a 
	JOIN sc_endpoint_category AS b ON a.endpoint_category_key = b.endpoint_category_key 
	JOIN sc_menu AS c ON a.menu_key = c.menu_key WHERE a.rec_status = 1`
	log.Println("====================>>>", query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err.Error())
		return http.StatusBadGateway, err
	}
	return http.StatusOK, nil
}

func GetDetailEndpointModels(c *ScEndpointDetail, EndPointKey string) (int, error) {
	query := `SELECT endpoint_key,endpoint_category_key,endpoint_code,endpoint_name,menu_key,endpoint_verb,endpoint_uri,endpoint_version,privileges_key FROM sc_endpoint WHERE endpoint_key = ` + EndPointKey
	log.Println("====================>>>", query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err.Error())
		return http.StatusBadGateway, err
	}
	return http.StatusOK, nil
}
func CreateEndpointSc(params map[string]string) (int, error) {
	query := "INSERT INTO sc_endpoint"
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

func UpdateEndpointSc(EndpointKey string, params map[string]string) (int, error) {
	query := "UPDATE sc_endpoint SET "
	// Get params
	i := 0
	for key, value := range params {
		if key != "endpoint_key" {

			query += key + " = '" + value + "'"

			if (len(params) - 2) > i {
				query += ", "
			}
			i++
		}
	}
	query += " WHERE endpoint_key = " + params["endpoint_key"]

	log.Println("UpdateEndpointSc:", query)

	resultSQL, err := db.Db.Exec(query)
	if err != nil {
		log.Println(err.Error())
		return http.StatusInternalServerError, err
	}
	rows, _ := resultSQL.RowsAffected()
	if rows < 1 {
		log.Println("nothing rows affected")
		err2 := fmt.Errorf("nothing rows affected")
		return http.StatusNotFound, err2
	}

	return http.StatusOK, nil
}

// func DeleteEndpoint(EndpointKey string, params map[string]string) (int, error) {
// 	query := `UPDATE sc_endpoint SET `
// 	var setClauses []string
// 	var values []interface{}

// 	for key, value := range params {
// 		if key != "endpoint_key" {
// 			setClauses = append(setClauses, key+" = ?")
// 			values = append(values, value)
// 		}
// 	}
// 	query += strings.Join(setClauses, ", ")
// 	query += ` WHERE endpoint_key = ?`
// 	values = append(values, EndpointKey)

// 	log.Println("========== Update ==========>>>", query)

// 	_, err := db.Db.Exec(query, values...)
// 	if err != nil {
// 		log.Println(err.Error())
// 		return http.StatusBadRequest, err
// 	}

// 	return http.StatusOK, nil
// }
