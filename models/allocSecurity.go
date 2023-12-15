package models

import (
	"database/sql"
	"fmt"
	"log"
	"mf-bo-api/db"
	"net/http"
	"strings"

	"github.com/shopspring/decimal"
)

type AllocSecurity struct {
	AllocSecKey   uint64           `db:"alloc_security_key"          json:"alloc_security_key"`
	ProductKey    uint64           `db:"product_key"          json:"product_key"`
	ProductName   *string          `db:"product_name"          json:"product_name"`
	PeriodeName   *string          `db:"periode_name"          json:"periode_name"`
	SecurityKey   uint64           `db:"sec_key"          json:"sec_key"`
	SecurityName  *string          `db:"sec_name" json:"sec_name"`
	SecurityValue *decimal.Decimal `db:"security_value" json:"security_value"`
	RecOrder      *int64           `db:"rec_order" json:"rec_order"`
}
type AllocSecDetail struct {
	AllocSecKey   uint64           `db:"alloc_security_key"          json:"alloc_security_key"`
	ProductKey    uint64           `db:"product_key"          json:"product_key"`
	ProductName   *string          `db:"product_name"          json:"product_name"`
	PeriodeKey    uint64           `db:"periode_key" json:"periode_key"`
	PeriodeName   string           `db:"periode_name"          json:"periode_name"`
	SecurityKey   uint64           `db:"sec_key"          json:"sec_key"`
	SecurityName  *string          `db:"sec_name" json:"sec_name"`
	SecurityValue *decimal.Decimal `db:"security_value" json:"security_value"`
	RecOrder      int64            `db:"rec_order" json:"rec_order"`
}

func GetAllocSecModels(c *[]AllocSecurity) (int, error) {
	query := `SELECT a.alloc_security_key, 
	a.product_key, 
	b.product_name ,
	c.periode_name, 
	a.sec_key, 
	d.sec_name,
	a.security_value 
	FROM ffs_alloc_security a 
	JOIN ms_product b ON a.product_key = b.product_key 
	JOIN ffs_periode c ON a.periode_key = c.periode_key 
	JOIN ms_securities d ON a.sec_key = d.sec_key 
	WHERE a.rec_status = 1`

	log.Println("====================>>>", query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err.Error())
		return http.StatusBadGateway, err
	}
	return http.StatusOK, nil
}
func GetAllocSecDetailModels(c *AllocSecDetail, AllocSecKey string) (int, error) {
	query := `SELECT a.alloc_security_key, 
	a.product_key, 
	b.product_name,
	a.periode_key, 
	c.periode_name, 
	a.sec_key, 
	d.sec_name,
	a.security_value,
	a.rec_order
	FROM ffs_alloc_security a 
	JOIN ms_product b ON a.product_key = b.product_key 
	JOIN ffs_periode c ON a.periode_key = c.periode_key 
	JOIN ms_securities d ON a.sec_key = d.sec_key 
	WHERE a.rec_status = 1 
	AND a.alloc_security_key =` + AllocSecKey

	log.Println("====================>>>", query)
	err := db.Db.Get(c, query)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("AllocSecKey not found")
			return http.StatusBadGateway, err
		}

		log.Println(err.Error())
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}
func CreateAllocSec(params map[string]interface{}) (int, error) {
	query := "INSERT INTO ffs_alloc_security"
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
func UpdateAllocSec(AllocSeKey string, params map[string]interface{}) (int, error) {
	query := `UPDATE ffs_alloc_security SET `
	var setClauses []string
	var values []interface{}

	for key, value := range params {
		if key != "alloc_security_key" {
			setClauses = append(setClauses, key+" = ?")
			values = append(values, value)
		}
	}
	query += strings.Join(setClauses, ", ")
	query += ` WHERE alloc_security_key = ?`
	values = append(values, AllocSeKey)

	log.Println("========== UpdateFFSAllocSecurity ==========>>>", query)

	_, err := db.Db.Exec(query, values...)
	if err != nil {
		log.Println(err.Error())
		return http.StatusBadRequest, err
	}
	return http.StatusOK, nil
}

func DeleteAllocSec(AllocSecKey string, params map[string]string) (int, error) {
	query := `UPDATE ffs_alloc_security SET `
	var setClauses []string
	var values []interface{}

	for key, value := range params {
		if key != "alloc_security_key" {
			setClauses = append(setClauses, key+" = ?")
			values = append(values, value)
		}
	}
	query += strings.Join(setClauses, ", ")
	query += ` WHERE alloc_security_key = ?`
	values = append(values, AllocSecKey)

	log.Println("========== UpdateRiskProfile ==========>>>", query)

	resultSQL, err := db.Db.Exec(query, values...)
	if err != nil {
		log.Println(err.Error())
		return http.StatusBadGateway, err
	}
	rows, _ := resultSQL.RowsAffected()
	if rows < 1 {
		log.Println("nothing rows affected")
		err2 := fmt.Errorf("nothing rows affected")
		return http.StatusNotFound, err2
	}

	return http.StatusOK, nil
}
