package models

import (
	"database/sql"
	"mf-bo-api/db"
	"net/http"
	"strings"

	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

type UdfValue struct {
	UdfValueKey uint64 `db:"udf_value_key"         json:"udf_value_key"`
	UdfInfoKey  uint64 `db:"udf_info_key"          json:"udf_info_key"`
	// UdfInfoKey1 uint64  `db:"udf_info_key1"          json:"udf_info_key1"`
	RowDataKey uint64  `db:"row_data_key"          json:"row_data_key"`
	UdfValues  *string `db:"udf_values"             json:"udf_values"`
}

func GetUdfValueIn(c *[]UdfValue, value []string, field string) (int, error) {
	inQuery := strings.Join(value, ",")
	query2 := `SELECT
				udf_value.* FROM 
				udf_value `
	query := query2 + " WHERE udf_value." + field + " IN(" + inQuery + ")"

	// Main query
	log.Println(query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func CreateMultipleUdfValue(params []interface{}) (int, error) {

	q := `INSERT INTO udf_value ( 
			udf_info_key,
			row_data_key,
			udf_values) VALUES `

	for i := 0; i < len(params); i++ {
		q += "(?)"
		if i < (len(params) - 1) {
			q += ","
		}
	}
	log.Info(q)
	query, args, err := sqlx.In(q, params...)
	if err != nil {
		return http.StatusBadGateway, err
	}

	query = db.Db.Rebind(query)
	_, err = db.Db.Query(query, args...)
	if err != nil {
		log.Error(err.Error())
		return http.StatusBadGateway, err
	}
	return http.StatusOK, nil
}

func DeleteUdfValue(field string, value string, valueIn []string) (int, error) {
	inQuery := strings.Join(valueIn, ",")
	query := `DELETE FROM mam_core.udf_value where ` + field + ` = "` + value + `" 
	AND udf_info_key IN (` + inQuery + `)`

	log.Info(query)

	tx, err := db.Db.Begin()
	if err != nil {
		log.Error(err)
		return http.StatusBadGateway, err
	}
	var ret sql.Result
	ret, err = tx.Exec(query)
	row, _ := ret.RowsAffected()
	if row > 0 {
		tx.Commit()
	} else {
		return http.StatusNotFound, err
	}
	if err != nil {
		log.Error(err)
		return http.StatusBadRequest, err
	}
	return http.StatusOK, nil
}

func GetAllUdfValue(c *[]UdfValue, params map[string]string) (int, error) {
	query := `SELECT
				u.*
			  FROM udf_value as u 
			  INNER JOIN udf_info as ui ON ui.udf_info_key = u.udf_info_key AND ui.rec_status = 1`
	var present bool
	var whereClause []string
	var condition string

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, field+" = '"+value+"'")
		}
	}

	// Combile where clause
	if len(whereClause) > 0 {
		condition += " WHERE "
		for index, where := range whereClause {
			condition += where
			if (len(whereClause) - 1) > index {
				condition += " AND "
			}
		}
	}

	// Check order by
	var orderBy string
	var orderType string
	if orderBy, present = params["orderBy"]; present == true {
		condition += " ORDER BY " + orderBy
		if orderType, present = params["orderType"]; present == true {
			condition += " " + orderType
		}
	}
	query += condition

	// Main query
	// log.Infoln("=============================QUERY GET UDF VALUE=================================")
	log.Infoln(query)
	// log.Infoln("==============================================================")

	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func CreateUdfValue(params map[string]string) (int, error) {
	query := "INSERT INTO udf_value"
	// Get params
	var fields, values string
	var bindvars []interface{}
	for key, value := range params {
		fields += key + ", "
		values += "?, "
		if value == "NULL" {
			var s *string
			bindvars = append(bindvars, s)
		} else {
			bindvars = append(bindvars, value)
		}

	}
	fields = fields[:(len(fields) - 2)]
	values = values[:(len(values) - 2)]

	// Combine params to build query
	query += "(" + fields + ") VALUES(" + values + ")"
	log.Info(query)

	// log.Println("================================ QUERYNYA ADALAH ===============================")
	// log.Println(query)
	// log.Info(query)
	// log.Println("================================================================================")

	tx, err := db.Db.Begin()
	if err != nil {
		log.Error(err)
		return http.StatusBadGateway, err
	}

	_, err = tx.Exec(query, bindvars...)
	tx.Commit()
	if err != nil {
		log.Error(err)
		return http.StatusBadRequest, err
	}
	return http.StatusOK, nil
}

func UpdateDeleteUdfValue(params map[string]string, valueIn []string, fieldIn string, rowDataKey string) (int, error) {
	query := "UPDATE udf_value SET "
	// Get params
	i := 0
	for key, value := range params {
		query += key + " = '" + value + "'"

		if (len(params) - 1) > i {
			query += ", "
		}
		i++
	}

	inQuery := strings.Join(valueIn, ",")
	query += " WHERE udf_info_key IN(" + inQuery + ")"
	query += " AND row_data_key = '" + rowDataKey + "'"

	log.Info(query)

	tx, err := db.Db.Begin()
	if err != nil {
		log.Error(err)
		return http.StatusBadGateway, err
	}
	var ret sql.Result
	ret, err = tx.Exec(query)
	row, _ := ret.RowsAffected()
	tx.Commit()
	if row > 0 {
	} else {
		return http.StatusNotFound, err
	}
	if err != nil {
		log.Error(err)
		return http.StatusBadRequest, err
	}
	return http.StatusOK, nil
}
