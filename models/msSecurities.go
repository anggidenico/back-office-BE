package models

import (
	"database/sql"
	"errors"
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
	SecClassificationName  *string `db:"sec_classification_name"    json:"sec_classification_name"`
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
	SecClassificationName  *string `db:"sec_classification_name"    json:"sec_classification_name"`
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
	a.sec_classification, 
	f.lkp_name sec_classification_name
	FROM ms_securities a 
	JOIN gen_lookup b ON a.securities_category = b.lookup_key
	JOIN gen_lookup c ON a.security_type = c.lookup_key
	left JOIN gen_lookup d ON a.security_status = d.lookup_key
	left JOIN ms_currency e ON a.currency_key = e.currency_key
	left JOIN gen_lookup f ON a.sec_classification = f.lookup_key
	WHERE a.rec_status = 1 ORDER BY a.rec_created_date DESC`

	// log.Println("====================>>>", query)
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

	// log.Println("========== DeleteMsSecurities ==========>>>", query)

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
	a.sec_classification, 
	f.lkp_name sec_classification_name
	FROM ms_securities a 
	JOIN gen_lookup b ON a.securities_category = b.lookup_key
	JOIN gen_lookup c ON a.security_type = c.lookup_key
	left JOIN gen_lookup d ON a.security_status = d.lookup_key
	left JOIN ms_currency e ON a.currency_key = e.currency_key
	left JOIN gen_lookup f ON a.sec_classification = f.lookup_key
	WHERE a.rec_status = 1 
	AND a.sec_key =` + SecKey

	// log.Println("====================>>>", query)
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

	// log.Println("========== UpdateMsSecurities ==========>>>", query)

	_, err := db.Db.Exec(query, values...)
	if err != nil {
		log.Println(err.Error())
		return http.StatusBadGateway, err
	}
	return http.StatusOK, nil
}

func CreateMsSecurities(params map[string]string) (int, error) {
	// Check for duplicate records
	duplicate, _, err := CheckDuplicateSecurities(params["sec_code"], params["sec_name"], params["security_type"])
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Jika duplikasi ditemukan, perbarui data yang sudah ada
	if duplicate {
		return http.StatusBadRequest, errors.New("data duplikat ditemukan")
	}

	// Jika tidak ada duplikasi, buat data baru
	fields := ""
	placeholders := ""
	var bindvars []interface{}

	for key, value := range params {
		fields += key + `, `
		if value == "NULL" {
			placeholders += `NULL, `
		} else {
			placeholders += `?, `
			bindvars = append(bindvars, value)
		}
	}

	fields = fields[:len(fields)-2]
	placeholders = placeholders[:len(placeholders)-2]

	query := "INSERT INTO ms_securities (" + fields + ") VALUES (" + placeholders + ")"

	tx, err := db.Db.Begin()
	if err != nil {
		return http.StatusBadGateway, err
	}

	_, err = tx.Exec(query, bindvars...)
	if err != nil {
		tx.Rollback()
		return http.StatusBadRequest, err
	}

	tx.Commit()

	return http.StatusOK, nil
}

func CheckDuplicateSecurities(SecCode, SecName, SecType string) (bool, string, error) {
	// Query to check for duplicates
	query := "SELECT sec_key FROM ms_securities WHERE sec_code = ? AND sec_name = ? AND security_type = ? LIMIT 1"
	var key string
	err := db.Db.QueryRow(query, SecCode, SecName, SecType).Scan(&key)

	if err != nil {
		if err == sql.ErrNoRows {
			// No duplicate found
			return false, "", nil
		}
		// Other error occurred
		return false, "", err
	}

	// Duplicate found
	return true, key, nil
}
