package models

import (
	"log"
	"mf-bo-api/db"
	"net/http"
	"strings"
)

type MsSecurities struct {
	SecKey             uint64  `db:"sec_key"               json:"sec_key"`
	SecCode            string  `db:"sec_code"              json:"sec_code"`
	SecName            string  `db:"sec_name"              json:"sec_name"`
	SecuritiesCategory uint64  `db:"securities_category"   json:"securities_category"`
	SecurityType       uint64  `db:"security_type"         json:"security_type"`
	DateIssued         *string `db:"date_issued"           json:"date_issued"`
	DateMatured        *string `db:"date_matured"          json:"date_matured"`
	CurrencyKey        *uint64 `db:"currency_key"          json:"currency_key"`
	SecurityStatus     uint64  `db:"security_status"       json:"security_status"`
	IsinCode           *string `db:"isin_code"             json:"isin_code"`
	SecClassification  uint64  `db:"sec_classification"    json:"sec_classification"`
	RecOrder           *uint64 `db:"rec_order"             json:"rec_order"`
	RecStatus          uint8   `db:"rec_status"            json:"rec_status"`
	RecCreatedDate     *string `db:"rec_created_date"      json:"rec_created_date"`
	RecCreatedBy       *string `db:"rec_created_by"        json:"rec_created_by"`
	RecModifiedDate    *string `db:"rec_modified_date"     json:"rec_modified_date"`
	RecModifiedBy      *string `db:"rec_modified_by"       json:"rec_modified_by"`
	RecImage1          *string `db:"rec_image1"            json:"rec_image1"`
	RecImage2          *string `db:"rec_image2"            json:"rec_image2"`
	RecApprovalStatus  *uint8  `db:"rec_approval_status"   json:"rec_approval_status"`
	RecApprovalStage   *uint64 `db:"rec_approval_stage"    json:"rec_approval_stage"`
	RecApprovedDate    *string `db:"rec_approved_date"     json:"rec_approved_date"`
	RecApprovedBy      *string `db:"rec_approved_by"       json:"rec_approved_by"`
	RecDeletedDate     *string `db:"rec_deleted_date"      json:"rec_deleted_date"`
	RecDeletedBy       *string `db:"rec_deleted_by"        json:"rec_deleted_by"`
	RecAttributeID1    *string `db:"rec_attribute_id1"     json:"rec_attribute_id1"`
	RecAttributeID2    *string `db:"rec_attribute_id2"     json:"rec_attribute_id2"`
	RecAttributeID3    *string `db:"rec_attribute_id3"     json:"rec_attribute_id3"`
}
type Securities struct {
	SecKey                 uint64  `db:"sec_key"               json:"sec_key"`
	SecCode                string  `db:"sec_code"              json:"sec_code"`
	SecName                string  `db:"sec_name"              json:"sec_name"`
	SecuritiesCategory     uint64  `db:"securities_category"   json:"securities_category"`
	SecuritiesCategoryName string  `db:"securities_category_name" json:"securities_category_name"`
	SecurityType           uint64  `db:"security_type"         json:"security_type"`
	SecurityTypeName       string  `db:"security_type_name"         json:"security_type_name"`
	DateIssued             *string `db:"date_issued"           json:"date_issued"`
	DateMatured            *string `db:"date_matured"          json:"date_matured"`
	CurrencyKey            *uint64 `db:"currency_key"          json:"currency_key"`
	CurrencyCode           *string `db:"currency_code"          json:"currency_code"`
	CurrencyName           *string `db:"currency_name"          json:"currency_name"`
	SecurityStatus         *uint64 `db:"security_status"       json:"security_status"`
	SecurityStatusName     *string `db:"security_status_name"       json:"security_status_name"`
	IsinCode               *string `db:"isin_code"             json:"isin_code"`
	SecClassification      *uint64 `db:"sec_classification"    json:"sec_classification"`
}

type SecuritiesDetail struct {
	SecKey                 string  `db:"sec_key"              json:"sec_key"`
	SecCode                string  `db:"sec_code"              json:"sec_code"`
	SecName                string  `db:"sec_name"              json:"sec_name"`
	SecuritiesCategory     uint64  `db:"securities_category"   json:"securities_category"`
	SecuritiesCategoryName string  `db:"securities_category_name" json:"securities_category_name"`
	SecurityType           uint64  `db:"security_type"         json:"security_type"`
	SecurityTypeName       string  `db:"security_type_name"         json:"security_type_name"`
	DateIssued             *string `db:"date_issued"           json:"date_issued"`
	DateMatured            *string `db:"date_matured"          json:"date_matured"`
	CurrencyKey            *uint64 `db:"currency_key"          json:"currency_key"`
	CurrencyCode           *string `db:"currency_code"          json:"currency_code"`
	CurrencyName           *string `db:"currency_name"          json:"currency_name"`
	SecurityStatus         *uint64 `db:"security_status"       json:"security_status"`
	SecurityStatusName     *string `db:"security_status_name"       json:"security_status_name"`
	IsinCode               *string `db:"isin_code"             json:"isin_code"`
	SecClassification      *uint64 `db:"sec_classification"    json:"sec_classification"`
}

func CreateMsSecurities(params map[string]string) (int, error) {
	query := "INSERT INTO ms_securities"
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
func GetSecuritiesModels(c *[]Securities) (int, error) {
	query := `SELECT a.sec_key,
	a.sec_code, 
	a.sec_name,
	a.securities_category,
	b.lkp_name securities_category_name,
	a.security_type,
	c.lkp_name security_type_name,
	a.date_issued,
	a.date_matured,
	a.currency_key,
	e.code currency_code,
	e.name currency_name,
	a.security_status,
	d.lkp_name security_status_name,
	a.isin_code,
	a.sec_classification 
	FROM ms_securities a 
	JOIN gen_lookup b ON a.securities_category = b.lookup_key
	JOIN gen_lookup c ON a.security_type = c.lookup_key
	JOIN gen_lookup d ON a.security_status = d.lookup_key
	JOIN ms_currency e ON a.currency_key = e.currency_key
	WHERE a.rec_status =1`

	log.Println("====================>>>", query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err.Error())
		return http.StatusBadGateway, err
	}
	return http.StatusOK, nil
}
func DeleteMsSecurities(SecKey string, params map[string]string) (int, error) {
	query := `UPDATE ms_securities SET `
	var setClauses []string
	var values []interface{}

	for key, value := range params {
		if key != "sec_key" {
			setClauses = append(setClauses, key+" = ?")
			values = append(values, value)
		}
	}
	query += strings.Join(setClauses, ", ")
	query += ` WHERE sec_key = ?`
	values = append(values, SecKey)

	log.Println("========== DeleteMsSecurities ==========>>>", query)

	_, err := db.Db.Exec(query, values...)
	if err != nil {
		log.Println(err.Error())
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}
func GetMsSecuritiesDetailModels(c *SecuritiesDetail, SecKey string) (int, error) {
	query := `SELECT a.sec_key,
	a.sec_code, 
	a.sec_name,
	a.securities_category,
	b.lkp_name securities_category_name,
	a.security_type,
	c.lkp_name security_type_name,
	a.date_issued,
	a.date_matured,
	a.currency_key,
	e.code currency_code,
	e.name currency_name,
	a.security_status,
	d.lkp_name security_status_name,
	a.isin_code,
	a.sec_classification 
	FROM ms_securities a 
	JOIN gen_lookup b ON a.securities_category = b.lookup_key
	JOIN gen_lookup c ON a.security_type = c.lookup_key
	JOIN gen_lookup d ON a.security_status = d.lookup_key
	JOIN ms_currency e ON a.currency_key = e.currency_key
	WHERE a.rec_status = 1 
	AND a.sec_key =` + SecKey

	log.Println("====================>>>", query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err.Error())
		return http.StatusBadGateway, err
	}
	return http.StatusOK, nil
}
func UpdateMsSecurities(SecKey string, params map[string]string) (int, error) {
	query := `UPDATE ms_securities SET `
	var setClauses []string
	var values []interface{}

	for key, value := range params {
		if key != "sec_key" {
			setClauses = append(setClauses, key+" = ?")
			values = append(values, value)
		}
	}
	query += strings.Join(setClauses, ", ")
	query += ` WHERE sec_key = ?`
	values = append(values, SecKey)

	log.Println("========== UpdateMsSecurities ==========>>>", query)

	_, err := db.Db.Exec(query, values...)
	if err != nil {
		log.Println(err.Error())
		return http.StatusBadGateway, err
	}
	return http.StatusOK, nil
}
